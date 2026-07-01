// SPDX-License-Identifier: MIT
import { invoke } from "@tauri-apps/api/core";
import { STORAGE_KEYS } from "../api/storage";

const TOKEN_STORAGE_KEY = STORAGE_KEYS.apiToken;
const PORT_STORAGE_KEY = STORAGE_KEYS.apiPort;

let _baseUrl: string | null = null;
let _cachedToken: string | null = null;
let _resolvePromise: Promise<string | null> | null = null;

export interface BackendInfo {
  base_url: string;
  token: string;
  ready: boolean;
}

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

export function configureBackendAuth(info: BackendInfo): boolean {
  if (!info.ready || !info.base_url || !info.base_url.trim()) return false;

  const baseUrl = info.base_url.trim();
  _baseUrl = baseUrl;

  if (typeof window !== "undefined") {
    const port = baseUrl.match(/:(\d+)\/api$/)?.[1];
    if (port) localStorage.setItem(PORT_STORAGE_KEY, port);
  }

  const token = info.token?.trim();
  if (token) {
    storeToken(token);
    _resolvePromise = Promise.resolve(token);
  } else {
    _resolvePromise = null;
  }

  return true;
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
  try {
    const info = await invoke<BackendInfo>("get_backend_info");
    if (configureBackendAuth(info) && info.token?.trim()) {
      return info.token.trim();
    }
  } catch {
    // ignore; fall through
  }

  let token = getStoredToken();
  if (token) return token;

  token = await fetchTokenFromBackend();
  if (token) storeToken(token);
  return token;
}

export async function rehydrateToken(): Promise<string | null> {
  _cachedToken = null;
  _resolvePromise = null;
  return resolveApiToken();
}

export async function initApiClientBase(): Promise<void> {
  if (typeof window === "undefined") {
    _baseUrl = "http://localhost:8080/api";
    return;
  }

  try {
    const info = await invoke<BackendInfo>("get_backend_info");
    if (configureBackendAuth(info)) {
      return;
    }
  } catch {
    // ignore
  }

  try {
    const storedPort = localStorage.getItem(PORT_STORAGE_KEY);
    if (storedPort) {
      _baseUrl = `http://localhost:${storedPort}/api`;
      return;
    }
  } catch {
    // ignore
  }

  _baseUrl = "http://localhost:8080/api";
}

export function getCurrentToken(): string | null {
  if (_cachedToken) return _cachedToken;
  if (typeof window !== "undefined") {
    return localStorage.getItem(TOKEN_STORAGE_KEY);
  }
  return null;
}

export function clearTokenStorageState(): void {
  _cachedToken = null;
  _resolvePromise = null;
  _baseUrl = null;
}
