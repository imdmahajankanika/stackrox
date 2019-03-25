package check62

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

const checkID = "PCI_DSS_3_2:6_2"

func init() {
	framework.MustRegisterNewCheck(
		framework.CheckMetadata{
			ID:                 checkID,
			Scope:              framework.ClusterKind,
			DataDependencies:   []string{"Images"},
			InterpretationText: interpretationText,
		},
		clusterIsCompliant)
}
func clusterIsCompliant(ctx framework.ComplianceContext) {
	common.CheckFixedCVES(ctx)
	common.CISBenchmarksSatisfied(ctx)
}
