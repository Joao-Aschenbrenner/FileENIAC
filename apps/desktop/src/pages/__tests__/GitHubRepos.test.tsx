// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import GitHubRepos from '../GitHubRepos';
import { getGitHubRepositories, importGitHubRepos } from '../../api/client';

vi.mock('../../api/client', () => ({
  getGitHubRepositories: vi.fn(),
  importGitHubRepos: vi.fn(),
  TimeoutError: class TimeoutError extends Error {
    constructor(ms: number) {
      super(`Request timed out after ${ms}ms`);
      this.name = 'TimeoutError';
    }
  },
}));

function MockTarget() {
  return <div data-testid="mock-target" />;
}

function renderWithRoutes(initialRoute = '/github/repos') {
  return render(
    <MemoryRouter initialEntries={[initialRoute]}>
      <Routes>
        <Route path="/github/repos" element={<GitHubRepos />} />
        <Route path="/github/orgs" element={<MockTarget />} />
        <Route path="/projects" element={<MockTarget />} />
      </Routes>
    </MemoryRouter>,
  );
}

const mockRepos = [
  { id: 1, name: 'repo-a', private: false, description: 'First repo', default_branch: 'main', language: 'TypeScript', imported: false },
  { id: 2, name: 'repo-b', private: true, description: 'Second repo', default_branch: 'main', language: 'Go', imported: true },
];

describe('GitHubRepos', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getGitHubRepositories).mockResolvedValue(mockRepos);
  });

  it('renders personal repos when no org param', async () => {
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('Meus Repositórios')).toBeInTheDocument();
    });
    expect(screen.getByText('repo-a')).toBeInTheDocument();
    expect(screen.getByText('repo-b')).toBeInTheDocument();
  });

  it('renders org repos when org param is present', async () => {
    renderWithRoutes('/github/repos?org=ENIACSystems');
    await waitFor(() => {
      expect(screen.getByText('Repositórios: ENIACSystems')).toBeInTheDocument();
    });
  });

  it('navigates to /github/orgs when clicking Voltar para organizações from personal repos', async () => {
    renderWithRoutes('/github/repos');
    await waitFor(() => expect(screen.getByText('Meus Repositórios')).toBeInTheDocument());

    fireEvent.click(screen.getByText('Voltar para organizações'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });

  it('navigates to /github/orgs when clicking Voltar para organizações from org repos', async () => {
    renderWithRoutes('/github/repos?org=ENIACSystems');
    await waitFor(() => expect(screen.getByText('Repositórios: ENIACSystems')).toBeInTheDocument());

    fireEvent.click(screen.getByText('Voltar para organizações'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });

  it('shows empty state when no repos', async () => {
    vi.mocked(getGitHubRepositories).mockResolvedValue([]);
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('Nenhum repositório encontrado.')).toBeInTheDocument();
    });
  });

  it('navigates to /projects after successful import', async () => {
    vi.mocked(importGitHubRepos).mockResolvedValue([{ id: 1, name: 'repo-a', error: false }]);
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('repo-a')).toBeInTheDocument());

    fireEvent.click(screen.getByText('repo-a'));
    fireEvent.click(screen.getByText('Importar (1)'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });
});
