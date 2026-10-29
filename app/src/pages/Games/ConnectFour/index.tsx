import { useEffect, useState } from 'react'
import type { MatchOptions } from '../../../hooks/useOnlineMatch'
import { useOnlineMatch } from '../../../hooks/useOnlineMatch'
import { ConnectFourLobby } from './ConnectFourLobby'
import { ConnectFourTable } from './ConnectFourTable'

export function ConnectFourGame() {
  const match = useOnlineMatch('connect-four')
  const connectMatch = match.connect
  const playerName = match.name
  const [pending, setPending] = useState<MatchOptions | null>(null)

  useEffect(() => {
    const roomId = new URLSearchParams(window.location.search).get('watch')
    if (!roomId) return
    const options: MatchOptions = { name: playerName, role: 'spectator', roomId }
    setPending(options)
    connectMatch(options)
  }, [connectMatch, playerName])

  function join(options: MatchOptions) { setPending(options); match.connect(options) }

  if (match.state) return <ConnectFourTable state={match.state} playerID={match.playerID} error={match.error} onMove={match.move} onChat={match.chat} onResign={match.resign} onRematch={match.rematch} onLeave={match.close} />
  if (match.phase === 'connecting' || match.phase === 'waiting') return <main className="c4-waiting"><div><i /><i /><i /><i /></div><span>{pending?.role === 'spectator' ? 'CONNECTING TO RAIL' : 'SCANNING LIVE TABLES'}</span><h1>{match.phase === 'waiting' ? 'Opponent search active.' : 'Opening the grid.'}</h1>{pending?.inviteCode && <p>PLAYER CODE <strong>{pending.inviteCode}</strong></p>}<small>{match.error || 'Matchmaking keeps modes, series length, and clock preferences exact.'}</small><button type="button" onClick={match.close}>Cancel</button></main>
  return <ConnectFourLobby name={match.name} onJoin={join} />
}

export default ConnectFourGame
