import { useEffect } from "react";
import { Lobby } from "../../features/Lobby"
import { useGameStore } from "../../stores/useGameStore"
import { Outlet, useNavigate, useParams } from "react-router-dom";

export function PlayQuick() {
  const { gameState } = useGameStore();
  const navigate = useNavigate();
  const { mode } = useParams();

  useEffect(() => {
    if (gameState) {
      navigate(`/play/${mode}/${gameState.gameID}`)
    }
  }, [gameState, navigate, mode])

  return (
    <div className="flex flex-col w-full h-[90vh] p-10 gap-4"> 
      {gameState ? <Outlet /> : <Lobby />}
    </div>
  )
}