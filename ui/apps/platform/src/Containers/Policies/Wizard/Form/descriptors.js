import {
    lifecycleStageLabels,
    portExposureLabels,
    envVarSrcLabels,
    rbacPermissionLabels,
    policyCriteriaCategories,
    mountPropagationLabels,
} from 'messages/common';
import { knownBackendFlags } from 'utils/featureFlags';
import { clientOnlyExclusionFieldNames } from './whitelistFieldNames';

const equalityOptions = [
    { label: 'Is greater than', value: '>' },
    {
        label: 'Is greater than or equal to',
        value: '>=',
    },
    { label: 'Is equal to', value: '=' },
    {
        label: 'Is less than or equal to',
        value: '<=',
    },
    { label: 'Is less than', value: '<' },
];

const cpuResource = (label, policy, field) => ({
    label,
    name: label,
    jsonpath: `fields.${policy}.${field}`,
    category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
    type: 'group',
    jsonpaths: [
        {
            jsonpath: `fields.${policy}.${field}.op`,
            type: 'select',
            options: equalityOptions,
            subpath: 'key',
        },
        {
            jsonpath: `fields.${policy}.${field}.value`,
            type: 'number',
            placeholder: '# of cores',
            min: 0,
            step: 0.1,
            subpath: 'value',
        },
    ],
    required: false,
    default: false,
    canBooleanLogic: true,
});

const capabilities = [
    'AUDIT_CONTROL',
    'AUDIT_READ',
    'AUDIT_WRITE',
    'BLOCK_SUSPEND',
    'CHOWN',
    'DAC_OVERRIDE',
    'DAC_READ_SEARCH',
    'FOWNER',
    'FSETID',
    'IPC_LOCK',
    'IPC_OWNER',
    'KILL',
    'LEASE',
    'LINUX_IMMUTABLE',
    'MAC_ADMIN',
    'MAC_OVERRIDE',
    'MKNOD',
    'NET_ADMIN',
    'NET_BIND_SERVICE',
    'NET_BROADCAST',
    'NET_RAW',
    'SETGID',
    'SETFCAP',
    'SETPCAP',
    'SETUID',
    'SYS_ADMIN',
    'SYS_BOOT',
    'SYS_CHROOT',
    'SYS_MODULE',
    'SYS_NICE',
    'SYS_PACCT',
    'SYS_PTRACE',
    'SYS_RAWIO',
    'SYS_RESOURCE',
    'SYS_TIME',
    'SYS_TTY_CONFIG',
    'SYSLOG',
    'WAKE_ALARM',
].map((cap) => ({ label: cap, value: cap }));

const memoryResource = (label, policy, field) => ({
    label,
    name: label,
    jsonpath: `fields.${policy}.${field}`,
    category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
    type: 'group',
    jsonpaths: [
        {
            jsonpath: `fields.${policy}.${field}.op`,
            type: 'select',
            options: equalityOptions,
            subpath: 'key',
        },
        {
            jsonpath: `fields.${policy}.${field}.value`,
            type: 'number',
            placeholder: '# MB',
            min: 0,
            subpath: 'value',
        },
    ],
    required: false,
    default: false,
    canBooleanLogic: true,
});

// A descriptor for every option on the policy creation page.
const policyStatusDescriptor = [
    {
        label: '',
        header: true,
        jsonpath: 'disabled',
        type: 'toggle',
        required: false,
        reverse: true,
        default: true,
    },
];

