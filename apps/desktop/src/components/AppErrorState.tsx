// SPDX-License-Identifier: MIT
import { useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import { AlertTriangle, RefreshCw, Copy, FileText } from "lucide-react";

interface AppErrorStateProps {
  title?: string;
  message: string;
  onRetry?: () => void;
}

interface DiagnosticsInfo {
  log_path: string;
  bootstrap_log_path?: string;
  startup_error: {
    message: string;
    exit_code: number | null;
  } | null;
}

export function AppErrorState({ title, message, onRetry }: AppErrorStateProps) {
  const [showDiag, setShowDiag] = useState(false);
  const [diag, setDiag] = useState<DiagnosticsInfo | null>(null);
  const [copied, setCopied] = useState(false);

  async function loadDiagnostics() {
    try {
      const info = await invoke<DiagnosticsInfo>("get_diagnostics");
      setDiag(info);
      setShowDiag(true);
    } catch {
      setShowDiag(true);
    }
  }

  async function copyDetails() {
    const lines = [
      `FileENIAC v0.1.7 - Diagnostico`,
      ``,
      `Erro: ${message}`,
    ];
    if (diag?.startup_error) {
      lines.push(`Detalhe: ${diag.startup_error.message}`);
      if (diag.startup_error.exit_code != null) {
        lines.push(`Codigo de saida: ${diag.startup_error.exit_code}`);
      }
    }
    if (diag?.log_path) {
      lines.push(`Log: ${diag.log_path}`);
    }
    if (diag?.bootstrap_log_path) {
      lines.push(`Bootstrap: ${diag.bootstrap_log_path}`);
    }
    try {
      await navigator.clipboard.writeText(lines.join("\n"));
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {}
  }

  return (
    <div className="fixed inset-0 bg-eniac-950 flex items-center justify-center">
      <div className="text-center max-w-md px-8">
        <div className="mx-auto mb-6 w-16 h-16 rounded-2xl bg-red-900/30 border border-red-700/30 flex items-center justify-center">
          <AlertTriangle className="h-8 w-8 text-red-400" />
        </div>
        <h2 className="text-lg font-semibold text-white mb-2">
          {title || "Nao foi possivel iniciar"}
        </h2>
        <p className="text-eniac-300 text-sm mb-6 leading-relaxed">
          {message}
        </p>
        <div className="flex flex-col gap-3">
          {onRetry && (
            <button
              onClick={onRetry}
              className="w-full py-3 px-6 bg-eniac-600 text-white rounded-lg font-semibold hover:bg-eniac-700 transition-colors flex items-center justify-center gap-2"
            >
              <RefreshCw className="h-4 w-4" />
              Tentar novamente
            </button>
          )}
          <button
            onClick={loadDiagnostics}
            className="w-full py-3 px-6 bg-eniac-900/60 border border-eniac-700/40 text-eniac-200 rounded-lg font-medium hover:bg-eniac-800/60 transition-colors flex items-center justify-center gap-2"
          >
            <FileText className="h-4 w-4" />
            Abrir diagnostico
          </button>
        </div>
        {showDiag && (
          <div className="mt-6 text-left bg-eniac-900/60 border border-eniac-700/30 rounded-lg p-4 max-h-48 overflow-y-auto">
            <p className="text-xs text-eniac-400 mb-2">Detalhes tecnicos:</p>
            {diag?.startup_error && (
              <p className="text-xs text-red-300 mb-1">{diag.startup_error.message}</p>
            )}
            {diag?.startup_error?.exit_code != null && (
              <p className="text-xs text-eniac-400 mb-1">Codigo: {diag.startup_error.exit_code}</p>
            )}
            {diag?.log_path && (
              <p className="text-xs text-eniac-500">Log: {diag.log_path}</p>
            )}
            {diag?.bootstrap_log_path && (
              <p className="text-xs text-eniac-500">Bootstrap: {diag.bootstrap_log_path}</p>
            )}
            {!diag && (
              <p className="text-xs text-eniac-500">Diagnostico local indisponivel.</p>
            )}
            <button
              onClick={copyDetails}
              className="mt-3 text-xs text-eniac-300 hover:text-white flex items-center gap-1 transition-colors"
            >
              <Copy className="h-3 w-3" />
              {copied ? "Copiado!" : "Copiar detalhes"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
