// SPDX-License-Identifier: MIT
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { getWorkspace, checkHealth } from "../api/client";
import { open } from "@tauri-apps/plugin-dialog";

async function pickFolder(): Promise<string | null> {
  try {
    const selected = await open({ directory: true, multiple: false, title: "Selecionar Workspace" });
    return selected || null;
  } catch {
    return null;
  }
}

export default function Onboarding() {
  const navigate = useNavigate();
  const [step, setStep] = useState<"welcome" | "config" | "ready">("welcome");
  const [wsPath, setWsPath] = useState("");
  const [error, setError] = useState("");
  const [wsInfo, setWsInfo] = useState<any>(null);
  const [checking, setChecking] = useState(false);
  const [connecting, setConnecting] = useState(false);

  const [healthRetries, setHealthRetries] = useState(0);
  const [maxRetries] = useState(15);

  async function handleCheckBackend() {
    setChecking(true);
    setError("");
    for (let i = 0; i < maxRetries; i++) {
      setHealthRetries(i + 1);
      const ok = await checkHealth();
      if (ok) {
        setStep("config");
        setChecking(false);
        return;
      }
      await new Promise((r) => setTimeout(r, 1000));
    }
    setError("Não foi possível iniciar o backend local. Tente novamente.");
    setChecking(false);
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
      setError(e.message);
    }
    setConnecting(false);
  }

  function handleEnter() {
    navigate("/dashboard");
  }

  if (step === "welcome") {
    return (
      <div className="flex h-screen items-center justify-center bg-gradient-to-br from-eniac-900 to-eniac-700">
        <div className="text-center max-w-md px-8">
          <h1 className="text-4xl font-bold text-white mb-2">FileENIAC</h1>
          <p className="text-eniac-200 mb-8">Desktop</p>
          <p className="text-white/80 mb-8 text-sm leading-relaxed">
            Gerencie seus projetos, deploys FTPS e mantenha seus workspaces
            organizados — tudo do seu desktop.
          </p>
          <button
            onClick={handleCheckBackend}
            disabled={checking}
            className="w-full py-3 px-6 bg-white text-eniac-900 rounded-lg font-semibold hover:bg-eniac-100 disabled:opacity-60 transition-colors"
          >
            {checking ? `Inicializando backend... (${healthRetries}/${maxRetries})` : "Começar"}
          </button>
          {error && (
            <p className="mt-4 text-red-300 text-sm">{error}</p>
          )}
        </div>
      </div>
    );
  }

  if (step === "config") {
    return (
      <div className="flex h-screen items-center justify-center bg-gradient-to-br from-eniac-900 to-eniac-700">
        <div className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-md">
          <h2 className="text-xl font-bold text-gray-800 mb-2">
            Conectar Workspace
          </h2>
          <p className="text-sm text-gray-500 mb-6">
            Informe o caminho do workspace FileENIAC que deseja gerenciar.
          </p>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Caminho do Workspace
          </label>
          <div className="flex gap-2">
            <input
              type="text"
              value={wsPath}
              onChange={(e) => setWsPath(e.target.value)}
              placeholder="C:/projetos/meu-workspace"
              className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500 focus:border-transparent"
            />
            <button
              type="button"
              onClick={async () => { const p = await pickFolder(); if (p) setWsPath(p); }}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Procurar
            </button>
          </div>
          {error && (
            <p className="mt-2 text-red-600 text-sm">{error}</p>
          )}
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setStep("welcome")}
              className="flex-1 py-2 px-4 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Voltar
            </button>
            <button
              onClick={handleConnect}
              disabled={connecting}
              className="flex-1 py-2 px-4 bg-eniac-600 text-white rounded-lg text-sm font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors"
            >
              {connecting ? "Conectando..." : "Conectar"}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen items-center justify-center bg-gradient-to-br from-eniac-900 to-eniac-700">
      <div className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-md text-center">
        <div className="text-4xl mb-4">✅</div>
        <h2 className="text-xl font-bold text-gray-800 mb-2">
          Workspace Conectado
        </h2>
        <div className="text-left bg-gray-50 rounded-lg p-4 mb-6 text-sm space-y-1">
          <p>
            <span className="font-medium text-gray-600">Nome:</span>{" "}
            {wsInfo?.name}
          </p>
          <p>
            <span className="font-medium text-gray-600">Projetos:</span>{" "}
            {wsInfo?.projects}
          </p>
          <p>
            <span className="font-medium text-gray-600">Servidores:</span>{" "}
            {wsInfo?.servers}
          </p>
          <p>
            <span className="font-medium text-gray-600">Deploys:</span>{" "}
            {wsInfo?.deploys}
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
