// SPDX-License-Identifier: MIT
import { check, type Update } from "@tauri-apps/plugin-updater";
import { useCallback, useRef, useState } from "react";

export type UpdateState =
  | { status: "idle" }
  | { status: "checking" }
  | { status: "available"; update: Update; version: string; body?: string }
  | { status: "downloading"; progress: number }
  | { status: "downloaded"; update: Update }
  | { status: "error"; message: string };

export function useUpdateCheck() {
  const [state, setState] = useState<UpdateState>({ status: "idle" });
  const updateRef = useRef<Update | null>(null);

  function dismiss() {
    setState({ status: "idle" });
  }

  const checkForUpdates = useCallback(async () => {
    setState({ status: "checking" });
    try {
      const update = await check();
      if (update) {
        updateRef.current = update;
        setState({
          status: "available",
          update,
          version: update.version,
          body: update.body,
        });
      } else {
        setState({ status: "idle" });
      }
    } catch (e: any) {
      setState({ status: "error", message: e?.message ?? "Falha ao verificar atualizacao" });
    }
  }, []);

  function startDownload() {
    const u = updateRef.current;
    if (!u || state.status !== "available") return;
    setState({ status: "downloading", progress: 0 });
    u.downloadAndInstall((event) => {
      if (event.event === "Progress") {
        setState((prev) =>
          prev.status === "downloading"
            ? { ...prev, progress: event.data.chunkLength }
            : prev
        );
      }
    }).catch((e: any) => {
      setState({ status: "error", message: e?.message ?? "Falha ao baixar atualizacao" });
    });
  }

  return { state, check: checkForUpdates, dismiss, startDownload };
}
