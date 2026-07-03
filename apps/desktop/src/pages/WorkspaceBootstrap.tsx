// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getGitHubStatus, getGitHubOrganizations } from "../api/client";
import { Badge } from "../components/ui/Badge";
import { Loader } from "../components/ui/Loader";

export default function WorkspaceBootstrap() {
  const navigate = useNavigate();
  const [status, setStatus] = useState<any>(null);
  const [orgs, setOrgs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const wsPath = localStorage.getItem("eniac_ws_path") || "";

  useEffect(() => {
    const activeWorkspace = localStorage.getItem("eniac_ws_path") || "";
    if (!activeWorkspace) {
      setLoading(false);
      return;
    }
    Promise.all([
      getGitHubStatus(),
      getGitHubOrganizations().catch(() => []),
    ])
      .then(([s, o]) => {
        setStatus(s);
        setOrgs(Array.isArray(o) ? o : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <Loader text="Verificando ambiente..." />;

  if (!wsPath) {
    return (
      <div className="max-w-xl mx-auto mt-8">
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-8 text-center">
          <h2 className="text-2xl font-bold text-gray-800">Nenhum workspace selecionado</h2>
          <p className="text-sm text-gray-500 mt-2">
            Escolha ou crie um workspace antes de conectar seus projetos.
          </p>
          <button onClick={() => navigate("/")}
            className="mt-6 px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Ir para Workspaces
          </button>
        </div>
      </div>
    );
  }

  const steps = [
    {
      label: "Workspace selecionado",
      done: true,
      detail: wsPath,
    },
    {
      label: "Conectar GitHub",
      done: status?.authenticated,
      action: "/github/login",
      detail: status?.authenticated ? `Autenticado como ${status.user}` : "Conectar GitHub",
    },
    {
      label: "Organizacoes",
      done: orgs.length > 0,
      action: "/github/orgs",
      detail: orgs.length > 0 ? `${orgs.length} organizacoes encontradas` : "Nenhuma organizacao",
    },
  ];

  const allDone = steps.every((s) => s.done);

  return (
    <div className="max-w-xl mx-auto mt-8">
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold text-gray-800">Configuracao do Ambiente</h2>
        <p className="text-sm text-gray-500 mt-1">Verifique o workspace e conecte provedores</p>
      </div>

      <div className="space-y-3 mb-8">
        {steps.map((step, idx) => (
          <div key={idx}
            className={`bg-white rounded-xl border p-5 shadow-sm flex items-center justify-between ${
              step.done ? "border-green-200" : "border-gray-200"
            }`}
          >
            <div className="flex items-center gap-4">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold ${
                step.done ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-500"
              }`}>
                {step.done ? "✓" : idx + 1}
              </div>
              <div>
                <p className="font-medium text-gray-800">{step.label}</p>
                <p className="text-xs text-gray-500">{step.detail}</p>
              </div>
            </div>
            {!step.done && step.action && (
              <button onClick={() => navigate(step.action)}
                className="px-3 py-1.5 text-xs bg-eniac-600 text-white rounded-lg hover:bg-eniac-700 transition-colors">
                Configurar
              </button>
            )}
            {step.done && <Badge variant="success">Feito</Badge>}
          </div>
        ))}
      </div>

      {!status?.authenticated && (
        <div className="text-center mt-6">
          <button onClick={() => navigate("/github/login")}
            className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Conectar GitHub
          </button>
        </div>
      )}

      {status?.authenticated && !allDone && (
        <div className="text-center mt-6">
          <button onClick={() => navigate("/github/orgs")}
            className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Explorar Organizacoes
          </button>
        </div>
      )}

      {allDone && (
        <div className="text-center mt-6 space-y-3">
          <div className="bg-green-50 border border-green-200 rounded-lg p-4 inline-block">
            <p className="text-green-700 font-medium">Ambiente Pronto!</p>
            <p className="text-green-600 text-sm mt-1">Workspace configurado</p>
          </div>
          <div>
            <button onClick={() => navigate("/projects")}
              className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
              Ir para Projetos
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
