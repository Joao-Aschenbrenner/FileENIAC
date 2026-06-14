import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Card } from '../Card';
import { describe, it, expect, vi } from 'vitest';

describe('Card', () => {
  it('renders title', () => {
    render(<Card title="Meu Card">content</Card>);
    expect(screen.getByText('Meu Card')).toBeInTheDocument();
  });

  it('renders subtitle', () => {
    render(<Card subtitle="Subtítulo">content</Card>);
    expect(screen.getByText('Subtítulo')).toBeInTheDocument();
  });

  it('renders children', () => {
    render(<Card><span data-testid="child">child</span></Card>);
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('calls onClick when clicked', async () => {
    const onClick = vi.fn();
    render(<Card onClick={onClick}>clickable</Card>);
    await userEvent.click(screen.getByText('clickable'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
