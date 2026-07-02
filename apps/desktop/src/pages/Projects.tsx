// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { listProjects, deleteProject, createProject } from "../api/client";

import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";

export default function Projects() {
  const [projects, setProjects] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [creating, setCreating] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);
  const [form, setForm] = useState({ name: "", local_path: "", remote_path: "/", branch: "main" });

  function loadProjects() {
    setLoading(true);
    setError("");
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    listProjects(wsPath)
      .then(setProjects)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    loadProjects();
  }, []);

  if (loading) return <Loader text="Carregando projetos..." />;
  if (error) return <ErrorState message={error} onRetry={loadProjects} />;

  async function handleCreate() {
    if (creating) return;
    setCreating(true);
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    try {
      await createProject(wsPath, form);
      setShowForm(false);
      setForm({ name: "", local_path: "", remote_path: "/", branch: "main" });
      loadProjects();
    } catch (e: any) {
      setError(e.message);
    }
    setCreating(false);
  }

  async function handleDelete(name: string) {
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    try {
      await deleteProject(wsPath, name);
      setDeleteConfirm(null);
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
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors"
        >
          {showForm ? "Cancelar" : "+ Novo Projeto"}
        </button>
      </div>

      {error && (
        <div className="mb-4 text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      {showForm && (
        <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm mb-6">
          <h3 className="font-semibold text-gray-700 mb-4">Novo Projeto</h3>
          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Nome</label>
              <input
                type="text" value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Caminho Local</label>
              <input
                type="text" value={form.local_path}
                onChange={(e) => setForm({ ...form, local_path: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Caminho Remoto</label>
              <input
                type="text" value={form.remote_path}
                onChange={(e) => setForm({ ...form, remote_path: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Branch</label>
              <input
                type="text" value={form.branch}
                onChange={(e) => setForm({ ...form, branch: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
              />
            </div>
          </div>
          <button
            onClick={handleCreate}
            disabled={creating}
            className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors disabled:opacity-50"
          >
            {creating ? "Salvando..." : "Salvar Projeto"}
          </button>
        </div>
      )}

      {projects.length === 0 ? (
        <div className="bg-white rounded-xl border border-gray-200 p-8 text-center">
          <p className="text-gray-500 text-sm">Nenhum projeto cadastrado.</p>
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
          <div className="absolute inset-0 bg-black/40" onClick={() => setDeleteConfirm(null)} />
          <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
            <h3 className="font-semibold text-gray-800 mb-2">Remover Projeto</h3>
            <p className="text-sm text-gray-600 mb-4">Tem certeza que deseja remover o projeto "{deleteConfirm}"? Servidores e histórico associados também serão removidos.</p>
            <div className="flex gap-3">
              <button onClick={() => setDeleteConfirm(null)} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">Cancelar</button>
              <button onClick={() => handleDelete(deleteConfirm)} className="flex-1 py-2 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700">Remover</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
