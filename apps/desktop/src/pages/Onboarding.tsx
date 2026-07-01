// SPDX-License-Identifier: MIT
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { invoke } from "@tauri-apps/api/core";
import { getWorkspace, checkHealth, configureApiClientFromBackendInfo } from "../api/client";
import type { BackendInfo } from "../auth/tokenStorage";
import { open } from "@tauri-apps/plugin-dialog";
import { AppLoadingState } from "../components/AppLoadingState";
import { AppErrorState } from "../components/AppErrorState";

async function pickFolder(): Promise<string | null> {
  try {
    const selected = await open({ directory: true, multiple: false, title: "Selecionar Workspace" });
    return selected || null;
  } catch {
    return null;
  }
}

type StartupStatus = "waiting" | "ready" | "failed";

export default function Onboarding() {
  const navigate = useNavigate();
  const [step, setStep] = useState<"welcome" | "config" | "ready">("welcome");
  const [wsPath, setWsPath] = useState("");
  const [error, setError] = useState("");
  const [wsInfo, setWsInfo] = useState<any>(null);
  const [connecting, setConnecting] = useState(false);
  const [startupStatus, setStartupStatus] = useState<StartupStatus>("waiting");
  const [startupError, setStartupError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function waitForBackend() {
      try {
        const info = await invoke<BackendInfo>("get_backend_info");

        if (cancelled) return;

        if (!configureApiClientFromBackendInfo(info)) {
          const diag = await invoke<{ startup_error: { message: string } | null }>("get_diagnostics").catch(() => null);
          const msg = diag?.startup_error?.message || "Nao foi possivel configurar o ambiente automaticamente.";
          setStartupStatus("failed");
          setStartupError(msg);
          return;
        }

        setStartupStatus("ready");

        let retries = 0;
        while (retries < 15) {
          if (cancelled) return;
          const ok = await checkHealth();
          if (ok) return;
          await new Promise((r) => setTimeout(r, 1000));
          retries++;
        }

        setStartupStatus("failed");
        setStartupError("O servico nao respondeu a tempo. Tente novamente.");
      } catch (err) {
        console.error("FileENIAC startup check failed", err);
        if (!cancelled) {
          setStartupStatus("failed");
          setStartupError("Erro ao verificar o ambiente.");
        }
      }
    }

    waitForBackend();
    return () => { cancelled = true; };
  }, []);

  async function handleRetryStartup() {
    setStartupStatus("waiting");
    setStartupError("");
    setError("");

    try {
      const info = await invoke<BackendInfo>("get_backend_info");
      if (configureApiClientFromBackendInfo(info)) {
        setStartupStatus("ready");
        let retries = 0;
        while (retries < 15) {
          const ok = await checkHealth();
          if (ok) return;
          await new Promise((r) => setTimeout(r, 1000));
          retries++;
        }
      }
      setStartupStatus("failed");
      setStartupError("O servico nao respondeu a tempo.");
    } catch (err) {
      console.error("FileENIAC startup retry failed", err);
      setStartupStatus("failed");
      setStartupError("Erro ao verificar o ambiente.");
    }
  }

  async function handleConnect() {
    if (!wsPath.trim()) {
      setError("Informe o caminho do workspace");
      return;
    }
    setConnecting(true);
    setError("");
    try {
      const info = await getWorkspace(wsPath.trim());
      setWsInfo(info);
      setStep("ready");
    } catch (e: any) {
      setError(e.message || "Nao foi possivel conectar ao workspace informado.");
    }
    setConnecting(false);
  }

  function handleEnter() {
    navigate("/dashboard");
  }

  if (startupStatus === "waiting") {
    return <AppLoadingState message="Configurando ambiente..." />;
  }

  if (startupStatus === "failed") {
    return (
      <AppErrorState
        message={startupError}
        onRetry={handleRetryStartup}
      />
    );
  }

  if (step === "welcome") {
    return (
      <div className="fixed inset-0 bg-eniac-950 flex items-center justify-center">
        <div className="text-center max-w-md px-8">
          <div className="mx-auto mb-6 w-16 h-16 rounded-2xl bg-eniac-900/60 border border-eniac-700/30 flex items-center justify-center">
            <svg className="h-8 w-8 text-eniac-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-white mb-2 tracking-tight">FileENIAC</h1>
          <p className="text-eniac-400 mb-8 text-sm">Gerencie projetos e deploys do seu desktop</p>
          <button
            onClick={() => setStep("config")}
            className="w-full py-3 px-6 bg-eniac-600 text-white rounded-lg font-semibold hover:bg-eniac-700 transition-colors"
          >
            Comecar
          </button>
        </div>
      </div>
    );
  }

  if (step === "config") {
    return (
      <div className="fixed inset-0 bg-eniac-950 flex items-center justify-center">
        <div className="bg-eniac-900/80 border border-eniac-700/30 rounded-2xl shadow-2xl p-8 w-full max-w-md">
          <h2 className="text-xl font-bold text-white mb-2">
            Conectar Workspace
          </h2>
          <p className="text-sm text-eniac-300 mb-6">
            Informe o caminho do workspace que deseja gerenciar.
          </p>
          <label className="block text-sm font-medium text-eniac-200 mb-1">
            Caminho do Workspace
          </label>
          <div className="flex gap-2">
            <input
              type="text"
              value={wsPath}
              onChange={(e) => setWsPath(e.target.value)}
              placeholder="C:/projetos/meu-workspace"
              className="flex-1 px-3 py-2.5 bg-eniac-950/60 border border-eniac-700/40 rounded-lg text-sm text-white placeholder-eniac-500 focus:outline-none focus:ring-2 focus:ring-eniac-500 focus:border-transparent"
            />
            <button
              type="button"
              onClick={async () => { const p = await pickFolder(); if (p) setWsPath(p); }}
              className="px-3 py-2.5 border border-eniac-700/40 rounded-lg text-sm font-medium text-eniac-200 hover:bg-eniac-800/60 transition-colors"
            >
              Procurar
            </button>
          </div>
          {error && (
            <p className="mt-2 text-red-400 text-sm">{error}</p>
          )}
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setStep("welcome")}
              className="flex-1 py-2.5 px-4 border border-eniac-700/40 rounded-lg text-sm font-medium text-eniac-200 hover:bg-eniac-800/60 transition-colors"
            >
              Voltar
            </button>
            <button
              onClick={handleConnect}
              disabled={connecting}
              className="flex-1 py-2.5 px-4 bg-eniac-600 text-white rounded-lg text-sm font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors"
            >
              {connecting ? "Conectando..." : "Conectar"}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-eniac-950 flex items-center justify-center">
      <div className="bg-eniac-900/80 border border-eniac-700/30 rounded-2xl shadow-2xl p-8 w-full max-w-md text-center">
        <div className="mx-auto mb-4 w-16 h-16 rounded-2xl bg-emerald-900/30 border border-emerald-700/30 flex items-center justify-center">
          <svg className="h-8 w-8 text-emerald-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 className="text-xl font-bold text-white mb-2">
          Workspace Conectado
        </h2>
        <div className="text-left bg-eniac-950/60 rounded-lg p-4 mb-6 text-sm space-y-1 border border-eniac-700/20">
          <p>
            <span className="font-medium text-eniac-300">Nome:</span>{" "}
            <span className="text-white">{wsInfo?.name}</span>
          </p>
          <p>
            <span className="font-medium text-eniac-300">Projetos:</span>{" "}
            <span className="text-white">{wsInfo?.projects}</span>
          </p>
          <p>
            <span className="font-medium text-eniac-300">Servidores:</span>{" "}
            <span className="text-white">{wsInfo?.servers}</span>
          </p>
          <p>
            <span className="font-medium text-eniac-300">Deploys:</span>{" "}
            <span className="text-white">{wsInfo?.deploys}</span>
          </p>
        </div>
        <button
          onClick={handleEnter}
          className="w-full py-3 px-6 bg-eniac-600 text-white rounded-lg font-semibold hover:bg-eniac-700 transition-colors"
        >
          Entrar no Workspace
        </button>
      </div>
    </div>
  );
}
