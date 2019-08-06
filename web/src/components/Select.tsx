import DropdownIcon from 'mdi-react/ChevronDownIcon'
import * as React from 'react'

/**
 * A drop-in replacement for native `<select>` elements that ensures proper cross-browser styling.
 */
export const Select: React.FunctionComponent<
    React.DetailedHTMLProps<React.SelectHTMLAttributes<HTMLSelectElement>, HTMLSelectElement>
> = ({ className = '', ...props }) => (
    <div className="select">
        {/* this is the ONLY allowed instance of <select> */}
        {/* eslint-disable-next-line react/forbid-elements */}
        <select {...props} className={`select__picker form-control ${className}`}>
            {props.children}
        </select>
        <DropdownIcon className="select__icon" />
    </div>
)
