// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  listProjects,
  createProject,
  checkHealth,
  TimeoutError,
  GET_TIMEOUT_MS,
  MUTATION_TIMEOUT_MS,
} from "../client";
import { clearTokenStorageState } from "../../auth/tokenStorage";
import { ApiError } from "../errors";

// @tauri-apps/api/core is unavailable under jsdom — the Tauri host only
// exists inside the desktop shell. Mock invoke() so resolveApiToken()
// short-circuits cleanly without falling through to the handshake fetch.
vi.mock("@tauri-apps/api/core", () => ({
  invoke: vi.fn().mockResolvedValue(undefined),
}));

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

beforeEach(() => {
  mockFetch.mockReset();
  localStorage.clear();
  clearTokenStorageState(); // reset module-level caches that localStorage.clear() doesn't touch
});

// Helper: extract headers from the *first* fetch invocation that hits the
// path. Honors the optional handshake behaviour — when there is no token
// in localStorage the client auto-issues an unauthenticated GET to
// /_handshake/token before the main request, so the call index of the
// actual call is 1 in that case. The spec doesn't care about token
// resolution per test, only about the data-fetch headers, so we accept
// either index.
function pickDataCallIdx(urlMatch: RegExp): number {
  for (let i = 0; i < mockFetch.mock.calls.length; i++) {
    if (urlMatch.test(String(mockFetch.mock.calls[i][0]))) return i;
  }
  return -1;
}

function lastHeadersContaining(urlMatch: RegExp): Record<string, string> {
  const idx = pickDataCallIdx(urlMatch);
  expect(idx).toBeGreaterThanOrEqual(0);
  const init = mockFetch.mock.calls[idx][1] as RequestInit;
  return (init?.headers ?? {}) as Record<string, string>;
}

describe("API client security headers", () => {
  it("all requests include User-Agent header FileENIAC/1.0.0", async () => {
    // Preset a token so resolveApiToken() does NOT issue an extra
    // handshake fetch — keeps the assertion on a single call.
    localStorage.setItem("eniac_api_token", "unit-test-token");
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: "p" }],
    });
    await listProjects("/ws");
    const headers = lastHeadersContaining(/\/projects/);
    expect(headers["User-Agent"]).toBe("FileENIAC/1.0.0");
  });

  it("all requests include X-Workspace header sourced from localStorage", async () => {
    localStorage.setItem("eniac_api_token", "unit-test-token");
    localStorage.setItem("eniac_ws_path", "/my/custom-workspace");
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: "p" }],
    });
    await listProjects("/my/custom-workspace");
    const headers = lastHeadersContaining(/\/projects/);
    expect(headers["X-Workspace"]).toBe("/my/custom-workspace");
  });

  it("X-Workspace is empty string when localStorage has no eniac_ws_path", async () => {
    localStorage.setItem("eniac_api_token", "unit-test-token");
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: "p" }],
    });
    await listProjects("/ws");
    const headers = lastHeadersContaining(/\/projects/);
    expect(headers["X-Workspace"]).toBe("");
  });

  it("Authorization Bearer header included when eniac_api_token is set", async () => {
    localStorage.setItem("eniac_api_token", "secret-token-xyz");
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: "p" }],
    });
    await listProjects("/ws");
    const headers = lastHeadersContaining(/\/projects/);
    expect(headers["Authorization"]).toBe("Bearer secret-token-xyz");
  });

  it("Authorization header absent when no token is stored (handshake fails too)", async () => {
    // No token in localStorage, invoke() returns undefined (Tauri stub),
    // so resolveApiToken() falls through to fetch(_handshake/token). We
    // mock that to return a 404 so the handshake also fails to produce a
    // token. No Authorization header should ever leak.
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404,
      json: async () => ({}),
    });
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: "p" }],
    });
    await listProjects("/ws");
    const headers = lastHeadersContaining(/\/projects/);
    expect(headers["Authorization"]).toBeUndefined();
    // sanity: two calls happened (handshake + projects)
    expect(mockFetch).toHaveBeenCalledTimes(2);
  });

  it("Authorization remains absent for /api/health when no token present", async () => {
    // /health is the public fallback path — the spec wants to ensure it
    // never accidentally carries auth when there's no token.
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404,
      json: async () => ({}),
    });
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ status: "ok" }),
    });
    const ok = await checkHealth();
    expect(ok).toBe(true);
    const headers = lastHeadersContaining(/\/health$/);
    expect(headers["Authorization"]).toBeUndefined();
  });

  it("CORS-relevant headers are NOT exposed (no proxy-style abuse surface)", async () => {
    localStorage.setItem("eniac_api_token", "unit-test-token");
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ name: "p" }],
    });
    await listProjects("/ws");
    const idx = pickDataCallIdx(/\/projects/);
    const init = mockFetch.mock.calls[idx][1] as RequestInit;
    const headers = (init.headers ?? {}) as Record<string, string>;
    // The browser controls Origin/Referer/Host for cross-origin requests;
    // the client itself must never set them manually or pre-emptively.
    expect(headers["Origin"]).toBeUndefined();
    expect(headers["Referer"]).toBeUndefined();
    expect(headers["Host"]).toBeUndefined();
    // credentials mode must not be "include" (would send cookies cross-origin)
    const creds = (init.credentials ?? "same-origin") as string;
    expect(creds).not.toBe("include");
  });

  it("AbortController fires on timeout — fetch receives a signal that aborts", async () => {
    vi.useFakeTimers();
    try {
      localStorage.setItem("eniac_api_token", "unit-test-token");
      let abortFired = false;
      mockFetch.mockImplementationOnce(
        (_url: unknown, init: RequestInit | undefined) =>
          new Promise((_resolve, reject) => {
            if (init?.signal) {
              init.signal.addEventListener("abort", () => {
                abortFired = true;
                reject(new DOMException("aborted", "AbortError"));
              });
              if (init.signal.aborted) {
                abortFired = true;
                reject(new DOMException("aborted", "AbortError"));
              }
            }
          }),
      );

      const promise = listProjects("/ws");
      const expectation = expect(promise).rejects.toBeInstanceOf(TimeoutError);
      await vi.advanceTimersByTimeAsync(GET_TIMEOUT_MS + 100);
      await expectation;
      expect(abortFired).toBe(true);
    } finally {
      vi.useRealTimers();
    }
  });

  it("TimeoutError is thrown when GET exceeds GET_TIMEOUT_MS (10s)", async () => {
    vi.useFakeTimers();
    try {
      localStorage.setItem("eniac_api_token", "unit-test-token");
      mockFetch.mockImplementationOnce(
        (_url: unknown, init: RequestInit | undefined) =>
          new Promise((_resolve, reject) => {
            if (init?.signal) {
              init.signal.addEventListener("abort", () => {
                reject(new DOMException("aborted", "AbortError"));
              });
            }
          }),
      );

      const promise = listProjects("/ws");
      const expectation = expect(promise).rejects.toBeInstanceOf(TimeoutError);
      await vi.advanceTimersByTimeAsync(GET_TIMEOUT_MS + 500);
      await expectation;
    } finally {
      vi.useRealTimers();
    }
  });

  it("TimeoutError is thrown when POST exceeds MUTATION_TIMEOUT_MS (30s)", async () => {
    vi.useFakeTimers();
    try {
      localStorage.setItem("eniac_api_token", "unit-test-token");
      mockFetch.mockImplementationOnce(
        (_url: unknown, init: RequestInit | undefined) =>
          new Promise((_resolve, reject) => {
            if (init?.signal) {
              init.signal.addEventListener("abort", () => {
                reject(new DOMException("aborted", "AbortError"));
              });
            }
          }),
      );

      const promise = createProject("/ws", { name: "newproject" });
      const expectation = expect(promise).rejects.toBeInstanceOf(TimeoutError);
      await vi.advanceTimersByTimeAsync(MUTATION_TIMEOUT_MS + 500);
      await expectation;
    } finally {
      vi.useRealTimers();
    }
  }, 15000);
});

