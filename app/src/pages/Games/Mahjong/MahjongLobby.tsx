import { Flower2, LockKeyhole, ShieldCheck, Swords, Users } from "lucide-react";
import { useState } from "react";
import mahjongArt from "../../../assets/mahjong.png";
import type {
  MahjongOptions,
  MahjongState,
} from "../../../hooks/useMahjongMatch";

const modes: Array<{
  id: MahjongState["mode"];
  title: string;
  label: string;
  description: string;
  detail: string;
}> = [
  {
    id: "garden",
    title: "Shared Garden",
    label: "Co-operative",
    description: "Clear one full 144-tile board together before time runs out.",
    detail: "12 minutes · shared board",
  },
  {
    id: "duel",
    title: "Mirror Duel",
    label: "Head-to-head",
    description: "Race through identical private layouts. Every pair matters.",
    detail: "8 minutes · first clear wins",
  },
  {
    id: "zen",
    title: "Zen Partners",
    label: "Untimed",
    description: "Share a compact board with no clock and no pressure.",
    detail: "Untimed · shared board",
  },
];

export function MahjongLobby({
  name,
  onJoin,
}: {
  name: string;
  onJoin: (options: MahjongOptions) => void;
}) {
  const [mode, setMode] = useState<MahjongState["mode"]>("garden");
  const [privateRoom, setPrivateRoom] = useState(false);
  return (
    <main className="mjh-lobby">
      <section className="mjh-lobby__story">
        <a href="/games">← Game collection</a>
        <div className="mjh-lobby__copy">
          <span><Flower2 size={14} /> THE TILE GARDEN</span>
          <h1>Uncover calm.<br /><em>Share the table.</em></h1>
          <p>
            Mahjong Solitaire rebuilt as a live two-player ritual: authentic
            tile families, guaranteed-solvable deals, and every move verified
            by the server.
          </p>
        </div>
        <figure>
          <img src={mahjongArt} alt="Traditional mahjong tiles arranged on a table" />
          <figcaption>144 standard tiles · no downloads · no bots</figcaption>
        </figure>
      </section>
      <form
        className="mjh-lobby__form"
        onSubmit={(event) => {
          event.preventDefault();
          const data = new FormData(event.currentTarget);
          const code = String(data.get("inviteCode") ?? "").trim().toUpperCase();
          onJoin({
            name: String(data.get("name") ?? name),
            mode,
            inviteCode: privateRoom
              ? code || `JADE${Math.random().toString(36).slice(2, 6).toUpperCase()}`
              : undefined,
          });
        }}
      >
        <header>
          <small>CHOOSE A TABLE</small>
          <h2>How would you like to play?</h2>
        </header>
        <fieldset>
          <legend>Game mode</legend>
          {modes.map((item) => (
            <button
              type="button"
              key={item.id}
              className={mode === item.id ? "is-selected" : ""}
              aria-pressed={mode === item.id}
              onClick={() => setMode(item.id)}
            >
              <i>{item.id === "duel" ? <Swords /> : <Users />}</i>
              <span><b>{item.title}</b><small>{item.label}</small></span>
              <p>{item.description}</p>
              <em>{item.detail}</em>
            </button>
          ))}
        </fieldset>
        <label className="mjh-field">
          <span>Display name</span>
          <input name="name" defaultValue={name} maxLength={24} autoComplete="nickname" />
        </label>
        <label className="mjh-private">
          <input type="checkbox" checked={privateRoom} onChange={(event) => setPrivateRoom(event.target.checked)} />
          <LockKeyhole size={17} />
          <span><b>Private garden</b><small>Match only with the same invitation code</small></span>
        </label>
        {privateRoom && (
          <label className="mjh-field">
            <span>Garden code</span>
            <input name="inviteCode" maxLength={12} placeholder="BAMBOO7" />
          </label>
        )}
        <button className="mjh-join" type="submit"><Flower2 size={18} /> Enter the tile garden</button>
        <footer><ShieldCheck size={14} /> Layouts and legal pairs are generated and checked on the server.</footer>
      </form>
    </main>
  );
}
