import { Moon, Sun } from "lucide-react";
import { useState } from "react";
import { applyTheme, effectiveTheme, readTheme } from "../lib/theme";

export function ThemeToggle() {
  const [theme, setTheme] = useState(readTheme);
  const current = effectiveTheme(theme);
  const next = current === "dark" ? "light" : "dark";
  return (
    <button
      className="btn btn--ghost btn--icon"
      type="button"
      aria-label={`Switch to ${next} mode`}
      title={`Switch to ${next} mode`}
      onClick={() => {
        applyTheme(next);
        setTheme(next);
      }}
    >
      {current === "dark" ? <Sun size={18} /> : <Moon size={18} />}
    </button>
  );
}
