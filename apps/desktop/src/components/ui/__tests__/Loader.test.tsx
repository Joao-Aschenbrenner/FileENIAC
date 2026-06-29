// SPDX-License-Identifier: MIT
import { render, screen } from '@testing-library/react';
import { Loader } from '../Loader';
import { describe, it, expect } from 'vitest';

describe('Loader', () => {
  it('renders default loading text', () => {
    render(<Loader />);
    expect(screen.getByText('Carregando...')).toBeInTheDocument();
  });

  it('renders custom text', () => {
    render(<Loader text="Aguarde..." />);
    expect(screen.getByText('Aguarde...')).toBeInTheDocument();
  });
});
