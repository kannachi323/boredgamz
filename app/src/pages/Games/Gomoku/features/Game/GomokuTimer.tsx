import { useEffect, useState } from "react";
import { Player } from "../../types";
import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";

export function Timer({ player }: { player: Player }) {
  const { gameState } = useGomokuStore();
  const [time, setTime] = useState(0);

  useEffect(() => {
    if (!player?.playerClock) return;

    const serverSeconds = Math.floor(player.playerClock.remaining / 1e9);
    console.log("changed server seconds to ", serverSeconds, " for ", player)
    setTime(serverSeconds)
  }, [gameState]);

  useEffect(() => {
    if (!gameState || gameState.status?.code !== "online") return;

    const interval = setInterval(() => {
      setTime((t) => {
        if (gameState.turn !== player.playerID) return t;

        return Math.max(0, t - 1);
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [gameState, player.playerID]);

  if (!player?.playerClock) return null;

  const formatTimer = (timer: number) => {
    const minutes = Math.floor(timer / 60);
    const seconds = timer % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  const isActive = gameState?.turn === player.playerID;

  return (
    <div
      className={`flex justify-center items-center rounded-lg transition-colors duration-300
        ${isActive ? "text-white" : " text-[#C3B299]"}`
      }
    >
      <b className="text-3xl">{formatTimer(time)}</b>
    </div>
  );
}
