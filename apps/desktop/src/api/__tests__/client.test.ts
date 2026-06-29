// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { checkHealth, listProjects } from '../client';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

beforeEach(() => {
  mockFetch.mockReset();
});

describe('checkHealth', () => {
  it('calls GET /health and returns true when status is ok', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ status: 'ok' }),
    });
    const result = await checkHealth();
    expect(result).toBe(true);
    expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/api/health');
  });

  it('returns false on network error', async () => {
    mockFetch.mockRejectedValueOnce(new Error('Network error'));
    const result = await checkHealth();
    expect(result).toBe(false);
  });
});

describe('listProjects', () => {
  it('calls GET /projects?workspace= with encoded path', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: 'proj1' }],
    });
    const result = await listProjects('/my/workspace');
    expect(result).toEqual([{ name: 'proj1' }]);
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/projects?workspace=%2Fmy%2Fworkspace'
    );
  });

  it('throws when server returns non-200', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: async () => ({ error: 'Internal error' }),
    });
    await expect(listProjects('/ws')).rejects.toThrow('Internal error');
  });
});
