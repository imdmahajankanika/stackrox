package rbac

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stackrox/rox/sensor/kubernetes/eventpipeline/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	// Roles:
	//  1. role-admin (all verbs on all resources)
	//  2. role-default (get)
	//  3. role-elevated (get, list) in 2 rules
	//  4. role-elevated-2 (get, list) in a single rule
	// Bindings:
	//  1. admin-subject      -> role-admin
	//  2. default-subject    -> role-default
	//  3. elevated-subject   -> role-elevated
	//  4. elevated-subject-2 -> role-elevated-2
	// Cluster Roles:
	//  1. cluster-admin (all verbs on all resources)
	//  2. cluster-elevated (get on all resources)
	//  3. cluster-elevated-2 (deletecollection)
	//  4. cluster-elevated-3 (deletecollection on pod duplicated)
	//  5. cluster-none (invalid verb on all resources in all API groups)
	//  6. cluster-elevated-admin (all verbs on all resources with additional rule)
	// Cluster Bindings:
	//  3. cluster-admin-subject    -> cluster-admin
	//  4. cluster-elevated-subject -> cluster-elevated
	//  5. cluster-elevated-admin   -> cluster-admin-2]
	//  6. cluster-elevated-2       -> cluster-elevated-subject-3
	//  7. cluster-elevated-3       -> cluster-elevated-subject-4
	//  8. cluster-none             -> cluster-none-subject
	roles = []*v1.Role{
		{
			ObjectMeta: meta("role-admin"),
		},
		{
			ObjectMeta: meta("role-default"),
		},
		{
			ObjectMeta: meta("role-admin"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			}},
		},
		{
			ObjectMeta: meta("role-default"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{""},
				Verbs:     []string{"get"},
			}},
		},
		{
			ObjectMeta: meta("role-elevated"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{""},
				Verbs:     []string{"get"},
			}, {
				APIGroups: []string{""},
				Resources: []string{""},
				Verbs:     []string{"list"},
			}},
		},
		{
			ObjectMeta: meta("role-elevated-2"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{""},
				Verbs:     []string{"get", "list"},
			}},
		},
	}
	bindings = []*v1.RoleBinding{
		{
			ObjectMeta: meta("b1"),
			RoleRef:    role("role-admin"),
		},
		{
			ObjectMeta: meta("b2"),
			RoleRef:    role("role-default"),
		},
		{
			ObjectMeta: meta("b1"),
			RoleRef:    role("role-admin"),
			Subjects: []v1.Subject{
				{
					Name:      "admin-subject",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
				{
					Name:      "cluster-namespace-subject",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
			},
		},
		{
			ObjectMeta: meta("b2"),
			RoleRef:    role("role-default"),
			Subjects: []v1.Subject{{
				Name:      "default-subject",
				Kind:      v1.ServiceAccountKind,
				Namespace: "n1",
			}},
		},
		{
			ObjectMeta: meta("b3"),
			RoleRef:    role("role-elevated"),
			Subjects: []v1.Subject{{
				Name:      "elevated-subject",
				Kind:      v1.ServiceAccountKind,
				Namespace: "n1",
			}},
		},
		{
			ObjectMeta: meta("b4"),
			RoleRef:    role("role-elevated-2"),
			Subjects: []v1.Subject{{
				Name:      "elevated-subject-2",
				Kind:      v1.ServiceAccountKind,
				Namespace: "n1",
			}},
		},
	}
	clusterRoles = []*v1.ClusterRole{
		{
			ObjectMeta: meta("cluster-admin"),
		},
		{
			ObjectMeta: meta("cluster-elevated"),
		},
		{
			ObjectMeta: meta("cluster-admin"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			}},
		},
		{
			ObjectMeta: meta("cluster-elevated-2"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"deletecollection"},
			}},
		},
		{
			ObjectMeta: meta("cluster-elevated-3"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"pod"},
				Verbs:     []string{"deletecollection"},
			}, {
				APIGroups: []string{""},
				Resources: []string{"pod"},
				Verbs:     []string{"deletecollection"},
			}},
		},
		{
			ObjectMeta: meta("cluster-none"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"invalidverb"},
			}},
		},
		{
			ObjectMeta: meta("cluster-elevated-admin"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"get"},
			}, {
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			}},
		},
		{
			ObjectMeta: meta("cluster-elevated"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"get"},
			}},
		},
		{
			ObjectMeta: meta("cluster-elevated"),
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"get"},
			}},
		},
	}
	clusterBindings = []*v1.ClusterRoleBinding{
		{
			ObjectMeta: meta("b3"),
			RoleRef:    clusterRole("cluster-admin"),
		},
		{
			ObjectMeta: meta("b4"),
			RoleRef:    clusterRole("cluster-elevated"),
		},
		{
			ObjectMeta: meta("b3"),
			RoleRef:    clusterRole("cluster-admin"),
			Subjects: []v1.Subject{{
				Name: "cluster-admin-subject",
				Kind: v1.ServiceAccountKind,
			}},
		},
		{
			ObjectMeta: meta("b4"),
			RoleRef:    clusterRole("cluster-elevated"),
			Subjects: []v1.Subject{
				{
					Name:      "cluster-elevated-subject",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
				{
					Name:      "cluster-elevated-subject-2",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
				{
					Name:      "cluster-namespace-subject",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
			},
		},
		{
			ObjectMeta: meta("b5"),
			RoleRef:    clusterRole("cluster-elevated-admin"),
			Subjects: []v1.Subject{
				{
					Name:      "cluster-admin-2",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
			},
		},
		{
			ObjectMeta: meta("b6"),
			RoleRef:    clusterRole("cluster-elevated-2"),
			Subjects: []v1.Subject{{
				Name: "cluster-elevated-subject-3",
				Kind: v1.ServiceAccountKind,
			}},
		},
		{
			ObjectMeta: meta("b7"),
			RoleRef:    clusterRole("cluster-elevated-3"),
			Subjects: []v1.Subject{{
				Name: "cluster-elevated-subject-4",
				Kind: v1.ServiceAccountKind,
			}},
		},
		{
			ObjectMeta: meta("b8"),
			RoleRef:    clusterRole("cluster-none"),
			Subjects: []v1.Subject{{
				Name: "cluster-none-subject",
				Kind: v1.ServiceAccountKind,
			}},
		},
		{
			ObjectMeta: meta("b4"),
			RoleRef:    clusterRole("cluster-elevated"),
			Subjects: []v1.Subject{
				{
					Name:      "cluster-elevated-subject",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
				{
					Name:      "cluster-elevated-subject-2",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
				{
					Name:      "cluster-namespace-subject",
					Kind:      v1.ServiceAccountKind,
					Namespace: "n1",
				},
			},
		},
	}
)

func TestStore(t *testing.T) {
	// Namespace: n1
	// Role: r1
	// Bindings:
	//  - b1 -> r1
	//  - b2 -> r1
	// Cluster role: r2
	// Cluster bindings:
	//  - b3 -> r2
	//  - b4 -> r2
	roles := []*v1.Role{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("r1"),
				Name:      "r1",
				Namespace: "n1",
			},
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{""},
				Verbs:     []string{"get"},
			}, {
				APIGroups: []string{""},
				Resources: []string{""},
				Verbs:     []string{"list"},
			}},
		},
	}
	clusterRoles := []*v1.ClusterRole{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("r2"),
				Name:      "r2",
				Namespace: "n1",
			},
		},
	}
	bindings := []*v1.RoleBinding{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("b1"),
				Name:      "b1",
				Namespace: "n1",
			},
			RoleRef: v1.RoleRef{
				Name:     "r1",
				Kind:     "Role",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("b2"),
				Name:      "b2",
				Namespace: "n1",
			},
			RoleRef: v1.RoleRef{
				Name:     "r1",
				Kind:     "Role",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("b5"),
				Name:      "b5",
				Namespace: "n1",
			},
			RoleRef: v1.RoleRef{
				Name:     "r2",
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
	}
	clusterBindings := []*v1.ClusterRoleBinding{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("b3"),
				Name:      "b3",
				Namespace: "n1",
			},
			RoleRef: v1.RoleRef{
				Name:     "r2",
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("b4"),
				Name:      "b4",
				Namespace: "n1",
			},
			RoleRef: v1.RoleRef{
				Name:     "r2",
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
	}

	tested := NewStore().(*storeImpl)
	dispatcher := NewDispatcher(tested)

	// Add a binding with no role, should get a binding update with no role id.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{
				{
					Id:     "b1",
					Action: central.ResourceAction_UPDATE_RESOURCE,
					Resource: &central.SensorEvent_Binding{
						Binding: &storage.K8SRoleBinding{
							Id:        "b1",
							Name:      "b1",
							Namespace: "n1",
							// No role ID since the role does not yet exist.
							CreatedAt: protoconv.ConvertTimeToTimestamp(bindings[0].GetCreationTimestamp().Time),
							Subjects:  []*storage.Subject{},
						},
					},
				}},
		},
		dispatcher.ProcessEvent(bindings[0], nil, central.ResourceAction_UPDATE_RESOURCE))

	// Upsert the role for the previous binding. The next binding update will get its ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{
				{
					Id:     "r1",
					Action: central.ResourceAction_CREATE_RESOURCE,
					Resource: &central.SensorEvent_Role{
						Role: &storage.K8SRole{
							Id:        "r1",
							Name:      "r1",
							Namespace: "n1",
							CreatedAt: protoconv.ConvertTimeToTimestamp(roles[0].GetCreationTimestamp().Time),
							Rules: []*storage.PolicyRule{{
								ApiGroups: []string{""},
								Resources: []string{""},
								Verbs:     []string{"get"},
							}, {
								ApiGroups: []string{""},
								Resources: []string{""},
								Verbs:     []string{"list"},
							}},
						},
					},
				}},
		},
		dispatcher.ProcessEvent(roles[0], nil, central.ResourceAction_CREATE_RESOURCE))

	// Add another binding for the first role. The binding update should contain the role ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "b2",
				Action: central.ResourceAction_UPDATE_RESOURCE,
				Resource: &central.SensorEvent_Binding{
					Binding: &storage.K8SRoleBinding{
						Id:        "b2",
						Name:      "b2",
						Namespace: "n1",
						RoleId:    "r1", // Note that the role ID is now filled in.
						CreatedAt: protoconv.ConvertTimeToTimestamp(bindings[1].GetCreationTimestamp().Time),
						Subjects:  []*storage.Subject{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(bindings[1], nil, central.ResourceAction_UPDATE_RESOURCE))

	// Add binding for the second role. The binding update should NOT contain the role ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "b5",
				Action: central.ResourceAction_UPDATE_RESOURCE,
				Resource: &central.SensorEvent_Binding{
					Binding: &storage.K8SRoleBinding{
						Id:          "b5",
						Name:        "b5",
						Namespace:   "n1",
						RoleId:      "",
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(bindings[2].GetCreationTimestamp().Time),
						Subjects:    []*storage.Subject{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(bindings[2], nil, central.ResourceAction_UPDATE_RESOURCE))

	// Add a cluster binding with no role, should get a cluster binding update with no role id.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "b3",
				Action: central.ResourceAction_CREATE_RESOURCE,
				Resource: &central.SensorEvent_Binding{
					Binding: &storage.K8SRoleBinding{
						Id:        "b3",
						Name:      "b3",
						Namespace: "n1",
						// No role ID since the role does not yet exist.
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(clusterBindings[0].GetCreationTimestamp().Time),
						Subjects:    []*storage.Subject{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(clusterBindings[0], nil, central.ResourceAction_CREATE_RESOURCE))

	// Upsert the role for the previous binding. The next binding update will get its ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "r2",
				Action: central.ResourceAction_UPDATE_RESOURCE,
				Resource: &central.SensorEvent_Role{
					Role: &storage.K8SRole{
						Id:          "r2",
						Name:        "r2",
						Namespace:   "n1",
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(clusterRoles[0].GetCreationTimestamp().Time),
						Rules:       []*storage.PolicyRule{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(clusterRoles[0], nil, central.ResourceAction_UPDATE_RESOURCE))

	// Upsert binding for the second role. The binding update should contain the role ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "b5",
				Action: central.ResourceAction_UPDATE_RESOURCE,
				Resource: &central.SensorEvent_Binding{
					Binding: &storage.K8SRoleBinding{
						Id:          "b5",
						Name:        "b5",
						Namespace:   "n1",
						RoleId:      "r2",
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(bindings[2].GetCreationTimestamp().Time),
						Subjects:    []*storage.Subject{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(bindings[2], nil, central.ResourceAction_UPDATE_RESOURCE))

	// Update the cluster binding to add a new Subject, should get a cluster binding update with the new role ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "b3",
				Action: central.ResourceAction_UPDATE_RESOURCE,
				Resource: &central.SensorEvent_Binding{
					Binding: &storage.K8SRoleBinding{
						Id:          "b3",
						Name:        "b3",
						Namespace:   "n1",
						RoleId:      "r2", // Note that the role ID is now filled in.
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(clusterBindings[0].GetCreationTimestamp().Time),
						Subjects:    []*storage.Subject{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(clusterBindings[0], nil, central.ResourceAction_UPDATE_RESOURCE))

	// Remove the role. The role should get removed and the binding should get updated with an empty role id.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "r2",
				Action: central.ResourceAction_REMOVE_RESOURCE,
				Resource: &central.SensorEvent_Role{
					Role: &storage.K8SRole{
						Id:          "r2",
						Name:        "r2",
						Namespace:   "n1",
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(clusterRoles[0].GetCreationTimestamp().Time),
						Rules:       []*storage.PolicyRule{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(clusterRoles[0], nil, central.ResourceAction_REMOVE_RESOURCE))

	// Update the cluster binding to add another Subject, should get a cluster binding update *without* role ID.
	assert.Equal(t,
		&component.ResourceEvent{
			ForwardMessages: []*central.SensorEvent{{
				Id:     "b3",
				Action: central.ResourceAction_UPDATE_RESOURCE,
				Resource: &central.SensorEvent_Binding{
					Binding: &storage.K8SRoleBinding{
						Id:        "b3",
						Name:      "b3",
						Namespace: "n1",
						// Note that the role ID is now absent.
						ClusterRole: true,
						CreatedAt:   protoconv.ConvertTimeToTimestamp(clusterBindings[0].GetCreationTimestamp().Time),
						Subjects:    []*storage.Subject{},
					},
				},
			}},
		},
		dispatcher.ProcessEvent(clusterBindings[0], nil, central.ResourceAction_UPDATE_RESOURCE))
}

type storeObjectCounts struct {
	roles      int
	bindings   int
	namespaces int
}

func (c storeObjectCounts) String() string {
	return fmt.Sprintf("Roles %v Bindings %v Namespaces %v", c.roles, c.bindings, c.namespaces)
}

// Taken from 4 customer debug dump metrics on 2021-09-09.
// The number of namespaces is unknown and entirely made up.
func getCustomerStoreObjectCounts() []storeObjectCounts {
	return []storeObjectCounts{
		{roles: 4_168, bindings: 5_281, namespaces: 50},
		{roles: 1_720, bindings: 10_306, namespaces: 100},
		{roles: 873, bindings: 66_258, namespaces: 500},
		{roles: 1_788, bindings: 351_582, namespaces: 1000},
	}
}

// Generates a store with the provided count of elements.
func generateStore(counts storeObjectCounts) Store {
	store := NewStore()

	for i := 0; i < counts.roles; i++ {
		store.UpsertRole(&v1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("some-role-%d", i),
				Namespace: fmt.Sprintf("some-namespace-%d", i%counts.namespaces),
				UID:       types.UID(uuid.NewV4().String()),
			},
			Rules: []v1.PolicyRule{{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get"},
			}},
		})
	}
	for i := 0; i < counts.bindings; i++ {
		store.UpsertBinding(&v1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("some-binding-%d", i),
				Namespace: fmt.Sprintf("some-namespace-%d", i%counts.namespaces),
				UID:       types.UID(uuid.NewV4().String()),
			},
			RoleRef: v1.RoleRef{
				Name: fmt.Sprintf("some-role-%d", i%counts.roles),
			},
			Subjects: []v1.Subject{{
				Name:      "default-subject",
				Kind:      v1.ServiceAccountKind,
				Namespace: fmt.Sprintf("some-namespace-%d", i%counts.namespaces),
			}},
		})
	}

	return store
}

func BenchmarkRBACStoreUpsertTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateStore(storeObjectCounts{roles: 1000, bindings: 10_000, namespaces: 10})
	}
}

func runRBACBenchmarkGetPermissionLevelForDeployment(b *testing.B, store Store, keepCache bool) {
	for n := 0; n < b.N; n++ {
		store.GetPermissionLevelForDeployment(
			&storage.Deployment{ServiceAccount: "default-subject", Namespace: "namespace0"})
		if !keepCache {
			// Important! We really want to call b.StopTimer() here and b.StartTimer() below the
			// UpsertRole call, but when we do this the Benchmarker hangs (see
			// https://stackoverflow.com/a/37624250 for more information). This means the UpsertRole
			// call will be included in the benchmark time.
			// Create a new role to trigger cache invalidation.
			store.UpsertRole(&v1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("roletoinvalidatecache%s", uuid.NewV4().String()),
					Namespace: "namespaceforcacheinvalidation",
					UID:       types.UID(uuid.NewV4().String()),
				},
				Rules: []v1.PolicyRule{{
					APIGroups: []string{""},
					Resources: []string{"pods"},
					Verbs:     []string{"get"},
				}},
			})
		}
	}
}

func BenchmarkRBACStoreAssignPermissionLevelToDeployment(b *testing.B) {
	for _, keepCache := range []bool{true, false} {
		for _, warmUpCache := range []bool{true, false} {
			for _, counts := range getCustomerStoreObjectCounts() {
				store := generateStore(counts)
				if warmUpCache {
					_ = store.GetPermissionLevelForDeployment(
						&storage.Deployment{ServiceAccount: "default-subject", Namespace: "namespace0"})
				}
				b.Run(fmt.Sprintf("KeepCache %t WarmUpCache %t %+v", keepCache, warmUpCache, counts), func(b *testing.B) {
					// The bucket evaluator is not built yet, we will build it initially
					// and keep using it.
					runRBACBenchmarkGetPermissionLevelForDeployment(b, store, keepCache)
				})
			}
		}
	}
}

func BenchmarkRBACUpsertExistingBinding(b *testing.B) {
	b.StopTimer()
	store := NewStore()
	binding := &v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "role",
			Namespace: "namespace",
			UID:       types.UID(uuid.NewV4().String()),
		},
		RoleRef: v1.RoleRef{
			Name: "role",
		},
	}
	store.UpsertBinding(binding)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		store.UpsertBinding(binding)
	}
}

func TestStoreGetPermissionLevelForDeployment(t *testing.T) {
	testCases := []struct {
		deployment *storage.Deployment
		expected   storage.PermissionLevel
	}{
		{expected: storage.PermissionLevel_ELEVATED_CLUSTER_WIDE, deployment: &storage.Deployment{ServiceAccount: "cluster-elevated-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_ELEVATED_CLUSTER_WIDE, deployment: &storage.Deployment{ServiceAccount: "cluster-elevated-subject-2", Namespace: "n1"}},
		{expected: storage.PermissionLevel_ELEVATED_CLUSTER_WIDE, deployment: &storage.Deployment{ServiceAccount: "cluster-elevated-subject-3"}},
		{expected: storage.PermissionLevel_ELEVATED_CLUSTER_WIDE, deployment: &storage.Deployment{ServiceAccount: "cluster-elevated-subject-4"}},
		{expected: storage.PermissionLevel_CLUSTER_ADMIN, deployment: &storage.Deployment{ServiceAccount: "cluster-admin-2", Namespace: "n1"}},
		{expected: storage.PermissionLevel_ELEVATED_CLUSTER_WIDE, deployment: &storage.Deployment{ServiceAccount: "cluster-namespace-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "cluster-elevated-subject"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "cluster-admin-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_CLUSTER_ADMIN, deployment: &storage.Deployment{ServiceAccount: "cluster-admin-subject"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "cluster-none-subject"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "cluster-none-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_ELEVATED_IN_NAMESPACE, deployment: &storage.Deployment{ServiceAccount: "admin-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_DEFAULT, deployment: &storage.Deployment{ServiceAccount: "default-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_ELEVATED_IN_NAMESPACE, deployment: &storage.Deployment{ServiceAccount: "elevated-subject", Namespace: "n1"}},
		{expected: storage.PermissionLevel_ELEVATED_IN_NAMESPACE, deployment: &storage.Deployment{ServiceAccount: "elevated-subject-2", Namespace: "n1"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "elevated-subject-2", Namespace: "n2"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "default-subject"}},
		{expected: storage.PermissionLevel_NONE, deployment: &storage.Deployment{ServiceAccount: "admin-subject"}},
	}
	store := setupStore(roles, clusterRoles, bindings, clusterBindings)
	storeWithNoRoles := setupStore(roles, clusterRoles, bindings, clusterBindings)
	for _, r := range roles {
		storeWithNoRoles.RemoveRole(r)
	}
	for _, r := range clusterRoles {
		storeWithNoRoles.RemoveClusterRole(r)
	}
	storeWithNoBindings := setupStore(roles, clusterRoles, bindings, clusterBindings)
	for _, b := range bindings {
		storeWithNoBindings.RemoveBinding(b)
	}
	for _, b := range clusterBindings {
		storeWithNoBindings.RemoveClusterBinding(b)
	}
	for _, tc := range testCases {
		tc := tc

		name := fmt.Sprintf("%q in namespace %q should have %q permision level",
			tc.deployment.ServiceAccount, tc.deployment.Namespace, tc.expected)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected.String(), store.GetPermissionLevelForDeployment(tc.deployment).String())
		})

		name = fmt.Sprintf("%q in namespace %q should have NO permisions after removing roles but keeping bindings",
			tc.deployment.ServiceAccount, tc.deployment.Namespace)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, storage.PermissionLevel_NONE.String(), storeWithNoRoles.GetPermissionLevelForDeployment(tc.deployment).String())
		})

		name = fmt.Sprintf("%q in namespace %q should have NO permisions after removing bindings but keeping roles",
			tc.deployment.ServiceAccount, tc.deployment.Namespace)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, storage.PermissionLevel_NONE.String(), storeWithNoBindings.GetPermissionLevelForDeployment(tc.deployment).String())
		})
	}
}

