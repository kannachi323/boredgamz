import { Clock3, Grid3X3, LockKeyhole, ScanLine, Skull, Sparkles, UsersRound } from 'lucide-react'
import { useState } from 'react'
import cover from '../../../assets/games/tictactoe-cover.webp'
import type { MatchOptions } from '../../../hooks/useOnlineMatch'

const modes = [
  { id: 'classic', name: 'Classic', detail: '3 × 3 · connect three', icon: Grid3X3 },
  { id: 'four-square', name: 'Four Square', detail: '4 × 4 · connect three', icon: ScanLine },
  { id: 'misere', name: 'Misère', detail: 'Making three loses', icon: Skull },
] as const

interface Props {
  defaultName: string
  onJoin: (options: MatchOptions) => void
}

export function TicTacToeLobby({ defaultName, onJoin }: Props) {
  const [mode, setMode] = useState('classic')
  const [bestOf, setBestOf] = useState(3)
  const [timeControl, setTimeControl] = useState('30s')
  const [isPrivate, setIsPrivate] = useState(false)

  function submit(form: HTMLFormElement) {
    const data = new FormData(form)
    const requestedCode = String(data.get('inviteCode') ?? '').trim().toUpperCase()
    const inviteCode = isPrivate ? requestedCode || makeInviteCode() : undefined
    onJoin({
      name: String(data.get('name') ?? defaultName),
      mode,
      bestOf,
      timeControl,
      inviteCode,
    })
  }

  return (
    <main className="ttt-lobby">
      <section className="ttt-lobby__intro">
        <a className="ttt-back" href="/games">← All games</a>
        <div className="ttt-kicker"><Sparkles size={14} /> BoredGamz quick table</div>
        <h1>Tic-Tac-Toe,<br /><em>properly played.</em></h1>
        <p>Fast rounds, uncommon variants, live opponents, and just enough table talk. No account required.</p>
        <div className="ttt-cover" aria-hidden="true">
          <img src={cover} alt="" />
          <span><UsersRound size={16} /> Live two-player tables</span>
        </div>
      </section>

      <form className="ttt-lobby__setup" onSubmit={(event) => { event.preventDefault(); submit(event.currentTarget) }}>
        <header>
          <span>Table setup</span>
          <strong>Choose your match</strong>
        </header>

        <fieldset>
          <legend>Ruleset</legend>
          <div className="ttt-mode-grid">
            {modes.map((item) => {
              const Icon = item.icon
              return (
                <button className={mode === item.id ? 'is-selected' : ''} key={item.id} type="button" onClick={() => setMode(item.id)} aria-pressed={mode === item.id}>
                  <Icon size={19} />
                  <span><b>{item.name}</b><small>{item.detail}</small></span>
                </button>
              )
            })}
          </div>
        </fieldset>

        <div className="ttt-option-row">
          <label>
            <span>Match length</span>
            <select value={bestOf} onChange={(event) => setBestOf(Number(event.target.value))}>
              <option value={1}>Single round</option>
              <option value={3}>Best of 3</option>
              <option value={5}>Best of 5</option>
              <option value={7}>Best of 7</option>
            </select>
          </label>
          <label>
            <span><Clock3 size={14} /> Move clock</span>
            <select value={timeControl} onChange={(event) => setTimeControl(event.target.value)}>
              <option value="untimed">Untimed</option>
              <option value="15s">15 seconds</option>
              <option value="30s">30 seconds</option>
              <option value="60s">60 seconds</option>
            </select>
          </label>
        </div>

        <label className="ttt-name-field">
          <span>Display name</span>
          <input name="name" defaultValue={defaultName} maxLength={24} autoComplete="nickname" />
        </label>

        <label className="ttt-private-toggle">
          <input type="checkbox" checked={isPrivate} onChange={(event) => setIsPrivate(event.target.checked)} />
          <span><LockKeyhole size={17} /><b>Private table</b><small>Only players with the same code are matched.</small></span>
        </label>

        {isPrivate && (
          <label className="ttt-name-field ttt-invite-field">
            <span>Invite code <small>Leave blank to generate one</small></span>
            <input name="inviteCode" maxLength={12} placeholder="FRIDAY" autoCapitalize="characters" />
          </label>
        )}

        <button className="ttt-primary-action" type="submit">Find my opponent <span>→</span></button>
        <p className="ttt-fine-print">By joining, you agree to keep chat friendly. Games recover automatically after short connection drops.</p>
      </form>
    </main>
  )
}

function makeInviteCode() {
  return `PLAY${Math.random().toString(36).slice(2, 6).toUpperCase()}`
}
