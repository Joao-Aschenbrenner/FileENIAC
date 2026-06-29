// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { checkHealth, listProjects } from '../client';
import { clearTokenStorageState } from '../../auth/tokenStorage';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

beforeEach(() => {
  mockFetch.mockReset();
  clearTokenStorageState();
});

function mockResponse(body: any, ok = true, status = 200) {
  return { ok, status, json: async () => body };
}

function setupFetch(handler: (url: string) => Promise<any> | any) {
  mockFetch.mockImplementation(handler);
}

describe('checkHealth', () => {
  it('calls GET /health and returns true when status is ok', async () => {
    setupFetch(async (url: string) => {
      if (url.includes('/_handshake/token')) return mockResponse({ token: 'test-token' });
      if (url === 'http://localhost:8080/api/health') return mockResponse({ status: 'ok' });
      return undefined;
    });
    const result = await checkHealth();
    expect(result).toBe(true);
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/health',
      expect.objectContaining({
        headers: expect.objectContaining({
          Accept: 'application/json',
          'User-Agent': 'FileENIAC/1.0.0',
          'X-Workspace': '',
        }),
        signal: expect.any(AbortSignal),
      }),
    );
  });

  it('returns false on network error', async () => {
    setupFetch(async () => {
      throw new Error('Network error');
    });
    const result = await checkHealth();
    expect(result).toBe(false);
  });
});

describe('listProjects', () => {
  it('calls GET /projects?workspace= with encoded path', async () => {
    setupFetch(async (url: string) => {
      if (url.includes('/_handshake/token')) return mockResponse({ token: 'test-token' });
      if (url === 'http://localhost:8080/api/projects?workspace=%2Fmy%2Fworkspace') {
        return mockResponse([{ name: 'proj1' }]);
      }
      return undefined;
    });
    const result = await listProjects('/my/workspace');
    expect(result).toEqual([{ name: 'proj1' }]);
  });

  it('throws when server returns non-200', async () => {
    setupFetch(async (url: string) => {
      if (url.includes('/_handshake/token')) return mockResponse({ token: 'test-token' });
      if (url === 'http://localhost:8080/api/projects?workspace=%2Fws') {
        return mockResponse({ error: 'Internal error' }, false, 500);
      }
      return undefined;
    });
    await expect(listProjects('/ws')).rejects.toThrow('Internal error');
  });
});
