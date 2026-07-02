// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { listProjects, getDiff, getSyncs, executeSync } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Table } from "../components/ui/Table";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Toast } from "../components/ui/Toast";

export default function SyncCenter() {
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState<string>("");
  const [diff, setDiff] = useState<any>(null);
  const [syncs, setSyncs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncLoading, setSyncLoading] = useState(false);
  const [syncing, setSyncing] = useState(false);
  const [error, setError] = useState("");
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" } | null>(null);
  const [confirmOpen, setConfirmOpen] = useState(false);

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

  if (loading) return <Loader text="Carregando projetos..." />;
  if (error) return <ErrorState message={error} onRetry={loadProjects} />;

  async function handleAnalyze() {
    if (!selectedProject) return;
    setSyncLoading(true);
    try {
      const [diffData, syncData] = await Promise.all([
        getDiff(selectedProject),
        getSyncs(selectedProject, 10),
      ]);
      setDiff(diffData);
      setSyncs(syncData);
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setSyncLoading(false);
  }

  async function handleSync() {
    if (!selectedProject) return;
    setSyncing(true);
    try {
      const result = await executeSync(selectedProject, "mirror_update");
      setToast({ message: `Sync concluído: ${result.manifest?.result || "ok"}`, type: "success" });
      setConfirmOpen(false);
      handleAnalyze();
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setSyncing(false);
  }

  const divergents = projects.filter((p) => p.divergence_status !== "sincronizado");

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Sync Center</h2>
        {divergents.length > 0 && (
          <Badge variant="warning">{divergents.length} projetos divergentes</Badge>
        )}
      </div>

      <div className="grid grid-cols-4 gap-4 mb-6">
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

        <Card title="Status Diff">
          {diff ? <Badge variant={diff.status === "identical" ? "success" : "warning"}>{diff.status}</Badge> : <p className="text-sm text-gray-400">-</p>}
        </Card>

        <Card title="Arquivos">
          <p className="text-2xl font-bold">{diff ? diff.files?.filter((f: any) => f.status !== "identical").length || 0 : "-"}</p>
        </Card>

        <Card title="Ação">
          <button
            onClick={handleAnalyze}
            disabled={!selectedProject || syncLoading}
            className="w-full py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 disabled:opacity-50 transition-colors"
          >
            {syncLoading ? "Analisando..." : "Analisar"}
          </button>
        </Card>
      </div>

      {diff && diff.status !== "identical" && (
        <Card title="Diferenças Detectadas" className="mb-4">
          <p className="text-sm text-gray-600 mb-3">{diff.files?.filter((f: any) => f.status !== "identical").length || 0} arquivo(s) divergente(s) encontrado(s).</p>
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-3 mb-4">
            <p className="text-xs text-amber-700">Sync irá atualizar o mirror com as alterações locais. Esta ação NÃO altera o servidor remoto.</p>
          </div>
          <button
            onClick={() => setConfirmOpen(true)}
            disabled={syncing}
            className="px-4 py-2 bg-emerald-600 text-white text-sm rounded-lg hover:bg-emerald-700 disabled:opacity-50 transition-colors"
          >
            {syncing ? "Sincronizando..." : "Executar Sync"}
          </button>
        </Card>
      )}

      {syncs.length > 0 && (
        <Card title="Histórico de Syncs">
          <Table
            columns={[
              { key: "id", header: "#" },
              { key: "manifest_id", header: "Manifest" },
              { key: "operation_type", header: "Operação" },
              { key: "files_count", header: "Arquivos" },
              { key: "result", header: "Resultado", render: (v: string) => <Badge variant={v === "completed" ? "success" : v === "failed" ? "danger" : "info"}>{v}</Badge> },
              { key: "created_at", header: "Data" },
            ]}
            data={syncs}
          />
        </Card>
      )}

      {confirmOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setConfirmOpen(false)} />
          <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
            <h3 className="font-semibold text-gray-800 mb-2">Confirmar Sync</h3>
            <p className="text-sm text-gray-600 mb-4">Sincronizar <strong>{selectedProject}</strong>? O mirror será atualizado.</p>
            <div className="flex gap-3">
              <button onClick={() => setConfirmOpen(false)} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">Cancelar</button>
              <button onClick={handleSync} disabled={syncing} className="flex-1 py-2 bg-emerald-600 text-white text-sm rounded-lg hover:bg-emerald-700 disabled:opacity-50">
                {syncing ? "Executando..." : "Confirmar Sync"}
              </button>
            </div>
          </div>
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
