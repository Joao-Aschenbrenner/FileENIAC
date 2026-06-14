import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { EmptyState } from '../EmptyState';
import { describe, it, expect, vi } from 'vitest';

describe('EmptyState', () => {
  it('renders title', () => {
    render(<EmptyState title="Nada aqui" />);
    expect(screen.getByText('Nada aqui')).toBeInTheDocument();
  });

  it('renders description', () => {
    render(<EmptyState title="vazio" description="descrição do estado vazio" />);
    expect(screen.getByText('descrição do estado vazio')).toBeInTheDocument();
  });

  it('renders action button and calls onClick', async () => {
    const onClick = vi.fn();
    render(<EmptyState title="vazio" action={{ label: "Criar", onClick }} />);
    await userEvent.click(screen.getByText('Criar'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
