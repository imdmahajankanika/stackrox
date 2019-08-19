import React, { Component } from 'react';
import PropTypes from 'prop-types';
import * as Icon from 'react-feather';

class Select extends Component {
    static propTypes = {
        options: PropTypes.arrayOf(PropTypes.object).isRequired,
        onChange: PropTypes.func.isRequired,
        placeholder: PropTypes.string,
        className: PropTypes.string,
        value: PropTypes.oneOfType([PropTypes.string, PropTypes.number])
    };

    static defaultProps = {
        placeholder: '',
        className:
            'block w-full border bg-base-200 border-base-400 text-base-600 p-3 pr-8 rounded-sm z-1 focus:border-base-500',
        value: ''
    };

    onClick = event => {
        const selectedOption = this.props.options.find(
            option => option.value === event.target.value
        );
        if (!selectedOption) {
            throw new Error('Selected ID does not match any known option in Select control.');
        }

        this.props.onChange(selectedOption);
    };

    render() {
        const { className, options, placeholder, value } = this.props;
        return (
            <div className="flex relative">
                <select
                    className={`${className} cursor-pointer`}
                    onChange={this.onClick}
                    value={value}
                    aria-label={placeholder}
                >
                    {placeholder && (
                        <option value="" disabled>
                            {placeholder}
                        </option>
                    )}
                    {options.map(option => (
                        <option key={option.label} value={option.value}>
                            {option.label}
                        </option>
                    ))}
                </select>
                <div className="flex items-center px-2 cursor-pointer z-0 pointer-events-none">
                    <Icon.ChevronDown className="h-4 w-4" />
                </div>
            </div>
        );
    }
}

export default Select;
