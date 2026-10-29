import { ArrowRight, Circle, Sparkles, Wifi, Zap } from 'lucide-react'
import { Link } from 'react-router-dom'
import { GameCard } from '../../components/Cards/GameCard'
import { games } from '../../data/games'

export default function Home() {
  return (
    <>
      <section className="hero">
        <div className="hero__glow" aria-hidden="true" />
        <div className="hero__content">
          <div className="status-pill"><span /> Tables are open</div>
          <p className="eyebrow">The neighborhood game room, online</p>
          <h1>Your table<br />is <em>ready.</em></h1>
          <p className="hero__lede">Six great games. No downloads, no rulebooks, no waiting around. Pick a seat and start playing.</p>
          <div className="hero__actions">
            <Link className="button button--primary" to="/games">Choose a game <ArrowRight size={18} /></Link>
            <Link className="button button--ghost" to="/games/gomoku">Play Gomoku</Link>
          </div>
          <div className="hero__proof">
            <span><Zap size={17} /> Instant play</span>
            <span><Wifi size={17} /> Live multiplayer</span>
            <span><Sparkles size={17} /> Always free</span>
          </div>
        </div>
        <div className="hero-board" aria-label="Decorative game board">
          <div className="hero-board__label"><span>Table 06</span><span>YOUR TURN</span></div>
          <div className="hero-board__grid">
            {Array.from({ length: 25 }, (_, index) => (
              <span key={index} className={index === 6 || index === 12 || index === 18 ? 'is-one' : index === 8 || index === 13 ? 'is-two' : ''}>
                {(index === 6 || index === 12 || index === 18 || index === 8 || index === 13) && <Circle fill="currentColor" />}
              </span>
            ))}
          </div>
          <div className="hero-board__players"><span><i className="avatar avatar--coral">M</i> Maya</span><b>VS</b><span><i className="avatar avatar--lime">J</i> You</span></div>
        </div>
      </section>

      <section className="section" id="featured-games">
        <div className="section-heading">
          <div><p className="eyebrow">Pick your pace</p><h2>Six ways to make a move.</h2></div>
          <Link className="text-link" to="/games">View game room <ArrowRight size={17} /></Link>
        </div>
        <div className="game-grid">{games.map((game, index) => <GameCard key={game.slug} game={game} index={index} />)}</div>
      </section>

      <section className="club-callout">
        <div><p className="eyebrow">No velvet rope</p><h2>A game room for everybody.</h2></div>
        <p>Start locally on one screen, or match with another player online. The rules live on the table, so all you need to bring is a little competitive spirit.</p>
        <Link className="button button--light" to="/games">Take a seat <ArrowRight size={18} /></Link>
      </section>
    </>
  )
}
