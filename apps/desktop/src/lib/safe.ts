export function asArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : [];
}

export function hasItems<T>(value: T[] | null | undefined): value is T[] {
  return Array.isArray(value) && value.length > 0;
}

export function safeString(value: string | null | undefined, fallback = ""): string {
  return typeof value === "string" ? value : fallback;
}

export function safeNumber(value: number | null | undefined, fallback = 0): number {
  return typeof value === "number" && !Number.isNaN(value) ? value : fallback;
}
