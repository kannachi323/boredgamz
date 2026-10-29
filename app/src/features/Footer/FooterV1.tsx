import { Dices, Github } from 'lucide-react'
import { Link } from 'react-router-dom'

export function FooterV1() {
  return (
    <footer className="footer">
      <div className="footer__inner">
        <div>
          <Link className="brand" to="/"><span className="brand__mark"><Dices size={18} /></span><span>Bored<span>Gamz</span></span></Link>
          <p>Six familiar games, one welcoming table.</p>
        </div>
        <div className="footer__links">
          <Link to="/games">All games</Link>
          <a href="https://github.com/kannachi323/playgomoku" aria-label="BoredGamz on GitHub"><Github size={18} /></a>
        </div>
      </div>
      <div className="footer__legal">© {new Date().getFullYear()} BoredGamz. Built for friendly competition.</div>
    </footer>
  )
}
