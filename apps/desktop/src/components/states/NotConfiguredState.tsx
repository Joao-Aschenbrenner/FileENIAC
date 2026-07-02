import { useNavigate } from "react-router-dom";

interface NotConfiguredStateProps {
  areaName: string;
  configPath?: string;
}

export function NotConfiguredState({ areaName, configPath = "/bootstrap" }: NotConfiguredStateProps) {
  const navigate = useNavigate();
  return (
    <div className="bg-eniac-900/20 border border-eniac-700/20 rounded-xl p-8 text-center">
      <div className="mx-auto mb-4 w-12 h-12 rounded-xl bg-eniac-800/50 border border-eniac-700/30 flex items-center justify-center">
        <span className="text-eniac-300 text-xl">?</span>
      </div>
      <h3 className="font-semibold text-eniac-200 mb-2">{areaName}</h3>
      <p className="text-sm text-eniac-400 mb-4">
        Você ainda não configurou esta área.
      </p>
      <button
        onClick={() => navigate(configPath)}
        className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors"
      >
        Configurar agora
      </button>
    </div>
  );
}

export function NotConfiguredInline({ areaName, onConfigure }: { areaName: string; onConfigure: () => void }) {
  return (
    <div className="bg-eniac-900/20 border border-eniac-700/20 rounded-lg p-5 text-center">
      <p className="text-sm text-eniac-400 mb-3">Você ainda não configurou {areaName}.</p>
      <button
        onClick={onConfigure}
        className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors"
      >
        Configurar agora
      </button>
    </div>
  );
}
