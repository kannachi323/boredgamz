import { useEffect } from "react";

import { Timer } from "../../../components/Timer"
import { PlayerBanner } from "../../../components/Banner"
import { GamePanel } from "./GamePanel"
import { useGameStore } from "../../../stores/useGameStore";
import { Board } from "../../../features/Board"; 
import { ChatBox } from "../../../features/Chat/ChatBox";
import { GameModal } from "./GomokuModal";

/*TODO: need to implement game state saving after refresh (use database :)
Matthew pls do this asap lol this is pretty important
*/


export default function GomokuGame() {
  const { gameState, setPlayer, setOpponent, player, opponent } = useGameStore();

  useEffect(() => {
    //This effect adds the player clocks from server
    if (!player || !gameState) return
    const p1 = gameState.players[0]
    const p2 = gameState.players[1]
    setPlayer(p1.playerID == player.playerID ? p1 : p2)
    setOpponent( p1.playerID == player.playerID ? p2 : p1)
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