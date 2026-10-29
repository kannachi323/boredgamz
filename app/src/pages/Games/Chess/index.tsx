import { useEffect, useState } from 'react'
import type { MatchOptions } from '../../../hooks/useOnlineMatch'
import { useOnlineMatch } from '../../../hooks/useOnlineMatch'
import { ChessLobby } from './ChessLobby'
import { ChessTable } from './ChessTable'
import type { ChessTheme } from './chessThemes'

export function ChessGame() {
  const match = useOnlineMatch<string>('chess')
  const connectMatch = match.connect
  const playerName = match.name
  const [theme, setTheme] = useState<ChessTheme>(() => (localStorage.getItem('chess-theme') as ChessTheme) || 'walnut')
  const [pending, setPending] = useState<MatchOptions | null>(null)

  useEffect(() => {
    const roomId = new URLSearchParams(location.search).get('watch')
    if (!roomId) return
    const options: MatchOptions = { name: playerName, role: 'spectator', roomId }
    setPending(options)
    connectMatch(options)
  }, [connectMatch, playerName])

  function selectTheme(nextTheme: ChessTheme) {
    localStorage.setItem('chess-theme', nextTheme)
    setTheme(nextTheme)
  }

  function join(options: MatchOptions, nextTheme: ChessTheme) {
    selectTheme(nextTheme)
    setPending(options)
    match.connect(options)
  }

  if (match.state) return <ChessTable state={match.state} theme={theme} playerID={match.playerID} error={match.error} onTheme={selectTheme} onMove={match.chessMove} onChat={match.chat} onResign={match.resign} onRematch={match.rematch} onOfferDraw={match.offerDraw} onAcceptDraw={match.acceptDraw} onDeclineDraw={match.declineDraw} onLeave={match.close} />
  if (match.phase === 'connecting' || match.phase === 'waiting') return <main className="chs-waiting"><div className="chs-waiting__pieces"><i>♔</i><span /><i>♚</i></div><small>{pending?.role === 'spectator' ? 'OPENING OBSERVER SEAT' : 'CHESS CLUB MATCHMAKING'}</small><h1>Preparing your table.</h1><p>{pending?.inviteCode ? <>Challenge code <b>{pending.inviteCode}</b></> : 'Finding the same format and time control.'}</p>{match.error && <p className="chs-error">{match.error}</p>}<button type="button" onClick={match.close}>Cancel search</button></main>
  return <ChessLobby name={match.name} onJoin={join} />
}

export default ChessGame
