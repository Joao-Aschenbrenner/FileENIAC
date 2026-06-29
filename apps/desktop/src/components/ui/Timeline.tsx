// SPDX-License-Identifier: MIT
import { Badge } from "./Badge";

interface TimelineItem {
  id: number | string;
  title: string;
  description?: string;
  timestamp: string;
  type?: string;
}

interface TimelineProps {
  items: TimelineItem[];
}

const typeLabels: Record<string, { label: string; variant: "success" | "warning" | "danger" | "info" }> = {
  DEPLOY_SUCCESS: { label: "Deploy OK", variant: "success" },
  DEPLOY_STARTED: { label: "Deploy", variant: "info" },
  DEPLOY_FAILED: { label: "Falha", variant: "danger" },
  ROLLBACK_SUCCESS: { label: "Rollback OK", variant: "success" },
  ROLLBACK_FAILED: { label: "Rollback", variant: "danger" },
  VERIFY_SUCCESS: { label: "Verificado", variant: "success" },
  VERIFY_FAILED: { label: "Verificação", variant: "danger" },
  SYNC_COMPLETED: { label: "Sync OK", variant: "success" },
  SYNC_FAILED: { label: "Sync", variant: "danger" },
  PROJECT_CREATED: { label: "Criado", variant: "info" },
  SERVER_ADDED: { label: "Servidor", variant: "info" },
  ALERT: { label: "Alerta", variant: "warning" },
  ERROR: { label: "Erro", variant: "danger" },
};

export function Timeline({ items }: TimelineProps) {
  if (items.length === 0) {
    return <p className="text-gray-400 text-sm text-center py-8">Nenhum evento registrado.</p>;
  }

  return (
    <div className="space-y-0">
      {items.map((item, idx) => {
        const info = typeLabels[item.type || ""] || { label: item.type || "Evento", variant: "neutral" as const };
        return (
          <div key={item.id} className="flex gap-4">
            <div className="flex flex-col items-center">
              <div className={`w-2.5 h-2.5 rounded-full mt-1.5 ${
                info.variant === "success" ? "bg-green-500" :
                info.variant === "danger" ? "bg-red-500" :
                info.variant === "warning" ? "bg-amber-500" : "bg-blue-500"
              }`} />
              {idx < items.length - 1 && <div className="w-px flex-1 bg-gray-200 my-1" />}
            </div>
            <div className="flex-1 pb-4">
              <div className="flex items-center gap-2">
                <p className="text-sm font-medium text-gray-800">{item.title}</p>
                <Badge variant={info.variant}>{info.label}</Badge>
              </div>
              {item.description && (
                <p className="text-xs text-gray-500 mt-0.5">{item.description}</p>
              )}
              <p className="text-xs text-gray-400 mt-0.5">{item.timestamp}</p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
