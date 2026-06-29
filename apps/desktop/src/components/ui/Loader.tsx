// SPDX-License-Identifier: MIT
export function Loader({ text = "Carregando..." }: { text?: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-12">
      <div className="w-8 h-8 border-4 border-eniac-200 border-t-eniac-600 rounded-full animate-spin" />
      <p className="mt-3 text-sm text-gray-500">{text}</p>
    </div>
  );
}
