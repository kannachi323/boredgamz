import { SendHorizontal } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import type { ChatMessage } from "../../hooks/useGameSocket";

interface Props {
  messages: ChatMessage[];
  viewerId?: string;
  onSend: (content: string) => void;
}

export function MatchChat({ messages, viewerId, onSend }: Props) {
  const listRef = useRef<HTMLOListElement>(null);
  const [draft, setDraft] = useState("");
  useEffect(() => {
    const list = listRef.current;
    if (list) list.scrollTop = list.scrollHeight;
  }, [messages.length]);

  return (
    <section className="chat" aria-label="Table chat">
      <ol className="chat__list" ref={listRef} aria-live="polite">
        {messages.length === 0 && <li className="chat__empty">Say hi.</li>}
        {messages.map((message) => {
          const mine = message.senderId === viewerId;
          return (
            <li key={message.id} className={mine ? "is-mine" : undefined}>
              {!mine && <b>{message.senderName}</b>}
              <p>{message.content}</p>
            </li>
          );
        })}
      </ol>
      <form
        className="chat__form"
        onSubmit={(event) => {
          event.preventDefault();
          if (!draft.trim()) return;
          onSend(draft);
          setDraft("");
        }}
      >
        <label className="sr-only" htmlFor="chat-input">Message</label>
        <input id="chat-input" className="field" value={draft} onChange={(event) => setDraft(event.target.value)} maxLength={240} placeholder="Message" autoComplete="off" />
        <button className="btn btn--ink btn--icon" type="submit" aria-label="Send" disabled={!draft.trim()}><SendHorizontal size={17} /></button>
      </form>
    </section>
  );
}
