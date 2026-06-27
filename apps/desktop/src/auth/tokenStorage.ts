/**
 * Isolated token storage and resolution.
 *
 * The auth token flows through 3 possible sources (in priority order):
 * 1. Tauri IPC command "get_api_token"  (set by the native shell on startup)
 * 2. Backend handshake endpoint         (unauthenticated, generates ephemeral token)
 * 3. localStorage (persisted bearer)    (from a previous successful handshake)
 *
 * All reads/writes go through this module so the rest of the app
 * never touches localStorage directly for auth concerns.
 */

import { invoke } from "@tauri-apps/api/core";
import { STORAGE_KEYS } from "../api/storage";

const TOKEN_STORAGE_KEY = STORAGE_KEYS.apiToken;
const PORT_STORAGE_KEY = STORAGE_KEYS.apiPort;

let _baseUrl: string | null = null;
let _cachedToken: string | null = null;
let _resolvePromise: Promise<string | null> | null = null;

export function getBaseUrl(): string {
  if (_baseUrl) return _baseUrl;
  return "http://localhost:8080/api";
}

async function fetchTokenFromBackend(): Promise<string | null> {
  try {
    const res = await fetch(`${getBaseUrl()}/_handshake/token`, {
      method: "GET",
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return null;
    const body = await res.json();
    return typeof body?.token === "string" && body.token.length > 0
      ? body.token
      : null;
  } catch {
    return null;
  }
}

function getStoredToken(): string | null {
  if (typeof window === "undefined") return null;
  const stored = localStorage.getItem(TOKEN_STORAGE_KEY);
  if (stored) _cachedToken = stored;
  return stored;
}

function storeToken(token: string): void {
  _cachedToken = token;
  if (typeof window !== "undefined") {
    localStorage.setItem(TOKEN_STORAGE_KEY, token);
  }
}

export function clearStoredToken(): void {
  _cachedToken = null;
  if (typeof window !== "undefined") {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
  }
}

export function isTokenCached(): boolean {
  return _cachedToken !== null;
}

export async function resolveApiToken(): Promise<string | null> {
  if (_resolvePromise) return _resolvePromise;

  _resolvePromise = _doResolve();
  return _resolvePromise;
}

async function _doResolve(): Promise<string | null> {
  let token = getStoredToken();
  if (token) return token;

  // Try Tauri IPC first.
  try {
    const fromTauri: string = await invoke("get_api_token");
    if (fromTauri && fromTauri.trim()) {
      storeToken(fromTauri.trim());
      return fromTauri.trim();
    }
  } catch {
    // ignore; fall through to handshake
  }

  // Last resort: ask the backend directly (unauthenticated handshake).
  token = await fetchTokenFromBackend();
  if (token) storeToken(token);
  return token;
}

/** Rehydrates the token after a wipe — tries Tauri then handshake. */
export async function rehydrateToken(): Promise<string | null> {
  _cachedToken = null;
  _resolvePromise = null; // clear so resolveApiToken retries fresh
  return resolveApiToken();
}

/** Initializes the base URL from Tauri or localStorage. Call once at boot. */
export async function initApiClientBase(): Promise<void> {
  if (typeof window === "undefined") {
    _baseUrl = "http://localhost:8080/api";
    return;
  }

  try {
    const storedPort = localStorage.getItem(PORT_STORAGE_KEY);
    if (storedPort) {
      _baseUrl = `http://localhost:${storedPort}/api`;
      return;
    }

    const port: string = await invoke("get_api_port");
    if (port && String(port).trim()) {
      _baseUrl = `http://localhost:${port}/api`;
      localStorage.setItem(PORT_STORAGE_KEY, String(port));
      return;
    }
  } catch {
    // ignore
  }

  _baseUrl = "http://localhost:8080/api";
}

/** Returns the cached token or null. Does NOT attempt re-resolution. */
export function getCurrentToken(): string | null {
  if (_cachedToken) return _cachedToken;
  if (typeof window !== "undefined") {
    return localStorage.getItem(TOKEN_STORAGE_KEY);
  }
  return null;
}

/** Clears all module-level state (for use in test beforeEach). */
export function clearTokenStorageState(): void {
  _cachedToken = null;
  _resolvePromise = null;
  _baseUrl = null;
}