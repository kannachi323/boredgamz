import type { CSSProperties } from "react";
import { useEffect, useMemo, useState } from "react";
import type { MahjongState } from "../../../hooks/useMahjongMatch";

interface Props {
  state: MahjongState;
  theme: string;
  onPair: (first: number, second: number) => void;
}

export function MahjongBoard({ state, theme, onPair }: Props) {
  const [selected, setSelected] = useState<number | null>(null);
  const free = useMemo(() => new Set(state.free), [state.free]);
  const hinted = useMemo(() => new Set(state.hint), [state.hint]);
  useEffect(() => setSelected(null), [state.roomId, state.removed]);
  function select(index: number) {
    if (selected === index) {
      setSelected(null);
      return;
    }
    if (selected === null) {
      setSelected(index);
      return;
    }
    onPair(selected, index);
    setSelected(null);
  }
  return (
    <div className={`mjh-board is-${state.mode === "garden" ? "full" : "compact"} theme-${theme}`} aria-label="Mahjong tile layout">
      <div className="mjh-board__mat" />
      {state.tiles.map((tile, index) => {
        const removed = state.removed[index];
        const style = {
          "--tile-x": tile.x,
          "--tile-y": tile.y,
          "--tile-z": tile.z,
        } as CSSProperties;
        return (
          <button
            type="button"
            key={tile.id}
            style={style}
            className={`mjh-tile group-${faceKind(tile.face)} ${removed ? "is-removed" : ""} ${free.has(index) ? "is-free" : ""} ${selected === index ? "is-selected" : ""} ${hinted.has(index) ? "is-hinted" : ""}`}
            disabled={removed || !free.has(index) || state.status !== "playing" || state.viewerSeat < 0}
            onClick={() => select(index)}
            aria-label={`${tile.face} tile${free.has(index) ? ", free" : ", blocked"}`}
          >
            <span>{tile.face}</span>
            <small>{tile.group === "season" ? "季" : tile.group === "flower" ? "花" : ""}</small>
          </button>
        );
      })}
    </div>
  );
}

function faceKind(face: string) {
  if (face.endsWith("萬") || face === "中") return "red";
  if (face.endsWith("筒") || face === "白") return "blue";
  if (face.endsWith("索") || face === "發") return "green";
  if (["春", "夏", "秋", "冬", "梅", "蘭", "菊", "竹"].includes(face)) return "special";
  return "ink";
}
