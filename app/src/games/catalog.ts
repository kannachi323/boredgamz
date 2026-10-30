export type GameSlug = "tictactoe" | "connect-four" | "gomoku" | "chess" | "poker" | "mahjong";
export type GameService = "arena" | "poker" | "mahjong";

export type ModeGlyph =
  | { kind: "grid"; rows: number; columns: number; strike?: boolean }
  | { kind: "clock"; minutes: number }
  | { kind: "cards"; count: number }
  | { kind: "tiles"; layers: number };

export interface GameMode {
  id: string;
  name: string;
  detail: string;
  glyph: ModeGlyph;
  /** Chess: the time control that belongs to this mode. */
  timeControl?: string;
}

export interface Choice<T extends string | number = string> {
  value: T;
  label: string;
}

export interface BoardTheme {
  id: string;
  label: string;
  light: string;
  dark: string;
}

export interface GameDefinition {
  slug: GameSlug;
  name: string;
  minutes: string;
  service: GameService;
  /** First digit of private codes, so a code alone finds its game. */
  codePrefix: string;
  modes: GameMode[];
  series?: Choice<number>[];
  clocks?: Choice[];
  defaultClock?: string;
  themes?: BoardTheme[];
}

const series = (...values: number[]): Choice<number>[] =>
  values.map((value) => ({ value, label: value === 1 ? "Single" : `Best of ${value}` }));

export const games: GameDefinition[] = [
  {
    slug: "tictactoe",
    name: "Tic-Tac-Toe",
    minutes: "2 min",
    service: "arena",
    codePrefix: "1",
    modes: [
      { id: "classic", name: "Classic", detail: "3 × 3", glyph: { kind: "grid", rows: 3, columns: 3 } },
      { id: "four-square", name: "Four Square", detail: "4 × 4, three in a row", glyph: { kind: "grid", rows: 4, columns: 4 } },
      { id: "misere", name: "Misère", detail: "Three in a row loses", glyph: { kind: "grid", rows: 3, columns: 3, strike: true } },
    ],
    series: series(1, 3, 5, 7),
    clocks: [
      { value: "untimed", label: "Off" },
      { value: "15s", label: "15s" },
      { value: "30s", label: "30s" },
      { value: "60s", label: "60s" },
    ],
    defaultClock: "30s",
  },
  {
    slug: "connect-four",
    name: "Connect Four",
    minutes: "5 min",
    service: "arena",
    codePrefix: "2",
    modes: [
      { id: "classic", name: "Classic", detail: "6 × 7, four to win", glyph: { kind: "grid", rows: 6, columns: 7 } },
      { id: "quick", name: "Quick", detail: "5 × 6, four to win", glyph: { kind: "grid", rows: 5, columns: 6 } },
      { id: "connect-five", name: "Connect Five", detail: "7 × 8, five to win", glyph: { kind: "grid", rows: 7, columns: 8 } },
    ],
    series: series(1, 3, 5, 7),
    clocks: [
      { value: "untimed", label: "Off" },
      { value: "15s", label: "15s" },
      { value: "30s", label: "30s" },
      { value: "60s", label: "60s" },
    ],
    defaultClock: "30s",
  },
  {
    slug: "gomoku",
    name: "Gomoku",
    minutes: "10 min",
    service: "arena",
    codePrefix: "3",
    modes: [
      { id: "freestyle", name: "Freestyle", detail: "15 × 15, five or more", glyph: { kind: "grid", rows: 7, columns: 7 } },
      { id: "pro", name: "Exact Five", detail: "Overlines don't count", glyph: { kind: "grid", rows: 7, columns: 7, strike: true } },
      { id: "compact", name: "Compact", detail: "13 × 13", glyph: { kind: "grid", rows: 5, columns: 5 } },
    ],
    series: series(1, 3, 5),
    clocks: [
      { value: "untimed", label: "Off" },
      { value: "30s", label: "30s" },
      { value: "60s", label: "60s" },
    ],
    defaultClock: "untimed",
    themes: [
      { id: "maple", label: "Maple", light: "#e8c995", dark: "#b8894f" },
      { id: "walnut", label: "Walnut", light: "#b5845a", dark: "#6e4b30" },
      { id: "slate", label: "Slate", light: "#8f9497", dark: "#4b5054" },
    ],
  },
  {
    slug: "chess",
    name: "Chess",
    minutes: "3–25 min",
    service: "arena",
    codePrefix: "4",
    modes: [
      { id: "blitz", name: "Blitz", detail: "3 min + 2s", glyph: { kind: "clock", minutes: 3 }, timeControl: "3m+2s" },
      { id: "rapid", name: "Rapid", detail: "10 min", glyph: { kind: "clock", minutes: 10 }, timeControl: "10m" },
      { id: "standard", name: "Classical", detail: "15 min + 10s", glyph: { kind: "clock", minutes: 15 }, timeControl: "15m+10s" },
    ],
    themes: [
      { id: "walnut", label: "Walnut", light: "#eadbc3", dark: "#a5774f" },
      { id: "sage", label: "Sage", light: "#e7ead9", dark: "#7d9466" },
      { id: "slate", label: "Slate", light: "#d9dde2", dark: "#6c7886" },
      { id: "paper", label: "Paper", light: "#f3f1ea", dark: "#bdb6a6" },
    ],
  },
  {
    slug: "poker",
    name: "Poker",
    minutes: "10 min",
    service: "poker",
    codePrefix: "5",
    modes: [
      { id: "holdem", name: "Hold'em", detail: "Heads-up, no limit", glyph: { kind: "cards", count: 2 } },
      { id: "five-card-draw", name: "Five-Card Draw", detail: "One draw, two bets", glyph: { kind: "cards", count: 5 } },
    ],
  },
  {
    slug: "mahjong",
    name: "Mahjong",
    minutes: "8–12 min",
    service: "mahjong",
    codePrefix: "6",
    modes: [
      { id: "garden", name: "Shared Garden", detail: "Clear 144 tiles together", glyph: { kind: "tiles", layers: 3 } },
      { id: "duel", name: "Mirror Duel", detail: "Race the same 72 tiles", glyph: { kind: "tiles", layers: 2 } },
      { id: "zen", name: "Zen", detail: "Together, no clock", glyph: { kind: "tiles", layers: 1 } },
    ],
  },
];

export function getGame(slug: string | undefined) {
  return games.find((game) => game.slug === slug);
}

/** Six digits: the game's digit, then five random ones (easy to read aloud and type on a phone). */
export function makeInviteCode(game: GameDefinition) {
  const bytes = new Uint32Array(5);
  crypto.getRandomValues(bytes);
  return game.codePrefix + Array.from(bytes, (value) => value % 10).join("");
}

export function normalizeCode(value: string) {
  return value.replace(/\D/g, "").slice(0, 6);
}

/** Finds which game a private code belongs to from its prefix. */
export function gameForCode(code: string) {
  const clean = normalizeCode(code);
  if (clean.length !== 6) return undefined;
  return games.find((game) => game.codePrefix === clean[0]);
}
