import React from 'react';
import { Breadcrumb, BreadcrumbItem } from '@patternfly/react-core';

import { Cluster } from 'types/cluster.proto';
import useURLSearch from 'hooks/useURLSearch';
// import useFetchClusterNamespaces from 'hooks/useFetchClusterNamespaces';
// import { Namespace } from 'hooks/useFetchNamespaceDeployments';
import useFetchNamespacesForClusterAndPermission from 'hooks/useFetchClusterNamespacesForPermission';
import useFetchNamespaceDeployments from 'hooks/useFetchNamespaceDeployments';
import { AccessLevel } from 'services/RolesService';
import ClusterSelector from './ClusterSelector';
import NamespaceSelector from './NamespaceSelector';
import DeploymentSelector from './DeploymentSelector';

type NetworkBreadcrumbsProps = {
    clusters: Cluster[];
    selectedCluster?: { name?: string; id?: string };
    selectedNamespaces: string[];
    selectedDeployments: string[];
};

function NetworkBreadcrumbs({
    clusters = [],
    selectedCluster = {},
    selectedNamespaces = [],
    selectedDeployments = [],
}: NetworkBreadcrumbsProps) {
    const { searchFilter, setSearchFilter } = useURLSearch();

    const networkGraphResource = 'NetworkGraph';
    const readAccess: AccessLevel = 'READ_ACCESS';
    // const { namespaces } = useFetchClusterNamespaces(selectedCluster?.id);
    const results = useFetchNamespacesForClusterAndPermission(
        networkGraphResource,
        readAccess,
        selectedCluster?.id
    );
    // const { namespaces } = results.namespaces;
    const selectedNamespaceIds = results.namespaces.reduce<string[]>((acc: string[], namespace) => {
        return selectedNamespaces.includes(namespace.metadata.name)
            ? [...acc, namespace.metadata.id]
            : acc;
    }, []);
    const { deploymentsByNamespace } = useFetchNamespaceDeployments(selectedNamespaceIds);

    return (
        <>
            <Breadcrumb>
                <BreadcrumbItem isDropdown>
                    <ClusterSelector
                        clusters={clusters}
                        selectedClusterName={selectedCluster?.name ?? ''}
                        searchFilter={searchFilter}
                        setSearchFilter={setSearchFilter}
                    />
                </BreadcrumbItem>
                <BreadcrumbItem isDropdown>
                    <NamespaceSelector
                        namespaces={results.namespaces}
                        selectedNamespaces={selectedNamespaces}
                        selectedDeployments={selectedDeployments}
                        deploymentsByNamespace={deploymentsByNamespace}
                        searchFilter={searchFilter}
                        setSearchFilter={setSearchFilter}
                    />
                </BreadcrumbItem>
                <BreadcrumbItem isDropdown>
                    <DeploymentSelector
                        deploymentsByNamespace={deploymentsByNamespace}
                        selectedDeployments={selectedDeployments}
                        searchFilter={searchFilter}
                        setSearchFilter={setSearchFilter}
                    />
                </BreadcrumbItem>
            </Breadcrumb>
        </>
    );
}

export default NetworkBreadcrumbs;
