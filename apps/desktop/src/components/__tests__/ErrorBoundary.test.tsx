// SPDX-License-Identifier: MIT
import { render, screen, cleanup, fireEvent } from '@testing-library/react';
import { afterEach, describe, it, expect, vi } from 'vitest';
import ErrorBoundary from '../ErrorBoundary';
import { ApiError } from '../../api/errors';
import { STORAGE_KEYS } from '../../api/storage';

function Bomb({ error }: { error: Error }): null {
  throw error;
}

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('ErrorBoundary 401 handling', () => {
  it('renders the friendly "Sessão expirada" copy for HTTP 401 errors', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new ApiError(401, '/api/projects', 'token expired')} />
      </ErrorBoundary>,
    );
    expect(screen.getByRole('heading', { name: /sessão expirada/i })).toBeInTheDocument();
    expect(screen.getByText(/Sessão expirada. Reinicie o aplicativo/i)).toBeInTheDocument();
  });

  it('swaps the primary button to "Voltar à seleção de sessões" on 401', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new ApiError(401, '/api/projects', 'unauthorized')} />
      </ErrorBoundary>,
    );
    expect(
      screen.getByRole('button', { name: /voltar à seleção de sessões/i }),
    ).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /^voltar ao início$/i })).not.toBeInTheDocument();
  });

  it('clears the stored auth token/session IDs when the 401 button is clicked', () => {
    localStorage.setItem(STORAGE_KEYS.apiToken, 't');
    localStorage.setItem(STORAGE_KEYS.sessionId, '7');
    localStorage.setItem(STORAGE_KEYS.workspacePath, '/ws');

    // Stub window.location so we don't actually navigate.
    const originalLocation = window.location;
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { ...originalLocation, href: '' },
    });

    render(
      <ErrorBoundary>
        <Bomb error={new ApiError(401, '/api/projects', 'unauthorized')} />
      </ErrorBoundary>,
    );
    fireEvent.click(screen.getByRole('button', { name: /voltar à seleção de sessões/i }));

    expect(localStorage.getItem(STORAGE_KEYS.apiToken)).toBeNull();
    expect(localStorage.getItem(STORAGE_KEYS.sessionId)).toBeNull();
    expect(localStorage.getItem(STORAGE_KEYS.workspacePath)).toBeNull();

    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation,
    });
  });

  it('falls back to the generic copy for non-401 errors', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new Error('HTTP 500 disk failure')} />
      </ErrorBoundary>,
    );
    expect(screen.getByRole('heading', { name: /erro inesperado/i })).toBeInTheDocument();
    expect(screen.getByText(/encontrou um erro interno/i)).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: /^voltar ao início$/i }),
    ).toBeInTheDocument();
  });

  it('falls back to the generic copy for null-pointer-shaped errors', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new Error("Cannot read properties of undefined (reading 'x')")} />
      </ErrorBoundary>,
    );
    expect(screen.getByText(/servidor retornou dados inesperados/i)).toBeInTheDocument();
  });

  it('renders raw message text when no friendly mapping exists', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new Error('mysterious gremlin')} />
      </ErrorBoundary>,
    );
    expect(screen.getByText('mysterious gremlin')).toBeInTheDocument();
  });

  it('"Tentar Novamente" exists as an affordance even on 401 (auth fallback path)', () => {
    render(
      <ErrorBoundary>
        <Bomb error={new ApiError(401, '/api/projects', 'unauthorized')} />
      </ErrorBoundary>,
    );
    // The button offers a recovery path regardless of error type.
    expect(screen.getByRole('button', { name: /tentar novamente/i })).toBeInTheDocument();
  });
});
