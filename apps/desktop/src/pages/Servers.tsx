// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { listServers, deleteServer, createServer, listProjects } from "../api/client";
import { Loader } from "../components/ui/Loader";

export default function Servers() {
  const [servers, setServers] = useState<any[]>([]);
  const [projects, setProjects] = useState<any[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState<number | null>(null);
  const [creating, setCreating] = useState(false);
  const [form, setForm] = useState({
    project_id: 0,
    name: "",
    host: "",
    port: 21,
    user: "",
    password: "",
    target_path: "/",
  });

  function loadData() {
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    setLoading(true);
    setError("");
    Promise.all([
      listServers(wsPath),
      listProjects(wsPath),
    ])
      .then(([s, p]) => { setServers(s); setProjects(p); })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    loadData();
  }, []);
  async function handleCreate() {
    if (creating) return;
    setCreating(true);
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    try {
      await createServer(wsPath, form);
      setShowForm(false);
      setForm({ project_id: 0, name: "", host: "", port: 21, user: "", password: "", target_path: "/" });
      loadData();
    } catch (e: any) {
      setError(e.message);
    }
    setCreating(false);
  }

  async function handleDelete(id: number) {
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    try {
      await deleteServer(wsPath, id);
      setDeleteConfirm(null);
      loadData();
    } catch (e: any) {
      setError(e.message);
      setDeleteConfirm(null);
    }
  }

  if (loading) return <Loader />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Servidores</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors"
        >
          {showForm ? "Cancelar" : "+ Novo Servidor"}
        </button>
      </div>

      {error && (
        <div className="mb-4 text-red-600 bg-red-50 border border-red-200 rounded-lg p-3 text-sm">
          {error}
        </div>
      )}

      {showForm && (
        <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm mb-6">
          <h3 className="font-semibold text-gray-700 mb-4">Novo Servidor</h3>
          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Projeto</label>
              <select value={form.project_id} onChange={(e) => setForm({ ...form, project_id: +e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500">
                <option value={0}>Selecione...</option>
                {projects.map((p) => (
                  <option key={p.id} value={p.id}>{p.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Nome</label>
              <input type="text" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Host</label>
              <input type="text" value={form.host} onChange={(e) => setForm({ ...form, host: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Porta</label>
              <input type="number" value={form.port} onChange={(e) => setForm({ ...form, port: +e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Usuário</label>
              <input type="text" value={form.user} onChange={(e) => setForm({ ...form, user: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Senha</label>
              <input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500" />
            </div>
            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-600 mb-1">Target Path</label>
              <input type="text" value={form.target_path} onChange={(e) => setForm({ ...form, target_path: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500" />
            </div>
          </div>
          <button onClick={handleCreate}
            className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Salvar Servidor
          </button>
        </div>
      )}

      {servers.length === 0 ? (
        <div className="bg-white rounded-xl border border-gray-200 p-8 text-center">
          <p className="text-gray-500 text-sm">Nenhum servidor cadastrado.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {servers.map((s) => {
            const proj = projects.find((p) => p.id === s.project_id);
            return (
              <div key={s.id} className="bg-white rounded-xl border border-gray-200 p-4 shadow-sm flex items-center justify-between">
                <div>
                  <p className="font-semibold text-gray-800">{s.name}</p>
                  <p className="text-sm text-gray-500">{s.host}:{s.port} · {proj ? proj.name : `Projeto #${s.project_id}`}</p>
                </div>
                <button onClick={() => setDeleteConfirm(s.id)}
                  className="text-xs text-red-600 hover:text-red-800 transition-colors">
                  Remover
                </button>
              </div>
            );
          })}
        </div>
      )}

      {deleteConfirm !== null && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setDeleteConfirm(null)} />
          <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
            <h3 className="font-semibold text-gray-800 mb-2">Remover Servidor</h3>
            <p className="text-sm text-gray-600 mb-4">Tem certeza que deseja remover este servidor?</p>
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
