import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import ReactRouterPropTypes from 'react-router-prop-types';
import URLService from 'modules/URLService';

import CollapsibleBanner from 'Components/CollapsibleBanner/CollapsibleBanner';
import ComplianceAcrossEntities from 'Containers/Compliance/widgets/ComplianceAcrossEntities';
import ControlsMostFailed from 'Containers/Compliance/widgets/ControlsMostFailed';
import SearchInput from './SearchInput';
import Header from './Header';
import ListTable from './Table';
import SidePanel from './SidePanel';

class ComplianceListPage extends Component {
    static propTypes = {
        match: ReactRouterPropTypes.match.isRequired,
        location: ReactRouterPropTypes.location.isRequired,
        params: PropTypes.shape({
            entityType: PropTypes.string.isRequired
        })
    };

    static defaultProps = {
        params: null
    };

    constructor(props) {
        super(props);
        this.state = {
            selectedRow: null
        };
    }

    updateSelectedRow = selectedRow => this.setState({ selectedRow });

    clearSelectedRow = () => {
        this.setState({ selectedRow: null });
    };

    render() {
        const { match, location } = this.props;
        const { selectedRow } = this.state;
        const params = URLService.getParams(match, location);
        return (
            <section className="flex flex-col h-full relative" id="capture-list">
                <Header searchComponent={<SearchInput categories={['COMPLIANCE']} />} />
                <CollapsibleBanner className="pdf-page">
                    <ComplianceAcrossEntities params={params} />
                    <ControlsMostFailed params={params} showEmpty />
                </CollapsibleBanner>
                <div className="flex flex-1 overflow-y-auto">
                    <ListTable
                        selectedRow={selectedRow}
                        params={params}
                        updateSelectedRow={this.updateSelectedRow}
                        pdfId="capture-list"
                    />
                    {selectedRow && (
                        <SidePanel
                            match={match}
                            location={location}
                            selectedRow={selectedRow}
                            clearSelectedRow={this.clearSelectedRow}
                        />
                    )}
                </div>
            </section>
        );
    }
}

export default withRouter(ComplianceListPage);
