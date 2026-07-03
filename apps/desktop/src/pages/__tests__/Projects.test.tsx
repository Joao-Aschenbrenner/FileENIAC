import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import Projects from '../Projects';
import { listProjects, deleteProject } from '../../api/client';

vi.mock('../../api/client', () => ({
  listProjects: vi.fn(),
  deleteProject: vi.fn(),
}));

function MockTarget() {
  return <div data-testid="mock-target" />;
}

function renderWithRoutes() {
  return render(
    <MemoryRouter initialEntries={['/projects']}>
      <Routes>
        <Route path="/projects" element={<Projects />} />
        <Route path="/github/orgs" element={<MockTarget />} />
      </Routes>
    </MemoryRouter>,
  );
}

const mockProjects = [
  { id: 1, name: 'FileENIAC', local_path: '/workspace/FileENIAC', environment: 'dev', divergence_status: 'sincronizado' },
  { id: 2, name: 'my-app', local_path: '/workspace/my-app', environment: 'dev', divergence_status: 'divergente' },
];

describe('Projects', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.setItem('eniac_ws_path', '/workspace');
    vi.mocked(listProjects).mockResolvedValue(mockProjects);
  });

  it('renders project list after loading', async () => {
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('FileENIAC')).toBeInTheDocument();
    });
    expect(screen.getByText('my-app')).toBeInTheDocument();
  });

  it('shows empty state when no projects', async () => {
    vi.mocked(listProjects).mockResolvedValue([]);
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('Nenhum repositorio adicionado')).toBeInTheDocument();
    });
  });

  it('shows add repositories button', async () => {
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('FileENIAC')).toBeInTheDocument();
    });
    expect(screen.getByText('+ Adicionar Repositórios')).toBeInTheDocument();
  });

  it('navigates to /github/orgs when clicking Adicionar Repositórios', async () => {
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('FileENIAC')).toBeInTheDocument());

    fireEvent.click(screen.getByText('+ Adicionar Repositórios'));
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });

  it('navigates to /github/orgs when clicking Adicionar Repositórios in empty state', async () => {
    vi.mocked(listProjects).mockResolvedValue([]);
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('Nenhum repositorio adicionado')).toBeInTheDocument());

    const buttons = screen.getAllByText('Adicionar Repositórios');
    fireEvent.click(buttons[buttons.length - 1]);
    await waitFor(() => {
      expect(screen.getByTestId('mock-target')).toBeInTheDocument();
    });
  });

  it('shows safe remove modal with delete files checkbox', async () => {
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('FileENIAC')).toBeInTheDocument());

    fireEvent.click(screen.getAllByText('Remover')[0]);
    await waitFor(() => {
      expect(screen.getByText(/nao apaga os arquivos locais/)).toBeInTheDocument();
    });
    expect(screen.getByText('Tambem apagar arquivos locais')).toBeInTheDocument();
  });

  it('ignores technical directories', async () => {
    const withTechDirs = [
      ...mockProjects,
      { id: 3, name: '.github', local_path: '/workspace/.github', environment: 'dev', divergence_status: 'sincronizado' },
      { id: 4, name: 'node_modules', local_path: '/workspace/node_modules', environment: 'dev', divergence_status: 'sincronizado' },
    ];
    vi.mocked(listProjects).mockResolvedValue(withTechDirs);
    renderWithRoutes();
    await waitFor(() => {
      expect(screen.getByText('FileENIAC')).toBeInTheDocument();
    });
    expect(screen.queryByText('.github')).not.toBeInTheDocument();
    expect(screen.queryByText('node_modules')).not.toBeInTheDocument();
  });

  it('calls deleteProject when confirming removal', async () => {
    vi.mocked(deleteProject).mockResolvedValue({});
    renderWithRoutes();
    await waitFor(() => expect(screen.getByText('FileENIAC')).toBeInTheDocument());

    fireEvent.click(screen.getAllByText('Remover')[0]);
    await waitFor(() => expect(screen.getByText(/remove o projeto/)).toBeInTheDocument());

    const modalButtons = screen.getAllByText('Remover');
    fireEvent.click(modalButtons[modalButtons.length - 1]);
    await waitFor(() => {
      expect(deleteProject).toHaveBeenCalledWith('/workspace', 'FileENIAC');
    });
  });
});
