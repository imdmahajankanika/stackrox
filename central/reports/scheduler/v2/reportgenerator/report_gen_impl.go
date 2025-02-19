package reportgenerator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	blobDS "github.com/stackrox/rox/central/blob/datastore"
	deploymentDS "github.com/stackrox/rox/central/deployment/datastore"
	notifierDS "github.com/stackrox/rox/central/notifier/datastore"
	reportConfigDS "github.com/stackrox/rox/central/reportconfigurations/datastore"
	"github.com/stackrox/rox/central/reports/common"
	reportSnapshotDS "github.com/stackrox/rox/central/reports/snapshot/datastore"
	collectionDS "github.com/stackrox/rox/central/resourcecollection/datastore"
	watchedImageDS "github.com/stackrox/rox/central/watchedimage/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/branding"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/mathutil"
	"github.com/stackrox/rox/pkg/notifier"
	"github.com/stackrox/rox/pkg/notifiers"
	"github.com/stackrox/rox/pkg/retry"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/resources"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/templates"
	"github.com/stackrox/rox/pkg/timestamp"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	log = logging.LoggerForModule()
)

type reportGeneratorImpl struct {
	reportConfigDatastore   reportConfigDS.DataStore
	reportSnapshotStore     reportSnapshotDS.DataStore
	deploymentDatastore     deploymentDS.DataStore
	watchedImageDatastore   watchedImageDS.DataStore
	collectionDatastore     collectionDS.DataStore
	collectionQueryResolver collectionDS.QueryResolver
	notifierDatastore       notifierDS.DataStore
	notificationProcessor   notifier.Processor
	blobStore               blobDS.Datastore

	Schema *graphql.Schema
}

func (rg *reportGeneratorImpl) ProcessReportRequest(req *ReportRequest) {
	// First do some basic validation checks on the request.
	err := ValidateReportRequest(req)
	if err != nil {
		rg.logAndUpsertError(errors.Wrap(err, "Invalid report request"), req)
		return
	}
	if req.Ctx == nil {
		log.Error("Request does not have valid non-nil context")
		return
	}

	if req.ReportConfig.GetVulnReportFilters().GetSinceLastSentScheduledReport() {
		req.DataStartTime, err = rg.lastSuccessfulScheduledReportTime(req.Ctx, req.ReportConfig)
		if err != nil {
			rg.logAndUpsertError(errors.Wrap(err, "Error finding last successful scheduled report time"), req)
			return
		}
	} else if req.ReportConfig.GetVulnReportFilters().GetSinceStartDate() != nil {
		req.DataStartTime = req.ReportConfig.GetVulnReportFilters().GetSinceStartDate()
	}

	// Change report status to PREPARING
	err = rg.updateReportStatus(req.Ctx, req.ReportSnapshot, storage.ReportStatus_PREPARING)
	if err != nil {
		rg.logAndUpsertError(errors.Wrap(err, "Error changing report status to PREPARING"), req)
		return
	}

	err = rg.generateReportAndNotify(req)
	if err != nil {
		rg.logAndUpsertError(err, req)
		return
	}

	// Change report status to SUCCESS
	req.ReportSnapshot.ReportStatus.CompletedAt = types.TimestampNow()
	err = rg.updateReportStatus(req.Ctx, req.ReportSnapshot, storage.ReportStatus_SUCCESS)
	if err != nil {
		rg.logAndUpsertError(errors.Wrap(err, "Error changing report status to SUCCESS"), req)
		return
	}
}

