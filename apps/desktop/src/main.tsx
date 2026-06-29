// SPDX-License-Identifier: MIT
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import { initApiClient, heartbeat } from "./api/client";
import "./index.css";

let heartbeatInterval: ReturnType<typeof setInterval> | null = null;
initApiClient().then(() => {
  heartbeatInterval = setInterval(heartbeat, 10000);
});

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
