import { describe, it, expect, vi, beforeEach } from 'vitest';
import { executeSync, executeSyncSafe, executeSyncWithDelete } from '../client';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

function lastPostBody(): any {
  const calls = mockFetch.mock.calls;
  const last = calls[calls.length - 1];
  const init = last[1] as RequestInit;
  return JSON.parse(String(init.body));
}

beforeEach(() => {
  mockFetch.mockReset();
  mockFetch.mockResolvedValue({
    ok: true,
    json: async () => ({ manifest: { result: 'ok' } }),
  });
  localStorage.clear();
  localStorage.setItem('eniac_api_token', 'unit-test-token');
  localStorage.setItem('eniac_ws_path', '/ws');
});

// H-15: the contract for executeSync MUST forward a confirm flag that
// matches what the caller asked for. The legacy implementation
// always sent confirm=false and silently dropped deletions.
describe('executeSync confirm contract (H-15)', () => {
  it('boolean true argument maps to confirm=true in body', async () => {
    await executeSync('proj', 'mirror_update', true);
    expect(lastPostBody()).toEqual({
      project: 'proj',
      action: 'mirror_update',
      confirm: true,
    });
  });

  it('boolean false argument maps to confirm=false in body', async () => {
    await executeSync('proj', 'mirror_update', false);
    expect(lastPostBody()).toEqual({
      project: 'proj',
      action: 'mirror_update',
      confirm: false,
    });
  });

  it('options { confirm: true } maps to confirm=true in body', async () => {
    await executeSync('proj', 'mirror_update', { confirm: true });
    expect(lastPostBody().confirm).toBe(true);
  });

  it('options { confirmDeleting: true } maps to confirm=true in body', async () => {
    await executeSync('proj', 'mirror_update', { confirmDeleting: true });
    expect(lastPostBody().confirm).toBe(true);
  });

  it('options { confirmDeleting: false, no confirm } defaults to confirm=false', async () => {
    await executeSync('proj', 'mirror_update', { confirmDeleting: false });
    expect(lastPostBody().confirm).toBe(false);
  });

  it('no third argument defaults to confirm=false (safe default)', async () => {
    await executeSync('proj', 'mirror_update');
    expect(lastPostBody()).toEqual({
      project: 'proj',
      action: 'mirror_update',
      confirm: false,
    });
  });

  it('explicit confirm beats confirmDeleting', async () => {
    await executeSync('proj', 'mirror_update', {
      confirm: false,
      confirmDeleting: true,
    });
    expect(lastPostBody().confirm).toBe(false);
  });

  it('executeSyncSafe forwards confirm=false in body', async () => {
    await executeSyncSafe('proj', 'mirror_update');
    expect(lastPostBody().confirm).toBe(false);
  });

  it('executeSyncWithDelete forwards confirm=true in body', async () => {
    await executeSyncWithDelete('proj', 'mirror_update');
    expect(lastPostBody().confirm).toBe(true);
  });

  it('POSTs to /sync with the project workspace in the URL', async () => {
    await executeSync('proj', 'mirror_update');
    const url = mockFetch.mock.calls[0][0] as string;
    expect(url).toContain('/api/sync');
    expect(url).toContain('workspace=');
  });
});
