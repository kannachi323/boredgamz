import { ArrowUpRight, Clock3, Users } from 'lucide-react'
import { Link } from 'react-router-dom'
import type { GameDefinition } from '../../data/games'

export function GameCard({ game, index = 0 }: { game: GameDefinition; index?: number }) {
  const Icon = game.icon
  return (
    <Link className="game-card" to={`/games/${game.slug}`} style={{ '--accent': game.accent, '--delay': `${index * 55}ms` } as React.CSSProperties}>
      <div className="game-card__top">
        <span className="game-card__icon"><Icon size={27} /></span>
        <ArrowUpRight className="game-card__arrow" size={22} />
      </div>
      <p className="eyebrow">{game.kicker}</p>
      <h3>{game.name}</h3>
      <p className="game-card__description">{game.description}</p>
      <div className="game-card__meta">
        <span><Users size={15} /> {game.players}</span>
        <span><Clock3 size={15} /> {game.duration}</span>
      </div>
      <div className="game-card__modes">{game.modes.map((mode) => <span key={mode}>{mode}</span>)}</div>
    </Link>
  )
}
