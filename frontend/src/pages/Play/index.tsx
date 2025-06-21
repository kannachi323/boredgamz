
import {GomokuBoard} from "../../components/Board";
import { PlayerBanner } from "../../components/Banner";
import { GamePanel } from "../../components/GamePanel";
import { ChatBox } from "../../components/ChatBox";
import { Lobby } from "./Lobby";

import { Timer } from "../../components/Timer";


import { useGameContext } from "../../hooks/useGameContext";



export default function Play() {
  const { gameState } = useGameContext();



  // call the find game in like a lobby or something idk
 

  if (!gameState) {
    return (
       <div className="flex flex-col items-start justify-evenly w-full h-[90vh] p-10 gap-4 ">
        <Lobby />
       </div>
    );
  }

  console.log("Game State: ", gameState);

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
        <GomokuBoard/>
      </div>
      <div className="col-span-7 row-span-1">
        <ChatBox username={gameState.players[0].playerID}/>
      </div>
    </div>

  );
}
