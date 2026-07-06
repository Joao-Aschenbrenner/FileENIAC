// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getGitHubStatus, listServers } from "../api/client";
import { IconGithub, IconGitlab, IconServer, IconSwitchWorkspace, IconCheck, IconWarning, IconPlus, IconExternalLink } from "../components/Icons";
import { Loader } from "../components/ui/Loader";
import { Badge } from "../components/ui/Badge";

export default function WorkspaceBootstrap() {
  const navigate = useNavigate();
  const [status, setStatus] = useState<any>(null);
  const [servers, setServers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showSwitchModal, setShowSwitchModal] = useState(false);
  const wsPath = localStorage.getItem("eniac_ws_path") || "";

  useEffect(() => {
    const activeWorkspace = localStorage.getItem("eniac_ws_path") || "";
    if (!activeWorkspace) {
      setLoading(false);
      return;
    }
    Promise.all([
      getGitHubStatus(),
      listServers(activeWorkspace).catch(() => []),
    ])
      .then(([s, srv]) => {
        setStatus(s);
        setServers(Array.isArray(srv) ? srv : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <Loader text="Verificando ambiente..." />;

  if (!wsPath) {
    return (
      <div className="max-w-xl mx-auto mt-8">
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-8 text-center">
          <h2 className="text-2xl font-bold text-gray-800">Nenhum workspace selecionado</h2>
          <p className="text-sm text-gray-500 mt-2">
            Escolha ou crie um workspace antes de configurar o ambiente.
          </p>
          <button onClick={() => navigate("/")}
            className="mt-6 px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors">
            Ir para Workspaces
          </button>
        </div>
      </div>
    );
  }

  const githubConnected = status?.authenticated;
  const serversCount = servers.length;

  return (
    <div className="max-w-2xl mx-auto">
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold text-gray-800">Configurar Ambiente</h2>
        <p className="text-sm text-gray-500 mt-1">Workspace, integracoes e servidores</p>
      </div>

      {/* Workspace Section */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm mb-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center">
              <IconServer className="text-gray-600 w-5 h-5" />
            </div>
            <div>
              <p className="font-medium text-gray-800">Workspace</p>
              <p className="text-xs text-gray-500 break-all">{wsPath}</p>
            </div>
          </div>
          <button
            onClick={() => setShowSwitchModal(true)}
            className="flex items-center gap-2 px-3 py-1.5 text-xs border border-gray-300 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors"
          >
            <IconSwitchWorkspace className="w-4 h-4" />
            Trocar workspace
          </button>
        </div>
      </div>

      {/* Integrations Section */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm mb-4">
        <h3 className="font-semibold text-gray-700 mb-4">Integracoes</h3>
        <div className="grid grid-cols-1 gap-3">

          {/* GitHub Card */}
          <div className="border border-gray-200 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 bg-gray-900 rounded-lg flex items-center justify-center">
                  <IconGithub className="text-white w-5 h-5" />
                </div>
                <div>
                  <p className="font-medium text-gray-800">GitHub</p>
                  {githubConnected ? (
                    <p className="text-xs text-gray-500">Conectado como {status.user}</p>
                  ) : (
                    <p className="text-xs text-gray-400">Nao conectado</p>
                  )}
                </div>
              </div>
              {githubConnected ? (
                <div className="flex items-center gap-2">
                  <Badge variant="success" className="text-xs">Conectado</Badge>
                  <button
                    onClick={() => navigate("/github/login")}
                    className="text-xs text-gray-500 hover:text-gray-700 underline"
                  >
                    Desconectar
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => navigate("/github/login")}
                  className="px-3 py-1.5 bg-eniac-600 text-white text-xs rounded-lg hover:bg-eniac-700 transition-colors"
                >
                  Conectar
                </button>
              )}
            </div>
            {githubConnected && (
              <div className="mt-3 pt-3 border-t border-gray-100 flex items-center gap-4">
                <button
                  onClick={() => navigate("/github/orgs")}
                  className="flex items-center gap-1 text-xs text-eniac-600 hover:text-eniac-700"
                >
                  <IconExternalLink className="w-3 h-3" />
                  Ver organizacoes
                </button>
                <button
                  onClick={() => navigate("/github/repos")}
                  className="flex items-center gap-1 text-xs text-eniac-600 hover:text-eniac-700"
                >
                  <IconExternalLink className="w-3 h-3" />
                  Meus repositorios
                </button>
              </div>
            )}
          </div>

          {/* GitLab Card */}
          <div className="border border-gray-200 rounded-lg p-4 bg-gray-50 opacity-75">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 bg-orange-500 rounded-lg flex items-center justify-center">
                  <IconGitlab className="text-white w-5 h-5" />
                </div>
                <div>
                  <p className="font-medium text-gray-600">GitLab</p>
                  <p className="text-xs text-gray-400">Disponivel em breve</p>
                </div>
              </div>
              <span className="text-xs text-gray-400 bg-gray-200 px-2 py-1 rounded-full">
                Em breve
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* FTPS/FTP Servers Section */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm mb-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold text-gray-700">Servidor FTPS/FTP</h3>
          <button
            onClick={() => navigate("/servers")}
            className="flex items-center gap-1 text-xs text-eniac-600 hover:text-eniac-700"
          >
            <IconPlus className="w-4 h-4" />
            Adicionar servidor
          </button>
        </div>

        {serversCount === 0 ? (
          <div className="border border-dashed border-gray-300 rounded-lg p-4 text-center">
            <p className="text-sm text-gray-500 mb-3">Nenhum servidor cadastrado</p>
            <button
              onClick={() => navigate("/servers")}
              className="px-4 py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700 transition-colors"
            >
              Cadastrar servidor
            </button>
          </div>
        ) : (
          <div className="space-y-2">
            {servers.slice(0, 3).map((s: any) => (
              <div key={s.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-100">
                <div className="flex items-center gap-3">
                  <IconServer className="w-4 h-4 text-gray-400" />
                  <div>
                    <p className="text-sm font-medium text-gray-700">{s.name}</p>
                    <p className="text-xs text-gray-400">{s.host}:{s.port}</p>
                  </div>
                </div>
                {s.is_default && <Badge variant="success" className="text-xs">Principal</Badge>}
              </div>
            ))}
            {serversCount > 3 && (
              <button
                onClick={() => navigate("/servers")}
                className="text-xs text-gray-500 hover:text-gray-700"
              >
                + {serversCount - 3} mais
              </button>
            )}
          </div>
        )}

        {serversCount > 0 && (
          <div className="mt-4 pt-4 border-t border-gray-100">
            <button
              onClick={() => navigate("/servers")}
              className="flex items-center gap-2 text-sm text-gray-600 hover:text-gray-800"
            >
              <IconExternalLink className="w-4 h-4" />
              Gerenciar servidores
            </button>
          </div>
        )}
      </div>

      {/* Status Summary */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm mb-6">
        <h3 className="font-semibold text-gray-700 mb-4">Resumo do Ambiente</h3>
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">Workspace</span>
            <span className="text-gray-700 flex items-center gap-1">
              <IconCheck className="w-4 h-4 text-green-500" /> Selecionado
            </span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">GitHub</span>
            <span className={githubConnected ? "text-green-600 flex items-center gap-1" : "text-amber-600 flex items-center gap-1"}>
              {githubConnected ? <><IconCheck className="w-4 h-4" /> Conectado</> : <><IconWarning className="w-4 h-4" /> Pendente</>}
            </span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">GitLab</span>
            <span className="text-gray-400 flex items-center gap-1">
              <IconWarning className="w-4 h-4" /> Em breve
            </span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">Servidores</span>
            <span className={serversCount > 0 ? "text-green-600 flex items-center gap-1" : "text-amber-600 flex items-center gap-1"}>
              {serversCount > 0 ? <><IconCheck className="w-4 h-4" /> {serversCount} cadastrado{serversCount > 1 ? "s" : ""}</> : <><IconWarning className="w-4 h-4" /> Nenhum</>}
            </span>
          </div>
        </div>

        {githubConnected && serversCount > 0 && (
          <div className="mt-4 pt-4 border-t border-gray-100 text-center">
            <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-3">
              <p className="text-green-700 font-medium">Ambiente Pronto!</p>
              <p className="text-green-600 text-sm mt-1">Tudo configurado</p>
            </div>
            <button
              onClick={() => navigate("/projects")}
              className="px-6 py-3 bg-eniac-600 text-white rounded-lg font-medium hover:bg-eniac-700 transition-colors"
            >
              Ir para Projetos
            </button>
          </div>
        )}
      </div>

      {/* Switch Workspace Modal */}
      {showSwitchModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setShowSwitchModal(false)} />
          <div className="relative bg-white rounded-xl shadow-2xl p-6 w-full max-w-sm mx-4">
            <h3 className="font-semibold text-gray-800 mb-2">Trocar workspace?</h3>
            <p className="text-sm text-gray-600 mb-4">
              Voce sera levado para a selecao de workspaces. O ambiente atual continuara salvo.
            </p>
            <div className="flex gap-3">
              <button onClick={() => setShowSwitchModal(false)} className="flex-1 py-2 border border-gray-300 text-sm rounded-lg hover:bg-gray-50">Cancelar</button>
              <button onClick={() => { setShowSwitchModal(false); navigate("/"); }} className="flex-1 py-2 bg-eniac-600 text-white text-sm rounded-lg hover:bg-eniac-700">Trocar workspace</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}