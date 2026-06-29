// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, act, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { STORAGE_KEYS } from "../../../api/storage";

function tickCheckbox(checkbox: HTMLElement) {
  // cast through HTMLElement — getByLabelText returns the bound element,
  // which for a native <input type="checkbox"> is an HTMLInputElement.
  // We coerce via unknown so the cast is explicit at the boundary.
  const input = checkbox as unknown as HTMLInputElement;
  // React's controlled <input type=checkbox>: the native click DOES toggle
  // the checked state because React 18 re-syncs from the synthetic onChange
  // event payload that the click triggers. We funnel through fireEvent.click
  // and wrap in act() to make sure subsequent assertions see the update.
  act(() => {
    fireEvent.click(input);
  });
}

// SessionProvider + many lazy pages import from "../api/client". Gate the
// network surface so the gate flow remains hermetic.
vi.mock("../../api/client", () => ({
  checkHealth: vi.fn().mockResolvedValue(true),
  listSessions: vi.fn().mockResolvedValue([]),
  activateSession: vi.fn().mockResolvedValue({}),
  clearSessionWorkspace: vi.fn().mockResolvedValue({}),
  deleteSession: vi.fn().mockResolvedValue({}),
  listProjects: vi.fn().mockResolvedValue([]),
  getProject: vi.fn().mockResolvedValue({}),
  createProject: vi.fn().mockResolvedValue({}),
  deleteProject: vi.fn().mockResolvedValue({}),
  listServers: vi.fn().mockResolvedValue([]),
  getServer: vi.fn().mockResolvedValue({}),
  createServer: vi.fn().mockResolvedValue({}),
  deleteServer: vi.fn().mockResolvedValue({}),
  getSettings: vi.fn().mockResolvedValue({}),
  updateSettings: vi.fn().mockResolvedValue({}),
  getHistory: vi.fn().mockResolvedValue([]),
  getEvents: vi.fn().mockResolvedValue([]),
  getDeploys: vi.fn().mockResolvedValue([]),
  executeDeploy: vi.fn().mockResolvedValue({}),
  executeRollback: vi.fn().mockResolvedValue({}),
  executeVerify: vi.fn().mockResolvedValue({}),
  getDiff: vi.fn().mockResolvedValue({}),
  getSyncs: vi.fn().mockResolvedValue([]),
  executeSync: vi.fn().mockResolvedValue({}),
  executeSyncSafe: vi.fn().mockResolvedValue({}),
  executeSyncWithDelete: vi.fn().mockResolvedValue({}),
  createMirror: vi.fn().mockResolvedValue({}),
  getHealthCheck: vi.fn().mockResolvedValue({}),
  getGitHubStatus: vi.fn().mockResolvedValue({}),
  gitHubLogin: vi.fn().mockResolvedValue({}),
  gitHubLogout: vi.fn().mockResolvedValue({}),
  getGitHubOrganizations: vi.fn().mockResolvedValue([]),
  getGitHubRepositories: vi.fn().mockResolvedValue([]),
  importGitHubRepos: vi.fn().mockResolvedValue([]),
  cloneGitHubRepo: vi.fn().mockResolvedValue({}),
  listRepositories: vi.fn().mockResolvedValue([]),
  getRepository: vi.fn().mockResolvedValue({}),
  createSession: vi.fn().mockResolvedValue({}),
  updateSession: vi.fn().mockResolvedValue({}),
  getSession: vi.fn().mockResolvedValue({}),
  getWorkspace: vi.fn().mockResolvedValue({}),
  closeWorkspace: vi.fn().mockResolvedValue({}),
  initApiClient: vi.fn().mockResolvedValue(undefined),
  heartbeat: vi.fn().mockResolvedValue(undefined),
}));

import App from "../../../App";

beforeEach(() => {
  localStorage.clear();
});

function mountApp(initialEntry = "/") {
  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <App />
    </MemoryRouter>,
  );
}

