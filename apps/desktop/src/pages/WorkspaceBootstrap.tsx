// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getGitHubStatus, getGitHubOrganizations, listRepositories, listProjects } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Loader } from "../components/ui/Loader";

export default function WorkspaceBootstrap() {
  const navigate = useNavigate();
  const [status, setStatus] = useState<any>(null);
  const [orgs, setOrgs] = useState<any[]>([]);
  const [repos, setRepos] = useState<any[]>([]);
  const [projects, setProjects] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      getGitHubStatus(),
      getGitHubOrganizations().catch(() => []),
      listRepositories().catch(() => []),
      listProjects(localStorage.getItem("eniac_ws_path") || "").catch(() => []),
    ])
      .then(([s, o, r, p]) => {
        setStatus(s);
        setOrgs(Array.isArray(o) ? o : []);
        setRepos(Array.isArray(r) ? r : []);
        setProjects(Array.isArray(p) ? p : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <Loader text="Verificando ambiente..." />;

  const steps = [
    {
      label: "Autenticação GitHub",
      done: status?.authenticated,
      action: "/github/login",
      detail: status?.authenticated ? `Autenticado como ${status.user}` : "Conectar GitHub",
    },
    {
      label: "Organizações",
      done: orgs.length > 0,
      action: "/github/orgs",
      detail: `${orgs.length} organizações encontradas`,
    },
    {
      label: "Repositórios Importados",
      done: repos.length > 0,
      action: "/github/repos",
      detail: `${repos.length} repositórios importados`,
    },
    {
      label: "Projetos no Workspace",
      done: projects.length > 0,
      action: "/projects",
      detail: `${projects.length} projetos registrados`,
    },
  ];

  return (
    <div className="max-w-xl mx-auto mt-8">
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold text-gray-800">Workspace Bootstrap</h2>
        <p className="text-sm text-gray-500 mt-1">Configure seu ambiente do zero em poucos passos</p>
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
            {!step.done && (
              <button onClick={() => navigate(step.action)}
                className="px-3 py-1.5 text-xs bg-eniac-600 text-white rounded-lg hover:bg-eniac-700 transition-colors">
                Configurar
              </button>
            )}
            {step.done && <Badge variant="success">Feito</Badge>}
          </div>
        ))}
      </div>

      <Card title="Resumo do Ambiente">
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-500">GitHub</span>
            <span>{status?.authenticated ? `Conectado (${status.user})` : "Desconectado"}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Organizações</span>
            <span>{orgs.length}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Repositórios importados</span>
            <span>{repos.length}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Projetos</span>
            <span>{projects.length}</span>
          </div>
        </div>
      </Card>

      {!status?.authenticated && (
        <div className="text-center mt-6">
          <button onClick={() => navigate("/github/login")}
            className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Iniciar Bootstrap
          </button>
        </div>
      )}

      {status?.authenticated && projects.length === 0 && (
        <div className="text-center mt-6">
          <button onClick={() => navigate("/github/orgs")}
            className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Importar Repositórios
          </button>
        </div>
      )}

      {status?.authenticated && projects.length > 0 && (
        <div className="text-center mt-6">
          <div className="bg-green-50 border border-green-200 rounded-lg p-4 inline-block">
            <p className="text-green-700 font-medium">Ambiente Pronto!</p>
            <p className="text-green-600 text-sm mt-1">Workspace configurado com {projects.length} projetos</p>
          </div>
        </div>
      )}
    </div>
  );
}
