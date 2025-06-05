import {GomokuBoard} from "../components/Board";
import { PlayerBanner } from "../components/Banner";

import { createConnection, sendData } from "../utils/connection";

import { useGameContext } from "../hooks/useGameContext";
export default function Play() {

  const {conn, setConn, gameState, update} = useGameContext();

  if (!gameState) {
    return;
  }

  return (
      <div className="w-full h-[90vh] grid grid-cols-17 grid-rows-1">
        <div className="col-span-4 h-full">
          <PlayerBanner />
        </div>
        <div className="col-span-9 h-full">
          <GomokuBoard/>
        </div>
        <div className="col-span-4">
          <PlayerBanner />
        </div>
     
        
      </div>

  );
}
