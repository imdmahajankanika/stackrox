package datastore

import (
	"fmt"
	"sort"

	"github.com/stackrox/rox/central/alert/index"
	"github.com/stackrox/rox/central/alert/search"
	"github.com/stackrox/rox/central/alert/store"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	searchCommon "github.com/stackrox/rox/pkg/search"
)

// datastoreImpl is a transaction script with methods that provide the domain logic for CRUD uses cases for Alert
// objects.
type datastoreImpl struct {
	storage  store.Store
	indexer  index.Indexer
	searcher search.Searcher
}

func (ds *datastoreImpl) SearchListAlerts(q *v1.Query) ([]*storage.ListAlert, error) {
	return ds.searcher.SearchListAlerts(q)
}

func (ds *datastoreImpl) ListAlerts(request *v1.ListAlertsRequest) ([]*storage.ListAlert, error) {
	var q *v1.Query
	if request.GetQuery() == "" {
		q = searchCommon.EmptyQuery()
	} else {
		var err error
		q, err = searchCommon.ParseRawQuery(request.GetQuery())
		if err != nil {
			return nil, err
		}
	}
	alerts, err := ds.SearchListAlerts(q)
	if err != nil {
		return nil, err
	}

	// Sort by descending timestamp.
	sort.SliceStable(alerts, func(i, j int) bool {
		if sI, sJ := alerts[i].GetTime().GetSeconds(), alerts[j].GetTime().GetSeconds(); sI != sJ {
			return sI > sJ
		}
		return alerts[i].GetTime().GetNanos() > alerts[j].GetTime().GetNanos()
	})
	return alerts, nil
}

// SearchAlerts returns search results for the given request.
func (ds *datastoreImpl) SearchAlerts(q *v1.Query) ([]*v1.SearchResult, error) {
	return ds.searcher.SearchAlerts(q)
}

// SearchRawAlerts returns search results for the given request in the form of a slice of alerts.
func (ds *datastoreImpl) SearchRawAlerts(q *v1.Query) ([]*storage.Alert, error) {
	return ds.searcher.SearchRawAlerts(q)
}

// GetAlertStore returns all the alerts. Mainly used for compliance checks.
func (ds *datastoreImpl) GetAlertStore() ([]*storage.ListAlert, error) {
	return ds.ListAlerts(nil)
}

// GetAlert returns an alert by id.
func (ds *datastoreImpl) GetAlert(id string) (*storage.Alert, bool, error) {
	return ds.storage.GetAlert(id)
}

// CountAlerts returns the number of alerts that are active
func (ds *datastoreImpl) CountAlerts() (int, error) {
	alerts, err := ds.searcher.SearchListAlerts(searchCommon.NewQueryBuilder().AddStrings(searchCommon.ViolationState, storage.ViolationState_ACTIVE.String()).ProtoQuery())
	return len(alerts), err
}

// AddAlert inserts an alert into storage and into the indexer
func (ds *datastoreImpl) AddAlert(alert *storage.Alert) error {
	if err := ds.storage.AddAlert(alert); err != nil {
		return err
	}
	return ds.indexer.AddAlert(alert)
}

// UpdateAlert updates an alert in storage and in the indexer
func (ds *datastoreImpl) UpdateAlert(alert *storage.Alert) error {
	if err := ds.storage.UpdateAlert(alert); err != nil {
		return err
	}
	return ds.indexer.AddAlert(alert)
}

func (ds *datastoreImpl) MarkAlertStale(id string) error {
	alert, exists, err := ds.GetAlert(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("alert with id '%s' does not exist", id)
	}
	alert.State = storage.ViolationState_RESOLVED
	return ds.UpdateAlert(alert)
}
