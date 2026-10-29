import { describe, expect, it } from 'vitest'
import { compareHands, createDeck, evaluateHand, type Card, type Rank, type Suit } from './poker'

const hand = (cards: Array<[Rank, Suit]>): Card[] => cards.map(([rank, suit]) => ({ rank, suit, id: `${rank}${suit}` }))

describe('poker rules', () => {
  it('builds a unique 52-card deck', () => {
    const deck = createDeck(() => 0.5)
    expect(deck).toHaveLength(52)
    expect(new Set(deck.map((card) => card.id)).size).toBe(52)
  })

  it('recognizes an ace-low straight', () => {
    const score = evaluateHand(hand([['A', '♠'], ['2', '♥'], ['3', '♦'], ['4', '♣'], ['5', '♠']]))
    expect(score.name).toBe('Straight')
    expect(score.values).toEqual([5])
  })

  it('ranks a full house above a flush', () => {
    const fullHouse = hand([['K', '♠'], ['K', '♥'], ['K', '♦'], ['2', '♣'], ['2', '♠']])
    const flush = hand([['A', '♥'], ['J', '♥'], ['8', '♥'], ['5', '♥'], ['3', '♥']])
    expect(compareHands(fullHouse, flush)).toBeGreaterThan(0)
  })

  it('uses kickers to break equal pairs', () => {
    const aceKicker = hand([['Q', '♠'], ['Q', '♥'], ['A', '♦'], ['6', '♣'], ['2', '♠']])
    const kingKicker = hand([['Q', '♦'], ['Q', '♣'], ['K', '♦'], ['6', '♥'], ['2', '♥']])
    expect(compareHands(aceKicker, kingKicker)).toBeGreaterThan(0)
  })
})
