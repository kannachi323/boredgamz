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
          className="bg-[#302e2e] aspect-square w-[256px] flex flex-col items-center justify-center
          rounded-lg border border-[#1b1918] hover:bg-[#524b4b] transition cursor-pointer p-2 gap-2"
          onClick={() => setConnection(size, player, handler)}
        >
          <p className="text-2xl text-[#C3B299]">{size}</p>
          <img src={img} alt="gomoku board" className="w-[256px] aspect-square rounded-md object-cover" />
        </div>
      ))}
    </>
  );

}