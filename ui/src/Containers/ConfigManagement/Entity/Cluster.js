import React, { useContext } from 'react';
import entityTypes from 'constants/entityTypes';
import queryService from 'modules/queryService';
import gql from 'graphql-tag';
import Query from 'Components/ThrowingQuery';
import Loader from 'Components/Loader';
import PageNotFound from 'Components/PageNotFound';
import CollapsibleSection from 'Components/CollapsibleSection';
import RelatedEntityListCount from 'Containers/ConfigManagement/Entity/widgets/RelatedEntityListCount';
import Metadata from 'Containers/ConfigManagement/Entity/widgets/Metadata';
import Tabs from 'Components/Tabs';
import TabContent from 'Components/TabContent';
import searchContext from 'Containers/searchContext';
import { NODE_FRAGMENT } from 'queries/node';
import { DEPLOYMENT_FRAGMENT } from 'queries/deployment';
import { NAMESPACE_FRAGMENT } from 'queries/namespace';
import { SUBJECT_WITH_CLUSTER_FRAGMENT } from 'queries/subject';
import { ROLE_FRAGMENT } from 'queries/role';
import { SECRET_FRAGMENT } from 'queries/secret';
import { SERVICE_ACCOUNT_FRAGMENT } from 'queries/serviceAccount';
import { CONTROL_FRAGMENT, CONTROLS_FRAGMENT } from 'queries/controls';
import { entityComponentPropTypes, entityComponentDefaultProps } from 'constants/entityPageProps';
import { POLICY_FRAGMENT } from 'queries/policy';
import { IMAGE_FRAGMENT } from 'queries/image';
import { standardLabels } from 'messages/standards';
import getSubListFromEntity from '../List/utilities/getSubListFromEntity';
import NodesWithFailedControls from './widgets/NodesWithFailedControls';
import DeploymentsWithFailedPolicies from './widgets/DeploymentsWithFailedPolicies';
import EntityList from '../List/EntityList';

