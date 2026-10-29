import { Club, LockKeyhole, ShieldCheck, Spade } from "lucide-react";
import { useState } from "react";
import pokerArt from "../../../assets/poker.png";
import type { PokerOptions, PokerState } from "../../../hooks/usePokerMatch";

const rooms: Array<{
  id: PokerState["mode"];
  name: string;
  kicker: string;
  rule: string;
}> = [
  {
    id: "holdem",
    name: "No-Limit Hold’em",
    kicker: "Community cards",
    rule: "10 / 20 blinds · 1,000 chips",
  },
  {
    id: "five-card-draw",
    name: "Five-Card Draw",
    kicker: "The original draw game",
    rule: "10 ante · 1,000 chips",
  },
];

export function PokerLobby({
  name,
  onJoin,
}: {
  name: string;
  onJoin: (options: PokerOptions) => void;
}) {
  const [mode, setMode] = useState<PokerState["mode"]>("holdem");
  const [privateRoom, setPrivateRoom] = useState(false);
  return (
    <main className="pkr-lobby">
      <section className="pkr-lobby__intro">
        <a href="/games">← Game collection</a>
        <div>
          <span>
            <Spade size={14} /> THE BOREDGAMZ CARD ROOM
          </span>
          <h1>
            Read the table.
            <br />
            <em>Not the software.</em>
          </h1>
          <p>
            Private hole cards, server-shuffled decks, exact wager validation,
            and a room built for the rhythm of a real heads-up game.
          </p>
        </div>
        <img src={pokerArt} alt="Poker chips and cards on dark green felt" />
      </section>
      <form
        className="pkr-lobby__form"
        onSubmit={(event) => {
          event.preventDefault();
          const data = new FormData(event.currentTarget);
          const code = String(data.get("inviteCode") ?? "")
            .trim()
            .toUpperCase();
          onJoin({
            name: String(data.get("name") ?? name),
            mode,
            inviteCode: privateRoom
              ? code ||
                `DEAL${Math.random().toString(36).slice(2, 6).toUpperCase()}`
              : undefined,
          });
        }}
      >
        <header>
          <small>TAKE A SEAT</small>
          <h2>Choose your game</h2>
          <p>
            Free-play chips. Two human players. No bots, purchases, or cash
            value.
          </p>
        </header>
        <fieldset>
          <legend>Table rules</legend>
          {rooms.map((room) => (
            <button
              type="button"
              key={room.id}
              className={mode === room.id ? "is-selected" : ""}
              aria-pressed={mode === room.id}
              onClick={() => setMode(room.id)}
            >
              <i>{room.id === "holdem" ? "♠" : "♦"}</i>
              <span>
                <b>{room.name}</b>
                <small>{room.kicker}</small>
              </span>
              <em>{room.rule}</em>
            </button>
          ))}
        </fieldset>
        <label className="pkr-field">
          <span>Display name</span>
          <input
            name="name"
            defaultValue={name}
            maxLength={24}
            autoComplete="nickname"
          />
        </label>
        <label className="pkr-private">
          <input
            type="checkbox"
            checked={privateRoom}
            onChange={(event) => setPrivateRoom(event.target.checked)}
          />
          <LockKeyhole size={16} />
          <span>
            <b>Private table</b>
            <small>Only match the same invitation code</small>
          </span>
        </label>
        {privateRoom && (
          <label className="pkr-field">
            <span>Table code</span>
            <input name="inviteCode" maxLength={12} placeholder="RIVER7" />
          </label>
        )}
        <button className="pkr-seat" type="submit">
          <Club size={17} /> Find a heads-up table
        </button>
        <footer>
          <ShieldCheck size={14} />
          <span>
            Cards are shuffled with operating-system cryptographic randomness
            and dealt only by the server.
          </span>
        </footer>
      </form>
    </main>
  );
}
