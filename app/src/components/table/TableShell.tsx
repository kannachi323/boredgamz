import { ArrowLeft, Check, Eye, Link2, MessageCircle, Palette, X } from "lucide-react";
import { useEffect, useState, type ReactNode } from "react";
import type { GameDefinition } from "../../games/catalog";
import type { ChatMessage } from "../../hooks/useGameSocket";
import { useShell } from "../../stores/useShell";
import { MatchChat } from "./MatchChat";

interface Props {
  game: GameDefinition;
  meta: string[];
  roomId: string;
  spectators: number;
  inviteCode?: string;
  reconnecting?: boolean;
  error?: string;
  rail: ReactNode;
  settings?: ReactNode;
  chat: { messages: ChatMessage[]; viewerId?: string; onSend: (text: string) => void };
  onLeave: () => void;
  children: ReactNode;
  className?: string;
}

export function TableShell(props: Props) {
  const setImmersive = useShell((shell) => shell.setImmersive);
  const [chatOpen, setChatOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const [seen, setSeen] = useState(props.chat.messages.length);
  const unread = chatOpen ? 0 : Math.max(0, props.chat.messages.length - seen);

  useEffect(() => {
    setImmersive(true);
    return () => setImmersive(false);
  }, [setImmersive]);
  useEffect(() => {
    if (chatOpen) setSeen(props.chat.messages.length);
  }, [chatOpen, props.chat.messages.length]);

  async function copyWatchLink() {
    const url = `${window.location.origin}/games/${props.game.slug}?watch=${props.roomId}`;
    try {
      await navigator.clipboard.writeText(url);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
    } catch {
      window.prompt("Copy this spectator link", url);
    }
  }

  return (
    <div className={`table ${props.className ?? ""}`}>
      <header className="table-bar">
        <button className="btn btn--ghost" type="button" onClick={props.onLeave}>
          <ArrowLeft size={17} /> <span className="hide-sm">Leave</span>
        </button>
        <div className="table-bar__title">
          <b>{props.game.name}</b>
          <span className="table-bar__meta">{props.meta.filter(Boolean).join(" · ")}</span>
        </div>
        <span className="topbar__spacer" />
        {props.inviteCode && <span className="chip hide-sm" title="Private table code">{props.inviteCode}</span>}
        <span className={`chip ${props.spectators ? "" : "hide-sm"}`} title="Spectators"><Eye size={14} /> {props.spectators}</span>
        <button className="btn btn--ghost btn--icon" type="button" onClick={() => void copyWatchLink()} aria-label="Copy spectator link" title="Copy spectator link">
          {copied ? <Check size={17} /> : <Link2 size={17} />}
        </button>
        {props.settings && (
          <div className="menu">
            <button className="btn btn--ghost btn--icon" type="button" aria-label="Board style" aria-expanded={settingsOpen} title="Board style" onClick={() => setSettingsOpen((open) => !open)}>
              <Palette size={17} />
            </button>
            {settingsOpen && <div className="menu__panel table-settings">{props.settings}</div>}
          </div>
        )}
        <button className="btn btn--ghost btn--icon chat-toggle" type="button" aria-label={chatOpen ? "Close chat" : "Open chat"} aria-expanded={chatOpen} onClick={() => setChatOpen((open) => !open)}>
          {chatOpen ? <X size={18} /> : <MessageCircle size={18} />}
          {unread > 0 && <span className="badge">{unread}</span>}
        </button>
      </header>
      {props.reconnecting && <div className="table-banner" role="status"><span className="spinner" /> Reconnecting…</div>}
      <div className="table-body">
        <aside className="table-rail">{props.rail}</aside>
        <section className="table-stage">
          {props.children}
          {props.error && <p className="alert table-error" role="alert">{props.error}</p>}
        </section>
        <aside className={`table-chat ${chatOpen ? "is-open" : ""}`}>
          <MatchChat {...props.chat} />
        </aside>
      </div>
    </div>
  );
}
