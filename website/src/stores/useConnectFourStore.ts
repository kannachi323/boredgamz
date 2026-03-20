import { create } from "zustand"

import {
  type ChatMessage,
  type ClientRequest,
  type GameState,
  type LobbyRequest,
  type Player,
  type ServerResponse,
} from "@/pages/Games/ConnectFour/types"

type ConnectFourConnectionStatus = "idle" | "connecting" | "queued" | "connected" | "error"

function createPlayer(): Player {
  return {
    playerID: "",
    playerName: "",
    color: "red",
    playerClock: { remaining: 0 },
  }
}

function createLobbyRequest(): LobbyRequest {
  return {
    type: "lobby",
    data: {
      timeControl: "Rapid",
      playerID: "",
      playerColor: "random",
      playerName: "",
    },
  }
}

interface ConnectFourStore {
  lobbyRequest: LobbyRequest
  gameState: GameState | null
  conn: WebSocket | null
  connectionStatus: ConnectFourConnectionStatus
  connectionError: string | null
  player: Player
  opponent: Player
  messages: ChatMessage[]

  setLobbyRequest: (lobbyRequest: LobbyRequest) => void
  setPlayer: (player: Player) => void
  setConnection: (lobbyRequest: LobbyRequest) => void
  closeConnection: () => void
  clearSession: () => void
  send: (req: ClientRequest) => void
  handler: (payload: ServerResponse) => void
  refreshPlayers: () => void
}

export const useConnectFourStore = create<ConnectFourStore>((set, get) => ({
  lobbyRequest: createLobbyRequest(),
  gameState: null,
  conn: null,
  connectionStatus: "idle",
  connectionError: null,
  player: createPlayer(),
  opponent: createPlayer(),
  messages: [],

  setLobbyRequest: (lobbyRequest) => set({ lobbyRequest }),
  setPlayer: (player) => set({ player }),

  setConnection: (lobbyRequest) => {
    const oldConn = get().conn
    if (oldConn && oldConn.readyState !== WebSocket.CLOSED) {
      oldConn.close(1000, "Switching session")
    }

    const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}/join-connectfour-lobby`)
    set({
      conn: socket,
      lobbyRequest,
      connectionStatus: "connecting",
      connectionError: null,
      gameState: null,
      messages: [],
    })

    socket.onopen = () => {
      if (get().conn !== socket) return
      set({ connectionStatus: "queued" })
      socket.send(JSON.stringify(lobbyRequest))
    }

    socket.onmessage = (event) => {
      if (get().conn !== socket) return
      const payload = JSON.parse(event.data) as ServerResponse
      get().handler(payload)
    }

    socket.onerror = () => {
      if (get().conn !== socket) return
      set({
        connectionStatus: "error",
        connectionError: "Could not connect to the Connect Four lobby.",
      })
    }

    socket.onclose = (event) => {
      if (get().conn !== socket) return
      const gameFinished = get().gameState?.status.code === "offline"
      const userClosed = event.code === 1000 && event.reason === "User cancelled"

      set({
        conn: null,
        connectionStatus: gameFinished || userClosed ? "idle" : "error",
        connectionError: gameFinished || userClosed ? null : event.reason || "Connection closed unexpectedly.",
      })
    }
  },

  closeConnection: () => {
    const conn = get().conn
    if (!conn) return

    if (conn.readyState === WebSocket.OPEN || conn.readyState === WebSocket.CONNECTING) {
      conn.close(1000, "User cancelled")
    }

    set({ conn: null, connectionStatus: "idle", connectionError: null })
  },

  clearSession: () => {
    get().closeConnection()
    set({
      gameState: null,
      messages: [],
      opponent: createPlayer(),
    })
  },

  send: (req) => {
    const socket = get().conn
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      set({
        connectionStatus: "error",
        connectionError: "Connection is not open.",
      })
      return
    }

    socket.send(JSON.stringify(req))
  },

  handler: (payload) => {
    switch (payload.type) {
      case "update":
        set({
          gameState: payload.data,
          messages: payload.data.messages || [],
          connectionStatus: "connected",
          connectionError: null,
        })
        get().refreshPlayers()
        break
      case "chat":
        set((state) => {
          const nextMessages = [...state.messages, payload.data]
          return {
            messages: nextMessages,
            gameState: state.gameState ? { ...state.gameState, messages: nextMessages } : state.gameState,
          }
        })
        break
      case "error":
        set({
          connectionStatus: "error",
          connectionError: payload.data.message,
        })
        break
    }
  },

  refreshPlayers: () => {
    const { gameState, player } = get()
    if (!gameState) return

    const p1 = gameState.players[0]
    const p2 = gameState.players[1]
    const current = p1.playerID === player.playerID ? p1 : p2
    const opponent = p1.playerID === player.playerID ? p2 : p1

    set({ player: current, opponent })
  },
}))
