package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ClusterRegistryMirrorSet --references=storage.Cluster

// PostgresSQL entries for cluster registry mirror sets should be deleted if the associated cluster is
// deleted. The storage.ClusterRegistryMirrorSet proto references storage.Cluster.Id as a foreign key,
// by doing so an `ON DELETE CASCADE` constraint is added to the table. The contraint makes it so that
// when a cluster is deleted from `clusters` the corresponding rows in `cluster_registry_mirror_sets`
// are also deleted.
