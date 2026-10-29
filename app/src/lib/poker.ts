export type Suit = '♠' | '♥' | '♦' | '♣'
export type Rank = '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | '10' | 'J' | 'Q' | 'K' | 'A'
export interface Card { suit: Suit; rank: Rank; id: string }
export interface HandScore { category: number; name: string; values: number[] }

const suits: Suit[] = ['♠', '♥', '♦', '♣']
const ranks: Rank[] = ['2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K', 'A']

export function createDeck(random: () => number = Math.random) {
  const deck = suits.flatMap((suit) => ranks.map((rank) => ({ suit, rank, id: `${rank}${suit}` })))
  for (let index = deck.length - 1; index > 0; index -= 1) {
    const swapIndex = Math.floor(random() * (index + 1))
    ;[deck[index], deck[swapIndex]] = [deck[swapIndex]!, deck[index]!]
  }
  return deck
}

export function evaluateHand(hand: Card[]): HandScore {
  const values = hand.map((card) => ranks.indexOf(card.rank) + 2).sort((a, b) => b - a)
  const counts = new Map<number, number>()
  values.forEach((value) => counts.set(value, (counts.get(value) ?? 0) + 1))
  const groups = [...counts.entries()].sort((a, b) => b[1] - a[1] || b[0] - a[0])
  const flush = hand.every((card) => card.suit === hand[0]?.suit)
  const unique = [...new Set(values)]
  const wheel = unique.join(',') === '14,5,4,3,2'
  const straight = unique.length === 5 && (values[0]! - values[4]! === 4 || wheel)
  const straightHigh = wheel ? 5 : values[0]!

  if (straight && flush) return { category: 8, name: 'Straight flush', values: [straightHigh] }
  if (groups[0]?.[1] === 4) return { category: 7, name: 'Four of a kind', values: groups.flatMap(([value, count]) => Array(count).fill(value)) }
  if (groups[0]?.[1] === 3 && groups[1]?.[1] === 2) return { category: 6, name: 'Full house', values: [groups[0][0], groups[1][0]] }
  if (flush) return { category: 5, name: 'Flush', values }
  if (straight) return { category: 4, name: 'Straight', values: [straightHigh] }
  if (groups[0]?.[1] === 3) return { category: 3, name: 'Three of a kind', values: groups.flatMap(([value, count]) => Array(count).fill(value)) }
  if (groups[0]?.[1] === 2 && groups[1]?.[1] === 2) return { category: 2, name: 'Two pair', values: groups.flatMap(([value, count]) => Array(count).fill(value)) }
  if (groups[0]?.[1] === 2) return { category: 1, name: 'One pair', values: groups.flatMap(([value, count]) => Array(count).fill(value)) }
  return { category: 0, name: 'High card', values }
}

export function compareHands(left: Card[], right: Card[]) {
  const a = evaluateHand(left)
  const b = evaluateHand(right)
  if (a.category !== b.category) return a.category - b.category
  for (let index = 0; index < Math.max(a.values.length, b.values.length); index += 1) {
    if ((a.values[index] ?? 0) !== (b.values[index] ?? 0)) return (a.values[index] ?? 0) - (b.values[index] ?? 0)
  }
  return 0
}

export function botHoldIndexes(hand: Card[]) {
  const score = evaluateHand(hand)
  const values = hand.map((card) => ranks.indexOf(card.rank) + 2)
  if (score.category >= 4) return [0, 1, 2, 3, 4]
  const counts = new Map<number, number>()
  values.forEach((value) => counts.set(value, (counts.get(value) ?? 0) + 1))
  return values.map((value, index) => ((counts.get(value) ?? 0) > 1 || value >= 12 ? index : -1)).filter((index) => index >= 0)
}
