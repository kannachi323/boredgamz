import { useNavigate } from "react-router-dom";
import { useEffect } from "react";

import { useAuthStore } from "../../stores/useAuthStore";
import { useGameStore } from "../../stores/useGameStore";
import { Player } from "../../types";
import { LobbyOptionsPanel } from "./LobbyOptionsPanel";

export function Lobby() {
  const { isAuthenticated, user, checkAuth } = useAuthStore();
  const { setPlayer, setConnection, player, handler } = useGameStore();
  const navigate = useNavigate();

  useEffect(() => {
    const check = async () => {
      const success = await checkAuth(() => navigate('/login'));
      if (!success) return;

      const user = useAuthStore.getState().user;
      if (!user) return;

      const player: Player = {
        playerID: user.id,
        color: 'white',
      };

      setPlayer(player);
    };
    check()
  }, [isAuthenticated, user, navigate, setPlayer, checkAuth])

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


