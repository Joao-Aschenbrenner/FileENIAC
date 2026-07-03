export const IGNORED_PROJECT_DIRS = new Set([
  ".git",
  ".github",
  ".eniac",
  ".vscode",
  ".idea",
  "node_modules",
  "dist",
  "build",
  "target",
  "vendor",
  "__pycache__",
]);

export function isProjectDirectory(name: string): boolean {
  return IGNORED_PROJECT_DIRS.has(name);
}
