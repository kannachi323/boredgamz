import { useState } from 'react';

import { Message } from '../types';

export function ChatBox({ username }: { username: string }) {
  const [messages, setMessages] = useState<Message[]>([])

  return (
    <div className="w-full h-full bg-[#363430] flex flex-col">
      <div className="flex-1 overflow-y-auto p-2 text-white">
        {messages.map((message, idx) => 
          <p key={idx}>{message.sender}: {message.content}</p>
        )}
      </div>

      <div className="w-full border-t border-[#7f7c7b] text-white px-2 py-1">
        <textarea
          placeholder="Type a message here..."
          className="w-full resize-none overflow-hidden bg-transparent text-white outline-none max-h-40"
          rows={1}
          onInput={(e) => {
            e.currentTarget.style.height = "auto";
            e.currentTarget.style.height = `${e.currentTarget.scrollHeight}px`;
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              const content = e.currentTarget.value.trim();
              if (content !== "") {
                setMessages(prev => [...prev, { content, sender: username || "Anonymous" }]);
                e.currentTarget.value = "";
                e.currentTarget.style.height = "auto";
              }
            }
          }}
        />
      </div>
    </div>
  );
}
