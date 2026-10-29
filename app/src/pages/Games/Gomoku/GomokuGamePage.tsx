import { useEffect, useState } from "react";
import type { MatchOptions } from "../../../hooks/useOnlineMatch";
import { useOnlineMatch } from "../../../hooks/useOnlineMatch";
import { GomokuLobby } from "./GomokuLobby";
import { GomokuTable } from "./GomokuTable";
import type { BoardTheme } from "./boardThemes";

export function GomokuGamePage() {
  const match = useOnlineMatch("gomoku");
  const connectMatch = match.connect;
  const playerName = match.name;
  const [theme, setTheme] = useState<BoardTheme>(
    () => (localStorage.getItem("gomoku-board-theme") as BoardTheme) || "mid",
  );
  const [pending, setPending] = useState<MatchOptions | null>(null);
  useEffect(() => {
    const roomId = new URLSearchParams(location.search).get("watch");
    if (!roomId) return;
    const options: MatchOptions = {
      name: playerName,
      role: "spectator",
      roomId,
    };
    setPending(options);
    connectMatch(options);
  }, [connectMatch, playerName]);
  function join(options: MatchOptions, nextTheme: BoardTheme) {
    localStorage.setItem("gomoku-board-theme", nextTheme);
    setTheme(nextTheme);
    setPending(options);
    match.connect(options);
  }
  if (match.state)
    return (
      <GomokuTable
        state={match.state}
        theme={theme}
        playerID={match.playerID}
        error={match.error}
        onMove={match.move}
        onChat={match.chat}
        onResign={match.resign}
        onRematch={match.rematch}
        onLeave={match.close}
      />
    );
  if (match.phase === "connecting" || match.phase === "waiting")
    return (
      <main className="gmk-waiting">
        <div>
          <i />
          <i />
        </div>
        <span>
          {pending?.role === "spectator"
            ? "Opening observer seat"
            : "Searching the club"}
        </span>
        <h1>Preparing the board.</h1>
        {pending?.inviteCode && (
          <p>
            Invite code <strong>{pending.inviteCode}</strong>
          </p>
        )}
        <small>
          {match.error || "Matched only with the same rule, series, and clock."}
        </small>
        <button onClick={match.close}>Cancel</button>
      </main>
    );
  return <GomokuLobby name={match.name} onJoin={join} />;
}
