// SPDX-License-Identifier: MIT
import { render, screen, cleanup, fireEvent } from '@testing-library/react';
import { afterEach, describe, it, expect, vi } from 'vitest';
import ErrorBoundary from '../ErrorBoundary';

function Bomb({ error }: { error: Error }): null {
  throw error;
}

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('ErrorBoundary', () => {
  it('renders the generic error fallback', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new Error('disk failure')} />
      </ErrorBoundary>,
    );
    expect(screen.getByRole('heading', { name: /erro inesperado/i })).toBeInTheDocument();
    expect(screen.getByText(/disk failure/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /voltar ao inicio/i })).toBeInTheDocument();
  });

  it('renders the fallback copy when the error message is empty', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new Error('')} />
      </ErrorBoundary>,
    );
    expect(screen.getByText(/Ocorreu um erro inesperado./)).toBeInTheDocument();
  });

  it('resets error state and navigates home when the button is clicked', () => {
    const originalLocation = window.location;
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { ...originalLocation, href: '' },
    });

    render(
      <ErrorBoundary>
        <Bomb error={new Error('boom')} />
      </ErrorBoundary>,
    );

    fireEvent.click(screen.getByRole('button', { name: /voltar ao inicio/i }));
    expect(window.location.href).toBe('/');

    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation,
    });
  });
});
