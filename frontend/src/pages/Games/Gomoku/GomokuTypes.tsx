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
  moves: [];
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
    lobbyType: string
    player: Player
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


//IMPORTANT: server always returns gamestate
export interface ServerResponse {
  type: string
  data: GameState
}



