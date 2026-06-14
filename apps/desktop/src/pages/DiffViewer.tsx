import { useEffect, useState } from "react";
import { listProjects, getDiff } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Table } from "../components/ui/Table";
import { Toast } from "../components/ui/Toast";

export default function DiffViewer() {
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState<string>("");
  const [diff, setDiff] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" | "info" } | null>(null);

  useEffect(() => {
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    listProjects(wsPath)
      .then(setProjects)
      .catch(() => {});
  }, []);

  async function handleLoadDiff() {
    if (!selectedProject) return;
    setLoading(true);
    try {
      const data = await getDiff(selectedProject);
      setDiff(data);
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setLoading(false);
  }

  const fileColumns = [
    { key: "path", header: "Arquivo" },
    { key: "status", header: "Status", render: (v: string) => {
      const variant: Record<string, "success" | "warning" | "danger" | "neutral"> = {
        identical: "success", modified: "warning", added: "neutral", removed: "danger",
      };
      return <Badge variant={variant[v] || "neutral"}>{v}</Badge>;
    }},
    { key: "local_hash", header: "Hash Local", className: "font-mono text-xs" },
    { key: "mirror_hash", header: "Hash Mirror", className: "font-mono text-xs" },
  ];

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Diff Viewer</h2>
      </div>

      <div className="grid grid-cols-4 gap-4 mb-6">
        <Card title="Projeto">
          <select
            value={selectedProject}
            onChange={(e) => setSelectedProject(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
          >
            <option value="">Selecione...</option>
            {projects.filter(p => p.divergence_status !== "sincronizado").length === 0 && (
              <option value="" disabled>Todos sincronizados</option>
            )}
            {projects.map((p) => (
              <option key={p.id} value={p.name}>{p.name} {p.divergence_status !== "sincronizado" ? "⚠️" : ""}</option>
            ))}
          </select>
        </Card>

        <Card title="Status">
          {diff ? (
            <Badge variant={diff.status === "identical" ? "success" : "warning"}>{diff.status}</Badge>
          ) : (
            <p className="text-sm text-gray-400">-</p>
          )}
        </Card>

        <Card title="Divergentes">
          <p className="text-2xl font-bold">{diff ? diff.files?.filter((f: any) => f.status !== "identical").length || 0 : "-"}</p>
        </Card>

        <Card title="Ação">
          <button
            onClick={handleLoadDiff}
            disabled={!selectedProject || loading}
            className="w-full py-2 bg-violet-600 text-white text-sm rounded-lg hover:bg-violet-700 disabled:opacity-50 transition-colors"
          >
            {loading ? "Carregando..." : "Carregar Diff"}
          </button>
        </Card>
      </div>

      {diff && (
        <Card title="Arquivos">
          <Table columns={fileColumns} data={diff.files || []} />
        </Card>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
