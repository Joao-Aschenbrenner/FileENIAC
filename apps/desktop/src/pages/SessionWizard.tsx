// SPDX-License-Identifier: MIT
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createSession, updateSession, gitHubLogin, getWorkspace } from "../api/client";
import { useSession } from "../context/SessionContext";
import { pickFolder } from "../utils/dialog";

type Step = "name" | "github" | "workspace";

const STEP_LABELS = ["Nome", "GitHub", "Workspace"];

function StepIndicator({ current }: { current: Step }) {
  const idx = STEP_LABELS.indexOf(current);
  return (
    <div className="flex items-center gap-2 mb-6">
      {STEP_LABELS.map((label, i) => (
        <div key={label} className="flex items-center gap-2">
          <div
            className={`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold ${
              i < idx
                ? "bg-green-500 text-white"
                : i === idx
                ? "bg-eniac-600 text-white"
                : "bg-gray-200 text-gray-400"
            }`}
          >
            {i < idx ? "✓" : i + 1}
          </div>
          {i < STEP_LABELS.length - 1 && <div className="w-8 h-0.5 bg-gray-200" />}
        </div>
      ))}
    </div>
  );
}

function StepName({
  name,
  setName,
  description,
  setDescription,
  onContinue,
  onCancel,
  saving,
  error,
}: {
  name: string;
  setName: (v: string) => void;
  description: string;
  setDescription: (v: string) => void;
  onContinue: () => void;
  onCancel: () => void;
  saving: boolean;
  error: string;
}) {
  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-1">Nome da Sessão</h2>
      <p className="text-sm text-gray-500 mb-6">
        Dê um nome para identificar esta sessão de trabalho.
      </p>
      <label className="block text-sm font-medium text-gray-700 mb-1">Nome</label>
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Ex: Projeto Cliente X"
        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500 mb-4"
      />
      <label className="block text-sm font-medium text-gray-700 mb-1">Descrição (opcional)</label>
      <textarea
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="Ex: Sessão para gerenciar deploys do projeto Cliente X"
        rows={2}
        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500 mb-4"
      />
      {error && <p className="text-red-600 text-sm mb-3">{error}</p>}
      <div className="flex gap-3">
        <button onClick={onCancel} className="flex-1 py-2 px-4 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
          Cancelar
        </button>
        <button onClick={onContinue} disabled={saving || !name.trim()} className="flex-1 py-2 px-4 bg-eniac-600 text-white rounded-lg text-sm font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors">
          {saving ? "Criando..." : "Continuar"}
        </button>
      </div>
    </div>
  );
}

