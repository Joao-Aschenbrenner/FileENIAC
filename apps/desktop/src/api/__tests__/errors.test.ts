import { describe, it, expect, vi, beforeEach } from "vitest";
import { listSessions } from "../client";
import { ApiError } from "../errors";

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

beforeEach(() => {
  mockFetch.mockReset();
  localStorage.clear();
  localStorage.setItem("eniac_api_token", "unit-test-token");
});

describe("ApiError", () => {
  it("is thrown with status=401 when the backend returns 401 and isUnauthorized() is true", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: async () => ({ error: "token expired" }),
    });
    let caught: unknown;
    try {
      await listSessions();
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    const err = caught as ApiError;
    expect(err.status).toBe(401);
    expect(err.isUnauthorized()).toBe(true);
    expect(err.isForbidden()).toBe(false);
    expect(err.isTimeout()).toBe(false);
    expect(err.message).toBe("token expired");
  });

  it("is thrown with status=500 when the backend returns 500 and isUnauthorized() is false", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: async () => ({ error: "Internal error" }),
    });
    let caught: unknown;
    try {
      await listSessions();
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    const err = caught as ApiError;
    expect(err.status).toBe(500);
    expect(err.isUnauthorized()).toBe(false);
    expect(err.message).toBe("Internal error");
  });

  it("falls back to 'HTTP <status>' message when body has no error field", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 503,
      json: async () => ({}),
    });
    let caught: unknown;
    try {
      await listSessions();
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(ApiError);
    expect((caught as ApiError).message).toBe("HTTP 503");
  });
});
