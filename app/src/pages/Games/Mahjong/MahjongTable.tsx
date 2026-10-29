import { Copy, DoorOpen, Eye, Lightbulb, Palette, RotateCcw, Sparkles } from "lucide-react";
import { useState } from "react";
import { MatchChat } from "../../../components/games/MatchChat";
import type { MahjongState } from "../../../hooks/useMahjongMatch";
import { MahjongBoard } from "./MahjongBoard";
import { MahjongClock } from "./MahjongClock";

interface Props {
  state: MahjongState;
  playerID: string;
  error: string;
  onPair: (first: number, second: number) => void;
  onHint: () => void;
  onChat: (text: string) => void;
  onRematch: () => void;
  onLeave: () => void;
}

const themes = ["jade", "ivory", "midnight", "rosewood"] as const;

export function MahjongTable(props: Props) {
  const { state } = props;
  const me = state.viewerSeat;
  const [settings, setSettings] = useState(false);
  const [theme, setTheme] = useState(
    () => localStorage.getItem("boredgamz-mahjong-theme") ?? "jade",
  );
  const ready = me >= 0 && state.ready.includes(me);
  const cleared = state.seats[me >= 0 ? me : 0]?.remaining ?? 0;

  async function share() {
    const url = new URL(location.href);
    url.search = "";
    url.searchParams.set("watch", state.roomId);
    await navigator.clipboard.writeText(url.toString());
  }

  function chooseTheme(next: string) {
    localStorage.setItem("boredgamz-mahjong-theme", next);
    setTheme(next);
  }

  return (
    <main className="mjh-table-page">
      <header className="mjh-topbar">
        <a href="/">Bored<span>Gamz</span> / TILE GARDEN</a>
        <div><b>{modeTitle(state.mode)}</b><i /> ROOM {state.roomId.slice(-6).toUpperCase()}</div>
        <button type="button" onClick={props.onLeave}><DoorOpen size={15} /> Leave garden</button>
      </header>
      <div className="mjh-room">
        <aside className="mjh-sidebar">
          <span className="mjh-overline"><Sparkles size={12} /> LIVE MAHJONG SOLITAIRE</span>
          <h1>{headline(state, me)}</h1>
          <p>{state.message}</p>
          <MahjongClock deadline={state.deadline} />
          <section className="mjh-players">
            {state.seats.map((seat, index) => (
              <article className={index === me ? "is-me" : ""} key={seat.id}>
                <i>{seat.name.slice(0, 1).toUpperCase()}</i>
                <span><b>{index === me ? "You" : seat.name}</b><small>{seat.connected ? `${seat.pairs} PAIRS` : "RECONNECTING"}</small></span>
                <strong>{seat.remaining}</strong>
              </article>
            ))}
          </section>
          <dl>
            <div><dt>Open tiles</dt><dd>{state.free.length}</dd></div>
            <div><dt>Observers</dt><dd><Eye size={13} /> {state.spectators}</dd></div>
            <div><dt>Reshuffles</dt><dd>{state.reshuffles}</dd></div>
          </dl>
          {state.inviteCode && <p className="mjh-invite">Private code <b>{state.inviteCode}</b></p>}
          <div className="mjh-utility">
            <button type="button" onClick={() => void share()}><Copy size={14} /> Watch link</button>
            <button type="button" onClick={() => setSettings((open) => !open)} aria-expanded={settings}><Palette size={14} /> Tile style</button>
          </div>
          {settings && (
            <section className="mjh-themes" aria-label="Tile style">
              {themes.map((item) => (
                <button type="button" key={item} className={theme === item ? "is-selected" : ""} onClick={() => chooseTheme(item)}>
                  <i className={`theme-${item}`} /> {item}
                </button>
              ))}
            </section>
          )}
          {props.error && <p className="mjh-error" role="alert">{props.error}</p>}
        </aside>
        <section className="mjh-stage">
          <header>
            <div><small>{state.mode === "duel" ? "YOUR MIRRORED LAYOUT" : "SHARED LAYOUT"}</small><b>{cleared} tiles remain</b></div>
            {me >= 0 && state.status === "playing" && (
              <button type="button" disabled={state.hintsRemaining <= 0} onClick={props.onHint}>
                <Lightbulb size={15} /> Hint · {state.hintsRemaining}
              </button>
            )}
          </header>
          <MahjongBoard state={state} theme={theme} onPair={props.onPair} />
          <footer>
            <span>Free tiles have an open left or right edge and no tile above.</span>
            {state.status === "finished" && me >= 0 && (
              <button type="button" disabled={ready} onClick={props.onRematch}>
                <RotateCcw size={15} /> {ready ? "Waiting for partner" : "Play another layout"}
              </button>
            )}
          </footer>
        </section>
        <MatchChat messages={state.chat} playerID={props.playerID} onSend={props.onChat} />
      </div>
    </main>
  );
}

function modeTitle(mode: MahjongState["mode"]) {
  return mode === "duel" ? "MIRROR DUEL" : mode === "zen" ? "ZEN PARTNERS" : "SHARED GARDEN";
}

function headline(state: MahjongState, me: number) {
  if (state.status === "playing") {
    return me < 0 ? "Watching the garden." : state.mode === "duel" ? "Find your rhythm." : "Clear it together.";
  }
  if (state.mode !== "duel") return state.seats[0].remaining === 0 ? "Garden complete." : "The garden rests.";
  if (state.winner < 0) return "A level duel.";
  return state.winner === me ? "You cleared the way." : `${state.seats[state.winner].name} takes the duel.`;
}
