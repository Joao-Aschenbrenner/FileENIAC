import { render, screen } from '@testing-library/react';
import { Toast } from '../Toast';
import { describe, it, expect, vi, afterEach } from 'vitest';

describe('Toast', () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders message', () => {
    render(<Toast message="Operação concluída" onClose={() => {}} />);
    expect(screen.getByText('Operação concluída')).toBeInTheDocument();
  });

  it('auto-dismisses after duration plus fade delay', () => {
    vi.useFakeTimers();
    const onClose = vi.fn();
    render(<Toast message="teste" onClose={onClose} duration={2000} />);
    vi.advanceTimersByTime(2000 + 300);
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
