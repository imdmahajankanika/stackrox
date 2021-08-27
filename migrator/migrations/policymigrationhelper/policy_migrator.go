package policymigrationhelper

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/bolthelpers"
	"github.com/stackrox/rox/migrator/log"
	"github.com/stackrox/rox/pkg/set"
	bolt "go.etcd.io/bbolt"
)

// PolicyDiff is an alternative to PolicyChanges that automatically constructs migrations based on diffs of policies.
type PolicyDiff struct {
	FieldsToCompare []FieldComparator
	PolicyFileName  string
}

// PolicyChanges lists the fields that must match before a policy is updated and what it should be updated to
type PolicyChanges struct {
	FieldsToCompare []FieldComparator
	// ToChange is the set of changes that must be made to this policy
	ToChange PolicyUpdates
}

// PolicyUpdates lists the possible fields of a policy that can be updated. Any nil fields will not be updated
// In order to change an item in an array (e.g. exclusions), remove the existing one and add the updated one back in.
type PolicyUpdates struct {
	// PolicySections is the new policy sections
	PolicySections []*storage.PolicySection
	// ExclusionsToAdd is a list of exclusions to insert (or append) to policy
	ExclusionsToAdd []*storage.Exclusion
	// ExclusionsToRemove is a list of exclusions to remove from policy
	ExclusionsToRemove []*storage.Exclusion
	// Remediation is the new remediation string
	Remediation *string
	// Rationale is the new rationale string
	Rationale *string
	// Description is the new description string
	Description *string
}

func (u *PolicyUpdates) applyToPolicy(policy *storage.Policy) {
	if u == nil {
		return
	}

	if u.ExclusionsToRemove != nil {
		for _, toRemove := range u.ExclusionsToRemove {
			if !removeExclusion(policy, toRemove) {
				log.WriteToStderrf("policy ID %s has already been altered because exclusion was already removed. Will not update.", policy.Id)
				continue
			}
		}
	}

	// Add new exclusions as needed
	if u.ExclusionsToAdd != nil {
		policy.Exclusions = append(policy.Exclusions, u.ExclusionsToAdd...)
	}

	// If policy section is to be updated, just clear the old one for the new
	if u.PolicySections != nil {
		policy.PolicySections = u.PolicySections
	}

	// Update string fields as needed
	if u.Rationale != nil {
		policy.Rationale = *u.Rationale
	}
	if u.Remediation != nil {
		policy.Remediation = *u.Remediation
	}
	if u.Description != nil {
		policy.Description = *u.Description
	}
}

func removeExclusion(policy *storage.Policy, exclusionToRemove *storage.Exclusion) bool {
	exclusions := policy.GetExclusions()
	for i, exclusion := range exclusions {
		if reflect.DeepEqual(exclusion, exclusionToRemove) {
			policy.Exclusions = append(exclusions[:i], exclusions[i+1:]...)
			return true
		}
	}
	return false
}

// FieldComparator should compare policies and return true if they match for a defined field
type FieldComparator func(first, second *storage.Policy) bool

// PolicySectionComparator compares the policySections of both policies and returns true if they are equal
func PolicySectionComparator(first, second *storage.Policy) bool {
	return reflect.DeepEqual(first.GetPolicySections(), second.GetPolicySections())
}

// ExclusionComparator compares the Exclusions of both policies and returns true if they are equal
func ExclusionComparator(first, second *storage.Policy) bool {
	return reflect.DeepEqual(first.GetExclusions(), second.GetExclusions())
}

// RemediationComparator compares the Remediation section of both policies and returns true if they are equal
func RemediationComparator(first, second *storage.Policy) bool {
	return first.GetRemediation() == second.GetRemediation()
}

// RationaleComparator compares the Rationale section of both policies and returns true if they are equal
func RationaleComparator(first, second *storage.Policy) bool {
	return first.GetRationale() == second.GetRationale()
}

// DescriptionComparator compares the Description of both policies and returns true if they are equal
func DescriptionComparator(first, second *storage.Policy) bool {
	return first.GetDescription() == second.GetDescription()
}

var (
	policyBucketName = []byte("policies")
)

func readPolicyFromFile(fs embed.FS, filePath string) (*storage.Policy, error) {
	contents, err := fs.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file %s", filePath)
	}
	var policy storage.Policy
	err = jsonpb.Unmarshal(bytes.NewReader(contents), &policy)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshal policy json at path %s", filePath)
	}
	return &policy, nil
}

func diffPolicies(beforePolicy, afterPolicy *storage.Policy) (PolicyUpdates, error) {
	// Clone policies because we mutate them.
	beforePolicy = beforePolicy.Clone()
	afterPolicy = afterPolicy.Clone()

	var updates PolicyUpdates

	matchedAfterExclusionsIdxs := set.NewIntSet()
	for _, beforeExclusion := range beforePolicy.GetExclusions() {
		var found bool
		for afterExclusionIdx, afterExclusion := range afterPolicy.GetExclusions() {
			if reflect.DeepEqual(beforeExclusion, afterExclusion) {
				found = true
				matchedAfterExclusionsIdxs.Add(afterExclusionIdx)
				break
			}
		}
		if !found {
			updates.ExclusionsToRemove = append(updates.ExclusionsToRemove, beforeExclusion)
		}
	}
	for i, exclusion := range afterPolicy.GetExclusions() {
		if !matchedAfterExclusionsIdxs.Contains(i) {
			updates.ExclusionsToAdd = append(updates.ExclusionsToAdd, exclusion)
		}
	}
	beforePolicy.Exclusions = nil
	afterPolicy.Exclusions = nil
	if !reflect.DeepEqual(beforePolicy, afterPolicy) {
		return PolicyUpdates{}, errors.New("policies have diff after nil-ing out fields we checked, please update this function " +
			"to be able to diff more fields")
	}
	return updates, nil
}

