/**
 * Application route paths constants.
 */

import { resourceTypes, standardEntityTypes, rbacConfigTypes } from 'constants/entityTypes';
import { IsFeatureFlagEnabled } from 'hooks/useFeatureFlags';
import { HasReadAccess } from 'hooks/usePermissions';
import { FeatureFlagEnvVar } from 'types/featureFlag';
import { ResourceName } from 'types/roleResources';

export const mainPath = '/main';
export const loginPath = '/login';
export const testLoginResultsPath = '/test-login-results';
export const authResponsePrefix = '/auth/response/';
export const authorizeRoxctlPath = '/authorize-roxctl';

// Add (related) path variables in alphabetical order to minimize merge conflicts when multiple people add routes.
export const accessControlBasePath = `${mainPath}/access-control`;
export const accessControlPath = `${accessControlBasePath}/:entitySegment?/:entityId?`;
export const apidocsPath = `${mainPath}/apidocs`;
export const clustersBasePath = `${mainPath}/clusters`;
export const clustersPathWithParam = `${clustersBasePath}/:clusterId?`;
export const clustersListPath = `${mainPath}/clusters-pf`;
export const clustersDelegateScanningPath = `${clustersBasePath}/delegate-scanning`;
export const collectionsBasePath = `${mainPath}/collections`;
export const collectionsPath = `${mainPath}/collections/:collectionId?`;
export const complianceBasePath = `${mainPath}/compliance`;
export const compliancePath = `${mainPath}/:context(compliance)`;
export const complianceEnhancedBasePath = `${mainPath}/compliance-enhanced`;
export const configManagementPath = `${mainPath}/configmanagement`;
export const dashboardPath = `${mainPath}/dashboard`;
export const dataRetentionPath = `${mainPath}/retention`;
export const integrationsPath = `${mainPath}/integrations`;
export const integrationCreatePath = `${integrationsPath}/:source/:type/create`;
export const integrationDetailsPath = `${integrationsPath}/:source/:type/view/:id`;
export const integrationEditPath = `${integrationsPath}/:source/:type/edit/:id`;
export const integrationsListPath = `${integrationsPath}/:source/:type`;
export const listeningEndpointsBasePath = `${mainPath}/listening-endpoints`;
export const networkBasePath = `${mainPath}/network-graph`;
export const networkPath = `${networkBasePath}/:detailType?/:detailId?`;
export const policyManagementBasePath = `${mainPath}/policy-management`;
export const policiesBasePath = `${policyManagementBasePath}/policies`;
export const policiesPath = `${policiesBasePath}/:policyId?/:command?`;
export const policyCategoriesPath = `${policyManagementBasePath}/policy-categories`;
export const deprecatedPoliciesBasePath = `${mainPath}/policies`;
export const deprecatedPoliciesPath = `${deprecatedPoliciesBasePath}/:policyId?/:command?`;
export const riskBasePath = `${mainPath}/risk`;
export const riskPath = `${riskBasePath}/:deploymentId?`;
export const searchPath = `${mainPath}/search`;
export const secretsPath = `${mainPath}/configmanagement/secrets/:secretId?`;
export const systemConfigPath = `${mainPath}/systemconfig`;
export const systemHealthPath = `${mainPath}/system-health`;
export const userBasePath = `${mainPath}/user`;
export const userRolePath = `${userBasePath}/roles/:roleName`;
export const violationsBasePath = `${mainPath}/violations`;
export const violationsPath = `${violationsBasePath}/:alertId?`;
export const vulnManagementPath = `${mainPath}/vulnerability-management`;
export const vulnManagementReportsPath = `${vulnManagementPath}/reports`;
export const vulnManagementRiskAcceptancePath = `${vulnManagementPath}/risk-acceptance`;
export const vulnerabilitiesBasePath = `${mainPath}/vulnerabilities`;
export const vulnerabilitiesWorkloadCvesPath = `${vulnerabilitiesBasePath}/workload-cves`;
export const vulnerabilityReportsPath = `${vulnerabilitiesBasePath}/reports`;

// Configuration Management paths for links from Search:

export const configManagementRolesPath = `${configManagementPath}/roles`;
export const configManagementSecretsPath = `${configManagementPath}/secrets`;
export const configManagementServiceAccountsPath = `${configManagementPath}/serviceaccounts`;

// Vulnerability Management 1.0 paths for links from Dashboard or Search:

export const vulnManagementImagesPath = `${vulnManagementPath}/images`;
export const vulnManagementNamespacesPath = `${vulnManagementPath}/namespaces`;
export const vulnManagementNodesPath = `${vulnManagementPath}/nodes`;

