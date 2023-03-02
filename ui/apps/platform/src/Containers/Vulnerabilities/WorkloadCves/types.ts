export type VulnerabilitySeverityLabel = 'Critical' | 'Important' | 'Moderate' | 'Low';
export type FixableStatus = 'Fixable' | 'Not fixable';

export type DefaultFilters = {
    Severity: VulnerabilitySeverityLabel[];
    Fixable: FixableStatus[];
};

export type VulnMgmtLocalStorage = {
    preferences: {
        defaultFilters: DefaultFilters;
    };
};

const detailsTabValues = ['Vulnerabilities', 'Resources'] as const;

export type DetailsTab = typeof detailsTabValues[number];

export function isDetailsTab(value: unknown): value is DetailsTab {
    return detailsTabValues.some((tab) => tab === value);
}
