import { useEffect, useState } from "react";
import { DeadlineClock } from "../../components/table/Clocks";
import { PlayerTag } from "../../components/table/PlayerTag";
import { ResultBar } from "../../components/table/ResultBar";
import { TableShell } from "../../components/table/TableShell";
import type { PokerCard, PokerState } from "../../hooks/usePokerMatch";
import type { GameDefinition } from "../catalog";
import { RematchButton } from "./shared";

const rankLabel: Record<number, string> = { 11: "J", 12: "Q", 13: "K", 14: "A" };
const suitGlyph = { s: "♠", h: "♥", d: "♦", c: "♣" } as const;
const suitName = { s: "spades", h: "hearts", d: "diamonds", c: "clubs" } as const;
const phaseLabel: Record<string, string> = {
  preflop: "Pre-flop", flop: "Flop", turn: "Turn", river: "River",
  "draw-bet": "First bet", draw: "Draw", "final-bet": "Final bet",
  showdown: "Showdown", "hand-finished": "Hand over",
};

interface Props {
  game: GameDefinition;
  state: PokerState;
  error: string;
  reconnecting: boolean;
  onAction: (action: string, amount?: number) => void;
  onDiscard: (positions: number[]) => void;
  onChat: (text: string) => void;
  onNext: () => void;
  onLeave: () => void;
}

export function PokerTable({ game, state, error, reconnecting, onAction, onDiscard, onChat, onNext, onLeave }: Props) {
  const me = state.viewerSeat;
  const seated = me >= 0;
  const bottom = seated ? me : 1;
  const top = 1 - bottom;
  const mySeat = seated ? state.seats[me] : null;
  const myTurn = seated && state.turn === me;
  const handOver = state.phase === "showdown" || state.phase === "hand-finished";
  const minRaise = state.currentBet + state.minimumRaise;
  const maxRaise = mySeat ? mySeat.roundBet + mySeat.stack : minRaise;
  const [raiseTo, setRaiseTo] = useState(minRaise);
  const [discards, setDiscards] = useState<number[]>([]);
  useEffect(() => {
    setRaiseTo(Math.min(minRaise, maxRaise));
    setDiscards([]);
  }, [state.hand, state.phase, state.turn, minRaise, maxRaise]);
  const pot = handOver ? state.awardedPot : state.pot + state.seats.reduce((sum, seat) => sum + seat.roundBet, 0);
  const ready = seated && state.ready.includes(me);
  const mode = game.modes.find((item) => item.id === state.mode);
  const result = describe(state, me);
  const canRaise = state.legalActions.includes("raise") && maxRaise > minRaise;

  function toggle(index: number) {
    if (state.phase !== "draw" || !myTurn) return;
    setDiscards((current) => (current.includes(index) ? current.filter((item) => item !== index) : current.length < 4 ? [...current, index] : current));
  }

  const rail = (
    <div className="players">
      {[top, bottom].map((index) => {
        const seat = state.seats[index];
        const active = !handOver && state.turn === index;
        return (
          <PlayerTag
            key={index}
            seat={index === 0 ? "one" : "two"}
            name={seat.name}
            you={index === me}
            active={active}
            connected={seat.connected}
            left={seat.left}
            status={[index === state.dealer ? "Dealer" : "", seat.folded ? "Folded" : seat.allIn ? "All in" : ""].filter(Boolean).join(" · ") || undefined}
            score={seat.stack.toLocaleString()}
            clock={active ? <DeadlineClock deadline={state.actionDeadline} urgentAt={8} /> : null}
          />
        );
      })}
    </div>
  );

  return (
    <TableShell
      game={game}
      className="table--poker"
      meta={[mode?.name ?? "", `Hand ${state.hand}`, state.mode === "holdem" ? "Blinds 10/20" : "Ante 10"]}
      roomId={state.roomId}
      spectators={state.spectators}
      inviteCode={state.inviteCode}
      reconnecting={reconnecting}
      error={error}
      rail={rail}
      chat={{ messages: state.chat, viewerId: state.viewerId, onSend: onChat }}
      onLeave={onLeave}
    >
      <div className="felt">
        <Seat state={state} index={top} hand={state.seats[top].cards} count={state.seats[top].cardCount} />
        <div className="felt__center">
          <div className="pot">
            <span className="label">{phaseLabel[state.phase] ?? state.phase}</span>
            <strong className="num">{pot.toLocaleString()}</strong>
          </div>
          {state.mode === "holdem" && (
            <div className="hand hand--board">
              {Array.from({ length: 5 }, (_, i) => (state.community[i] ? <Card key={state.community[i].code} card={state.community[i]} /> : <span key={i} className="card is-slot" />))}
            </div>
          )}
          {state.message && <p className="felt__message">{state.message}</p>}
        </div>
        <Seat state={state} index={bottom} hand={state.seats[bottom].cards} count={state.seats[bottom].cardCount} mine={seated} selected={discards} onCard={myTurn && state.phase === "draw" ? toggle : undefined} />
      </div>

      {!handOver && myTurn && (
        <div className="bet-bar">
          {state.phase === "draw" ? (
            <button className="btn btn--primary btn--lg" type="button" onClick={() => onDiscard(discards)}>
              {discards.length ? `Draw ${discards.length}` : "Stand pat"}
            </button>
          ) : (
            <>
              {state.legalActions.includes("fold") && <button className="btn btn--lg" type="button" onClick={() => onAction("fold")}>Fold</button>}
              {state.legalActions.includes("check") && <button className="btn btn--lg" type="button" onClick={() => onAction("check")}>Check</button>}
              {state.legalActions.includes("call") && <button className="btn btn--primary btn--lg" type="button" onClick={() => onAction("call")}>Call {state.toCall}</button>}
              {canRaise && (
                <div className="raise">
                  <label className="sr-only" htmlFor="raise-to">Raise to</label>
                  <input id="raise-to" type="range" min={minRaise} max={maxRaise} step={10} value={raiseTo} onChange={(event) => setRaiseTo(Number(event.target.value))} />
                  <button className="btn btn--ink btn--lg" type="button" onClick={() => onAction("raise", raiseTo)}>Raise to <span className="num">{raiseTo}</span></button>
                </div>
              )}
              {state.legalActions.includes("all-in") && <button className="btn btn--lg" type="button" onClick={() => onAction("all-in")}>All in</button>}
            </>
          )}
        </div>
      )}
      {!handOver && seated && !myTurn && <p className="turn-line">{state.seats[top].name} to act…</p>}
      {handOver && (
        <ResultBar title={result.title} detail={state.seats.some((seat) => seat.stack === 0) ? `${result.detail ? `${result.detail} · ` : ""}Stacks reset next hand` : result.detail} tone={result.tone}>
          {seated && !state.seats[top].left && <RematchButton sent={ready} label="Next hand" onClick={onNext} />}
          <button className={`btn ${state.seats[top].left ? "btn--primary" : ""}`} type="button" onClick={onLeave}>{seated && state.seats[top].left ? "Back to lobby" : "Leave"}</button>
        </ResultBar>
      )}
    </TableShell>
  );
}

