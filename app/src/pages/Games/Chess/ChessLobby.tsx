import { Crown, LockKeyhole, ShieldCheck, Timer } from 'lucide-react'
import { useState } from 'react'
import chessArt from '../../../assets/chess.png'
import type { MatchOptions } from '../../../hooks/useOnlineMatch'
import { chessThemes, type ChessTheme } from './chessThemes'

const formats = [
  { id: 'standard', name: 'Classical', clock: '15m+10s', detail: '15 min + 10 sec' },
  { id: 'rapid', name: 'Rapid', clock: '10m', detail: '10 minutes' },
  { id: 'blitz', name: 'Blitz', clock: '3m+2s', detail: '3 min + 2 sec' },
] as const

export function ChessLobby({ name, onJoin }: { name: string; onJoin: (options: MatchOptions, theme: ChessTheme) => void }) {
  const [format, setFormat] = useState<(typeof formats)[number]>(formats[1])
  const [theme, setTheme] = useState<ChessTheme>(() => (localStorage.getItem('chess-theme') as ChessTheme) || 'walnut')
  const [privateTable, setPrivateTable] = useState(false)
  return <main className="chs-lobby">
    <section className="chs-lobby__hero"><a href="/games">← All tables</a><div><span><Crown size={15} /> BoredGamz Chess Club</span><h1>A better room<br />for the <em>royal game.</em></h1><p>Server-validated rules, real clocks, thoughtful board themes, spectators, and conversation—without the tournament-site clutter.</p></div><img src={chessArt} alt="Brown chess board with pieces arranged for a game" /></section>
    <form className="chs-setup" onSubmit={(event) => { event.preventDefault(); const data = new FormData(event.currentTarget); const requested = String(data.get('inviteCode') ?? '').trim().toUpperCase(); localStorage.setItem('chess-theme', theme); onJoin({ name: String(data.get('name') ?? name), mode: format.id, timeControl: format.clock, bestOf: 1, inviteCode: privateTable ? requested || `MATE${Math.random().toString(36).slice(2, 6).toUpperCase()}` : undefined }, theme) }}>
      <header><small>New match</small><h2>Choose your table</h2></header>
      <fieldset><legend><Timer size={14} /> Time format</legend><div className="chs-formats">{formats.map((item) => <button key={item.id} type="button" className={format.id === item.id ? 'is-selected' : ''} aria-pressed={format.id === item.id} onClick={() => setFormat(item)}><b>{item.name}</b><span>{item.detail}</span></button>)}</div></fieldset>
      <fieldset><legend>Board colors</legend><div className="chs-theme-picker">{(Object.keys(chessThemes) as ChessTheme[]).map((key) => <button key={key} type="button" className={theme === key ? 'is-selected' : ''} aria-pressed={theme === key} onClick={() => setTheme(key)}><i style={{ background: `linear-gradient(135deg,${chessThemes[key].light} 50%,${chessThemes[key].dark} 50%)` }} /><span>{chessThemes[key].name}</span></button>)}</div></fieldset>
      <label className="chs-field"><span>Display name</span><input name="name" defaultValue={name} maxLength={24} autoComplete="nickname" /></label>
      <label className="chs-private"><input type="checkbox" checked={privateTable} onChange={(e) => setPrivateTable(e.target.checked)} /><span><LockKeyhole size={16} /><b>Private challenge</b><small>Match only with the same code</small></span></label>
      {privateTable && <label className="chs-field"><span>Challenge code</span><input name="inviteCode" maxLength={12} placeholder="MATE24" /></label>}
      <button className="chs-find" type="submit"><ShieldCheck size={17} /> Find an opponent</button><p>Legal moves, castling, en passant, promotion, checkmate, repetition, and material draws are enforced by the server.</p>
    </form>
  </main>
}
