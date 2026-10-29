import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, Copy, DoorOpen, Eye, Flag, History, RotateCcw } from "lucide-react";
import { useEffect, useState } from "react";
import { MatchChat } from "../../../components/games/MatchChat";
import type { ArenaPlayer, ArenaState } from "../../../hooks/useOnlineMatch";
import { TurnClock } from "../TicTacToe/TurnClock";
import { GomokuBoard } from "./GomokuBoard";
import type { BoardTheme } from "./boardThemes";

interface Props {
  state: ArenaState;
  theme: BoardTheme;
  playerID: string;
  error: string;
  onMove: (at: number) => void;
  onChat: (text: string) => void;
  onResign: () => void;
  onRematch: () => void;
  onLeave: () => void;
}

export function GomokuTable(props: Props) {
  const { state } = props;
  const [reviewPly, setReviewPly] = useState<number | null>(null);
  useEffect(() => setReviewPly(null), [state.roomId, state.round, state.gridMoves.length]);
  const me = state.players.find((p) => p?.mark === state.viewerMark) ?? null;
  const reviewing = reviewPly !== null;
  const canMove = !reviewing && state.status === "playing" && me?.mark === state.turn;
  const displayBoard = reviewing ? replayBoard(state.board.length, state.gridMoves, reviewPly) : state.board;
  const displayLastMove = reviewing && reviewPly >= 0 ? state.gridMoves[reviewPly]?.position ?? -1 : state.lastMove;
  const ready = me ? state.rematchReady.includes(me.mark) : false;
  async function share() {
    const url = new URL(window.location.href);
    url.search = "";
    url.searchParams.set("watch", state.roomId);
    await navigator.clipboard.writeText(url.toString());
  }
  return (
    <main className="gmk-table-page">
      <header className="gmk-topbar">
        <a href="/">
          Bored<span>Gamz</span> / GOMOKU
        </a>
        <div>
          <span>{labelMode(state.settings.mode)}</span>
          <i /> ROUND {state.round}
          <i /> BEST OF {state.settings.bestOf}
        </div>
        <button type="button" onClick={props.onLeave}>
          <DoorOpen size={15} /> Leave
        </button>
      </header>
      <div className="gmk-layout">
        <aside className="gmk-match-panel">
          <span className="gmk-overline">Match room</span>
          <h1>{headline(state, me)}</h1>
          <div className="gmk-players">
            <Player
              player={state.players[0]}
              active={state.status === "playing" && state.turn === "one"}
              score={state.score.one}
              dark
            />
            <Player
              player={state.players[1]}
              active={state.status === "playing" && state.turn === "two"}
              score={state.score.two}
            />
          </div>
          <div className="gmk-details">
            <div>
              <span>Move clock</span>
              <TurnClock
                deadline={state.turnDeadline}
                control={state.settings.timeControl}
              />
            </div>
            <div>
              <span>Observers</span>
              <b>
                <Eye size={13} /> {state.spectators}
              </b>
            </div>
            <div>
              <span>Rule</span>
              <b>{labelMode(state.settings.mode)}</b>
            </div>
          </div>
          {state.settings.inviteCode && (
            <p className="gmk-code">
              Code <strong>{state.settings.inviteCode}</strong>
            </p>
          )}
          <button
            className="gmk-share"
            type="button"
            onClick={() => void share()}
          >
            <Copy size={14} /> Copy spectator link
          </button>
          <div className="gmk-actions">
            {state.status === "playing" && me && (
              <button onClick={props.onResign}>
                <Flag size={14} /> Resign
              </button>
            )}
            {state.status !== "playing" && me && (
              <button
                className="is-primary"
                disabled={ready}
                onClick={props.onRematch}
              >
                <RotateCcw size={14} /> {ready ? "Waiting" : "Next game"}
              </button>
            )}
          </div>
          <section className="gmk-replay" aria-label="Move replay">
            <header><span><History size={13} /> Move history</span><b>{state.gridMoves.length} moves</b></header>
            {state.gridMoves.length === 0 ? <p>No stones have been placed.</p> : <>
              <ol>{state.gridMoves.map((move, index) => <li key={`${index}-${move.position}`}><button className={reviewPly === index ? "is-active" : ""} type="button" onClick={() => setReviewPly(index)}><i className={move.mark === "one" ? "is-dark" : ""} /><span>{index + 1}. {coordinate(move.position, state.columns)}</span></button></li>)}</ol>
              <div className="gmk-replay__controls">
                <button type="button" aria-label="Review from empty board" onClick={() => setReviewPly(-1)}><ChevronsLeft size={16} /></button>
                <button type="button" aria-label="Previous move" disabled={reviewPly === null || reviewPly < 0} onClick={() => setReviewPly(Math.max(-1, (reviewPly ?? 0) - 1))}><ChevronLeft size={16} /></button>
                <button type="button" aria-label="Next move" disabled={reviewPly === null || reviewPly >= state.gridMoves.length - 1} onClick={() => setReviewPly(Math.min(state.gridMoves.length - 1, (reviewPly ?? -1) + 1))}><ChevronRight size={16} /></button>
                <button type="button" aria-label="Return to live board" disabled={!reviewing} onClick={() => setReviewPly(null)}><ChevronsRight size={16} /></button>
              </div>
            </>}
          </section>
          {props.error && <p className="gmk-error">{props.error}</p>}
        </aside>
        <section className="gmk-board-stage">
          <div className="gmk-board-caption">
            <span>
              {reviewing ? `Reviewing move ${(reviewPly ?? -1) + 1}` : state.status === "playing" ? "Live board" : "Position final"}
            </span>
            <p>{prompt(state, me)}</p>
          </div>
          <GomokuBoard
            state={{ ...state, board: displayBoard, lastMove: displayLastMove }}
            theme={props.theme}
            canMove={canMove}
            onMove={props.onMove}
          />
          <p className="gmk-coordinate-note">
            Coordinates run A–{String.fromCharCode(64 + state.columns)} and 1–
            {state.rows}. Last move carries a cyan pin.
          </p>
        </section>
        <MatchChat
          messages={state.chat}
          playerID={props.playerID}
          onSend={props.onChat}
        />
      </div>
    </main>
  );
}

