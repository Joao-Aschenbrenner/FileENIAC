// SPDX-License-Identifier: MIT
import { render, screen, cleanup, fireEvent } from '@testing-library/react';
import { afterEach, describe, it, expect, vi } from 'vitest';
import ErrorBoundary from '../ErrorBoundary';

function Bomb({ error }: { error: Error }): null {
  throw error;
}

function setupLocalStorage(wsPath?: string) {
  const store: Record<string, string> = {};
  if (wsPath) store['eniac_ws_path'] = wsPath;
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation((key: string) => store[key] ?? null);
}

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('ErrorBoundary', () => {
  it('renders fallback with human message (no technical details by default)', () => {
    setupLocalStorage();
    render(
      <ErrorBoundary>
        <Bomb error={new Error('disk failure')} />
      </ErrorBoundary>,
    );
    expect(screen.getByRole('heading', { name: /Ops, algo deu errado/i })).toBeInTheDocument();
    expect(screen.queryByText(/disk failure/)).not.toBeInTheDocument();
    expect(screen.getByText('Detalhes técnicos')).toBeInTheDocument();
  });

  it('shows button "Voltar ao início" when no workspace exists', () => {
    setupLocalStorage();
    render(
      <ErrorBoundary>
        <Bomb error={new Error('x')} />
      </ErrorBoundary>,
    );
    expect(screen.getByRole('button', { name: /Voltar ao início/i })).toBeInTheDocument();
  });

  it('shows button "Voltar ao dashboard" when workspace exists', () => {
    setupLocalStorage('/some/workspace');
    render(
      <ErrorBoundary>
        <Bomb error={new Error('x')} />
      </ErrorBoundary>,
    );
    expect(screen.getByRole('button', { name: /Voltar ao dashboard/i })).toBeInTheDocument();
  });

  it('shows technical details when clicking "Detalhes técnicos"', () => {
    setupLocalStorage();
    render(
      <ErrorBoundary>
        <Bomb error={new Error('disk failure')} />
      </ErrorBoundary>,
    );
    fireEvent.click(screen.getByText('Detalhes técnicos'));
    expect(screen.getByText('Detalhes técnicos:')).toBeInTheDocument();
    const matches = screen.getAllByText(/disk failure/);
    expect(matches.length).toBeGreaterThanOrEqual(1);
  });

  it('navigates to / when no workspace and button is clicked', () => {
    setupLocalStorage();
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

    fireEvent.click(screen.getByRole('button', { name: /Voltar ao início/i }));
    expect(window.location.href).toBe('/');

    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation,
    });
  });

  it('navigates to /dashboard when workspace exists', () => {
    setupLocalStorage('/some/workspace');
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

    fireEvent.click(screen.getByRole('button', { name: /Voltar ao dashboard/i }));
    expect(window.location.href).toBe('/dashboard');

    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation,
    });
  });
});
