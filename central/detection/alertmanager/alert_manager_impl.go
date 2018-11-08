package alertmanager

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	ptypes "github.com/gogo/protobuf/types"
	alertDataStore "github.com/stackrox/rox/central/alert/datastore"
	"github.com/stackrox/rox/central/detection/runtime"
	notifierProcessor "github.com/stackrox/rox/central/notifier/processor"
	"github.com/stackrox/rox/central/searchbasedpolicies/builders"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
)

var (
	logger = logging.LoggerForModule()
)

type alertManagerImpl struct {
	notifier        notifierProcessor.Processor
	alerts          alertDataStore.DataStore
	runtimeDetector runtime.Detector
}

func (d *alertManagerImpl) AlertAndNotify(currentAlerts []*v1.Alert, oldAlertFilters ...AlertFilterOption) error {
	// Merge the old and the new alerts.
	newAlerts, updatedAlerts, staleAlerts, err := d.mergeManyAlerts(currentAlerts, oldAlertFilters...)
	if err != nil {
		return err
	}

	// Mark any old alerts no longer generated as stale, and insert new alerts.
	if err := d.notifyAndUpdateBatch(newAlerts); err != nil {
		return err
	}
	if err := d.updateBatch(updatedAlerts); err != nil {
		return err
	}
	return d.markAlertsStale(staleAlerts)
}

// UpdateBatch updates all of the alerts in the datastore.
func (d *alertManagerImpl) updateBatch(alertsToMark []*v1.Alert) error {
	errList := errorhelpers.NewErrorList("Error updating alerts: ")
	for _, existingAlert := range alertsToMark {
		errList.AddError(d.alerts.UpdateAlert(existingAlert))
	}
	return errList.ToError()
}

// MarkAlertsStale marks all of the input alerts stale in the input datastore.
func (d *alertManagerImpl) markAlertsStale(alertsToMark []*v1.Alert) error {
	errList := errorhelpers.NewErrorList("Error marking alerts as stale: ")
	for _, existingAlert := range alertsToMark {
		errList.AddError(d.alerts.MarkAlertStale(existingAlert.GetId()))
	}
	return errList.ToError()
}

// NotifyAndUpdateBatch runs the notifier on the input alerts then stores them.
func (d *alertManagerImpl) notifyAndUpdateBatch(alertsToMark []*v1.Alert) error {
	for _, existingAlert := range alertsToMark {
		d.notifier.ProcessAlert(existingAlert)
	}
	return d.updateBatch(alertsToMark)
}

// It is the caller's responsibility to not call this with an empty slice,
// else this function will panic.
func lastTimestamp(processes []*v1.ProcessIndicator) *ptypes.Timestamp {
	return processes[len(processes)-1].GetSignal().GetTime()
}

// This depends on the fact that we sort process indicators in increasing order of timestamp.
func getMostRecentProcessTimestampInAlerts(alerts ...*v1.Alert) (ts *ptypes.Timestamp) {
	for _, alert := range alerts {
		for _, violation := range alert.GetViolations() {
			processes := violation.GetProcesses()
			if len(processes) == 0 {
				continue
			}
			lastTimeStampInViolation := lastTimestamp(processes)
			if protoconv.CompareProtoTimestamps(ts, lastTimeStampInViolation) < 0 {
				ts = lastTimeStampInViolation
			}
		}
	}
	return
}

