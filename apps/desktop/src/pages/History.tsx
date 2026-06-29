// SPDX-License-Identifier: MIT
import { useEffect, useState, useCallback } from "react";
import { getHistory, getEvents } from "../api/client";
import { Timeline } from "../components/ui/Timeline";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Badge } from "../components/ui/Badge";

const eventTypes = [
  { value: "", label: "Todos" },
  { value: "DEPLOY_SUCCESS", label: "Deploy OK" },
  { value: "DEPLOY_FAILED", label: "Deploy Falha" },
  { value: "ROLLBACK_SUCCESS", label: "Rollback OK" },
  { value: "ROLLBACK_FAILED", label: "Rollback Falha" },
  { value: "VERIFY_SUCCESS", label: "Verificação OK" },
  { value: "VERIFY_FAILED", label: "Verificação Falha" },
  { value: "SYNC_COMPLETED", label: "Sync OK" },
  { value: "SYNC_FAILED", label: "Sync Falha" },
  { value: "PROJECT_CREATED", label: "Projeto" },
  { value: "SERVER_ADDED", label: "Servidor" },
  { value: "ALERT", label: "Alerta" },
];

export default function History() {
  const [events, setEvents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filterType, setFilterType] = useState("");
  const [limit, setLimit] = useState(50);

  const load = useCallback(() => {
    setLoading(true);
    setError("");
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    const fetch = filterType ? getEvents(wsPath, { type: filterType, limit }) : getHistory(wsPath, { limit });
    fetch
      .then((data) => setEvents(Array.isArray(data) ? data : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [filterType, limit]);

  useEffect(() => { load() }, [load]);

  const timelineItems = events.map((e: any) => ({
    id: e.id,
    title: e.description || e.event_type,
    description: e.metadata || "",
    timestamp: e.created_at,
    type: e.event_type,
  }));

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Histórico</h2>
        <div className="flex items-center gap-2">
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="px-3 py-1.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-eniac-500"
          >
            {eventTypes.map((t) => (
              <option key={t.value} value={t.value}>{t.label}</option>
            ))}
          </select>
          <Badge variant="info">{events.length} eventos</Badge>
        </div>
      </div>

      {error && <ErrorState message={error} onRetry={load} />}
      {loading ? <Loader /> : (
        <>
          <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm">
            <Timeline items={timelineItems} />
          </div>
          {events.length >= limit && (
            <div className="text-center mt-4">
              <button
                onClick={() => setLimit((l) => l + 50)}
                className="px-4 py-2 text-sm text-eniac-600 hover:text-eniac-800 font-medium"
              >
                Carregar mais
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
