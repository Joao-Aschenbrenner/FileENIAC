// SPDX-License-Identifier: MIT
import { describe, it, expect, beforeEach, vi } from "vitest";
import { render, screen, act } from "@testing-library/react";
import { ThemeProvider, useTheme } from "../ThemeContext";
import { STORAGE_KEYS } from "../../api/storage";

function Probe() {
  const { mode, resolvedTheme } = useTheme();
  return (
    <div>
      <span data-testid="mode">{mode}</span>
      <span data-testid="resolved">{resolvedTheme}</span>
    </div>
  );
}

function renderProbe() {
  return render(
    <ThemeProvider>
      <Probe />
    </ThemeProvider>,
  );
}

beforeEach(() => {
  localStorage.clear();
  document.documentElement.classList.remove("light", "dark");
  // jsdom does not implement matchMedia; provide a deterministic stub.
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: query === "(prefers-color-scheme: dark)" ? false : false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
});

describe("ThemeContext", () => {
  it("defaults to 'system' on first load when localStorage is empty", () => {
    renderProbe();
    expect(screen.getByTestId("mode").textContent).toBe("system");
  });

  it("reads persisted 'light' theme from localStorage on mount", () => {
    localStorage.setItem(STORAGE_KEYS.themeMode, "light");
    renderProbe();
    expect(screen.getByTestId("mode").textContent).toBe("light");
  });

  it("reads persisted 'dark' theme from localStorage on mount", () => {
    localStorage.setItem(STORAGE_KEYS.themeMode, "dark");
    renderProbe();
    expect(screen.getByTestId("mode").textContent).toBe("dark");
  });

  it("reads persisted 'system' theme from localStorage on mount", () => {
    localStorage.setItem(STORAGE_KEYS.themeMode, "system");
    renderProbe();
    expect(screen.getByTestId("mode").textContent).toBe("system");
  });

  it("setMode updates the DOM (adds 'dark' class on <html>)", () => {
    function Setter() {
      const { setMode } = useTheme();
      return (
        <button data-testid="dark-btn" onClick={() => setMode("dark")}>
          dark
        </button>
      );
    }
    render(
      <ThemeProvider>
        <Setter />
        <Probe />
      </ThemeProvider>,
    );

    act(() => {
      screen.getByTestId("dark-btn").click();
    });

    expect(document.documentElement.classList.contains("dark")).toBe(true);
    expect(document.documentElement.classList.contains("light")).toBe(false);
    expect(screen.getByTestId("mode").textContent).toBe("dark");
  });

  it("setMode removes the 'dark' class when switching from dark to light", () => {
    function Setter() {
      const { setMode } = useTheme();
      return (
        <>
          <button data-testid="dark-btn" onClick={() => setMode("dark")}>d</button>
          <button data-testid="light-btn" onClick={() => setMode("light")}>l</button>
        </>
      );
    }
    render(
      <ThemeProvider>
        <Setter />
      </ThemeProvider>,
    );
    act(() => {
      screen.getByTestId("dark-btn").click();
    });
    expect(document.documentElement.classList.contains("dark")).toBe(true);

    act(() => {
      screen.getByTestId("light-btn").click();
    });
    expect(document.documentElement.classList.contains("dark")).toBe(false);
    expect(document.documentElement.classList.contains("light")).toBe(true);
  });

  it("prefers-color-scheme changes are detected when mode='system'", () => {
    // jsdom default matchMedia: jsdom's matchMedia is a no-op stub that
    // returns matches=false. We replace it with one whose `matches` and
    // the listener-firing behavior we can drive ourselves.
    type Listener = (e: MediaQueryListEvent) => void;
    const listeners = new Set<Listener>();
    // The real DOM's MediaQueryList.matches is readonly, but under tests we
    // need to mutate it to simulate the OS theme flipping. Use a typed alias
    // so the mutation is acknowledged at the boundary — this block does NOT
    // touch production code, only the test-local stub.
    interface MutableMQL {
      matches: boolean;
      media: string;
      onchange: ((e: MediaQueryListEvent) => void) | null;
      addEventListener(t: string, cb: EventListener): void;
      removeEventListener(t: string, cb: EventListener): void;
      addListener(cb: Listener): void;
      removeListener(cb: Listener): void;
      dispatchEvent(e: Event): boolean;
    }
    const mql: MutableMQL = {
      matches: false,
      media: "(prefers-color-scheme: dark)",
      onchange: null,
      addEventListener: (_t: string, cb: EventListener) => {
        listeners.add(cb as unknown as Listener);
      },
      removeEventListener: (_t: string, cb: EventListener) => {
        listeners.delete(cb as unknown as Listener);
      },
      addListener: (cb: Listener) => listeners.add(cb),
      removeListener: (cb: Listener) => listeners.delete(cb),
      dispatchEvent: () => true,
    };
    window.matchMedia = vi.fn().mockReturnValue(mql as unknown as MediaQueryList);

    renderProbe();
    expect(screen.getByTestId("resolved").textContent).toBe("light");
    expect(document.documentElement.classList.contains("light")).toBe(true);

    // Simulate the OS flipping to dark.
    mql.matches = true;
    act(() => {
      listeners.forEach((cb) => cb({ matches: true } as MediaQueryListEvent));
    });

    expect(screen.getByTestId("resolved").textContent).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });

  it("does not throw when localStorage is unavailable", () => {
    // Force all storage access to throw; the provider must still mount.
    const originalGetItem = Storage.prototype.getItem;
    const originalSetItem = Storage.prototype.setItem;
    Storage.prototype.getItem = () => {
      throw new Error("storage disabled");
    };
    Storage.prototype.setItem = () => {
      throw new Error("storage disabled");
    };
    try {
      let didThrow = false;
      try {
        renderProbe();
      } catch {
        didThrow = true;
      }
      expect(didThrow).toBe(false);
      expect(screen.getByTestId("mode").textContent).toBeDefined();
    } finally {
      Storage.prototype.getItem = originalGetItem;
      Storage.prototype.setItem = originalSetItem;
    }
  });
});
