import { GameCard } from '../../components/Cards/GameCard'
import { games } from '../../data/games'

export default function Games() {
  return (
    <div className="library-page">
      <header className="page-heading">
        <p className="eyebrow">Game room</p>
        <h1>What are we playing?</h1>
        <p>From two-minute tactics to a quiet tile-clearing session, there is a table ready for you.</p>
      </header>
      <div className="game-grid">{games.map((game, index) => <GameCard key={game.slug} game={game} index={index} />)}</div>
    </div>
  )
}
