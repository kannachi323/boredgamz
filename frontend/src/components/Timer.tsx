import { useEffect, useState } from "react";
import { Player } from "../types";
import { useGameStore } from "../stores/useGameStore";

export function Timer({ player }: { player: Player }) {
  const { gameState } = useGameStore();
  const [time, setTime] = useState(0);

  useEffect(() => {
    if (!gameState || !player?.playerClock) return;

    // Sync with server's remaining time on every update
    const serverSeconds = Math.floor(player.playerClock.remaining / 1e9);
    setTime(serverSeconds);

    const interval = setInterval(() => {
      if (gameState.turn === player.playerID) {
        setTime((t) => Math.max(0, t - 1));
        if (time <= 0) return
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [gameState?.turn, player?.playerClock?.remaining, player?.playerID]);

  if (!player?.playerClock) return null;

  const formatTimer = (timer: number) => {
    const minutes = Math.floor(timer / 60);
    const seconds = timer % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  const isActive = gameState?.turn === player.playerID;

  return (
    <div
      className={`w-1/3 p-2 flex justify-center items-center rounded-lg transition-colors duration-300
        ${isActive ? "bg-green-700 text-white" : "bg-[#363430] text-gray-300"}`}
    >
      <b className="text-3xl">{formatTimer(time)}</b>
    </div>
  );
}
