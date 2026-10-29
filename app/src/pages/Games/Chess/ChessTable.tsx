import { Copy, DoorOpen, Eye, Flag, Handshake, Palette, RotateCcw, Settings2, X } from 'lucide-react'
import { useState } from 'react'
import { MatchChat } from '../../../components/games/MatchChat'
import type { ArenaPlayer, ArenaState } from '../../../hooks/useOnlineMatch'
import { ChessBoard } from './ChessBoard'
import { ChessClock } from './ChessClock'
import { chessThemes, type ChessTheme } from './chessThemes'

interface Props {
  state: ArenaState<string>
  theme: ChessTheme
  playerID: string
  error: string
  onTheme: (theme: ChessTheme) => void
  onMove: (from: string, to: string, promotion?: string) => void
  onChat: (text: string) => void
  onResign: () => void
  onRematch: () => void
  onOfferDraw: () => void
  onAcceptDraw: () => void
  onDeclineDraw: () => void
  onLeave: () => void
}

export function ChessTable(props: Props) {
  const { state } = props
  const [settingsOpen, setSettingsOpen] = useState(false)
  const me = state.players.find((player) => player?.mark === state.viewerMark) ?? null
  const ready = me ? state.rematchReady.includes(me.mark) : false
  const white = state.players[0]
  const black = state.players[1]

  async function share() {
    const url = new URL(window.location.href)
    url.search = ''
    url.searchParams.set('watch', state.roomId)
    await navigator.clipboard.writeText(url.toString())
  }

  return <main className="chs-table-page">
    <header className="chs-topbar">
      <a href="/">Bored<span>Gamz</span> <i /> CHESS CLUB</a>
      <div><span>{formatLabel(state.settings.mode)}</span><i /> {formatClock(state.settings.timeControl)}<i /> ROUND {state.round}</div>
      <div className="chs-topbar__actions"><button type="button" onClick={() => setSettingsOpen(true)}><Settings2 size={15} /> Settings</button><button type="button" onClick={props.onLeave}><DoorOpen size={15} /> Leave</button></div>
    </header>
    <div className="chs-room">
      <aside className="chs-ledger">
        <div className="chs-room-id"><span>Club table</span><b>#{state.roomId.slice(-6).toUpperCase()}</b></div>
        <PlayerCard player={black} mark="two" active={state.status === 'playing' && state.turn === 'two'} milliseconds={state.clocks?.two} score={state.score.two} />
        <section className="chs-moves">
          <header><span>Move ledger</span><small>{state.chess?.moves.length ?? 0} plies</small></header>
          <ol aria-label="Move history">{pairMoves(state.chess?.moves ?? []).map((pair, index) => <li key={index}><b>{index + 1}.</b><span>{pair[0]?.san ?? ''}</span><span>{pair[1]?.san ?? ''}</span></li>)}</ol>
          {!state.chess?.moves.length && <p>The ledger is waiting for White.</p>}
        </section>
        <PlayerCard player={white} mark="one" active={state.status === 'playing' && state.turn === 'one'} milliseconds={state.clocks?.one} score={state.score.one} />
        <div className="chs-room-tools"><span><Eye size={14} /> {state.spectators} watching</span><button type="button" onClick={() => void share()}><Copy size={14} /> Spectator link</button></div>
      </aside>
      <section className="chs-stage">
        <div className="chs-stage__heading"><div><small>{state.status === 'playing' ? 'LIVE POSITION' : 'GAME COMPLETE'}</small><h1>{headline(state, me)}</h1></div><p>{positionNote(state, me)}</p></div>
        <ChessBoard state={state} theme={props.theme} myMark={me?.mark} onMove={props.onMove} />
        <footer className="chs-stage__footer"><span>{state.chess?.fen.split(' ')[1] === 'w' ? 'White' : 'Black'} to move</span><i /> <span>{state.chess?.method && state.chess.method !== 'NoMethod' ? state.chess.method : 'Server-verified position'}</span></footer>
        <div className="chs-game-actions">
          {state.status === 'playing' && me && <button type="button" onClick={props.onResign}><Flag size={15} /> Resign</button>}
          {state.status === 'playing' && me && !state.drawOfferedBy && state.turn !== me.mark && Boolean(state.chess?.moves.length) && <button type="button" onClick={props.onOfferDraw}><Handshake size={15} /> Offer draw</button>}
          {state.status === 'playing' && me && state.drawOfferedBy === me.mark && <span className="chs-draw-pending"><Handshake size={14} /> Draw offered</span>}
          {state.status !== 'playing' && me && <button className="is-primary" type="button" disabled={ready} onClick={props.onRematch}><RotateCcw size={15} /> {ready ? 'Opponent must ready' : 'Play again'}</button>}
          {state.settings.inviteCode && <span>Challenge code <b>{state.settings.inviteCode}</b></span>}
        </div>
        {state.status === 'playing' && me && state.drawOfferedBy && state.drawOfferedBy !== me.mark && <aside className="chs-draw-offer" role="status"><span><Handshake size={17} /><b>Your opponent offers a draw.</b></span><div><button type="button" onClick={props.onDeclineDraw}>Decline</button><button className="is-primary" type="button" onClick={props.onAcceptDraw}>Accept draw</button></div></aside>}
        {props.error && <p className="chs-error" role="alert">{props.error}</p>}
      </section>
      <MatchChat messages={state.chat} playerID={props.playerID} onSend={props.onChat} />
    </div>
    {settingsOpen && <div className="chs-settings-backdrop" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) setSettingsOpen(false) }}><section className="chs-settings" role="dialog" aria-modal="true" aria-labelledby="chess-settings-title"><header><div><Palette size={17} /><span><small>PERSONAL DISPLAY</small><h2 id="chess-settings-title">Board settings</h2></span></div><button type="button" aria-label="Close settings" onClick={() => setSettingsOpen(false)}><X size={18} /></button></header><p>Color choices apply only to your board and can be changed during play.</p><div className="chs-settings__themes">{(Object.keys(chessThemes) as ChessTheme[]).map((key) => <button key={key} type="button" className={props.theme === key ? 'is-selected' : ''} aria-pressed={props.theme === key} onClick={() => props.onTheme(key)}><i style={{ background: `linear-gradient(135deg,${chessThemes[key].light} 50%,${chessThemes[key].dark} 50%)` }} /><span><b>{chessThemes[key].name}</b><small>{key === 'midnight' ? 'Low-light contrast' : key === 'sage' ? 'Quiet tournament green' : key === 'parchment' ? 'Warm reading room' : 'Classic club wood'}</small></span></button>)}</div><footer><span>Piece orientation follows your assigned color.</span><button type="button" onClick={() => setSettingsOpen(false)}>Done</button></footer></section></div>}
  </main>
}

