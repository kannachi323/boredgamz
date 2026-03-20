import { useState } from 'react';
import { useGomokuStore } from '@/stores/Gomoku/useGomokuStore';

export function ChatBox() {
  const [draft, setDraft] = useState('')
  const { messages, send, conn, player } = useGomokuStore()

  const submitMessage = () => {
    const content = draft.trim()
    if (!content || !conn || conn.readyState !== WebSocket.OPEN) {
      return
    }

    send({
      type: "chat",
      data: { content },
    })
    setDraft('')
  }

  return (
    <div className="w-full h-full bg-[#363430] flex flex-col">
      <div className="flex-1 overflow-y-auto p-2 text-white">
        {messages.map((message, idx) => 
          <p key={`${message.sentAt}-${idx}`}>
            <span className={message.senderID === player.playerID ? "text-[#C3B299] font-semibold" : "text-white font-semibold"}>
              {message.senderName}
            </span>
            : {message.content}
          </p>
        )}
      </div>

      <div className="w-full border-t border-[#7f7c7b] text-white px-2 py-1">
        <textarea
          placeholder="Type a message here..."
          className="w-full resize-none overflow-hidden bg-transparent text-white outline-none max-h-40"
          rows={1}
          value={draft}
          onChange={(e) => setDraft(e.currentTarget.value)}
          onInput={(e) => {
            e.currentTarget.style.height = "auto";
            e.currentTarget.style.height = `${e.currentTarget.scrollHeight}px`;
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              submitMessage();
              e.currentTarget.style.height = "auto";
            }
          }}
        />
      </div>
    </div>
  );
}
