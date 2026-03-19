import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore"

import SMALL_BOARD from '@/assets/small-board.jpg'
import MID_BOARD from '@/assets/mid-board.jpg'
import LARGE_BOARD from '@/assets/large-board.jpg'


const boards = [
  { label: "9x9", img: SMALL_BOARD },
  { label: "13x13", img: MID_BOARD },
  { label: "19x19", img: LARGE_BOARD },
];

export function GomokuLobbyBoards() {
  const { lobbyRequest, setLobbyRequest } = useGomokuStore();

  return (
    <>
      {boards.map(({ img, label }) => (
        <div
          key={label}
          className={`bg-[#302e2e] aspect-square w-[256px] flex flex-col items-center justify-center
          rounded-lg border border-[#1b1918] hover:bg-[#524b4b] transition cursor-pointer p-2 gap-2
          ${lobbyRequest.data.name === label
            ? "border-[#7DCFB6] shadow-[0_0_8px_#7DCFB6]"
            : "border-[#C3B299] hover:border-[#7DCFB6] hover:bg-[#524b4b]"
          }
          `}
          onClick={() => setLobbyRequest({...lobbyRequest, data: { ...lobbyRequest.data, name: label}})}
        >
          <p className="text-2xl text-[#C3B299]">{label}</p>
          <img src={img} alt="gomoku board" className="w-[256px] aspect-square rounded-md object-cover" />
        </div>
      ))}
    </>
  );

}