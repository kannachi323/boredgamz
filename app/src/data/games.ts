import type { LucideIcon } from "lucide-react";
import { Blocks, Brain, CircleDot, Club, Crown, Grid3X3 } from "lucide-react";

export type GameSlug =
  "tictactoe" | "connect-four" | "gomoku" | "chess" | "poker" | "mahjong";

export interface GameDefinition {
  slug: GameSlug;
  name: string;
  kicker: string;
  description: string;
  players: string;
  duration: string;
  modes: string[];
  accent: string;
  icon: LucideIcon;
}

export const games: GameDefinition[] = [
  {
    slug: "tictactoe",
    name: "Tic-Tac-Toe",
    kicker: "Quick thinking",
    description:
      "A tiny live board with nowhere to hide. Find the fork before your rival does.",
    players: "2 online",
    duration: "2 min",
    modes: ["Classic", "Four Square", "Misère"],
    accent: "#ff7a5c",
    icon: Grid3X3,
  },
  {
    slug: "connect-four",
    name: "Connect Four",
    kicker: "Vertical tactics",
    description:
      "Drop a disc, shape a threat, and connect before the board fills up.",
    players: "2 online",
    duration: "5 min",
    modes: ["Classic", "Quick", "Connect Five"],
    accent: "#ffd447",
    icon: CircleDot,
  },
  {
    slug: "gomoku",
    name: "Gomoku",
    kicker: "Five in a row",
    description:
      "Build intersecting threats on a live wood board with exact server rules.",
    players: "2 online",
    duration: "10 min",
    modes: ["Freestyle", "Exact Five", "Compact"],
    accent: "#c8ff5a",
    icon: Blocks,
  },
  {
    slug: "chess",
    name: "Chess",
    kicker: "The royal game",
    description:
      "Server-verified chess with real clocks, notation, chat, and personal board themes.",
    players: "2 online",
    duration: "3–25 min",
    modes: ["Blitz", "Rapid", "Classical"],
    accent: "#a98bff",
    icon: Crown,
  },
  {
    slug: "poker",
    name: "Poker",
    kicker: "Read the table",
    description:
      "A private-card heads-up room with exact wagers and cryptographic dealing.",
    players: "2 online",
    duration: "10 min",
    modes: ["Hold’em", "Five-Card Draw"],
    accent: "#58d6c7",
    icon: Club,
  },
  {
    slug: "mahjong",
    name: "Mahjong Solitaire",
    kicker: "Quiet focus",
    description:
      "Clear a guaranteed-solvable live tile garden together, or race identical mirrored layouts.",
    players: "2 online",
    duration: "8–12 min",
    modes: ["Shared Garden", "Mirror Duel", "Zen"],
    accent: "#ff9fca",
    icon: Brain,
  },
];

export function getGame(slug: string) {
  return games.find((game) => game.slug === slug);
}
