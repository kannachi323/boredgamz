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
      <div className="w-full h-[90vh] grid grid-cols-17 bg-red-50">
        <div className="col-span-4">
          <PlayerBanner />
        </div>
        <div className="col-span-9 bg-blue-500">
          <GomokuBoard/>
        </div>
        <div className="col-span-4">
          <PlayerBanner />
        </div>
        <div className="col-span-4">
          <button className="border-2 outline-2 w-full"
          
            onClick={() => setConn(createConnection(gameState.players.p1, update))}
          >
            join game
          </button>
          <button className="border-2 outline-2 w-full"
          
            onClick={() => {
              if (!conn) return;
              sendData(conn, {data: "aldkfja;sdkljas;dkl"})
            }}
          >
            send data
          </button>
        </div>
        
      </div>

  );
}
