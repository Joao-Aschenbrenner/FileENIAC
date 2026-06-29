// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Onboarding from '../Onboarding';

vi.mock('@tauri-apps/plugin-dialog', () => ({
  open: vi.fn().mockResolvedValue('/mock/workspace'),
}));

vi.mock('../../api/client', () => ({
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

  it('renders the welcome screen', () => {
    renderOnboarding();
    expect(screen.getByRole('heading', { name: /FileENIAC/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Começar/i })).toBeInTheDocument();
  });

  it('moves to config step when backend is healthy', async () => {
    renderOnboarding();
    fireEvent.click(screen.getByRole('button', { name: /Começar/i }));
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /Conectar Workspace/i })).toBeInTheDocument();
    });
  });

  it('connects workspace and shows summary', async () => {
    renderOnboarding();
    fireEvent.click(screen.getByRole('button', { name: /Começar/i }));
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
