export interface Conn {
  roomID: string;
  player: Player
  type: string
}

export interface Board {
  stones: Stone[][];
  size: number;
  numStones: number;
}

export interface Move {
  row: number;
  col: number;
  color: string;
}

export interface GameState {
  gameID: string;
  board: Board;
  size: number;
  openingRule: "freestyle" | "standard" | "renju";
  swapRuleEnabled: boolean;
  firstMoveCenterEnabled: boolean;
  players: Player[];
  turn: string;
  status: GameStatus;
  lastMove: Move | null;
  moves: Move[];
  messages: ChatMessage[];
}

export interface AnalysisState {
  moves: Move[];
  board: Board | null
  active: boolean
  index: number
}

export interface GameStatus {
  result: "win" | "draw" | "loss";
  code: "online" | "offline";
  winner: Player | null;
}

export interface User {
  id: string;
  username: string;
}

export interface Player {
  playerID: string;
  playerName: string;
  color: string;
  playerClock: PlayerClock | null
}

export interface PlayerClock {
  remaining: number;
}

export interface Stone {
  color: string | null;
}

export interface ChatMessage {
  senderID: string
  senderName: string
  content: string
  sentAt: string
}

export interface LobbyRequest {
  type: "lobby"
  data: {
    name: string
    timeControl: string;
    mode: string 
    openingRule: "freestyle" | "standard" | "renju"
    swapRuleEnabled: boolean
    firstMoveCenterEnabled: boolean
    botDifficulty?: "beginner" | "intermediate" | "advanced"
    playerID: string;
    playerName: string;
    playerColor: string;
  }
}

export interface ReconnectRequest {
  type: "reconnect"
  data: {
    lobbyID: string
    playerID: string
  }
}

export interface ErrorData {
  message: string
}

export interface MoveRequest {
  type: "move"
  data: {
    move: Move
  }
}

export interface SwapRequest {
  type: "swap"
  data: Record<string, never>
}

export interface ChatRequest {
  type: "chat"
  data: {
    content: string
  }
}

export type ClientRequest = 
  | MoveRequest
  | SwapRequest
  | ChatRequest
  | LobbyRequest
  | ReconnectRequest


export type ServerResponse =
  | { type: "update"; data: GameState }
  | { type: "error"; data: ErrorData }
  | { type: "chat"; data: ChatMessage }
  | { type: string; data: unknown }



//MODELS

export interface PlayerRow {
  playerID: string;
  playerName: string;
  color: string;
}

export interface GameStateRow {
  gameID: string;
  boardSize: number;
  players: PlayerRow[]
  moves: Move[];
  messages: ChatMessage[];
  result: "win" | "draw" | "loss";
  winner: PlayerRow | null;
}


