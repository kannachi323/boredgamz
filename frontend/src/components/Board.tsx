import { useState } from "react";
import { Stone } from "../types";
import { useGameContext } from "../hooks/useGameContext";
import { makeMove } from "../utils/game";

export function GomokuBoard() {
  const { gameState, conn } = useGameContext();

  const [hoveredIndex, setHoveredIndex] = useState<[number, number] | null>(null);

  if (!gameState || !gameState.board) {
    return null;
  }

  return (
    <div className="flex justify-center h-full w-full relative">
      <img src="/small-board.jpg" alt="gomoku board" className="absolute h-full w-full z-0" />
      <div className="absolute h-full w-full grid grid-cols-9 grid-rows-9 z-10 p-4">
        {gameState.board.stones.map((row, rowIdx) => (
          <>
            {row.map((stone, colIdx) => (
              <div
                key={colIdx}
                className="h-full w-full z-20 flex justify-center items-center"
                onClick={() => makeMove(conn, gameState, 
                  { r: rowIdx, c: colIdx, color: gameState.turn === "P1" ? gameState.players[0].color : gameState.players[1].color}
                )}
                onMouseEnter={() => setHoveredIndex([rowIdx, colIdx])}
                onMouseLeave={() => setHoveredIndex(null)}
              >
                <StonePiece
                  stone={stone}
                  isHovered={
                    hoveredIndex &&
                    hoveredIndex[0] === rowIdx &&
                    hoveredIndex[1] === colIdx &&
                    !stone.color
                  }
                />
              </div>
            ))}
          </>
        ))}
      </div>
    </div>
  );
}

function StonePiece({ stone, isHovered }: { stone: Stone; isHovered: boolean | null}) {
  if (stone.color) {
    return (
      <img
        src={`/${stone.color}.svg`}
        alt={`${stone.color} stone`}
        className="h-full w-full opacity-100"
      />
    );
  }

  if (isHovered) {
    return (
      <img
        src={`/black.svg`}
        alt={`preview stone`}
        className="h-14 w-14 opacity-50"
      />
    );
  }

  return null;
}
