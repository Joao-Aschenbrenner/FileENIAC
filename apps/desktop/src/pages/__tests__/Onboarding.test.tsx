// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Onboarding from '../Onboarding';

vi.mock('@tauri-apps/api/core', () => ({
  invoke: vi.fn().mockResolvedValue({
    base_url: 'http://127.0.0.1:12345/api',
    token: 'test-token',
    ready: true,
  }),
}));

vi.mock('@tauri-apps/plugin-dialog', () => ({
  open: vi.fn().mockResolvedValue('/mock/workspace'),
}));

vi.mock('../../api/client', () => ({
  configureApiClientFromBackendInfo: vi.fn().mockReturnValue(true),
  checkHealth: vi.fn().mockResolvedValue(true),
  getWorkspace: vi.fn().mockResolvedValue({
    name: 'mock-workspace',
    projects: 2,
    servers: 1,
    deploys: 0,
  }),
}));

function renderOnboarding() {
  return render(
    <MemoryRouter>
      <Onboarding />
    </MemoryRouter>,
  );
}

describe('Onboarding', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the welcome screen after startup', async () => {
    renderOnboarding();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Comecar/i })).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: /FileENIAC/i })).toBeInTheDocument();
  });

  it('moves to config step when Comecar is clicked', async () => {
    renderOnboarding();
    const btn = await screen.findByRole('button', { name: /Comecar/i });
    fireEvent.click(btn);
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /Conectar Workspace/i })).toBeInTheDocument();
    });
  });

  it('connects workspace and shows summary', async () => {
    renderOnboarding();
    const btn = await screen.findByRole('button', { name: /Comecar/i });
    fireEvent.click(btn);
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /Conectar Workspace/i })).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText(/C:\/projetos\/meu-workspace/);
    fireEvent.change(input, { target: { value: '/mock/workspace' } });

    fireEvent.click(screen.getByRole('button', { name: /Conectar/i }));
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /Workspace Conectado/i })).toBeInTheDocument();
    });
    expect(screen.getByText(/mock-workspace/)).toBeInTheDocument();
  });
});