function StepGitHub({
  token,
  setToken,
  onSave,
  onSkip,
  saving,
  error,
}: {
  token: string;
  setToken: (v: string) => void;
  onSave: () => void;
  onSkip: () => void;
  saving: boolean;
  error: string;
}) {
  const [showSteps, setShowSteps] = useState(false);
  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-1">GitHub</h2>
      <p className="text-sm text-gray-500 mb-4">
        Conecte sua conta GitHub para importar repositórios.
      </p>

      {/* Instruções colapsáveis */}
      <button
        onClick={() => setShowSteps(!showSteps)}
        className="w-full text-left bg-blue-50 border border-blue-200 rounded-lg px-4 py-3 mb-4 hover:bg-blue-100 transition-colors"
      >
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium text-blue-800">
            Como gerar o token (clique para ver)
          </span>
          <span className="text-blue-500 text-xs">{showSteps ? "▲" : "▼"}</span>
        </div>
      </button>

      {showSteps && (
        <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 mb-4 text-sm space-y-3">
          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-eniac-600 text-white rounded-full flex items-center justify-center text-xs font-bold">1</span>
            <div>
              <p className="font-medium text-gray-800">Acesse github.com/settings/tokens</p>
              <p className="text-gray-500 text-xs mt-0.5">
                Clique em{" "}
                <a href="https://github.com/settings/tokens" target="_blank" rel="noopener noreferrer" className="text-eniac-600 underline">
                  github.com/settings/tokens
                </a>
              </p>
            </div>
          </div>

          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-eniac-600 text-white rounded-full flex items-center justify-center text-xs font-bold">2</span>
            <div>
              <p className="font-medium text-gray-800">Clique em "Generate new token" → <span className="text-green-700">"Generate new token (classic)"</span></p>
              <p className="text-gray-500 text-xs mt-0.5">
                IMPORTANTE: deve ser o <strong>classic</strong>, não o fine-grained.
              </p>
            </div>
          </div>

          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-eniac-600 text-white rounded-full flex items-center justify-center text-xs font-bold">3</span>
            <div>
              <p className="font-medium text-gray-800">Marque as permissões:</p>
              <div className="flex flex-wrap gap-1 mt-1">
                {["repo", "workflow", "delete:packages", "write:packages"].map((scope) => (
                  <code key={scope} className="bg-green-100 text-green-800 px-2 py-0.5 rounded text-xs font-mono">{scope}</code>
                ))}
              </div>
              <p className="text-gray-500 text-xs mt-1">
                Em "Select scopes", marque <strong>repo</strong> (acesso completo a repositórios),{" "}
                <strong>workflow</strong>, <strong>delete:packages</strong> e <strong>write:packages</strong>.
              </p>
            </div>
          </div>

          <div className="flex gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-eniac-600 text-white rounded-full flex items-center justify-center text-xs font-bold">4</span>
            <div>
              <p className="font-medium text-gray-800">Clique em "Generate token" e copie o valor</p>
              <p className="text-gray-500 text-xs mt-0.5">
                O token aparece apenas uma vez. Cole ele no campo abaixo.
              </p>
            </div>
          </div>
        </div>
      )}

      <label className="block text-sm font-medium text-gray-700 mb-1">Token</label>
      <input
        type="password"
        value={token}
        onChange={(e) => setToken(e.target.value)}
        placeholder="ghp_..."
        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500 mb-1"
      />
      <p className="text-xs text-gray-400 mb-4">
        Deve começar com <code className="bg-gray-100 px-1 rounded">ghp_</code> e ser do tipo <strong>classic</strong>.
      </p>

      {error && <p className="text-red-600 text-sm mb-3">{error}</p>}

      <div className="flex gap-3">
        <button onClick={onSkip} className="flex-1 py-2 px-4 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
          Pular
        </button>
        <button onClick={onSave} disabled={saving || !token.trim()} className="flex-1 py-2 px-4 bg-eniac-600 text-white rounded-lg text-sm font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors">
          {saving ? "Validando..." : "Salvar"}
        </button>
      </div>
    </div>
  );
}

function StepWorkspace({
  wsPath,
  setWsPath,
  onFinish,
  onSkip,
  creating,
  error,
}: {
  wsPath: string;
  setWsPath: (v: string) => void;
  onFinish: () => void;
  onSkip: () => void;
  creating: boolean;
  error: string;
}) {
  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-1">Workspace</h2>
      <p className="text-sm text-gray-500 mb-4">
        O workspace é a pasta onde seus projetos e configurações ficam armazenados.
      </p>

      <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 mb-4 text-sm space-y-2">
        <p className="font-medium text-gray-800">Exemplo de estrutura:</p>
        <pre className="text-xs text-gray-600 font-mono bg-white p-2 rounded border">
{`meu-workspace/
├── projects/
│   ├── projeto-a/
│   └── projeto-b/
└── servers.db`}
        </pre>
        <p className="text-gray-500 text-xs">
          Cada sessão pode apontar para um workspace diferente.
        </p>
      </div>

      <label className="block text-sm font-medium text-gray-700 mb-1">Caminho do Workspace</label>
      <div className="flex gap-2 mb-2">
        <input
          type="text"
          value={wsPath}
          onChange={(e) => setWsPath(e.target.value)}
          placeholder="C:/projetos/meu-workspace"
          className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
        />
        <button
          type="button"
          onClick={async () => { const p = await pickFolder(); if (p) setWsPath(p); }}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          Procurar
        </button>
      </div>

      {error && <p className="text-red-600 text-sm mb-3">{error}</p>}

      <p className="text-xs text-gray-400 mb-4">
        Se houver repositórios Git dentro desta pasta, eles serão reconhecidos nas próximas etapas.
        Caso não existam, crie um workspace com o comando
        <code className="bg-gray-100 px-1 rounded mx-1">fileeniac workspace init --name "MeuWorkspace" --path&nbsp;
          {"<caminho>"}
        </code>
        ou inicialize um repositório Git (<code>git init</code>) e faça o push para o GitHub.
      </p>

      <div className="flex gap-3">
        <button onClick={onSkip} className="flex-1 py-2 px-4 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
          Pular
        </button>
        <button onClick={onFinish} disabled={creating} className="flex-1 py-2 px-4 bg-eniac-600 text-white rounded-lg text-sm font-semibold hover:bg-eniac-700 disabled:opacity-60 transition-colors">
          {creating ? "Configurando..." : "Finalizar"}
        </button>
      </div>
    </div>
  );
}

