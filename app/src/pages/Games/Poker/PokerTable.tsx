import { Copy, DoorOpen, Eye, RotateCcw } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { MatchChat } from "../../../components/games/MatchChat";
import type { PokerState } from "../../../hooks/usePokerMatch";
import { PokerCard } from "./PokerCard";
import { PokerClock } from "./PokerClock";

interface Props {
  state: PokerState;
  playerID: string;
  error: string;
  onAction: (action: string, amount?: number) => void;
  onDiscard: (positions: number[]) => void;
  onChat: (text: string) => void;
  onNext: () => void;
  onLeave: () => void;
}

export function PokerTable(props: Props) {
  const { state } = props;
  const me = state.seats.findIndex((seat) => seat.id === props.playerID);
  const opponent = me === 0 ? 1 : 0;
  const topSeat = me < 0 ? 0 : opponent;
  const bottomSeat = me < 0 ? 1 : me;
  const mySeat = me >= 0 ? state.seats[me] : null;
  const [selected, setSelected] = useState<number[]>([]);
  const minimumTarget = state.currentBet + state.minimumRaise;
  const maximumTarget = mySeat ? mySeat.roundBet + mySeat.stack : minimumTarget;
  const [raiseTo, setRaiseTo] = useState(minimumTarget);
  useEffect(() => {
    setRaiseTo(Math.min(maximumTarget, minimumTarget));
    setSelected([]);
  }, [state.hand, state.phase, state.turn, minimumTarget, maximumTarget]);
  const handFinished =
    state.phase === "showdown" || state.phase === "hand-finished";
  const totalPot = handFinished
    ? state.awardedPot
    : state.pot + state.seats.reduce((total, seat) => total + seat.roundBet, 0);
  const myTurn = me >= 0 && state.turn === me;
  const ready = me >= 0 && state.ready.includes(me);
  const community = useMemo(
    () => Array.from({ length: 5 }, (_, index) => state.community[index]),
    [state.community],
  );

  async function share() {
    const url = new URL(location.href);
    url.search = "";
    url.searchParams.set("watch", state.roomId);
    await navigator.clipboard.writeText(url.toString());
  }
  function toggleCard(index: number) {
    if (state.phase !== "draw" || !myTurn) return;
    setSelected((current) =>
      current.includes(index)
        ? current.filter((item) => item !== index)
        : current.length < 4
          ? [...current, index]
          : current,
    );
  }

  return (
    <main className="pkr-table-page">
      <header className="pkr-topbar">
        <a href="/">
          Bored<span>Gamz</span> / CARD ROOM
        </a>
        <div>
          <b>
            {state.mode === "holdem" ? "NO-LIMIT HOLD’EM" : "FIVE-CARD DRAW"}
          </b>
          <i /> HAND {state.hand}
          <i /> 1,000 CHIP BUY-IN
        </div>
        <button type="button" onClick={props.onLeave}>
          <DoorOpen size={15} /> Leave table
        </button>
      </header>
      <div className="pkr-room">
        <aside className="pkr-sidebar">
          <span className="pkr-overline">
            TABLE #{state.roomId.slice(-6).toUpperCase()}
          </span>
          <h1>{roomHeadline(state, me)}</h1>
          <p>{state.message}</p>
          <dl>
            <div>
              <dt>Pot</dt>
              <dd>{totalPot.toLocaleString()}</dd>
            </div>
            <div>
              <dt>Blinds / ante</dt>
              <dd>{state.mode === "holdem" ? "10 / 20" : "10"}</dd>
            </div>
            <div>
              <dt>Observers</dt>
              <dd>
                <Eye size={13} /> {state.spectators}
              </dd>
            </div>
          </dl>
          {state.inviteCode && (
            <p className="pkr-invite">
              Private code <b>{state.inviteCode}</b>
            </p>
          )}
          <button
            className="pkr-share"
            type="button"
            onClick={() => void share()}
          >
            <Copy size={14} /> Copy spectator link
          </button>
          <section className="pkr-integrity">
            <b>Fair table contract</b>
            <ul>
              <li>Server-only deck and wagers</li>
              <li>Private cards sent per seat</li>
              <li>30-second action clock</li>
              <li>Reconnect grace before fold</li>
            </ul>
          </section>
          {props.error && (
            <p className="pkr-error" role="alert">
              {props.error}
            </p>
          )}
        </aside>
        <section className="pkr-stage">
          <div className="pkr-felt">
            <Seat
              state={state}
              index={topSeat}
              active={!handFinished && state.turn === topSeat}
              deadline={state.actionDeadline}
            />
            <div className="pkr-board">
              <div className="pkr-pot">
                <small>{phaseLabel(state.phase)}</small>
                <b>
                  {totalPot.toLocaleString()} <span>chips</span>
                </b>
              </div>
              <div className="pkr-community">
                {community.map((card, index) =>
                  card ? (
                    <PokerCard key={card.code} card={card} compact />
                  ) : (
                    <span className="pkr-card-slot" key={index} />
                  ),
                )}
              </div>
            </div>
            <Seat
              state={state}
              index={bottomSeat}
              active={!handFinished && me >= 0 && state.turn === bottomSeat}
              deadline={state.actionDeadline}
              mine={me >= 0}
              onCard={toggleCard}
              selected={selected}
            />
          </div>
          <div className="pkr-controls">
            {state.phase === "draw" && myTurn && (
              <div className="pkr-draw-controls">
                <span>
                  {selected.length
                    ? `${selected.length} selected for replacement`
                    : "Select cards to replace, or stand pat"}
                </span>
                <button
                  className="is-primary"
                  type="button"
                  onClick={() => props.onDiscard(selected)}
                >
                  {selected.length ? `Draw ${selected.length}` : "Stand pat"}
                </button>
              </div>
            )}
            {state.legalActions.length > 0 && (
              <div className="pkr-bet-controls">
                {state.legalActions.includes("fold") && (
                  <button type="button" onClick={() => props.onAction("fold")}>
                    Fold
                  </button>
                )}
                {state.legalActions.includes("check") && (
                  <button type="button" onClick={() => props.onAction("check")}>
                    Check
                  </button>
                )}
                {state.legalActions.includes("call") && (
                  <button
                    className="is-primary"
                    type="button"
                    onClick={() => props.onAction("call")}
                  >
                    Call {state.toCall}
                  </button>
                )}
                {state.legalActions.includes("raise") && (
                  <label>
                    <span>
                      Raise to <b>{raiseTo}</b>
                    </span>
                    <input
                      type="range"
                      min={minimumTarget}
                      max={maximumTarget}
                      step={10}
                      value={raiseTo}
                      onChange={(event) =>
                        setRaiseTo(Number(event.target.value))
                      }
                    />
                    <button
                      type="button"
                      onClick={() => props.onAction("raise", raiseTo)}
                    >
                      Raise
                    </button>
                  </label>
                )}
                {state.legalActions.includes("all-in") && (
                  <button
                    type="button"
                    onClick={() => props.onAction("all-in")}
                  >
                    All in · {mySeat?.stack}
                  </button>
                )}
              </div>
            )}
            {(state.phase === "showdown" || state.phase === "hand-finished") &&
              me >= 0 && (
                <button
                  className="pkr-next"
                  type="button"
                  disabled={ready}
                  onClick={props.onNext}
                >
                  <RotateCcw size={15} />{" "}
                  {ready ? "Waiting for opponent" : "Deal next hand"}
                </button>
              )}
          </div>
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

function Seat({
  state,
  index,
  active,
  deadline,
  mine = false,
  onCard,
  selected = [],
}: {
  state: PokerState;
  index: number;
  active: boolean;
  deadline?: string;
  mine?: boolean;
  onCard?: (index: number) => void;
  selected?: number[];
}) {
  const seat = state.seats[index];
  return (
    <article
      className={`pkr-player is-${mine ? "mine" : "opponent"} ${active ? "is-active" : ""} ${seat.folded ? "is-folded" : ""}`}
    >
      <header>
        <i>{seat.name.slice(0, 1).toUpperCase()}</i>
        <span>
          <b>{seat.name}</b>
          <small>
            {index === state.dealer ? "DEALER · " : ""}
            {seat.connected
              ? seat.allIn
                ? "ALL IN"
                : active
                  ? "TO ACT"
                  : "AT TABLE"
              : "RECONNECTING"}
          </small>
        </span>
        <strong>{seat.stack.toLocaleString()}</strong>
        <PokerClock deadline={deadline} active={active} />
      </header>
      <div className="pkr-hand">
        {Array.from({ length: seat.cardCount }, (_, cardIndex) =>
          mine ? (
            <button
              type="button"
              key={seat.cards?.[cardIndex]?.code ?? cardIndex}
              aria-pressed={selected.includes(cardIndex)}
              disabled={state.phase !== "draw" || !active}
              onClick={() => onCard?.(cardIndex)}
            >
              <PokerCard
                card={seat.cards?.[cardIndex]}
                selected={selected.includes(cardIndex)}
              />
            </button>
          ) : (
            <PokerCard
              key={seat.cards?.[cardIndex]?.code ?? cardIndex}
              card={seat.cards?.[cardIndex]}
              hidden={!seat.cards?.[cardIndex]}
              compact
            />
          ),
        )}
      </div>
      {seat.roundBet > 0 && (
        <span className="pkr-bet-chip">{seat.roundBet}</span>
      )}
    </article>
  );
}
function phaseLabel(phase: string) {
  return (
    (
      {
        preflop: "PREFLOP",
        flop: "FLOP",
        turn: "TURN",
        river: "RIVER",
        "draw-bet": "OPENING BET",
        draw: "DRAW",
        "final-bet": "FINAL BET",
        showdown: "SHOWDOWN",
        "hand-finished": "HAND COMPLETE",
      } as Record<string, string>
    )[phase] ?? phase.toUpperCase()
  );
}
function roomHeadline(state: PokerState, me: number) {
  if (state.phase === "showdown")
    return state.winners.includes(me)
      ? `You win with ${state.winningHand}.`
      : state.winners.length > 1
        ? `Split pot · ${state.winningHand}.`
        : `${state.winningHand} wins.`;
  if (state.phase === "hand-finished")
    return state.winners.includes(me)
      ? "The fold is yours."
      : "The pot is conceded.";
  if (me < 0) return "Watching the action.";
  return state.turn === me ? "The action is yours." : "Watch the other seat.";
}
