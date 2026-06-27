/**
 * AuthProvider — centralises all authentication state.
 *
 * Responsibilities (SRP):
 * - Track auth validity (valid / expired / loading)
 * - Expose rehydrate() so 401 handlers can attempt token recovery
 *   before nuking the session
 * - Expose wipe() for explicit sign-out
 *
 * Does NOT manage sessions or workspace — that's SessionContext's job.
 */

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import {
  initApiClientBase,
  resolveApiToken,
  rehydrateToken,
  clearStoredToken,
} from "./tokenStorage";

export type AuthState = "valid" | "expired" | "loading";

interface AuthContextValue {
  state: AuthState;
  rehydrate: () => Promise<string | null>;
  wipe: () => void;
}

const AuthContext = createContext<AuthContextValue>({
  state: "loading",
  rehydrate: async () => null,
  wipe: () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>("loading");

  const init = useCallback(async () => {
    setState("loading");
    await initApiClientBase();
    const token = await resolveApiToken();
    setState(token ? "valid" : "expired");
    return token;
  }, []);

  const rehydrate = useCallback(async (): Promise<string | null> => {
    const token = await rehydrateToken();
    setState(token ? "valid" : "expired");
    return token;
  }, []);

  const wipe = useCallback(() => {
    clearStoredToken();
    setState("expired");
  }, []);

  useEffect(() => {
    init();
  }, [init]);

  return (
    <AuthContext.Provider value={{ state, rehydrate, wipe }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  return useContext(AuthContext);
}