function PlayerCard({ player, mark, active, milliseconds, score }: { player: ArenaPlayer | null; mark: 'one' | 'two'; active: boolean; milliseconds?: number; score: number }) {
  return <article className={`chs-player ${active ? 'is-active' : ''}`}><i className={`chs-player__piece is-${mark}`} aria-hidden="true">{mark === 'one' ? '♔' : '♚'}</i><span><b>{player?.name ?? 'Waiting for player'}</b><small>{player?.connected === false ? 'Reconnecting…' : mark === 'one' ? 'White' : 'Black'} · {score} pts</small></span>{milliseconds === undefined ? <time>∞</time> : <ChessClock milliseconds={milliseconds} active={active} />}</article>
}

function pairMoves<T>(moves: T[]) { const rows: [T?, T?][] = []; for (let index = 0; index < moves.length; index += 2) rows.push([moves[index], moves[index + 1]]); return rows }
function formatLabel(mode: string) { return mode === 'blitz' ? 'Blitz' : mode === 'rapid' ? 'Rapid' : 'Classical' }
function formatClock(value: string) { return value.replace(/m/g, ' min').replace(/s/g, ' sec') }
function headline(state: ArenaState<string>, me: ArenaPlayer | null) {
  if (state.status === 'playing') return me ? state.turn === me.mark ? 'Your move.' : 'Opponent is thinking.' : 'Observing the table.'
  if (!state.winner) return 'The game is drawn.'
  return !me ? `${state.winner === 'one' ? 'White' : 'Black'} wins.` : state.winner === me.mark ? 'Checkmate. Well played.' : 'The game is decided.'
}
function positionNote(state: ArenaState<string>, me: ArenaPlayer | null) {
  if (state.status !== 'playing') return state.chess?.method || 'Both players may request another game.'
  if (!me) return 'Spectator mode · board follows White’s perspective'
  return state.turn === me.mark ? 'Select one of your pieces to reveal legal destinations.' : 'The clock and position stay live while you wait.'
}
