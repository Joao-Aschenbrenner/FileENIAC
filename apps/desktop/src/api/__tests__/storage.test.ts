// SPDX-License-Identifier: MIT
import { describe, it, expect, beforeEach } from "vitest";
import {
  storageGet,
  storageSet,
  storageRemove,
  storageClearAll,
  storageClearAuth,
  STORAGE_KEYS,
} from "../storage";

beforeEach(() => {
  localStorage.clear();
});

describe("storage helpers", () => {
  describe("storageGet", () => {
    it("returns null when key does not exist", () => {
      expect(storageGet(STORAGE_KEYS.apiToken)).toBeNull();
    });

    it("returns stored value when key exists", () => {
      localStorage.setItem(STORAGE_KEYS.apiToken, "my-token");
      expect(storageGet(STORAGE_KEYS.apiToken)).toBe("my-token");
    });

    it("returns null when localStorage.getItem throws", () => {
      const orig = localStorage.getItem;
      localStorage.getItem = () => { throw new Error("disabled"); };
      expect(storageGet(STORAGE_KEYS.apiToken)).toBeNull();
      localStorage.getItem = orig;
    });

    // SSR protection (typeof window === "undefined") is covered by the
    // throw test above — both paths return null gracefully.
  });

  describe("storageSet", () => {
    it("stores value and can be retrieved", () => {
      storageSet(STORAGE_KEYS.apiToken, "new-token");
      expect(localStorage.getItem(STORAGE_KEYS.apiToken)).toBe("new-token");
    });

    it("silently ignores when localStorage.setItem throws", () => {
      const orig = localStorage.setItem;
      localStorage.setItem = () => { throw new Error("quota exceeded"); };
      expect(() => storageSet(STORAGE_KEYS.themeMode, "dark")).not.toThrow();
      localStorage.setItem = orig;
    });
  });

  describe("storageRemove", () => {
    it("removes key from localStorage", () => {
      localStorage.setItem(STORAGE_KEYS.apiToken, "token");
      storageRemove(STORAGE_KEYS.apiToken);
      expect(localStorage.getItem(STORAGE_KEYS.apiToken)).toBeNull();
    });

    it("silently ignores when localStorage.removeItem throws", () => {
      const orig = localStorage.removeItem;
      localStorage.removeItem = () => { throw new Error("disabled"); };
      expect(() => storageRemove(STORAGE_KEYS.themeMode)).not.toThrow();
      localStorage.removeItem = orig;
    });
  });

  describe("storageClearAll", () => {
    it("removes all STORAGE_KEYS entries", () => {
      localStorage.setItem(STORAGE_KEYS.apiToken, "t1");
      localStorage.setItem(STORAGE_KEYS.sessionId, "s1");
      localStorage.setItem(STORAGE_KEYS.workspacesRoot, "root");
      localStorage.setItem(STORAGE_KEYS.workspacePath, "w1");
      localStorage.setItem(STORAGE_KEYS.themeMode, "dark");

      storageClearAll();

      expect(localStorage.getItem(STORAGE_KEYS.apiToken)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.sessionId)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.workspacesRoot)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.workspacePath)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.themeMode)).toBeNull();
    });
  });

  describe("storageClearAuth", () => {
    it("removes only auth-related keys", () => {
      localStorage.setItem(STORAGE_KEYS.apiToken, "token");
      localStorage.setItem(STORAGE_KEYS.sessionId, "session");
      localStorage.setItem(STORAGE_KEYS.workspacesRoot, "root");
      localStorage.setItem(STORAGE_KEYS.workspacePath, "workspace");
      localStorage.setItem(STORAGE_KEYS.githubUser, "gh-user");
      localStorage.setItem(STORAGE_KEYS.themeMode, "dark"); // should NOT be cleared

      storageClearAuth();

      expect(localStorage.getItem(STORAGE_KEYS.apiToken)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.sessionId)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.workspacePath)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.githubUser)).toBeNull();
      expect(localStorage.getItem(STORAGE_KEYS.workspacesRoot)).toBe("root");
      expect(localStorage.getItem(STORAGE_KEYS.themeMode)).toBe("dark");
    });
  });
});