const (
	policyDiffParentDirName = "policies_before_and_after"
	beforeDirName           = policyDiffParentDirName + "/before"
	afterDirName            = policyDiffParentDirName + "/after"
)

// MigratePoliciesWithDiffs migrates policies with the given diffs.
// The policyDiffFS should be an embedded FS that satisfies the following conditions:
// 1. It must contain a top-level directory called "policies_before_and_after".
// 2. That directory must contain two subdirectories: "before" and "after".
// 3. For each policy being migrated, there must be one copy in the "before" directory and one in the "after" directory.
// 4. The file names for a policy should match the PolicyFileName in the corresponding PolicyDiff passed in the third argument.
// This function then automatically computes the diff for each policy, and executes the migration.
func MigratePoliciesWithDiffs(db *bolt.DB, policyDiffFS embed.FS, policyDiffs []PolicyDiff) error {
	policiesToMigrate := make(map[string]PolicyChanges, len(policyDiffs))
	preMigrationPolicies := make(map[string]*storage.Policy, len(policyDiffs))
	for _, diff := range policyDiffs {
		beforePolicy, err := readPolicyFromFile(policyDiffFS, filepath.Join(beforeDirName, diff.PolicyFileName))
		if err != nil {
			return err
		}
		afterPolicy, err := readPolicyFromFile(policyDiffFS, filepath.Join(afterDirName, diff.PolicyFileName))
		if err != nil {
			return err
		}
		if beforePolicy.GetId() == "" || beforePolicy.GetId() != afterPolicy.GetId() {
			return errors.Errorf("policies in file %s don't both have the same, non-empty, id", diff.PolicyFileName)
		}
		updates, err := diffPolicies(beforePolicy, afterPolicy)
		if err != nil {
			return err
		}
		policiesToMigrate[beforePolicy.GetId()] = PolicyChanges{FieldsToCompare: diff.FieldsToCompare, ToChange: updates}
		preMigrationPolicies[beforePolicy.GetId()] = beforePolicy
	}
	return MigratePolicies(db, policiesToMigrate, preMigrationPolicies)
}

// MigratePoliciesWithPreMigrationFS is a variant of MigratePolicies that takes in an embed.FS with the pre migration policies.
// `preMigFS` is expected to have a directory called `preMigDirName`, which has one JSON file per policy.
// Each JSON file is expected to have the filename <policy_id>.json.
func MigratePoliciesWithPreMigrationFS(db *bolt.DB, policiesToMigrate map[string]PolicyChanges, preMigFS embed.FS, preMigDirName string) error {
	comparisonPolicies := make(map[string]*storage.Policy)
	for policyID := range policiesToMigrate {
		path := filepath.Join(preMigDirName, fmt.Sprintf("%s.json", policyID))
		policy, err := readPolicyFromFile(preMigFS, path)
		if err != nil {
			return err
		}
		comparisonPolicies[policyID] = policy
	}
	return MigratePolicies(db, policiesToMigrate, comparisonPolicies)
}

// MigratePolicies will migrate all policies in the db as specified by policiesToMigrate assuming the policies in the db
// matches the policies within comparisonPolicies.
func MigratePolicies(db *bolt.DB, policiesToMigrate map[string]PolicyChanges, comparisonPolicies map[string]*storage.Policy) error {
	if exists, err := bolthelpers.BucketExists(db, policyBucketName); err != nil {
		return errors.Wrapf(err, "getting bucket with name %q", policyBucketName)
	} else if !exists {
		return errors.Errorf("unable to find policy bucket with name %s", policyBucketName)
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(policyBucketName)

		// Migrate and update policies one by one. Abort the transaction, and hence
		// the migration, in case of any error.
		for policyID, updateDetails := range policiesToMigrate {
			v := bucket.Get([]byte(policyID))
			if v == nil {
				log.WriteToStderrf("no policy exists for ID %s in policy migration. Continuing", policyID)
				continue
			}

			var policy storage.Policy
			if err := proto.Unmarshal(v, &policy); err != nil {
				// Unable to recover, so abort transaction
				return errors.Wrapf(err, "unmarshaling migrated policy with id %q", policyID)
			}

			// Fetch the saved policy state to compare with
			comparePolicy, ok := comparisonPolicies[policyID]
			if !ok || comparePolicy == nil {
				return errors.Errorf("policy cannot be compared because comparison policy doesn't exist for %q", policyID)
			}

			// Validate all the required fields to ensure policy hasn't been updated
			if !checkIfPoliciesMatch(updateDetails.FieldsToCompare, comparePolicy, &policy) {
				log.WriteToStderrf("policy ID %s has already been altered. Will not update.", policyID)
				continue
			}

			// Update policy as needed
			updateDetails.ToChange.applyToPolicy(&policy)

			policyBytes, err := proto.Marshal(&policy)
			if err != nil {
				return errors.Wrapf(err, "marshaling migrated policy %q with id %q", policy.GetName(), policy.GetId())
			}
			if err := bucket.Put([]byte(policyID), policyBytes); err != nil {
				return errors.Wrapf(err, "writing migrated policy with id %q to the store", policy.GetId())
			}
		}

		return nil
	})
}

func checkIfPoliciesMatch(fieldsToCompare []FieldComparator, first *storage.Policy, second *storage.Policy) bool {
	for _, field := range fieldsToCompare {
		if !field(first, second) {
			return false
		}
	}
	return true
}
