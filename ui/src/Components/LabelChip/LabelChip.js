import React from 'react';
import PropTypes from 'prop-types';
import { colorTypes, defaultColorType } from 'constants/visuals/colors';

const LabelChip = ({ text, type }) => {
    let className = 'border px-2 py-1 rounded font-600 text-base';
    const colorType = colorTypes.find(datum => datum === type) || defaultColorType;
    className = `${className} bg-${colorType}-200 border-${colorType}-400 text-${colorType}-800`;
    return <span className={className}>{text}</span>;
};

LabelChip.propTypes = {
    text: PropTypes.string.isRequired,
    type: PropTypes.oneOf(colorTypes)
};

LabelChip.defaultProps = {
    type: defaultColorType
};

export default LabelChip;
