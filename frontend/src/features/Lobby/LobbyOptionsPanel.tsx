import { useGameStore } from "../../stores/useGameStore";

export function LobbyOptionsPanel() {
  const { player, setPlayer } = useGameStore();

  if (player == null) return

  return (
    <>
      <div className="flex flex-row items-center justify-evenly gap-2">
        <p className="text-2xl">Color:</p>
        <img
          src="/white.svg"
          alt="preview stone"
          className={`h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${player.color === "white" && " border-green-300"}`}
          onClick={() => setPlayer({ ...player, color: "white" })}
        />
        <img
          src={`/black.svg`}
          alt={`preview stone`}
          className={`h-14 w-14 rounded-full cursor-pointer border-2 transition-all duration-300
            ${player.color === "black" && " border-green-300"}`}
          onClick={() => setPlayer({...player, color: "black" })}
        />
      </div>
    </>
  )
}