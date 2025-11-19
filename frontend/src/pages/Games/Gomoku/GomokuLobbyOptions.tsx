import { useState } from "react";

import { useGameStore } from "../../../stores/useGomokuStore.ts";
import { convertTime } from '../../../utils.ts'


export function LobbyOptionsPanel() {
  const { player, setPlayer } = useGameStore();
  const [colorSelected, setColorSelected] = useState(player.color)
  const [timeControlSelected, setTimeControlSelected] = useState("Rapid")

  function handleColorSelect(color: string) {
    setColorSelected(color)
    setPlayer({...player, color: color })
  }

  function handleTimeControlSelect(mode : string) {
    setTimeControlSelected(mode)
    let time : number;
    switch (mode) {
      case "Rapid":
        time = 5
        break
      case "Blitz":
        time = 3
        break
      case "Bullet":
        time = 1
        break
      case "Hyperbullet":
        time = 0.5
        break
      default:
        time = 5
    }

    setPlayer({...player, playerClock: {...player.playerClock, remaining: convertTime(time, "minutes", "nanoseconds")}})
  }

  return (
    <>
      <div className="flex flex-row items-center justify-evenly gap-2">
        <p className="text-2xl">Color:</p>
        <img
          src="/white.svg"
          alt="preview stone"
          className={`h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${colorSelected === "white" && " border-green-300"}`}
          onClick={() => handleColorSelect("white")}
        />
        <img
          src={`/black.svg`}
          alt={`preview stone`}
          className={`h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${colorSelected === "black" && " border-green-300"}`}
          onClick={() => handleColorSelect("black")}
        />
      </div>
     <div className="flex flex-row items-center justify-center gap-4 py-4">
        {/* Placeholder: You would define timeControlSelected and handleTimeControlSelect */}
        {/* For example: const [timeControlSelected, setTimeControlSelected] = useState('Rapid'); */}
        
        <p
          className={`text-xl font-semibold cursor-pointer border-2 rounded-xl px-4 py-2 transition-all duration-300 
            ${timeControlSelected === "Rapid" 
              ? "bg-[#585858] text-white" 
              : "bg-gray-100 text-gray-700 border-gray-300 hover:bg-gray-200"
            }`}
          onClick={() => handleTimeControlSelect("Rapid")}
        >
          Rapid (5 min)
        </p>

        <p
          className={`text-xl font-semibold cursor-pointer border-2 rounded-xl px-4 py-2 transition-all duration-300 
            ${timeControlSelected === "Blitz" 
              ? "bg-[#585858] text-white" 
              : "bg-gray-100 text-gray-700 border-gray-300 hover:bg-gray-200"
            }`}
          onClick={() => handleTimeControlSelect("Blitz")}
        >
          Blitz (3 min)
        </p>

        <p
          className={`text-xl font-semibold cursor-pointer border-2 rounded-xl px-4 py-2 transition-all duration-300 
            ${timeControlSelected === "Bullet" 
              ? "bg-[#585858] text-white" 
              : "bg-gray-100 text-gray-700 border-gray-300 hover:bg-gray-200"
            }`}
          onClick={() => handleTimeControlSelect("Bullet")}
        >
          Bullet (1 min)
        </p>

        {/* Note: Hyperbullet is often 30 seconds or less, but keeping your original text */}
        <p
          className={`text-xl font-semibold cursor-pointer border-2 rounded-xl px-4 py-2 transition-all duration-300 
            ${timeControlSelected === "Hyperbullet" 
              ? "bg-[#585858] text-white" 
              : "bg-gray-100 text-gray-700 border-gray-300 hover:bg-gray-200"
            }`}
          onClick={() => handleTimeControlSelect("Hyperbullet")}
        >
          Hyperbullet (30 sec)
        </p>
      </div>
    </>
  )
}