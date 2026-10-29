export type GridKind = 'tictactoe' | 'connect-four' | 'gomoku'
export type GridMark = '' | 'one' | 'two'

export interface GridConfig {
  kind: GridKind
  rows: number
  columns: number
  connect: number
}

export const gridConfigs: Record<GridKind, GridConfig> = {
  tictactoe: { kind: 'tictactoe', rows: 3, columns: 3, connect: 3 },
  'connect-four': { kind: 'connect-four', rows: 6, columns: 7, connect: 4 },
  gomoku: { kind: 'gomoku', rows: 15, columns: 15, connect: 5 },
}

export function placeGridMove(board: GridMark[], config: GridConfig, position: number, mark: Exclude<GridMark, ''>) {
  const next = [...board]
  let target = position
  if (config.kind === 'connect-four') {
    if (position < 0 || position >= config.columns) return null
    target = -1
    for (let row = config.rows - 1; row >= 0; row -= 1) {
      const candidate = row * config.columns + position
      if (!next[candidate]) {
        target = candidate
        break
      }
    }
  }
  if (target < 0 || target >= next.length || next[target]) return null
  next[target] = mark
  return { board: next, position: target }
}

export function gridWinner(board: GridMark[], config: GridConfig, position: number): GridMark {
  const mark = board[position]
  if (!mark) return ''
  const row = Math.floor(position / config.columns)
  const column = position % config.columns
  const directions = [[0, 1], [1, 0], [1, 1], [1, -1]]
  for (const [rowStep, columnStep] of directions) {
    let count = 1
    count += countLine(board, config, row, column, rowStep, columnStep, mark)
    count += countLine(board, config, row, column, -rowStep, -columnStep, mark)
    if (count >= config.connect) return mark
  }
  return ''
}

function countLine(board: GridMark[], config: GridConfig, startRow: number, startColumn: number, rowStep: number, columnStep: number, mark: GridMark) {
  let count = 0
  let row = startRow + rowStep
  let column = startColumn + columnStep
  while (row >= 0 && row < config.rows && column >= 0 && column < config.columns && board[row * config.columns + column] === mark) {
    count += 1
    row += rowStep
    column += columnStep
  }
  return count
}
