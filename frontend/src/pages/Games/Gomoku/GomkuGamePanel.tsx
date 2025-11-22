import { useEffect } from "react";

import { Handshake, RefreshCw, Plus, Flag } from "lucide-react";
import { ChevronsLeft, ChevronLeft, ChevronRight, ChevronsRight } from "lucide-react";

import { useGomokuStore } from "../../../stores/useGomokuStore";
import { Move } from "./GomokuTypes";


export function GamePanel() {
  const { analysis, setAnalysisIndex, gameState, exitAnalysis, startAnalysis } = useGomokuStore();

  useEffect(() => {
    exitAnalysis()
  }, [gameState])

  return (
    <div className="text-[#C3B299] rounded-lg flex flex-col w-full h-full">
      
      {/* Move History (takes remaining space) */}
      <div className="flex-1 overflow-y-auto text-sm border-b rounded-md p-2">
        {analysis.moves.length === 0 ? (
          <p className="text-[#C3B299]/60 text-center">No moves yet</p>
        ) : (
          analysis.moves.map((move: Move, idx: number) => (
            <p
              key={idx}
              className={idx === analysis.index ? "font-bold text-green-300" : ""}
            >
              {idx + 1}. {move.color === "black" ? "B" : "W"} → ({move.row + 1},{" "}
              {move.col + 1})
            </p>
          ))
        )}
      </div>

      {/* Move Navigation */}
      <div className="flex items-center justify-center gap-4 py-2">

        <button
          className="p-1 rounded hover:bg-[#524b4b] disabled:opacity-40"
          disabled={analysis.index === -1}
          onClick={() => startAnalysis()}
        >
          <ChevronsLeft size={28} className="text-[#C3B299]" />
        </button>

        <button
          className="p-1 rounded hover:bg-[#524b4b] disabled:opacity-40"
          disabled={analysis.index === -1}
          onClick={() => setAnalysisIndex(Math.max(analysis.index - 1, -1))}
        >
          <ChevronLeft size={28} className="text-[#C3B299]" />
        </button>

        <button
          className="p-1 rounded hover:bg-[#524b4b] disabled:opacity-40"
          disabled={analysis.index >= analysis.moves.length - 1}
          onClick={() =>
            setAnalysisIndex(Math.min(analysis.index + 1, analysis.moves.length - 1))
          }
        >
          <ChevronRight size={28} className="text-[#C3B299]" />
        </button>

        <button
          className="p-1 rounded hover:bg-[#524b4b] disabled:opacity-40"
          disabled={analysis.index >= analysis.moves.length - 1}
          onClick={() => exitAnalysis()}
        >
          <ChevronsRight size={28} className="text-[#C3B299]" />
        </button>

      </div>


      {/* Game Controls */}
      <div className="flex flex-row justify-center items-center gap-6 p-3 border-t border-[#1b1918]/60">
        <button className="hover:text-[#C3B299]/70 transition flex flex-row gap-2 items-center">
          <p className="font-bold">Draw</p>
          <Handshake size={19} />
        </button>

        <button className="hover:text-red-400 transition flex flex-row gap-2 items-center">
          <p className="font-bold">Resign</p>
          <Flag size={19} />
        </button>

        <button className="hover:text-[#C3B299]/70 transition flex flex-row gap-2 items-center">
          <p className="font-bold">Rematch</p>
          <RefreshCw size={19} />
        </button>

        <button className="hover:text-green-400 transition flex flex-row gap-2 items-center">
          <p className="font-bold">New</p>
          <Plus size={24} />
        </button>
      </div>
    </div>
  );
}
