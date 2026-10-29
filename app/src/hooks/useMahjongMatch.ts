import { useCallback, useEffect, useRef, useState } from "react";
import { getPlayerConnectionID } from "../lib/playerIdentity";
import { getSocketRoot } from "../lib/socketEndpoints";
import { useAuthStore } from "../stores/useAuthStore";
import type { ArenaChatMessage } from "./useOnlineMatch";

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
}

export interface MahjongState {
  roomId: string;
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
  chat: ArenaChatMessage[];
  spectators: number;
  ready: number[];
  inviteCode?: string;
  startedAt: string;
  deadline?: string;
  reshuffles: number;
}

export interface MahjongOptions {
  name: string;
  mode?: MahjongState["mode"];
  inviteCode?: string;
  role?: "player" | "spectator";
  roomId?: string;
}

function socketURL(playerID: string, options: MahjongOptions) {
  const root = getSocketRoot("mahjong");
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

export function useMahjongMatch() {
  const socketRef = useRef<WebSocket | null>(null);
  const optionsRef = useRef<MahjongOptions | null>(null);
  const reconnectRef = useRef<number | null>(null);
  const intentional = useRef(false);
  const reconnectAction = useRef<(options: MahjongOptions) => void>(
    () => undefined,
  );
  const [playerID] = useState(getPlayerConnectionID);
  const [name, setName] = useState(
    () => localStorage.getItem("boredgamz-player-name") ?? "Guest",
  );
  const accountName = useAuthStore((auth) => auth.user?.username);
  const [state, setState] = useState<MahjongState | null>(null);
  const [phase, setPhase] = useState<
    "idle" | "connecting" | "waiting" | "playing"
  >("idle");
  const [error, setError] = useState("");

  const close = useCallback(() => {
    intentional.current = true;
    optionsRef.current = null;
    if (reconnectRef.current !== null) clearTimeout(reconnectRef.current);
    const socket = socketRef.current;
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: "leave" }));
    }
    socket?.close(1000, "leaving mahjong table");
    socketRef.current = null;
    setState(null);
    setPhase("idle");
    setError("");
  }, []);

  const connect = useCallback(
    (input: MahjongOptions) => {
      const options = { ...input, name: (accountName ?? input.name.trim()) || "Guest" };
      intentional.current = false;
      optionsRef.current = options;
      setPhase("connecting");
      setError("");
      localStorage.setItem("boredgamz-player-name", options.name);
      setName(options.name);
      const socket = new WebSocket(socketURL(playerID, options));
      const previous = socketRef.current;
      socketRef.current = socket;
      previous?.close();
      socket.onmessage = (event) => {
        if (socketRef.current !== socket) return;
        try {
          const message = JSON.parse(String(event.data)) as {
            type: string;
            state?: MahjongState;
            message?: string;
          };
          if (message.type === "waiting") setPhase("waiting");
          if (message.type === "state" && message.state) {
            setState(message.state);
            setPhase("playing");
            setError("");
          }
          if (message.type === "error") {
            setError(message.message ?? "That tile action was not accepted.");
          }
        } catch {
          setError("The garden sent an unreadable update.");
        }
      };
      socket.onerror = () => {
        if (socketRef.current === socket) {
          setError("The mahjong room is temporarily unavailable.");
        }
      };
      socket.onclose = () => {
        if (
          socketRef.current !== socket ||
          intentional.current ||
          !optionsRef.current
        ) {
          return;
        }
        setError("Connection lost. Rejoining your place…");
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
  const pair = useCallback(
    (first: number, second: number) => send({ type: "pair", first, second }),
    [send],
  );
  const hint = useCallback(() => send({ type: "hint" }), [send]);
  const chat = useCallback(
    (content: string) => send({ type: "chat", content }),
    [send],
  );
  const rematch = useCallback(() => send({ type: "rematch" }), [send]);
  useEffect(() => close, [close]);

  return {
    playerID,
    name: accountName ?? name,
    state,
    phase,
    error,
    connect,
    close,
    pair,
    hint,
    chat,
    rematch,
  };
}
