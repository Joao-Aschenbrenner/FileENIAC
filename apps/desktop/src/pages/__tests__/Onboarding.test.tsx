// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Onboarding from '../Onboarding';
import { createWorkspace, enterWorkspace, listWorkspaces } from '../../api/client';

vi.mock('@tauri-apps/api/core', () => ({
  invoke: vi.fn().mockResolvedValue({
    base_url: 'http://127.0.0.1:12345/api',
    token: 'test-token',
    ready: true,
  }),
}));

vi.mock('@tauri-apps/plugin-dialog', () => ({
  open: vi.fn().mockResolvedValue('/mock/base'),
}));

vi.mock('../../api/client', () => ({
  configureApiClientFromBackendInfo: vi.fn().mockReturnValue(true),
  checkHealth: vi.fn().mockResolvedValue(true),
  listWorkspaces: vi.fn().mockResolvedValue([]),
  createWorkspace: vi.fn().mockResolvedValue({
    name: 'Cliente A',
    path: '/mock/base/Cliente-A',
    projects: 0,
    servers: 0,
    deploys: 0,
  }),
  enterWorkspace: vi.fn().mockResolvedValue({
    name: 'Cliente A',
    path: '/mock/base/Cliente-A',
  }),
}));

function renderOnboarding() {
  return render(
    <MemoryRouter>
      <Onboarding />
    </MemoryRouter>,
  );
}

async function goToWorkspaceManager() {
  const start = await screen.findByRole('button', { name: /Começar/i });
  fireEvent.click(start);
  await waitFor(() => {
    expect(screen.getByRole('heading', { name: /Escolher Pasta dos Workspaces/i })).toBeInTheDocument();
  });

  fireEvent.change(screen.getByPlaceholderText(/C:\/projetos\/ENIAC_SYSTEMS/i), {
    target: { value: '/mock/base' },
  });
  fireEvent.click(screen.getByRole('button', { name: /Continuar/i }));

  await waitFor(() => {
    expect(screen.getByRole('heading', { name: /Criar ou Entrar em um Workspace/i })).toBeInTheDocument();
  });
}

describe('Onboarding', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    vi.mocked(listWorkspaces).mockResolvedValue([]);
    vi.mocked(createWorkspace).mockResolvedValue({
      name: 'Cliente A',
      path: '/mock/base/Cliente-A',
      projects: 0,
      servers: 0,
      deploys: 0,
    });
    vi.mocked(enterWorkspace).mockResolvedValue({
      name: 'Cliente A',
      path: '/mock/base/Cliente-A',
    });
  });

  it('renders the welcome screen after startup', async () => {
    renderOnboarding();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Começar/i })).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: /FileENIAC/i })).toBeInTheDocument();
  });

  it('stores the base folder and lists workspaces without creating a workspace at the base', async () => {
    renderOnboarding();
    await goToWorkspaceManager();

    expect(localStorage.getItem('eniac_workspaces_root')).toBe('/mock/base');
    expect(localStorage.getItem('eniac_ws_path')).toBeNull();
    expect(listWorkspaces).toHaveBeenCalledWith('/mock/base');
    expect(createWorkspace).not.toHaveBeenCalled();
  });

  it('creates a workspace inside the base folder and enters it only when selected', async () => {
    renderOnboarding();
    await goToWorkspaceManager();

    fireEvent.change(screen.getByPlaceholderText(/Cliente X/i), { target: { value: 'Cliente A' } });
    fireEvent.click(screen.getByRole('button', { name: /Criar Workspace/i }));

    await waitFor(() => {
      expect(createWorkspace).toHaveBeenCalledWith('/mock/base/Cliente-A', { name: 'Cliente A', description: '' });
    });
    expect(localStorage.getItem('eniac_ws_path')).toBeNull();
    expect(screen.getByText(/Workspace criado/i)).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /Entrar/i }));
    await waitFor(() => {
      expect(enterWorkspace).toHaveBeenCalledWith('/mock/base/Cliente-A');
    });
  });

  it('shows existing workspaces from the selected base folder', async () => {
    vi.mocked(listWorkspaces).mockResolvedValue([
      { name: 'Cliente B', path: '/mock/base/ClienteB' },
    ]);

    renderOnboarding();
    await goToWorkspaceManager();

    expect(screen.getByText('Cliente B')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /Entrar/i }));
    await waitFor(() => {
      expect(enterWorkspace).toHaveBeenCalledWith('/mock/base/ClienteB');
    });
  });
});
