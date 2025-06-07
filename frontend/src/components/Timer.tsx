import { useLocalTimer } from "../hooks/useLocalTimer";

export function Timer({ seconds } : { seconds: number }) {
  const timer = useLocalTimer(seconds)

  function formatTimer(timer : number) {
    const minutes = Math.floor(timer / 60)
    const seconds = timer % 60

    return `${minutes}:${seconds.toString().padStart(2, '0')}`
  }

  return (
    <div className="bg-[#363430] text-white w-1/3 p-2 flex flex-row justify-center items-center rounded-lg">
      <b className="text-3xl">{formatTimer(timer)}</b>
    </div>
  )
}