export default function SessionWizard() {
  const navigate = useNavigate();
  const { refresh } = useSession();
  const [step, setStep] = useState<Step>("name");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [token, setToken] = useState("");
  const [wsPath, setWsPath] = useState("");
  const [sessionId, setSessionId] = useState<number | null>(null);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [creating, setCreating] = useState(false);

  async function handleCreateSession() {
    if (!name.trim()) { setError("Nome da sessão é obrigatório"); return; }
    setSaving(true);
    setError("");
    try {
      const sess = await createSession({ name: name.trim(), description: description.trim() });
      setSessionId(sess.id);
      setStep("github");
    } catch (e: any) {
      setError(e.message);
    }
    setSaving(false);
  }

  async function handleSaveGitHub() {
    if (!token.trim()) { setError("Token é obrigatório"); return; }
    if (!token.trim().startsWith("ghp_")) {
      setError("O token deve começar com 'ghp_' — certifique-se de que é um token classic");
      return;
    }
    setSaving(true);
    setError("");
    try {
      const result = await gitHubLogin(token.trim());
      if (sessionId) {
        await updateSession(sessionId, { github_token: token.trim(), github_user: result.user || "" });
      }
      setStep("workspace");
    } catch (e: any) {
      setError(e.message);
    }
    setSaving(false);
  }

  async function handleConfigureWorkspace() {
    setCreating(true);
    setError("");
    try {
      if (wsPath.trim()) {
        await getWorkspace(wsPath.trim());
        if (sessionId) {
          await updateSession(sessionId, { workspace_path: wsPath.trim() });
        }
      }
      await refresh();
      navigate("/");
    } catch (e: any) {
      setError(e.message);
    }
    setCreating(false);
  }

  async function handleSkipWorkspace() {
    await refresh();
    navigate("/");
  }

  return (
    <div className="flex h-screen items-center justify-center bg-gradient-to-br from-eniac-900 to-eniac-700 p-4">
      <div className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-lg">
        <StepIndicator current={step} />

        {step === "name" && (
          <StepName
            name={name}
            setName={setName}
            description={description}
            setDescription={setDescription}
            onContinue={handleCreateSession}
            onCancel={() => navigate("/")}
            saving={saving}
            error={error}
          />
        )}

        {step === "github" && (
          <StepGitHub
            token={token}
            setToken={setToken}
            onSave={handleSaveGitHub}
            onSkip={() => setStep("workspace")}
            saving={saving}
            error={error}
          />
        )}

        {step === "workspace" && (
          <StepWorkspace
            wsPath={wsPath}
            setWsPath={setWsPath}
            onFinish={handleConfigureWorkspace}
            onSkip={handleSkipWorkspace}
            creating={creating}
            error={error}
          />
        )}
      </div>
    </div>
  );
}
