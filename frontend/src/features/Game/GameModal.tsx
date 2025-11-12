import { useGameStore } from "../../stores/useGameStore";

export function GameModal() {
  const { gameState, player } = useGameStore();

  if (!gameState || gameState.status.code === "online") return null;

  let modal: React.ReactNode;
  switch (gameState.status.result) {
    case "win":
      if (player.playerID === gameState.status.winner) {
        modal = <GameWinModal />;
      } else {
        modal = <GameLossModal />;
      }
      break;
    case "draw":
      modal = <GameDrawModal />;
      break;
    default:
      return null;
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50">
      {modal}
    </div>
  );
}

function GameWinModal() {
  return (
    <div className="bg-[#262322] p-10 rounded-md text-white w-1/2 text-center">
      YOU WON 🎉
    </div>
  );
}

function GameLossModal() {
  return (
    <div className="bg-[#262322] p-10 rounded-md text-white w-1/2 text-center">
      YOU LOST 😢
    </div>
  );
}

function GameDrawModal() {
  return (
    <div className="bg-[#262322] p-10 rounded-md text-white w-1/2 text-center">
      DRAW 🤝
    </div>
  );
}
