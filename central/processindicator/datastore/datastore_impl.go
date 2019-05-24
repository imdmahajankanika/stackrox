package datastore

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/processindicator/index"
	"github.com/stackrox/rox/central/processindicator/pruner"
	"github.com/stackrox/rox/central/processindicator/search"
	"github.com/stackrox/rox/central/processindicator/store"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/sac"
	pkgSearch "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
)

var (
	indicatorSAC = sac.ForResource(resources.Indicator)
)

type datastoreImpl struct {
	storage       store.Store
	indexer       index.Indexer
	searcher      search.Searcher
	prunerFactory pruner.Factory

	stopSig, stoppedSig concurrency.Signal
}

func (ds *datastoreImpl) Search(ctx context.Context, q *v1.Query) ([]pkgSearch.Result, error) {
	return ds.searcher.Search(ctx, q)
}

func (ds *datastoreImpl) SearchRawProcessIndicators(ctx context.Context, q *v1.Query) ([]*storage.ProcessIndicator, error) {
	return ds.searcher.SearchRawProcessIndicators(ctx, q)
}

func (ds *datastoreImpl) GetProcessIndicator(ctx context.Context, id string) (*storage.ProcessIndicator, bool, error) {
	indicator, exists, err := ds.storage.GetProcessIndicator(id)
	if err != nil || !exists {
		return nil, false, err
	}

	if ok, err := indicatorSAC.ScopeChecker(ctx, storage.Access_READ_ACCESS).ForNamespaceScopedObject(indicator).Allowed(ctx); err != nil || !ok {
		return nil, false, err
	}

	return indicator, true, nil
}

func (ds *datastoreImpl) GetProcessIndicators(ctx context.Context) ([]*storage.ProcessIndicator, error) {
	if ok, err := indicatorSAC.ReadAllowed(ctx); err != nil {
		return nil, err
	} else if ok {
		return ds.storage.GetProcessIndicators()
	}
	return ds.SearchRawProcessIndicators(ctx, pkgSearch.EmptyQuery())
}

func (ds *datastoreImpl) AddProcessIndicators(ctx context.Context, indicators ...*storage.ProcessIndicator) error {
	if ok, err := indicatorSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	removedIndicators, err := ds.storage.AddProcessIndicators(indicators...)
	if err != nil {
		return err
	}

	// If there are no indicators to remove, short-circuit the rest of the code path.
	if len(removedIndicators) == 0 {
		return ds.indexer.AddProcessIndicators(indicators)
	}

	removedIndicatorsSet := set.NewStringSet(removedIndicators...)

	// We want to filter out indicators in the current batch which were dropped.
	filteredIndicators := indicators[:0]
	for _, indicator := range indicators {
		if removedIndicatorsSet.Contains(indicator.GetId()) {
			removedIndicatorsSet.Remove(indicator.GetId())
			continue
		}
		filteredIndicators = append(filteredIndicators, indicator)
	}

	// This removes indicators that previously existed in the index.
	if removedIndicatorsSet.Cardinality() > 0 {
		if err := ds.indexer.DeleteProcessIndicators(removedIndicatorsSet.AsSlice()...); err != nil {
			return err
		}
	}
	return ds.indexer.AddProcessIndicators(filteredIndicators)
}

func (ds *datastoreImpl) AddProcessIndicator(ctx context.Context, i *storage.ProcessIndicator) error {
	if ok, err := indicatorSAC.ScopeChecker(ctx, storage.Access_READ_WRITE_ACCESS).ForNamespaceScopedObject(i).Allowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}
	removedIndicator, err := ds.storage.AddProcessIndicator(i)
	if err != nil {
		return errors.Wrap(err, "adding indicator to bolt")
	}
	if removedIndicator != "" {
		if err := ds.indexer.DeleteProcessIndicator(removedIndicator); err != nil {
			return errors.Wrap(err, "removing process indicator")
		}
	}
	if err := ds.indexer.AddProcessIndicator(i); err != nil {
		return errors.Wrap(err, "adding indicator to index")
	}
	return nil
}

func (ds *datastoreImpl) removeMatchingIndicators(results []pkgSearch.Result) error {
	idsToDelete := make([]string, 0, len(results))
	for _, r := range results {
		idsToDelete = append(idsToDelete, r.ID)
	}
	return ds.removeIndicators(idsToDelete)
}

func (ds *datastoreImpl) removeIndicators(ids []string) error {
	for _, id := range ids {
		if err := ds.storage.RemoveProcessIndicator(id); err != nil {
			log.Warnf("Failed to remove process indicator %q: %v", id, err)
		}
	}
	return ds.indexer.DeleteProcessIndicators(ids...)
}

func (ds *datastoreImpl) RemoveProcessIndicatorsByDeployment(ctx context.Context, id string) error {
	if ok, err := indicatorSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}
	q := pkgSearch.NewQueryBuilder().AddExactMatches(pkgSearch.DeploymentID, id).ProtoQuery()
	results, err := ds.Search(ctx, q)
	if err != nil {
		return err
	}
	return ds.removeMatchingIndicators(results)
}

func (ds *datastoreImpl) RemoveProcessIndicatorsOfStaleContainers(ctx context.Context, deploymentID string, currentContainerIDs []string) error {
	if ok, err := indicatorSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}
	queries := make([]*v1.Query, 0, len(currentContainerIDs)+1)
	queries = append(queries, pkgSearch.NewQueryBuilder().AddExactMatches(pkgSearch.DeploymentID, deploymentID).ProtoQuery())

	for _, containerID := range currentContainerIDs {
		queries = append(queries, pkgSearch.NewQueryBuilder().AddStrings(pkgSearch.ContainerID, pkgSearch.NegateQueryString(containerID)).ProtoQuery())
	}

	results, err := ds.Search(ctx, pkgSearch.ConjunctionQuery(queries...))
	if err != nil {
		return err
	}
	return ds.removeMatchingIndicators(results)
}

func (ds *datastoreImpl) prunePeriodically() {
	defer ds.stoppedSig.Signal()

	if ds.prunerFactory == nil {
		return
	}

	t := time.NewTicker(ds.prunerFactory.Period())
	defer t.Stop()
	for !ds.stopSig.IsDone() {
		select {
		case <-t.C:
			ds.prune()
		case <-ds.stopSig.Done():
			return
		}
	}
}

func (ds *datastoreImpl) prune() {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Prune, "ProcessIndicator")
	processInfoToArgs, err := ds.storage.GetProcessInfoToArgs()
	if err != nil {
		log.Errorf("Error while pruning processes: couldn't retrieve process info to args: %s", err)
		return
	}

	pruner := ds.prunerFactory.StartPruning()
	defer pruner.Finish()
	for _, args := range processInfoToArgs {
		idsToRemove := pruner.Prune(args)
		if len(idsToRemove) > 0 {
			if err := ds.removeIndicators(idsToRemove); err != nil {
				log.Errorf("Error while pruning processes: %s", err)
			} else {
				incrementPrunedProcessesMetric(len(idsToRemove))
			}
		}
	}
}

func (ds *datastoreImpl) Stop() bool {
	return ds.stopSig.Signal()
}

func (ds *datastoreImpl) Wait(cancelWhen concurrency.Waitable) bool {
	return concurrency.WaitInContext(&ds.stoppedSig, cancelWhen)
}
