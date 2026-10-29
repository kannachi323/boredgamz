import maple from "../../../assets/games/gomoku-maple.webp";
import oak from "../../../assets/games/gomoku-oak.webp";
import walnut from "../../../assets/games/gomoku-walnut.webp";

export type BoardTheme = "small" | "mid" | "large";
export const boardImages: Record<BoardTheme, string> = {
  small: maple,
  mid: walnut,
  large: oak,
};