describe("401 authentication error propagation", () => {
  it("401 on GET throws ApiError with isUnauthorized=true and does not retry", async () => {
    localStorage.setItem("eniac_api_token", "expired-token");
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: async () => ({ error: "token expired" }),
    });
    let caught: unknown;
    try {
      await listProjects("/ws");
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    const err = caught as ApiError;
    expect(err.status).toBe(401);
    expect(err.isUnauthorized()).toBe(true);
    expect(mockFetch).toHaveBeenCalledTimes(1);
  });

  it("401 on POST throws ApiError without infinite retry loop", async () => {
    localStorage.setItem("eniac_api_token", "expired-token");
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: async () => ({ error: "invalid token" }),
    });
    let caught: unknown;
    try {
      await createProject("/ws", { name: "testproject" });
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    const err = caught as ApiError;
    expect(err.status).toBe(401);
    expect(err.isUnauthorized()).toBe(true);
    expect(mockFetch).toHaveBeenCalledTimes(1);
  });

  it("401 response does not expose raw token or credentials in error message", async () => {
    localStorage.setItem("eniac_api_token", "super-secret-token-12345");
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: async () => ({ error: "token expired" }),
    });
    let caught: unknown;
    try {
      await listProjects("/ws");
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    const err = caught as ApiError;
    expect(err.message).toBe("token expired");
    expect(err.message).not.toContain("super-secret");
    expect(err.message).not.toContain("super-secret-token-12345");
  });

  it("multiple rapid 401 calls each throw exactly once (no deduplication)", async () => {
    localStorage.setItem("eniac_api_token", "expired-token");
    mockFetch.mockResolvedValue({
      ok: false,
      status: 401,
      json: async () => ({ error: "token expired" }),
    });
    const results = await Promise.allSettled([
      listProjects("/ws"),
      listProjects("/ws"),
      listProjects("/ws"),
    ]);
    expect(results.filter((r) => r.status === "rejected")).toHaveLength(3);
    expect(mockFetch).toHaveBeenCalledTimes(3);
  });

  it("403 Forbidden is distinct from 401 — isUnauthorized returns false", async () => {
    localStorage.setItem("eniac_api_token", "valid-token");
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 403,
      json: async () => ({ error: "access denied" }),
    });
    let caught: unknown;
    try {
      await listProjects("/ws");
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    const err = caught as ApiError;
    expect(err.status).toBe(403);
    expect(err.isUnauthorized()).toBe(false);
    expect(err.isForbidden()).toBe(true);
  });
});
