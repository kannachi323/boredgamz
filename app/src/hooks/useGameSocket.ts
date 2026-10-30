import { useCallback, useEffect, useRef, useState } from "react";
import { getPlayerConnectionID } from "../lib/playerIdentity";
import { getSocketRoot, type SocketService } from "../lib/socketEndpoints";
import { useAuthStore } from "../stores/useAuthStore";

export type SocketPhase = "idle" | "connecting" | "waiting" | "live";

export interface JoinRequest {
  mode?: string;
  inviteCode?: string;
  role?: "player" | "spectator";
  roomId?: string;
  bestOf?: number;
  timeControl?: string;
}

export interface ChatMessage {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  createdAt: string;
}

interface ServerMessage<S> {
  type: string;
  state?: S;
  message?: string;
}

const maxBackoff = 10_000;

interface SavedSeat {
  request: JoinRequest;
  live: boolean;
}

// The tab ID lives in sessionStorage, so after a reload the server recognises
// this tab and hands back its seat. We only need to remember what we joined.
function readSeat(key: string): SavedSeat | null {
  try {
    const raw = sessionStorage.getItem(key);
    return raw ? (JSON.parse(raw) as SavedSeat) : null;
  } catch {
    return null;
  }
}

function writeSeat(key: string, seat: SavedSeat | null) {
  try {
    if (seat) sessionStorage.setItem(key, JSON.stringify(seat));
    else sessionStorage.removeItem(key);
  } catch {
    // Without storage a reload simply returns to the lobby.
  }
}

function socketURL(service: SocketService, path: string, playerID: string, request: JoinRequest) {
  const params = new URLSearchParams({
    playerId: request.role === "spectator" ? `${playerID}:spectator` : playerID,
  });
  for (const key of ["mode", "inviteCode", "role", "roomId", "timeControl"] as const) {
    const value = request[key];
    if (value) params.set(key, value);
  }
  if (request.bestOf) params.set("bestOf", String(request.bestOf));
  const root = getSocketRoot(service).replace(/\/$/, "");
  return `${root}${path}?${params}`;
}

/**
 * One live table connection. Handles waiting/matched messages, reconnects
 * with backoff (the server keeps the seat during a grace period), and
 * refreshes the session once when the handshake is refused.
 */
