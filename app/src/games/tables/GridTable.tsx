import { ChevronLeft, ChevronRight, ChevronsRight } from "lucide-react";
import { useEffect, useState } from "react";
import { DeadlineClock } from "../../components/table/Clocks";
import { PlayerTag } from "../../components/table/PlayerTag";
import { ResultBar } from "../../components/table/ResultBar";
import { TableShell } from "../../components/table/TableShell";
import { useBoardTheme } from "../../components/room/useBoardTheme";
import type { ArenaPlayer, ArenaState, GridMark } from "../../hooks/useOnlineMatch";
import type { GameDefinition } from "../catalog";
import { clockLabel, seriesLabel } from "../labels";
import { RematchButton, ResignButton } from "./shared";

type Mark = Exclude<GridMark, "">;

interface Props {
  game: GameDefinition;
  state: ArenaState;
  error: string;
  reconnecting: boolean;
  onMove: (position: number) => void;
  onChat: (text: string) => void;
  onResign: () => void;
  onRematch: () => void;
  onLeave: () => void;
}

export function GridTable({ game, state, error, reconnecting, onMove, onChat, onResign, onRematch, onLeave }: Props) {
  const me = state.players.find((player) => player?.mark === state.viewerMark) ?? null;
  const [theme, setTheme] = useBoardTheme(game);
  const [replay, setReplay] = useState<number | null>(null);
  useEffect(() => setReplay(null), [state.round, state.gridMoves.length]);
  const reviewing = replay !== null;
  const canMove = !reviewing && state.status === "playing" && me?.mark === state.turn;
  const board = reviewing ? replayBoard(state.board.length, state.gridMoves, replay) : state.board;
  const lastMove = reviewing ? state.gridMoves[replay]?.position ?? -1 : state.lastMove;
  const mode = game.modes.find((item) => item.id === state.settings.mode);
  const outcome = describe(state, me);
  const rematchSent = me ? state.rematchReady.includes(me.mark) : false;
  const opponentLeft = Boolean(me && state.players.find((player) => player && player.mark !== me.mark)?.left);

  const rail = (
    <>
      <div className="players">
        {state.players.map((player, index) => {
          const mark: Mark = index === 0 ? "one" : "two";
          const active = state.status === "playing" && state.turn === mark;
          return (
            <PlayerTag
              key={mark}
              seat={mark}
              name={player?.name}
              you={player?.mark === state.viewerMark}
              active={active}
              connected={player?.connected ?? true}
              left={player?.left}
              score={state.score[mark]}
              status={active ? "To move" : undefined}
              piece={<Piece kind={game.slug} mark={mark} />}
              clock={active ? <DeadlineClock deadline={state.turnDeadline} /> : null}
            />
          );
        })}
      </div>
      {game.slug === "gomoku" && state.gridMoves.length > 0 && (
        <div className="replay" aria-label="Move replay">
          <span className="label">Moves {state.gridMoves.length}</span>
          <div className="replay__controls">
            <button className="btn btn--icon btn--sm" type="button" aria-label="Previous move" disabled={replay === -1} onClick={() => setReplay((value) => Math.max(-1, (value ?? state.gridMoves.length - 1) - 1))}><ChevronLeft size={16} /></button>
            <span className="num">{reviewing ? replay + 1 : state.gridMoves.length}</span>
            <button className="btn btn--icon btn--sm" type="button" aria-label="Next move" disabled={!reviewing} onClick={() => setReplay((value) => (value === null || value + 1 >= state.gridMoves.length - 1 ? null : value + 1))}><ChevronRight size={16} /></button>
            <button className="btn btn--icon btn--sm" type="button" aria-label="Back to live" disabled={!reviewing} onClick={() => setReplay(null)}><ChevronsRight size={16} /></button>
          </div>
        </div>
      )}
      {state.status === "playing" && me && (
        <div className="rail-actions"><ResignButton onResign={onResign} /></div>
      )}
    </>
  );

  const settings = game.themes && (
    <div className="swatches" role="group" aria-label="Board style">
      {game.themes.map((item) => (
        <button key={item.id} type="button" className="swatch" aria-pressed={theme === item.id} aria-label={item.label} title={item.label} onClick={() => setTheme(item.id)}>
          <i style={{ background: `linear-gradient(135deg, ${item.light} 50%, ${item.dark} 50%)` }} />
        </button>
      ))}
    </div>
  );

  return (
    <TableShell
      game={game}
      className={`table--${game.slug}`}
      meta={[mode?.name ?? "", seriesLabel(state.settings.bestOf), state.settings.bestOf > 1 ? `Round ${state.round}` : "", clockLabel(state.settings.timeControl)]}
      roomId={state.roomId}
      spectators={state.spectators}
      inviteCode={state.settings.inviteCode}
      reconnecting={reconnecting}
      error={error}
      rail={rail}
      settings={settings || undefined}
      chat={{ messages: state.chat, viewerId: state.viewerId, onSend: onChat }}
      onLeave={onLeave}
    >
      <p className={`turn-line ${canMove ? "is-mine" : ""}`} aria-live="polite">{reviewing ? `Reviewing move ${replay + 1}` : outcome.line}</p>
      {game.slug === "connect-four" ? (
        <ConnectFourBoard state={state} board={board} lastMove={lastMove} canMove={canMove} onMove={onMove} />
      ) : (
        <div
          className={`board-grid board-grid--${game.slug} ${game.slug === "gomoku" ? `wood-${theme}` : ""} ${canMove ? "can-move" : ""}`}
          style={{ "--cols": state.columns, "--rows": state.rows } as React.CSSProperties}
        >
          {board.map((mark, index) => (
            <button
              key={index}
              type="button"
              className={`cell ${mark ? `is-${mark}` : ""} ${lastMove === index ? "is-last" : ""}`}
              disabled={!canMove || Boolean(mark)}
              onClick={() => onMove(index)}
              aria-label={cellLabel(game.slug, index, state.columns, mark)}
            >
              {mark && <span className="mark" aria-hidden="true" />}
            </button>
          ))}
        </div>
      )}
      {state.status !== "playing" && (
        <ResultBar title={outcome.title} detail={outcome.detail} tone={outcome.tone}>
          {me && !opponentLeft && <RematchButton sent={rematchSent} label={state.status === "finished" ? "Rematch" : "Next round"} onClick={onRematch} />}
          <button className={`btn ${opponentLeft ? "btn--primary" : ""}`} type="button" onClick={onLeave}>{opponentLeft ? "Back to lobby" : "Leave"}</button>
        </ResultBar>
      )}
    </TableShell>
  );
}

