import { useState } from "react";
import { Stone } from "../types";
import { placeStone } from "../utils/game";
import { useGameContext } from "../hooks/useGameContext";

export function GomokuBoard() {
  const { gameState, conn } = useGameContext();
  const [board, setBoard] = useState<Stone[]>(
    Array(gameState.size * gameState.size).fill({ colorName: null, color: null })
  );
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);

  if (!gameState || !board) {
    return null; // Or some loading UI
  }

  return (
    <div className="flex justify-center h-full w-full relative">
      <img src="/small-board.jpg" alt="gomoku board" className="absolute h-full w-full z-0" />
      <div className="absolute h-full w-full grid grid-cols-9 grid-rows-9 gap-6 z-10 p-10">
        {board.map((stone, idx) => (
          <div
            key={idx}
            className="h-14 w-14 z-20 flex justify-center items-center"
            onClick={() => {
              if (!stone.colorName) placeStone(conn, idx, setBoard);
            }}
            onMouseEnter={() => setHoveredIndex(idx)}
            onMouseLeave={() => setHoveredIndex(null)}
          >
            <StonePiece stone={stone} isHovered={hoveredIndex === idx && !stone.colorName} />
          </div>
        ))}
      </div>
    </div>
  );
}

function StonePiece({ stone, isHovered }: { stone: Stone; isHovered: boolean }) {
  if (stone.colorName) {
    return (
      <img
        src={`/${stone.colorName}.svg`}
        alt={`${stone.colorName} stone`}
        className="h-14 w-14 opacity-100"
      />
    );
  }

  // Show hover ghost stone
  if (isHovered) {
    return (
      <img
        src={`/black.svg`} // or `/white.svg` depending on current turn
        alt={`preview stone`}
        className="h-14 w-14 opacity-50"
      />
    );
  }

  return null;
}
