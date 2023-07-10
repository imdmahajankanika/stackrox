import React, { useMemo, useState, ReactElement, useCallback, useEffect } from 'react';
import {
    Button,
    ButtonVariant,
    Flex,
    FlexItem,
    Select,
    SelectOption,
    SelectProps,
    SelectVariant,
    ValidatedOptions,
} from '@patternfly/react-core';
import sortBy from 'lodash/sortBy';
import uniqBy from 'lodash/uniqBy';

import FormLabelGroup from 'Components/PatternFly/FormLabelGroup';
import { Collection, listCollections } from 'services/CollectionsService';
import CollectionsFormModal, {
    CollectionFormModalAction,
} from 'Containers/Collections/CollectionFormModal';
import { useCollectionFormSubmission } from 'Containers/Collections/hooks/useCollectionFormSubmission';
import useSelectToggle from 'hooks/patternfly/useSelectToggle';
import { usePaginatedQuery } from 'hooks/usePaginatedQuery';
import { ReportScope } from 'hooks/useFetchReport';

const COLLECTION_PAGE_SIZE = 10;

type CollectionSelectionProps = {
    selectedScopeId: string;
    initialReportScope: ReportScope | null;
    onChange: (selection: string) => void;
    allowCreate: boolean;
};

function CollectionSelection({
    selectedScopeId,
    initialReportScope,
    onChange,
    allowCreate,
}: CollectionSelectionProps): ReactElement {
    const { isOpen, onToggle } = useSelectToggle();
    const [modalAction, setModalAction] = useState<CollectionFormModalAction>({ type: 'create' });
    const [isCollectionModalOpen, setIsCollectionModalOpen] = useState(false);

    const { configError, setConfigError, onSubmit } = useCollectionFormSubmission(modalAction);
    const [search, setSearch] = useState('');

    const requestFn = useCallback(
        (page: number) => {
            return listCollections(
                { 'Collection Name': search },
                { field: 'Collection Name', reversed: false },
                page,
                COLLECTION_PAGE_SIZE
            ).request;
        },
        [search]
    );

    const { data, isEndOfResults, isFetchingNextPage, fetchNextPage } = usePaginatedQuery(
        requestFn,
        COLLECTION_PAGE_SIZE
    );

    const isLegacyReportScopeSelected =
        initialReportScope?.type === 'AccessControlScope' &&
        initialReportScope?.id === selectedScopeId;

    // Combines the server-side fetched pages of collections data with the local cache
    // of created collections to create a flattened array sorted by name. This is intended to keep
    // the collection dropdown up to date with any collections that the user creates while in the form.
    //
    // Previously this was not needed since the component would refetch _all_ access scopes
    // upon creation of a new access scope, but we cannot do that efficiently since the collection dropdown
    // is paginated.
    //
    // This functionality can likely be removed if we move to a library based method of data fetching.
    const [createdCollections, setCreatedCollections] = useState<Collection[]>([]);
    const sortedCollections = useMemo(() => {
        const availableScopes: Pick<Collection, 'id' | 'name' | 'description'>[] = [
            ...data.flat(),
            ...createdCollections,
        ];
        // Adding the initial report scope, if available, allows the collection name to be displayed even
        // if it has not yet been fetched via the dropdown's pagination.
        if (initialReportScope && initialReportScope.type === 'CollectionScope') {
            availableScopes.push(initialReportScope);
        }

        // This is inefficient due to the multiple loops and the fact that we are already tracking
        // uniqueness for the _server side_ values, but need to do it twice to handle possible client
        // side values. However, 'N' should be small here and we are memoizing the result.
        const sorted = sortBy(availableScopes, ({ name }) => name.toLowerCase());
        return uniqBy(sorted, 'id');
    }, [data, createdCollections, initialReportScope]);

    // This makes sure that if a collection was deleted then we clear the scopeId
    useEffect(() => {
        const selectedCollection = sortedCollections.find(
            (collection) => collection.id === selectedScopeId
        );
        if (!selectedCollection) {
            onChange('');
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [data]);

    function onOpenViewCollectionModal() {
        setModalAction({ type: 'view', collectionId: selectedScopeId });
        setIsCollectionModalOpen((current) => !current);
    }

    function onOpenCreateCollectionModal() {
        setModalAction({ type: 'create' });
        setIsCollectionModalOpen((current) => !current);
    }

    function onScopeChange(_id, selection) {
        onToggle(false);
        onChange(selection);
    }

    let selectLoadingVariant: SelectProps['loadingVariant'];

    if (isFetchingNextPage) {
        selectLoadingVariant = 'spinner';
    } else if (!isEndOfResults) {
        selectLoadingVariant = {
            text: 'View more',
            onClick: () => fetchNextPage(),
        };
    }

    return (
        <>
            <FormLabelGroup
                isRequired
                label="Configure report scope"
                fieldId="scopeId"
                touched={isLegacyReportScopeSelected ? { scopeId: true } : {}}
                errors={
                    isLegacyReportScopeSelected
                        ? { scopeId: 'Choose a new collection to use as the report scope' }
                        : {}
                }
            >
                <Flex
                    direction={{ default: 'row' }}
                    spaceItems={{ default: 'spaceItemsNone' }}
                    alignItems={{ default: 'alignItemsFlexEnd' }}
                >
                    <FlexItem>
                        <Select
                            id="scopeId"
                            onSelect={onScopeChange}
                            selections={isLegacyReportScopeSelected ? '' : selectedScopeId}
                            placeholderText="Select a collection"
                            variant={SelectVariant.typeahead}
                            isOpen={isOpen}
                            onToggle={onToggle}
                            onTypeaheadInputChanged={setSearch}
                            loadingVariant={selectLoadingVariant}
                            onBlur={() => setSearch('')}
                            style={{
                                maxHeight: '275px',
                                overflowY: 'auto',
                            }}
                            validated={
                                isLegacyReportScopeSelected
                                    ? ValidatedOptions.error
                                    : ValidatedOptions.default
                            }
                        >
                            {sortedCollections.map(({ id, name, description }) => (
                                <SelectOption key={id} value={id} description={description}>
                                    {name}
                                </SelectOption>
                            ))}
                        </Select>
                    </FlexItem>
                    <FlexItem spacer={{ default: 'spacerMd' }}>
                        <Button
                            variant={ButtonVariant.tertiary}
                            onClick={onOpenViewCollectionModal}
                            isDisabled={selectedScopeId === ''}
                        >
                            View
                        </Button>
                    </FlexItem>
                    {allowCreate && (
                        <FlexItem>
                            <Button
                                variant={ButtonVariant.secondary}
                                onClick={onOpenCreateCollectionModal}
                            >
                                Create collection
                            </Button>
                        </FlexItem>
                    )}
                </Flex>
            </FormLabelGroup>
            {isCollectionModalOpen && (
                <CollectionsFormModal
                    hasWriteAccessForCollections={allowCreate}
                    modalAction={modalAction}
                    onClose={() => setIsCollectionModalOpen(false)}
                    configError={configError}
                    setConfigError={setConfigError}
                    onSubmit={(collection) =>
                        onSubmit(collection).then((collectionResponse) => {
                            onChange(collectionResponse.id);
                            setIsCollectionModalOpen(false);
                            setCreatedCollections((oldCollections) => [
                                ...oldCollections,
                                collectionResponse,
                            ]);
                        })
                    }
                />
            )}
        </>
    );
}

export default CollectionSelection;
