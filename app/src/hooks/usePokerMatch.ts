import { useCallback } from "react";
import { useGameSocket, type ChatMessage } from "./useGameSocket";

export interface PokerCard {
  rank: number;
  suit: "s" | "h" | "d" | "c";
  code: string;
}
export interface PokerSeat {
  id: string;
  name: string;
  stack: number;
  roundBet: number;
  totalBet: number;
  folded: boolean;
  allIn: boolean;
  connected: boolean;
  left?: boolean;
  cards?: PokerCard[];
  cardCount: number;
}
export interface PokerState {
  roomId: string;
  viewerId?: string;
  viewerSeat: number;
  mode: "holdem" | "five-card-draw";
  phase: string;
  seats: [PokerSeat, PokerSeat];
  community: PokerCard[];
  dealer: number;
  turn: number;
  pot: number;
  awardedPot: number;
  currentBet: number;
  minimumRaise: number;
  toCall: number;
  legalActions: string[];
  hand: number;
  winners: number[];
  winningHand?: string;
  message: string;
  chat: ChatMessage[];
  spectators: number;
  ready: number[];
  inviteCode?: string;
  actionDeadline?: string;
}
const pokerOver = (state: PokerState) => state.phase === "showdown" || state.phase === "hand-finished";

export function usePokerMatch() {
  const socket = useGameSocket<PokerState>("poker", "", pokerOver);
  const { send } = socket;
  return {
    ...socket,
    action: useCallback((action: string, amount = 0) => send({ type: "action", action, amount }), [send]),
    discard: useCallback((positions: number[]) => send({ type: "discard", positions }), [send]),
    nextHand: useCallback(() => send({ type: "next-hand" }), [send]),
  };
}
