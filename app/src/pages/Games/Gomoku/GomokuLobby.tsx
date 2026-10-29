import {
  CircleDot,
  LockKeyhole,
  Ruler,
  ShieldCheck,
  Timer,
} from "lucide-react";
import { useState } from "react";
import type { MatchOptions } from "../../../hooks/useOnlineMatch";
import { boardImages, type BoardTheme } from "./boardThemes";

export function GomokuLobby({
  name,
  onJoin,
}: {
  name: string;
  onJoin: (options: MatchOptions, theme: BoardTheme) => void;
}) {
  const [mode, setMode] = useState("freestyle");
  const [theme, setTheme] = useState<BoardTheme>("mid");
  const [bestOf, setBestOf] = useState(1);
  const [timeControl, setTimeControl] = useState("60s");
  const [privateTable, setPrivateTable] = useState(false);
  return (
    <main className="gmk-lobby">
      <header>
        <a href="/games">← Game room</a>
        <span>Traditional strategy / live table</span>
        <h1>Gomoku</h1>
        <p>
          Build quiet pressure across a wooden grid. The first exact or unbroken
          line of five decides the table.
        </p>
      </header>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          const data = new FormData(event.currentTarget);
          const code = String(data.get("inviteCode") ?? "")
            .trim()
            .toUpperCase();
          onJoin(
            {
              name: String(data.get("name") ?? name),
              mode,
              bestOf,
              timeControl,
              inviteCode: privateTable
                ? code ||
                  `FIVE${Math.random().toString(36).slice(2, 6).toUpperCase()}`
                : undefined,
            },
            theme,
          );
        }}
      >
        <section>
          <div className="gmk-section-title">
            <CircleDot size={16} />
            <span>
              <b>Rules</b>
              <small>Choose how five is judged</small>
            </span>
          </div>
          <div className="gmk-choice-row">
            <button
              type="button"
              className={mode === "freestyle" ? "is-selected" : ""}
              aria-pressed={mode === "freestyle"}
              onClick={() => setMode("freestyle")}
            >
              <b>Freestyle</b>
              <small>15 × 15 · five or more</small>
            </button>
            <button
              type="button"
              className={mode === "pro" ? "is-selected" : ""}
              aria-pressed={mode === "pro"}
              onClick={() => setMode("pro")}
            >
              <b>Exact Five</b>
              <small>15 × 15 · overlines do not win</small>
            </button>
            <button
              type="button"
              className={mode === "compact" ? "is-selected" : ""}
              aria-pressed={mode === "compact"}
              onClick={() => setMode("compact")}
            >
              <b>Compact</b>
              <small>13 × 13 · faster pressure</small>
            </button>
          </div>
        </section>
        <section>
          <div className="gmk-section-title">
            <Ruler size={16} />
            <span>
              <b>Board</b>
              <small>Choose your table grain</small>
            </span>
          </div>
          <div className="gmk-board-themes">
            {(Object.keys(boardImages) as BoardTheme[]).map((item) => (
              <button
                type="button"
                key={item}
                className={theme === item ? "is-selected" : ""}
                aria-pressed={theme === item}
                onClick={() => setTheme(item)}
              >
                <img src={boardImages[item]} alt="" />
                <span>
                  {item === "small"
                    ? "Pale maple"
                    : item === "mid"
                      ? "Club walnut"
                      : "Deep oak"}
                </span>
              </button>
            ))}
          </div>
        </section>
        <div className="gmk-form-row">
          <label>
            <span>Series</span>
            <select
              value={bestOf}
              onChange={(e) => setBestOf(Number(e.target.value))}
            >
              <option value={1}>Single game</option>
              <option value={3}>Best of 3</option>
              <option value={5}>Best of 5</option>
            </select>
          </label>
          <label>
            <span>
              <Timer size={13} /> Move clock
            </span>
            <select
              value={timeControl}
              onChange={(e) => setTimeControl(e.target.value)}
            >
              <option value="untimed">Untimed</option>
              <option value="30s">30 seconds</option>
              <option value="60s">60 seconds</option>
            </select>
          </label>
        </div>
        <label className="gmk-field">
          <span>Display name</span>
          <input name="name" maxLength={24} defaultValue={name} />
        </label>
        <label className="gmk-private">
          <input
            type="checkbox"
            checked={privateTable}
            onChange={(e) => setPrivateTable(e.target.checked)}
          />
          <span>
            <LockKeyhole size={15} /> Private match code
          </span>
        </label>
        {privateTable && (
          <label className="gmk-field">
            <span>Invite code</span>
            <input name="inviteCode" maxLength={12} placeholder="FIVE5" />
          </label>
        )}
        <button className="gmk-play" type="submit">
          <ShieldCheck size={18} /> Find a live table
        </button>
      </form>
    </main>
  );
}
