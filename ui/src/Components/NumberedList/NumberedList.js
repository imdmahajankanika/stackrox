import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import Tooltip from 'rc-tooltip';

import TooltipOverlay from 'Components/TooltipOverlay';

const leftSideClasses = 'text-sm text-primary-800 font-600 truncate';

const NumberedList = ({ data, linkLeftOnly }) => {
    // eslint-disable-next-line no-unused-vars
    const list = data.map(({ text, subText, url, component, tooltip }, i) => {
        const className = `flex items-center py-2 px-2 ${
            i !== 0 ? 'border-t border-base-300' : ''
        } ${url ? 'hover:bg-base-200' : ''}`;
        let leftSide = (
            <>
                {i + 1}.&nbsp;{text}&nbsp;
                {subText && <span className="text-base-500 text-xs">{subText}</span>}
            </>
        );
        if (url && linkLeftOnly) {
            leftSide = (
                <Link
                    data-testid="numbered-list-item-name"
                    className={`${leftSideClasses} no-underline w-full`}
                    to={url}
                >
                    {leftSide}
                </Link>
            );
        } else {
            leftSide = (
                <span data-testid="numbered-list-item-name" className={leftSideClasses}>
                    {leftSide}
                </span>
            );
        }
        let content = (
            <>
                {leftSide}
                <div className="flex flex-1 justify-end ml-4 whitespace-no-wrap">{component}</div>
            </>
        );
        if (url && !linkLeftOnly) {
            content = (
                <Link className="flex items-center no-underline w-full" to={url}>
                    {content}
                </Link>
            );
        }
        return (
            <>
                {tooltip && (
                    <li key={text} className={className}>
                        <Tooltip
                            placement="top"
                            overlay={
                                <TooltipOverlay
                                    title={tooltip.title}
                                    body={tooltip.body}
                                    footer={tooltip.footer}
                                />
                            }
                            mouseLeaveDelay={0}
                            overlayClassName="opacity-100 max-w-1/3"
                        >
                            {content}
                        </Tooltip>
                    </li>
                )}
                {!tooltip && (
                    <li key={text} className={className}>
                        {content}
                    </li>
                )}
            </>
        );
    });
    return <ul className="list-reset leading-loose">{list}</ul>;
};

NumberedList.propTypes = {
    data: PropTypes.arrayOf(
        PropTypes.shape({
            text: PropTypes.string.isRequired,
            subText: PropTypes.string,
            components: PropTypes.element,
            url: PropTypes.string
        })
    ),
    linkLeftOnly: PropTypes.bool
};

NumberedList.defaultProps = {
    data: [],
    linkLeftOnly: false
};

export default NumberedList;
