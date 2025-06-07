import {createContext } from 'react';
import { GameState, ServerResponse } from '../types';


interface GameContext {
  gameState: GameState | null;
  setGameState: (gameState: GameState) => void;
  update: (payload: ServerResponse) => void;
  conn: WebSocket | null;
  setConn: (conn: WebSocket) => void;
}

export const GameContext = createContext<GameContext | undefined>(undefined);

