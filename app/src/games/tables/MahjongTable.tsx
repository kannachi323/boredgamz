import { Lightbulb } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { DeadlineClock } from "../../components/table/Clocks";
import { PlayerTag } from "../../components/table/PlayerTag";
import { ResultBar } from "../../components/table/ResultBar";
import { TableShell } from "../../components/table/TableShell";
import type { MahjongState } from "../../hooks/useMahjongMatch";
import type { GameDefinition } from "../catalog";
import { RematchButton } from "./shared";

interface Props {
  game: GameDefinition;
  state: MahjongState;
  error: string;
  reconnecting: boolean;
  onPair: (first: number, second: number) => void;
  onHint: () => void;
  onChat: (text: string) => void;
  onRematch: () => void;
  onLeave: () => void;
}

export function MahjongTable({ game, state, error, reconnecting, onPair, onHint, onChat, onRematch, onLeave }: Props) {
  const me = state.viewerSeat;
  const playing = state.status === "playing";
  const mode = game.modes.find((item) => item.id === state.mode);
  const remaining = state.seats[me >= 0 ? me : 0]?.remaining ?? 0;
  const ready = me >= 0 && state.ready.includes(me);
  const result = describe(state, me);

  const rail = (
    <>
      <div className="players">
        {state.seats.map((seat, index) => (
          <PlayerTag
            key={seat.id || index}
            seat={index === 0 ? "one" : "two"}
            name={seat.name}
            you={index === me}
            connected={seat.connected}
            left={seat.left}
            status={`${seat.pairs} ${seat.pairs === 1 ? "pair" : "pairs"}`}
            score={state.mode === "duel" ? seat.remaining : undefined}
          />
        ))}
      </div>
      <div className="stat-row">
        <div><span className="label">Tiles left</span><strong className="num">{remaining}</strong></div>
        <div><span className="label">Free tiles</span><strong className="num">{state.free.length}</strong></div>
        {state.deadline && <div><span className="label">Time</span><DeadlineClock deadline={state.deadline} urgentAt={60} /></div>}
      </div>
      {me >= 0 && playing && (
        <div className="rail-actions">
          <button className="btn" type="button" disabled={state.hintsRemaining <= 0} onClick={onHint}>
            <Lightbulb size={15} /> Hint <span className="num">{state.hintsRemaining}</span>
          </button>
        </div>
      )}
    </>
  );

  return (
    <TableShell
      game={game}
      className="table--mahjong"
      meta={[mode?.name ?? "", state.reshuffles ? `${state.reshuffles} reshuffle${state.reshuffles > 1 ? "s" : ""}` : ""]}
      roomId={state.roomId}
      spectators={state.spectators}
      inviteCode={state.inviteCode}
      reconnecting={reconnecting}
      error={error}
      rail={rail}
      chat={{ messages: state.chat, viewerId: state.viewerId, onSend: onChat }}
      onLeave={onLeave}
    >
      <p className="turn-line" aria-live="polite">{playing ? (state.message || "Match two free tiles") : result.title}</p>
      <MahjongBoard state={state} onPair={onPair} />
      {!playing && (
        <ResultBar title={result.title} detail={result.detail} tone={result.tone}>
          {me >= 0 && !state.seats[1 - me]?.left && <RematchButton sent={ready} label="New layout" onClick={onRematch} />}
          <button className="btn" type="button" onClick={onLeave}>Leave</button>
        </ResultBar>
      )}
    </TableShell>
  );
}

function MahjongBoard({ state, onPair }: { state: MahjongState; onPair: (first: number, second: number) => void }) {
  const [selected, setSelected] = useState<number | null>(null);
  const free = useMemo(() => new Set(state.free), [state.free]);
  const hinted = useMemo(() => new Set(state.hint), [state.hint]);
  useEffect(() => setSelected(null), [state.removed]);
  const columns = Math.max(...state.tiles.map((tile) => tile.x)) + 1;
  const rows = Math.max(...state.tiles.map((tile) => tile.y)) + 1;
  const canPlay = state.status === "playing" && state.viewerSeat >= 0;

  function select(index: number) {
    if (selected === null) setSelected(index);
    else if (selected === index) setSelected(null);
    else {
      onPair(selected, index);
      setSelected(null);
    }
  }

  return (
    <div className="mahjong" style={{ "--cols": columns, "--rows": rows } as React.CSSProperties}>
      {state.tiles.map((tile, index) => {
        if (state.removed[index]) return null;
        const open = free.has(index);
        return (
          <button
            key={tile.id}
            type="button"
            className={`mj-tile mj-tile--${suit(tile.face)} ${open ? "is-free" : ""} ${selected === index ? "is-selected" : ""} ${hinted.has(index) ? "is-hinted" : ""}`}
            style={{ "--x": tile.x, "--y": tile.y, "--z": tile.z } as React.CSSProperties}
            disabled={!canPlay || !open}
            aria-pressed={selected === index}
            aria-label={`${tile.face}${open ? "" : ", blocked"}`}
            onClick={() => select(index)}
          >
            <span>{tile.face}</span>
          </button>
        );
      })}
    </div>
  );
}

function suit(face: string) {
  if (face.endsWith("萬") || face === "中") return "red";
  if (face.endsWith("筒") || face === "白") return "blue";
  if (face.endsWith("索") || face === "發") return "green";
  return "ink";
}

function describe(state: MahjongState, me: number) {
  if (state.mode !== "duel") {
    const cleared = state.seats[0].remaining === 0;
    return { title: cleared ? "Board cleared" : "Time's up", detail: `${state.seats[0].pairs + state.seats[1].pairs} pairs cleared together`, tone: cleared ? ("win" as const) : ("neutral" as const) };
  }
  if (state.winner < 0) return { title: "Dead even", detail: undefined, tone: "neutral" as const };
  if (me < 0) return { title: `${state.seats[state.winner].name} wins`, detail: undefined, tone: "neutral" as const };
  const won = state.winner === me;
  return { title: won ? "You win the duel" : `${state.seats[state.winner].name} wins`, detail: undefined, tone: won ? ("win" as const) : ("loss" as const) };
}
