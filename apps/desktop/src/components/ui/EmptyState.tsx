// SPDX-License-Identifier: MIT
interface EmptyStateProps {
  title: string;
  description?: string;
  action?: { label: string; onClick: () => void };
}

export function EmptyState({ title, description, action }: EmptyStateProps) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-8 text-center">
      <div className="text-3xl mb-3">📭</div>
      <h3 className="font-semibold text-gray-700 mb-1">{title}</h3>
      {description && <p className="text-sm text-gray-500 mb-4">{description}</p>}
      {action && (
        <button
          onClick={action.onClick}
          className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors"
        >
          {action.label}
        </button>
      )}
    </div>
  );
}
