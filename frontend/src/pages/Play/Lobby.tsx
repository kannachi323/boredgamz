import { createConnection } from "../../utils/connection";
import { useGameContext } from "../../hooks/useGameContext";
import { useEffect } from "react";
import { useAuthContext } from "../../hooks/useAuthContext";

export function Lobby() {
  const { setPlayer } = useGameContext();
  const { user } = useAuthContext();

  useEffect(() => {
    if (!user) return;
    setPlayer({
      playerID: user.id,
      color: "white",
    });
  }, [user, setPlayer])
  return (
    <>

      <div className="bg-[#302e2e] flex flex-row items-center justify-evenly w-full gap-5 p-5">
        <OptionsPanel />
      </div>
      <div className="flex flex-row items-center justify-evenly p-5 w-full gap-5">
        <QuickPlay />
      </div>

    </>
  )
}

function OptionsPanel() {
  const { player, setPlayer } = useGameContext();

  return (
    <>
      <div className="flex flex-row items-center justify-evenly gap-2">
        <p className="text-2xl">Color:</p>
        <img
          src="/white.svg"
          alt="preview stone"
          className={`h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${player.color === "white" && " border-green-300"}`}
          onClick={() => setPlayer({ ...player, color: "white" })}
        />
        <img
          src={`/black.svg`}
          alt={`preview stone`}
          className={`h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${player.color === "black" && " border-green-300"}`}
          onClick={() => setPlayer({...player, color: "black" })}
        />
      </div>
      
      
    
    </>
  )
}

function QuickPlay() {
  const { setConn, update, player } = useGameContext();

  return (
    <>
      <div className="bg-[#302e2e] w-1/3 p-5 flex flex-col items-center justify-center gap-4 rounded-lg border-2 border-[#1b1918]
        hover:bg-[#3d3a3a] transition-colors duration-300 cursor-pointer"
        
        onClick={() => setConn(createConnection("9x9", player, update))}
      >
        <p className="text-5xl">9x9</p>
        <img src="/small-board.jpg" alt="gomoku board" className="w-full h-auto" />
      </div>

      <div className="bg-[#302e2e] w-1/3 p-5 flex flex-col items-center justify-center gap-4 rounded-lg border-2 border-[#1b1918]
        hover:bg-[#3d3a3a] transition-colors duration-300 cursor-pointer">
        <p className="text-5xl">13x13</p>
        <img src="/mid-board.jpg" alt="gomoku board" className="w-full h-auto" />
      </div>



      <div className="bg-[#302e2e] w-1/3 p-5 flex flex-col items-center justify-center gap-4 rounded-lg border-2 border-[#1b1918]
        hover:bg-[#3d3a3a] transition-colors duration-300 cursor-pointer">
        <p className="text-5xl">19x19</p>
        <img src="/large-board.jpg" alt="gomoku board" className="w-full h-auto" />
      </div>



     

      
    
    
    
    </>
  )
}