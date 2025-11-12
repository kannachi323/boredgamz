import { useGameStore } from "../../stores/useGameStore";
import { LobbyOptionsPanel } from "./LobbyOptionsPanel";

export function Lobby() {
  const { setConnection, player, handler } = useGameStore();

  

  if (!player) {
    //TODO: navigate to error screen
    return
  }
  
  return (
    <>
      <div className="bg-[#302e2e] flex flex-row items-center justify-evenly w-full gap-5 p-5">
        <LobbyOptionsPanel />
      </div>
      <div className="bg-[#302e2e] flex flex-row items-center justify-evenly w-full gap-5 p-5">
        <div className="bg-[#302e2e] w-1/3 p-5 flex flex-col items-center justify-center gap-4 rounded-lg border-2 border-[#1b1918]
          hover:bg-[#3d3a3a] transition-colors duration-300 cursor-pointer"
          
          onClick={() => setConnection("9x9", player, handler)}
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
      </div>
      
    </>
  )
}


