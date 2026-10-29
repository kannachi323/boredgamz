import { describe, expect, it } from 'vitest'
import { gridConfigs, gridWinner, placeGridMove, type GridMark } from './gridGames'

describe('grid game rules', () => {
  it('detects a tic-tac-toe diagonal', () => {
    const board: GridMark[] = ['one', '', '', '', 'one', '', '', '', 'one']
    expect(gridWinner(board, gridConfigs.tictactoe, 8)).toBe('one')
  })

  it('rejects occupied grid spaces', () => {
    const board: GridMark[] = Array(9).fill('')
    board[4] = 'two'
    expect(placeGridMove(board, gridConfigs.tictactoe, 4, 'one')).toBeNull()
  })

  it('drops Connect Four discs and detects a vertical win', () => {
    let board: GridMark[] = Array(42).fill('')
    let finalPosition = -1
    for (let move = 0; move < 4; move += 1) {
      const result = placeGridMove(board, gridConfigs['connect-four'], 2, 'one')!
      board = result.board
      finalPosition = result.position
    }
    expect(finalPosition).toBe(16)
    expect(gridWinner(board, gridConfigs['connect-four'], finalPosition)).toBe('one')
  })

  it('does not allow a move in a full Connect Four column', () => {
    let board: GridMark[] = Array(42).fill('')
    for (let move = 0; move < 6; move += 1) board = placeGridMove(board, gridConfigs['connect-four'], 0, move % 2 ? 'one' : 'two')!.board
    expect(placeGridMove(board, gridConfigs['connect-four'], 0, 'one')).toBeNull()
  })
})
