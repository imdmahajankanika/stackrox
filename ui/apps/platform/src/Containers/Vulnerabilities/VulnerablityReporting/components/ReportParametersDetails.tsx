import {
    DescriptionList,
    DescriptionListDescription,
    DescriptionListGroup,
    DescriptionListTerm,
    Flex,
    FlexItem,
    Title,
} from '@patternfly/react-core';
import React, { ReactElement } from 'react';

import { ReportFormValues } from 'Containers/Vulnerabilities/VulnerablityReporting/forms/useReportFormValues';
import { fixabilityLabels } from 'constants/reportConstants';
import { getDate } from 'utils/dateUtils';
import {
    cvesDiscoveredSinceLabelMap,
    imageTypeLabelMap,
} from 'Containers/Vulnerabilities/VulnerablityReporting/utils';

import VulnerabilitySeverityIconText from 'Components/PatternFly/IconText/VulnerabilitySeverityIconText';

export type ReportParametersDetailsProps = {
    formValues: ReportFormValues;
};

function ReportParametersDetails({ formValues }: ReportParametersDetailsProps): ReactElement {
    const cveSeverities =
        formValues.reportParameters.cveSeverities.length !== 0 ? (
            formValues.reportParameters.cveSeverities.map((severity) => (
                <li key={severity}>
                    <VulnerabilitySeverityIconText severity={severity} />
                </li>
            ))
        ) : (
            <li>None</li>
        );
    const cveStatuses =
        formValues.reportParameters.cveStatus.length !== 0 ? (
            formValues.reportParameters.cveStatus.map((status) => (
                <li key={status}>{fixabilityLabels[status]}</li>
            ))
        ) : (
            <li>None</li>
        );
    const imageTypes =
        formValues.reportParameters.imageType.length !== 0 ? (
            formValues.reportParameters.imageType.map((type) => (
                <li key={type}>{imageTypeLabelMap[type]}</li>
            ))
        ) : (
            <li>None</li>
        );

    return (
        <Flex direction={{ default: 'column' }}>
            <FlexItem>
                <Title headingLevel="h3">Report parameters</Title>
            </FlexItem>
            <FlexItem flex={{ default: 'flexNone' }}>
                <DescriptionList
                    columnModifier={{
                        default: '2Col',
                        md: '2Col',
                        sm: '1Col',
                    }}
                >
                    <DescriptionListGroup>
                        <DescriptionListTerm>Report name</DescriptionListTerm>
                        <DescriptionListDescription>
                            {formValues.reportParameters.reportName || 'None'}
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Description</DescriptionListTerm>
                        <DescriptionListDescription>
                            {formValues.reportParameters.description || 'None'}
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>CVE severity</DescriptionListTerm>
                        <DescriptionListDescription>
                            <ul>{cveSeverities}</ul>
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>CVE status</DescriptionListTerm>
                        <DescriptionListDescription>
                            <ul>{cveStatuses}</ul>
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Report scope</DescriptionListTerm>
                        <DescriptionListDescription>
                            {formValues.reportParameters.reportScope?.name || 'None'}
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Image type</DescriptionListTerm>
                        <DescriptionListDescription>
                            <ul>{imageTypes}</ul>
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>CVEs discovered since</DescriptionListTerm>
                        <DescriptionListDescription>
                            {formValues.reportParameters.cvesDiscoveredSince === 'START_DATE' &&
                            !!formValues.reportParameters.cvesDiscoveredStartDate
                                ? getDate(formValues.reportParameters.cvesDiscoveredStartDate)
                                : cvesDiscoveredSinceLabelMap[
                                      formValues.reportParameters.cvesDiscoveredSince
                                  ]}
                        </DescriptionListDescription>
                    </DescriptionListGroup>
                </DescriptionList>
            </FlexItem>
        </Flex>
    );
}

export default ReportParametersDetails;