function ConnectFourBoard({ state, board, lastMove, canMove, onMove }: { state: ArenaState; board: GridMark[]; lastMove: number; canMove: boolean; onMove: (column: number) => void }) {
  const columns = Array.from({ length: state.columns }, (_, column) => column);
  return (
    <div className={`c4 ${canMove ? "can-move" : ""}`} style={{ "--cols": state.columns, "--rows": state.rows } as React.CSSProperties}>
      {columns.map((column) => {
        const full = Boolean(board[column]);
        return (
          <button key={column} type="button" className="c4__column" disabled={!canMove || full} onClick={() => onMove(column)} aria-label={`Drop in column ${column + 1}${full ? " (full)" : ""}`}>
            {Array.from({ length: state.rows }, (_, row) => {
              const index = row * state.columns + column;
              const mark = board[index];
              return <span key={row} className={`c4__hole ${mark ? `is-${mark}` : ""} ${lastMove === index ? "is-last" : ""}`} style={{ "--row": row } as React.CSSProperties} />;
            })}
          </button>
        );
      })}
    </div>
  );
}

function Piece({ kind, mark }: { kind: string; mark: Mark }) {
  if (kind === "tictactoe") return <>{mark === "one" ? "×" : "○"}</>;
  return <span className={`piece-dot is-${mark} piece-dot--${kind}`} />;
}

function cellLabel(kind: string, index: number, columns: number, mark: GridMark) {
  const row = Math.floor(index / columns) + 1;
  const column = (index % columns) + 1;
  const where = kind === "gomoku" ? `${String.fromCharCode(64 + column)}${row}` : `Row ${row}, column ${column}`;
  if (!mark) return `${where}, empty`;
  const names: Record<string, [string, string]> = { tictactoe: ["X", "O"], gomoku: ["black stone", "white stone"] };
  const [one, two] = names[kind] ?? ["first player", "second player"];
  return `${where}, ${mark === "one" ? one : two}`;
}

function replayBoard(size: number, moves: { position: number; mark: Mark }[], upTo: number) {
  const board: GridMark[] = Array(size).fill("");
  for (let i = 0; i <= upTo; i++) board[moves[i].position] = moves[i].mark;
  return board;
}

function describe(state: ArenaState, me: ArenaPlayer | null) {
  const nameOf = (mark?: string) => state.players.find((player) => player?.mark === mark)?.name ?? "Player";
  if (state.status === "playing") {
    const line = me ? (state.turn === me.mark ? "Your move" : `${nameOf(state.turn)} is thinking…`) : `${nameOf(state.turn)} to move`;
    return { line, title: "", tone: "neutral" as const };
  }
  const matchOver = state.status === "finished";
  const winner = matchOver && state.seriesWinner ? state.seriesWinner : state.winner;
  const series = state.settings.bestOf > 1;
  const unit = !series ? "" : matchOver ? " the match" : " the round";
  const detail = series ? `${state.score.one} – ${state.score.two}` : undefined;
  if (!winner) return { line: "Draw", title: "Draw", detail, tone: "neutral" as const };
  if (!me) return { line: `${nameOf(winner)} wins`, title: `${nameOf(winner)} wins${unit}`, detail, tone: "neutral" as const };
  const won = winner === me.mark;
  return { line: won ? "You won" : "You lost", title: won ? `You win${unit}` : `${nameOf(winner)} wins${unit}`, detail, tone: won ? ("win" as const) : ("loss" as const) };
}