const Cluster = ({ id, entityListType, query, entityContext }) => {
    const searchParam = useContext(searchContext);

    const queryObject = { ...query[searchParam] };

    if (entityListType === entityTypes.POLICY) queryObject['Lifecycle Stage'] = 'DEPLOY';
    if (!queryObject.Standard) queryObject.Standard = 'CIS';

    const variables = {
        id,
        query: queryService.objectToWhereClause(queryObject)
    };

    const QUERY = gql`
        query getCluster($id: ID!${entityListType ? ', $query: String' : ''}) {
            cluster(id: $id) {
                id
                name
                admissionController
                centralApiEndpoint
                ${
                    entityListType === entityTypes.IMAGE
                        ? 'images(query: $query) { ...imageFields }'
                        : 'imageCount'
                }
                ${
                    entityListType === entityTypes.NODE
                        ? 'nodes(query: $query) { ...nodeFields }'
                        : 'nodeCount'
                }
                ${
                    entityListType === entityTypes.DEPLOYMENT
                        ? 'deployments(query: $query) { ...deploymentFields }'
                        : 'deploymentCount'
                }
                ${
                    entityListType === entityTypes.NAMESPACE
                        ? 'namespaces(query: $query) { ...namespaceFields }'
                        : 'namespaceCount'
                }
                ${
                    entityListType === entityTypes.SUBJECT
                        ? 'subjects(query: $query) { ...subjectWithClusterFields }'
                        : 'subjectCount'
                }
                ${
                    entityListType === entityTypes.ROLE
                        ? 'k8sroles(query: $query) { ...k8roleFields }'
                        : 'k8sroleCount'
                }
                ${
                    entityListType === entityTypes.SECRET
                        ? 'secrets(query: $query) { ...secretFields }'
                        : 'secretCount'
                }
                ${
                    entityListType === entityTypes.POLICY
                        ? 'policies(query: $query) { ...policyFields }'
                        : 'policyCount(query: "Lifecycle Stage:DEPLOY")'
                }
                ${
                    entityListType === entityTypes.SERVICE_ACCOUNT
                        ? 'serviceAccounts(query: $query) { ...serviceAccountFields }'
                        : 'serviceAccountCount'
                }
                ${
                    entityListType === entityTypes.CONTROL
                        ? 'complianceResults(query: $query) { ...controlFields } controls(query: $query) { ...controlsListFields}'
                        : 'controls(query: "Standard:CIS") { ...controlsListFields}'
                }
                status {
                    orchestratorMetadata {
                        version
                        buildDate
                    }
                }
            }
        }
        ${entityListType === entityTypes.IMAGE ? IMAGE_FRAGMENT : ''}
        ${entityListType === entityTypes.NODE ? NODE_FRAGMENT : ''}
        ${entityListType === entityTypes.DEPLOYMENT ? DEPLOYMENT_FRAGMENT : ''}
        ${entityListType === entityTypes.NAMESPACE ? NAMESPACE_FRAGMENT : ''}
        ${entityListType === entityTypes.SUBJECT ? SUBJECT_WITH_CLUSTER_FRAGMENT : ''}
        ${entityListType === entityTypes.ROLE ? ROLE_FRAGMENT : ''}
        ${entityListType === entityTypes.SERVICE_ACCOUNT ? SERVICE_ACCOUNT_FRAGMENT : ''}
        ${entityListType === entityTypes.SECRET ? SECRET_FRAGMENT : ''}
        ${entityListType === entityTypes.POLICY ? POLICY_FRAGMENT : ''}
        ${entityListType === entityTypes.CONTROL ? CONTROL_FRAGMENT : ''}
        ${CONTROLS_FRAGMENT}
    `;

    return (
        <Query query={QUERY} variables={variables}>
            {({ loading, data }) => {
                if (loading || !data) return <Loader transparent />;
                const { cluster: entity } = data;
                if (!entity) return <PageNotFound resourceType={entityTypes.CLUSTER} />;

                const { complianceResults = [], controls } = entity;

                if (entityListType) {
                    let listData = getSubListFromEntity(entity, entityListType);
                    if (entityListType === entityTypes.CONTROL) {
                        const failedComplianceResults = complianceResults
                            .filter(cr => cr.value.overallState === 'COMPLIANCE_STATE_FAILURE')
                            .map(cr => ({
                                ...cr,
                                control: {
                                    ...cr.control,
                                    standard: standardLabels[cr.control.standardId]
                                }
                            }));
                        listData = controls.map(control => ({
                            ...control,
                            standard: standardLabels[control.standardId],
                            control: `${control.name} - ${control.description}`,
                            passing: !failedComplianceResults.find(
                                cr => cr.control.id === control.id
                            )
                        }));
                    }
                    return (
                        <EntityList
                            entityListType={entityListType}
                            data={listData}
                            entityContext={{ ...entityContext, [entityTypes.CLUSTER]: id }}
                            query={query}
                        />
                    );
                }

                const {
                    name,
                    nodeCount,
                    deploymentCount,
                    namespaceCount,
                    subjectCount,
                    serviceAccountCount,
                    k8sroleCount,
                    secretCount,
                    imageCount,
                    status: { orchestratorMetadata = null }
                } = entity;

                const { version = 'N/A' } = orchestratorMetadata;

                const metadataKeyValuePairs = [
                    {
                        key: 'K8s version',
                        value: version
                    }
                ];

                return (
                    <div className="w-full" id="capture-dashboard-stretch">
                        <CollapsibleSection title="Cluster Details">
                            <div className="flex flex-wrap pdf-page">
                                <Metadata
                                    className="mx-4 min-w-48 bg-base-100 h-48 mb-4"
                                    keyValuePairs={metadataKeyValuePairs}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Nodes"
                                    value={nodeCount}
                                    entityType={entityTypes.NODE}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Namespaces"
                                    value={namespaceCount}
                                    entityType={entityTypes.NAMESPACE}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Deployments"
                                    value={deploymentCount}
                                    entityType={entityTypes.DEPLOYMENT}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Secrets"
                                    value={secretCount}
                                    entityType={entityTypes.SECRET}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Images"
                                    value={imageCount}
                                    entityType={entityTypes.IMAGE}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Users & Groups"
                                    value={subjectCount}
                                    entityType={entityTypes.SUBJECT}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Service Accounts"
                                    value={serviceAccountCount}
                                    entityType={entityTypes.SERVICE_ACCOUNT}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Roles"
                                    value={k8sroleCount}
                                    entityType={entityTypes.ROLE}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="CIS Controls"
                                    value={controls.length}
                                    entityType={entityTypes.CONTROL}
                                />
                            </div>
                        </CollapsibleSection>
                        <CollapsibleSection title="Cluster Findings">
                            <div className="flex pdf-page pdf-stretch rounded relative rounded mb-4 ml-4 mr-4">
                                <Tabs
                                    hasTabSpacing
                                    headers={[{ text: 'Policies' }, { text: 'CIS Controls' }]}
                                >
                                    <TabContent>
                                        <DeploymentsWithFailedPolicies
                                            query={queryService.objectToWhereClause({
                                                Cluster: name
                                            })}
                                            message="No deployments violating policies in this cluster"
                                        />
                                    </TabContent>
                                    <TabContent>
                                        <NodesWithFailedControls />
                                    </TabContent>
                                </Tabs>
                            </div>
                        </CollapsibleSection>
                    </div>
                );
            }}
        </Query>
    );
};

Cluster.propTypes = entityComponentPropTypes;
Cluster.defaultProps = entityComponentDefaultProps;

export default Cluster;
