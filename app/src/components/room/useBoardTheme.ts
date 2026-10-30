import { useCallback, useState } from "react";
import type { GameDefinition } from "../../games/catalog";

/** Personal board theme per game, remembered on this device. */
export function useBoardTheme(game: GameDefinition) {
  const key = `boredgamz-board-${game.slug}`;
  const fallback = game.themes?.[0]?.id ?? "";
  const [theme, setThemeState] = useState(() => {
    try {
      const stored = localStorage.getItem(key);
      return game.themes?.some((item) => item.id === stored) ? stored! : fallback;
    } catch {
      return fallback;
    }
  });
  const setTheme = useCallback(
    (next: string) => {
      setThemeState(next);
      try {
        localStorage.setItem(key, next);
      } catch {
        // Not persisted in private mode; still applies for this visit.
      }
    },
    [key],
  );
  return [theme, setTheme] as const;
}