// If a user has resolved some process indicators from a particular alert, we don't want to display them
// if an alert was violated by a new indicator. This function trims old resolved processes from an alert.
// It returns a bool indicating whether the alert contained only resolved processes -- in which case
// we don't want to generate an alert at all.
func (d *alertManagerImpl) trimResolvedProcessesFromRuntimeAlert(alert *v1.Alert) (isFullyResolved bool) {
	if alert.GetLifecycleStage() != v1.LifecycleStage_RUNTIME {
		return false
	}

	q := search.NewQueryBuilder().
		AddStrings(search.ViolationState, v1.ViolationState_RESOLVED.String()).
		AddExactMatches(search.DeploymentID, alert.GetDeployment().GetId()).
		AddExactMatches(search.PolicyID, alert.GetPolicy().GetId()).
		ProtoQuery()

	oldRunTimeAlerts, err := d.alerts.SearchRawAlerts(q)
	// If there's an error, just log it, and assume there was no previously resolved alert.
	if err != nil {
		logger.Errorf("Failed to retrieve resolved runtime alerts corresponding to %+v: %s", alert, err)
		return false
	}
	if len(oldRunTimeAlerts) == 0 {
		return false
	}

	mostRecentResolvedTimestamp := getMostRecentProcessTimestampInAlerts(oldRunTimeAlerts...)

	var newProcessFound bool
	for _, violation := range alert.GetViolations() {
		if len(violation.GetProcesses()) == 0 {
			continue
		}
		filtered := violation.GetProcesses()[:0]
		for _, process := range violation.GetProcesses() {
			if protoconv.CompareProtoTimestamps(process.GetSignal().GetTime(), mostRecentResolvedTimestamp) > 0 {
				newProcessFound = true
				filtered = append(filtered, process)
			}
		}
		violation.Processes = filtered
		builders.UpdateRuntimeAlertViolationMessage(violation)
	}
	isFullyResolved = !newProcessFound
	return
}

// Some processes in the old alert might have been deleted from the process store because of our pruning,
// which means they only exist in the old alert, and will not be in the new generated alert.
// We don't want to lose them, though, so we keep all the processes from the old alert, and add ones from the new, if any.
// Note that the old alert _was_ active which means that all the processes in it are guaranteed to violate the policy.
func mergeProcessesFromOldIntoNew(old, newAlert *v1.Alert) (newAlertHasNewProcesses bool) {
	// There is exactly one sub-object which has processes.
	var oldViolationWithProcesses *v1.Alert_Violation
	for _, violation := range old.GetViolations() {
		if len(violation.GetProcesses()) > 0 {
			oldViolationWithProcesses = violation
			break
		}
	}
	if oldViolationWithProcesses == nil {
		logger.Errorf("UNEXPECTED: found no old violation with processes for runtime alert %s", proto.MarshalTextString(old))
		newAlertHasNewProcesses = true
		return
	}

	for i, newViolation := range newAlert.GetViolations() {
		if len(newViolation.GetProcesses()) == 0 {
			continue
		}
		newProcessesSlice := oldViolationWithProcesses.GetProcesses()
		// De-dupe processes using timestamps.
		timestamp := lastTimestamp(oldViolationWithProcesses.GetProcesses())
		for _, process := range newViolation.GetProcesses() {
			if protoconv.CompareProtoTimestamps(process.GetSignal().GetTime(), timestamp) > 0 {
				newAlertHasNewProcesses = true
				newProcessesSlice = append(newProcessesSlice, process)
			}
		}
		// If there are no new processes, we'll just use the old alert.
		if !newAlertHasNewProcesses {
			return
		}
		newAlert.Violations[i].Processes = newProcessesSlice
		builders.UpdateRuntimeAlertViolationMessage(newAlert.Violations[i])
		return
	}

	logger.Errorf("UNEXPECTED: found no new violation with processes for runtime alert %s", proto.MarshalTextString(newAlert))
	return
}

// MergeAlerts merges two alerts.
func mergeAlerts(old, newAlert *v1.Alert) *v1.Alert {
	if old.GetLifecycleStage() == v1.LifecycleStage_RUNTIME && newAlert.GetLifecycleStage() == v1.LifecycleStage_RUNTIME {
		newAlertHasNewProcesses := mergeProcessesFromOldIntoNew(old, newAlert)
		// This ensures that we don't keep updating the timestamp of an old runtime alert.
		if !newAlertHasNewProcesses {
			return old
		}
	}

	newAlert.Id = old.GetId()
	// Updated deploy-time alerts continue to have the same enforcement action.
	if newAlert.GetLifecycleStage() == v1.LifecycleStage_DEPLOY && old.GetLifecycleStage() == v1.LifecycleStage_DEPLOY {
		newAlert.Enforcement = old.GetEnforcement()
	}
	newAlert.FirstOccurred = old.GetFirstOccurred()
	return newAlert
}

