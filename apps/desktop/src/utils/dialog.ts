// SPDX-License-Identifier: MIT
import { open } from "@tauri-apps/plugin-dialog";

export async function pickFolder(title = "Selecionar pasta"): Promise<string | null> {
  try {
    const selected = await open({ directory: true, multiple: false, title });
    return selected || null;
  } catch {
    return null;
  }
}

export async function pickFile(filters?: { name: string; extensions: string[] }[]): Promise<string | null> {
  try {
    const selected = await open({ multiple: false, filters });
    return selected || null;
  } catch {
    return null;
  }
}

export async function pickFiles(filters?: { name: string; extensions: string[] }[]): Promise<string[]> {
  try {
    const selected = await open({ multiple: true, filters });
    return Array.isArray(selected) ? selected : selected ? [selected] : [];
  } catch {
    return [];
  }
}

export async function saveFile(defaultPath?: string, filters?: { name: string; extensions: string[] }[]): Promise<string | null> {
  try {
    const selected = await open({ directory: false, multiple: false, defaultPath, filters, title: "Salvar arquivo" });
    return selected || null;
  } catch {
    return null;
  }
}