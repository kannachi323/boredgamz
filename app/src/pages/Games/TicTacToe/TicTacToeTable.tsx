import { Copy, DoorOpen, Eye, Flag, RotateCcw, Settings2 } from 'lucide-react'
import { useState } from 'react'
import { MatchChat } from '../../../components/games/MatchChat'
import type { ArenaPlayer, ArenaState } from '../../../hooks/useOnlineMatch'
import { TicTacToeBoard } from './TicTacToeBoard'
import { TurnClock } from './TurnClock'

interface Props {
  state: ArenaState
  playerID: string
  error: string
  onMove: (position: number) => void
  onChat: (content: string) => void
  onResign: () => void
  onRematch: () => void
  onLeave: () => void
}

export function TicTacToeTable({ state, playerID, error, onMove, onChat, onResign, onRematch, onLeave }: Props) {
  const [copied, setCopied] = useState(false)
  const currentPlayer = state.players.find((player) => player?.mark === state.viewerMark) ?? null
  const isSpectator = !currentPlayer
  const canMove = state.status === 'playing' && currentPlayer?.mark === state.turn
  const headline = tableHeadline(state, currentPlayer)
  const rematchSent = currentPlayer ? state.rematchReady.includes(currentPlayer.mark) : false

  async function copySpectatorLink() {
    const url = new URL(window.location.href)
    url.search = ''
    url.searchParams.set('watch', state.roomId)
    await navigator.clipboard.writeText(url.toString())
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1600)
  }

  return (
    <main className="ttt-table-page">
      <header className="ttt-table-bar">
        <a href="/" className="ttt-table-brand"><i>BG</i><span>Bored<span>Gamz</span></span></a>
        <div className="ttt-table-bar__meta"><span>{modeLabel(state.settings.mode)}</span><i /> <span>Best of {state.settings.bestOf}</span><i /> <span>Round {state.round}</span></div>
        <button type="button" onClick={onLeave}><DoorOpen size={16} /> Leave table</button>
      </header>

      <div className="ttt-table-layout">
        <aside className="ttt-match-panel">
          <div className="ttt-match-panel__label">Live match</div>
          <h1>{headline}</h1>
          <div className="ttt-scoreboard">
            <PlayerCard player={state.players[0]} active={state.turn === 'one' && state.status === 'playing'} score={state.score.one} symbol="×" />
            <div className="ttt-scoreboard__versus">vs</div>
            <PlayerCard player={state.players[1]} active={state.turn === 'two' && state.status === 'playing'} score={state.score.two} symbol="○" />
          </div>
          <div className="ttt-match-detail">
            <div><span>Move clock</span><TurnClock deadline={state.turnDeadline} control={state.settings.timeControl} /></div>
            <div><span>Spectators</span><b><Eye size={14} /> {state.spectators}</b></div>
            <div><span>Table</span><b>{state.settings.visibility}</b></div>
          </div>
          <button className="ttt-copy-button" type="button" onClick={() => void copySpectatorLink()}><Copy size={15} /> {copied ? 'Link copied' : 'Invite spectators'}</button>
          {state.settings.inviteCode && <p className="ttt-room-code">Player code <strong>{state.settings.inviteCode}</strong></p>}
          <div className="ttt-match-actions">
            {state.status === 'playing' && !isSpectator && <button type="button" onClick={onResign}><Flag size={15} /> Resign</button>}
            {state.status !== 'playing' && !isSpectator && <button className="is-primary" type="button" disabled={rematchSent} onClick={onRematch}><RotateCcw size={15} /> {rematchSent ? 'Waiting for rival' : 'Play again'}</button>}
            <button type="button" disabled title="Board themes arrive with player profiles"><Settings2 size={15} /> Table settings</button>
          </div>
        </aside>

        <section className="ttt-board-stage">
          <div className="ttt-board-stage__status">
            <span className={state.status === 'playing' ? 'is-live' : ''}>{state.status === 'playing' ? 'Live' : 'Round complete'}</span>
            <p>{boardPrompt(state, currentPlayer)}</p>
          </div>
          <TicTacToeBoard state={state} canMove={canMove} onMove={onMove} />
          {error && <p className="ttt-table-error" role="alert">{error}</p>}
          <p className="ttt-rules-note">{ruleText(state.settings.mode)}</p>
        </section>

        <MatchChat messages={state.chat} playerID={playerID} onSend={onChat} />
      </div>
    </main>
  )
}

function PlayerCard({ player, active, score, symbol }: { player: ArenaPlayer | null; active: boolean; score: number; symbol: string }) {
  return (
    <div className={`ttt-player-card ${active ? 'is-active' : ''}`}>
      <div className="ttt-player-card__symbol">{symbol}</div>
      <div><strong>{player?.name ?? 'Waiting…'}</strong><span>{player?.connected === false ? 'Reconnecting' : active ? 'Thinking' : 'Ready'}</span></div>
      <b>{score}</b>
    </div>
  )
}

function tableHeadline(state: ArenaState, player: ArenaPlayer | null) {
  if (state.status === 'playing') return player?.mark === state.turn ? 'Your move.' : player ? 'Their move.' : 'Watching live.'
  if (!state.winner) return 'A fine draw.'
  if (!player) return `${state.players.find((item) => item?.mark === state.winner)?.name ?? 'A player'} wins.`
  return state.winner === player.mark ? 'Round won.' : 'Round lost.'
}

function boardPrompt(state: ArenaState, player: ArenaPlayer | null) {
  if (state.status !== 'playing') return state.seriesWinner ? 'The series is decided.' : 'Both players can ready up for the next round.'
  if (!player) return `${state.players.find((item) => item?.mark === state.turn)?.name ?? 'A player'} is choosing a square.`
  return player.mark === state.turn ? 'Find the fork. The table is yours.' : 'Your rival is choosing a square.'
}

function modeLabel(mode: string) {
  if (mode === 'four-square') return 'Four Square'
  if (mode === 'misere') return 'Misère'
  return 'Classic'
}

function ruleText(mode: string) {
  if (mode === 'misere') return 'Misère rules: complete a line of three and your opponent wins the round.'
  if (mode === 'four-square') return 'Four Square rules: the board is larger, but a line of three still takes the round.'
  return 'Classic rules: make a line of three horizontally, vertically, or diagonally.'
}
