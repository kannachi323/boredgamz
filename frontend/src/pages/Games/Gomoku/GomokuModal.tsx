import { useGameStore } from "../../../stores/useGameStore";

export function GameModal() {
  const { gameState } = useGameStore();

  if (!gameState || gameState.status.code === "online") return null;

  let modal: React.ReactNode;
  switch (gameState.status.result) {
    case "win":
      modal = <GameWinModal />;
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
  const { gameState } = useGameStore();
  if (!gameState || !gameState.status.winner) return

  const winnerColor = gameState.status.winner.color

  return (
    <div className="bg-[#262322] p-10 rounded-md text-white w-1/2 text-center">
      {winnerColor === 'black' ?
        <b>Black Won</b> : <b>White Won</b>
      }
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
