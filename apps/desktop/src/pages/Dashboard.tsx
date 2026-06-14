import { useEffect, useRef, useState } from "react";
import { getHealthCheck } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";

export default function Dashboard() {
  const [data, setData] = useState<any>(null);
  const [error, setError] = useState("");
  const mounted = useRef(true);

  function load() {
    setError("");
    getHealthCheck()
      .then((d) => { if (mounted.current) setData(d); })
      .catch((e) => { if (mounted.current) setError(e.message); });
  }

  useEffect(() => { load(); return () => { mounted.current = false; }; }, []);

  if (error) return <ErrorState message={error} onRetry={load} />;
  if (!data) return <Loader text="Carregando dashboard..." />;

  const stats = [
    { label: "Projetos", value: data.projects, color: "bg-eniac-500" },
    { label: "Servidores", value: data.servers, color: "bg-emerald-500" },
    { label: "Divergentes", value: data.divergent, color: "bg-amber-500" },
    { label: "Eventos Recentes", value: data.last_events?.length || 0, color: "bg-violet-500" },
  ];

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-gray-800">Dashboard</h2>
          <p className="text-sm text-gray-500 mt-0.5">Visão geral do workspace</p>
        </div>
        <Badge variant={data.status === "healthy" ? "success" : "warning"}>
          {data.status === "healthy" ? "Saudável" : "Degradado"}
        </Badge>
      </div>

      <div className="grid grid-cols-4 gap-4 mb-6">
        {stats.map((s) => (
          <Card key={s.label}>
            <div className={`w-3 h-3 rounded-full ${s.color} mb-3`} />
            <p className="text-2xl font-bold text-gray-800">{s.value}</p>
            <p className="text-sm text-gray-500">{s.label}</p>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-2 gap-6">
        <Card title="Status Geral">
          <div className="space-y-3">
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Projetos cadastrados</span>
              <span className="font-medium">{data.projects}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Servidores configurados</span>
              <span className="font-medium">{data.servers}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Projetos divergentes</span>
              <span className="font-medium">{data.divergent}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Estado</span>
              <Badge variant={data.status === "healthy" ? "success" : "warning"}>
                {data.status}
              </Badge>
            </div>
          </div>
        </Card>

        <Card title="Últimos Eventos">
          {data.last_events?.length > 0 ? (
            <div className="space-y-2">
              {data.last_events.map((ev: any) => (
                <div key={ev.id} className="flex items-center justify-between text-sm">
                  <span className="text-gray-700 truncate max-w-[200px]">{ev.description || ev.event_type}</span>
                  <span className="text-xs text-gray-400">{ev.created_at}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-400">Nenhum evento recente.</p>
          )}
        </Card>
      </div>
    </div>
  );
}
