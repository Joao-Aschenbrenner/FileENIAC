// SPDX-License-Identifier: MIT
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ErrorState } from '../ErrorState';
import { describe, it, expect, vi } from 'vitest';

describe('ErrorState', () => {
  it('renders message', () => {
    render(<ErrorState message="Algo deu errado" />);
    expect(screen.getByText('Algo deu errado')).toBeInTheDocument();
  });

  it('renders retry button and calls onClick', async () => {
    const onRetry = vi.fn();
    render(<ErrorState message="erro" onRetry={onRetry} />);
    await userEvent.click(screen.getByText('Tentar novamente'));
    expect(onRetry).toHaveBeenCalledTimes(1);
  });
});
