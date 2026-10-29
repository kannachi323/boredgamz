import { Copy, DoorOpen, Eye, Flag, RotateCcw } from 'lucide-react'
import { useState } from 'react'
import { MatchChat } from '../../../components/games/MatchChat'
import type { ArenaPlayer, ArenaState } from '../../../hooks/useOnlineMatch'
import { TurnClock } from '../TicTacToe/TurnClock'
import { ConnectFourBoard } from './ConnectFourBoard'

interface Props {
  state: ArenaState
  playerID: string
  error: string
  onMove: (column: number) => void
  onChat: (content: string) => void
  onResign: () => void
  onRematch: () => void
  onLeave: () => void
}

export function ConnectFourTable(props: Props) {
  const { state, playerID } = props
  const [copied, setCopied] = useState(false)
  const me = state.players.find((player) => player?.mark === state.viewerMark) ?? null
  const canMove = state.status === 'playing' && state.turn === me?.mark
  const ready = me ? state.rematchReady.includes(me.mark) : false

  async function copyLink() {
    const url = new URL(window.location.href)
    url.search = ''
    url.searchParams.set('watch', state.roomId)
    await navigator.clipboard.writeText(url.toString())
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1400)
  }

  return (
    <main className="c4-table-page">
      <header className="c4-topbar">
        <a href="/" className="c4-brand"><i>4</i><span>CONNECT / BOREDGAMZ</span></a>
        <div><b>{modeName(state.settings.mode)}</b><span>ROUND {state.round}</span><span>FIRST TO {Math.floor(state.settings.bestOf / 2) + 1}</span></div>
        <button type="button" onClick={props.onLeave}><DoorOpen size={16} /> Exit</button>
      </header>
      <div className="c4-table-layout">
        <aside className="c4-score-panel">
          <span className="c4-panel-label">Match telemetry</span>
          <h1>{headline(state, me)}</h1>
          <div className="c4-score">
            <C4Player player={state.players[0]} score={state.score.one} active={state.status === 'playing' && state.turn === 'one'} color="coral" />
            <C4Player player={state.players[1]} score={state.score.two} active={state.status === 'playing' && state.turn === 'two'} color="gold" />
          </div>
          <dl>
            <div><dt>Drop clock</dt><dd><TurnClock deadline={state.turnDeadline} control={state.settings.timeControl} /></dd></div>
            <div><dt>Observers</dt><dd><Eye size={13} /> {state.spectators}</dd></div>
            <div><dt>Room</dt><dd>{state.settings.visibility}</dd></div>
          </dl>
          {state.settings.inviteCode && <p className="c4-code">Player code <strong>{state.settings.inviteCode}</strong></p>}
          <button className="c4-copy" type="button" onClick={() => void copyLink()}><Copy size={14} /> {copied ? 'Copied' : 'Share spectator rail'}</button>
          <div className="c4-actions">
            {state.status === 'playing' && me && <button type="button" onClick={props.onResign}><Flag size={14} /> Concede game</button>}
            {state.status !== 'playing' && me && <button className="is-primary" disabled={ready} type="button" onClick={props.onRematch}><RotateCcw size={14} /> {ready ? 'Opponent ready?' : 'Ready next round'}</button>}
          </div>
          {props.error && <p className="c4-error" role="alert">{props.error}</p>}
        </aside>
        <section className="c4-board-stage">
          <div className="c4-board-stage__prompt"><span>{state.status === 'playing' ? 'DROP ACTIVE' : 'BOARD LOCKED'}</span><p>{prompt(state, me)}</p></div>
          <ConnectFourBoard state={state} canMove={canMove} onMove={props.onMove} />
          <p className="c4-rule">{rule(state.settings.mode)}</p>
        </section>
        <MatchChat messages={state.chat} playerID={playerID} onSend={props.onChat} />
      </div>
    </main>
  )
}

function C4Player({ player, score, active, color }: { player: ArenaPlayer | null; score: number; active: boolean; color: string }) {
  return <div className={`c4-player ${active ? 'is-active' : ''}`}><i className={`is-${color}`} /><span><b>{player?.name ?? 'Waiting'}</b><small>{player?.connected === false ? 'SIGNAL LOST' : active ? 'TO DROP' : 'STANDING BY'}</small></span><strong>{score.toString().padStart(2, '0')}</strong></div>
}

function modeName(mode: string) { return mode === 'quick' ? 'QUICK STACK' : mode === 'connect-five' ? 'LONG GAME' : 'CLASSIC DROP' }
function rule(mode: string) { return mode === 'connect-five' ? 'Eight columns. Five connected discs win.' : mode === 'quick' ? 'A compact board for faster tactical collisions.' : 'Drop into any open column. Four connected discs win.' }
function headline(state: ArenaState, me: ArenaPlayer | null) {
  if (state.status === 'playing') return me?.mark === state.turn ? 'Your drop.' : me ? 'Track the threat.' : 'Watching the stack.'
  if (!state.winner) return 'Gridlocked.'
  if (!me) return `${state.players.find((player) => player?.mark === state.winner)?.name ?? 'Player'} connected.`
  return state.winner === me.mark ? 'Connection made.' : 'Line conceded.'
}
function prompt(state: ArenaState, me: ArenaPlayer | null) {
  if (state.status !== 'playing') return state.seriesWinner ? 'Series complete.' : 'Both players must ready the next board.'
  if (!me) return 'Live board feed'
  return state.turn === me.mark ? 'Select any column with an open slot.' : 'Opponent is lining up a drop.'
}
