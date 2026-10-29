import type { ReactNode } from 'react'
import { useLocation } from 'react-router-dom'
import { NavBarV1 } from './components/NavBar/NavBarV1'
import { FooterV1 } from './features/Footer/FooterV1'

export function AppLayout({ children }: { children: ReactNode }) {
  const location = useLocation()
  const immersive = ['/games/tictactoe', '/games/connect-four', '/games/gomoku', '/games/chess', '/games/poker', '/games/mahjong'].includes(location.pathname)
  if (immersive) {
    return <div className="site-shell"><a className="skip-link" href="#main-content">Skip to content</a><div id="main-content">{children}</div></div>
  }
  return (
    <div className="site-shell">
      <a className="skip-link" href="#main-content">Skip to content</a>
      <NavBarV1 />
      <main id="main-content">{children}</main>
      <FooterV1 />
    </div>
  )
}
