import {createContext } from 'react';
import { GameState, ServerMessage } from '../types';


interface GameContext {
  gameState: GameState;
  setGameState?: (gameState: GameState) => void;
  update: (payload: ServerMessage) => void;
  conn: WebSocket | null;
  setConn: (conn: WebSocket) => void;
}

export const GameContext = createContext<GameContext | undefined>(undefined);

