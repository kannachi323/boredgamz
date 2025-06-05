export interface Conn {
  roomID: string;
  player: Player
  type: string
}

export interface GameState {
  board: Stone[];
  size: number;
  players: {
    p1: Player;
    p2: Player;
  }
  turn: string;
  status: string;
}

export interface Player {
    playerID: string;
    color: string;
}

export interface Stone {
  color: number
  colorName: string
}

export interface ServerMessage {
  type: string
  data: GameState
}