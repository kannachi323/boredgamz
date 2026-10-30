import { useCallback, useEffect, useRef, useState } from "react";
import { useSearchParams } from "react-router-dom";
import type { JoinRequest } from "../../hooks/useGameSocket";
import { normalizeCode } from "../../games/catalog";
import { useAuthStore } from "../../stores/useAuthStore";

/**
 * Wires a game socket to the URL: ?join=CODE seats you at a private table,
 * ?watch=ROOM opens it as a spectator. Both wait until you're signed in.
 */
export function useRoomFlow(connect: (request: JoinRequest) => void, close: () => void, resume: () => JoinRequest | null) {
  const [params, setParams] = useSearchParams();
  const signedIn = useAuthStore((auth) => auth.status === "authenticated");
  const [request, setRequest] = useState<JoinRequest | null>(null);
  const handled = useRef("");

  const join = useCallback(
    (next: JoinRequest) => {
      setRequest(next);
      connect(next);
    },
    [connect],
  );

  useEffect(() => {
    if (!signedIn) return;
    const code = params.get("join");
    const watch = params.get("watch");
    const key = `${code}|${watch}`;
    if (handled.current === key) return;
    // Deferred and cancellable so a StrictMode remount can't half-open a socket.
    const timer = window.setTimeout(() => {
      handled.current = key;
      if (watch) join({ role: "spectator", roomId: watch.slice(0, 64) });
      else if (code) join({ inviteCode: normalizeCode(code) });
      else {
        const resumed = resume();
        if (resumed) setRequest(resumed);
      }
    }, 0);
    return () => window.clearTimeout(timer);
  }, [join, params, resume, signedIn]);

  const leave = useCallback(() => {
    close();
    setRequest(null);
    handled.current = "";
    if (params.has("join") || params.has("watch")) setParams({}, { replace: true });
  }, [close, params, setParams]);

  return { request, join, leave };
}
