// SPDX-License-Identifier: MIT
interface ErrorStateProps {
  message: string;
  onRetry?: () => void;
}

export function ErrorState({ message, onRetry }: ErrorStateProps) {
  return (
    <div className="bg-red-900/20 border border-red-700/30 rounded-xl p-5 text-center">
      <div className="text-3xl mb-2">⚠️</div>
      <p className="text-sm text-red-300 mb-3">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="px-4 py-2 bg-red-600 text-white text-sm rounded-lg font-medium hover:bg-red-700 transition-colors"
        >
          Tentar novamente
        </button>
      )}
    </div>
  );
}
