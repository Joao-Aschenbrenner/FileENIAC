import { describe, it, expect } from 'vitest';
import { isProjectDirectory, IGNORED_PROJECT_DIRS } from '../projectUtils';

describe('isProjectDirectory', () => {
  it('returns true for .git', () => {
    expect(isProjectDirectory('.git')).toBe(true);
  });

  it('returns true for .github', () => {
    expect(isProjectDirectory('.github')).toBe(true);
  });

  it('returns true for node_modules', () => {
    expect(isProjectDirectory('node_modules')).toBe(true);
  });

  it('returns false for a real project name', () => {
    expect(isProjectDirectory('my-project')).toBe(false);
  });

  it('returns false for FileENIAC', () => {
    expect(isProjectDirectory('FileENIAC')).toBe(false);
  });

  it('rejects all entries in IGNORED_PROJECT_DIRS', () => {
    for (const dir of IGNORED_PROJECT_DIRS) {
      expect(isProjectDirectory(dir)).toBe(true);
    }
  });
});
