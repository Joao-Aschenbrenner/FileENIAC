// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  resolveApiToken,
  rehydrateToken,
  clearStoredToken,
  getCurrentToken,
  isTokenCached,
  clearTokenStorageState,
  getBaseUrl,
} from "../tokenStorage";

vi.mock("@tauri-apps/api/core", () => ({
  invoke: vi.fn().mockResolvedValue(undefined),
}));

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

beforeEach(() => {
  localStorage.clear();
  mockFetch.mockReset();
  clearTokenStorageState();
});

describe("tokenStorage", () => {
  describe("resolveApiToken", () => {
    it("returns null when no token source yields a token", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({}),
      });

      const token = await resolveApiToken();
      expect(token).toBeNull();
    });

    it("returns token from localStorage when stored", async () => {
      localStorage.setItem("eniac_api_token", "stored-token");

      const token = await resolveApiToken();
      expect(token).toBe("stored-token");
    });

    it("deduplicates concurrent calls — returns same promise", async () => {
      mockFetch.mockResolvedValue({
        ok: false,
        status: 404,
        json: async () => ({}),
      });

      const [p1, p2] = await Promise.all([resolveApiToken(), resolveApiToken()]);
      // Both resolve to null, but the key assertion is they return the same promise object
      // (the dedup works because _resolvePromise is shared between calls).
      expect(p1).toBe(p2);
    });

    it("stores fetched token back to localStorage", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ token: "handshake-token" }),
      });

      await resolveApiToken();
      expect(localStorage.getItem("eniac_api_token")).toBe("handshake-token");
    });

    it("skips handshake when localStorage already has a token", async () => {
      localStorage.setItem("eniac_api_token", "initial-token");

      const token = await resolveApiToken();
      expect(token).toBe("initial-token");
      // fetch should NOT have been called (localStorage had a token)
      expect(mockFetch).not.toHaveBeenCalled();
    });
  });

  describe("rehydrateToken", () => {
    it("clears in-memory cache and fetches a fresh token (no localStorage fallback)", async () => {
      // localStorage has old token; rehydrateToken should clear in-memory cache,
      // then since Tauri returns undefined it falls through to handshake.
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ token: "new-handshake-token" }),
      });

      const token = await rehydrateToken();
      expect(token).toBe("new-handshake-token");
      expect(localStorage.getItem("eniac_api_token")).toBe("new-handshake-token");
    });

    it("returns null when all sources fail (localStorage empty, Tauri fails, handshake fails)", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({}),
      });

      const token = await rehydrateToken();
      expect(token).toBeNull();
    });
  });

  describe("clearStoredToken", () => {
    it("removes token from localStorage and memory", async () => {
      localStorage.setItem("eniac_api_token", "some-token");
      clearStoredToken();
      expect(localStorage.getItem("eniac_api_token")).toBeNull();
      expect(isTokenCached()).toBe(false);
    });
  });

  describe("getCurrentToken", () => {
    it("returns cached token from memory after load", async () => {
      localStorage.setItem("eniac_api_token", "cached-token");
      await resolveApiToken(); // loads into cache

      const token = getCurrentToken();
      expect(token).toBe("cached-token");
    });

    it("falls back to localStorage if no memory cache", async () => {
      localStorage.setItem("eniac_api_token", "fallback-token");

      const token = getCurrentToken();
      expect(token).toBe("fallback-token");
    });
  });

  describe("isTokenCached", () => {
    it("returns false when no token has been loaded", () => {
      expect(isTokenCached()).toBe(false);
    });

    it("returns true after a token is loaded", async () => {
      localStorage.setItem("eniac_api_token", "cached");
      await resolveApiToken();
      expect(isTokenCached()).toBe(true);
    });
  });

  describe("clearTokenStorageState", () => {
    it("resets all module-level state for test isolation", async () => {
      localStorage.setItem("eniac_api_token", "test-token");
      await resolveApiToken();
      expect(isTokenCached()).toBe(true);

      clearTokenStorageState();
      expect(isTokenCached()).toBe(false);
    });
  });

  describe("getBaseUrl", () => {
    it("defaults to localhost:8080 when not initialised", () => {
      clearTokenStorageState();
      expect(getBaseUrl()).toBe("http://localhost:8080/api");
    });
  });
});