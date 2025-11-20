import { useState } from "react";

import { useGomokuStore } from "../../../stores/useGomokuStore.ts";
import { convertTime } from '../../../utils.ts'
import BLACK from '../../../assets/black.svg'
import WHITE from '../../../assets/white.svg'


export function GomokuLobbyOptions() {
  const { player, setPlayer } = useGomokuStore();
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
        <p className="text-2xl text-[#C3B299] font-bold">Color:</p>

        {/* White Stone */}
        <img
          src={WHITE}
          alt="white stone"
          className={`
            h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${colorSelected === "white"
              ? "border-[#7DCFB6] shadow-[0_0_8px_#7DCFB6]"
              : "border-[#C3B299] hover:border-[#7DCFB6]"
            }
          `}
          onClick={() => handleColorSelect("white")}
        />

        {/* Black Stone */}
        <img
          src={BLACK}
          alt="black stone"
          className={`
            h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${colorSelected === "black"
              ? "border-[#7DCFB6] shadow-[0_0_8px_#7DCFB6]"
              : "border-[#C3B299] hover:border-[#7DCFB6]"
            }
          `}
          onClick={() => handleColorSelect("black")}
        />
      </div>

      <div className="flex flex-row items-center justify-center gap-4">
        {[
          { label: "Rapid (5 min)", mode: "Rapid" },
          { label: "Blitz (3 min)", mode: "Blitz" },
          { label: "Bullet (1 min)", mode: "Bullet" },
          { label: "Hyperbullet (30 sec)", mode: "Hyperbullet" },
        ].map(({ label, mode }) => (
          <p
            key={mode}
            className={`
              text-xl font-semibold cursor-pointer rounded-xl px-5 py-2 transition-all duration-300
              border-2
              ${
                timeControlSelected === mode
                  ? "bg-[#f4c97f] text-black border-[#7DCFB6]"
                  : "bg-[#2B2825] text-[#C3B299] border-[#514C47] hover:bg-[#514C47]"
              }
            `}
            onClick={() => handleTimeControlSelect(mode)}
          >
            {label}
          </p>
        ))}
      </div>

    </>
  )
}