export function useGameSocket<S>(service: SocketService, path = "", isOver: (state: S) => boolean = () => false) {
  const socketRef = useRef<WebSocket | null>(null);
  const requestRef = useRef<JoinRequest | null>(null);
  const retryTimer = useRef<number | null>(null);
  const attempts = useRef(0);
  const connectRef = useRef<(request: JoinRequest) => void>(() => undefined);
  const [playerID] = useState(getPlayerConnectionID);
  const [state, setState] = useState<S | null>(null);
  const [phase, setPhase] = useState<SocketPhase>("idle");
  const [error, setError] = useState("");
  const [reconnecting, setReconnecting] = useState(false);
  const refreshSession = useAuthStore((auth) => auth.refresh);
  const seatKey = `boredgamz-seat:${service}${path}`;
  const resumingLive = useRef(false);
  const latest = useRef<{ state: S | null; isOver: (state: S) => boolean }>({ state: null, isOver });
  latest.current = { state, isOver };

  const clearRetry = () => {
    if (retryTimer.current !== null) window.clearTimeout(retryTimer.current);
    retryTimer.current = null;
  };

  const close = useCallback(() => {
    writeSeat(seatKey, null);
    resumingLive.current = false;
    requestRef.current = null;
    clearRetry();
    const socket = socketRef.current;
    socketRef.current = null;
    if (socket?.readyState === WebSocket.OPEN) socket.send(JSON.stringify({ type: "leave" }));
    socket?.close(1000, "leaving table");
    setState(null);
    setPhase("idle");
    setError("");
    setReconnecting(false);
  }, [seatKey]);

  const connect = useCallback(
    (request: JoinRequest) => {
      clearRetry();
      requestRef.current = request;
      if (!readSeat(seatKey)) writeSeat(seatKey, { request, live: false });
      setError("");
      setPhase((current) => (current === "live" ? current : "connecting"));
      const socket = new WebSocket(socketURL(service, path, playerID, request));
      const previous = socketRef.current;
      socketRef.current = socket;
      previous?.close();
      let opened = false;

      socket.onopen = () => {
        opened = true;
        attempts.current = 0;
        setReconnecting(false);
        setError("");
      };
      socket.onmessage = (event) => {
        if (socketRef.current !== socket) return;
        let message: ServerMessage<S>;
        try {
          message = JSON.parse(String(event.data)) as ServerMessage<S>;
        } catch {
          return;
        }
        if (message.type === "waiting") {
          if (resumingLive.current) {
            // The match we were resuming is over; don't queue for a new one.
            requestRef.current = null;
            writeSeat(seatKey, null);
            resumingLive.current = false;
            socket.close(1000, "match ended");
            setPhase("idle");
            setError("That match has ended.");
            return;
          }
          setPhase("waiting");
        }
        if (message.type === "state" && message.state) {
          resumingLive.current = false;
          writeSeat(seatKey, { request, live: true });
          setState(message.state);
          setPhase("live");
          setError("");
        }
        if (message.type === "error") setError(message.message ?? "That move wasn't accepted.");
      };
      socket.onclose = () => {
        if (socketRef.current !== socket || !requestRef.current) return;
        const pending = requestRef.current;
        attempts.current += 1;
        setReconnecting(true);
        const retry = () => {
          if (requestRef.current === pending) connectRef.current(pending);
        };
        if (!opened && attempts.current <= 3) {
          // A refused handshake is usually an expired access cookie.
          void refreshSession().then((result) => {
            if (result === "offline") {
              setError("Can't reach the server. Retrying…");
              retryTimer.current = window.setTimeout(retry, Math.min(maxBackoff, 1500 * attempts.current));
              return;
            }
            if (result === "signed-out") {
              requestRef.current = null;
              writeSeat(seatKey, null);
              setReconnecting(false);
              setPhase("idle");
              setState(null);
              setError("Your session ended. Sign in to keep playing.");
              return;
            }
            retryTimer.current = window.setTimeout(retry, 300);
          });
          return;
        }
        if (!opened) setError("Can't reach the server. Retrying…");
        const delay = Math.min(maxBackoff, 800 * 2 ** Math.min(attempts.current, 4));
        retryTimer.current = window.setTimeout(retry, delay);
      };
    },
    [path, playerID, refreshSession, seatKey, service],
  );
  connectRef.current = connect;

  /** Rejoins the table this tab was at before a reload. */
  const resume = useCallback(() => {
    const seat = readSeat(seatKey);
    if (!seat) return null;
    resumingLive.current = seat.live;
    connect(seat.request);
    return seat.request;
  }, [connect, seatKey]);

  const send = useCallback((payload: object) => {
    const socket = socketRef.current;
    if (socket?.readyState !== WebSocket.OPEN) {
      setError("Reconnecting. Try again in a moment.");
      return;
    }
    setError("");
    socket.send(JSON.stringify(payload));
  }, []);

  const chat = useCallback((content: string) => send({ type: "chat", content }), [send]);

  // Unmounting (Back button, route change) only drops the socket: the server
  // holds the seat for a grace period and returning to the page resumes it.
  // Only the explicit Leave button forfeits.
  useEffect(
    () => () => {
      clearRetry();
      requestRef.current = null;
      const socket = socketRef.current;
      socketRef.current = null;
      const { state: last, isOver: over } = latest.current;
      if (last && over(last)) {
        // Nothing left to resume: free the seat instead of coming back to an old result.
        writeSeat(seatKey, null);
        if (socket?.readyState === WebSocket.OPEN) socket.send(JSON.stringify({ type: "leave" }));
      }
      socket?.close(1000, "page closed");
    },
    [seatKey],
  );

  return { playerID, state, phase, error, reconnecting, connect, close, resume, send, chat };
}
