// SPDX-License-Identifier: MIT
import { InputHTMLAttributes, forwardRef } from "react";

export type InputSizeVariant = "sm" | "md" | "lg";

interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  label?: string;
  error?: string;
  hint?: string;
  inputSize?: InputSizeVariant;
  icon?: React.ReactNode;
  iconPosition?: "left" | "right";
}

const sizeClasses: Record<InputSizeVariant, string> = {
  sm: "px-3 py-1.5 text-sm",
  md: "px-3 py-2 text-sm",
  lg: "px-4 py-3 text-base",
};

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, inputSize = "md", icon, iconPosition = "left", className = "", id, ...props }, ref) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, "-");

    return (
      <div className="w-full">
        {label && (
          <label htmlFor={inputId} className="block text-sm font-medium text-gray-700 mb-1">
            {label}
          </label>
        )}
        <div className="relative">
          {icon && iconPosition === "left" && (
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
              {icon}
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={`
              w-full rounded-lg border transition-colors duration-150
              focus:outline-none focus:ring-2 focus:ring-eniac-500 focus:border-transparent
              disabled:bg-gray-50 disabled:cursor-not-allowed
              ${error ? "border-red-300 focus:ring-red-500" : "border-gray-300"}
              ${icon && iconPosition === "left" ? "pl-10" : ""}
              ${icon && iconPosition === "right" ? "pr-10" : ""}
              ${sizeClasses[inputSize]}
              ${className}
            `}
            aria-invalid={error ? "true" : "false"}
            aria-describedby={error ? `${inputId}-error` : hint ? `${inputId}-hint` : undefined}
            {...props}
          />
          {icon && iconPosition === "right" && (
            <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-gray-400">
              {icon}
            </div>
          )}
        </div>
        {error && (
          <p id={`${inputId}-error`} className="mt-1 text-sm text-red-600" role="alert">
            {error}
          </p>
        )}
        {hint && !error && (
          <p id={`${inputId}-hint`} className="mt-1 text-sm text-gray-500">
            {hint}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = "Input";