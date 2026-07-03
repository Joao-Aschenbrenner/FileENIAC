// SPDX-License-Identifier: MIT
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import GitHubOrgs from '../GitHubOrgs';
import { getGitHubOrganizations } from '../../api/client';

vi.mock('../../api/client', () => ({
  getGitHubOrganizations: vi.fn(),
}));

function MockTarget() {
  return <div data-testid="mock-target" />;
}

function renderWithRoutes(initialRoute = '/github/orgs') {
  return render(
    <MemoryRouter initialEntries={[initialRoute]}>
      <Routes>
        <Route path="/github/orgs" element={<GitHubOrgs />} />
        <Route path="/github/repos" element={<MockTarget />} />
        <Route path="/projects" element={<MockTarget />} />
      </Routes>
    </MemoryRouter>,
  );
}

const mockOrgs = [
  { login: 'ENIACSystems', url: 'https://api.github.com/orgs/ENIACSystems' },
  { login: 'SkyMetron', url: 'https://api.github.com/orgs/SkyMetron' },
];

describe('GitHubOrgs', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getGitHubOrganizations).mockResolvedValue(mockOrgs);
  });

  it('renders organization list after loading', async () => {
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('ENIACSystems')).toBeInTheDocument();
    });
    expect(screen.getByText('SkyMetron')).toBeInTheDocument();
  });

  it('navigates to /github/repos when clicking Ver meus repositórios pessoais', async () => {
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('ENIACSystems')).toBeInTheDocument());

    fireEvent.click(screen.getByText('Ver meus repositórios pessoais'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });

  it('navigates to /projects when clicking Voltar', async () => {
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('ENIACSystems')).toBeInTheDocument());

    fireEvent.click(screen.getByRole('button', { name: /^Voltar$/ }));
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });

  it('shows empty state when no organizations', async () => {
    vi.mocked(getGitHubOrganizations).mockResolvedValue([]);
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('Nenhuma organização encontrada.')).toBeInTheDocument();
    });
  });
});
