import { TimeoutError } from "../api/client";
import { ApiError } from "../api/errors";

export interface ClassifiedError {
  title: string;
  description: string;
  actionLabel?: string;
  technical: string;
}

export function classifyError(err: unknown): ClassifiedError {
  if (err instanceof TimeoutError) {
    return {
      title: "A operação demorou muito",
      description: "O servidor não respondeu dentro do tempo esperado. Pode ser que a operação tenha sido concluída em segundo plano.",
      actionLabel: "Tentar novamente",
      technical: `Timeout after ${err.ms}ms`,
    };
  }

  if (err instanceof ApiError) {
    if (err.isTimeout()) {
      return {
        title: "A operação demorou muito",
        description: "O servidor não respondeu dentro do tempo esperado.",
        actionLabel: "Tentar novamente",
        technical: `API timeout: ${err.message}`,
      };
    }
    if (err.isUnauthorized()) {
      return {
        title: "Sessão expirada",
        description: "Sua sessão expirou. Faça login novamente para continuar.",
        actionLabel: "Fazer login",
        technical: `HTTP 401: ${err.message}`,
      };
    }
    if (err.isForbidden()) {
      return {
        title: "Acesso negado",
        description: "Você não tem permissão para realizar esta operação.",
        technical: `HTTP 403: ${err.message}`,
      };
    }
    if (err.status >= 500) {
      return {
        title: "Erro no servidor",
        description: "O servidor encontrou um erro interno. Tente novamente em alguns instantes.",
        actionLabel: "Tentar novamente",
        technical: `HTTP ${err.status}: ${err.message}`,
      };
    }
    return {
      title: "Erro na requisição",
      description: "Não foi possível completar a operação.",
      actionLabel: "Tentar novamente",
      technical: `API error: ${err.message}`,
    };
  }

  if (err instanceof TypeError && err.message?.includes("Failed to fetch")) {
    return {
      title: "Servidor não encontrado",
      description: "Não foi possível conectar ao servidor. Verifique se o FileENIAC está em execução.",
      actionLabel: "Verificar conexão",
      technical: "Network error: Failed to fetch",
    };
  }

  if (err instanceof TypeError && err.message?.includes("null")) {
    return {
      title: "Erro inesperado",
      description: "Ocorreu um erro inesperado ao carregar os dados. Tente novamente.",
      actionLabel: "Tentar novamente",
      technical: `Null reference: ${err.message}`,
    };
  }

  const msg = err instanceof Error ? err.message : String(err);
  return {
    title: "Erro inesperado",
    description: "Ocorreu um erro inesperado. Tente novamente.",
    actionLabel: "Tentar novamente",
    technical: msg,
  };
}
