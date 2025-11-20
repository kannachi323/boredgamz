import { useEffect, useState } from "react";

import { useNavigate } from "react-router-dom";
import { useGomokuStore } from "../../../stores/useGomokuStore";
import { GomokuLobbyOptions } from "./GomokuLobbyOptions";
import { GomokuLobbyBoards } from "./GomokuLobbyBoards";
import { GomokuLobbyModes } from "./GomokuLobbyModes";
import { GomokuModeModal } from "./GomokuModeModal";


export function GomokuLobby() {
  const { gameState } = useGomokuStore();
  const navigate = useNavigate();
  const [activeMode, setActiveMode] = useState('')
  const [showModeModal, setShowModeModal] = useState(false)

  useEffect(() => {
    if (gameState?.gameID && gameState?.status.code === "online") {
      navigate(`/games/gomoku/${gameState.gameID}`)
    }
  }, [gameState])

  return (
    <>
      <h1 className="text-7xl text-[#C3B299] font-bold">Gomoku</h1>

      <section className="flex flex-col justify-center items-center w-full h-full gap-2">
        <p className="text-xl text-[#C3B299] font-bold">Mode</p>
        <div className="bg-[#433d3a] w-4/5 p-4 rounded-2xl">
          <GomokuLobbyModes onSelect={setActiveMode} onOpen={() => setShowModeModal(true)} />
        </div>
      </section>



      <section className="flex flex-col justify-center items-center w-full h-full gap-2">
        <p className="text-xl text-[#C3B299] font-bold">Game</p>
        <div className="bg-[#433d3a] flex flex-row items-center justify-evenly w-4/5 h-1/8 gap-5 p-3 rounded-2xl">
          <GomokuLobbyOptions />
        </div>
      </section>

             
    

      <section className="flex flex-col justify-center items-center w-full h-full gap-2">
        <p className="text-xl text-[#C3B299] font-bold">Board</p>
        <div className="bg-[#433d3a] flex flex-row items-center justify-evenly w-4/5 gap-5 p-5 rounded-2xl">
          <GomokuLobbyBoards />
        </div>
      </section>

      {showModeModal && (
        <GomokuModeModal 
          mode={activeMode}
          onClose={() => setShowModeModal(false)}
        />
      )}



      
      

      <div>
        {/* TODO: ads */}


      </div>
    </>
  )
}


