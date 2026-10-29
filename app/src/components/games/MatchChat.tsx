import { SendHorizontal } from 'lucide-react'
import { useEffect, useRef } from 'react'
import type { ArenaChatMessage } from '../../hooks/useOnlineMatch'

interface Props {
  messages: ArenaChatMessage[]
  playerID: string
  onSend: (content: string) => void
}

export function MatchChat({ messages, playerID, onSend }: Props) {
  const endRef = useRef<HTMLDivElement | null>(null)
  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
  }, [messages.length])

  return (
    <section className="match-chat">
      <header><span>Table talk</span><small>{messages.length ? `${messages.length} messages` : 'Be the first to say hello'}</small></header>
      <div className="match-chat__messages" aria-live="polite">
        {messages.length === 0 && <p className="match-chat__empty">A little sportsmanship goes a long way.</p>}
        {messages.map((message) => (
          <article className={message.senderId === playerID ? 'is-mine' : ''} key={message.id}>
            <div><b>{message.senderId === playerID ? 'You' : message.senderName}</b><time>{formatTime(message.createdAt)}</time></div>
            <p>{message.content}</p>
          </article>
        ))}
        <div ref={endRef} />
      </div>
      <form onSubmit={(event) => {
        event.preventDefault()
        const field = event.currentTarget.elements.namedItem('message') as HTMLInputElement
        if (!field.value.trim()) return
        onSend(field.value)
        field.value = ''
      }}>
        <label className="sr-only" htmlFor="table-message">Message</label>
        <input id="table-message" name="message" maxLength={240} placeholder="Write a message…" autoComplete="off" />
        <button type="submit" aria-label="Send message"><SendHorizontal size={17} /></button>
      </form>
    </section>
  )
}

function formatTime(value: string) {
  return new Intl.DateTimeFormat(undefined, { hour: 'numeric', minute: '2-digit' }).format(new Date(value))
}
