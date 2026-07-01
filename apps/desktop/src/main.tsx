// SPDX-License-Identifier: MIT
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import { initApiClient, heartbeat } from "./api/client";
import "./index.css";
import "./styles/themes.css";

let heartbeatInterval: ReturnType<typeof setInterval> | null = null;

async function boot() {
  await initApiClient();
  heartbeatInterval = setInterval(heartbeat, 10000);
}

boot();

window.addEventListener("beforeunload", () => {
  if (heartbeatInterval) clearInterval(heartbeatInterval);
});

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);
