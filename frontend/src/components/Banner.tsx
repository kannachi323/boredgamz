
import { Player } from "../pages/Games/Gomoku/GomokuTypes";


interface Props {
  player : Player;
}


export function PlayerBanner({ player } : Props) {
  return (
    <>
      <div className="flex flex-row justify-center items-center gap-2">
        <img src={player.color} alt="user's profile picture" className="h-16 w-16 rounded-md"/>
        <h2 className="text-xl font-bold self-start">{player.playerName}</h2>

      </div>
    </>
  );
}