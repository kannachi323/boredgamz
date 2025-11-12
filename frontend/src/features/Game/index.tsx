import { useEffect } from "react";

import { Timer } from "../../components/Timer"
import { PlayerBanner } from "../../components/Banner"
import { GamePanel } from "./GamePanel"
import { useGameStore } from "../../stores/useGameStore";
import { useAuthStore } from "../../stores/useAuthStore";
import { Board } from "../Board"; 
import { ChatBox } from "../Chat/ChatBox";
import { GameModal } from "./GameModal";

export default function Game() {
  const { gameState, setPlayer, setOpponent, player, opponent } = useGameStore();
  const { user } = useAuthStore();

  useEffect(() => {
    //This effect adds the player clocks from server
    if (!user || !gameState) return
    const p1 = gameState.players[0]
    const p2 = gameState.players[1]
    const player = p1.playerID == user.id ? p1 : p2
    const opponent = p1.playerID == user.id ? p2 : p1
    setPlayer(player)
    setOpponent(opponent)
  }, [gameState])

  return (
    <div className="w-full h-[90vh] grid grid-cols-26 grid-rows-1 gap-10 p-10">
      <div className="col-span-7 row-span-1 flex flex-col justify-center gap-10">
        
        <div className="w-full h-1/2 flex flex-col items-center justify-center gap-2">
          <PlayerBanner player={opponent}/>
          <Timer player={opponent}/>
        </div>

        <GamePanel />

        <div className="w-full h-1/2 flex flex-col items-center justify-center gap-2">
          <Timer player={player}/>
          <PlayerBanner player={player}/>
        </div>

      </div>

      <div className="col-span-12 row-span-1">
        <Board/>
      </div>

      <div className="col-span-7 row-span-1">
        <ChatBox username={player.playerName}/>
      </div>


      <GameModal />

    </div>
  )
}