package singleton

import (
	"time"

	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/license/validator"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	utils.Must(
		validatorInstance.RegisterSigningKey(
			validator.EC256,
			// projects/stackrox-dev/locations/global/keyRings/licensing-demos/cryptoKeys/demo-license-signer/cryptoKeyVersions/1
			[]byte{
				0x30, 0x59, 0x30, 0x13, 0x06, 0x07, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x02,
				0x01, 0x06, 0x08, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07, 0x03,
				0x42, 0x00, 0x04, 0x00, 0xcd, 0x7d, 0xbd, 0x6d, 0xb5, 0xc1, 0x13, 0x4b,
				0x54, 0x7a, 0xe2, 0x80, 0x51, 0x51, 0x80, 0x38, 0x24, 0xd0, 0xd3, 0xb0,
				0x88, 0x02, 0x06, 0xb0, 0x69, 0xbd, 0x1d, 0x1b, 0x94, 0xb0, 0xbb, 0x8b,
				0xac, 0x63, 0x9e, 0xe7, 0x87, 0x20, 0xbb, 0x8f, 0x1b, 0x13, 0x42, 0xdb,
				0x0d, 0x23, 0xb3, 0x00, 0xd0, 0xb0, 0xe4, 0x28, 0xeb, 0x7b, 0x64, 0x1f,
				0x8b, 0x10, 0xb6, 0x3a, 0x89, 0xe7, 0xa8,
			},
			validator.SigningKeyRestrictions{
				EarliestNotValidBefore:        buildinfo.BuildTimestamp(),
				LatestNotValidAfter:           buildinfo.BuildTimestamp().Add(90 * 24 * time.Hour),
				MaxDuration:                   14 * 24 * time.Hour,
				AllowOffline:                  true,
				MaxNodeLimit:                  50,
				AllowNoBuildFlavorRestriction: true,
				DeploymentEnvironments:        []string{"gcp/ultra-current-825", "azure/66c57ff5-f49f-4510-ae04-e26d3ad2ee63", "aws/051999192406"},
			}),
	)
}
