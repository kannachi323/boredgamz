
import { Timer } from "../../components/Timer"
import { PlayerBanner } from "../../components/Banner"
import { GamePanel } from "./GamePanel"
import { useGameStore } from "../../stores/useGameStore";
import { Board } from "../Board"; 
import { ChatBox } from "../Chat/ChatBox";

export default function Game() {
  const { gameState } = useGameStore();

  if (!gameState) return

  return (
    <div className="w-full h-[90vh] grid grid-cols-26 grid-rows-1 gap-10 p-10">
      <div className="col-span-7 row-span-1 flex flex-col justify-center gap-10">
        <div className="w-full h-1/2 flex flex-col items-center justify-center gap-2">
        <PlayerBanner player={gameState.players[1]}/>
        <Timer seconds={60}/>
        </div>

        <GamePanel />

        <div className="w-full h-1/2 flex flex-col items-center justify-center gap-2">
        <Timer seconds={60}/>
        <PlayerBanner player={gameState.players[0]}/>
        </div>
      </div>
      <div className="col-span-12 row-span-1">
        <Board/>
      </div>
      <div className="col-span-7 row-span-1">
        <ChatBox username={gameState.players[0].playerID}/>
      </div>
      </div>
  )
}