// MergeManyAlerts merges two alerts.
func (d *alertManagerImpl) mergeManyAlerts(presentAlerts []*v1.Alert, oldAlertFilters ...AlertFilterOption) (newAlerts, updatedAlerts, staleAlerts []*v1.Alert, err error) {
	qb := search.NewQueryBuilder().AddStrings(search.ViolationState, v1.ViolationState_ACTIVE.String())
	for _, filter := range oldAlertFilters {
		filter.apply(qb)
	}
	previousAlerts, err := d.alerts.SearchRawAlerts(qb.ProtoQuery())
	if err != nil {
		err = fmt.Errorf("couldn't load previous alerts (query was %s): %s", qb.Query(), err)
		return
	}

	// Merge any alerts that have new and old alerts.
	for _, alert := range presentAlerts {
		// Don't generate a new alert if it was a resolved runtime alerts.
		isFullyResolvedRuntimeAlert := d.trimResolvedProcessesFromRuntimeAlert(alert)
		if isFullyResolvedRuntimeAlert {
			continue
		}
		if matchingOld := findAlert(alert, previousAlerts); matchingOld != nil {
			updatedAlerts = append(updatedAlerts, mergeAlerts(matchingOld, alert))
			continue
		}

		alert.FirstOccurred = ptypes.TimestampNow()
		newAlerts = append(newAlerts, alert)
	}

	// Find any old alerts no longer being produced.
	for _, alert := range previousAlerts {
		if d.shouldMarkAlertStale(alert, presentAlerts, oldAlertFilters...) {
			staleAlerts = append(staleAlerts, alert)
		}
	}
	return
}

func (d *alertManagerImpl) shouldMarkAlertStale(alert *v1.Alert, presentAlerts []*v1.Alert, oldAlertFilters ...AlertFilterOption) bool {
	// If the alert is still being produced, don't mark it stale.
	if matchingNew := findAlert(alert, presentAlerts); matchingNew != nil {
		return false
	}

	// Only runtime alerts should not be marked stale when they are no longer produced.
	// (Deploy time alerts should disappear along with deployments, for example.)
	if alert.GetLifecycleStage() != v1.LifecycleStage_RUNTIME {
		return true
	}

	// We only want to mark runtime alerts as stale if a policy update causes them to no longer be produced.
	// To determine if this is a policy update, we check if there is a filter on policy ids here.
	specifiedPolicyIDs := set.NewStringSet()
	for _, filter := range oldAlertFilters {
		if filterSpecified := filter.specifiedPolicyID(); filterSpecified != "" {
			specifiedPolicyIDs.Add(filterSpecified)
		}
	}
	if specifiedPolicyIDs.Cardinality() == 0 {
		return false
	}

	// Some other policies were updated, we don't want to mark this alert stale in response.
	if !specifiedPolicyIDs.Contains(alert.GetPolicy().GetId()) {
		return false
	}

	// If the deployment is whitelisted for the policy now, we should mark the alert stale, otherwise we will keep it around.
	return d.runtimeDetector.DeploymentWhitelistedForPolicy(alert.GetDeployment().GetId(), alert.GetPolicy().GetId())
}

func findAlert(toFind *v1.Alert, alerts []*v1.Alert) *v1.Alert {
	for _, alert := range alerts {
		if alertsAreForSamePolicyAndDeployment(alert, toFind) {
			return alert
		}
	}
	return nil
}

func alertsAreForSamePolicyAndDeployment(a1, a2 *v1.Alert) bool {
	return a1.GetPolicy().GetId() == a2.GetPolicy().GetId() && a1.GetDeployment().GetId() == a2.GetDeployment().GetId()
}
