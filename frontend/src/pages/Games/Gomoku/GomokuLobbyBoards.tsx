import { useGomokuStore } from "../../../stores/useGomokuStore"

import SMALL_BOARD from '../../../assets/small-board.jpg'
import MID_BOARD from '../../../assets/mid-board.jpg'
import LARGE_BOARD from '../../../assets/large-board.jpg'


const boards = [
  { size: "9x9", img: SMALL_BOARD },
  { size: "13x13", img: MID_BOARD },
  { size: "19x19", img: LARGE_BOARD },
];

export function GomokuLobbyBoards() {
  const { setConnection, player, handler } = useGomokuStore();

  return (
    <>
      {boards.map(({ size, img }) => (
        <div
          key={size}
          className="bg-[#302e2e] w-1/7  p-4 flex flex-col items-center justify-center gap-3
          rounded-lg border border-[#1b1918] hover:bg-[#524b4b] transition-colors duration-300 cursor-pointer"
          onClick={() => setConnection(size, player, handler)}
        >
          <p className="text-3xl text-[#C3B299]">{size}</p>
          <img src={img} alt="gomoku board" className="w-full h-auto rounded-md" />
        </div>
      ))}
    </>
  )
}