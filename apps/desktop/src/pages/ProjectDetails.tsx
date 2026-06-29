// SPDX-License-Identifier: MIT
import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getProject, listServers, getDeploys, executeDeploy, executeRollback, executeVerify, getDiff } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Table } from "../components/ui/Table";
import { Modal } from "../components/ui/Modal";
import { Toast } from "../components/ui/Toast";

export default function ProjectDetails() {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const [project, setProject] = useState<any>(null);
  const [servers, setServers] = useState<any[]>([]);
  const [deploys, setDeploys] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [actionModal, setActionModal] = useState<"deploy" | "rollback" | "verify" | null>(null);
  const [running, setRunning] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" } | null>(null);
  const [diffData, setDiffData] = useState<any>(null);
  const [diffModal, setDiffModal] = useState(false);
  const [diffLoading, setDiffLoading] = useState(false);

  const load = useCallback(() => {
    if (!name) return;
    setLoading(true);
    setError("");
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    Promise.all([
      getProject(wsPath, name),
      listServers(wsPath, name),
      getDeploys(name, 5),
    ])
      .then(([p, s, d]) => {
        setProject(p);
        setServers(s);
        setDeploys(d);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [name]);

  useEffect(() => { load() }, [load]);

  async function handleAction(action: string) {
    if (!name) return;
    setRunning(true);
    try {
      if (action === "deploy") await executeDeploy(name);
      else if (action === "rollback") await executeRollback(name);
      else if (action === "verify") await executeVerify(name);
      setToast({ message: `${action} executado com sucesso`, type: "success" });
      setActionModal(null);
      load();
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setRunning(false);
  }

  async function handleLoadDiff() {
    if (!name) return;
    setDiffLoading(true);
    try {
      const data = await getDiff(name);
      setDiffData(data);
      setDiffModal(true);
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setDiffLoading(false);
  }

  if (!name) { navigate("/projects"); return null; }
  if (error) return <ErrorState message={error} onRetry={load} />;
  if (loading) return <Loader />;
  if (!project) return <ErrorState message="Projeto não encontrado" onRetry={load} />;

  const statusVariant = project.divergence_status === "sincronizado" ? "success" : "warning";

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <button onClick={() => navigate("/projects")} className="text-gray-400 hover:text-gray-600">&larr;</button>
        <h2 className="text-2xl font-bold text-gray-800">{project.name}</h2>
        <Badge variant={statusVariant}>{project.divergence_status}</Badge>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <Card title="Informações">
          <div className="space-y-2 text-sm">
            <p><span className="text-gray-500">Path:</span> {project.local_path}</p>
            <p><span className="text-gray-500">Branch:</span> {project.branch}</p>
            <p><span className="text-gray-500">Ambiente:</span> {project.environment}</p>
            <p><span className="text-gray-500">Remote:</span> {project.remote_path}</p>
          </div>
        </Card>

        <Card title="Ações Rápidas">
          <div className="flex flex-wrap gap-2">
            <button onClick={() => setActionModal("deploy")} className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 transition-colors">Deploy</button>
            <button onClick={() => setActionModal("rollback")} className="px-4 py-2 bg-amber-600 text-white text-sm rounded-lg hover:bg-amber-700 transition-colors">Rollback</button>
            <button onClick={() => setActionModal("verify")} className="px-4 py-2 bg-emerald-600 text-white text-sm rounded-lg hover:bg-emerald-700 transition-colors">Verificar</button>
            <button onClick={handleLoadDiff} disabled={diffLoading} className="px-4 py-2 bg-violet-600 text-white text-sm rounded-lg hover:bg-violet-700 transition-colors">
              {diffLoading ? "Carregando..." : "Diff"}
            </button>
          </div>
        </Card>

        <Card title="Servidores">
          {servers.length === 0 ? (
            <p className="text-sm text-gray-400">Nenhum servidor</p>
          ) : (
            <div className="space-y-1">
              {servers.map((s) => (
                <p key={s.id} className="text-sm text-gray-700">{s.name} ({s.host})</p>
              ))}
            </div>
          )}
        </Card>
      </div>

      <Card title="Últimos Deploys">
        <Table
          columns={[
            { key: "id", header: "#" },
            { key: "status", header: "Status", render: (v: string) => <Badge variant={v === "success" ? "success" : v === "running" ? "info" : "danger"}>{v}</Badge> },
            { key: "version", header: "Versão" },
            { key: "created_at", header: "Data" },
          ]}
          data={deploys}
        />
      </Card>

      <Modal open={actionModal !== null} onClose={() => setActionModal(null)} title={`Confirmar ${actionModal}`}>
        <p className="text-sm text-gray-600 mb-4">Tem certeza que deseja executar <strong>{actionModal}</strong> no projeto <strong>{name}</strong>?</p>
        <div className="flex gap-3 justify-end">
          <button onClick={() => setActionModal(null)} className="px-4 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">Cancelar</button>
          <button onClick={() => handleAction(actionModal!)} disabled={running} className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 disabled:opacity-50">
            {running ? "Executando..." : "Confirmar"}
          </button>
        </div>
      </Modal>

      <Modal open={diffModal} onClose={() => setDiffModal(false)} title={`Diff: ${name}`}>
        {diffData ? (
          <div className="space-y-3">
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Status:</span>
              <Badge variant={diffData.status === "identical" ? "success" : "warning"}>{diffData.status}</Badge>
            </div>
            {diffData.files && (
              <div>
                <p className="text-sm font-medium text-gray-700 mb-2">Arquivos ({diffData.files.length})</p>
                <div className="max-h-60 overflow-y-auto space-y-1">
                  {diffData.files.map((f: any, i: number) => (
                    <div key={i} className="text-xs text-gray-600 flex justify-between">
                      <span>{f.path}</span>
                      <Badge variant={f.status === "identical" ? "success" : f.status === "modified" ? "warning" : "danger"}>{f.status}</Badge>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : <p className="text-sm text-gray-400">Nenhum dado de diff.</p>}
      </Modal>

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
