import { Stone } from "../../types";
import BLACK from "@/assets/black.svg"
import WHITE from "@/assets/white.svg"

export function GomokuStone({ stone, isHovered }: { stone: Stone; isHovered: boolean | null}) {
  if (stone.color) {
    return (
      <img
        src={stone.color === 'black' ? BLACK : WHITE}
        alt={`${stone.color} stone`}
        className="h-full w-full opacity-100"
      />
    );
  }

  if (isHovered) {
    return (
      <img
        src={BLACK}
        alt={`preview stone`}
        className="h-14 w-14 opacity-50"
      />
    );
  }

  return null;
}