// SPDX-License-Identifier: MIT
import { NavLink } from "react-router-dom";

const links = [
  { to: "/dashboard", label: "Dashboard", icon: "◉" },
  { to: "/bootstrap", label: "Configurar", icon: "🚀" },
  { to: "/projects", label: "Projetos", icon: "📁" },
  { to: "/servers", label: "Servidores", icon: "🖥" },
  { to: "/github/login", label: "GitHub", icon: "🐙" },
  { to: "/deploy", label: "Deploy", icon: "🚀" },
  { to: "/rollback", label: "Rollback", icon: "⏪" },
  { to: "/sync", label: "Sync", icon: "🔄" },
  { to: "/diff", label: "Diff", icon: "📊" },
  { to: "/history", label: "Histórico", icon: "📋" },
  { to: "/health", label: "Saúde", icon: "❤️" },
];

export default function Sidebar() {
  return (
    <aside className="w-60 bg-eniac-900 text-white flex flex-col">
      <div className="p-4 border-b border-eniac-700">
        <h1 className="text-lg font-bold tracking-tight">FileENIAC</h1>
        <p className="text-xs text-eniac-300 mt-0.5">Desktop</p>
      </div>
      <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                isActive
                  ? "bg-eniac-700 text-white font-medium"
                  : "text-eniac-200 hover:bg-eniac-800 hover:text-white"
              }`
            }
          >
            <span className="text-lg">{link.icon}</span>
            {link.label}
          </NavLink>
        ))}
      </nav>
      <div className="p-4 border-t border-eniac-700 text-xs text-eniac-400">
        v0.3.0
      </div>
    </aside>
  );
}
