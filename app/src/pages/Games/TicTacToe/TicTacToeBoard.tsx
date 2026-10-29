import type { ArenaState } from '../../../hooks/useOnlineMatch'

interface Props {
  state: ArenaState
  canMove: boolean
  onMove: (position: number) => void
}

export function TicTacToeBoard({ state, canMove, onMove }: Props) {
  return (
    <div className={`ttt-board ttt-board--${state.rows}`} style={{ gridTemplateColumns: `repeat(${state.columns}, 1fr)` }} role="grid" aria-label="Tic-Tac-Toe board">
      {state.board.map((mark, index) => {
        const row = Math.floor(index / state.columns) + 1
        const column = index % state.columns + 1
        const disabled = !canMove || Boolean(mark)
        return (
          <button
            className={`${mark ? `is-${mark}` : ''} ${state.lastMove === index ? 'is-last' : ''}`}
            disabled={disabled}
            key={index}
            onClick={() => onMove(index)}
            type="button"
            role="gridcell"
            aria-label={`Row ${row}, column ${column}${mark ? `, ${mark === 'one' ? 'cross' : 'circle'}` : ', empty'}`}
          >
            {mark && <span aria-hidden="true" />}
          </button>
        )
      })}
    </div>
  )
}
