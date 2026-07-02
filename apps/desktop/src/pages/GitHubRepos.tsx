// SPDX-License-Identifier: MIT
import { useEffect, useState, useCallback } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { getGitHubRepositories, importGitHubRepos, TimeoutError } from "../api/client";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Badge } from "../components/ui/Badge";
import { Toast } from "../components/ui/Toast";

export default function GitHubRepos() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const org = searchParams.get("org") || "";
  const [repos, setRepos] = useState<any[]>([]);
  const [selected, setSelected] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(true);
  const [importing, setImporting] = useState(false);
  const [error, setError] = useState("");
  const [toast, setToast] = useState<{ message: string; type: "success" | "error" } | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError("");
    getGitHubRepositories(org || undefined)
      .then((data) => setRepos(Array.isArray(data) ? data : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [org]);

  useEffect(() => { load() }, [load]);

  function toggle(id: number) {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  async function handleImport() {
    const toImport = repos.filter((r) => selected.has(r.id));
    if (toImport.length === 0) return;

    setImporting(true);
    setError("");
    try {
      const results = await importGitHubRepos(toImport);
      const success = results.filter((r: any) => !r.error);
      const failed = results.filter((r: any) => r.error);
      if (failed.length > 0) {
        setToast({ message: `${success.length} importados, ${failed.length} falhas`, type: "error" });
      } else {
        setToast({ message: `${success.length} projetos importados com sucesso!`, type: "success" });
        setSelected(new Set());
        load();
      }
    } catch (e: any) {
      if (e instanceof TimeoutError) {
        setError("A importação demorou mais que o esperado. Pode ser que alguns projetos já tenham sido importados. Clique em \"Tentar novamente\" para verificar.");
      } else {
        setError(e.message);
      }
    }
    setImporting(false);
  }

  if (loading) return <Loader text="Carregando repositórios..." />;
  if (error) return <ErrorState message={error} onRetry={load} />;

  const imported = repos.filter((r) => r.imported);
  const unimported = repos.filter((r) => !r.imported);

  return (
    <div className="max-w-2xl mx-auto mt-8">
      <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-xl font-bold text-gray-800">
              {org ? `Repositórios: ${org}` : "Meus Repositórios"}
            </h2>
            <p className="text-sm text-gray-500 mt-0.5">
              {unimported.length} disponíveis · {imported.length} já importados · {selected.size} selecionados
            </p>
          </div>
          <div className="flex gap-2">
            <button onClick={() => navigate(org ? "/github/orgs" : "/bootstrap")}
              className="px-3 py-1.5 border border-gray-300 text-sm rounded-lg hover:bg-gray-50 transition-colors">
              Voltar
            </button>
            <button onClick={handleImport} disabled={selected.size === 0 || importing}
              className="px-4 py-1.5 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 disabled:opacity-50 transition-colors">
              {importing ? "Importando..." : `Importar (${selected.size})`}
            </button>
          </div>
        </div>

        {repos.length === 0 && (
          <p className="text-center text-sm text-gray-400 py-8">Nenhum repositório encontrado.</p>
        )}

        <div className="space-y-1 max-h-96 overflow-y-auto">
          {unimported.map((repo: any) => (
            <label key={repo.id}
              className={`flex items-center gap-3 px-4 py-2.5 border rounded-lg cursor-pointer transition-all ${
                selected.has(repo.id) ? "border-eniac-400 bg-eniac-50" : "border-gray-200 hover:bg-gray-50"
              }`}
            >
              <input
                type="checkbox"
                checked={selected.has(repo.id)}
                onChange={() => toggle(repo.id)}
                className="rounded border-gray-300 text-eniac-600 focus:ring-eniac-500"
              />
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <p className="font-medium text-gray-800 truncate">{repo.name}</p>
                  {repo.private && <Badge variant="warning">Privado</Badge>}
                </div>
                <p className="text-xs text-gray-500 truncate">{repo.description || "Sem descrição"}</p>
                <p className="text-xs text-gray-400 mt-0.5 font-mono">{repo.default_branch} · {repo.language || "N/A"}</p>
              </div>
            </label>
          ))}
          {imported.map((repo: any) => (
            <div key={repo.id}
              className="flex items-center gap-3 px-4 py-2.5 border border-gray-100 rounded-lg bg-gray-50 opacity-60"
            >
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <p className="font-medium text-gray-600 truncate">{repo.name}</p>
                  <Badge variant="success">Importado</Badge>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