const policyDetailsFormDescriptor = [
    {
        label: 'Name',
        hideInnerLabel: true,
        jsonpath: 'name',
        type: 'text',
        required: true,
        default: true,
    },
    {
        label: 'Severity',
        hideInnerLabel: true,
        jsonpath: 'severity',
        type: 'select',
        options: [
            { label: 'Critical', value: 'CRITICAL_SEVERITY' },
            { label: 'High', value: 'HIGH_SEVERITY' },
            { label: 'Medium', value: 'MEDIUM_SEVERITY' },
            { label: 'Low', value: 'LOW_SEVERITY' },
        ],
        placeholder: 'Select a severity level',
        required: true,
        default: true,
    },
    {
        label: 'Lifecycle Stages',
        jsonpath: 'lifecycleStages',
        type: 'multiselect',
        options: Object.keys(lifecycleStageLabels).map((key) => ({
            label: lifecycleStageLabels[key],
            value: key,
        })),
        required: true,
        default: true,
    },
    {
        label: 'Description',
        jsonpath: 'description',
        type: 'textarea',
        placeholder: 'What does this policy do?',
        required: false,
        default: true,
    },
    {
        label: 'Rationale',
        jsonpath: 'rationale',
        type: 'textarea',
        placeholder: 'Why does this policy exist?',
        required: false,
        default: true,
    },
    {
        label: 'Remediation',
        jsonpath: 'remediation',
        type: 'textarea',
        placeholder: 'What can an operator do to resolve any violations?',
        required: false,
        default: true,
    },
    {
        label: 'Categories',
        jsonpath: 'categories',
        type: 'multiselect-creatable',
        options: [],
        required: true,
        default: true,
    },
    {
        label: 'Notifications',
        jsonpath: 'notifiers',
        type: 'multiselect',
        options: [],
        required: false,
        default: true,
    },
    {
        label: 'Restrict to Scope',
        jsonpath: 'scope',
        type: 'scope',
        options: [],
        required: false,
        default: true,
    },
    {
        label: 'Exclude by Scope',
        jsonpath: clientOnlyExclusionFieldNames.EXCLUDED_DEPLOYMENT_SCOPES,
        type: 'whitelistScope',
        options: [],
        required: false,
        default: true,
    },
    {
        label: 'Excluded Images (Build Lifecycle only)',
        jsonpath: clientOnlyExclusionFieldNames.EXCLUDED_IMAGE_NAMES,
        type: 'multiselect-creatable',
        options: [],
        required: false,
        default: true,
    },
];

