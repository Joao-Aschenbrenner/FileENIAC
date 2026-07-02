import { describe, it, expect } from 'vitest';
import { AppError, UserFacingError } from '../AppError';

describe('AppError', () => {
  it('creates error with user-facing fields', () => {
    const err = new AppError('disk failure', 'Erro no disco', 'O disco falhou');
    expect(err.message).toBe('disk failure');
    expect(err.userTitle).toBe('Erro no disco');
    expect(err.userDescription).toBe('O disco falhou');
    expect(err.name).toBe('AppError');
  });
});

describe('UserFacingError', () => {
  it('creates error with user-facing message', () => {
    const err = new UserFacingError('Título', 'Descrição');
    expect(err.userTitle).toBe('Título');
    expect(err.userDescription).toBe('Descrição');
    expect(err.message).toBe('Título: Descrição');
  });
});
