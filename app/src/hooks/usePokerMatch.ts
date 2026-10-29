import { useCallback, useEffect, useRef, useState } from "react";
import { getPlayerConnectionID } from "../lib/playerIdentity";
import { getSocketRoot } from "../lib/socketEndpoints";
import { useAuthStore } from "../stores/useAuthStore";
import type { ArenaChatMessage } from "./useOnlineMatch";

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
  cards?: PokerCard[];
  cardCount: number;
}
export interface PokerState {
  roomId: string;
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
  chat: ArenaChatMessage[];
  spectators: number;
  ready: number[];
  inviteCode?: string;
  actionDeadline?: string;
}
export interface PokerOptions {
  name: string;
  mode?: PokerState["mode"];
  inviteCode?: string;
  role?: "player" | "spectator";
  roomId?: string;
}

function pokerSocketURL(playerID: string, options: PokerOptions) {
  const root = getSocketRoot("poker");
  const params = new URLSearchParams({
    playerId: options.role === "spectator" ? `${playerID}:spectator` : playerID,
    name: options.name,
  });
  if (options.mode) params.set("mode", options.mode);
  if (options.inviteCode) params.set("inviteCode", options.inviteCode);
  if (options.role) params.set("role", options.role);
  if (options.roomId) params.set("roomId", options.roomId);
  return `${root}?${params}`;
}

export function usePokerMatch() {
  const socketRef = useRef<WebSocket | null>(null);
  const optionsRef = useRef<PokerOptions | null>(null);
  const reconnectRef = useRef<number | null>(null);
  const intentional = useRef(false);
  const reconnectAction = useRef<(options: PokerOptions) => void>(
    () => undefined,
  );
  const [playerID] = useState(getPlayerConnectionID);
  const [name, setName] = useState(
    () => localStorage.getItem("boredgamz-player-name") ?? "Guest",
  );
  const accountName = useAuthStore((auth) => auth.user?.username);
  const [state, setState] = useState<PokerState | null>(null);
  const [phase, setPhase] = useState<
    "idle" | "connecting" | "waiting" | "playing"
  >("idle");
  const [error, setError] = useState("");

  const close = useCallback(() => {
    intentional.current = true;
    optionsRef.current = null;
    if (reconnectRef.current !== null) clearTimeout(reconnectRef.current);
    const socket = socketRef.current;
    if (socket?.readyState === WebSocket.OPEN)
      socket.send(JSON.stringify({ type: "leave" }));
    socket?.close(1000, "leaving poker table");
    socketRef.current = null;
    setState(null);
    setPhase("idle");
    setError("");
  }, []);

  const connect = useCallback(
    (input: PokerOptions) => {
      const options = { ...input, name: (accountName ?? input.name.trim()) || "Guest" };
      intentional.current = false;
      optionsRef.current = options;
      setPhase("connecting");
      setError("");
      localStorage.setItem("boredgamz-player-name", options.name);
      setName(options.name);
      const socket = new WebSocket(pokerSocketURL(playerID, options));
      const previous = socketRef.current;
      socketRef.current = socket;
      previous?.close();
      socket.onmessage = (event) => {
        if (socketRef.current !== socket) return;
        try {
          const message = JSON.parse(String(event.data)) as {
            type: string;
            state?: PokerState;
            message?: string;
          };
          if (message.type === "waiting") setPhase("waiting");
          if (message.type === "state" && message.state) {
            setState(message.state);
            setPhase("playing");
            setError("");
          }
          if (message.type === "error")
            setError(message.message ?? "That poker action was not accepted.");
        } catch {
          setError("The table sent an unreadable update.");
        }
      };
      socket.onerror = () => {
        if (socketRef.current === socket)
          setError("The poker room is temporarily unavailable.");
      };
      socket.onclose = () => {
        if (
          socketRef.current !== socket ||
          intentional.current ||
          !optionsRef.current
        )
          return;
        setError("Connection lost. Rejoining your seat…");
        reconnectRef.current = window.setTimeout(() => {
          if (optionsRef.current) reconnectAction.current(optionsRef.current);
        }, 1200);
      };
    },
    [accountName, playerID],
  );
  reconnectAction.current = connect;

  const send = useCallback((payload: object) => {
    setError("");
    if (socketRef.current?.readyState !== WebSocket.OPEN) {
      setError("The table is reconnecting.");
      return;
    }
    socketRef.current.send(JSON.stringify(payload));
  }, []);
  const action = useCallback(
    (actionName: string, amount = 0) =>
      send({ type: "action", action: actionName, amount }),
    [send],
  );
  const discard = useCallback(
    (positions: number[]) => send({ type: "discard", positions }),
    [send],
  );
  const chat = useCallback(
    (content: string) => send({ type: "chat", content }),
    [send],
  );
  const nextHand = useCallback(() => send({ type: "next-hand" }), [send]);
  useEffect(() => close, [close]);
  return {
    playerID,
    name: accountName ?? name,
    state,
    phase,
    error,
    connect,
    close,
    action,
    discard,
    chat,
    nextHand,
  };
}
