package multipliers

import (
	"testing"

	"github.com/stackrox/rox/central/risk/getters"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/assert"
)

func TestViolationsScore(t *testing.T) {
	cases := []struct {
		name     string
		alerts   []*v1.ListAlert
		expected *storage.Risk_Result
	}{
		{
			name:     "No alerts",
			alerts:   nil,
			expected: nil,
		},
		{
			name: "One critical",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 1",
					},
				},
			},
			expected: &storage.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []*storage.Risk_Result_Factor{
					{Message: "Policy 1 (severity: Critical)"},
				},
				Score: 1.96,
			},
		},
		{
			name: "Two critical",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 1",
					},
				},
			},
			expected: &storage.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []*storage.Risk_Result_Factor{
					{Message: "Policy 1 (severity: Critical)"},
				},
				Score: 1.96,
			},
		},
		{
			name: "Mix of severities (1)",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_HIGH_SEVERITY,
						Name:     "Policy 1",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_MEDIUM_SEVERITY,
						Name:     "Policy 2",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_LOW_SEVERITY,
						Name:     "Policy 3",
					},
				},
			},
			expected: &storage.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []*storage.Risk_Result_Factor{
					{Message: "Policy 1 (severity: High)"},
					{Message: "Policy 2 (severity: Medium)"},
					{Message: "Policy 3 (severity: Low)"},
				},
				Score: 1.84,
			},
		},
		{
			name: "Mix of severities (2)",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 1",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_HIGH_SEVERITY,
						Name:     "Policy 2",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_LOW_SEVERITY,
						Name:     "Policy 3",
					},
				},
			},
			expected: &storage.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []*storage.Risk_Result_Factor{
					{Message: "Policy 1 (severity: Critical)"},
					{Message: "Policy 2 (severity: High)"},
					{Message: "Policy 3 (severity: Low)"},
				},
				Score: 2.56,
			},
		},
		{
			name: "Don't include stale alerts",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 3",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_HIGH_SEVERITY,
						Name:     "Policy 2",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_LOW_SEVERITY,
						Name:     "Policy 1",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_CRITICAL_SEVERITY,
						Name:     "Policy Don't Show Me!",
					},
					State: v1.ViolationState_RESOLVED,
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_HIGH_SEVERITY,
						Name:     "Policy Don't Show Me!",
					},
					State: v1.ViolationState_RESOLVED,
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: storage.Severity_LOW_SEVERITY,
						Name:     "Policy Don't Show Me!",
					},
					State: v1.ViolationState_RESOLVED,
				},
			},
			expected: &storage.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []*storage.Risk_Result_Factor{
					{Message: "Policy 3 (severity: Critical)"},
					{Message: "Policy 2 (severity: High)"},
					{Message: "Policy 1 (severity: Low)"},
				},
				Score: 2.56,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mult := NewViolations(&getters.MockAlertsGetter{
				Alerts: c.alerts,
			})
			deployment := getMockDeployment()
			result := mult.Score(deployment)
			assert.Equal(t, c.expected, result)
		})
	}
}
