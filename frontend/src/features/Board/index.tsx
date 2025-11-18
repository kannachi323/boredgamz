import { useState } from "react";
import { BoardStonePiece } from "./BoardStonePiece";
import { useGameStore } from "../../stores/useGomokuStore";
import { sendMove } from "./board";

export function Board() {
  const { gameState } = useGameStore();

  const [hoveredIndex, setHoveredIndex] = useState<[number, number] | null>(null);

  if (!gameState || !gameState.board) {
    return null;
  }

  return (
    <div className="flex justify-center h-full w-full relative">
      <img src="/small-board.jpg" alt="gomoku board" className="absolute h-full w-full z-0" />
      <div className="absolute h-full w-full grid grid-cols-9 grid-rows-9 z-10 p-4">
        {gameState.board.stones.flatMap((row, rowIdx) => 
          row.map((stone, colIdx) => (
            <div
              key={`${rowIdx}-${colIdx}`}
              className="h-full w-full z-20 flex justify-center items-center"
              onClick={() => sendMove(rowIdx, colIdx)}
              onMouseEnter={() => setHoveredIndex([rowIdx, colIdx])}
              onMouseLeave={() => setHoveredIndex(null)}
            >
              <BoardStonePiece
                stone={stone}
                isHovered={
                  hoveredIndex &&
                  hoveredIndex[0] === rowIdx &&
                  hoveredIndex[1] === colIdx &&
                  !stone.color
                }
              />
            </div>
          ))
        )}
      </div>
    </div>
  );
}

