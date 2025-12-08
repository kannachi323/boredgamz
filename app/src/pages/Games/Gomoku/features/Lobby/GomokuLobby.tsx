import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";
import { GomokuLobbyOptions } from "./GomokuLobbyOptions";
import { GomokuLobbyBoards } from "./GomokuLobbyBoards"
import { GomokuLobbyModes } from "@Gomoku/features/Lobby/GomokuLobbyModes";
import { GomokuLobbyWaiting } from "./GomokuLobbyWaiting";
import { GomokuScrollTooltip } from "./GomokuScrollToolTip";
import { useAuthStore } from "@/stores/useAuthStore";

export function GomokuLobby() {
  const { user } = useAuthStore();
  const { gameState, setConnection, lobbyRequest, setLobbyRequest, closeConnection } = useGomokuStore();
  const navigate = useNavigate();
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    if (!user?.id || !user?.username) return;
    setLobbyRequest({...lobbyRequest, data: {...lobbyRequest.data, 
      playerName: user.username,
      playerID: user.id, 
    }})
  }, [])

  useEffect(() => {
    if (gameState?.gameID && gameState?.status.code === "online") {
      const timer = window.setTimeout(() => {
        setIsSearching(false); 
        navigate(`/games/gomoku/gomoku-${lobbyRequest.data.mode}-${lobbyRequest.data.name}`);
      }, 2000);
      return () => clearTimeout(timer);
    }
  }, [gameState]);

  function handlePlayNow() {
    setIsSearching(true)
    setConnection(lobbyRequest);
  }

  function handlePlayCancel() {
    closeConnection()
    setIsSearching(false)
  }

  return (
    <div className="relative flex flex-col justify-center items-center p-10 gap-10">
      <h1 className="text-6xl text-[#C3B299] font-bold">Gomoku</h1>

      {/* Game Options */}
      <section className="flex flex-col items-center gap-1">
        <p className="text-lg text-[#C3B299] font-bold mb-1">Game</p>
        <div className="bg-[#433d3a] flex flex-row items-center justify-evenly p-3 rounded-xl gap-5">
          <GomokuLobbyOptions />
        </div>
      </section>

      {/* Mode */}
      <section className="flex flex-col items-center gap-1">
        <p className="text-lg text-[#C3B299] font-bold">Mode</p>
        <div className="bg-[#433d3a] p-3 rounded-xl flex flex-row justify-evenly gap-5">
          <GomokuLobbyModes/>
        </div>
      </section>


      {/* Board */}
      <section className="flex flex-col items-center gap-1">
        <p className="text-lg text-[#C3B299] font-bold mb-1">Board</p>
        <div className="bg-[#433d3a] flex flex-row items-center justify-center p-3 rounded-xl gap-5">
          <GomokuLobbyBoards />
        </div>
      </section>


      {/* tooltip */}
      <div className="fixed bottom-4 right-4 text-[#C3B299] rounded-lg shadow-lg text-sm 
        font-semibold flex items-center gap-3 cursor-pointer">
        <GomokuScrollTooltip />
      </div>

      {/* Start / Play Button */}
      <button
        className="px-10 py-3 bg-[#C3B299] text-[#433d3a] font-bold rounded-lg hover:bg-[#d7c9b8] transition"
        onClick={() => handlePlayNow()}
      >
        Play Now
      </button>

      <div>{/* TODO: ads */}</div>

      <div>{isSearching && <GomokuLobbyWaiting onCancel={() => handlePlayCancel()} />} </div>
    </div>
  );
}