/* Report generation helper functions */
func (rg *reportGeneratorImpl) generateReportAndNotify(req *ReportRequest) error {
	// Get the results of running the report query
	deployedImgData, watchedImgData, err := rg.getReportData(req.Ctx, req.ReportConfig, req.Collection, req.DataStartTime)
	if err != nil {
		return err
	}

	// Format results into CSV
	zippedCSVData, err := common.Format(deployedImgData, watchedImgData)
	if err != nil {
		return err
	}

	// If it is an empty report, do not send an attachment in the final notification email and the email body
	// will indicate that no vulns were found
	templateStr := vulnReportEmailTemplate
	if zippedCSVData == nil {
		// If it is an empty report, the email body will indicate that no vulns were found
		templateStr = noVulnsFoundEmailTemplate
	}

	messageText, err := formatMessage(req.DataStartTime, templateStr, req.ReportConfig.GetVulnReportFilters().GetImageTypes())
	if err != nil {
		return errors.Wrap(err, "error formatting the report email text")
	}

	switch req.ReportSnapshot.ReportStatus.ReportNotificationMethod {
	case storage.ReportStatus_DOWNLOAD:
		if err = rg.saveReportData(req.Ctx, req.ReportSnapshot.GetReportConfigurationId(),
			req.ReportSnapshot.GetReportId(), zippedCSVData); err != nil {
			return errors.Wrap(err, "error persisting blob")
		}
	case storage.ReportStatus_EMAIL:
		errorList := errorhelpers.NewErrorList("Error sending email notifications: ")
		for _, notifierConfig := range req.ReportConfig.GetNotifiers() {
			nf := rg.notificationProcessor.GetNotifier(req.Ctx, notifierConfig.GetEmailConfig().GetNotifierId())
			reportNotifier, ok := nf.(notifiers.ReportNotifier)
			if !ok {
				errorList.AddError(errors.Errorf("incorrect type of notifier '%s'", notifierConfig.GetEmailConfig().GetNotifierId()))
				continue
			}
			err := rg.retryableSendReportResults(req.Ctx, reportNotifier, notifierConfig.GetEmailConfig().GetMailingLists(),
				zippedCSVData, messageText)
			if err != nil {
				errorList.AddError(errors.Errorf("Error sending email for notifier '%s': %s",
					notifierConfig.GetEmailConfig().GetNotifierId(), err))
			}
		}
		if !errorList.Empty() {
			return errorList.ToError()
		}
	}
	return nil
}

func (rg *reportGeneratorImpl) saveReportData(ctx context.Context, configID, reportID string, data *bytes.Buffer) error {
	if data == nil {
		data = bytes.NewBuffer([]byte{})
	}

	// Augment ctx with write access to Administration for Upsert.
	ctx = sac.WithGlobalAccessScopeChecker(ctx,
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Administration)))

	// Store downloadable report in blob storage
	b := &storage.Blob{
		Name:         common.GetReportBlobPath(configID, reportID),
		LastUpdated:  types.TimestampNow(),
		ModifiedTime: types.TimestampNow(),
		Length:       int64(data.Len()),
	}
	return rg.blobStore.Upsert(ctx, b, data)
}

func (rg *reportGeneratorImpl) getReportData(ctx context.Context, rc *storage.ReportConfiguration, collection *storage.ResourceCollection,
	dataStartTime *types.Timestamp) ([]common.DeployedImagesResult, []common.WatchedImagesResult, error) {
	var deployedImgResults []common.DeployedImagesResult
	var watchedImgResults []common.WatchedImagesResult
	rQuery, err := rg.buildReportQuery(ctx, rc, collection, dataStartTime)
	if err != nil {
		return nil, nil, err
	}

	if filterOnImageType(rc.GetVulnReportFilters().GetImageTypes(), storage.VulnerabilityReportFilters_DEPLOYED) {
		// We first get deploymentIDs using a DeploymentsQuery and then again run graphQL queries with deploymentIDs to get the deployment objects.
		// Why do we not directly create a queryString directly from the collection and pass that to graphQL?
		// The  query language we support for graphQL has some limitations that prevent us from doing that.
		// DeploymentsQuery is of type *v1.Query and can support complex queries like the one below.
		// [(Cluster: c1 AND Namespace: n1 AND Deployment: d1) OR (Cluster: c2 AND Namespace: n2 AND Deployment: d2)]
		// This query is a 'disjunction of conjunctions' where all conjunctions involve same fields.
		// Current query language for graphQL does not have semantics to define such a query. Due to this we need to fetch deploymentIDs first
		// and then pass them to graphQL.
		deploymentIds, err := rg.getDeploymentIDs(ctx, rQuery.DeploymentsQuery)
		if err != nil {
			return nil, nil, err
		}
		result, err := rg.runPaginatedDeploymentsQuery(ctx, rQuery.CveFieldsQuery, deploymentIds)
		if err != nil {
			return nil, nil, err
		}
		result.Deployments = orderByClusterAndNamespace(result.Deployments)
		deployedImgResults = append(deployedImgResults, result)
	}

	if filterOnImageType(rc.GetVulnReportFilters().GetImageTypes(), storage.VulnerabilityReportFilters_WATCHED) {
		watchedImages, err := rg.getWatchedImages(ctx)
		if err != nil {
			return nil, nil, err
		}
		result, err := rg.runPaginatedImagesQuery(ctx, rQuery.CveFieldsQuery, watchedImages)
		if err != nil {
			return nil, nil, err
		}
		watchedImgResults = append(watchedImgResults, result)
	}

	return deployedImgResults, watchedImgResults, nil
}

