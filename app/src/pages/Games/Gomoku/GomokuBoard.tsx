import type { ArenaState } from "../../../hooks/useOnlineMatch";
import type { BoardTheme } from "./boardThemes";

export function GomokuBoard({
  state,
  theme,
  canMove,
  onMove,
}: {
  state: ArenaState;
  theme: BoardTheme;
  canMove: boolean;
  onMove: (position: number) => void;
}) {
  return (
    <div
      className={`gmk-board gmk-board--${theme}`}
      style={{ gridTemplateColumns: `repeat(${state.columns}, 1fr)` }}
      role="grid"
      aria-label="Gomoku board"
    >
      {state.board.map((mark, index) => (
        <button
          key={index}
          type="button"
          role="gridcell"
          disabled={!canMove || Boolean(mark)}
          className={`${mark ? `is-${mark}` : ""} ${state.lastMove === index ? "is-last" : ""}`}
          onClick={() => onMove(index)}
          aria-label={`${String.fromCharCode(65 + (index % state.columns))}${Math.floor(index / state.columns) + 1}${mark ? `, ${mark === "one" ? "black" : "white"} stone` : ", empty"}`}
        >
          <span />
        </button>
      ))}
    </div>
  );
}
