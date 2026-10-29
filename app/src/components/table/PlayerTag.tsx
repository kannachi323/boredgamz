import type { ReactNode } from "react";

interface Props {
  name?: string;
  seat: "one" | "two";
  you?: boolean;
  active?: boolean;
  connected?: boolean;
  left?: boolean;
  score?: ReactNode;
  status?: string;
  clock?: ReactNode;
  piece?: ReactNode;
}

export function PlayerTag({ name, seat, you, active, connected = true, left, score, status, clock, piece }: Props) {
  const line = left ? "Left the table" : !connected ? "Reconnecting…" : status;
  return (
    <div className={`player ${active ? "is-active" : ""} ${connected && !left ? "" : "is-away"}`}>
      <span className={`avatar ${seat === "two" ? "avatar--two" : ""}`} aria-hidden="true">{piece ?? (name ?? "?").trim().slice(0, 1).toLocaleUpperCase()}</span>
      <span className="player__who">
        <b>{name ?? "Waiting…"}{you && <em> you</em>}</b>
        {line && <small>{line}</small>}
      </span>
      {clock}
      {score !== undefined && <strong className="player__score num">{score}</strong>}
    </div>
  );
}
