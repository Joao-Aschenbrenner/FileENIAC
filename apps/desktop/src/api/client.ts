import { invoke } from "@tauri-apps/api/core";

let BASE_URL = "http://localhost:8080/api";

export async function initApiClient(): Promise<void> {
  try {
    const port = await invoke<string>("get_api_port");
    BASE_URL = `http://localhost:${port}/api`;
  } catch {
    BASE_URL = "http://localhost:8080/api";
  }
}

async function get(path: string): Promise<any> {
  const res = await fetch(`${BASE_URL}${path}`);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json();
}

async function post(path: string, body: any): Promise<any> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const errBody = await res.json().catch(() => ({}));
    throw new Error(errBody.error || `HTTP ${res.status}`);
  }
  return res.json();
}

async function del(path: string): Promise<any> {
  const res = await fetch(`${BASE_URL}${path}`, { method: "DELETE" });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json();
}

function ws(path: string, extraParams?: string): string {
  const wsPath = localStorage.getItem("eniac_ws_path") || "";
  const base = `${path}?workspace=${encodeURIComponent(wsPath)}`;
  return extraParams ? `${base}&${extraParams}` : base;
}

export async function checkHealth(): Promise<boolean> {
  try {
    const data = await get("/health");
    return data.status === "ok";
  } catch {
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

export async function executeSync(project: string, action: string): Promise<any> {
  return post(ws("/sync"), { project, action });
}

export async function createMirror(project: string): Promise<any> {
  return post(ws("/mirror"), { project });
}

export async function getHealthCheck(): Promise<any> {
  return get(ws("/health/check"));
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
