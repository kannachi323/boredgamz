import { Gauge, LockKeyhole, Maximize2, Sparkles, Timer, Zap } from 'lucide-react'
import { useState } from 'react'
import boardArt from '../../../assets/connect4.png'
import type { MatchOptions } from '../../../hooks/useOnlineMatch'

const variants = [
  { id: 'classic', name: 'Classic Drop', detail: '6 × 7 · connect four', icon: Sparkles },
  { id: 'quick', name: 'Quick Stack', detail: '5 × 6 · connect four', icon: Zap },
  { id: 'connect-five', name: 'Long Game', detail: '7 × 8 · connect five', icon: Maximize2 },
] as const

export function ConnectFourLobby({ name, onJoin }: { name: string; onJoin: (options: MatchOptions) => void }) {
  const [mode, setMode] = useState('classic')
  const [bestOf, setBestOf] = useState(3)
  const [timeControl, setTimeControl] = useState('30s')
  const [privateTable, setPrivateTable] = useState(false)

  return (
    <main className="c4-lobby">
      <section className="c4-lobby__hero">
        <a href="/games">← Game room</a>
        <div className="c4-lobby__copy">
          <span><Gauge size={15} /> Live vertical tactics</span>
          <h1>Every drop<br />changes <em>everything.</em></h1>
          <p>Shape two threats at once, control the center, and listen for the satisfying clack of a plan coming together.</p>
        </div>
        <div className="c4-lobby__art" aria-hidden="true">
          <img src={boardArt} alt="" />
          <i className="is-red" /><i className="is-yellow" /><i className="is-red" /><i className="is-yellow" /><i className="is-red" />
        </div>
      </section>

      <form className="c4-setup" onSubmit={(event) => {
        event.preventDefault()
        const data = new FormData(event.currentTarget)
        const requested = String(data.get('inviteCode') ?? '').trim().toUpperCase()
        onJoin({ name: String(data.get('name') ?? name), mode, bestOf, timeControl, inviteCode: privateTable ? requested || randomCode() : undefined })
      }}>
        <div className="c4-setup__heading"><span>01</span><div><small>Match configuration</small><h2>Set the drop zone</h2></div></div>
        <fieldset>
          <legend>Board format</legend>
          <div className="c4-variants">
            {variants.map((variant) => {
              const Icon = variant.icon
              return <button key={variant.id} type="button" className={mode === variant.id ? 'is-selected' : ''} onClick={() => setMode(variant.id)} aria-pressed={mode === variant.id}><Icon size={18} /><span><b>{variant.name}</b><small>{variant.detail}</small></span></button>
            })}
          </div>
        </fieldset>
        <div className="c4-setup__row">
          <label><span>Series</span><select value={bestOf} onChange={(event) => setBestOf(Number(event.target.value))}><option value={1}>One game</option><option value={3}>Best of 3</option><option value={5}>Best of 5</option><option value={7}>Best of 7</option></select></label>
          <label><span><Timer size={13} /> Drop timer</span><select value={timeControl} onChange={(event) => setTimeControl(event.target.value)}><option value="untimed">Untimed</option><option value="15s">15 seconds</option><option value="30s">30 seconds</option><option value="60s">60 seconds</option></select></label>
        </div>
        <label className="c4-field"><span>Player tag</span><input name="name" maxLength={24} defaultValue={name} autoComplete="nickname" /></label>
        <label className="c4-lock"><input type="checkbox" checked={privateTable} onChange={(event) => setPrivateTable(event.target.checked)} /><span><LockKeyhole size={16} /><b>Use a private match code</b></span></label>
        {privateTable && <label className="c4-field"><span>Match code <small>Blank creates one</small></span><input name="inviteCode" maxLength={12} placeholder="DROP4" autoCapitalize="characters" /></label>}
        <button className="c4-launch" type="submit"><span>Enter matchmaking</span><i>↓</i></button>
        <p>Authoritative rules · reconnect protection · live spectators</p>
      </form>
    </main>
  )
}

function randomCode() {
  return `DROP${Math.random().toString(36).slice(2, 6).toUpperCase()}`
}
