// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { checkHealth, createWorkspace, enterWorkspace, listProjects, listWorkspaces } from '../client';
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

describe('workspace flow', () => {
  it('lists workspaces under the selected base folder without storing active workspace path', async () => {
    localStorage.setItem('eniac_api_token', 'unit-test-token');
    setupFetch(async (url: string) => {
      if (url === 'http://localhost:8080/api/workspaces?root=C%3A%2FUsers%2FUSUARIO%2FDesktop%2FPROJETOS%2FENIAC_SYSTEMS') {
        return mockResponse([{ name: 'Cliente A', path: 'C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A' }]);
      }
      return undefined;
    });

    const result = await listWorkspaces('C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS');

    expect(result).toHaveLength(1);
    expect(localStorage.getItem('eniac_ws_path')).toBeNull();
  });

  it('creates a workspace at the requested workspace path without entering it', async () => {
    localStorage.setItem('eniac_api_token', 'unit-test-token');
    setupFetch(async (url: string) => {
      if (url === 'http://localhost:8080/api/workspace') {
        return mockResponse({ name: 'Cliente A', path: 'C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A', projects: 0 });
      }
      return undefined;
    });

    const result = await createWorkspace('C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A', { name: 'Cliente A' });

    expect(result.name).toBe('Cliente A');
    expect(localStorage.getItem('eniac_ws_path')).toBeNull();
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/workspace',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ path: 'C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A', name: 'Cliente A' }),
      }),
    );
  });

  it('stores active workspace path only when entering a workspace', async () => {
    localStorage.setItem('eniac_api_token', 'unit-test-token');
    setupFetch(async (url: string) => {
      if (url === 'http://localhost:8080/api/workspace?workspace=C%3A%2FUsers%2FUSUARIO%2FDesktop%2FPROJETOS%2FENIAC_SYSTEMS%2FCliente-A') {
        return mockResponse({ name: 'Cliente A', path: 'C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A', projects: 0 });
      }
      return undefined;
    });

    const result = await enterWorkspace('C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A');

    expect(result.name).toBe('Cliente A');
    expect(localStorage.getItem('eniac_ws_path')).toBe('C:/Users/USUARIO/Desktop/PROJETOS/ENIAC_SYSTEMS/Cliente-A');
  });
});