describe("LGPDConsent gate behavior", () => {
  it("renders the gate when no eniac_lgpd_consent is in localStorage", async () => {
    mountApp();
    // Gate gate step is "intro" — heading "Bienvenido ao FileENIAC" is the
    // distinctive copy unique to this screen.
    await waitFor(() => {
      expect(
        screen.getByRole("heading", { name: /Bienvenido ao FileENIAC/i }),
      ).toBeInTheDocument();
    });
    // Footer "Termos de Uso" / "Política de Privacidade" call-to-action
    // confirms the LGPD landing rendered, not the SessionSelector.
    expect(screen.getByText("Termos de Uso")).toBeInTheDocument();
    expect(screen.getByText("Política de Privacidade")).toBeInTheDocument();
  });

  it("skips the gate when consent has already been given", async () => {
    localStorage.setItem(
      STORAGE_KEYS.lgpdConsent,
      JSON.stringify({
        agreed: true,
        termsAccepted: true,
        privacyAccepted: true,
        dataProcessingAccepted: true,
        consentedAt: new Date().toISOString(),
      }),
    );
    mountApp();
    // The SessionSelector at "/" shows the "Nova sessão" CTA; we use a
    // distinctive heading to confirm we passed the gate.
    await waitFor(() => {
      // SessionSelector has a heading "Sessões" / "Workspace"; gate intro
      // heading must NOT be present.
      expect(
        screen.queryByRole("heading", { name: /Bienvenido ao FileENIAC/i }),
      ).not.toBeInTheDocument();
    });
  });

  it("after accepting all 3 checkboxes, writes localStorage eniac_lgpd_consent with agreed:true and an ISO consentedAt", async () => {
    mountApp();

    // Step 1: intro. Click the "Ler e Aceitar Termos" button -> "terms".
    await waitFor(() =>
      screen.getByRole("button", { name: /Ler e Aceitar Termos/i }),
    );
    fireEvent.click(screen.getByRole("button", { name: /Ler e Aceitar Termos/i }));

    // Step 2: terms step. Tick termsAccepted then click "Continuar para Privacidade".
    await waitFor(() =>
      screen.getByRole("heading", { name: /Termos de Uso/i }),
    );
    const termsCheckbox = screen.getByLabelText(/Li e aceito os Termos de Uso/i);
    tickCheckbox(termsCheckbox);

    const continueBtn = screen.getByRole("button", {
      name: /Continuar para Privacidade/i,
    });
    // Sanity: not disabled once the box is checked
    expect(continueBtn).not.toBeDisabled();
    fireEvent.click(continueBtn);

    // Step 3: privacy step. Tick BOTH remaining boxes then click "Concluir".
    await waitFor(() =>
      screen.getByRole("heading", { name: /Pol[ií]tica de Privacidade/i }),
    );

    const privacyCheckbox = screen.getByLabelText(
      /Li e aceito a Pol[ií]tica de Privacidade/i,
    );
    const dataCheckbox = screen.getByLabelText(
      /Consinto com o armazenamento e processamento/i,
    );

    tickCheckbox(privacyCheckbox);
    tickCheckbox(dataCheckbox);

    const concludeBtn = screen.getByRole("button", { name: /Concluir/i });
    expect(concludeBtn).not.toBeDisabled();

    await act(async () => {
      fireEvent.click(concludeBtn);
    });

    const stored = localStorage.getItem(STORAGE_KEYS.lgpdConsent);
    expect(stored).not.toBeNull();
    const parsed = JSON.parse(stored!);
    expect(parsed.agreed).toBe(true);
    expect(parsed.termsAccepted).toBe(true);
    expect(parsed.privacyAccepted).toBe(true);
    expect(parsed.dataProcessingAccepted).toBe(true);
    expect(typeof parsed.consentedAt).toBe("string");
    // Round-trip check: must be a valid ISO date.
    const parsedDate = new Date(parsed.consentedAt);
    expect(Number.isNaN(parsedDate.getTime())).toBe(false);
  });

  it("does not store consent if any checkbox is unchecked (Concluir button stays disabled)", async () => {
    mountApp();
    // Navigate to privacy step.
    await waitFor(() =>
      screen.getByRole("button", { name: /Ler e Aceitar Termos/i }),
    );
    fireEvent.click(screen.getByRole("button", { name: /Ler e Aceitar Termos/i }));

    await waitFor(() =>
      screen.getByRole("heading", { name: /Termos de Uso/i }),
    );
    tickCheckbox(screen.getByLabelText(/Li e aceito os Termos de Uso/i));
    fireEvent.click(
      screen.getByRole("button", { name: /Continuar para Privacidade/i }),
    );

    await waitFor(() =>
      screen.getByRole("heading", { name: /Pol[ií]tica de Privacidade/i }),
    );

    const concludeBtn = screen.getByRole("button", { name: /Concluir/i });
    // All three unchecked -> disabled. (termsAccepted IS already true from
    // the previous step, but privacy + data are still false here.)
    expect(concludeBtn).toBeDisabled();

    // Tick only 1 of the 2 remaining boxes — still disabled.
    tickCheckbox(
      screen.getByLabelText(/Li e aceito a Pol[ií]tica de Privacidade/i),
    );
    expect(concludeBtn).toBeDisabled();

    // localStorage must NOT have been written yet — disabled button means
    // a real user cannot reach handleComplete().
    expect(localStorage.getItem(STORAGE_KEYS.lgpdConsent)).toBeNull();

    // Tick the second remaining box — now all three true, button enables.
    tickCheckbox(
      screen.getByLabelText(/Consinto com o armazenamento e processamento/i),
    );
    expect(concludeBtn).not.toBeDisabled();

    // Untick one again — disabled again, demonstrating the gate is enforced
    // by the canProceed flag and not just by initial render.
    tickCheckbox(
      screen.getByLabelText(/Consinto com o armazenamento e processamento/i),
    );
    expect(concludeBtn).toBeDisabled();
  });
});
