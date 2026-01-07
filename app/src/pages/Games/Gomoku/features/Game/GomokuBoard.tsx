import { useState } from "react";
import { GomokuStone } from "./GomokuStone";
import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore"
import SMALL_BOARD from "@/assets/small-board.jpg"
import { Move } from "../../types";

export function GomokuBoard() {
  const { gameState, send, player, analysis } = useGomokuStore();

  const [hoveredIndex, setHoveredIndex] = useState<[number, number] | null>(null);

  if (!gameState || !gameState.board) {
    return null;
  }

  function sendMove(row: number, col: number) {
    const move : Move = {
      color: player.color,
      row: row,
      col: col
    }
    send({type: "move", data: {move: move}}, )
  }

  const board = analysis.active ? analysis.board : gameState.board
  return (
    <div className="flex justify-center h-full w-full relative">
      <img src={SMALL_BOARD} alt="gomoku board" className="absolute h-full w-full z-0" />
      <div className="absolute h-full w-full grid grid-cols-9 grid-rows-9 z-10 p-4">
        {board?.stones.flatMap((row, rowIdx) => 
          row.map((stone, colIdx) => (
            <div
              key={`${rowIdx}-${colIdx}`}
              className="h-full w-full z-20 flex justify-center items-center"
              onClick={() => sendMove(rowIdx, colIdx)}
              onMouseEnter={() => setHoveredIndex([rowIdx, colIdx])}
              onMouseLeave={() => setHoveredIndex(null)}
            >
              <GomokuStone
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

