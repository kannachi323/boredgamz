export interface Stone {
  color: string
}

export interface Board {
  rows: number
  cols: number
  stones: (Stone | null)[][]
}

export interface Move {
  row: number
  col: number
  color: string
}

export interface PlayerClock {
  remaining: number
}

export interface Player {
  playerID: string
  playerName: string
  color: string
  playerClock: PlayerClock | null
}

export interface ChatMessage {
  senderID: string
  senderName: string
  content: string
  sentAt: string
}

export interface GameStatus {
  result: string
  code: "online" | "offline"
  winner: Player | null
}

export interface GameState {
  gameID: string
  board: Board
  players: Player[]
  status: GameStatus
  lastMove: Move | null
  turn: string
  moves: Move[]
  messages: ChatMessage[]
}

export interface LobbyRequest {
  type: "lobby"
  data: {
    timeControl: string
    playerID: string
    playerColor: "red" | "yellow" | "random"
    playerName: string
  }
}

export interface MoveRequest {
  type: "move"
  data: {
    column: number
  }
}

export interface ChatRequest {
  type: "chat"
  data: {
    content: string
  }
}

export type ClientRequest = LobbyRequest | MoveRequest | ChatRequest

export type ServerResponse =
  | { type: "update"; data: GameState }
  | { type: "chat"; data: ChatMessage }
  | { type: "error"; data: { message: string } }
