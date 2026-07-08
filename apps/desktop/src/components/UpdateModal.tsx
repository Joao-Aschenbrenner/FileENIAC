// SPDX-License-Identifier: MIT
import type { UpdateState } from "../hooks/useUpdateCheck";

type Props = {
  state: UpdateState;
  onDismiss: () => void;
  onInstall: () => void;
};

export default function UpdateModal({ state, onDismiss, onInstall }: Props) {
  if (state.status === "idle" || state.status === "checking") return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onDismiss} />
      <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
        <h3 className="font-semibold text-gray-800 mb-2">Atualizacao</h3>

        {state.status === "available" && (
          <>
            <p className="text-sm text-gray-600 mb-1">
              Nova versao <strong>{state.version}</strong> disponivel.
            </p>
            {state.body && (
              <p className="text-xs text-gray-500 mb-4 whitespace-pre-wrap max-h-24 overflow-y-auto">
                {state.body}
              </p>
            )}
            <div className="flex gap-3 mt-4">
              <button onClick={onDismiss} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">
                Mais tarde
              </button>
              <button onClick={onInstall} className="flex-1 py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700">
                Atualizar agora
              </button>
            </div>
          </>
        )}

        {state.status === "downloading" && (
          <div className="text-center py-4">
            <p className="text-sm text-gray-500">Baixando atualizacao... isso pode levar alguns minutos.</p>
          </div>
        )}

        {state.status === "error" && (
          <>
            <p className="text-sm text-red-600 mb-4">{state.message}</p>
            <button onClick={onDismiss} className="w-full py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">
              Fechar
            </button>
          </>
        )}
      </div>
    </div>
  );
}
