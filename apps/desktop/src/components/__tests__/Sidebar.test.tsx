import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Sidebar from '../Sidebar';
import { describe, it, expect } from 'vitest';

function renderSidebar(initialRoute = '/') {
  return render(
    <MemoryRouter initialEntries={[initialRoute]}>
      <Sidebar />
    </MemoryRouter>
  );
}

describe('Sidebar', () => {
  it('renders all navigation links', () => {
    renderSidebar();
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Bootstrap')).toBeInTheDocument();
    expect(screen.getByText('Projetos')).toBeInTheDocument();
    expect(screen.getByText('Servidores')).toBeInTheDocument();
    expect(screen.getByText('GitHub')).toBeInTheDocument();
    expect(screen.getByText('Deploy')).toBeInTheDocument();
    expect(screen.getByText('Rollback')).toBeInTheDocument();
    expect(screen.getByText('Sync')).toBeInTheDocument();
    expect(screen.getByText('Diff')).toBeInTheDocument();
    expect(screen.getByText('Histórico')).toBeInTheDocument();
    expect(screen.getByText('Saúde')).toBeInTheDocument();
  });

  it('renders links with correct hrefs', () => {
    renderSidebar();
    expect(screen.getByText('Dashboard').closest('a')).toHaveAttribute('href', '/dashboard');
    expect(screen.getByText('Projetos').closest('a')).toHaveAttribute('href', '/projects');
    expect(screen.getByText('GitHub').closest('a')).toHaveAttribute('href', '/github/login');
    expect(screen.getByText('Histórico').closest('a')).toHaveAttribute('href', '/history');
  });

  it('highlights active route', () => {
    renderSidebar('/dashboard');
    const dashboardLink = screen.getByText('Dashboard').closest('a');
    expect(dashboardLink!.className).toContain('font-medium');
  });

  it('does not highlight inactive route', () => {
    renderSidebar('/other');
    const dashboardLink = screen.getByText('Dashboard').closest('a');
    expect(dashboardLink!.className).not.toContain('font-medium');
  });
});
