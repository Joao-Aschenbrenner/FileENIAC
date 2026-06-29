// SPDX-License-Identifier: MIT
import { useEffect, useState, useCallback } from "react";
import { getHealthCheck, checkHealth, listProjects } from "../api/client";
import { Card } from "../components/ui/Card";
import { Badge } from "../components/ui/Badge";
import { Loader } from "../components/ui/Loader";
import { ErrorState } from "../components/ui/ErrorState";
import { Table } from "../components/ui/Table";

export default function HealthCenter() {
  const [health, setHealth] = useState<any>(null);
  const [projects, setProjects] = useState<any[]>([]);
  const [apiOnline, setApiOnline] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const load = useCallback(() => {
    setLoading(true);
    setError("");
    const wsPath = localStorage.getItem("eniac_ws_path") || "";
    Promise.all([
      getHealthCheck(),
      listProjects(wsPath),
      checkHealth(),
    ])
      .then(([h, p, ok]) => {
        setHealth(h);
        setProjects(p);
        setApiOnline(ok);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load() }, [load]);

  if (loading) return <Loader />;
  if (error) return <ErrorState message={error} onRetry={load} />;

  const checks = [
    { label: "API Server", status: apiOnline, description: apiOnline ? "Respondendo na porta 8080" : "Sem resposta" },
    { label: "Workspace", status: !!health, description: health ? `${health.projects} projetos, ${health.servers} servidores` : "Não conectado" },
    { label: "Projetos", status: health && health.projects > 0, description: health && health.projects > 0 ? `${health.projects} cadastrados` : "Nenhum projeto" },
    { label: "Divergentes", status: health && health.divergent === 0, description: health && health.divergent > 0 ? `${health.divergent} divergentes` : "Nenhuma divergência" },
  ];

  const projectHealthColumns = [
    { key: "name", header: "Projeto" },
    { key: "divergence_status", header: "Status", render: (v: string) => (
      <Badge variant={v === "sincronizado" ? "success" : "warning"}>{v}</Badge>
    )},
    { key: "environment", header: "Ambiente" },
    { key: "branch", header: "Branch" },
  ];

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Health Center</h2>
        <Badge variant={health?.status === "healthy" ? "success" : "warning"}>
          {health?.status === "healthy" ? "Sistema Saudável" : "Degradado"}
        </Badge>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-6">
        {checks.map((c) => (
          <Card key={c.label}>
            <div className="flex items-center gap-3 mb-2">
              <div className={`w-3 h-3 rounded-full ${c.status ? "bg-green-500" : "bg-red-500"}`} />
              <p className="font-semibold text-gray-800">{c.label}</p>
            </div>
            <p className="text-sm text-gray-500">{c.description}</p>
          </Card>
        ))}
      </div>

      <Card title="Projetos">
        <Table columns={projectHealthColumns} data={projects} />
      </Card>

      {health?.last_events && health.last_events.length > 0 && (
        <Card title="Eventos Recentes" className="mt-4">
          <div className="space-y-2">
            {health.last_events.map((ev: any) => (
              <div key={ev.id} className="flex items-center justify-between text-sm py-1">
                <span className="text-gray-700">{ev.description || ev.event_type}</span>
                <span className="text-xs text-gray-400">{ev.created_at}</span>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
}