func (rg *reportGeneratorImpl) buildReportQuery(ctx context.Context, rc *storage.ReportConfiguration,
	collection *storage.ResourceCollection, dataStartTime *types.Timestamp) (*common.ReportQuery, error) {
	qb := common.NewVulnReportQueryBuilder(collection, rc.GetVulnReportFilters(), rg.collectionQueryResolver,
		timestamp.FromProtobuf(dataStartTime).GoTime())
	rQuery, err := qb.BuildQuery(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error building report query")
	}
	return rQuery, nil
}

// Returns vuln report data from deployments matched by the collection.
func (rg *reportGeneratorImpl) runPaginatedDeploymentsQuery(ctx context.Context, cveQuery string, deploymentIds []string) (common.DeployedImagesResult, error) {
	offset := paginatedQueryStartOffset
	var resultData common.DeployedImagesResult
	for {
		if offset >= len(deploymentIds) {
			break
		}
		scopeQuery := fmt.Sprintf("%s:%s", search.DeploymentID.String(),
			strings.Join(deploymentIds[offset:mathutil.MinInt(offset+paginationLimit, len(deploymentIds))], ","))
		r, err := execQuery[common.DeployedImagesResult](ctx, rg, deployedImagesReportQuery, deployedImagesReportQueryOpName,
			scopeQuery, cveQuery, nil)
		if err != nil {
			return r, err
		}
		resultData.Deployments = append(resultData.Deployments, r.Deployments...)
		offset += paginationLimit
	}
	return resultData, nil
}

// Returns vuln report data for watched images
func (rg *reportGeneratorImpl) runPaginatedImagesQuery(ctx context.Context, cveQuery string, watchedImages []string) (common.WatchedImagesResult, error) {
	offset := paginatedQueryStartOffset
	var resultData common.WatchedImagesResult
	for {
		if offset >= len(watchedImages) {
			break
		}
		scopeQuery := fmt.Sprintf("%s:%s", search.ImageName.String(),
			strings.Join(watchedImages[offset:mathutil.MinInt(offset+paginationLimit, len(watchedImages))], ","))
		sortOpt := map[string]interface{}{
			"field": search.ImageName.String(),
			"aggregateBy": map[string]interface{}{
				"aggregateFunc": "",
				"distinct":      true,
			},
		}
		r, err := execQuery[common.WatchedImagesResult](ctx, rg, watchedImagesReportQuery, watchedImagesReportQueryOpName,
			scopeQuery, cveQuery, sortOpt)
		if err != nil {
			return r, err
		}
		resultData.Images = append(resultData.Images, r.Images...)
		offset += paginationLimit
	}
	return resultData, nil
}

func execQuery[T any](ctx context.Context, rg *reportGeneratorImpl, gqlQuery, opName, scopeQuery, cveQuery string,
	sortOpt map[string]interface{}) (T, error) {
	pagination := map[string]interface{}{
		"offset": paginatedQueryStartOffset,
		"limit":  paginationLimit,
	}
	if sortOpt != nil {
		pagination["sortOptions"] = []interface{}{
			sortOpt,
		}
	}

	response := rg.Schema.Exec(ctx,
		gqlQuery, opName, map[string]interface{}{
			"scopequery": scopeQuery,
			"cvequery":   cveQuery,
			"pagination": pagination,
		})
	if len(response.Errors) > 0 {
		log.Errorf("error running graphql query: %s", response.Errors[0].Message)
		return getZero[T](), response.Errors[0].Err
	}
	var res T
	if err := json.Unmarshal(response.Data, &res); err != nil {
		return getZero[T](), err
	}
	return res, nil
}

/* Utility Functions */

func (rg *reportGeneratorImpl) retryableSendReportResults(ctx context.Context, reportNotifier notifiers.ReportNotifier, mailingList []string,
	zippedCSVData *bytes.Buffer, messageText string) error {
	return retry.WithRetry(func() error {
		return reportNotifier.ReportNotify(ctx, zippedCSVData, mailingList, messageText)
	},
		retry.OnlyRetryableErrors(),
		retry.Tries(3),
		retry.BetweenAttempts(func(previousAttempt int) {
			wait := time.Duration(previousAttempt * previousAttempt * 100)
			time.Sleep(wait * time.Millisecond)
		}),
	)
}

func (rg *reportGeneratorImpl) lastSuccessfulScheduledReportTime(ctx context.Context, config *storage.ReportConfiguration) (*types.Timestamp, error) {
	query := search.NewQueryBuilder().
		AddExactMatches(search.ReportConfigID, config.GetId()).
		AddExactMatches(search.ReportRequestType, storage.ReportStatus_SCHEDULED.String()).
		AddExactMatches(search.ReportState, storage.ReportStatus_SUCCESS.String()).
		WithPagination(search.NewPagination().
			AddSortOption(search.NewSortOption(search.ReportCompletionTime).Reversed(true)).
			Limit(1)).
		ProtoQuery()
	results, err := rg.reportSnapshotStore.SearchReportSnapshots(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "Error finding last successful scheduled report time")
	}
	if len(results) > 1 {
		return nil, errors.Errorf("Received %d records when only one record is expected", len(results))
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].GetReportStatus().GetCompletedAt(), nil
}

