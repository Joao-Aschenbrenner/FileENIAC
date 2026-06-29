// SPDX-License-Identifier: MIT
import { useEffect, useState } from "react";
import { checkHealth } from "../api/client";

export default function Header() {
  const [backendOnline, setBackendOnline] = useState(false);

  useEffect(() => {
    const check = async () => {
      const ok = await checkHealth();
      setBackendOnline(ok);
    };
    check();
    const interval = setInterval(check, 10000);
    return () => clearInterval(interval);
  }, []);

  return (
    <header className="flex items-center justify-between px-6 py-3 bg-white border-b border-gray-200">
      <h2 className="text-lg font-semibold text-gray-800">Workspace</h2>
      <div className="flex items-center gap-2 text-sm">
        <span
          className={`w-2 h-2 rounded-full ${
            backendOnline ? "bg-green-500" : "bg-red-500"
          }`}
        />
        <span className="text-gray-500">
          {backendOnline ? "Backend online" : "Backend offline"}
        </span>
      </div>
    </header>
  );
}
