import type { ArenaState } from '../../../hooks/useOnlineMatch'

export function ConnectFourBoard({ state, canMove, onMove }: { state: ArenaState; canMove: boolean; onMove: (column: number) => void }) {
  return (
    <div className="c4-board-wrap">
      <div className="c4-column-arrows" aria-hidden="true" style={{ gridTemplateColumns: `repeat(${state.columns}, 1fr)` }}>
        {Array.from({ length: state.columns }, (_, column) => <button key={column} tabIndex={-1} disabled={!canMove} onClick={() => onMove(column)}>↓</button>)}
      </div>
      <div className={`c4-board c4-board--${state.columns}`} role="grid" aria-label="Connect Four board" style={{ gridTemplateColumns: `repeat(${state.columns}, 1fr)` }}>
        {state.board.map((mark, index) => {
          const row = Math.floor(index / state.columns) + 1
          const column = index % state.columns
          return (
            <button
              key={index}
              className={`${mark ? `is-${mark}` : ''} ${state.lastMove === index ? 'is-last' : ''}`}
              type="button"
              role="gridcell"
              disabled={!canMove}
              onClick={() => onMove(column)}
              aria-label={`Column ${column + 1}, row ${row}${mark ? `, ${mark === 'one' ? 'coral disc' : 'gold disc'}` : ', empty'}`}
            ><span aria-hidden="true" /></button>
          )
        })}
      </div>
      <div className="c4-board-feet" aria-hidden="true"><i /><i /></div>
    </div>
  )
}
