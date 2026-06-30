// SPDX-License-Identifier: MIT
/**
 * SessionContext — manages active session, workspace path, and session list.
 *
 * Responsibilities (SRP):
 * - Track which session is active and its workspace path
 * - Persist session selection to localStorage
 * - Handle 401 via AuthProvider.rehydrate() before wiping
 *
 * Auth token state is owned by AuthProvider (inside SessionProvider).
 * This provider purely manages session selection.
 */

import { createContext, useContext, useEffect, useState, useCallback, useRef, type ReactNode } from "react";
import { listSessions, activateSession, checkHealth, heartbeat, clearSessionWorkspace } from "../api/client";
import { ApiError } from "../api/errors";
import type { Session } from "../types";
import { STORAGE_KEYS, storageGet, storageSet, storageClearAuth } from "../api/storage";
import { useDocumentVisible } from "../hooks/useDocumentVisible";
import { AuthProvider, useAuth } from "../auth/AuthProvider";

interface SessionContextValue {
  sessions: Session[];
  activeSession: Session | null;
  workspacePath: string;
  loading: boolean;
  error: string;
  authExpired: boolean;
  authState: "valid" | "expired" | "loading";
  backendOnline: boolean;
  refresh: () => Promise<void>;
  switchSession: (id: number) => Promise<void>;
  removeWorkspace: (id: number) => Promise<void>;
  clearAuthExpired: () => void;
}

const SessionContext = createContext<SessionContextValue>({
  sessions: [],
  activeSession: null,
  workspacePath: "",
  loading: true,
  error: "",
  authExpired: false,
  authState: "loading",
  backendOnline: false,
  refresh: async () => {},
  switchSession: async () => {},
  removeWorkspace: async () => {},
  clearAuthExpired: () => {},
});

function parseError(err: unknown): string {
  if (err instanceof TypeError && err.message === "Failed to fetch") {
    return "Não foi possível conectar ao backend. Verifique se o servidor está rodando.";
  }
  if (err instanceof Error) return err.message;
  return "Erro desconhecido";
}

export function SessionProvider({ children }: { children: ReactNode }) {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [activeSession, setActiveSession] = useState<Session | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [workspacePath, setWorkspacePath] = useState("");
  const [backendOnline, setBackendOnline] = useState(false);

  const auth = useAuth();
  const refreshGenerationRef = useRef(0);
  const refreshInFlightRef = useRef(false);
  const isVisible = useDocumentVisible();

  const refresh = useCallback(async () => {
    if (refreshInFlightRef.current) return;
    refreshInFlightRef.current = true;

    const generation = ++refreshGenerationRef.current;
    setError("");
    try {
      const ok = await checkHealth();
      if (generation !== refreshGenerationRef.current) return;
      setBackendOnline(ok);
      if (!ok) {
        setError("Backend offline. Tente reiniciar o aplicativo.");
        setSessions([]);
        setActiveSession(null);
        setWorkspacePath("");
        return;
      }
      heartbeat().catch(() => {});
      const all = await listSessions();
      if (generation !== refreshGenerationRef.current) return;
      const list = Array.isArray(all) ? all : [];
      setSessions(list);
      const active = list.find((s) => s.is_active);
      if (active) {
        setActiveSession(active);
        setWorkspacePath(active.workspace_path || "");
        storageSet(STORAGE_KEYS.sessionId, String(active.id));
        storageSet(STORAGE_KEYS.workspacePath, active.workspace_path || "");
        return;
      }
      const storedId = storageGet(STORAGE_KEYS.sessionId);
      if (storedId) {
        const found = list.find((s) => String(s.id) === storedId);
        if (found) {
          setActiveSession(found);
          setWorkspacePath(found.workspace_path || "");
          storageSet(STORAGE_KEYS.workspacePath, found.workspace_path || "");
        }
      }
    } catch (e: unknown) {
      if (generation !== refreshGenerationRef.current) return;
      if (e instanceof ApiError && e.isUnauthorized()) {
        const recovered = await auth.rehydrate();
        if (recovered) {
          return;
        }
        storageClearAuth();
        setActiveSession(null);
        setSessions([]);
        setWorkspacePath("");
        setError("Sessão expirada. Selecione uma sessão novamente.");
      } else {
        setError(parseError(e));
      }
    } finally {
      refreshInFlightRef.current = false;
      if (generation === refreshGenerationRef.current) {
        setLoading(false);
      }
    }
  }, [auth]);

  const switchSession = useCallback(async (id: number) => {
    setError("");
    refreshGenerationRef.current++;
    try {
      const sess = await activateSession(id);
      refreshGenerationRef.current++;
      setActiveSession(sess);
      setWorkspacePath(sess.workspace_path || "");
      storageSet(STORAGE_KEYS.sessionId, String(sess.id));
      storageSet(STORAGE_KEYS.workspacePath, sess.workspace_path || "");
    } catch (e: unknown) {
      if (e instanceof ApiError && e.isUnauthorized()) {
        const recovered = await auth.rehydrate();
        if (recovered) {
          return;
        }
        storageClearAuth();
        setActiveSession(null);
        setWorkspacePath("");
        setError("Sessão expirada. Selecione uma sessão novamente.");
        return;
      }
      setError(parseError(e));
      throw e;
    }
  }, [auth]);

  const removeWorkspace = useCallback(async (id: number) => {
    setError("");
    try {
      await clearSessionWorkspace(id);
      await refresh();
    } catch (e: unknown) {
      setError(parseError(e));
    }
  }, [refresh]);

  const clearAuthExpired = useCallback(() => {
    auth.rehydrate();
    setError("");
  }, [auth]);

  useEffect(() => {
    setLoading(true);
    refresh();
    const interval = setInterval(() => {
      if (isVisible) {
        refresh();
      }
    }, 30000);
    return () => {
      clearInterval(interval);
    };
  }, [refresh, isVisible]);

  return (
    <SessionContext.Provider
      value={{
        sessions,
        activeSession,
        workspacePath,
        loading,
        error,
        authExpired: auth.state === "expired",
        authState: auth.state,
        backendOnline,
        refresh,
        switchSession,
        removeWorkspace,
        clearAuthExpired,
      }}
    >
      {children}
    </SessionContext.Provider>
  );
}

export function useSession() {
  return useContext(SessionContext);
}

export { AuthProvider };