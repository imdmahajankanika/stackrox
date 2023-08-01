import React, { ReactElement } from 'react';
import { ActionsColumn, TableComposable, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { Bullseye, Card, Divider, Spinner } from '@patternfly/react-core';

import { ReportConfiguration, ReportRequestType } from 'services/ReportsService.types';
import { SlimUser } from 'types/user.proto';
import { getDateTime } from 'utils/dateUtils';
import { getReportFormValuesFromConfiguration } from 'Containers/Vulnerabilities/VulnerablityReporting/utils';
import useSet from 'hooks/useSet';

import NotFoundMessage from 'Components/NotFoundMessage/NotFoundMessage';
import EmptyStateTemplate from 'Components/PatternFly/EmptyStateTemplate/EmptyStateTemplate';
import { CubesIcon } from '@patternfly/react-icons';
import { saveFile } from 'services/DownloadService';
import LastRunStatusState from '../VulnReports/LastRunStatusState';
import ReportParametersDetails from '../components/ReportParametersDetails';
import DeliveryDestinationsDetails from '../components/DeliveryDestinationsDetails';
import ScheduleDetails from '../components/ScheduleDetails';
import useFetchReportHistory from '../api/useFetchReportHistory';
import useDeleteDownloadModal from '../hooks/useDeleteDownloadModal';
import DeleteModal from '../components/DeleteModal';

export type RunTypeProps = {
    reportRequestType: ReportRequestType;
    user: SlimUser;
};

function RunType({ reportRequestType, user }: RunTypeProps): ReactElement {
    if (reportRequestType === 'ON_DEMAND') {
        return (
            <div>
                On-demand / <span className="pf-u-color-200">Requested by {user.name}</span>
            </div>
        );
    }
    if (reportRequestType === 'SCHEDULED') {
        return <div>Scheduled</div>;
    }
    return <div>-</div>;
}

export type RunHistoryProps = {
    reportId: string;
};

function RunHistory({ reportId }: RunHistoryProps) {
    const { reportSnapshots, isLoading, error, fetchReportSnapshots } = useFetchReportHistory({
        id: reportId,
    });
    const expandedRowSet = useSet<string>();

    const {
        openDeleteDownloadModal,
        isDeleteDownloadModalOpen,
        closeDeleteDownloadModal,
        isDeletingDownload,
        onDeleteDownload,
        deleteDownloadError,
    } = useDeleteDownloadModal({
        onCompleted: fetchReportSnapshots,
    });

    if (isLoading) {
        return (
            <Bullseye>
                <Spinner isSVG />
            </Bullseye>
        );
    }

    if (error) {
        return (
            <NotFoundMessage
                title="Error fetching report history"
                message={error || 'No data available'}
            />
        );
    }

    if (!reportSnapshots.length) {
        return (
            <Bullseye>
                <EmptyStateTemplate title="No run history" headingLevel="h2" icon={CubesIcon} />
            </Bullseye>
        );
    }

    return (
        <>
            <TableComposable aria-label="Simple table" variant="compact">
                <Thead>
                    <Tr>
                        <Td>{/* Header for expanded column */}</Td>
                        <Th>Run time</Th>
                        <Th>Status</Th>
                        <Th>Run type</Th>
                        <Td>{/* Header for table actions column */}</Td>
                    </Tr>
                </Thead>
                {reportSnapshots.map(
                    (
                        {
                            id,
                            name,
                            description,
                            vulnReportFilters,
                            collectionSnapshot,
                            schedule,
                            notifiers,
                            reportStatus,
                            user,
                        },
                        rowIndex
                    ) => {
                        const isExpanded = expandedRowSet.has(id);
                        const reportConfiguration: ReportConfiguration = {
                            id,
                            name,
                            description,
                            type: 'VULNERABILITY',
                            vulnReportFilters,
                            notifiers,
                            schedule,
                            resourceScope: {
                                collectionScope: {
                                    collectionId: collectionSnapshot.id,
                                    collectionName: collectionSnapshot.name,
                                },
                            },
                        };
                        const formValues =
                            getReportFormValuesFromConfiguration(reportConfiguration);
                        const hasDownloadableReport =
                            reportStatus.runState === 'SUCCESS' &&
                            reportStatus.reportNotificationMethod === 'DOWNLOAD';
                        const rowActions = [
                            {
                                title: 'Download report',
                                onClick: (event) => {
                                    event.preventDefault();
                                    // @TODO: This will fail until API provides the "reportId"
                                    return saveFile({
                                        method: 'get',
                                        url: `/v2/reports/jobs/${id}/download`,
                                        data: null,
                                        timeout: 300000,
                                        name: `${name}-report.txt`,
                                    });
                                },
                            },
                            {
                                title: (
                                    <span className="pf-u-danger-color-100">Delete download</span>
                                ),
                                onClick: (event) => {
                                    event.preventDefault();
                                    // @TODO: This will fail until API provides the "reportId"
                                    openDeleteDownloadModal(id);
                                },
                            },
                        ];

                        return (
                            <Tbody key={id} isExpanded={isExpanded}>
                                <Tr>
                                    <Td
                                        expand={{
                                            rowIndex,
                                            isExpanded,
                                            onToggle: () => expandedRowSet.toggle(id),
                                        }}
                                    />
                                    <Td dataLabel="Run time">
                                        {getDateTime(reportStatus.completedAt)}
                                    </Td>
                                    <Td dataLabel="Status">
                                        <LastRunStatusState reportStatus={reportStatus} />
                                    </Td>
                                    <Td dataLabel="Run type">
                                        <RunType
                                            reportRequestType={reportStatus.reportRequestType}
                                            user={user}
                                        />
                                    </Td>
                                    <Td isActionCell>
                                        {hasDownloadableReport && (
                                            <ActionsColumn items={rowActions} />
                                        )}
                                    </Td>
                                </Tr>
                                <Tr isExpanded={isExpanded}>
                                    <Td colSpan={4}>
                                        <Card className="pf-u-m-md pf-u-p-md" isFlat>
                                            <ReportParametersDetails formValues={formValues} />
                                            <Divider component="div" className="pf-u-py-md" />
                                            <DeliveryDestinationsDetails formValues={formValues} />
                                            <Divider component="div" className="pf-u-py-md" />
                                            <ScheduleDetails formValues={formValues} />
                                        </Card>
                                    </Td>
                                </Tr>
                            </Tbody>
                        );
                    }
                )}
            </TableComposable>
            <DeleteModal
                title="Delete downloadable report?"
                isOpen={isDeleteDownloadModalOpen}
                onClose={closeDeleteDownloadModal}
                isDeleting={isDeletingDownload}
                onDelete={onDeleteDownload}
                error={deleteDownloadError}
            >
                All data in this downloadable report will be deleted. Regenerating a downloadable
                report will requre the download process to start over.
            </DeleteModal>
        </>
    );
}

export default RunHistory;
