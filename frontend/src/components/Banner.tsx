
import { Player } from "../types";


interface Props {
  player : Player;
}


export function PlayerBanner({ player } : Props) {
  return (
    <div className="bg-[#363430] text-white p-2 rounded-lg shadow-md w-full flex flex-row justify-between">
      <img src={player.color} alt="user's profile picture" className="h-16 w-16 bg-red-50 rounded-md"/>
      <h2 className="text-xl font-bold mb-2">{player.playerName}</h2>
    </div>
  );
}