import type { ReactNode } from "react";
import type { GameSlug, ModeGlyph } from "./catalog";

/* Hand-drawn SVG board illustrations. Colours come from CSS (.art-*) so the
   art follows light/dark themes. */

const arts: Record<GameSlug, () => ReactNode> = {
  tictactoe: TicTacToeArt,
  "connect-four": ConnectFourArt,
  gomoku: GomokuArt,
  chess: ChessArt,
  poker: PokerArt,
  mahjong: MahjongArt,
};

export function BoardArt({ slug, className }: { slug: GameSlug; className?: string }) {
  const Art = arts[slug];
  return (
    <svg className={`art ${className ?? ""}`} viewBox="0 0 400 300" preserveAspectRatio="xMidYMid slice" aria-hidden="true" focusable="false">
      <Art />
    </svg>
  );
}

function Cross({ x, y, size = 18 }: { x: number; y: number; size?: number }) {
  return <path className="art-stroke-ink" d={`M${x - size} ${y - size}L${x + size} ${y + size}M${x + size} ${y - size}L${x - size} ${y + size}`} strokeWidth={9} />;
}

function TicTacToeArt() {
  const cells = [140, 200, 260];
  return (
    <g transform="rotate(-6 200 150)">
      <rect className="art-paper" x="96" y="46" width="208" height="208" rx="22" />
      <path className="art-stroke-line" strokeWidth={6} d="M170 66V234M230 66V234M116 120H284M116 180H284" />
      <Cross x={cells[0]} y={90} />
      <Cross x={cells[1]} y={90} />
      <Cross x={cells[2]} y={150} />
      {[[cells[2], 90], [cells[1], 150], [cells[0], 210]].map(([x, y]) => (
        <circle key={`${x}-${y}`} className="art-stroke-accent" cx={x} cy={y} r={19} strokeWidth={9} />
      ))}
      <path className="art-stroke-accent" strokeWidth={7} d="M282 68L118 232" opacity={0.9} />
    </g>
  );
}

function ConnectFourArt() {
  const filled: Record<string, "p1" | "p2"> = {
    "5-1": "p1", "5-2": "p2", "5-3": "p1", "5-4": "p2", "5-5": "p1",
    "4-2": "p1", "4-3": "p2", "4-4": "p2", "3-3": "p1", "3-4": "p2", "2-4": "p1",
  };
  const holes = [];
  for (let row = 0; row < 6; row++) {
    for (let column = 0; column < 7; column++) {
      const fill = filled[`${row}-${column}`];
      holes.push(<circle key={`${row}-${column}`} className={fill ? `art-${fill}` : "art-hole"} cx={110 + column * 30} cy={98 + row * 28} r={11} />);
    }
  }
  return (
    <g>
      <rect className="art-frame" x="88" y="76" width="224" height="186" rx="20" />
      {holes}
      <rect className="art-frame" x="96" y="258" width="16" height="22" rx="4" />
      <rect className="art-frame" x="288" y="258" width="16" height="22" rx="4" />
      <g className="art-drop">
        <path className="art-stroke-line" strokeWidth={3} d="M250 20V36M260 14V34M240 14V34" opacity={0.6} />
        <circle className="art-p2" cx="250" cy="52" r="12" />
      </g>
    </g>
  );
}

function GomokuArt() {
  const step = 21;
  const lines = [];
  for (let i = 1; i < 10; i++) lines.push(`M${95 + i * step} 66V234M116 ${45 + i * step}H284`);
  const black = [[2, 7], [3, 6], [4, 5], [5, 4], [6, 3]];
  const white = [[4, 4], [5, 6], [3, 5], [6, 5], [4, 7]];
  const at = (n: number) => 95 + n * step;
  return (
    <g transform="rotate(4 200 150)">
      <rect className="art-wood" x="95" y="45" width="210" height="210" rx="12" />
      <path className="art-wood-line" strokeWidth={1.4} d={lines.join("")} />
      {white.map(([c, r]) => <circle key={`w${c}${r}`} className="art-stone-white" cx={at(c)} cy={at(r) - 50} r={9.5} />)}
      {black.map(([c, r]) => <circle key={`b${c}${r}`} className="art-stone-black" cx={at(c)} cy={at(r) - 50} r={9.5} />)}
      <circle className="art-stroke-accent" cx={at(6)} cy={at(3) - 50} r={14} strokeWidth={2.5} />
    </g>
  );
}

