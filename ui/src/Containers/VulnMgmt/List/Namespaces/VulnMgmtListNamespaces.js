import React from 'react';
import gql from 'graphql-tag';
import pluralize from 'pluralize';

import queryService from 'modules/queryService';
import TableCellLink from 'Components/TableCellLink';
import CVEStackedPill from 'Components/CVEStackedPill';
import LabelChip from 'Components/LabelChip';
import DateTimeField from 'Components/DateTimeField';
import { sortDate } from 'sorters/sorters';
import { defaultHeaderClassName, defaultColumnClassName } from 'Components/Table';
import entityTypes from 'constants/entityTypes';
import { generateURLToFromTable } from 'modules/URLReadWrite';
import WorkflowListPage from 'Containers/Workflow/WorkflowListPage';
import { NAMESPACE_LIST_FRAGMENT } from 'Containers/VulnMgmt/VulnMgmt.fragments';
import { workflowListPropTypes, workflowListDefaultProps } from 'constants/entityPageProps';

export function getNamespaceTableColumns(workflowState) {
    const tableColumns = [
        {
            Header: 'Id',
            headerClassName: 'hidden',
            className: 'hidden',
            accessor: 'id'
        },
        {
            Header: `Namespace`,
            headerClassName: `w-1/6 ${defaultHeaderClassName}`,
            className: `w-1/6 ${defaultColumnClassName}`,
            accessor: 'metadata.name'
        },
        {
            Header: `CVEs`,
            headerClassName: `w-1/6 ${defaultHeaderClassName}`,
            className: `w-1/6 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => {
                const { vulnCounter, id } = original;
                const url = generateURLToFromTable(workflowState, id, entityTypes.CVE);
                return <CVEStackedPill vulnCounter={vulnCounter} url={url} pdf={pdf} />;
            }
        },
        {
            Header: `Cluster`,
            headerClassName: `w-1/8 ${defaultHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => {
                const { metadata } = original;
                const { clusterName, clusterId, id } = metadata;
                const url = generateURLToFromTable(
                    workflowState,
                    id,
                    entityTypes.CLUSTER,
                    clusterId
                );
                return <TableCellLink pdf={pdf} url={url} text={clusterName} />;
            }
        },
        {
            Header: `Deployments`,
            headerClassName: `w-1/8 ${defaultHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => {
                const { deploymentCount, metadata } = original;
                const url = generateURLToFromTable(
                    workflowState,
                    metadata.id,
                    entityTypes.DEPLOYMENT
                );
                const text = `${deploymentCount} ${pluralize(
                    entityTypes.DEPLOYMENT.toLowerCase(),
                    deploymentCount
                )}`;
                return <TableCellLink pdf={pdf} url={url} text={text} />;
            }
        },
        {
            Header: `Images`,
            headerClassName: `w-1/8 ${defaultHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => {
                const { imageCount, metadata } = original;
                const url = generateURLToFromTable(workflowState, metadata.id, entityTypes.IMAGE);
                const text = `${imageCount} ${pluralize(
                    entityTypes.IMAGE.toLowerCase(),
                    imageCount
                )}`;
                return <TableCellLink pdf={pdf} url={url} text={text} />;
            }
        },
        {
            Header: `Policies`,
            headerClassName: `w-1/8 ${defaultHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => {
                const { policyCount, metadata } = original;
                const url = generateURLToFromTable(workflowState, metadata.id, entityTypes.POLICY);
                const text = `${policyCount} ${pluralize(
                    entityTypes.POLICY.toLowerCase(),
                    policyCount
                )}`;
                return <TableCellLink pdf={pdf} url={url} text={text} />;
            }
        },
        {
            Header: `Policy status`,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original }) => {
                const { policyStatus } = original;
                return policyStatus.status === 'pass' ? (
                    <LabelChip text="Pass" type="success" />
                ) : (
                    <LabelChip text="Fail" type="alert" />
                );
            },
            id: 'policyStatus'
        },
        {
            Header: `Latest violation`,
            headerClassName: `w-1/8 ${defaultHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            Cell: ({ original }) => {
                const { latestViolation } = original;
                return <DateTimeField date={latestViolation} />;
            },
            accessor: 'latestViolation',
            sortMethod: sortDate
        },
        {
            Header: `Risk`,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            accessor: 'metadata.priority'
        }
    ];
    return tableColumns.filter(col => col);
}

const VulnMgmtNamespaces = ({ selectedRowId, search, sort, page }) => {
    const query = gql`
        query getNamespaces {
            results: namespaces {
                ...namespaceListFields
            }
        }
        ${NAMESPACE_LIST_FRAGMENT}
    `;

    const queryOptions = {
        variables: {
            query: queryService.objectToWhereClause(search)
        }
    };

    const defaultNamespaceSort = [
        {
            id: 'metadata.priority',
            desc: false
        }
    ];

    return (
        <WorkflowListPage
            query={query}
            queryOptions={queryOptions}
            entityListType={entityTypes.NAMESPACE}
            getTableColumns={getNamespaceTableColumns}
            selectedRowId={selectedRowId}
            idAttribute="metadata.id"
            search={search}
            sort={sort}
            page={page}
            defaultSorted={defaultNamespaceSort}
        />
    );
};

VulnMgmtNamespaces.propTypes = workflowListPropTypes;
VulnMgmtNamespaces.defaultProps = workflowListDefaultProps;

export default VulnMgmtNamespaces;
