import { useState } from "react";

import { useGameStore } from "../../stores/useGameStore";


export function LobbyOptionsPanel() {
  const { player, setPlayer } = useGameStore();
  const [colorSelected, setColorSelected] = useState(player.color)

  function handleColorSelect(color: string) {
    setColorSelected(color)
    setPlayer({...player, color: color })
  }

  console.log(player)

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
    </>
  )
}