// SPDX-License-Identifier: MIT
import { Component, ErrorInfo, ReactNode } from 'react';

interface Props { children: ReactNode; }
interface State { hasError: boolean; error?: Error; showTechnical: boolean; }

class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, showTechnical: false };
  }
  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error, showTechnical: false };
  }
  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught:', error, errorInfo);
  }

  get hasWorkspace(): boolean {
    try {
      return !!localStorage.getItem("eniac_ws_path");
    } catch { return false; }
  }

  handleBack() {
    this.setState({ hasError: false, showTechnical: false });
    if (this.hasWorkspace) {
      window.location.href = '/dashboard';
    } else {
      window.location.href = '/';
    }
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-eniac-950 text-white p-8">
          <div className="text-center max-w-lg">
            <div className="mx-auto mb-6 w-16 h-16 rounded-2xl bg-red-900/30 border border-red-700/30 flex items-center justify-center">
              <span className="text-2xl">!</span>
            </div>
            <h1 className="text-2xl font-bold mb-4">Ops, algo deu errado</h1>
            <p className="text-eniac-300 mb-6">
              Ocorreu um erro inesperado. Você pode voltar e tentar novamente.
            </p>
            <div className="flex flex-col items-center gap-3">
              <button
                onClick={() => this.handleBack()}
                className="px-6 py-2 bg-eniac-600 hover:bg-eniac-700 rounded-lg transition-colors"
              >
                {this.hasWorkspace ? "Voltar ao dashboard" : "Voltar ao início"}
              </button>
              <button
                onClick={() => this.setState((s) => ({ showTechnical: !s.showTechnical }))}
                className="text-sm text-eniac-400 hover:text-eniac-200 transition-colors"
              >
                {this.state.showTechnical ? "Ocultar detalhes técnicos" : "Detalhes técnicos"}
              </button>
            </div>
            {this.state.showTechnical && this.state.error && (
              <div className="mt-6 text-left bg-eniac-900/60 border border-eniac-700/30 rounded-lg p-4 max-h-48 overflow-y-auto">
                <p className="text-xs text-eniac-400 mb-2">Detalhes técnicos:</p>
                <p className="text-xs text-red-300 font-mono">{this.state.error.name}: {this.state.error.message}</p>
                {this.state.error.stack && (
                  <pre className="text-xs text-eniac-500 mt-2 whitespace-pre-wrap font-mono">{this.state.error.stack}</pre>
                )}
              </div>
            )}
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}
export default ErrorBoundary;
