// SPDX-License-Identifier: MIT
import { NavLink } from "react-router-dom";
import {
  IconDashboard,
  IconSettings,
  IconFolder,
  IconServer,
  IconDeploy,
  IconRollback,
  IconSync,
  IconDiff,
  IconHistory,
  IconHealth,
} from "./Icons";

const links = [
  { to: "/dashboard", label: "Dashboard", Icon: IconDashboard },
  { to: "/configurar", label: "Configurar", Icon: IconSettings },
  { to: "/projects", label: "Projetos", Icon: IconFolder },
  { to: "/servers", label: "Servidores", Icon: IconServer },
  { to: "/deploy", label: "Deploy", Icon: IconDeploy },
  { to: "/rollback", label: "Rollback", Icon: IconRollback },
  { to: "/sync", label: "Sync", Icon: IconSync },
  { to: "/diff", label: "Diff", Icon: IconDiff },
  { to: "/history", label: "Historico", Icon: IconHistory },
  { to: "/health", label: "Saude", Icon: IconHealth },
];

export default function Sidebar() {
  return (
    <aside className="w-60 bg-gray-900 text-white flex flex-col">
      <div className="p-4 border-b border-gray-700">
        <h1 className="text-lg font-bold tracking-tight text-white">FileENIAC</h1>
        <p className="text-xs text-gray-400 mt-0.5">Desktop</p>
      </div>
      <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
        {links.map(({ to, label, Icon }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                isActive
                  ? "bg-gray-700 text-white font-medium"
                  : "text-gray-400 hover:bg-gray-800 hover:text-white"
              }`
            }
          >
            <Icon className="flex-shrink-0 opacity-80" />
            {label}
          </NavLink>
        ))}
      </nav>
      <div className="p-4 border-t border-gray-700 text-xs text-gray-500">
        v0.2.0
      </div>
    </aside>
  );
}