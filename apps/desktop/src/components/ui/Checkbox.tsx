// SPDX-License-Identifier: MIT
import { InputHTMLAttributes, forwardRef } from "react";

interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  label?: string;
}

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  ({ label, className = "", id, ...props }, ref) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, "-");

    return (
      <div className="flex items-center gap-2">
        <input
          ref={ref}
          type="checkbox"
          id={inputId}
          className={`w-4 h-4 rounded border-gray-300 text-eniac-600 focus:ring-eniac-500 cursor-pointer ${className}`}
          {...props}
        />
        {label && (
          <label htmlFor={inputId} className="text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            {label}
          </label>
        )}
      </div>
    );
  }
);

Checkbox.displayName = "Checkbox";