// SPDX-License-Identifier: MIT
import { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { listProjects, getDeploys, getProject, executeDeploy } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Table } from "../components/ui/Table";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Toast } from "../components/ui/Toast";

export default function DeployCenter() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState<string>("");
  const [deploys, setDeploys] = useState<any[]>([]);
  const [projectInfo, setProjectInfo] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [deploying, setDeploying] = useState(false);
  const [error, setError] = useState("");
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" } | null>(null);

  function loadProjects() {
    setLoading(true);
    setError("");
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    listProjects(wsPath)
      .then(setProjects)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => { loadProjects(); }, []);

  const loadDeploys = useCallback(() => {
    if (!selectedProject) return;
    getProject(localStorage.getItem("eniac_ws_path") || "", selectedProject)
      .then(setProjectInfo)
      .catch(() => {});
    getDeploys(selectedProject, 10)
      .then(setDeploys)
      .catch(() => {});
  }, [selectedProject]);

  useEffect(() => { loadDeploys() }, [loadDeploys]);

  async function handleDeploy(useFallback = false) {
    if (!selectedProject) return;
    setDeploying(true);
    try {
      const result = await executeDeploy(selectedProject, useFallback);
      setToast({ message: `Deploy concluído: ${result.status || "ok"}`, type: "success" });
      loadDeploys();
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setDeploying(false);
  }

  if (loading) return <Loader />;
  if (error) return <ErrorState message={error} onRetry={loadProjects} />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Deploy Center</h2>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <Card title="Projeto">
          <select
            value={selectedProject}
            onChange={(e) => setSelectedProject(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
          >
            <option value="">Selecione...</option>
            {projects.map((p) => (
              <option key={p.id} value={p.name}>{p.name}</option>
            ))}
          </select>
        </Card>

        <Card title="Estado Atual">
          {projectInfo ? (
            <div className="text-sm space-y-1">
              <p><span className="text-gray-500">Status:</span> <Badge variant={projectInfo.divergence_status === "sincronizado" ? "success" : "warning"}>{projectInfo.divergence_status}</Badge></p>
              <p><span className="text-gray-500">Branch:</span> {projectInfo.branch}</p>
            </div>
          ) : (
            <p className="text-sm text-gray-400">Selecione um projeto</p>
          )}
        </Card>

        <Card title="Ações">
          {selectedProject ? (
            <div className="flex flex-col gap-2">
              <button onClick={() => handleDeploy(false)} disabled={deploying}
                className="w-full py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 disabled:opacity-50 transition-colors">
                {deploying ? "Deployando..." : "Executar Deploy"}
              </button>
              <button onClick={() => handleDeploy(true)} disabled={deploying}
                className="w-full py-2 bg-amber-600 text-white text-sm rounded-lg hover:bg-amber-700 disabled:opacity-50 transition-colors">
                Deploy (Fallback)
              </button>
            </div>
          ) : (
            <p className="text-sm text-gray-400">Selecione um projeto</p>
          )}
        </Card>
      </div>

      {selectedProject && (
        <Card title="Histórico de Deploys">
          <Table
            columns={[
              { key: "id", header: "#" },
              { key: "status", header: "Status", render: (v: string) => <Badge variant={v === "success" ? "success" : v === "running" ? "info" : "danger"}>{v}</Badge> },
              { key: "version", header: "Versão" },
              { key: "created_at", header: "Data" },
              { key: "duration_ms", header: "Duração", render: (v: number) => v ? `${(v / 1000).toFixed(1)}s` : "-" },
            ]}
            data={deploys}
            onRowClick={() => navigate(`/projects/${selectedProject}`)}
          />
        </Card>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
