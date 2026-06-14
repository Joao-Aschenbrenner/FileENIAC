import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { gitHubLogin } from "../api/client";

export default function GitHubLogin() {
  const navigate = useNavigate();
  const [token, setToken] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleLogin() {
    if (!token.trim()) {
      setError("Token é obrigatório");
      return;
    }
    setLoading(true);
    setError("");
    try {
      const result = await gitHubLogin(token.trim());
      localStorage.setItem("github_user", result.user);
      navigate("/github/orgs");
    } catch (e: any) {
      setError(e.message);
    }
    setLoading(false);
  }

  return (
    <div className="max-w-md mx-auto mt-12">
      <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-8">
        <div className="text-center mb-6">
          <div className="text-4xl mb-2">🐙</div>
          <h2 className="text-xl font-bold text-gray-800">Conectar GitHub</h2>
          <p className="text-sm text-gray-500 mt-1">Informe seu Personal Access Token do GitHub</p>
        </div>

        <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 mb-6 text-xs text-gray-600 space-y-1">
          <p className="font-medium text-gray-700">Como gerar um token:</p>
          <p>1. Acesse github.com/settings/tokens</p>
          <p>2. Clique em "Generate new token (classic)"</p>
          <p>3. Selecione escopos: <code className="bg-gray-200 px-1 rounded">repo</code>, <code className="bg-gray-200 px-1 rounded">read:org</code></p>
          <p>4. Copie o token e cole abaixo</p>
        </div>

        <label className="block text-sm font-medium text-gray-700 mb-1">Personal Access Token</label>
        <input
          type="password"
          value={token}
          onChange={(e) => setToken(e.target.value)}
          placeholder="ghp_xxxxxxxxxxxxxxxxxxxx"
          className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500 font-mono"
        />

        {error && <p className="mt-2 text-sm text-red-600">{error}</p>}

        <div className="flex gap-3 mt-6">
          <button onClick={() => navigate("/dashboard")} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50 transition-colors">
            Cancelar
          </button>
          <button onClick={handleLogin} disabled={loading}
            className="flex-1 py-2 bg-eniac-600 text-white text-sm rounded-lg font-medium hover:bg-eniac-700 disabled:opacity-50 transition-colors">
            {loading ? "Validando..." : "Conectar"}
          </button>
        </div>
      </div>
    </div>
  );
}
