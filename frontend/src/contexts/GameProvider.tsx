
import { useState } from 'react';

import { GameContext } from './GameContext'
import { ServerResponse, GameState, Player } from '../types';


export const GameProvider = ({ children }: { children: React.ReactNode }) => {
  
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [conn, setConn] = useState<WebSocket | null>(null);
  const [player, setPlayer] = useState<Player>({
    playerID: 'annonymous',
    color: 'white',
  })

  function update(payload: ServerResponse) {
    if (payload.type === "update") {
      setGameState(payload.data)
      console.log("Game updated:", payload);
    } else if (payload.type === "chat") {
      console.log("Chat message:", payload);
    } else {
      setGameState(payload.data);
    }
  }


  return (
    <GameContext.Provider value={{ 
      gameState, setGameState, update, conn, setConn, player, setPlayer
    }}>
      {children}
    </GameContext.Provider>
  );
};

