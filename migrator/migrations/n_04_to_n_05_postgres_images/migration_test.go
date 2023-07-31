// Code originally generated by pg-bindings generator.

//go:build sql_integration

package n4ton5

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations/n_04_to_n_05_postgres_images/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_04_to_n_05_postgres_images/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	ctx context.Context

	legacyDB   *rocksdb.RocksDB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), false)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestImageMigration() {
	newStore := pgStore.New(s.postgresDB.DB, true)
	dacky, err := dackbox.NewRocksDBDackBox(s.legacyDB, nil, []byte("graph"), []byte("dirty"), []byte("valid"))
	s.NoError(err)
	legacyStore := legacy.New(dacky, concurrency.NewKeyFence(), false)

	// Prepare data and write to legacy DB
	images := []*storage.Image{
		{
			Id: "sha256:sha1",
			Name: &storage.ImageName{
				FullName: "name1",
			},
			Metadata: &storage.ImageMetadata{
				V1: &storage.V1Metadata{
					Created: types.TimestampNow(),
				},
			},
			Scan: &storage.ImageScan{
				ScanTime:        types.TimestampNow(),
				OperatingSystem: "cloud",
				Components: []*storage.EmbeddedImageScanComponent{
					{
						Name:    "comp1",
						Version: "ver1",
						HasLayerIndex: &storage.EmbeddedImageScanComponent_LayerIndex{
							LayerIndex: 1,
						},
						Vulns: []*storage.EmbeddedVulnerability{},
					},
					{
						Name:    "comp2",
						Version: "ver2",
						HasLayerIndex: &storage.EmbeddedImageScanComponent_LayerIndex{
							LayerIndex: 3,
						},
						Vulns: []*storage.EmbeddedVulnerability{
							{
								Cve:                "cve1",
								VulnerabilityType:  storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
								VulnerabilityTypes: []storage.EmbeddedVulnerability_VulnerabilityType{storage.EmbeddedVulnerability_IMAGE_VULNERABILITY},
							},
							{
								Cve:                "cve2",
								VulnerabilityType:  storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
								VulnerabilityTypes: []storage.EmbeddedVulnerability_VulnerabilityType{storage.EmbeddedVulnerability_IMAGE_VULNERABILITY},
								SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
									FixedBy: "ver3",
								},
							},
						},
					},
					{
						Name:    "comp3",
						Version: "ver1",
						HasLayerIndex: &storage.EmbeddedImageScanComponent_LayerIndex{
							LayerIndex: 2,
						},
						Vulns: []*storage.EmbeddedVulnerability{
							{
								Cve:                "cve1",
								VulnerabilityType:  storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
								VulnerabilityTypes: []storage.EmbeddedVulnerability_VulnerabilityType{storage.EmbeddedVulnerability_IMAGE_VULNERABILITY},
								SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
									FixedBy: "ver2",
								},
							},
							{
								Cve:                "cve2",
								VulnerabilityType:  storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
								VulnerabilityTypes: []storage.EmbeddedVulnerability_VulnerabilityType{storage.EmbeddedVulnerability_IMAGE_VULNERABILITY},
							},
						},
					},
				},
			},
			RiskScore: 30,
		},
		{
			Id: "sha256:sha2",
			Name: &storage.ImageName{
				FullName: "name2",
			},
		},
	}

	for _, image := range images {
		s.NoError(legacyStore.Upsert(s.ctx, image))
	}

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify Count
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(images), count)

	// Verify Image
	for _, image := range images {
		fetchedImage, found, err := newStore.Get(s.ctx, image.GetId())
		s.NoError(err)
		s.True(found)

		fetchedScan := fetchedImage.GetScan()
		scan := image.GetScan()
		fetchedImage.Scan = nil
		image.Scan = nil
		s.Equal(image, fetchedImage)

		s.Len(fetchedScan.GetComponents(), len(scan.GetComponents()))
		for ci, component := range scan.GetComponents() {
			fetchedComponent := fetchedScan.GetComponents()[ci]
			s.Len(fetchedComponent.GetVulns(), len(component.GetVulns()))
			for vi, vuln := range component.GetVulns() {
				fetchedVuln := fetchedComponent.GetVulns()[vi]
				fetchedVuln.FirstImageOccurrence = nil
				fetchedVuln.FirstSystemOccurrence = nil
				s.Equal(vuln, fetchedVuln)
			}
			s.Equal(component, fetchedComponent)
		}
	}
}
