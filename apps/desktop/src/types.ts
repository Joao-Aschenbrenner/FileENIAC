// SPDX-License-Identifier: MIT
export interface Session {
  id: number;
  name: string;
  description: string;
  workspace_path: string;
  is_active: boolean;
  github_user?: string;
}

export interface Workspace {
  name: string;
  description: string;
  path: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: number;
  name: string;
  local_path: string;
  git_url: string;
  branch: string;
  divergence_status: string;
  server_count: number;
  import_status?: string;
  organization?: string;
  repo_name?: string;
  clone_path?: string;
  created_at: string;
  updated_at: string;
}

export interface Server {
  id: number;
  project_id: number;
  project_name?: string;
  name: string;
  host: string;
  port: number;
  user: string;
  password?: string;
  target_path: string;
  is_active: boolean;
  created_at: string;
}

export interface DeployLog {
  id: number;
  project_id: number;
  version: string;
  status: string;
  artifact_hash: string;
  started_at: string;
  completed_at: string;
}

export interface SyncLog {
  id: number;
  project_id: number;
  project_name?: string;
  action: string;
  result: string;
  created_at: string;
}

export interface HistoryEvent {
  id: number;
  type: string;
  description: string;
  created_at: string;
}

export interface GitHubRepo {
  id: number;
  name: string;
  full_name: string;
  clone_url: string;
  description: string;
  default_branch: string;
  imported?: boolean;
}

export interface GitHubOrg {
  login: string;
  description: string;
}

export interface Repository {
  id: number;
  project_id: number;
  name: string;
  full_name: string;
  clone_url: string;
  default_branch: string;
  import_status: string;
  clone_path: string;
}

export interface Settings {
  key: string;
  value: string;
}

export interface DiffFile {
  path: string;
  status: string;
  local_hash: string;
  mirror_hash: string;
}

export interface HealthCheck {
  status: string;
  projects_total: number;
  servers_total: number;
  divergent_total: number;
  last_events: HistoryEvent[];
}

export interface RefreshResult {
  organizations: number;
  repositories: number;
  changes_found: number;
}

export interface ReadinessCheck {
  ready: boolean;
  checks: { name: string; passed: boolean; error?: string }[];
}

export interface RepairReport {
  orphaned: number;
  broken_paths: number;
  fixed: number;
  warnings: string[];
}

export interface BackgroundHealth {
  timestamp: string;
  status: string;
  github_token_valid: boolean;
  projects_total: number;
  servers_total: number;
  divergent_total: number;
}