func Test_FindSubjectsFromNamespacedRole(t *testing.T) {
	testCase := map[string]struct {
		queryNamespace, queryRoleName string
		expectedServiceAccounts       []string
	}{
		"role-admin-id with two Service Account subjects": {
			queryNamespace:          "n1",
			queryRoleName:           "role-admin",
			expectedServiceAccounts: []string{"admin-subject", "cluster-namespace-subject"},
		},
		"role-default-id with one Service Account subject": {
			queryNamespace:          "n1",
			queryRoleName:           "role-default",
			expectedServiceAccounts: []string{"default-subject"},
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			store := setupStore(roles, clusterRoles, bindings, clusterBindings)
			subjects := store.FindSubjectForRole(tc.queryNamespace, tc.queryRoleName)
			require.Len(t, subjects, len(tc.expectedServiceAccounts))

			var flatSubjects []string
			for _, subj := range subjects {
				_, sa, err := subj.decode()
				require.NoError(t, err, "should be able to decode namespace#name value")
				flatSubjects = append(flatSubjects, sa)
			}

			sort.Strings(flatSubjects)
			sort.Strings(tc.expectedServiceAccounts)

			assert.Equal(t, tc.expectedServiceAccounts, flatSubjects)
		})

	}
}

func Test_FindSubjectForBinding(t *testing.T) {
	testCase := map[string]struct {
		queryNamespace, queryBinding string
		expectedServiceAccounts      []string
	}{
		"b1-id with two Service Account subjects": {
			queryNamespace:          "n1",
			queryBinding:            "b1-id",
			expectedServiceAccounts: []string{"admin-subject", "cluster-namespace-subject"},
		},
		"b2-id with one Service Account subject": {
			queryNamespace:          "n1",
			queryBinding:            "b2-id",
			expectedServiceAccounts: []string{"default-subject"},
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			store := setupStore(roles, clusterRoles, bindings, clusterBindings)
			subjects := store.FindSubjectForBindingID(tc.queryNamespace, tc.queryBinding)
			assert.Len(t, subjects, len(tc.expectedServiceAccounts))

			var flatSubjects []string
			for _, subj := range subjects {
				_, sa, err := subj.decode()
				require.NoError(t, err, "should be able to decode namespace#name value")
				flatSubjects = append(flatSubjects, sa)
			}

			sort.Strings(flatSubjects)
			sort.Strings(tc.expectedServiceAccounts)

			assert.Equal(t, tc.expectedServiceAccounts, flatSubjects)
		})

	}
}

func role(name string) v1.RoleRef {
	return roleRef(name, "Role")
}

func clusterRole(name string) v1.RoleRef {
	return roleRef(name, "ClusterRole")
}

func roleRef(name, kind string) v1.RoleRef {
	return v1.RoleRef{
		Name: name, Kind: kind, APIGroup: "rbac.authorization.k8s.io",
	}
}

func meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: name, UID: types.UID(name + "-id"), Namespace: "n1",
	}
}

func setupStore(roles []*v1.Role, clusterRoles []*v1.ClusterRole, bindings []*v1.RoleBinding, clusterBindings []*v1.ClusterRoleBinding) Store {
	tested := NewStore()
	for _, r := range roles {
		tested.UpsertRole(r)
	}
	for _, b := range bindings {
		tested.UpsertBinding(b)
	}
	for _, r := range clusterRoles {
		tested.UpsertClusterRole(r)
	}
	for _, b := range clusterBindings {
		tested.UpsertClusterBinding(b)
	}
	return tested
}
