package postgres

//go:generate pg-table-bindings-wrapper --type=storage.ClusterRegistryMirrorSet --references=storage.Cluster

// PostgresSQL entries for ClusterRegistryMirrorSet's should be deleted if the cluster
// they belong to is deleted. The storage.ClusterRegistryMirrorSet proto references
// storage.Cluster.Id as a foreign key constraint which sets ON DELETE CASCADE for the keyj.
// The contraint makes it so that when a cluster is deleted from `clusters` the corresponding rows in
// `cluster_registry_mirror_sets` are also deleted.
