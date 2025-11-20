import { useNavigate } from "react-router-dom";
import { useGomokuStore } from "../../../stores/useGomokuStore";

export function GameModal() {
  const { gameState } = useGomokuStore();
  if (!gameState || gameState.status.code === "online") return null;

  return (
    <div className="fixed inset-0 backdrop-blur-sm bg-black/40 flex justify-center items-center z-50 p-4">
      <EndGameCard />
    </div>
  );
}

function EndGameCard() {
  const { gameState } = useGomokuStore();
  const navigate = useNavigate();

  if (!gameState) return null;

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
        bg-[#1e1d1b] border border-gray-700 rounded-2xl shadow-xl 
        px-10 py-8 w-full max-w-md text-center
        animate-fade-in
      "
    >

      {/* Emoji / Icon */}
      <div className="text-6xl mb-4">{emoji}</div>

      {/* Title */}
      <h1 className="text-3xl font-bold text-white mb-6 tracking-wide">
        {title}
      </h1>

      {/* Buttons */}
      <div className="flex flex-col gap-4">

        {/* Return to Lobby */}
        <button
          onClick={() => navigate("/games/gomoku")}
          className="
            w-full py-3 text-lg font-semibold 
            bg-indigo-600 hover:bg-indigo-700 
            text-white rounded-xl transition
          "
        >
          Return to Lobby
        </button>

        {/* Optional: Play Again */}
        {/* <button className="
          w-full py-3 text-lg font-semibold 
          bg-gray-700 hover:bg-gray-600 
          text-gray-200 rounded-xl transition
        ">
          Play Again
        </button> */}
      </div>

    </div>
  );
}