function Seat({ state, index, hand, count, mine, selected = [], onCard }: { state: PokerState; index: number; hand?: PokerCard[]; count: number; mine?: boolean; selected?: number[]; onCard?: (index: number) => void }) {
  const seat = state.seats[index];
  return (
    <div className={`seat ${mine ? "is-mine" : ""} ${seat.folded ? "is-folded" : ""}`}>
      <div className="hand">
        {Array.from({ length: count }, (_, i) => {
          const card = hand?.[i];
          const picked = selected.includes(i);
          if (onCard && card) {
            return (
              <button key={card.code} type="button" className={`card-button ${picked ? "is-picked" : ""}`} aria-pressed={picked} aria-label={`${cardName(card)}${picked ? ", discarding" : ""}`} onClick={() => onCard(i)}>
                <Card card={card} />
              </button>
            );
          }
          return card ? <Card key={card.code} card={card} /> : <span key={i} className="card is-back" aria-label="Hidden card" />;
        })}
      </div>
      {seat.roundBet > 0 && <span className="bet-chip num">{seat.roundBet}</span>}
    </div>
  );
}

function Card({ card }: { card: PokerCard }) {
  const red = card.suit === "h" || card.suit === "d";
  return (
    <span className={`card ${red ? "is-red" : ""}`} role="img" aria-label={cardName(card)}>
      <b>{rankLabel[card.rank] ?? card.rank}</b>
      <i>{suitGlyph[card.suit]}</i>
    </span>
  );
}

function cardName(card: PokerCard) {
  return `${rankLabel[card.rank] ?? card.rank} of ${suitName[card.suit]}`;
}

function describe(state: PokerState, me: number) {
  const winners = state.winners;
  const hand = state.winningHand;
  if (winners.length > 1) return { title: "Split pot", detail: hand, tone: "neutral" as const };
  const winner = winners[0];
  if (winner === undefined) return { title: "Hand over", detail: hand, tone: "neutral" as const };
  if (me < 0) return { title: `${state.seats[winner].name} wins`, detail: hand, tone: "neutral" as const };
  const won = winner === me;
  const folded = state.phase === "hand-finished";
  return {
    title: won ? `You win ${state.awardedPot.toLocaleString()}` : `${state.seats[winner].name} wins ${state.awardedPot.toLocaleString()}`,
    detail: folded ? "Folded" : hand,
    tone: won ? ("win" as const) : ("loss" as const),
  };
}
