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
  const [isRemoving, setIsRemoving] = useState(false);
  const [success, setSuccess] = useState("");

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

  useEffect(() => {
    if (!success) return;
    const t = setTimeout(() => setSuccess(""), 4000);
    return () => clearTimeout(t);
  }, [success]);

  if (loading) return <Loader text="Carregando projetos..." />;
  if (error) return <ErrorState message={error} onRetry={loadProjects} />;

  function openDeleteModal(name: string) {
    setDeleteConfirm(name);
    setDeleteFiles(false);
    setError("");
  }

  function closeDeleteModal() {
    if (isRemoving) return;
    setDeleteConfirm(null);
    setDeleteFiles(false);
  }

  async function handleDelete() {
    if (!deleteConfirm) return;
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    const name = deleteConfirm;
    setIsRemoving(true);
    setError("");
    try {
      const res = await deleteProject(wsPath, name, { deleteLocalFiles: deleteFiles });
      const localResult = res?.local_files ?? "skipped";
      setDeleteConfirm(null);
      setDeleteFiles(false);
      loadProjects();
      if (deleteFiles && localResult === "deleted") {
        setSuccess(`Projeto removido e pasta local apagada.`);
      } else if (deleteFiles && String(localResult).startsWith("failed")) {
        setSuccess(`Projeto removido do workspace, mas a pasta local nao pode ser apagada.`);
      } else {
        setSuccess(`Projeto removido do workspace.`);
      }
    } catch (e: any) {
      const msg = e?.message ?? "Erro ao remover projeto";
      if (msg.includes("FOREIGN KEY")) {
        setError("Nao consegui remover este projeto. Tente novamente.");
      } else {
        setError(msg);
      }
    } finally {
      setIsRemoving(false);
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

      {success && (
        <div className="mb-4 bg-green-50 border border-green-200 text-green-800 rounded-lg p-3 text-sm">
          {success}
        </div>
      )}

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
                  onClick={() => openDeleteModal(p.name)}
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
          <div className="absolute inset-0 bg-black/40" onClick={closeDeleteModal} />
          <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
            <h3 className="font-semibold text-gray-800 mb-2">Remover projeto do workspace?</h3>
            <p className="text-sm text-gray-600 mb-4">
              Isso remove o projeto &ldquo;{deleteConfirm}&rdquo; da lista do FileENIAC. O repositorio no
              GitHub/GitLab nao sera apagado.
            </p>
            <label className="flex items-start gap-2 mb-2 text-sm text-gray-600 cursor-pointer">
              <input
                type="checkbox"
                checked={deleteFiles}
                onChange={(e) => setDeleteFiles(e.target.checked)}
                disabled={isRemoving}
                className="mt-0.5 rounded border-gray-300 text-red-600 focus:ring-red-500"
              />
              <span>Tambem apagar a pasta local deste projeto</span>
            </label>
            {deleteFiles && (
              <p className="text-xs text-amber-700 bg-amber-50 border border-amber-200 rounded p-2 mb-4">
                Esta acao apagara os arquivos locais do projeto neste computador. O repositorio remoto no
                GitHub/GitLab nao sera apagado.
              </p>
            )}
            {!deleteFiles && <div className="mb-4" />}
            <div className="flex gap-3">
              <button
                onClick={closeDeleteModal}
                disabled={isRemoving}
                className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50 disabled:opacity-50"
              >
                Cancelar
              </button>
              <button
                onClick={handleDelete}
                disabled={isRemoving}
                className="flex-1 py-2 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700 disabled:opacity-50"
              >
                {isRemoving
                  ? "Removendo..."
                  : deleteFiles
                  ? "Remover e apagar pasta local"
                  : "Remover do workspace"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