// Source of truth for conditional rendering of Body route paths and NavigationSidebar links.

type RouteDescription = {
    featureFlagDependency?: FeatureFlagEnvVar[]; // assume multiple feature flags imply all must be enabled
    resourceAccessRequirements: ResourceName[]; // assume READ_ACCESS and multiple resource names imply must have access to all
};

// Add path variables in alphabetical order to minimize merge conflicts when multiple people add routes.
const routeDescriptionMap: Record<string, RouteDescription> = {
    [accessControlPath]: {
        resourceAccessRequirements: [],
    },
    [apidocsPath]: {
        resourceAccessRequirements: [],
    },
    [clustersDelegateScanningPath]: {
        resourceAccessRequirements: [],
    },
    [clustersPathWithParam]: {
        resourceAccessRequirements: [],
    },
    [collectionsPath]: {
        resourceAccessRequirements: ['WorkflowAdministration'],
    },
    [compliancePath]: {
        resourceAccessRequirements: [],
    },
    [complianceEnhancedBasePath]: {
        resourceAccessRequirements: [],
    },
    [configManagementPath]: {
        resourceAccessRequirements: [],
    },
    [dashboardPath]: {
        resourceAccessRequirements: [],
    },
    [integrationsPath]: {
        resourceAccessRequirements: [],
    },
    [listeningEndpointsBasePath]: {
        resourceAccessRequirements: ['Deployment', 'DeploymentExtension'],
    },
    [networkPath]: {
        resourceAccessRequirements: [],
    },
    [policyManagementBasePath]: {
        resourceAccessRequirements: [],
    },
    [riskPath]: {
        resourceAccessRequirements: [],
    },
    [searchPath]: {
        resourceAccessRequirements: [],
    },
    [systemConfigPath]: {
        resourceAccessRequirements: [],
    },
    [systemHealthPath]: {
        resourceAccessRequirements: [],
    },
    [userBasePath]: {
        resourceAccessRequirements: [],
    },
    [violationsPath]: {
        resourceAccessRequirements: [],
    },
    // Reporting and Risk Acceptance must precede generic Vulnerability Management in Body and so here for consistency.
    [vulnManagementReportsPath]: {
        resourceAccessRequirements: ['WorkflowAdministration'],
    },
    [vulnManagementRiskAcceptancePath]: {
        resourceAccessRequirements: [],
    },
    [vulnManagementPath]: {
        resourceAccessRequirements: [],
    },
    [vulnerabilitiesWorkloadCvesPath]: {
        featureFlagDependency: ['ROX_VULN_MGMT_WORKLOAD_CVES'],
        resourceAccessRequirements: [],
    },
    [vulnerabilityReportsPath]: {
        featureFlagDependency: ['ROX_VULN_MGMT_REPORTING_ENHANCEMENTS'],
        resourceAccessRequirements: ['WorkflowAdministration'],
    },
};

type RoutePredicates = {
    hasReadAccess: HasReadAccess;
    isFeatureFlagEnabled: IsFeatureFlagEnabled;
};

export function isRouteEnabled(
    { hasReadAccess, isFeatureFlagEnabled }: RoutePredicates,
    path: string
) {
    const routeDescription = routeDescriptionMap[path];

    if (!routeDescription) {
        // eslint-disable-next-line no-console
        console.warn(`isRouteEnabled for unknown path ${path}`);
        return false; // better to find mistakes than allow loopholes
    }

    const { featureFlagDependency, resourceAccessRequirements } = routeDescription;

    if (Array.isArray(featureFlagDependency)) {
        if (
            !featureFlagDependency.every((featureFlagEnvVar) =>
                isFeatureFlagEnabled(featureFlagEnvVar)
            )
        ) {
            return false;
        }
    }

    return resourceAccessRequirements.every((resourceName) => hasReadAccess(resourceName));
}

/**
 * New Framwork-related route paths
 */

