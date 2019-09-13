import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import entityTypes from 'constants/entityTypes';
import URLService from 'modules/URLService';
import { resourceLabels } from 'messages/common';
import capitalize from 'lodash/capitalize';
import Widget from 'Components/Widget';
import Query from 'Components/ThrowingQuery';
import Loader from 'Components/Loader';
import CountWidget from 'Components/CountWidget';
import { SEARCH_WITH_CONTROLS as QUERY } from 'queries/search';
import queryService from 'modules/queryService';
import { getResourceCountFromAggregatedResults } from 'modules/complianceUtils';
import ReactRouterPropTypes from 'react-router-prop-types';
import { withRouter } from 'react-router-dom';
import appContexts from 'constants/appContextTypes';
import searchContext from 'Containers/searchContext';

const ResourceCount = ({
    match,
    location,
    entityType,
    relatedToResourceType,
    relatedToResource,
    count
}) => {
    const searchParam = useContext(searchContext);
    function getUrl() {
        if (entityType === entityTypes.SECRET) {
            return URLService.getURL(match, location)
                .set('context', appContexts.SECRET)
                .query({
                    [searchParam]: {
                        [`${capitalize(relatedToResourceType)}`]: relatedToResource.name
                    }
                })
                .url();
        }

        return URLService.getURL(match, location)
            .base(relatedToResourceType, relatedToResource.id)
            .push(entityType)
            .url();
    }

    function getVariables() {
        let query;
        switch (relatedToResourceType) {
            case entityTypes.NAMESPACE:
                query = {
                    namespace: relatedToResource.name,
                    cluster: relatedToResource.clusterName
                };
                break;
            default:
                query = { [`${capitalize(relatedToResourceType)} ID`]: relatedToResource.id };
        }

        return {
            query: queryService.objectToWhereClause(query),
            categories: []
        };
    }

    const variables = getVariables();
    const headerText = `${resourceLabels[entityType]} Count`;
    const url = getUrl();

    if (count || count === 0) {
        return (
            <Widget header={headerText} bodyClassName="p-2">
                <CountWidget title={headerText} count={count} linkUrl={url} />;
            </Widget>
        );
    }
    return (
        <Query query={QUERY} variables={variables}>
            {({ loading, data }) => {
                const contents = <Loader />;
                if (!loading && data) {
                    const resourceCount = getResourceCountFromAggregatedResults(entityType, data);

                    return <CountWidget title={headerText} count={resourceCount} linkUrl={url} />;
                }
                return (
                    <Widget header={headerText} bodyClassName="p-2">
                        {contents}
                    </Widget>
                );
            }}
        </Query>
    );
};

ResourceCount.propTypes = {
    match: ReactRouterPropTypes.match.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    entityType: PropTypes.string,
    relatedToResourceType: PropTypes.string.isRequired,
    relatedToResource: PropTypes.shape({}),
    count: PropTypes.number
};

ResourceCount.defaultProps = {
    entityType: null,
    relatedToResource: null,
    count: null
};

export default withRouter(ResourceCount);