function Pawn({ x, y, className }: { x: number; y: number; className: string }) {
  return (
    <g transform={`translate(${x} ${y})`} className={className}>
      <rect x="-22" y="-11" width="44" height="11" rx="4" />
      <path d="M-14 -11C-12 -28 -7 -37 -7 -41H7C7 -37 12 -28 14 -11Z" />
      <rect x="-13" y="-47" width="26" height="7" rx="3.5" />
      <circle cy="-58" r="12" />
    </g>
  );
}

function King({ x, y, className }: { x: number; y: number; className: string }) {
  return (
    <g transform={`translate(${x} ${y})`} className={className}>
      <rect x="-26" y="-13" width="52" height="13" rx="5" />
      <path d="M-19 -13C-17 -38 -11 -58 -11 -66H11C11 -58 17 -38 19 -13Z" />
      <rect x="-18" y="-74" width="36" height="8" rx="4" />
      <path d="M-15 -74C-17 -87 -9 -93 0 -93S17 -87 15 -74Z" />
      <rect x="-3" y="-110" width="6" height="19" rx="2" />
      <rect x="-9" y="-104" width="18" height="6" rx="2" />
    </g>
  );
}

function ChessArt() {
  const squares = [];
  for (let row = 0; row < 4; row++) {
    for (let column = 0; column < 4; column++) {
      squares.push(<rect key={`${row}${column}`} className={(row + column) % 2 ? "art-sq-dark" : "art-sq-light"} x={96 + column * 52} y={46 + row * 52} width={52} height={52} />);
    }
  }
  return (
    <g>
      <rect className="art-sq-dark" x="90" y="40" width="220" height="220" rx="16" />
      <g>{squares}</g>
      <rect className="art-accent-soft" x="200" y="150" width="52" height="52" />
      <rect className="art-accent-soft" x="200" y="46" width="52" height="52" />
      <circle className="art-dot" cx="174" cy="176" r="7" />
      <King x={174} y={140} className="art-piece-dark" />
      <Pawn x={226} y={194} className="art-piece-light" />
      <Pawn x={278} y={90} className="art-piece-dark" />
    </g>
  );
}

const spade = "M0 -22C11 -9 24 0 24 11C24 22 13 26 4 19L8 31H-8L-4 19C-13 26 -24 22 -24 11C-24 0 -11 -9 0 -22Z";
const heart = "M0 25C-27 6 -25 -15 -12 -19C-5 -21 0 -15 0 -10C0 -15 5 -21 12 -19C25 -15 27 6 0 25Z";

function Card({ x, y, angle, rank, suit, red }: { x: number; y: number; angle: number; rank: string; suit: string; red?: boolean }) {
  return (
    <g transform={`translate(${x} ${y}) rotate(${angle})`}>
      <rect className="art-card" x="-56" y="-78" width="112" height="156" rx="12" />
      <text className={red ? "art-text-accent" : "art-text-ink"} x="-42" y="-50" fontSize="26" fontWeight="650">{rank}</text>
      <path className={red ? "art-accent" : "art-ink"} d={suit} transform="translate(0 6) scale(1.25)" />
    </g>
  );
}

function PokerArt() {
  return (
    <g>
      <ellipse className="art-felt" cx="200" cy="262" rx="230" ry="70" />
      <Card x={170} y={140} angle={-11} rank="A" suit={spade} />
      <Card x={232} y={146} angle={9} rank="K" suit={heart} red />
      <g transform="translate(316 232)">
        {[0, 1, 2].map((i) => (
          <g key={i} transform={`translate(0 ${-i * 9})`}>
            <ellipse className="art-chip" rx="26" ry="10" cy="4" />
            <ellipse className="art-accent" rx="26" ry="10" />
            <ellipse className="art-chip-ring" rx="17" ry="6" strokeWidth={2.5} strokeDasharray="5 5" />
          </g>
        ))}
      </g>
      <g transform="translate(88 238)">
        <ellipse className="art-chip" rx="24" ry="9" cy="4" />
        <ellipse className="art-ink" rx="24" ry="9" />
        <ellipse className="art-chip-ring-paper" rx="15" ry="5.5" strokeWidth={2.5} strokeDasharray="5 5" />
      </g>
    </g>
  );
}

function Tile({ x, y, children }: { x: number; y: number; children?: ReactNode }) {
  return (
    <g transform={`translate(${x} ${y})`}>
      <rect className="art-tile-side" x="5" y="7" width="56" height="74" rx="9" />
      <rect className="art-tile" width="56" height="74" rx="9" />
      {children}
    </g>
  );
}

