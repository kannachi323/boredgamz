import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";
import { X } from "lucide-react";

export function GameEnd() {
  const { gameState } = useGomokuStore();

  const [showGameEndModal, setShowGameEndModal] = useState(false);

  useEffect(() => {
    if (!gameState || gameState.status?.code === "online") return;
    setShowGameEndModal(true)
  
  }, [gameState?.status?.code])

  if (!gameState || gameState.status?.code === "online" || !showGameEndModal) return null;

  return (
    <div className="fixed inset-0 bg-black/40 flex justify-center items-center z-50 p-4">
      <GameEndCard onClose={() => setShowGameEndModal(false)} />
    </div>
  );
}

function GameEndCard({ onClose }: { onClose: () => void }) {
  const { gameState } = useGomokuStore();
  const navigate = useNavigate();

  if (!gameState || !gameState.status) return null;

  let title = "";
  let emoji = "";

  if (gameState.status.result === "win" && gameState.status.winner) {
    emoji = gameState.status.winner.color === "black" ? "⚫️" : "⚪️";
    title = `${gameState.status.winner.color === "black" ? "Black" : "White"} Wins!`;
  } else if (gameState.status.result === "draw") {
    emoji = "🤝";
    title = "Draw";
  }

  return (
    <div
      className="
        relative bg-[#302e2e] border border-[#1b1918] rounded-2xl shadow-xl
        px-10 py-8 w-full max-w-md text-center animate-fade-in
      "
    >
      {/* Close icon */}
      <button
        onClick={onClose}
        className="absolute top-4 right-4 text-[#C3B299] hover:text-white transition"
      >
        <X size={24} />
      </button>

      {/* Emoji */}
      <div className="text-6xl mb-4">{emoji}</div>

      {/* Title */}
      <h1 className="text-3xl font-bold text-[#C3B299] mb-6 tracking-wide">
        {title}
      </h1>

      {/* Buttons */}
      <div className="flex flex-col gap-4">
        <button
          onClick={() => navigate("/games/gomoku")}
          className="
            w-full py-3 text-lg font-semibold 
            bg-[#433d3a] hover:bg-[#524b4b] 
            text-[#C3B299] rounded-xl transition
          "
        >
          Return to Lobby
        </button>
      </div>
    </div>
  );
}
