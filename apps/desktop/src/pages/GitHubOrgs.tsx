import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getGitHubOrganizations } from "../api/client";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";

export default function GitHubOrgs() {
  const navigate = useNavigate();
  const [orgs, setOrgs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  function load() {
    setLoading(true);
    setError("");
    getGitHubOrganizations()
      .then((data) => setOrgs(Array.isArray(data) ? data : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => { load() }, []);

  if (loading) return <Loader text="Carregando organizações..." />;
  if (error) return <ErrorState message={error} onRetry={load} />;

  return (
    <div className="max-w-lg mx-auto mt-12">
      <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-8">
        <div className="text-center mb-6">
          <div className="text-4xl mb-2">🏢</div>
          <h2 className="text-xl font-bold text-gray-800">Selecionar Organização</h2>
          <p className="text-sm text-gray-500 mt-1">Escolha a organização para importar repositórios</p>
        </div>

        {orgs.length === 0 && (
          <p className="text-center text-sm text-gray-400 py-4">Nenhuma organização encontrada.</p>
        )}

        <div className="space-y-2 mb-6">
          {orgs.map((org) => (
            <button
              key={org.login}
              onClick={() => navigate(`/github/repos?org=${encodeURIComponent(org.login)}`)}
              className="w-full flex items-center gap-3 px-4 py-3 border border-gray-200 rounded-lg hover:bg-gray-50 hover:border-eniac-300 transition-all text-left"
            >
              <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center text-sm font-bold text-gray-600">
                {org.login.charAt(0).toUpperCase()}
              </div>
              <div>
                <p className="font-medium text-gray-800">{org.login}</p>
                <p className="text-xs text-gray-500">{org.url ? org.url.replace("https://api.github.com/orgs/", "github.com/") : org.login}</p>
              </div>
            </button>
          ))}
        </div>

        <button onClick={() => navigate("/github/repos")}
          className="w-full py-2 text-sm text-eniac-600 hover:text-eniac-800 font-medium">
          Ver meus repositórios pessoais
        </button>
      </div>
    </div>
  );
}
