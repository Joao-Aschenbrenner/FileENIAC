import { useTheme } from "../../context/ThemeContext";
import { Sun, Moon, Monitor } from "lucide-react";

export function ThemeToggle() {
  const { mode, setMode } = useTheme();

  const options = [
    { value: "light" as const, label: "Claro", icon: Sun },
    { value: "dark" as const, label: "Escuro", icon: Moon },
    { value: "system" as const, label: "Sistema", icon: Monitor },
  ];

  return (
    <div className="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
      {options.map((opt) => (
        <button
          key={opt.value}
          onClick={() => setMode(opt.value)}
          className={`flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-sm font-medium transition-all duration-150 ${
            mode === opt.value
              ? "bg-white text-eniac-700 shadow-sm"
              : "text-gray-600 hover:text-gray-900 hover:bg-gray-200"
          }`}
          aria-pressed={mode === opt.value}
        >
          <opt.icon className="h-4 w-4" aria-hidden="true" />
          <span>{opt.label}</span>
        </button>
      ))}
    </div>
  );
}