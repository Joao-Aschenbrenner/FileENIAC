// SPDX-License-Identifier: MIT
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useSession } from "../context/SessionContext";
import { Loader, Modal, Button } from "../components/ui";
import { deleteSession } from "../api/client";
import { Trash2, RefreshCw, Plus, X } from "lucide-react";

export default function SessionSelector() {
  const navigate = useNavigate();
  const { sessions, loading, error, refresh, switchSession, removeWorkspace } = useSession();
  const [deleteTarget, setDeleteTarget] = useState<{ id: number; name: string } | null>(null);
  const [deleting, setDeleting] = useState(false);

  async function handleSelect(id: number) {
    try {
      await switchSession(id);
      const selected = sessions.find(s => s.id === id);
      if (selected?.workspace_path && selected.workspace_path.trim() !== "") {
        navigate("/dashboard");
      }
    } catch {
      // error is displayed by context state
    }
  }

  async function handleConfirmDelete() {
    if (!deleteTarget) return;
    setDeleting(true);
    try {
      await deleteSession(deleteTarget.id);
      await refresh();
    } catch (err: any) {
      console.error("Falha ao excluir sessão:", err.message);
    }
    setDeleting(false);
    setDeleteTarget(null);
  }

  async function handleRetry() {
    await refresh();
  }

  if (loading && sessions.length === 0) return <Loader text="Carregando sessões..." />;

  return (
    <div className="flex h-screen items-center justify-center bg-gradient-to-br from-eniac-900 to-eniac-700">
      <div className="w-full max-w-lg px-8">
        <div className="text-center mb-10">
          <h1 className="text-4xl font-bold text-white mb-2">FileENIAC</h1>
          <p className="text-eniac-200 text-sm">Selecione ou crie uma sessão de trabalho</p>
        </div>

        {error && (
          <div className="bg-red-500/20 border border-red-400/40 text-red-200 px-4 py-3 rounded-lg text-sm mb-6 text-center">
            <p>{error}</p>
            <Button
              onClick={handleRetry}
              variant="ghost"
              size="sm"
              className="mt-2 text-red-200 hover:text-white"
            >
              Tentar novamente
            </Button>
          </div>
        )}

        {!error && sessions.length === 0 && (
          <div className="bg-eniac-800/40 border border-eniac-600/40 text-eniac-200 px-4 py-6 rounded-lg text-sm mb-6 text-center">
            <p className="mb-1">Nenhuma sessão encontrada.</p>
            <p className="text-eniac-300">Crie uma nova sessão para começar.</p>
          </div>
        )}

        <div className="space-y-3 mb-6">
          {sessions.map((s) => (
            <div
              key={s.id}
              className="bg-white/10 hover:bg-white/20 border border-white/20 rounded-xl p-5 transition-all group"
            >
              <div className="flex items-center justify-between">
                <div className="flex-1 cursor-pointer" onClick={() => handleSelect(s.id!)}>
                  <h3 className="text-white font-semibold text-lg">{s.name}</h3>
                  {s.description && (
                    <p className="text-eniac-200 text-sm mt-0.5">{s.description}</p>
                  )}
                  <p className="text-eniac-300 text-xs mt-1 flex items-center gap-2">
                    <span className={s.workspace_path ? "" : "italic"}>
                      {s.workspace_path || "Nenhum workspace configurado"}
                    </span>
                    {s.workspace_path && (
                      <button
                        type="button"
                        onClick={(e) => { e.stopPropagation(); removeWorkspace(s.id!); }}
                        className="text-eniac-300 hover:text-red-400 transition-colors"
                        title="Remover workspace"
                        aria-label="Remover workspace"
                      >
                        <X className="h-3.5 w-3.5" />
                      </button>
                    )}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  {s.github_user && (
                    <span className="text-xs text-eniac-300 bg-white/10 px-2 py-1 rounded">
                      {s.github_user}
                    </span>
                  )}
                  <Button
                    onClick={() => handleSelect(s.id!)}
                    variant="primary"
                    size="sm"
                  >
                    Selecionar
                  </Button>
                  <button
                    onClick={(e) => { e.stopPropagation(); setDeleteTarget({ id: s.id!, name: s.name }); }}
                    className="p-2 text-eniac-300 hover:text-red-400 hover:bg-red-500/20 rounded-lg transition-colors"
                    title="Excluir sessão"
                    aria-label={`Excluir sessão ${s.name}`}
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        <Button
          onClick={() => navigate("/wizard")}
          disabled={!!error}
          className="w-full"
          size="lg"
          icon={<Plus className="h-5 w-5" />}
          iconPosition="left"
        >
          Criar Nova Sessão
        </Button>

        <div className="mt-8 text-center">
          <Button
            onClick={handleRetry}
            variant="ghost"
            size="sm"
            className="text-eniac-300 hover:text-white"
            icon={<RefreshCw className="h-3.5 w-3.5" />}
            iconPosition="left"
          >
            Atualizar lista
          </Button>
        </div>
      </div>

      <Modal
        open={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        title="Excluir Sessão"
        size="sm"
      >
        <p className="text-sm text-gray-600 mb-4">
          Tem certeza que deseja excluir a sessão <strong>{deleteTarget?.name}</strong>? Esta ação não pode ser desfeita.
        </p>
        <div className="flex gap-3 justify-end">
          <Button variant="secondary" onClick={() => setDeleteTarget(null)}>
            Cancelar
          </Button>
          <Button variant="danger" onClick={handleConfirmDelete} loading={deleting}>
            {deleting ? "Excluindo..." : "Excluir"}
          </Button>
        </div>
      </Modal>
    </div>
  );
}
