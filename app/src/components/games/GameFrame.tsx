import { ArrowLeft, RotateCcw } from 'lucide-react'
import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import type { GameDefinition } from '../../data/games'

export function GameFrame({ game, status, children, aside, onReset }: { game: GameDefinition; status: string; children: ReactNode; aside?: ReactNode; onReset?: () => void }) {
  return (
    <div className="game-page" style={{ '--accent': game.accent } as React.CSSProperties}>
      <header className="game-header">
        <div>
          <Link className="back-link" to="/games"><ArrowLeft size={16} /> Game room</Link>
          <p className="eyebrow">{game.kicker}</p>
          <h1>{game.name}</h1>
        </div>
        <div className="game-header__status"><span /> {status}</div>
      </header>
      <div className={`game-stage ${aside ? '' : 'game-stage--centered'}`}>
        <section className="game-surface">{children}</section>
        {aside && <aside className="game-sidebar">{aside}</aside>}
      </div>
      {onReset && <button className="button button--ghost game-reset" type="button" onClick={onReset}><RotateCcw size={17} /> New game</button>}
    </div>
  )
}
