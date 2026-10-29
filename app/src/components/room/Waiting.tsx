import { Check, Link2 } from "lucide-react";
import { useState } from "react";
import { BoardArt } from "../../games/BoardArt";
import type { GameDefinition } from "../../games/catalog";
import type { JoinRequest } from "../../hooks/useGameSocket";

interface Props {
  game: GameDefinition;
  request: JoinRequest | null;
  connecting: boolean;
  error?: string;
  onCancel: () => void;
}

export function Waiting({ game, request, connecting, error, onCancel }: Props) {
  const [copied, setCopied] = useState(false);
  const code = request?.inviteCode;
  const watching = request?.role === "spectator";
  const title = watching ? "Opening the table…" : connecting ? "Pulling up a chair…" : code ? "Waiting for your friend" : "Finding an opponent";

  async function copyInvite() {
    if (!code) return;
    const url = `${window.location.origin}/games/${game.slug}?join=${code}`;
    try {
      await navigator.clipboard.writeText(url);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
    } catch {
      window.prompt("Copy this invite link", url);
    }
  }

  return (
    <section className="waiting" aria-live="polite">
      <div className="waiting__card">
        <div className="waiting__art"><BoardArt slug={game.slug} /></div>
        <span className="pulse-dots" aria-hidden="true"><i /><i /><i /></span>
        <h1>{title}</h1>
        {code && !watching && (
          <div className="code-card">
            <span className="label">Table code</span>
            <strong>{code}</strong>
            <button className="btn btn--ink" type="button" onClick={() => void copyInvite()}>
              {copied ? <><Check size={16} /> Link copied</> : <><Link2 size={16} /> Copy invite link</>}
            </button>
          </div>
        )}
        {error && <p className="alert" role="alert">{error}</p>}
        <button className="btn btn--ghost" type="button" onClick={onCancel}>Cancel</button>
      </div>
    </section>
  );
}
