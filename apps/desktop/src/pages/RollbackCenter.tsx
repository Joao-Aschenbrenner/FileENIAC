import { useEffect, useState } from "react";
import { listProjects, executeRollback } from "../api/client";
import { Card } from "../components/ui/Card";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Toast } from "../components/ui/Toast";

export default function RollbackCenter() {
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [rolling, setRolling] = useState(false);
  const [error, setError] = useState("");
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" } | null>(null);
  const [confirmOpen, setConfirmOpen] = useState(false);

  useEffect(() => {
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    listProjects(wsPath)
      .then(setProjects)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  async function handleRollback() {
    if (!selectedProject) return;
    setRolling(true);
    try {
      await executeRollback(selectedProject);
      setToast({ message: "Rollback executado com sucesso", type: "success" });
      setConfirmOpen(false);
    } catch (e: any) {
      setToast({ message: e.message, type: "error" });
    }
    setRolling(false);
  }

  if (loading) return <Loader />;
  if (error) return <ErrorState message={error} onRetry={() => window.location.reload()} />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Rollback Center</h2>
      </div>

      <div className="max-w-lg mx-auto space-y-4">
        <Card title="Selecionar Projeto">
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

        {selectedProject && (
          <Card title="Confirmação">
            <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 mb-4">
              <p className="text-sm text-amber-800 font-medium">Atenção</p>
              <p className="text-xs text-amber-700 mt-1">O rollback irá reverter o projeto <strong>{selectedProject}</strong> para a versão anterior. Esta ação pode causar indisponibilidade temporária.</p>
            </div>
            <button
              onClick={() => setConfirmOpen(true)}
              disabled={rolling}
              className="w-full py-2 bg-amber-600 text-white text-sm rounded-lg hover:bg-amber-700 disabled:opacity-50 transition-colors"
            >
              {rolling ? "Executando Rollback..." : "Executar Rollback"}
            </button>
          </Card>
        )}

        {confirmOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div className="absolute inset-0 bg-black/40" onClick={() => setConfirmOpen(false)} />
            <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
              <h3 className="font-semibold text-gray-800 mb-2">Confirmar Rollback</h3>
              <p className="text-sm text-gray-600 mb-4">Tem certeza que deseja reverter <strong>{selectedProject}</strong>?</p>
              <div className="flex gap-3">
                <button onClick={() => setConfirmOpen(false)} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">Cancelar</button>
                <button onClick={handleRollback} disabled={rolling} className="flex-1 py-2 bg-amber-600 text-white text-sm rounded-lg hover:bg-amber-700 disabled:opacity-50">
                  {rolling ? "Executando..." : "Confirmar"}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
