/* eslint-disable no-nested-ternary */
/* eslint-disable react/jsx-no-bind */
import React, { ReactElement, useEffect, useState } from 'react';
import { useHistory, useLocation, useParams } from 'react-router-dom';
import {
    Alert,
    AlertActionCloseButton,
    AlertVariant,
    Badge,
    Bullseye,
    Button,
    Spinner,
    Title,
    Toolbar,
    ToolbarContent,
    ToolbarGroup,
    ToolbarItem,
} from '@patternfly/react-core';

import {
    AccessScope,
    Role,
    createAccessScope,
    deleteAccessScope,
    fetchAccessScopes,
    fetchRolesAsArray,
    updateAccessScope,
} from 'services/RolesService';

import AccessControlNav from '../AccessControlNav';
import AccessControlPageTitle from '../AccessControlPageTitle';
import { getEntityPath, getQueryObject } from '../accessControlPaths';

import AccessScopeForm from './AccessScopeForm';
import AccessScopesList from './AccessScopesList';

import './AccessScopes.css';

const accessScopeNew: AccessScope = {
    id: '',
    name: '',
    description: '',
    rules: {
        includedClusters: [],
        includedNamespaces: [],
        clusterLabelSelectors: [],
        namespaceLabelSelectors: [],
    },
};

const entityType = 'ACCESS_SCOPE';

function AccessScopes(): ReactElement {
    const history = useHistory();
    const { search } = useLocation();
    const queryObject = getQueryObject(search);
    const { action } = queryObject;
    const { entityId } = useParams();

    const [isFetching, setIsFetching] = useState(false);
    const [accessScopes, setAccessScopes] = useState<AccessScope[]>([]);
    const [alertAccessScopes, setAlertAccessScopes] = useState<ReactElement | null>(null);
    const [roles, setRoles] = useState<Role[]>([]);
    const [alertRoles, setAlertRoles] = useState<ReactElement | null>(null);

    useEffect(() => {
        // The primary request has fetching spinner and unclosable alert.
        setIsFetching(true);
        setAlertAccessScopes(null);
        fetchAccessScopes()
            .then((accessScopesFetched) => {
                setAccessScopes(accessScopesFetched);
            })
            .catch((error) => {
                setAlertAccessScopes(
                    <Alert
                        title="Fetch access scopes failed"
                        variant={AlertVariant.danger}
                        isInline
                    >
                        {error.message}
                    </Alert>
                );
            })
            .finally(() => {
                setIsFetching(false);
            });

        // TODO Until secondary requests succeed, disable Create and Edit because selections might be incomplete?
        setAlertRoles(null);
        fetchRolesAsArray()
            .then((rolesFetched) => {
                setRoles(rolesFetched);
            })
            .catch((error) => {
                const actionClose = <AlertActionCloseButton onClose={() => setAlertRoles(null)} />;
                setAlertRoles(
                    <Alert
                        title="Fetch roles failed"
                        variant={AlertVariant.warning}
                        isInline
                        actionClose={actionClose}
                    >
                        {error.message}
                    </Alert>
                );
            });
    }, []);

    function onClickCreate() {
        history.push(getEntityPath(entityType, undefined, { action: 'create' }));
    }

    function handleDelete(idDelete: string) {
        return deleteAccessScope(idDelete).then(() => {
            // Remove the deleted entity.
            setAccessScopes(accessScopes.filter(({ id }) => id !== idDelete));
        }); // TODO catch error display alert
    }

    function handleEdit() {
        history.push(getEntityPath(entityType, entityId, { action: 'update' }));
    }

    function handleCancel() {
        // Go back from action=create to list or go back from action=update to entity.
        history.goBack();
    }

    function handleSubmit(values: AccessScope): Promise<null> {
        return action === 'create'
            ? createAccessScope(values).then((entityCreated) => {
                  // Append the created entity.
                  setAccessScopes([...accessScopes, entityCreated]);

                  // Replace path which had action=create with plain entity path.
                  history.replace(getEntityPath(entityType, entityCreated.id));

                  return null; // because the form has only catch and finally
              })
            : updateAccessScope(values).then(() => {
                  // Replace the updated entity with values because response is empty object.
                  setAccessScopes(
                      accessScopes.map((entity) => (entity.id === values.id ? values : entity))
                  );

                  // Replace path which had action=update with plain entity path.
                  history.replace(getEntityPath(entityType, entityId));

                  return null; // because the form has only catch and finally
              });
    }

    const accessScope = accessScopes.find(({ id }) => id === entityId) || accessScopeNew;
    const isActionable = true; // TODO does it depend on user role?
    const hasAction = Boolean(action);
    const isEntity = hasAction || Boolean(entityId);

    return (
        <>
            <AccessControlPageTitle entityType={entityType} isEntity={isEntity} />
            <AccessControlNav entityType={entityType} />
            {alertAccessScopes}
            {alertRoles}
            {isFetching ? (
                <Bullseye>
                    <Spinner />
                </Bullseye>
            ) : isEntity ? (
                <AccessScopeForm
                    isActionable={isActionable}
                    action={action}
                    accessScope={accessScope}
                    accessScopes={accessScopes}
                    handleCancel={handleCancel}
                    handleEdit={handleEdit}
                    handleSubmit={handleSubmit}
                />
            ) : (
                <>
                    <Toolbar inset={{ default: 'insetNone' }}>
                        <ToolbarContent>
                            <ToolbarGroup spaceItems={{ default: 'spaceItemsMd' }}>
                                <ToolbarItem>
                                    <Title headingLevel="h2">Access scopes</Title>
                                </ToolbarItem>
                                <ToolbarItem>
                                    <Badge isRead>{accessScopes.length}</Badge>
                                </ToolbarItem>
                            </ToolbarGroup>
                            <ToolbarItem alignment={{ default: 'alignRight' }}>
                                <Button
                                    variant="primary"
                                    onClick={onClickCreate}
                                    isDisabled={isFetching}
                                    isSmall
                                >
                                    Create access scope
                                </Button>
                            </ToolbarItem>
                        </ToolbarContent>
                    </Toolbar>
                    {accessScopes.length !== 0 && (
                        <AccessScopesList
                            entityId={entityId}
                            accessScopes={accessScopes}
                            roles={roles}
                            handleDelete={handleDelete}
                        />
                    )}
                </>
            )}
        </>
    );
}

export default AccessScopes;
