export type ThemePreference = "system" | "light" | "dark";
const key = "boredgamz-theme";

export function readTheme(): ThemePreference {
  try {
    const value = localStorage.getItem(key);
    return value === "light" || value === "dark" ? value : "system";
  } catch {
    return "system";
  }
}

export function applyTheme(theme: ThemePreference) {
  const root = document.documentElement;
  if (theme === "system") root.removeAttribute("data-theme");
  else root.setAttribute("data-theme", theme);
  try {
    if (theme === "system") localStorage.removeItem(key);
    else localStorage.setItem(key, theme);
  } catch {
    // Storage can be unavailable (private mode); the choice just won't persist.
  }
}

export function effectiveTheme(theme: ThemePreference): "light" | "dark" {
  if (theme !== "system") return theme;
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}
