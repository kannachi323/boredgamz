import { Stone } from "../../types";

export function BoardStonePiece({ stone, isHovered }: { stone: Stone; isHovered: boolean | null}) {
  if (stone.color) {
    return (
      <img
        src={`/${stone.color}.svg`}
        alt={`${stone.color} stone`}
        className="h-full w-full opacity-100"
      />
    );
  }

  if (isHovered) {
    return (
      <img
        src={`/black.svg`}
        alt={`preview stone`}
        className="h-14 w-14 opacity-50"
      />
    );
  }

  return null;
}