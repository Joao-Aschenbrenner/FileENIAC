import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, waitFor, act } from "@testing-library/react";
import { AuthProvider, useAuth } from "../AuthProvider";
import * as tokenStorage from "../tokenStorage";

vi.mock("../tokenStorage");

function AuthProbe() {
  const auth = useAuth();
  return (
    <div>
      <span data-testid="state">{auth.state}</span>
      <button data-testid="rehydrate" onClick={() => auth.rehydrate()}>rehydrate</button>
      <button data-testid="wipe" onClick={() => auth.wipe()}>wipe</button>
    </div>
  );
}

describe("AuthProvider", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("starts in loading state", async () => {
    vi.mocked(tokenStorage.initApiClientBase).mockResolvedValue(undefined);
    vi.mocked(tokenStorage.resolveApiToken).mockImplementation(
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      () => new Promise<string | null>(() => {}),
    );

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    expect(screen.getByTestId("state").textContent).toBe("loading");
  });

  it("transitions to valid state when token resolves", async () => {
    vi.mocked(tokenStorage.initApiClientBase).mockResolvedValue(undefined);
    vi.mocked(tokenStorage.resolveApiToken).mockResolvedValue("valid-token");

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("state").textContent).toBe("valid");
    });
  });

  it("transitions to expired state when no token is found", async () => {
    vi.mocked(tokenStorage.initApiClientBase).mockResolvedValue(undefined);
    vi.mocked(tokenStorage.resolveApiToken).mockResolvedValue(null);

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("state").textContent).toBe("expired");
    });
  });

  it("rehydrate() calls tokenStorage.rehydrateToken and transitions to valid", async () => {
    vi.mocked(tokenStorage.initApiClientBase).mockResolvedValue(undefined);
    vi.mocked(tokenStorage.resolveApiToken).mockResolvedValue(null);
    vi.mocked(tokenStorage.rehydrateToken).mockResolvedValue("recovered-token");

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("state").textContent).toBe("expired");
    });

    act(() => {
      screen.getByTestId("rehydrate").click();
    });

    await waitFor(() => {
      expect(screen.getByTestId("state").textContent).toBe("valid");
    });
    expect(tokenStorage.rehydrateToken).toHaveBeenCalled();
  });

  it("wipe() calls clearStoredToken and transitions to expired", async () => {
    vi.mocked(tokenStorage.initApiClientBase).mockResolvedValue(undefined);
    vi.mocked(tokenStorage.resolveApiToken).mockResolvedValue("valid-token");

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("state").textContent).toBe("valid");
    });

    act(() => {
      screen.getByTestId("wipe").click();
    });

    expect(tokenStorage.clearStoredToken).toHaveBeenCalled();
    await waitFor(() => {
      expect(screen.getByTestId("state").textContent).toBe("expired");
    });
  });
});