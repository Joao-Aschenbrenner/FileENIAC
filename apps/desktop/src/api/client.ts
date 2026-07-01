// SPDX-License-Identifier: MIT
import { invoke } from "@tauri-apps/api/core";
import { BackendInfo, configureBackendAuth, resolveApiToken } from "../auth/tokenStorage";
import { ApiError } from "./errors";
import { STORAGE_KEYS, storageGet } from "./storage";

let BASE_URL = "http://localhost:8080/api";

export const GET_TIMEOUT_MS = 10_000;
export const MUTATION_TIMEOUT_MS = 30_000;

export class TimeoutError extends Error {
  constructor(public ms: number) {
    super(`Request timed out after ${ms}ms`);
    this.name = "TimeoutError";
  }
}

export async function initApiClient(): Promise<void> {
  try {
    const info = await invoke<BackendInfo>("get_backend_info");
    if (configureApiClientFromBackendInfo(info)) {
      return;
    }
  } catch {
    BASE_URL = "http://localhost:8080/api";
  }
}

export function configureApiClientFromBackendInfo(info: BackendInfo): boolean {
  if (!configureBackendAuth(info)) return false;
  BASE_URL = info.base_url.trim();
  return true;
}

async function get(path: string): Promise<any> {
  const token = await resolveApiToken();
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), GET_TIMEOUT_MS);
  try {
    const res = await fetch(`${BASE_URL}${path}`, {
      signal: controller.signal,
      headers: {
        "Accept": "application/json",
        "User-Agent": "FileENIAC/1.0.0",
        "X-Workspace": storageGet(STORAGE_KEYS.workspacePath) || "",
        ...(token ? { "Authorization": `Bearer ${token}` } : {}),
      },
    });
    clearTimeout(timer);
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new ApiError(res.status, `${BASE_URL}${path}`, body.error || `HTTP ${res.status}`);
    }
    return res.json();
  } catch (err) {
    clearTimeout(timer);
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new TimeoutError(GET_TIMEOUT_MS);
    }
    throw err;
  }
}