func (rg *reportGeneratorImpl) getDeploymentIDs(ctx context.Context, deploymentsQuery *v1.Query) ([]string, error) {
	results, err := rg.deploymentDatastore.Search(ctx, deploymentsQuery)
	if err != nil {
		return nil, err
	}
	return search.ResultsToIDs(results), nil
}

func (rg *reportGeneratorImpl) getWatchedImages(ctx context.Context) ([]string, error) {
	watched, err := rg.watchedImageDatastore.GetAllWatchedImages(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]string, 0, len(watched))
	for _, img := range watched {
		results = append(results, img.GetName())
	}
	return results, nil
}

func (rg *reportGeneratorImpl) updateReportStatus(ctx context.Context, snapshot *storage.ReportSnapshot, status storage.ReportStatus_RunState) error {
	snapshot.ReportStatus.RunState = status
	return rg.reportSnapshotStore.UpdateReportSnapshot(ctx, snapshot)
}

func (rg *reportGeneratorImpl) logAndUpsertError(reportErr error, req *ReportRequest) {
	if req == nil || req.ReportConfig == nil {
		utils.Should(errors.New("Request does not have non-nil report configuration"))
		return
	}
	if req.Ctx == nil {
		utils.Should(errors.New("Request does not have valid non-nil context"))
		return
	}
	if req.ReportSnapshot == nil || req.ReportSnapshot.ReportStatus == nil {
		utils.Should(errors.New("Request does not have non-nil report snapshot with a non-nil report status"))
		return
	}
	if reportErr != nil {
		log.Errorf("Error while running report for config '%s': %s", req.ReportConfig.GetName(), reportErr)
		req.ReportSnapshot.ReportStatus.ErrorMsg = reportErr.Error()
	}
	req.ReportSnapshot.ReportStatus.CompletedAt = types.TimestampNow()
	err := rg.updateReportStatus(req.Ctx, req.ReportSnapshot, storage.ReportStatus_FAILURE)
	if err != nil {
		log.Errorf("Error changing report status to FAILURE for report config '%s', report ID '%s': %s",
			req.ReportConfig.GetName(), req.ReportSnapshot.GetReportId(), err)
	}
}

func filterOnImageType(imageTypes []storage.VulnerabilityReportFilters_ImageType,
	target storage.VulnerabilityReportFilters_ImageType) bool {
	for _, typ := range imageTypes {
		if typ == target {
			return true
		}
	}
	return false
}

func formatMessage(dataStartTime *types.Timestamp, emailTemplate string,
	imageTypes []storage.VulnerabilityReportFilters_ImageType) (string, error) {
	var imageTypesStr string
	includeDeployed := filterOnImageType(imageTypes, storage.VulnerabilityReportFilters_DEPLOYED)
	includeWatched := filterOnImageType(imageTypes, storage.VulnerabilityReportFilters_WATCHED)
	if includeDeployed && includeWatched {
		imageTypesStr = "deployed and watched"
	} else if includeDeployed {
		imageTypesStr = "deployed"
	} else {
		imageTypesStr = "watched"
	}

	data := &reportEmailFormat{
		BrandedProductName: branding.GetProductName(),
		WhichVulns:         "for all vulnerabilities",
		DateStr:            time.Now().Format("January 02, 2006"),
		ImageTypes:         imageTypesStr,
	}
	if dataStartTime != nil {
		data.WhichVulns = fmt.Sprintf("for new vulnerabilities since %s",
			timestamp.FromProtobuf(dataStartTime).GoTime().Format("January 02, 2006"))
	}
	tmpl, err := template.New("emailBody").Parse(emailTemplate)
	if err != nil {
		return "", err
	}
	return templates.ExecuteToString(tmpl, data)
}

func orderByClusterAndNamespace(deployments []*common.Deployment) []*common.Deployment {
	sort.SliceStable(deployments, func(i, j int) bool {
		if deployments[i].Cluster.GetName() == deployments[j].Cluster.GetName() {
			return deployments[i].Namespace < deployments[j].Namespace
		}
		return deployments[i].Cluster.GetName() < deployments[j].Cluster.GetName()
	})
	return deployments
}

func getZero[T any]() T {
	var result T
	return result
}