const policyConfigurationDescriptor = [
    {
        label: 'Image Registry',
        name: 'Image Registry',
        longName: 'Image pulled from registry',
        negatedName: 'Image not pulled from registry',
        jsonpath: 'fields.imageName.registry',
        category: policyCriteriaCategories.IMAGE_REGISTRY,
        type: 'text',
        placeholder: 'docker.io',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Image Remote',
        name: 'Image Remote',
        longName: 'Image name in the registry',
        negatedName: `Image name in the registry doesn't match`,
        jsonpath: 'fields.imageName.remote',
        category: policyCriteriaCategories.IMAGE_REGISTRY,
        type: 'text',
        placeholder: 'library/nginx',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Image Tag',
        name: 'Image Tag',
        negatedName: `Image tag doesn't match`,
        jsonpath: 'fields.imageName.tag',
        category: policyCriteriaCategories.IMAGE_REGISTRY,
        type: 'text',
        placeholder: 'latest',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Days since image was created',
        name: 'Image Age',
        longName: 'Minimum days since image was built',
        jsonpath: 'fields.imageAgeDays',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'number',
        placeholder: '1',
        required: false,
        default: false,
        canBooleanLogic: false,
    },
    {
        label: 'Days since image was last scanned',
        name: 'Image Scan Age',
        longName: 'Minimum days since last image scan',
        jsonpath: 'fields.scanAgeDays',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'number',
        placeholder: '1',
        required: false,
        default: false,
        canBooleanLogic: false,
    },
    {
        label: 'Image User',
        name: 'Image User',
        negatedName: `Image user is not`,
        jsonpath: 'fields.imageUser',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'text',
        placeholder: '0',
        required: false,
        default: false,
        canBooleanLogic: false,
    },
    {
        label: 'Dockerfile Line',
        name: 'Dockerfile Line',
        longName: 'Disallowed Dockerfile line',
        jsonpath: 'fields.lineRule',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.lineRule.instruction',
                type: 'select',
                options: [
                    { label: 'FROM', value: 'FROM' },
                    { label: 'LABEL', value: 'LABEL' },
                    { label: 'RUN', value: 'RUN' },
                    { label: 'CMD', value: 'CMD' },
                    { label: 'EXPOSE', value: 'EXPOSE' },
                    { label: 'ENV', value: 'ENV' },
                    { label: 'ADD', value: 'ADD' },
                    { label: 'COPY', value: 'COPY' },
                    { label: 'ENTRYPOINT', value: 'ENTRYPOINT' },
                    { label: 'VOLUME', value: 'VOLUME' },
                    { label: 'USER', value: 'USER' },
                    { label: 'WORKDIR', value: 'WORKDIR' },
                    { label: 'ONBUILD', value: 'ONBUILD' },
                ],
                label: 'Instruction',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.lineRule.value',
                name: 'value',
                type: 'text',
                label: 'Arguments',
                placeholder: 'Any',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Image is NOT Scanned',
        name: 'Unscanned Image',
        longName: 'Image Scan Status',
        jsonpath: 'fields.noScanExists',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Scanned',
                value: false,
            },
            {
                text: 'Not scanned',
                value: true,
            },
        ],
        required: false,
        default: false,
        defaultValue: true,
        disabled: true,
        canBooleanLogic: false,
    },
    {
        label: 'CVSS',
        name: 'CVSS',
        longName: 'Common Vulnerability Scoring System (CVSS) Score',
        jsonpath: 'fields.cvss',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.cvss.op',
                type: 'select',
                options: equalityOptions,
                subpath: 'key',
            },
            {
                jsonpath: 'fields.cvss.value',
                type: 'number',
                placeholder: '0-10',
                max: 10,
                min: 0,
                step: 0.1,
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Fixed By',
        name: 'Fixed By',
        longName: 'Version in which vulnerability is fixed',
        negatedName: `Version in which vulnerability is fixed doesn't match`,
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        jsonpath: 'fields.fixedBy',
        type: 'text',
        placeholder: '.*',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'CVE',
        name: 'CVE',
        longName: 'Common Vulnerabilities and Exposures (CVE) identifier',
        negatedName: `Common Vulnerabilities and Exposures (CVE) identifier doesn't match`,
        jsonpath: 'fields.cve',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'text',
        placeholder: 'CVE-2017-11882',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Image Component',
        name: 'Image Component',
        jsonpath: 'fields.component',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.component.name',
                type: 'text',
                label: 'Component Name',
                placeholder: 'example',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.component.version',
                type: 'text',
                label: 'Version',
                placeholder: 'Any',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Image OS',
        name: 'Image OS',
        longName: 'Image Operating System',
        negatedName: `Image Operating System doesn't match`,
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'text',
        placeholder: 'ubuntu:19.04',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Environment Variable',
        name: 'Environment Variable',
        jsonpath: 'fields.env',
        category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.env.key',
                type: 'text',
                label: 'Key',
                placeholder: 'Any',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.env.value',
                type: 'text',
                label: 'Value',
                placeholder: 'Any',
                subpath: 'value',
            },
            {
                jsonpath: 'fields.env.envVarSource',
                type: 'select',
                options: Object.keys(envVarSrcLabels).map((key) => ({
                    label: envVarSrcLabels[key],
                    value: key,
                })),
                label: 'Value From',
                placeholder: 'Select one',
                subpath: 'source',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Disallowed Annotation',
        name: 'Disallowed Annotation',
        jsonpath: 'fields.disallowedAnnotation',
        category: policyCriteriaCategories.DEPLOYMENT_METADATA,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.disallowedAnnotation.key',
                type: 'text',
                label: 'Key',
                placeholder: 'Any',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.disallowedAnnotation.value',
                type: 'text',
                label: 'Value',
                placeholder: 'Any',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Required Label',
        name: 'Required Label',
        longName: 'Required Deployment Label',
        jsonpath: 'fields.requiredLabel',
        category: policyCriteriaCategories.DEPLOYMENT_METADATA,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.requiredLabel.key',
                type: 'text',
                label: 'Key',
                placeholder: 'owner',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.requiredLabel.value',
                type: 'text',
                label: 'Value',
                placeholder: '.*',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Required Annotation',
        name: 'Required Annotation',
        longName: 'Required Deployment Annotation',
        jsonpath: 'fields.requiredAnnotation',
        category: policyCriteriaCategories.DEPLOYMENT_METADATA,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.requiredAnnotation.key',
                type: 'text',
                label: 'Key',
                placeholder: 'owner',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.requiredAnnotation.value',
                type: 'text',
                label: 'Value',
                placeholder: '.*',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Volume Name',
        name: 'Volume Name',
        negatedName: `Volume name doesn't match`,
        jsonpath: 'fields.volumePolicy.name',
        category: policyCriteriaCategories.STORAGE,
        type: 'text',
        placeholder: 'docker-socket',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Volume Source',
        name: 'Volume Source',
        negatedName: `Volume source doesn't match`,
        jsonpath: 'fields.volumePolicy.source',
        category: policyCriteriaCategories.STORAGE,
        type: 'text',
        placeholder: '/var/run/docker.sock',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Volume Destination',
        name: 'Volume Destination',
        negatedName: `Volume destination doesn't match`,
        jsonpath: 'fields.volumePolicy.destination',
        category: policyCriteriaCategories.STORAGE,
        type: 'text',
        placeholder: '/var/run/docker.sock',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Volume Type',
        name: 'Volume Type',
        negatedName: `Volume type doesn't match`,
        jsonpath: 'fields.volumePolicy.type',
        category: policyCriteriaCategories.STORAGE,
        type: 'text',
        placeholder: 'bind, secret',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Writable Mounted Volume',
        name: 'Writable Mounted Volume',
        longName: 'Mounted Volume Writability',
        jsonpath: 'fields.volumePolicy.readOnly',
        category: policyCriteriaCategories.STORAGE,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Writable',
                value: true,
            },
            {
                text: 'Read-only',
                value: false,
            },
        ],
        required: false,
        default: false,
        defaultValue: false,
        reverse: true,
        canBooleanLogic: false,
    },
    {
        label: 'Mount Propagation',
        name: 'Mount Propagation',
        negatedName: 'Mount Propagation is not',
        jsonpath: 'fields.volumePolicy.mountPropagation',
        category: policyCriteriaCategories.STORAGE,
        type: 'multiselect',
        options: Object.keys(mountPropagationLabels).map((key) => ({
            label: mountPropagationLabels[key],
            value: key,
        })),
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Protocol',
        name: 'Exposed Port Protocol',
        negatedName: `Exposed Port Protocol doesn't match`,
        jsonpath: 'fields.portPolicy.protocol',
        category: policyCriteriaCategories.NETWORKING,
        type: 'text',
        placeholder: 'tcp',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Exposed Node Port',
        name: 'Exposed Node Port',
        negatedName: `Exposed node port doesn't match`,
        jsonpath: 'fields.nodePortPolicy.exposedNodePort',
        category: policyCriteriaCategories.NETWORKING,
        type: 'number',
        placeholder: '22',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Port',
        name: 'Exposed Port',
        negatedName: `Exposed port doesn't match`,
        jsonpath: 'fields.portPolicy.port',
        category: policyCriteriaCategories.NETWORKING,
        type: 'number',
        placeholder: '22',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    cpuResource('Container CPU Request', 'containerResourcePolicy', 'cpuResourceRequest'),
    cpuResource('Container CPU Limit', 'containerResourcePolicy', 'cpuResourceLimit'),
    memoryResource('Container Memory Request', 'containerResourcePolicy', 'memoryResourceRequest'),
    memoryResource('Container Memory Limit', 'containerResourcePolicy', 'memoryResourceLimit'),
    {
        label: 'Privileged',
        name: 'Privileged Container',
        longName: 'Privileged Container Status',
        jsonpath: 'fields.privileged',
        category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Privileged Container',
                value: true,
            },
            {
                text: 'Not a Privileged Container',
                value: false,
            },
        ],
        required: false,
        default: false,
        defaultValue: true,
        disabled: true,
        canBooleanLogic: false,
    },
    {
        label: 'Read-Only Root Filesystem',
        name: 'Read-Only Root Filesystem',
        longName: 'Root Filesystem Writability',
        jsonpath: 'fields.readOnlyRootFs',
        category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Read-Only',
                value: true,
            },
            {
                text: 'Writable',
                value: false,
            },
        ],
        required: false,
        default: false,
        defaultValue: false,
        disabled: true,
        canBooleanLogic: false,
    },
    {
        label: 'Share Host PID Namespace',
        name: 'Host PID',
        longName: 'Host PID',
        jsonpath: 'fields.hostPid',
        category: policyCriteriaCategories.DEPLOYMENT_METADATA,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Uses Host PID namespace',
                value: true,
            },
            {
                text: 'Does Not Use Host PID namespace',
                value: false,
            },
        ],
        required: false,
        default: false,
        defaultValue: true,
        disabled: true,
        canBooleanLogic: false,
    },
    {
        label: 'Share Host IPC Namespace',
        name: 'Host IPC',
        longName: 'Host IPC',
        jsonpath: 'fields.hostIPC',
        category: policyCriteriaCategories.DEPLOYMENT_METADATA,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Uses Host IPC namespace',
                value: true,
            },
            {
                text: 'Does Not Use Host IPC namespace',
                value: false,
            },
        ],
        required: false,
        default: false,
        defaultValue: true,
        disabled: true,
        canBooleanLogic: false,
    },
    {
        label: 'Drop Capabilities',
        name: 'Drop Capabilities',
        jsonpath: 'fields.dropCapabilities',
        category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
        type: 'select',
        options: [...capabilities],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Add Capabilities',
        name: 'Add Capabilities',
        jsonpath: 'fields.addCapabilities',
        category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
        type: 'select',
        options: [...capabilities],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Process Name',
        name: 'Process Name',
        negatedName: `Process name doesn't match`,
        jsonpath: 'fields.processPolicy.name',
        category: policyCriteriaCategories.PROCESS_ACTIVITY,
        type: 'text',
        placeholder: 'apt-get',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Process Ancestor',
        name: 'Process Ancestor',
        negatedName: `Process ancestor doesn't match`,
        jsonpath: 'fields.processPolicy.ancestor',
        category: policyCriteriaCategories.PROCESS_ACTIVITY,
        type: 'text',
        placeholder: 'java',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Process Args',
        name: 'Process Arguments',
        negatedName: `Process arguments don't match`,
        jsonpath: 'fields.processPolicy.args',
        category: policyCriteriaCategories.PROCESS_ACTIVITY,
        type: 'text',
        placeholder: 'install nmap',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Process UID',
        name: 'Process UID',
        negatedName: `Process UID doesn't match`,
        jsonpath: 'fields.processPolicy.uid',
        category: policyCriteriaCategories.PROCESS_ACTIVITY,
        type: 'text',
        placeholder: '0',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Port Exposure',
        name: 'Port Exposure Method',
        negatedName: 'Port Exposure Method is not',
        jsonpath: 'fields.portExposurePolicy.exposureLevels',
        category: policyCriteriaCategories.NETWORKING,
        type: 'select',
        options: Object.keys(portExposureLabels)
            .filter((key) => key !== 'INTERNAL')
            .map((key) => ({
                label: portExposureLabels[key],
                value: key,
            })),
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Writable Host Mount',
        name: 'Writable Host Mount',
        longName: 'Host Mount Writability',
        jsonpath: 'fields.hostMountPolicy.readOnly',
        category: policyCriteriaCategories.STORAGE,
        type: 'radioGroup',
        radioButtons: [
            {
                text: 'Writable',
                value: true,
            },
            {
                text: 'Read-only',
                value: false,
            },
        ],
        required: false,
        default: false,
        defaultValue: false,
        reverse: true,
        disabled: true,
        canBooleanLogic: false,
    },
    {
        label: 'Process Baselining Enabled',
        name: 'Unexpected Process Executed',
        longName: 'Process Baselining Status',
        jsonpath: 'fields.whitelistEnabled',
        category: policyCriteriaCategories.PROCESS_ACTIVITY,
        type: 'radioGroup',
        radioButtons: [
            { text: 'Unexpected Process', value: true },
            { text: 'Expected Process', value: false },
        ],
        required: false,
        default: false,
        defaultValue: false,
        reverse: false,
        canBooleanLogic: false,
    },
    {
        label: 'Service Account',
        name: 'Service Account',
        longName: 'Service Account Name',
        negatedName: `Service Account Name doesn't match`,
        category: policyCriteriaCategories.KUBERNETES_ACCESS,
        type: 'text',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Minimum RBAC Permissions',
        name: 'Minimum RBAC Permissions',
        longName: 'RBAC permission level is at least',
        negatedName: 'RBAC permission level is less than',
        jsonpath: 'fields.permissionPolicy.permissionLevel',
        category: policyCriteriaCategories.KUBERNETES_ACCESS,
        type: 'select',
        options: Object.keys(rbacPermissionLabels).map((key) => ({
            label: rbacPermissionLabels[key],
            value: key,
        })),
        required: false,
        default: false,
        canBooleanLogic: false,
    },
    {
        label: 'Required Image Label',
        name: 'Required Image Label',
        jsonpath: 'fields.requiredImageLabel',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.requiredImageLabel.key',
                type: 'text',
                label: 'Key',
                placeholder: 'requiredLabelKey.*',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.requiredImageLabel.value',
                type: 'text',
                label: 'Value',
                placeholder: 'requiredValue.*',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Disallowed Image Label',
        name: 'Disallowed Image Label',
        jsonpath: 'fields.disallowedImageLabel',
        category: policyCriteriaCategories.IMAGE_CONTENTS,
        type: 'group',
        jsonpaths: [
            {
                jsonpath: 'fields.disallowedImageLabel.key',
                type: 'text',
                label: 'Key',
                placeholder: 'disallowedLabelKey.*',
                subpath: 'key',
            },
            {
                jsonpath: 'fields.disallowedImageLabel.value',
                type: 'text',
                label: 'Value',
                placeholder: 'disallowedValue.*',
                subpath: 'value',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Namespace',
        name: 'Namespace',
        longName: 'Namespace',
        negatedName: `Namespace doesn't match`,
        category: policyCriteriaCategories.DEPLOYMENT_METADATA,
        type: 'text',
        placeholder: 'default',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
    {
        label: 'Container Name',
        name: 'Container Name',
        longName: 'Container Name',
        negatedName: `Container name doesn't match`,
        category: policyCriteriaCategories.CONTAINER_CONFIGURATION,
        type: 'text',
        placeholder: 'default',
        required: false,
        default: false,
        canBooleanLogic: true,
    },
];

const k8sEventsDescriptor = [
    {
        label: 'Kubernetes Action',
        name: 'Kubernetes Resource',
        longName: 'Kubernetes Action',
        shortName: 'Kubernetes Action',
        category: policyCriteriaCategories.KUBERNETES_EVENTS,
        type: 'select',
        options: [
            {
                label: 'Pod Exec',
                value: 'PODS_EXEC',
            },
            {
                label: 'Pods Port Forward',
                value: 'PODS_PORTFORWARD',
            },
        ],
        required: false,
        default: false,
        canBooleanLogic: false,
    },
];

const networkDetectionDescriptor = [
    {
        label: 'Network Baselining Enabled',
        name: 'Unexpected Network Flow Detected',
        longName: 'Network Baselining Status',
        jsonpath: 'fields.networkbaselineenabled',
        category: policyCriteriaCategories.NETWORKING,
        type: 'radioGroup',
        radioButtons: [
            { text: 'Unexpected Network Flow', value: true },
            { text: 'Expected Network Flow', value: false },
        ],
        required: false,
        default: false,
        defaultValue: false,
        reverse: false,
        canBooleanLogic: false,
    },
];

export const policyStatus = {
    header: 'Enable Policy',
    descriptor: policyStatusDescriptor,
    dataTestId: 'policyStatusField',
};

export const policyDetails = {
    header: 'Policy Summary',
    descriptor: policyDetailsFormDescriptor,
    dataTestId: 'policyDetailsFields',
};

export const policyConfiguration = {
    header: 'Policy Criteria',
    descriptor: policyConfigurationDescriptor,
    dataTestId: 'policyConfigurationFields',
};

let isFirstLoad = true;

export const getPolicyConfiguration = (featureFlags) => {
    const newPolicyConfiguration = { ...policyConfiguration };
    if (featureFlags[knownBackendFlags.ROX_K8S_EVENTS_DETECTION] && isFirstLoad) {
        newPolicyConfiguration.descriptor.push(...k8sEventsDescriptor);
    }
    if (featureFlags[knownBackendFlags.ROX_NETWORK_DETECTION_BASELINE_VIOLATION] && isFirstLoad) {
        newPolicyConfiguration.descriptor.push(...networkDetectionDescriptor);
    }
    isFirstLoad = false;
    return newPolicyConfiguration;
};
