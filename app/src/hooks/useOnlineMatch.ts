import { useCallback } from "react";
import { useGameSocket, type ChatMessage } from "./useGameSocket";

export type GridMark = "" | "one" | "two";
export type GridKind = "tictactoe" | "connect-four" | "gomoku";

export interface ArenaPlayer {
  id: string;
  name: string;
  mark: Exclude<GridMark, "">;
  connected: boolean;
  left?: boolean;
}

export interface ArenaSettings {
  mode: string;
  bestOf: number;
  timeControl: string;
  visibility: "public" | "private";
  inviteCode?: string;
}

export interface ArenaChessMove {
  from: string;
  to: string;
  promotion?: string;
  san: string;
}

export interface ArenaGridMove {
  position: number;
  mark: Exclude<GridMark, "">;
}

export interface ArenaChessState {
  fen: string;
  moves: ArenaChessMove[];
  lastFrom?: string;
  lastTo?: string;
  outcome?: string;
  method?: string;
}

export type OnlineGameKind = GridKind | "chess";

export interface ArenaState<TCell extends string = GridMark> {
  roomId: string;
  viewerId?: string;
  kind: OnlineGameKind;
  status: "playing" | "round-finished" | "finished";
  rows: number;
  columns: number;
  board: TCell[];
  turn: Exclude<GridMark, "">;
  winner?: Exclude<GridMark, "">;
  seriesWinner?: Exclude<GridMark, "">;
  players: [ArenaPlayer | null, ArenaPlayer | null];
  viewerMark: Exclude<GridMark, ""> | "spectator";
  lastMove: number;
  settings: ArenaSettings;
  round: number;
  score: Record<Exclude<GridMark, "">, number>;
  rematchReady: string[];
  chat: ChatMessage[];
  spectators: number;
  turnDeadline?: string;
  chess?: ArenaChessState;
  clocks?: Record<Exclude<GridMark, "">, number>;
  drawOfferedBy?: Exclude<GridMark, "">;
  gridMoves: ArenaGridMove[];
}

const arenaOver = (state: ArenaState<string>) => state.status !== "playing";

export function useOnlineMatch<TCell extends string = GridMark>(kind: OnlineGameKind) {
  const socket = useGameSocket<ArenaState<TCell>>("arena", `/${kind}`, arenaOver);
  const { send } = socket;
  return {
    ...socket,
    move: useCallback((position: number) => send({ type: "move", position }), [send]),
    chessMove: useCallback((from: string, to: string, promotion = "") => send({ type: "move", from, to, promotion }), [send]),
    resign: useCallback(() => send({ type: "resign" }), [send]),
    rematch: useCallback(() => send({ type: "rematch" }), [send]),
    offerDraw: useCallback(() => send({ type: "offer-draw" }), [send]),
    acceptDraw: useCallback(() => send({ type: "accept-draw" }), [send]),
    declineDraw: useCallback(() => send({ type: "decline-draw" }), [send]),
  };
}
