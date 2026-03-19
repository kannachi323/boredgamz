import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";

import BLACK from '@/assets/black.svg'
import WHITE from '@/assets/white.svg'


export function GomokuLobbyOptions() {
  const { setLobbyRequest, lobbyRequest } = useGomokuStore();

  function handleColorSelect(color: string) {
    setLobbyRequest({...lobbyRequest, data: { ...lobbyRequest.data, playerColor: color}})
  }

  function handleTimeControlSelect(timeControl : string) {
    setLobbyRequest({...lobbyRequest, data: { ...lobbyRequest.data, timeControl: timeControl}})
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
            ${lobbyRequest.data.playerColor === "white"
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
            ${lobbyRequest.data.playerColor === "black"
              ? "border-[#7DCFB6] shadow-[0_0_8px_#7DCFB6]"
              : "border-[#C3B299] hover:border-[#7DCFB6]"
            }
          `}
          onClick={() => handleColorSelect("black")}
        />
      </div>

      <div className="flex flex-row items-center justify-center gap-5">
        <p className="text-2xl text-[#C3B299] font-bold">Time:</p>

        {[
          { label: "Rapid (5 min)", timeControl: "Rapid" },
          { label: "Blitz (3 min)", timeControl: "Blitz" },
          { label: "Bullet (1 min)", timeControl: "Bullet" },
          { label: "Hyperbullet (30 sec)", timeControl: "Hyperbullet" },
        ].map(({ label, timeControl }) => (
          <p
            key={timeControl}
            className={`
              text-xl font-semibold cursor-pointer rounded-xl px-5 py-2 transition-all duration-300
              border-2 text-[#C3B299] bg-[#302e2e]
              ${
                lobbyRequest.data.timeControl === timeControl
                  ? "border-[#7DCFB6] shadow-[0_0_8px_#7DCFB6]"
                  : "border-[#C3B299] hover:border-[#7DCFB6] hover:bg-[#524b4b]"
              }
            `}
            onClick={() => handleTimeControlSelect(timeControl)}
          >
            {label}
          </p>
        ))}
      </div>

    </>
  )
}