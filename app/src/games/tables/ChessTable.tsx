import { Chess, type Square } from "chess.js";
import { Handshake } from "lucide-react";
import { useMemo, useState } from "react";
import { BankClock } from "../../components/table/Clocks";
import { PlayerTag } from "../../components/table/PlayerTag";
import { ResultBar } from "../../components/table/ResultBar";
import { TableShell } from "../../components/table/TableShell";
import { useBoardTheme } from "../../components/room/useBoardTheme";
import type { ArenaPlayer, ArenaState } from "../../hooks/useOnlineMatch";
import type { GameDefinition } from "../catalog";
import { RematchButton, ResignButton } from "./shared";

// Solid glyphs for both colours (coloured in CSS); U+FE0E keeps them as text, not emoji.
const glyph: Record<string, string> = { k: "♚", q: "♛", r: "♜", b: "♝", n: "♞", p: "♟" };
const pieceNames: Record<string, string> = { p: "pawn", n: "knight", b: "bishop", r: "rook", q: "queen", k: "king" };

interface Props {
  game: GameDefinition;
  state: ArenaState<string>;
  error: string;
  reconnecting: boolean;
  onMove: (from: string, to: string, promotion?: string) => void;
  onChat: (text: string) => void;
  onResign: () => void;
  onRematch: () => void;
  onOfferDraw: () => void;
  onAcceptDraw: () => void;
  onDeclineDraw: () => void;
  onLeave: () => void;
}

export function ChessTable(props: Props) {
  const { game, state } = props;
  const me = state.players.find((player) => player?.mark === state.viewerMark) ?? null;
  const [theme, setTheme] = useBoardTheme(game);
  const colors = game.themes?.find((item) => item.id === theme) ?? game.themes![0];
  const flipped = me?.mark === "two";
  const order: ("one" | "two")[] = flipped ? ["one", "two"] : ["two", "one"];
  const playing = state.status === "playing";
  const moves = state.chess?.moves ?? [];
  const mode = game.modes.find((item) => item.id === state.settings.mode);
  const outcome = describe(state, me);
  const rematchSent = me ? state.rematchReady.includes(me.mark) : false;
  const opponentLeft = Boolean(me && state.players.find((player) => player && player.mark !== me.mark)?.left);
  const opponentOffers = playing && me && state.drawOfferedBy && state.drawOfferedBy !== me.mark;

  const rail = (
    <>
      <div className="players">
        {order.map((mark) => {
          const player = state.players[mark === "one" ? 0 : 1];
          const active = playing && state.turn === mark;
          return (
            <PlayerTag
              key={mark}
              seat={mark}
              name={player?.name}
              you={player?.mark === state.viewerMark}
              active={active}
              connected={player?.connected ?? true}
              left={player?.left}
              piece={<span className={`chess-piece is-${mark === "one" ? "w" : "b"}`}>{glyph.k}{"︎"}</span>}
              status={mark === "one" ? "White" : "Black"}
              clock={<BankClock milliseconds={state.clocks?.[mark]} running={active && moves.length > 0} />}
            />
          );
        })}
      </div>
      <ol className="move-list" aria-label="Moves">
        {pairs(moves.map((move) => move.san)).map(([white, black], index) => (
          <li key={index}><span className="num">{index + 1}.</span><b>{white}</b><b>{black}</b></li>
        ))}
      </ol>
      {playing && me && (
        <div className="rail-actions">
          {!state.drawOfferedBy && state.turn !== me.mark && moves.length > 0 && (
            <button className="btn" type="button" onClick={props.onOfferDraw}><Handshake size={15} /> Offer draw</button>
          )}
          {state.drawOfferedBy === me.mark && <span className="chip">Draw offered</span>}
          <ResignButton onResign={props.onResign} />
        </div>
      )}
    </>
  );

  const settings = (
    <div className="swatches" role="group" aria-label="Board colours">
      {game.themes!.map((item) => (
        <button key={item.id} type="button" className="swatch" aria-pressed={theme === item.id} aria-label={item.label} title={item.label} onClick={() => setTheme(item.id)}>
          <i style={{ background: `linear-gradient(135deg, ${item.light} 50%, ${item.dark} 50%)` }} />
        </button>
      ))}
    </div>
  );

  return (
    <TableShell
      game={game}
      className="table--chess"
      meta={[mode?.name ?? "", mode?.detail ?? ""]}
      roomId={state.roomId}
      spectators={state.spectators}
      inviteCode={state.settings.inviteCode}
      reconnecting={props.reconnecting}
      error={props.error}
      rail={rail}
      settings={settings}
      chat={{ messages: state.chat, viewerId: state.viewerId, onSend: props.onChat }}
      onLeave={props.onLeave}
    >
      <p className={`turn-line ${playing && me?.mark === state.turn ? "is-mine" : ""}`} aria-live="polite">{outcome.line}</p>
      <ChessBoard state={state} myMark={me?.mark} flipped={flipped} light={colors.light} dark={colors.dark} onMove={props.onMove} />
      {opponentOffers && (
        <ResultBar title="Draw offered">
          <button className="btn" type="button" onClick={props.onDeclineDraw}>Decline</button>
          <button className="btn btn--primary" type="button" onClick={props.onAcceptDraw}>Accept</button>
        </ResultBar>
      )}
      {!playing && (
        <ResultBar title={outcome.title} detail={outcome.detail} tone={outcome.tone}>
          {me && !opponentLeft && <RematchButton sent={rematchSent} label="Rematch" onClick={props.onRematch} />}
          <button className={`btn ${opponentLeft ? "btn--primary" : ""}`} type="button" onClick={props.onLeave}>{opponentLeft ? "Back to lobby" : "Leave"}</button>
        </ResultBar>
      )}
    </TableShell>
  );
}