function replayBoard(size: number, moves: ArenaState["gridMoves"], index: number) {
  const board: ArenaState["board"] = Array(size).fill("");
  for (const move of moves.slice(0, index + 1)) board[move.position] = move.mark;
  return board;
}

function coordinate(position: number, columns: number) {
  return `${String.fromCharCode(65 + (position % columns))}${Math.floor(position / columns) + 1}`;
}

function Player({
  player,
  active,
  score,
  dark = false,
}: {
  player: ArenaPlayer | null;
  active: boolean;
  score: number;
  dark?: boolean;
}) {
  return (
    <div className={`gmk-player ${active ? "is-active" : ""}`}>
      <i className={dark ? "is-dark" : ""} />
      <span>
        <b>{player?.name ?? "Waiting"}</b>
        <small>
          {player?.connected === false
            ? "Reconnecting"
            : active
              ? "To move"
              : "Ready"}
        </small>
      </span>
      <strong>{score}</strong>
    </div>
  );
}
function labelMode(mode: string) {
  return mode === "pro"
    ? "Exact five"
    : mode === "compact"
      ? "Compact"
      : "Freestyle";
}
function headline(state: ArenaState, me: ArenaPlayer | null) {
  if (state.status === "playing")
    return me?.mark === state.turn
      ? "Your stone."
      : me
        ? "Read the board."
        : "Quietly watching.";
  if (!state.winner) return "Even position.";
  return !me
    ? "Five connected."
    : state.winner === me.mark
      ? "Five is yours."
      : "Line conceded.";
}
function prompt(state: ArenaState, me: ArenaPlayer | null) {
  if (state.status !== "playing")
    return state.seriesWinner
      ? "The match is complete."
      : "Both players may ready the next board.";
  return !me
    ? "A live club position."
    : state.turn === me.mark
      ? "Build a threat or answer theirs."
      : "Your opponent is considering the grid.";
}
