import { useEffect, useState } from "react";
import { Player } from "../../types";
import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";

export function GomokuTimer({ player }: { player: Player }) {
  // Initialize with the prop value to avoid a "0:00" flash
  const [time, setTime] = useState(() => 
    player?.playerClock ? Math.floor(player.playerClock.remaining / 1e9) : 0
  );
  
  const { gameState } = useGomokuStore();

  // 1. Sync Effect: Updates local time when Server sends a new GameState
  useEffect(() => {
    if (!gameState || !player) return;

    // CRITICAL FIX: Find the player in the NEW gameState. 
    // Do not trust the 'player' prop to be up-to-date inside this effect.
    const freshPlayer = gameState.players.find((p) => p.playerID === player.playerID);

    if (freshPlayer?.playerClock) {
      const serverSeconds = Math.floor(freshPlayer.playerClock.remaining / 1e9);
      setTime(serverSeconds);
    }
  }, [gameState, player.playerID]); // Depend on ID, not the whole player object

  // 2. Countdown Effect: Ticks down locally every second
  useEffect(() => {
    if (!gameState || gameState.status?.code !== "online") return;
    if (gameState.turn !== player.playerID) return;

    const interval = setInterval(() => {
      setTime((t) => Math.max(0, t - 1));
    }, 1000);

    return () => clearInterval(interval);
  }, [gameState?.status?.code, gameState?.turn, player.playerID]);

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
        ${isActive ? "text-white" : "text-[#C3B299]"}`
      }
    >
      <b className="text-3xl">{formatTimer(time)}</b>
    </div>
  );
}