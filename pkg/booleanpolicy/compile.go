package booleanpolicy

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/set"
)

func sectionToQuery(section *storage.PolicySection, stage storage.LifecycleStage) (*query.Query, error) {
	fieldQueries := make([]*query.FieldQuery, 0, len(section.GetPolicyGroups()))
	for _, group := range section.GetPolicyGroups() {
		fqs, err := policyGroupToFieldQueries(group)
		if err != nil {
			return nil, err
		}
		fieldQueries = append(fieldQueries, fqs...)
	}
	contextQueries := constructRemainingContextQueries(fieldQueries, stage)
	fieldQueries = append(fieldQueries, contextQueries...)

	return &query.Query{FieldQueries: fieldQueries}, nil
}

func policyGroupToFieldQueries(group *storage.PolicyGroup) ([]*query.FieldQuery, error) {
	if len(group.GetValues()) == 0 {
		return nil, errors.New("no values")
	}

	metadata := fieldsToQB[group.GetFieldName()]
	if metadata == nil || metadata.qb == nil {
		return nil, errors.Errorf("no QB known for group %q", group.GetFieldName())
	}

	if metadata.negationForbidden && group.GetNegate() {
		return nil, errors.Errorf("invalid group: negation not allowed for field %s", group.GetFieldName())
	}
	if metadata.operatorsForbidden && len(group.GetValues()) != 1 {
		return nil, errors.Errorf("invalid group: operators not allowed for field %s", group.GetFieldName())
	}

	fqs := metadata.qb.FieldQueriesForGroup(group)
	if len(fqs) == 0 {
		return nil, errors.New("invalid group: no queries formed")
	}

	return fqs, nil
}

func matchAllQueryForField(fieldName string) *query.FieldQuery {
	return &query.FieldQuery{
		Field:    fieldName,
		MatchAll: true,
	}
}

func constructRemainingContextQueries(fieldQueries []*query.FieldQuery, stage storage.LifecycleStage) []*query.FieldQuery {
	fieldSet := set.NewStringSet()
	for _, query := range fieldQueries {
		fieldSet.Add(query.Field)
	}
	contextFieldSet := set.NewStringSet()
	for _, field := range fieldSet.AsSlice() {
		if metadata, ok := fieldsToQB[field]; ok {
			if contextFieldsToAdd, ok := metadata.contextFields[stage]; ok {
				for _, contextField := range contextFieldsToAdd.AsSlice() {
					contextFieldSet.Add(contextField)
				}
			}
		}
	}
	fieldsToAdd := contextFieldSet.Difference(fieldSet)
	contextQueries := make([]*query.FieldQuery, 0, len(fieldsToAdd))
	for contextField := range fieldsToAdd {
		contextQueries = append(contextQueries, matchAllQueryForField(contextField))
	}
	return contextQueries
}
