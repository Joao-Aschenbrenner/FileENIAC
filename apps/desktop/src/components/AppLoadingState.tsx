// SPDX-License-Identifier: MIT
export function AppLoadingState({ message }: { message?: string }) {
  return (
    <div className="fixed inset-0 bg-eniac-950 flex items-center justify-center">
      <div className="text-center max-w-sm px-8">
        <div className="mx-auto mb-6 w-16 h-16 rounded-2xl bg-eniac-900/60 border border-eniac-700/30 flex items-center justify-center">
          <svg className="animate-spin h-8 w-8 text-eniac-400" viewBox="0 0 24 24" fill="none">
            <circle className="opacity-20" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" />
            <path className="opacity-80" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        </div>
        <h2 className="text-lg font-semibold text-eniac-100 mb-1">FileENIAC</h2>
        <p className="text-eniac-300 text-sm">{message || "Configurando ambiente..."}</p>
      </div>
    </div>
  );
}