export const urlEntityListTypes = {
    [resourceTypes.NAMESPACE]: 'namespaces',
    [resourceTypes.CLUSTER]: 'clusters',
    [resourceTypes.NODE]: 'nodes',
    [resourceTypes.DEPLOYMENT]: 'deployments',
    [resourceTypes.IMAGE]: 'images',
    [resourceTypes.SECRET]: 'secrets',
    [resourceTypes.POLICY]: 'policies', // TODO verify if used for Configuration Management
    [resourceTypes.CVE]: 'cves', // TODO verify obsolete because non-postgres
    [resourceTypes.IMAGE_CVE]: 'image-cves',
    [resourceTypes.NODE_CVE]: 'node-cves',
    [resourceTypes.CLUSTER_CVE]: 'cluster-cves',
    [resourceTypes.COMPONENT]: 'components', // TODO verify obsolete because non-postgres
    [resourceTypes.NODE_COMPONENT]: 'node-components',
    [resourceTypes.IMAGE_COMPONENT]: 'image-components',
    [standardEntityTypes.CONTROL]: 'controls',
    [rbacConfigTypes.SERVICE_ACCOUNT]: 'serviceaccounts',
    [rbacConfigTypes.SUBJECT]: 'subjects',
    [rbacConfigTypes.ROLE]: 'roles',
};

export const urlEntityTypes = {
    [resourceTypes.NAMESPACE]: 'namespace',
    [resourceTypes.CLUSTER]: 'cluster',
    [resourceTypes.NODE]: 'node',
    [resourceTypes.DEPLOYMENT]: 'deployment',
    [resourceTypes.IMAGE]: 'image',
    [resourceTypes.SECRET]: 'secret',
    [resourceTypes.POLICY]: 'policy', // TODO verify if used for Configuration Management
    [resourceTypes.CVE]: 'cve', // TODO verify obsolete because non-postgres
    [resourceTypes.IMAGE_CVE]: 'image-cve',
    [resourceTypes.NODE_CVE]: 'node-cve',
    [resourceTypes.CLUSTER_CVE]: 'cluster-cve',
    [resourceTypes.COMPONENT]: 'component', // TODO verify obsolete because non-postgres
    [resourceTypes.NODE_COMPONENT]: 'node-component',
    [resourceTypes.IMAGE_COMPONENT]: 'image-component',
    [standardEntityTypes.CONTROL]: 'control',
    [standardEntityTypes.STANDARD]: 'standard',
    [rbacConfigTypes.SERVICE_ACCOUNT]: 'serviceaccount',
    [rbacConfigTypes.SUBJECT]: 'subject',
    [rbacConfigTypes.ROLE]: 'role',
};

const vulnManagementPathToLabelMap = {
    [vulnManagementPath]: 'Dashboard',
    // TODO: add mapping for Deferrals
    [vulnManagementReportsPath]: 'Reporting',
    [vulnManagementRiskAcceptancePath]: 'Risk Acceptance',
};

const vulnerabilitiesPathToLabelMap = {
    [vulnerabilitiesBasePath]: 'Vulnerabilities',
    [vulnerabilitiesWorkloadCvesPath]: 'Workload CVEs',
    [vulnerabilityReportsPath]: 'Vulnerability Reporting',
};

export const basePathToLabelMap = {
    [dashboardPath]: 'Dashboard',
    [networkBasePath]: 'Network Graph',
    [listeningEndpointsBasePath]: 'Listening Endpoints',
    [violationsBasePath]: 'Violations',
    [complianceBasePath]: 'Compliance',
    [complianceEnhancedBasePath]: 'Compliance (2.0)',
    ...vulnerabilitiesPathToLabelMap,
    ...vulnManagementPathToLabelMap,
    [configManagementPath]: 'Configuration Management',
    [riskBasePath]: 'Risk',
    [apidocsPath]: 'API Reference',
    [clustersBasePath]: 'Clusters',
    [policyManagementBasePath]: 'Policy Management',
    [policiesBasePath]: 'Policy Management',
    [policyCategoriesPath]: 'Policy Categories',
    [collectionsBasePath]: 'Collections',
    [integrationsPath]: 'Integrations',
    [accessControlPath]: 'Access Control',
    [accessControlBasePath]: 'Access Control',
    [systemConfigPath]: 'System Configuration',
    [systemHealthPath]: 'System Health',
    [loginPath]: 'Log In',
    [userBasePath]: 'User Profile',
};

const entityListTypeMatcher = `(${Object.values(urlEntityListTypes).join('|')})`;
const entityTypeMatcher = `(${Object.values(urlEntityTypes).join('|')})`;

export const workflowPaths = {
    DASHBOARD: `${mainPath}/:context`,
    LIST: `${mainPath}/:context/:pageEntityListType${entityListTypeMatcher}/:entityId1?/:entityType2?/:entityId2?`,
    ENTITY: `${mainPath}/:context/:pageEntityType${entityTypeMatcher}/:pageEntityId?/:entityType1?/:entityId1?/:entityType2?/:entityId2?`,
};
