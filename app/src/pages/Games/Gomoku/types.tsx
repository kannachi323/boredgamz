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
  players: Player[];
  turn: string;
  status: GameStatus;
  lastMove: Move | null;
  moves: Move[];
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
  type: "msg"
  data: {
    sender: string
    content: string
  }
}

export interface LobbyRequest {
  type: "lobby"
  data: {
    name: string
    timeControl: string;
    mode: string 
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

export interface MoveRequest {
  type: "move"
  data: {
    move: Move
  }
}

export type ClientRequest = 
  | MoveRequest
  | LobbyRequest
  | ReconnectRequest


//IMPORTANT: server always returns gamestate
export interface ServerResponse {
  type: string
  data: GameState
}



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
  result: "win" | "draw" | "loss";
  winner: PlayerRow | null;
}



