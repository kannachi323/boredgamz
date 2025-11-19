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
  r: number;
  c: number;
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

export interface ServerResponse {
  type: string
  data: GameState
}

export interface ClientRequest {
  type: string
  data: GameState
}

export interface LobbyRequest {
  lobbyType: string 
  player: Player
}

export interface Message {
  sender: string
  content: string
}