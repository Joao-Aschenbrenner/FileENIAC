// SPDX-License-Identifier: MIT
/**
 * Centralized storage keys and typed helpers.
 * All localStorage keys are defined here as single source of truth.
 * This prevents string literal typos and makes refactoring safe.
 */

export const STORAGE_KEYS = {
  apiToken: "eniac_api_token",
  apiPort: "eniac_api_port",
  sessionId: "eniac_session_id",
  workspacesRoot: "eniac_workspaces_root",
  workspacePath: "eniac_ws_path",
  githubUser: "github_user",
  lgpdConsent: "eniac_lgpd_consent",
  themeMode: "eniac_theme_mode",
} as const;

export type StorageKey = (typeof STORAGE_KEYS)[keyof typeof STORAGE_KEYS];

/** Returns the value from localStorage or null if not found / SSR / storage error. */
export function storageGet(key: StorageKey): string | null {
  if (typeof window === "undefined") return null;
  try {
    return localStorage.getItem(key);
  } catch {
    return null;
  }
}

/** Sets a value in localStorage (no-op in SSR or on storage error). */
export function storageSet(key: StorageKey, value: string): void {
  if (typeof window === "undefined") return;
  try {
    localStorage.setItem(key, value);
  } catch {
    // Silently ignore storage quota exceeded or unavailable storage.
  }
}

/** Removes a key from localStorage (no-op in SSR or on storage error). */
export function storageRemove(key: StorageKey): void {
  if (typeof window === "undefined") return;
  try {
    localStorage.removeItem(key);
  } catch {
    // Silently ignore.
  }
}

/** Clears ALL FileENIAC keys from localStorage. */
export function storageClearAll(): void {
  if (typeof window === "undefined") return;
  (Object.values(STORAGE_KEYS) as string[]).forEach((k) => {
    localStorage.removeItem(k);
  });
}

/** Clears only auth-related keys (token, session, workspace, github user). */
export function storageClearAuth(): void {
  storageRemove(STORAGE_KEYS.apiToken);
  storageRemove(STORAGE_KEYS.sessionId);
  storageRemove(STORAGE_KEYS.workspacePath);
  storageRemove(STORAGE_KEYS.githubUser);
}
