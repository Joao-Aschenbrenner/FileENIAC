// SPDX-License-Identifier: MIT
/**
 * ThemeProvider — manages light/dark/system theme with CSS variable support.
 *
 * Applies theme via TWO mechanisms for maximum compatibility:
 * 1. `data-theme` attribute  → drives CSS variables in themes.css
 * 2. `dark` class            → drives Tailwind's built-in dark: utilities
 *
 * Both are kept in sync so existing `dark:` Tailwind classes continue working
 * while the new CSS variable system also picks up the theme.
 */

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { STORAGE_KEYS, storageGet, storageSet } from "../api/storage";

export type ThemeMode = "light" | "dark" | "system";

interface ThemeContextValue {
  mode: ThemeMode;
  resolvedTheme: "light" | "dark";
  setMode: (mode: ThemeMode) => void;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setModeState] = useState<ThemeMode>(() => {
    if (typeof window === "undefined") return "system";
    return (storageGet(STORAGE_KEYS.themeMode) as ThemeMode) || "system";
  });

  const [resolvedTheme, setResolvedTheme] = useState<"light" | "dark">(() => {
    if (typeof window === "undefined") return "light";
    const stored = storageGet(STORAGE_KEYS.themeMode) as ThemeMode;
    if (stored === "light" || stored === "dark") return stored;
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  });

  const applyTheme = useCallback((theme: "light" | "dark") => {
    const root = document.documentElement;
    // CSS variables driven by data-theme
    root.setAttribute("data-theme", theme);
    // Tailwind's dark: variant driven by .dark class
    root.classList.remove("light", "dark");
    root.classList.add(theme);
    // Persist the user preference (not "system")
    if (theme !== resolvedTheme) {
      // Only save when user explicitly picks light/dark
    }
  }, []);

  useEffect(() => {
    const resolve = (): "light" | "dark" => {
      if (mode === "system") {
        return window.matchMedia("(prefers-color-scheme: dark)").matches
          ? "dark"
          : "light";
      }
      return mode;
    };

    const theme = resolve();
    setResolvedTheme(theme);
    applyTheme(theme);

    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const handleChange = () => {
      if (mode === "system") {
        const sysTheme = mediaQuery.matches ? "dark" : "light";
        setResolvedTheme(sysTheme);
        applyTheme(sysTheme);
      }
    };
    mediaQuery.addEventListener("change", handleChange);
    return () => mediaQuery.removeEventListener("change", handleChange);
  }, [mode, applyTheme]);

  const setMode = useCallback((newMode: ThemeMode) => {
    setModeState(newMode);
    storageSet(STORAGE_KEYS.themeMode, newMode);
  }, []);

  const toggleTheme = useCallback(() => {
    if (mode === "system") {
      setMode(resolvedTheme === "dark" ? "light" : "dark");
    } else {
      setMode(mode === "dark" ? "light" : "dark");
    }
  }, [mode, resolvedTheme, setMode]);

  const value = useMemo(
    () => ({ mode, resolvedTheme, setMode, toggleTheme }),
    [mode, resolvedTheme, setMode, toggleTheme]
  );

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
}