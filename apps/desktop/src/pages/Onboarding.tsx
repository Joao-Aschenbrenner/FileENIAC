// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { invoke } from "@tauri-apps/api/core";
import { open } from "@tauri-apps/plugin-dialog";
import {
  checkHealth,
  configureApiClientFromBackendInfo,
  createWorkspace,
  enterWorkspace,
  listWorkspaces,
} from "../api/client";
import { STORAGE_KEYS, storageGet, storageSet } from "../api/storage";
import type { BackendInfo } from "../auth/tokenStorage";
import { AppLoadingState } from "../components/AppLoadingState";
import { AppErrorState } from "../components/AppErrorState";

type StartupStatus = "waiting" | "ready" | "failed";
type Step = "welcome" | "base" | "workspaces";

type WorkspaceSummary = {
  name: string;
  description?: string;
  path: string;
};

async function pickFolder(): Promise<string | null> {
  try {
    const selected = await open({ directory: true, multiple: false, title: "Selecionar pasta dos workspaces" });
    return selected || null;
  } catch {
    return null;
  }
}

function workspaceFolderName(name: string): string {
  const safe = name
    .trim()
    .replace(/[<>:"/\\|?*\x00-\x1F]/g, "-")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
  return safe || "workspace";
}

function joinWorkspacePath(root: string, name: string): string {
  const base = root.trim().replace(/[\\/]+$/g, "");
  return `${base}/${workspaceFolderName(name)}`;
}

export default function Onboarding() {
  const navigate = useNavigate();
  const [step, setStep] = useState<Step>("welcome");
  const [rootPath, setRootPath] = useState(() => storageGet(STORAGE_KEYS.workspacesRoot) || "");
  const [workspaceName, setWorkspaceName] = useState("");
  const [workspaceDescription, setWorkspaceDescription] = useState("");
  const [workspaces, setWorkspaces] = useState<WorkspaceSummary[]>([]);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loadingWorkspaces, setLoadingWorkspaces] = useState(false);
  const [creating, setCreating] = useState(false);
  const [enteringPath, setEnteringPath] = useState("");
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

  async function loadWorkspaceList(path = rootPath) {
    if (!path.trim()) return;
    setLoadingWorkspaces(true);
    setError("");
    try {
      const items = await listWorkspaces(path.trim());
      setWorkspaces(Array.isArray(items) ? items : []);
    } catch (e: any) {
      setError(e.message || "Nao foi possivel carregar os workspaces dessa pasta.");
    }
    setLoadingWorkspaces(false);
  }

  async function handleStart() {
    setError("");
    setMessage("");
    if (rootPath.trim()) {
      setStep("workspaces");
      await loadWorkspaceList(rootPath.trim());
      return;
    }
    setStep("base");
  }

  async function handleSaveBaseFolder() {
    const path = rootPath.trim();
    if (!path) {
      setError("Informe a pasta onde os workspaces serao alocados");
      return;
    }
    storageSet(STORAGE_KEYS.workspacesRoot, path);
    setRootPath(path);
    setStep("workspaces");
    setMessage("");
    await loadWorkspaceList(path);
  }

  async function handleCreateWorkspace() {
    const name = workspaceName.trim();
    if (!rootPath.trim()) {
      setError("Informe primeiro a pasta-base dos workspaces");
      setStep("base");
      return;
    }
    if (!name) {
      setError("Informe o nome do workspace");
      return;
    }

    setCreating(true);
    setError("");
    setMessage("");
    try {
      const path = joinWorkspacePath(rootPath, name);
      const info = await createWorkspace(path, { name, description: workspaceDescription.trim() });
      setWorkspaces((current) => {
        const withoutDuplicate = current.filter((item) => item.path !== info.path);
        return [...withoutDuplicate, info].sort((a, b) => a.name.localeCompare(b.name));
      });
      setWorkspaceName("");
      setWorkspaceDescription("");
      setMessage("Workspace criado. Voce pode criar outro ou entrar nele.");
    } catch (e: any) {
      setError(e.message || "Nao foi possivel criar o workspace.");
    }
    setCreating(false);
  }

  async function handleEnterWorkspace(path: string) {
    setEnteringPath(path);
    setError("");
    try {
      await enterWorkspace(path);
      navigate("/configurar");
    } catch (e: any) {
      setError(e.message || "Nao foi possivel entrar no workspace.");
    }
    setEnteringPath("");
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
            onClick={handleStart}
            className="w-full py-3 px-6 bg-eniac-600 text-white rounded-lg font-semibold hover:bg-eniac-700 transition-colors"
          >
            Começar
          </button>
        </div>
      </div>
    );
  }

  if (step === "base") {
    return (
      <div className="fixed inset-0 bg-eniac-950 flex items-center justify-center">
        <div className="bg-eniac-900/80 border border-eniac-700/30 rounded-2xl shadow-2xl p-8 w-full max-w-md">
          <h2 className="text-xl font-bold text-white mb-2">Escolher Pasta dos Workspaces</h2>
          <p className="text-sm text-eniac-300 mb-6">
            Escolha a pasta-base onde o FileENIAC vai guardar seus workspaces. Cada workspace sera criado dentro dela.
          </p>
          <label className="block text-sm font-medium text-eniac-200 mb-1">Pasta-base</label>
          <div className="flex gap-2">
            <input
              type="text"
              value={rootPath}
              onChange={(e) => setRootPath(e.target.value)}
              placeholder="C:/projetos/ENIAC_SYSTEMS"
              className="flex-1 px-3 py-2.5 bg-eniac-950/60 border border-eniac-700/40 rounded-lg text-sm text-white placeholder-eniac-500 focus:outline-none focus:ring-2 focus:ring-eniac-500 focus:border-transparent"
            />
            <button
              type="button"
              onClick={async () => { const p = await pickFolder(); if (p) setRootPath(p); }}
              className="px-3 py-2.5 border border-eniac-700/40 rounded-lg text-sm font-medium text-eniac-200 hover:bg-eniac-800/60 transition-colors"
            >
              Procurar
            </button>
          </div>
          {error && <p className="mt-2 text-red-400 text-sm">{error}</p>}
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setStep("welcome")}
              className="flex-1 py-2.5 px-4 border border-eniac-700/40 rounded-lg text-sm font-medium text-eniac-200 hover:bg-eniac-800/60 transition-colors"
            >
              Voltar
            </button>
            <button
              onClick={handleSaveBaseFolder}
              disabled={loadingWorkspaces}
              className="flex-1 py-2.5 px-4 bg-eniac-600 text-white rounded-lg text-sm font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors"
            >
              {loadingWorkspaces ? "Carregando..." : "Continuar"}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-eniac-950 overflow-y-auto">
      <div className="min-h-full flex items-center justify-center p-6">
        <div className="bg-eniac-900/80 border border-eniac-700/30 rounded-2xl shadow-2xl p-8 w-full max-w-3xl">
          <div className="flex items-start justify-between gap-4 mb-6">
            <div>
              <h2 className="text-xl font-bold text-white mb-2">Criar ou Entrar em um Workspace</h2>
              <p className="text-sm text-eniac-300">
                Pasta-base: <span className="text-white break-all">{rootPath}</span>
              </p>
            </div>
            <button
              onClick={() => { setStep("base"); setError(""); setMessage(""); }}
              className="px-3 py-2 border border-eniac-700/40 rounded-lg text-sm font-medium text-eniac-200 hover:bg-eniac-800/60 transition-colors"
            >
              Trocar pasta
            </button>
          </div>

          {error && <p className="mb-4 text-red-400 text-sm">{error}</p>}
          {message && <p className="mb-4 text-emerald-300 text-sm">{message}</p>}

          <div className="grid gap-6 lg:grid-cols-[1fr_1fr]">
            <section className="bg-eniac-950/50 border border-eniac-700/30 rounded-xl p-5">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-white font-semibold">Workspaces nesta pasta</h3>
                <button
                  onClick={() => loadWorkspaceList()}
                  disabled={loadingWorkspaces}
                  className="text-xs text-eniac-300 hover:text-white disabled:opacity-60"
                >
                  {loadingWorkspaces ? "Atualizando..." : "Atualizar"}
                </button>
              </div>

              {workspaces.length === 0 && !loadingWorkspaces && (
                <div className="border border-dashed border-eniac-700/50 rounded-lg p-4 text-sm text-eniac-300">
                  Nenhum workspace criado nessa pasta ainda.
                </div>
              )}

              <div className="space-y-3 max-h-80 overflow-y-auto">
                {workspaces.map((ws) => (
                  <div key={ws.path} className="bg-eniac-900/80 border border-eniac-700/30 rounded-lg p-4">
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0">
                        <p className="text-white font-medium truncate">{ws.name}</p>
                        {ws.description && <p className="text-xs text-eniac-300 mt-1">{ws.description}</p>}
                        <p className="text-xs text-eniac-400 mt-1 break-all">{ws.path}</p>
                      </div>
                      <button
                        onClick={() => handleEnterWorkspace(ws.path)}
                        disabled={!!enteringPath}
                        className="px-3 py-1.5 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 disabled:opacity-60 transition-colors"
                      >
                        {enteringPath === ws.path ? "Entrando..." : "Entrar"}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </section>

            <section className="bg-eniac-950/50 border border-eniac-700/30 rounded-xl p-5">
              <h3 className="text-white font-semibold mb-2">Criar novo workspace</h3>
              <p className="text-sm text-eniac-300 mb-4">
                Crie quantos workspaces quiser dentro da pasta-base. Depois escolha em qual deseja entrar.
              </p>
              <label className="block text-sm font-medium text-eniac-200 mb-1">Nome do workspace</label>
              <input
                type="text"
                value={workspaceName}
                onChange={(e) => setWorkspaceName(e.target.value)}
                placeholder="Cliente X"
                className="w-full px-3 py-2.5 bg-eniac-950/60 border border-eniac-700/40 rounded-lg text-sm text-white placeholder-eniac-500 focus:outline-none focus:ring-2 focus:ring-eniac-500 focus:border-transparent mb-4"
              />
              <label className="block text-sm font-medium text-eniac-200 mb-1">Descricao opcional</label>
              <textarea
                value={workspaceDescription}
                onChange={(e) => setWorkspaceDescription(e.target.value)}
                placeholder="Projetos do cliente, ambiente interno, etc."
                rows={3}
                className="w-full px-3 py-2.5 bg-eniac-950/60 border border-eniac-700/40 rounded-lg text-sm text-white placeholder-eniac-500 focus:outline-none focus:ring-2 focus:ring-eniac-500 focus:border-transparent mb-4"
              />
              <button
                onClick={handleCreateWorkspace}
                disabled={creating}
                className="w-full py-3 px-6 bg-eniac-600 text-white rounded-lg font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors"
              >
                {creating ? "Criando..." : "Criar Workspace"}
              </button>
            </section>
          </div>
        </div>
      </div>
    </div>
  );
}
