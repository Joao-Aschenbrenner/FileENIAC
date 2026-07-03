// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { listProjects, deleteProject } from "../api/client";
import { isProjectDirectory } from "../lib/projectUtils";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";

export default function Projects() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);
  const [deleteFiles, setDeleteFiles] = useState(false);

  function loadProjects() {
    setLoading(true);
    setError("");
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    listProjects(wsPath)
      .then((data) => setProjects(Array.isArray(data) ? data.filter((p) => !isProjectDirectory(p.name)) : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    loadProjects();
  }, []);

  if (loading) return <Loader text="Carregando projetos..." />;
  if (error) return <ErrorState message={error} onRetry={loadProjects} />;

  async function handleDelete(name: string) {
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    try {
      await deleteProject(wsPath, name);
      setDeleteConfirm(null);
      setDeleteFiles(false);
      loadProjects();
    } catch (e: any) {
      setError(e.message);
      setDeleteConfirm(null);
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Projetos</h2>
        <button
          onClick={() => navigate("/github/orgs")}
          className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors"
        >
          + Adicionar Repositórios
        </button>
      </div>

      {error && (
        <div className="mb-4 text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      {projects.length === 0 ? (
        <div className="bg-white rounded-xl border border-gray-200 p-12 text-center">
          <h3 className="text-lg font-semibold text-gray-700 mb-2">Nenhum repositorio adicionado</h3>
          <p className="text-sm text-gray-500 mb-6">
            Importe repositorios do GitHub para comecar a gerenciar seus projetos.
          </p>
          <button
            onClick={() => navigate("/github/orgs")}
            className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors"
          >
            Adicionar Repositórios
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          {projects.map((p) => (
            <div
              key={p.id}
              className="bg-white rounded-xl border border-gray-200 p-4 shadow-sm flex items-center justify-between"
            >
              <div>
                <p className="font-semibold text-gray-800">{p.name}</p>
                <p className="text-sm text-gray-500">
                  {p.local_path} · {p.environment}
                </p>
              </div>
              <div className="flex items-center gap-3">
                <span
                  className={`text-xs px-2 py-1 rounded-full font-medium ${
                    p.divergence_status === "sincronizado"
                      ? "bg-green-100 text-green-700"
                      : "bg-amber-100 text-amber-700"
                  }`}
                >
                  {p.divergence_status}
                </span>
                <button
                  onClick={() => setDeleteConfirm(p.name)}
                  className="text-xs text-red-600 hover:text-red-800 transition-colors"
                >
                  Remover
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {deleteConfirm !== null && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => { setDeleteConfirm(null); setDeleteFiles(false); }} />
          <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
            <h3 className="font-semibold text-gray-800 mb-2">Remover Projeto</h3>
            <p className="text-sm text-gray-600 mb-4">
              Isso remove o projeto &ldquo;{deleteConfirm}&rdquo; da lista do FileENIAC, mas nao apaga os arquivos locais.
            </p>
            <label className="flex items-center gap-2 mb-4 text-sm text-gray-600 cursor-pointer">
              <input
                type="checkbox"
                checked={deleteFiles}
                onChange={(e) => setDeleteFiles(e.target.checked)}
                className="rounded border-gray-300 text-red-600 focus:ring-red-500"
              />
              Tambem apagar arquivos locais
            </label>
            <div className="flex gap-3">
              <button onClick={() => { setDeleteConfirm(null); setDeleteFiles(false); }} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">Cancelar</button>
              <button onClick={() => handleDelete(deleteConfirm)} className="flex-1 py-2 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700">Remover</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
