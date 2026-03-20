import { useEffect, useState } from "react"
import { RouteObject } from "react-router-dom"

import { useAuthStore } from "@/stores/useAuthStore"
import { useConnectFourStore } from "@/stores/useConnectFourStore"

function formatClock(milliseconds: number | undefined) {
  if (!milliseconds || milliseconds < 0) {
    return "00:00"
  }

  const totalSeconds = Math.floor(milliseconds / 1_000_000_000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`
}

function ConnectFourBoard() {
  const { gameState, player, send } = useConnectFourStore()

  if (!gameState) return null

  const isMyTurn = gameState.turn === player.playerID && gameState.status.code === "online"

  return (
    <div className="w-full max-w-3xl rounded-[2rem] bg-[#154c79] p-4 shadow-[0_20px_60px_rgba(0,0,0,0.35)]">
      <div className="grid grid-cols-7 gap-3">
        {Array.from({ length: gameState.board.cols }).map((_, colIdx) => (
          <button
            key={colIdx}
            type="button"
            disabled={!isMyTurn || gameState.board.stones[0][colIdx] !== null}
            onClick={() => send({ type: "move", data: { column: colIdx } })}
            className="space-y-3 rounded-3xl bg-[#0f3b5f] p-3 transition enabled:hover:-translate-y-1 disabled:cursor-not-allowed disabled:opacity-70"
          >
            {gameState.board.stones.map((row, rowIdx) => {
              const stone = row[colIdx]
              const color = stone?.color === "red" ? "bg-[#f55555]" : stone?.color === "yellow" ? "bg-[#ffd84d]" : "bg-[#e8f0f6]"
              return <div key={`${rowIdx}-${colIdx}`} className={`aspect-square w-full rounded-full border border-black/10 ${color}`} />
            })}
          </button>
        ))}
      </div>
    </div>
  )
}

function ConnectFourChat() {
  const [draft, setDraft] = useState("")
  const { messages, player, send, conn } = useConnectFourStore()

  const submit = () => {
    const content = draft.trim()
    if (!content || !conn || conn.readyState !== WebSocket.OPEN) return

    send({ type: "chat", data: { content } })
    setDraft("")
  }

  return (
    <div className="flex h-full min-h-0 flex-col rounded-[1.75rem] bg-[#102a43] p-4 text-[#f4f0ea]">
      <div className="mb-3 text-sm font-semibold uppercase tracking-[0.25em] text-[#9fb3c8]">Table Chat</div>
      <div className="flex-1 space-y-2 overflow-y-auto rounded-2xl bg-[#0b1f33] p-3">
        {messages.length === 0 ? (
          <p className="text-sm text-[#9fb3c8]">Messages will show up here once the match starts.</p>
        ) : (
          messages.map((message, idx) => (
            <p key={`${message.sentAt}-${idx}`} className="text-sm leading-6">
              <span className={message.senderID === player.playerID ? "font-semibold text-[#ffd84d]" : "font-semibold text-[#f4f0ea]"}>
                {message.senderName}
              </span>
              : {message.content}
            </p>
          ))
        )}
      </div>
      <textarea
        value={draft}
        onChange={(e) => setDraft(e.currentTarget.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault()
            submit()
          }
        }}
        rows={2}
        placeholder="Send a message"
        className="mt-3 resize-none rounded-2xl border border-[#486581] bg-[#0b1f33] px-3 py-2 text-sm outline-none"
      />
    </div>
  )
}

function ConnectFourLobby() {
  const { user } = useAuthStore()
  const { lobbyRequest, setLobbyRequest, setConnection, connectionStatus, connectionError } = useConnectFourStore()
  const canJoin = Boolean(user) && connectionStatus !== "connecting" && connectionStatus !== "queued"

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6 rounded-[2rem] bg-[#102a43] p-8 text-[#f4f0ea] shadow-[0_24px_80px_rgba(0,0,0,0.35)]">
      <div>
        <p className="text-sm uppercase tracking-[0.3em] text-[#9fb3c8]">Connect Four</p>
        <h1 className="mt-2 text-4xl font-black">Queue a live match</h1>
        <p className="mt-3 text-[#bcccdc]">Pick a color, join the lobby, and the same websocket carries gameplay and chat once you’re matched.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <label className="rounded-2xl bg-[#0b1f33] p-4">
          <div className="mb-2 text-sm font-semibold text-[#9fb3c8]">Color</div>
          <select
            value={lobbyRequest.data.playerColor}
            onChange={(e) =>
              setLobbyRequest({
                ...lobbyRequest,
                data: {
                  ...lobbyRequest.data,
                  playerColor: e.target.value as "red" | "yellow" | "random",
                },
              })
            }
            className="w-full rounded-xl border border-[#486581] bg-[#102a43] px-3 py-2 outline-none"
          >
            <option value="random">Random</option>
            <option value="red">Red</option>
            <option value="yellow">Yellow</option>
          </select>
        </label>

        <label className="rounded-2xl bg-[#0b1f33] p-4">
          <div className="mb-2 text-sm font-semibold text-[#9fb3c8]">Time Control</div>
          <select
            value={lobbyRequest.data.timeControl}
            onChange={(e) =>
              setLobbyRequest({
                ...lobbyRequest,
                data: {
                  ...lobbyRequest.data,
                  timeControl: e.target.value,
                },
              })
            }
            className="w-full rounded-xl border border-[#486581] bg-[#102a43] px-3 py-2 outline-none"
          >
            <option value="Rapid">Rapid</option>
            <option value="Blitz">Blitz</option>
            <option value="Bullet">Bullet</option>
          </select>
        </label>
      </div>

      {connectionError && <p className="rounded-2xl bg-[#7b341e] px-4 py-3 text-sm">{connectionError}</p>}
      {!user && <p className="rounded-2xl bg-[#486581] px-4 py-3 text-sm">Loading your player session before joining the queue.</p>}

      <button
        type="button"
        disabled={!canJoin}
        onClick={() => {
          if (!user) return
          setConnection({
            ...lobbyRequest,
            data: {
              ...lobbyRequest.data,
              playerID: user.id,
              playerName: user.username,
            },
          })
        }}
        className="rounded-2xl bg-[#ffd84d] px-5 py-3 text-base font-bold text-[#102a43] transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-70"
      >
        {connectionStatus === "connecting" || connectionStatus === "queued" ? "Joining Queue..." : "Join Queue"}
      </button>
    </div>
  )
}

export function ConnectFour() {
  const { checkAuth, isAuthenticated } = useAuthStore()
  const { setPlayer, player, gameState, opponent, connectionStatus, clearSession } = useConnectFourStore()

  useEffect(() => {
    const check = async () => {
      const success = await checkAuth(() => {})
      if (!success) return

      const user = useAuthStore.getState().user
      if (!user) return

      setPlayer({ ...player, playerID: user.id, playerName: user.username })
    }

    check()
  }, [])

  useEffect(() => {
    return () => clearSession()
  }, [clearSession])

  if (!isAuthenticated) {
    return <div className="flex h-screen items-center justify-center"></div>
  }

  const statusText =
    connectionStatus === "queued"
      ? "Waiting for an opponent..."
      : gameState?.status.code === "offline"
        ? gameState.status.winner
          ? `${gameState.status.winner.playerName} wins`
          : "Draw"
        : gameState?.turn === player.playerID
          ? "Your turn"
          : opponent.playerName
            ? `${opponent.playerName}'s turn`
            : "Match in progress"

  return (
    <div className="min-h-[90vh] bg-[radial-gradient(circle_at_top,#1f6f8b_0%,#0b1f33_50%,#08131f_100%)] px-4 py-10 text-[#f4f0ea]">
      {!gameState ? (
        <ConnectFourLobby />
      ) : (
        <div className="mx-auto grid max-w-7xl gap-6 lg:grid-cols-[280px_minmax(0,1fr)_320px]">
          <section className="rounded-[1.75rem] bg-[#102a43] p-5">
            <p className="text-sm uppercase tracking-[0.25em] text-[#9fb3c8]">Players</p>
            <div className="mt-4 space-y-4">
              <div className="rounded-2xl bg-[#0b1f33] p-4">
                <p className="text-xs uppercase tracking-[0.2em] text-[#9fb3c8]">You</p>
                <p className="mt-1 text-xl font-bold">{player.playerName}</p>
                <p className="mt-1 text-sm capitalize text-[#bcccdc]">{player.color}</p>
                <p className="mt-3 text-sm text-[#ffd84d]">{formatClock(player.playerClock?.remaining)}</p>
              </div>
              <div className="rounded-2xl bg-[#0b1f33] p-4">
                <p className="text-xs uppercase tracking-[0.2em] text-[#9fb3c8]">Opponent</p>
                <p className="mt-1 text-xl font-bold">{opponent.playerName || "Waiting..."}</p>
                <p className="mt-1 text-sm capitalize text-[#bcccdc]">{opponent.color || "Unknown"}</p>
                <p className="mt-3 text-sm text-[#ffd84d]">{formatClock(opponent.playerClock?.remaining)}</p>
              </div>
            </div>
            <div className="mt-6 rounded-2xl bg-[#0b1f33] p-4">
              <p className="text-xs uppercase tracking-[0.2em] text-[#9fb3c8]">Status</p>
              <p className="mt-2 text-lg font-semibold">{statusText}</p>
            </div>
          </section>

          <section className="flex flex-col items-center justify-center gap-5 rounded-[1.75rem] bg-[#102a43] p-5">
            <ConnectFourBoard />
            {gameState.status.code === "offline" && (
              <button
                type="button"
                onClick={() => clearSession()}
                className="rounded-2xl bg-[#ffd84d] px-5 py-3 font-bold text-[#102a43]"
              >
                Back to Lobby
              </button>
            )}
          </section>

          <ConnectFourChat />
        </div>
      )}
    </div>
  )
}

export default function ConnectFourRoutes(): RouteObject {
  return {
    path: "/games/connectfour",
    element: <ConnectFour />,
  }
}
