import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useGomokuStore } from "@/stores/useGomokuStore";
import { GomokuLobbyOptions } from "./GomokuLobbyOptions";
import { GomokuLobbyBoards } from "./GomokuLobbyBoards";
import { GomokuLobbyModes } from "./GomokuLobbyModes";
import { GomokuModeModal } from "./GomokuModeModal";
import { GomokuScrollTooltip } from "./GomokuScrollToolTip";

export function GomokuLobby() {
  const { gameState } = useGomokuStore();
  const navigate = useNavigate();
  const [activeMode, setActiveMode] = useState("");
  const [showModeModal, setShowModeModal] = useState(false);

  useEffect(() => {
    if (gameState?.gameID && gameState?.status.code === "online") {
      navigate(`/games/gomoku/${gameState.gameID}`);
    }
  }, [gameState]);

  return (
    <div className="relative flex flex-col justify-center items-center p-10 gap-10">
      <h1 className="text-6xl text-[#C3B299] font-bold">Gomoku</h1>

      {/* Game Options */}
      <section className="flex flex-col items-center gap-1">
        <p className="text-lg text-[#C3B299] font-bold mb-1">Game</p>
        <div className="bg-[#433d3a] flex flex-row items-center justify-evenly p-3 rounded-xl gap-3">
          <GomokuLobbyOptions />
        </div>
      </section>

      {/* Mode */}
      <section className="flex flex-col items-center gap-1">
        <p className="text-lg text-[#C3B299] font-bold">Mode</p>
        <div className="bg-[#433d3a] p-3 rounded-xl flex flex-row justify-evenly gap-3">
          <GomokuLobbyModes
            onSelect={setActiveMode}
            onOpen={() => setShowModeModal(true)}
          />
        </div>
      </section>


      {/* Board */}
      <section className="flex flex-col items-center gap-1">
        <p className="text-lg text-[#C3B299] font-bold mb-1">Board</p>
        <div className="bg-[#433d3a] flex flex-row items-center justify-center p-3 rounded-xl gap-3">
          <GomokuLobbyBoards />
        </div>
      </section>

      {/* Mode popup */}
      {showModeModal && (
        <GomokuModeModal
          mode={activeMode}
          onClose={() => setShowModeModal(false)}
        />
      )}

      {/* tooltip */}
      <div className="fixed bottom-4 right-4 text-[#C3B299] rounded-lg shadow-lg text-sm 
        font-semibold flex items-center gap-3 cursor-pointer">
        <GomokuScrollTooltip />
      </div>

      {/* Start / Play Button */}
      <button
        className="px-10 py-3 bg-[#C3B299] text-[#433d3a] font-bold rounded-lg hover:bg-[#d7c9b8] transition"
        onClick={() => {
          console.log("Play Now clicked!");
        }}
      >
        Play Now
      </button>

      <div>{/* TODO: ads */}</div>
    </div>
  );
}
