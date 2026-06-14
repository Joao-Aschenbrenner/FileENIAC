import { render, screen } from '@testing-library/react';
import { Timeline } from '../Timeline';
import { describe, it, expect } from 'vitest';

describe('Timeline', () => {
  it('renders empty state when no items', () => {
    render(<Timeline items={[]} />);
    expect(screen.getByText('Nenhum evento registrado.')).toBeInTheDocument();
  });

  it('renders items with titles', () => {
    const items = [{ id: 1, title: 'Deploy concluído', timestamp: '2024-01-01', type: 'DEPLOY_SUCCESS' }];
    render(<Timeline items={items} />);
    expect(screen.getByText('Deploy concluído')).toBeInTheDocument();
  });

  it('renders item timestamps', () => {
    const items = [{ id: 1, title: 'Evento', timestamp: '10/01/2024', type: 'ALERT' }];
    render(<Timeline items={items} />);
    expect(screen.getByText('10/01/2024')).toBeInTheDocument();
  });

  it('renders type badges with correct labels', () => {
    const items = [
      { id: 1, title: 'Deploy', timestamp: '2024', type: 'DEPLOY_SUCCESS' },
      { id: 2, title: 'Falha', timestamp: '2024', type: 'DEPLOY_FAILED' },
    ];
    render(<Timeline items={items} />);
    expect(screen.getByText('Deploy OK')).toBeInTheDocument();
    const badges = screen.getAllByText('Falha');
    expect(badges.length).toBe(2);
  });

  it('renders item description when provided', () => {
    const items = [{ id: 1, title: 'Test', timestamp: '2024', description: 'detalhe do evento', type: 'INFO' }];
    render(<Timeline items={items} />);
    expect(screen.getByText('detalhe do evento')).toBeInTheDocument();
  });
});
