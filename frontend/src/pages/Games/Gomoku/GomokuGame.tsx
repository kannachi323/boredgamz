import { useEffect } from "react";
import { useParams } from "react-router-dom";

import { PlayerBanner } from "../../../components/Banner";
import { GamePanel } from "./GomkuGamePanel";
import { useGomokuStore } from "../../../stores/useGomokuStore";
import { GomokuBoard } from "./GomokuBoard";
import { ChatBox } from "../../../features/Chat/ChatBox";
import { GameEnd } from "./GomokuGameEndModal";
import { Timer } from "../../../components/Timer";


export default function GomokuGame() {
  const { gameState, setPlayer, setOpponent, player, opponent, loadGame } = useGomokuStore();
  const { gameID } = useParams();

  useEffect(() => {
    if (!player || !gameState) return;

    const p1 = gameState.players[0];
    const p2 = gameState.players[1];

    setPlayer(p1.playerID === player.playerID ? p1 : p2);
    setOpponent(p1.playerID === player.playerID ? p2 : p1);
  }, [gameState]);

  useEffect(() => {
    if (!gameID) return;
    loadGame(gameID)
    console.log(gameState)
  }, [gameID])

  return (
    <div className="h-[90vh] w-full grid grid-cols-26 gap-6 p-6 bg-[#1b1918] overflow-hidden">

      {/* LEFT PANEL */}
      <section className="col-span-7 flex flex-col gap-2 bg-[#433d3a] p-2 rounded-xl border border-[#1b1918] min-h-0">

        {/* Opponent Info — FIXED HEIGHT */}
        <div className="bg-[#302e2e] h-20 p-3 rounded-xl flex flex-row items-center justify-between border border-[#1b1918]">
          <PlayerBanner player={opponent} />
          <Timer player={opponent} />
        </div>

        {/* Game Panel — FLEXES, CAN SCROLL INSIDE */}
        <div className="flex-1 min-h-0 bg-[#302e2e] rounded-xl border border-[#1b1918]">
          <GamePanel />
        </div>

        {/* Player Info — FIXED HEIGHT */}
        <div className="bg-[#302e2e] h-20 p-3 rounded-xl flex flex-row items-center justify-between border border-[#1b1918]">
          <PlayerBanner player={player} />
          <Timer player={player} />
        </div>

      </section>

      {/* BOARD SECTION */}
      <section className="col-span-12 bg-[#433d3a] p-2 rounded-xl border border-[#1b1918] flex justify-center items-center min-h-0">
        <div className="bg-[#302e2e] p-2 rounded-xl border border-[#1b1918] h-full w-full overflow-hidden">
          <GomokuBoard />
        </div>
      </section>

      {/* CHAT SECTION */}
      <section className="col-span-7 bg-[#433d3a] p-2 rounded-xl border border-[#1b1918] flex flex-col min-h-0">
        <ChatBox username={player.playerName} />
      </section>

      <GameEnd />
    </div>
  );
}
