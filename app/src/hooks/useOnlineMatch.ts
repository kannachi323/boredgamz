import { useCallback, useEffect, useRef, useState } from "react";
import type { GridKind, GridMark } from "../lib/gridGames";
import { getPlayerConnectionID } from "../lib/playerIdentity";
import { getSocketRoot } from "../lib/socketEndpoints";
import { useAuthStore } from "../stores/useAuthStore";

export interface ArenaPlayer {
  id: string;
  name: string;
  mark: Exclude<GridMark, "">;
  connected: boolean;
}

export interface ArenaChatMessage {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  createdAt: string;
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
  chat: ArenaChatMessage[];
  spectators: number;
  turnDeadline?: string;
  chess?: ArenaChessState;
  clocks?: Record<Exclude<GridMark, "">, number>;
  drawOfferedBy?: Exclude<GridMark, "">;
  gridMoves: ArenaGridMove[];
}

export interface MatchOptions {
  name: string;
  mode?: string;
  bestOf?: number;
  timeControl?: string;
  inviteCode?: string;
  role?: "player" | "spectator";
  roomId?: string;
}

export type MatchPhase =
  "idle" | "connecting" | "waiting" | "playing" | "round-finished" | "finished";

function socketURL(
  kind: OnlineGameKind,
  playerID: string,
  options: MatchOptions,
) {
  const base = getSocketRoot("arena");
  const params = new URLSearchParams({
    playerId: options.role === "spectator" ? `${playerID}:spectator` : playerID,
    name: options.name,
  });
  if (options.mode) params.set("mode", options.mode);
  if (options.bestOf) params.set("bestOf", String(options.bestOf));
  if (options.timeControl) params.set("timeControl", options.timeControl);
  if (options.inviteCode) params.set("inviteCode", options.inviteCode);
  if (options.role === "spectator") params.set("role", "spectator");
  if (options.roomId) params.set("roomId", options.roomId);
  return `${base.replace(/\/$/, "")}/${kind}?${params}`;
}

export function useOnlineMatch<TCell extends string = GridMark>(
  kind: OnlineGameKind,
) {
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectRef = useRef<number | null>(null);
  const reconnectActionRef = useRef<(options: MatchOptions) => void>(
    () => undefined,
  );
  const connectionRef = useRef<MatchOptions | null>(null);
  const intentionalClose = useRef(false);
  const [state, setState] = useState<ArenaState<TCell> | null>(null);
  const [phase, setPhase] = useState<MatchPhase>("idle");
  const [error, setError] = useState("");
  const [name, setName] = useState(
    () => window.localStorage.getItem("boredgamz-player-name") ?? "Guest",
  );
  const accountName = useAuthStore((auth) => auth.user?.username);
  const [playerID] = useState(getPlayerConnectionID);

  const close = useCallback(() => {
    intentionalClose.current = true;
    connectionRef.current = null;
    if (reconnectRef.current !== null)
      window.clearTimeout(reconnectRef.current);
    reconnectRef.current = null;
    const socket = socketRef.current;
    if (socket?.readyState === WebSocket.OPEN)
      socket.send(JSON.stringify({ type: "leave" }));
    socket?.close(1000, "leaving table");
    socketRef.current = null;
    setState(null);
    setPhase("idle");
    setError("");
  }, []);

  const connect = useCallback(
    (input: string | MatchOptions = name) => {
      const options: MatchOptions =
        typeof input === "string" ? { name: input } : input;
      const cleanName = (accountName ?? options.name.trim()) || "Guest";
      const nextOptions = { ...options, name: cleanName };
      intentionalClose.current = false;
      connectionRef.current = nextOptions;
      window.localStorage.setItem("boredgamz-player-name", cleanName);
      setName(cleanName);
      setError("");
      setPhase("connecting");

      const previous = socketRef.current;
      const socket = new WebSocket(socketURL(kind, playerID, nextOptions));
      socketRef.current = socket;
      previous?.close();

      socket.onmessage = (event) => {
        if (socketRef.current !== socket) return;
        try {
          const message = JSON.parse(String(event.data)) as {
            type: string;
            state?: ArenaState<TCell>;
            message?: string;
          };
          if (message.type === "waiting") setPhase("waiting");
          if (message.type === "state" && message.state) {
            setState(message.state);
            setPhase(message.state.status);
            setError("");
          }
          if (message.type === "error")
            setError(message.message ?? "That action was not accepted.");
        } catch {
          setError(
            "The table sent an unreadable update. Reconnecting may resolve it.",
          );
        }
      };
      socket.onerror = () => {
        if (socketRef.current === socket)
          setError("The live table is temporarily unavailable.");
      };
      socket.onclose = () => {
        if (
          socketRef.current !== socket ||
          intentionalClose.current ||
          connectionRef.current === null
        )
          return;
        setError("Connection lost. Rejoining the table…");
        reconnectRef.current = window.setTimeout(() => {
          const latest = connectionRef.current;
          if (latest) reconnectActionRef.current(latest);
        }, 1200);
      };
    },
    [accountName, kind, name, playerID],
  );

  reconnectActionRef.current = connect;

  const send = useCallback((message: object) => {
    setError("");
    if (socketRef.current?.readyState !== WebSocket.OPEN) {
      setError("The table is reconnecting. Try again in a moment.");
      return;
    }
    socketRef.current.send(JSON.stringify(message));
  }, []);

  const move = useCallback(
    (position: number) => send({ type: "move", position }),
    [send],
  );
  const chat = useCallback(
    (content: string) => send({ type: "chat", content }),
    [send],
  );
  const resign = useCallback(() => send({ type: "resign" }), [send]);
  const rematch = useCallback(() => send({ type: "rematch" }), [send]);
  const chessMove = useCallback(
    (from: string, to: string, promotion = "") =>
      send({ type: "move", from, to, promotion }),
    [send],
  );
  const offerDraw = useCallback(() => send({ type: "offer-draw" }), [send]);
  const acceptDraw = useCallback(() => send({ type: "accept-draw" }), [send]);
  const declineDraw = useCallback(() => send({ type: "decline-draw" }), [send]);

  useEffect(() => close, [close]);
  return {
    state,
    phase,
    error,
    name: accountName ?? name,
    playerID,
    connect,
    close,
    move,
    chessMove,
    chat,
    resign,
    rematch,
    offerDraw,
    acceptDraw,
    declineDraw,
  };
}