function MahjongArt() {
  const dots = (points: number[][], cls = "art-ink") => points.map(([cx, cy], i) => <circle key={i} className={cls} cx={cx} cy={cy} r={6} />);
  const bamboo = (xs: number[]) => xs.map((x) => <rect key={x} className="art-bamboo" x={x} y="14" width="7" height="46" rx="3.5" />);
  return (
    <g>
      <Tile x={84} y={150}>{dots([[28, 37]], "art-accent")}</Tile>
      <Tile x={144} y={150}>{bamboo([16, 33])}</Tile>
      <Tile x={204} y={150}>{dots([[16, 18], [28, 37], [40, 56]])}</Tile>
      <Tile x={264} y={150}>{dots([[18, 22], [38, 22], [18, 52], [38, 52]], "art-accent")}</Tile>
      <Tile x={124} y={88}>{bamboo([12, 25, 38])}</Tile>
      <Tile x={184} y={88}>{dots([[18, 26], [38, 26]], "art-accent")}{dots([[28, 50]])}</Tile>
      <Tile x={244} y={88}>{dots([[28, 37]])}</Tile>
      <g className="art-lift">
        <Tile x={172} y={20}>
          <circle className="art-stroke-accent" cx="28" cy="37" r="15" strokeWidth={5} />
          <circle className="art-accent" cx="28" cy="37" r="5" />
        </Tile>
      </g>
    </g>
  );
}

/* Small glyphs used on mode buttons. */
export function ModeIcon({ glyph }: { glyph: ModeGlyph }) {
  if (glyph.kind === "grid") {
    const size = 30;
    const cell = size / Math.max(glyph.rows, glyph.columns);
    const lines = [];
    for (let c = 1; c < glyph.columns; c++) lines.push(`M${c * cell} 0V${glyph.rows * cell}`);
    for (let r = 1; r < glyph.rows; r++) lines.push(`M0 ${r * cell}H${glyph.columns * cell}`);
    return (
      <svg viewBox={`-2 -2 ${size + 4} ${size + 4}`} aria-hidden="true">
        <rect x="0" y="0" width={glyph.columns * cell} height={glyph.rows * cell} rx="3" fill="none" stroke="currentColor" strokeWidth="1.6" />
        <path d={lines.join("")} stroke="currentColor" strokeWidth="1" opacity="0.55" />
        {glyph.strike && <path d={`M2 ${glyph.rows * cell - 2}L${glyph.columns * cell - 2} 2`} stroke="var(--accent)" strokeWidth="2.4" strokeLinecap="round" />}
      </svg>
    );
  }
  if (glyph.kind === "clock") {
    const angle = (glyph.minutes / 15) * 2 * Math.PI;
    const x = 15 + 12 * Math.sin(angle);
    const y = 15 - 12 * Math.cos(angle);
    const large = glyph.minutes > 7.5 ? 1 : 0;
    return (
      <svg viewBox="0 0 30 30" aria-hidden="true">
        <circle cx="15" cy="15" r="12" fill="none" stroke="currentColor" strokeWidth="1.6" opacity="0.4" />
        <path d={glyph.minutes >= 15 ? "M15 3A12 12 0 1 1 14.99 3Z" : `M15 15L15 3A12 12 0 ${large} 1 ${x} ${y}Z`} fill="var(--accent)" />
      </svg>
    );
  }
  if (glyph.kind === "cards") {
    return (
      <svg viewBox="0 0 30 30" aria-hidden="true">
        {Array.from({ length: glyph.count }, (_, i) => {
          const spread = glyph.count === 2 ? 14 : 9;
          const angle = (i - (glyph.count - 1) / 2) * spread;
          return <rect key={i} x="10" y="5" width="11" height="16" rx="2" fill={i === glyph.count - 1 ? "var(--accent)" : "var(--surface)"} stroke="currentColor" strokeWidth="1.3" transform={`rotate(${angle} 15 26)`} />;
        })}
      </svg>
    );
  }
  return (
    <svg viewBox="0 0 30 30" aria-hidden="true">
      {Array.from({ length: glyph.layers }, (_, i) => (
        <rect key={i} x={8 - i * 2} y={9 - i * 3 + (glyph.layers - 1) * 3} width="14" height="18" rx="2.5" fill={i === glyph.layers - 1 ? "var(--surface)" : "var(--surface-3)"} stroke="currentColor" strokeWidth="1.3" transform={`translate(${i * 4} ${-i * 3})`} />
      ))}
    </svg>
  );
}