function ChessBoard({ state, myMark, flipped, light, dark, onMove }: { state: ArenaState<string>; myMark?: string; flipped: boolean; light: string; dark: string; onMove: (from: string, to: string, promotion?: string) => void }) {
  const chess = useMemo(() => new Chess(state.chess?.fen), [state.chess?.fen]);
  const [selected, setSelected] = useState<Square | null>(null);
  const [promotion, setPromotion] = useState<{ from: Square; to: Square } | null>(null);
  const myColor = myMark === "one" ? "w" : myMark === "two" ? "b" : null;
  const canMove = state.status === "playing" && state.turn === myMark;
  const legal = selected && canMove ? chess.moves({ square: selected, verbose: true }) : [];
  const inCheck = chess.inCheck() ? chess.turn() : null;

  function choose(square: Square) {
    if (!canMove) return;
    const targets = legal.filter((move) => move.to === square);
    if (selected && targets.length) {
      if (targets.some((move) => move.promotion)) setPromotion({ from: selected, to: square });
      else onMove(selected, square);
      setSelected(null);
      return;
    }
    const piece = chess.get(square);
    setSelected(piece && piece.color === myColor && square !== selected ? square : null);
  }

  const squares = Array.from({ length: 64 }, (_, shown) => (flipped ? 63 - shown : shown));
  return (
    <div className="chess" style={{ "--sq-light": light, "--sq-dark": dark } as React.CSSProperties}>
      <div className="chess__board">
        {squares.map((index, shown) => {
          const square = `${String.fromCharCode(97 + (index % 8))}${8 - Math.floor(index / 8)}` as Square;
          const code = state.board[index];
          const isLight = (Math.floor(index / 8) + (index % 8)) % 2 === 0;
          const target = legal.find((move) => move.to === square);
          const last = state.chess?.lastFrom === square || state.chess?.lastTo === square;
          const checked = code && code[1] === "k" && code[0] === inCheck;
          return (
            <button
              key={square}
              type="button"
              className={`sq ${isLight ? "is-light" : "is-dark"} ${selected === square ? "is-selected" : ""} ${target ? (code ? "is-capture" : "is-target") : ""} ${last ? "is-last" : ""} ${checked ? "is-check" : ""}`}
              onClick={() => choose(square)}
              aria-label={`${square}${code ? `, ${code[0] === "w" ? "white" : "black"} ${pieceNames[code[1]]}` : ""}`}
              aria-pressed={selected === square}
            >
              {code && <span className={`chess-piece is-${code[0]}`} aria-hidden="true">{glyph[code[1]]}{"︎"}</span>}
              {shown % 8 === 0 && <small className="sq__rank" aria-hidden="true">{square[1]}</small>}
              {shown >= 56 && <small className="sq__file" aria-hidden="true">{square[0]}</small>}
            </button>
          );
        })}
      </div>
      {promotion && (
        <div className="promotion" role="dialog" aria-label="Promote pawn">
          {(["q", "r", "b", "n"] as const).map((piece) => (
            <button key={piece} type="button" className="btn btn--icon" aria-label={pieceNames[piece]} onClick={() => { onMove(promotion.from, promotion.to, piece); setPromotion(null); }}>
              <span className={`chess-piece is-${myColor}`}>{glyph[piece]}{"︎"}</span>
            </button>
          ))}
          <button className="btn btn--ghost btn--sm" type="button" onClick={() => setPromotion(null)}>Cancel</button>
        </div>
      )}
    </div>
  );
}

function pairs(moves: string[]) {
  const rows: [string, string?][] = [];
  for (let i = 0; i < moves.length; i += 2) rows.push([moves[i], moves[i + 1]]);
  return rows;
}

function describe(state: ArenaState<string>, me: ArenaPlayer | null) {
  const colorName = (mark?: string) => (mark === "one" ? "White" : "Black");
  const method = state.chess?.method && state.chess.method !== "NoMethod" ? humanize(state.chess.method) : undefined;
  if (state.status === "playing") {
    const check = state.chess?.fen && new Chess(state.chess.fen).inCheck() ? " · Check" : "";
    const line = me ? (state.turn === me.mark ? `Your move${check}` : `${colorName(state.turn)} is thinking…`) : `${colorName(state.turn)} to move${check}`;
    return { line, title: "", tone: "neutral" as const };
  }
  if (!state.winner) return { line: "Draw", title: "Draw", detail: method, tone: "neutral" as const };
  if (!me) return { line: `${colorName(state.winner)} wins`, title: `${colorName(state.winner)} wins`, detail: method, tone: "neutral" as const };
  const won = state.winner === me.mark;
  return { line: won ? "You won" : "You lost", title: won ? "You win" : "You lose", detail: method, tone: won ? ("win" as const) : ("loss" as const) };
}

function humanize(method: string) {
  return method.replace(/([a-z])([A-Z])/g, "$1 $2").replace(/^./, (c) => c.toUpperCase());
}
