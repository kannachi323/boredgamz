import { ArrowLeft, Lock, Play } from "lucide-react";
import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { BoardArt, ModeIcon } from "../../games/BoardArt";
import { gameForCode, makeInviteCode, normalizeCode, type GameDefinition } from "../../games/catalog";
import type { JoinRequest } from "../../hooks/useGameSocket";
import { useAuthStore } from "../../stores/useAuthStore";
import { useBoardTheme } from "./useBoardTheme";

interface Props {
  game: GameDefinition;
  error?: string;
  onJoin: (request: JoinRequest) => void;
}

export function GameLobby({ game, error, onJoin }: Props) {
  const signedIn = useAuthStore((auth) => auth.status === "authenticated");
  const navigate = useNavigate();
  const location = useLocation();
  const [mode, setMode] = useState(game.modes[0].id);
  const [bestOf, setBestOf] = useState(game.series?.[1]?.value ?? 1);
  const [clock, setClock] = useState(game.defaultClock ?? "");
  const [theme, setTheme] = useBoardTheme(game);
  const [code, setCode] = useState("");
  const [codeError, setCodeError] = useState("");
  const next = encodeURIComponent(location.pathname + location.search);

  const settings = (): JoinRequest => {
    const selected = game.modes.find((item) => item.id === mode)!;
    return {
      mode,
      bestOf: game.series ? bestOf : undefined,
      timeControl: selected.timeControl ?? (game.clocks ? clock : undefined),
    };
  };

  function joinCode(event: React.FormEvent) {
    event.preventDefault();
    const owner = gameForCode(code);
    if (!owner) {
      setCodeError("That code doesn't look right.");
      return;
    }
    if (owner.slug !== game.slug) {
      navigate(`/games/${owner.slug}?join=${normalizeCode(code)}`);
      return;
    }
    onJoin({ inviteCode: normalizeCode(code) });
  }

  return (
    <div className="room">
      <div className="room__stage">
        <BoardArt slug={game.slug} />
      </div>

      <div className="setup">
        <Link className="back" to="/"><ArrowLeft size={16} /> All games</Link>
        <div className="setup__title">
          <h1>{game.name}</h1>
          <span className="chip">2 players · {game.minutes}</span>
        </div>

        <div className="setup__group" role="group" aria-labelledby="mode-label">
          <span className="label" id="mode-label">Mode</span>
          <div className="modes">
            {game.modes.map((item) => (
              <button key={item.id} type="button" className="mode" aria-label={`${item.name}, ${item.detail}`} aria-pressed={mode === item.id} onClick={() => setMode(item.id)}>
                <span className="mode__glyph"><ModeIcon glyph={item.glyph} /></span>
                <span><b>{item.name}</b><small>{item.detail}</small></span>
                <span className="mode__check" aria-hidden="true" />
              </button>
            ))}
          </div>
        </div>

        {game.series && (
          <Segmented label="Match" options={game.series} value={bestOf} onChange={setBestOf} />
        )}
        {game.clocks && (
          <Segmented label="Move clock" options={game.clocks} value={clock} onChange={setClock} />
        )}
        {game.themes && (
          <div className="setup__group" role="group" aria-labelledby="board-label">
            <span className="label" id="board-label">Board</span>
            <div className="swatches">
              {game.themes.map((item) => (
                <button key={item.id} type="button" className="swatch" aria-pressed={theme === item.id} aria-label={item.label} title={item.label} onClick={() => setTheme(item.id)}>
                  <i style={{ background: `linear-gradient(135deg, ${item.light} 50%, ${item.dark} 50%)` }} />
                </button>
              ))}
            </div>
          </div>
        )}

        {error && <p className="alert" role="alert">{error}</p>}

        {signedIn ? (
          <div className="actions">
            <button className="btn btn--primary btn--lg btn--block" type="button" onClick={() => onJoin(settings())}>
              <Play size={18} fill="currentColor" /> Play
            </button>
            <button className="btn btn--lg btn--block" type="button" onClick={() => onJoin({ ...settings(), inviteCode: makeInviteCode(game) })}>
              <Lock size={16} /> Create private table
            </button>
            <div className="actions__divider">or join a friend</div>
            <form className="actions__row" onSubmit={joinCode}>
              <label className="sr-only" htmlFor="room-code">Private table code</label>
              <input id="room-code" className="field field--code actions__code" placeholder="Code" value={code} maxLength={6} autoComplete="off" inputMode="numeric" pattern="[0-9]*" spellCheck={false} aria-invalid={Boolean(codeError)} onChange={(event) => { setCode(normalizeCode(event.target.value)); setCodeError(""); }} />
              <button className="btn btn--ink" type="submit" disabled={code.length < 6}>Join</button>
            </form>
            {codeError && <p className="alert" role="alert">{codeError}</p>}
          </div>
        ) : (
          <div className="actions">
            <Link className="btn btn--primary btn--lg btn--block" to={`/signup?next=${next}`}>Create an account to play</Link>
            <Link className="btn btn--lg btn--block" to={`/login?next=${next}`}>Sign in</Link>
          </div>
        )}
      </div>
    </div>
  );
}

function Segmented<T extends string | number>({ label, options, value, onChange }: { label: string; options: { value: T; label: string }[]; value: T; onChange: (value: T) => void }) {
  const id = `seg-${label.replace(/\s+/g, "-").toLowerCase()}`;
  return (
    <div className="setup__group" role="group" aria-labelledby={id}>
      <span className="label" id={id}>{label}</span>
      <div className="seg">
        {options.map((option) => (
          <button key={String(option.value)} type="button" aria-pressed={value === option.value} onClick={() => onChange(option.value)}>{option.label}</button>
        ))}
      </div>
    </div>
  );
}