async function post(path: string, body: any): Promise<any> {
  const token = await resolveApiToken();
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), MUTATION_TIMEOUT_MS);
  try {
    const res = await fetch(`${BASE_URL}${path}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "User-Agent": "FileENIAC/1.0.0",
        "X-Workspace": storageGet(STORAGE_KEYS.workspacePath) || "",
        ...(token ? { "Authorization": `Bearer ${token}` } : {}),
      },
      body: JSON.stringify(body),
      signal: controller.signal,
    });
    clearTimeout(timer);
    if (!res.ok) {
      const errBody = await res.json().catch(() => ({}));
      throw new ApiError(res.status, `${BASE_URL}${path}`, errBody.error || `HTTP ${res.status}`);
    }
    return res.json();
  } catch (err) {
    clearTimeout(timer);
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new TimeoutError(MUTATION_TIMEOUT_MS);
    }
    throw err;
  }
}

async function del(path: string): Promise<any> {
  const token = await resolveApiToken();
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), MUTATION_TIMEOUT_MS);
  try {
    const res = await fetch(`${BASE_URL}${path}`, {
      method: "DELETE",
      headers: {
        "Accept": "application/json",
        "User-Agent": "FileENIAC/1.0.0",
        "X-Workspace": storageGet(STORAGE_KEYS.workspacePath) || "",
        ...(token ? { "Authorization": `Bearer ${token}` } : {}),
      },
      signal: controller.signal,
    });
    clearTimeout(timer);
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new ApiError(res.status, `${BASE_URL}${path}`, body.error || `HTTP ${res.status}`);
    }
    return res.json();
  } catch (err) {
    clearTimeout(timer);
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new TimeoutError(MUTATION_TIMEOUT_MS);
    }
    throw err;
  }
}

function ws(path: string, extraParams?: string): string {
  const wsPath = localStorage.getItem("eniac_ws_path") || "";
  const base = `${path}?workspace=${encodeURIComponent(wsPath)}`;
  return extraParams ? `${base}&${extraParams}` : base;
}

export async function checkHealth(): Promise<boolean> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), GET_TIMEOUT_MS);
  try {
    const res = await fetch(`${BASE_URL}/health`, {
      signal: controller.signal,
      headers: { Accept: "application/json" },
    });
    clearTimeout(timer);
    if (!res.ok) return false;
    const data = await res.json();
    return data.status === "ok";
  } catch {
    clearTimeout(timer);
    return false;
  }
}

export async function heartbeat(): Promise<void> {
  try {
    await post("/heartbeat", {});
  } catch {
    // ignore
  }
}

export async function getWorkspace(wsPath: string): Promise<any> {
  return get(`/workspace?workspace=${encodeURIComponent(wsPath)}`);
}

export async function listProjects(wsPath: string): Promise<any[]> {
  return get(`/projects?workspace=${encodeURIComponent(wsPath)}`);
}

export async function getProject(wsPath: string, name: string): Promise<any> {
  return get(`/projects/${encodeURIComponent(name)}?workspace=${encodeURIComponent(wsPath)}`);
}

export async function createProject(wsPath: string, project: any): Promise<any> {
  return post(`/projects?workspace=${encodeURIComponent(wsPath)}`, project);
}

export async function deleteProject(wsPath: string, name: string): Promise<any> {
  return del(`/projects/${encodeURIComponent(name)}?workspace=${encodeURIComponent(wsPath)}`);
}

export async function listServers(wsPath: string, project?: string): Promise<any[]> {
  let path = `/servers?workspace=${encodeURIComponent(wsPath)}`;
  if (project) path += `&project=${encodeURIComponent(project)}`;
  return get(path);
}

export async function getServer(wsPath: string, id: number): Promise<any> {
  return get(`/servers/${id}?workspace=${encodeURIComponent(wsPath)}`);
}

export async function createServer(wsPath: string, server: any): Promise<any> {
  return post(`/servers?workspace=${encodeURIComponent(wsPath)}`, server);
}

export async function deleteServer(wsPath: string, id: number): Promise<any> {
  return del(`/servers/${id}?workspace=${encodeURIComponent(wsPath)}`);
}

export async function getSettings(wsPath: string): Promise<Record<string, string>> {
  return get(`/settings?workspace=${encodeURIComponent(wsPath)}`);
}

export async function updateSettings(wsPath: string, settings: Record<string, string>): Promise<any> {
  return post(`/settings?workspace=${encodeURIComponent(wsPath)}`, settings);
}

export async function getHistory(wsPath: string, params: { project?: string; type?: string; limit?: number; offset?: number }): Promise<any[]> {
  const qs = new URLSearchParams();
  if (params.project) qs.set("project", params.project);
  if (params.type) qs.set("type", params.type);
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.offset) qs.set("offset", String(params.offset));
  const extra = qs.toString();
  return get(extra ? `/history?workspace=${encodeURIComponent(wsPath)}&${extra}` : `/history?workspace=${encodeURIComponent(wsPath)}`);
}

export async function getEvents(wsPath: string, params: { type?: string; limit?: number; offset?: number }): Promise<any[]> {
  const qs = new URLSearchParams();
  if (params.type) qs.set("type", params.type);
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.offset) qs.set("offset", String(params.offset));
  const extra = qs.toString();
  return get(extra ? `/events?workspace=${encodeURIComponent(wsPath)}&${extra}` : `/events?workspace=${encodeURIComponent(wsPath)}`);
}

export async function getDeploys(project: string, limit?: number): Promise<any[]> {
  const qs = new URLSearchParams({ project });
  if (limit) qs.set("limit", String(limit));
  return get(ws("/deploys", qs.toString()));
}

export async function executeDeploy(project: string, useFallback?: boolean): Promise<any> {
  return post(ws("/deploy"), { project, use_fallback: useFallback ?? false });
}

export async function executeRollback(project: string): Promise<any> {
  return post(ws("/rollback"), { project });
}

export async function executeVerify(project: string): Promise<any> {
  return post(ws("/verify"), { project });
}

export async function getDiff(project: string): Promise<any> {
  return get(`${ws("/diff")}&project=${encodeURIComponent(project)}`);
}

export async function getSyncs(project?: string, limit?: number): Promise<any[]> {
  const qs = new URLSearchParams();
  if (project) qs.set("project", project);
  if (limit) qs.set("limit", String(limit));
  return get(ws("/syncs", qs.toString()));
}

export async function executeSync(project: string, action: string, confirm?: boolean | { confirm?: boolean; confirmDeleting?: boolean }): Promise<any> {
  let confirmValue = false;
  if (typeof confirm === "boolean") {
    confirmValue = confirm;
  } else if (confirm && typeof confirm === "object") {
    confirmValue = confirm.confirm ?? (confirm.confirmDeleting ?? false);
  }
  return post(ws("/sync"), { project, action, confirm: confirmValue });
}

export async function executeSyncSafe(project: string, action: string): Promise<any> {
  return executeSync(project, action, false);
}

export async function executeSyncWithDelete(project: string, action: string): Promise<any> {
  return executeSync(project, action, true);
}

export async function createMirror(project: string): Promise<any> {
  return post(ws("/mirror"), { project });
}

export async function getHealthCheck(): Promise<any> {
  return get(ws("/health/check"));
}

export async function listSessions(): Promise<any[]> {
  return get(ws("/sessions"));
}

export async function activateSession(id: number): Promise<any> {
  return post(ws(`/sessions/${id}/activate`), {});
}

export async function clearSessionWorkspace(id: number): Promise<any> {
  return post(ws(`/sessions/${id}/clear-workspace`), {});
}

export async function deleteSession(id: number): Promise<any> {
  return del(ws(`/sessions/${id}`));
}

export async function createSession(data: { name: string; description: string }): Promise<any> {
  return post(ws("/sessions"), data);
}

export async function updateSession(id: number, data: Record<string, any>): Promise<any> {
  return post(ws(`/sessions/${id}`), data);
}

// GitHub endpoints
export async function getGitHubStatus(): Promise<any> {
  return get(ws("/github/status"));
}

export async function gitHubLogin(token: string): Promise<any> {
  return post(ws("/github/login"), { token });
}

export async function gitHubLogout(): Promise<any> {
  return post(ws("/github/logout"), {});
}

export async function getGitHubOrganizations(): Promise<any[]> {
  return get(ws("/github/organizations"));
}

export async function getGitHubRepositories(org?: string): Promise<any[]> {
  const extra = org ? `org=${encodeURIComponent(org)}` : "";
  return get(ws("/github/repositories", extra));
}

export async function importGitHubRepos(repos: any[], cloneDir?: string): Promise<any[]> {
  return post(ws("/github/import"), { repos, clone_dir: cloneDir });
}

export async function cloneGitHubRepo(repoId: number, projectId: number, cloneUrl: string, branch: string, cloneDir: string): Promise<any> {
  return post(ws("/github/clone"), { repo_id: repoId, project_id: projectId, clone_url: cloneUrl, branch, clone_dir: cloneDir });
}

export async function listRepositories(org?: string): Promise<any[]> {
  const wsPath = localStorage.getItem("eniac_ws_path") || "";
  let path = `/repositories?workspace=${encodeURIComponent(wsPath)}`;
  if (org) path += `&org=${encodeURIComponent(org)}`;
  return get(path);
}

export async function getRepository(githubId: number): Promise<any> {
  return get(`/repositories/${githubId}?workspace=${encodeURIComponent(localStorage.getItem("eniac_ws_path") || "")}`);
}
