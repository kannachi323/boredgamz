
import { useState } from 'react';

import { GameContext } from './GameContext'
import { ServerMessage, GameState } from '../types';
import { createGame } from '../utils/game';


export const GameProvider = ({ children }: { children: React.ReactNode }) => {
  const [gameState, setGameState] = useState<GameState>(createGame());
  const [conn, setConn] = useState<WebSocket | null>(null);


  function update(payload: ServerMessage) {
    if (payload.type === "gameUpdate") {
      setGameState(payload.data.gameState)
      console.log("Game updated:", payload);
    } else if (payload.type === "chat") {
      console.log("Chat message:", payload);
    }
  }

  return (
    <GameContext.Provider value={{ gameState, update, conn, setConn }}>
      {children}
    </GameContext.Provider>
  );
};

