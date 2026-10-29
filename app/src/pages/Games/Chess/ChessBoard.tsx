import { Chess, type Square } from 'chess.js'
import { useMemo, useState } from 'react'
import type { ArenaState } from '../../../hooks/useOnlineMatch'
import { chessThemes, type ChessTheme } from './chessThemes'

const glyphs: Record<string, string> = { wp: '♙', wn: '♘', wb: '♗', wr: '♖', wq: '♕', wk: '♔', bp: '♟', bn: '♞', bb: '♝', br: '♜', bq: '♛', bk: '♚' }
type Pending = { from: string; to: string }

export function ChessBoard({ state, theme, myMark, onMove }: { state: ArenaState<string>; theme: ChessTheme; myMark?: string; onMove: (from: string, to: string, promotion?: string) => void }) {
  const chess = useMemo(() => new Chess(state.chess?.fen ?? undefined), [state.chess?.fen])
  const [selected, setSelected] = useState<Square | null>(null)
  const [pending, setPending] = useState<Pending | null>(null)
  const legal = selected ? chess.moves({ square: selected, verbose: true }) : []
  const flipped = myMark === 'two'
  const display = Array.from({ length: 64 }, (_, shown) => flipped ? 63 - shown : shown)
  const colors = chessThemes[theme]

  function choose(index: number) {
    if (state.status !== 'playing' || state.turn !== myMark) return
    const square = indexToSquare(index)
    const candidates = legal.filter((move) => move.to === square)
    if (selected && candidates.length) {
      if (candidates.some((move) => move.promotion)) setPending({ from: selected, to: square })
      else onMove(selected, square)
      setSelected(null)
      return
    }
    const piece = chess.get(square)
    const expected = myMark === 'one' ? 'w' : 'b'
    setSelected(piece?.color === expected ? square : null)
  }

  return <div className="chs-board-shell" style={{ '--chs-light': colors.light, '--chs-dark': colors.dark, '--chs-accent': colors.accent } as React.CSSProperties}>
    <div className="chs-board" role="grid" aria-label={`Chess board${flipped ? ', black perspective' : ', white perspective'}`}>
      {display.map((index, shown) => { const code = state.board[index]; const square = indexToSquare(index); const row = Math.floor(index / 8); const column = index % 8; const target = legal.some((move) => move.to === square); const light = (row + column) % 2 === 0; return <button key={square} type="button" role="gridcell" aria-label={`${square}${code ? ` ${code[0] === 'w' ? 'white' : 'black'} ${pieceName(code[1])}` : ', empty'}`} className={`${light ? 'is-light' : 'is-dark'} ${selected === square ? 'is-selected' : ''} ${target ? 'is-legal' : ''} ${state.chess?.lastFrom === square || state.chess?.lastTo === square ? 'is-last' : ''}`} onClick={() => choose(index)}>{code && <span className={`chs-piece is-${code[0]}`}>{glyphs[code]}</span>}{shown % 8 === 0 && <small className="chs-rank">{8 - row}</small>}{shown >= 56 && <small className="chs-file">{String.fromCharCode(97 + column)}</small>}</button> })}
    </div>
    {pending && <div className="chs-promotion" role="dialog" aria-label="Choose promotion piece"><p>Promote pawn to</p><div>{(['q','r','b','n'] as const).map((piece) => <button key={piece} type="button" onClick={() => { onMove(pending.from, pending.to, piece); setPending(null) }}>{glyphs[`${myMark === 'one' ? 'w' : 'b'}${piece}`]}</button>)}</div><button type="button" onClick={() => setPending(null)}>Cancel</button></div>}
  </div>
}

function indexToSquare(index: number) { return `${String.fromCharCode(97 + index % 8)}${8 - Math.floor(index / 8)}` as Square }
function pieceName(code?: string) { return ({ p:'pawn', n:'knight', b:'bishop', r:'rook', q:'queen', k:'king' } as Record<string,string>)[code ?? ''] ?? 'piece' }
