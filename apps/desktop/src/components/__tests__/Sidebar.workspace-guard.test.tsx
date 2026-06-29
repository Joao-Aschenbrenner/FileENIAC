// SPDX-License-Identifier: MIT
import { render, screen, cleanup } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, describe, it, expect, vi } from 'vitest';
import Sidebar from '../Sidebar';
import type { Session } from '../../types';

// Mock the session context module so we can vary the workspace/session
// state independently of a real SessionProvider.
vi.mock('../../context/SessionContext', () => ({
  useSession: vi.fn(),
}));

import { useSession } from '../../context/SessionContext';
const useSessionMock = vi.mocked(useSession);

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

function renderSidebar() {
  return render(
    <MemoryRouter>
      <Sidebar />
    </MemoryRouter>,
  );
}

const baseSession: Session = {
  id: 1,
  name: 'Demo',
  description: '',
  workspace_path: '/some/workspace',
  is_active: true,
};

describe('Sidebar M-11 workspace guard', () => {
  it('shows the warning banner when no active session exists', () => {
    useSessionMock.mockReturnValue({
      activeSession: null,
      workspacePath: '',
      loading: false,
      sessions: [],
      error: '',
      authExpired: false,
      backendOnline: false,
      authState: "valid",
      refresh: async () => {},
      switchSession: async () => {},
      removeWorkspace: async () => {},
      clearAuthExpired: () => {},
    });
    renderSidebar();
    const banner = screen.getByTestId('sidebar-workspace-warning');
    expect(banner).toBeInTheDocument();
    expect(banner).toHaveTextContent('Nenhuma sessão ativa');
  });

  it('shows the warning banner when session has empty workspace_path', () => {
    useSessionMock.mockReturnValue({
      activeSession: { ...baseSession, workspace_path: '' },
      workspacePath: '',
      loading: false,
      sessions: [baseSession],
      error: '',
      authExpired: false,
      backendOnline: false,
      authState: "valid",
      refresh: async () => {},
      switchSession: async () => {},
      removeWorkspace: async () => {},
      clearAuthExpired: () => {},
    });
    renderSidebar();
    const banner = screen.getByTestId('sidebar-workspace-warning');
    expect(banner).toBeInTheDocument();
    expect(banner).toHaveTextContent('Workspace inválido');
  });

  it('shows the warning banner when workspacePath is just whitespace', () => {
    useSessionMock.mockReturnValue({
      activeSession: { ...baseSession, workspace_path: '   ' },
      workspacePath: '   ',
      loading: false,
      sessions: [baseSession],
      error: '',
      authExpired: false,
      backendOnline: false,
      authState: "valid",
      refresh: async () => {},
      switchSession: async () => {},
      removeWorkspace: async () => {},
      clearAuthExpired: () => {},
    });
    renderSidebar();
    expect(screen.getByTestId('sidebar-workspace-warning')).toBeInTheDocument();
  });

  it('does NOT show the banner when a session has a valid workspace', () => {
    useSessionMock.mockReturnValue({
      activeSession: baseSession,
      workspacePath: '/some/workspace',
      loading: false,
      sessions: [baseSession],
      error: '',
      authExpired: false,
      backendOnline: false,
      authState: "valid",
      refresh: async () => {},
      switchSession: async () => {},
      removeWorkspace: async () => {},
      clearAuthExpired: () => {},
    });
    renderSidebar();
    expect(screen.queryByTestId('sidebar-workspace-warning')).not.toBeInTheDocument();
  });

  it('still hides the banner while the session is loading', () => {
    useSessionMock.mockReturnValue({
      activeSession: null,
      workspacePath: '',
      loading: true,
      sessions: [],
      error: '',
      authExpired: false,
      backendOnline: false,
      authState: "valid",
      refresh: async () => {},
      switchSession: async () => {},
      removeWorkspace: async () => {},
      clearAuthExpired: () => {},
    });
    renderSidebar();
    expect(screen.queryByTestId('sidebar-workspace-warning')).not.toBeInTheDocument();
  });

  it('marks the navigation inert via aria-disabled when warning is up', () => {
    useSessionMock.mockReturnValue({
      activeSession: null,
      workspacePath: '',
      loading: false,
      sessions: [],
      error: '',
      authExpired: false,
      backendOnline: false,
      authState: "valid",
      refresh: async () => {},
      switchSession: async () => {},
      removeWorkspace: async () => {},
      clearAuthExpired: () => {},
    });
    renderSidebar();
    const nav = screen.getByRole('navigation');
    expect(nav).toHaveAttribute('aria-disabled', 'true');
    // Links should be present but unreachable (tabIndex=-1).
    const link = screen.getByText('Dashboard').closest('a');
    expect(link).toHaveAttribute('tabindex', '-1');
  });
});
