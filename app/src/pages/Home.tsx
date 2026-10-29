import { ArrowRight } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { BoardArt } from "../games/BoardArt";
import { gameForCode, games, normalizeCode } from "../games/catalog";
import { getApiRoot } from "../lib/api";

export default function Home() {
  const live = useLiveTables();
  return (
    <div className="lobby">
      <header className="lobby__head">
        <div>
          <h1>Pull up a chair.</h1>
          <p>Six classic games, live with a friend or a stranger.</p>
        </div>
        <JoinCode />
      </header>
      <div className="tiles">
        {games.map((game, index) => (
          <Link key={game.slug} to={`/games/${game.slug}`} className="tile" style={{ "--i": index } as React.CSSProperties} aria-label={`Play ${game.name}`}>
            <div className="tile__art"><BoardArt slug={game.slug} /></div>
            <div className="tile__foot">
              <div className="tile__name">
                <h2>{game.name}</h2>
                <span>2 players · {game.minutes}</span>
              </div>
              <span className="tile__play" aria-hidden="true"><ArrowRight size={18} /></span>
            </div>
            {index === 0 && live > 0 && <span className="chip tile__live"><span className="dot" /> {live} live</span>}
          </Link>
        ))}
      </div>
    </div>
  );
}

function JoinCode() {
  const navigate = useNavigate();
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        const game = gameForCode(code);
        if (!game) {
          setError("That code doesn't look right.");
          return;
        }
        navigate(`/games/${game.slug}?join=${normalizeCode(code)}`);
      }}
    >
      <label className="sr-only" htmlFor="join-code">Private table code</label>
      <div className="join-code">
        <input id="join-code" className="field field--code" placeholder="Have a code?" value={code} maxLength={6} autoComplete="off" spellCheck={false} onChange={(event) => { setCode(normalizeCode(event.target.value)); setError(""); }} aria-invalid={Boolean(error)} aria-describedby="join-code-hint" />
        <button className="btn btn--ink" type="submit" disabled={code.length < 6}>Join</button>
      </div>
      <p id="join-code-hint" className="join-code__hint" role={error ? "alert" : undefined}>{error || " "}</p>
    </form>
  );
}

/** Total live tables across all games, refreshed every 20 seconds. */
function useLiveTables() {
  const [live, setLive] = useState(0);
  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      try {
        const response = await fetch(`${getApiRoot()}/v1/status`);
        if (!response.ok) return;
        const body = (await response.json()) as Record<string, { rooms?: number }>;
        const total = ["arena", "poker", "mahjong"].reduce((sum, key) => sum + (body[key]?.rooms ?? 0), 0);
        if (!cancelled) setLive(total);
      } catch {
        // Status is decorative; ignore failures.
      }
    };
    void load();
    const timer = window.setInterval(load, 20_000);
    return () => {
      cancelled = true;
      window.clearInterval(timer);
    };
  }, []);
  return live;
}
