import { describe, it, expect } from 'vitest';
import { classifyError } from '../errorClassifier';
import { ApiError } from '../../api/errors';
import { TimeoutError } from '../../api/client';

describe('classifyError', () => {
  it('classifies TimeoutError', () => {
    const result = classifyError(new TimeoutError(5000));
    expect(result.title).toBe('A operação demorou muito');
    expect(result.technical).toContain('5000');
    expect(result.actionLabel).toBe('Tentar novamente');
  });

  it('classifies ApiError 401', () => {
    const err = new ApiError(401, '/api/test', 'Unauthorized');
    const result = classifyError(err);
    expect(result.title).toBe('Sessão expirada');
    expect(result.actionLabel).toBe('Fazer login');
  });

  it('classifies ApiError 403', () => {
    const err = new ApiError(403, '/api/test', 'Forbidden');
    const result = classifyError(err);
    expect(result.title).toBe('Acesso negado');
  });

  it('classifies ApiError 500', () => {
    const err = new ApiError(500, '/api/test', 'Internal error');
    const result = classifyError(err);
    expect(result.title).toBe('Erro no servidor');
    expect(result.actionLabel).toBe('Tentar novamente');
  });

  it('classifies generic ApiError', () => {
    const err = new ApiError(400, '/api/test', 'Bad request');
    const result = classifyError(err);
    expect(result.title).toBe('Erro na requisição');
  });

  it('classifies TypeError "Failed to fetch"', () => {
    const err = new TypeError('Failed to fetch');
    const result = classifyError(err);
    expect(result.title).toBe('Servidor não encontrado');
  });

  it('classifies TypeError with null', () => {
    const err = new TypeError("Cannot read properties of null (reading 'length')");
    const result = classifyError(err);
    expect(result.title).toBe('Erro inesperado');
    expect(result.technical).toContain('null');
  });

  it('classifies unknown error', () => {
    const result = classifyError('raw string');
    expect(result.title).toBe('Erro inesperado');
    expect(result.technical).toBe('raw string');
  });
});
