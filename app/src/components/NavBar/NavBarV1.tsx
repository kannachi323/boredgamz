import { Dices, LogIn, LogOut, Menu, UserRoundPlus, X } from 'lucide-react'
import { useState } from 'react'
import { Link, NavLink } from 'react-router-dom'
import { useAuthStore } from '../../stores/useAuthStore'

const links = [
  { to: '/', label: 'Home' },
  { to: '/games', label: 'Game room' },
]

export function NavBarV1() {
  const [open, setOpen] = useState(false)
  const status = useAuthStore((state) => state.status)
  const user = useAuthStore((state) => state.user)
  const logout = useAuthStore((state) => state.logout)
  return (
    <header className="topbar">
      <div className="topbar__inner">
        <Link className="brand" to="/" aria-label="BoredGamz home" onClick={() => setOpen(false)}>
          <span className="brand__mark"><Dices size={20} strokeWidth={2.4} /></span>
          <span>Bored<span>Gamz</span></span>
        </Link>
        <nav className={`nav ${open ? 'nav--open' : ''}`} aria-label="Primary navigation">
          {links.map((link) => (
            <NavLink key={link.to} to={link.to} end={link.to === '/'} onClick={() => setOpen(false)}>
              {link.label}
            </NavLink>
          ))}
          {status === 'authenticated' && user ? <>
            <span className="nav__player">Playing as <b>{user.username}</b></span>
            <button className="nav__auth" type="button" onClick={() => { setOpen(false); void logout() }}><LogOut size={15} /> Sign out</button>
          </> : <>
            <Link className="nav__auth" to="/login" onClick={() => setOpen(false)}><LogIn size={15} /> Sign in</Link>
            <Link className="button button--small button--primary" to="/signup" onClick={() => setOpen(false)}><UserRoundPlus size={15} /> Join</Link>
          </>}
        </nav>
        <button className="nav-toggle" type="button" aria-expanded={open} aria-label="Toggle menu" onClick={() => setOpen((value) => !value)}>
          {open ? <X /> : <Menu />}
        </button>
      </div>
    </header>
  )
}
