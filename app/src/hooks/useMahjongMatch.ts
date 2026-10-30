import { useCallback } from "react";
import { useGameSocket, type ChatMessage } from "./useGameSocket";

export interface MahjongTile {
  id: number;
  face: string;
  group: string;
  x: number;
  y: number;
  z: number;
}

export interface MahjongSeat {
  id: string;
  name: string;
  pairs: number;
  remaining: number;
  connected: boolean;
  left?: boolean;
}

export interface MahjongState {
  roomId: string;
  viewerId?: string;
  mode: "garden" | "duel" | "zen";
  status: "playing" | "finished";
  tiles: MahjongTile[];
  removed: boolean[];
  removedBy: number[];
  free: number[];
  seats: [MahjongSeat, MahjongSeat];
  viewerSeat: number;
  winner: number;
  message: string;
  hintsRemaining: number;
  hint: number[];
  chat: ChatMessage[];
  spectators: number;
  ready: number[];
  inviteCode?: string;
  startedAt: string;
  deadline?: string;
  reshuffles: number;
}

const mahjongOver = (state: MahjongState) => state.status === "finished";

export function useMahjongMatch() {
  const socket = useGameSocket<MahjongState>("mahjong", "", mahjongOver);
  const { send } = socket;
  return {
    ...socket,
    pair: useCallback((first: number, second: number) => send({ type: "pair", first, second }), [send]),
    hint: useCallback(() => send({ type: "hint" }), [send]),
    rematch: useCallback(() => send({ type: "rematch" }), [send]),
  };
}
