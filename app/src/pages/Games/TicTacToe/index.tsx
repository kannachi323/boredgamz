import { useEffect, useState } from 'react'
import type { MatchOptions } from '../../../hooks/useOnlineMatch'
import { useOnlineMatch } from '../../../hooks/useOnlineMatch'
import { TicTacToeLobby } from './TicTacToeLobby'
import { TicTacToeTable } from './TicTacToeTable'

export function TicTacToeGame() {
  const match = useOnlineMatch('tictactoe')
  const connectMatch = match.connect
  const playerName = match.name
  const [pendingOptions, setPendingOptions] = useState<MatchOptions | null>(null)

  useEffect(() => {
    const roomId = new URLSearchParams(window.location.search).get('watch')
    if (!roomId) return
    const options: MatchOptions = { name: playerName, role: 'spectator', roomId }
    setPendingOptions(options)
    connectMatch(options)
  }, [connectMatch, playerName])

  function join(options: MatchOptions) {
    setPendingOptions(options)
    match.connect(options)
  }

  if (match.state) {
    return <TicTacToeTable state={match.state} playerID={match.playerID} error={match.error} onMove={match.move} onChat={match.chat} onResign={match.resign} onRematch={match.rematch} onLeave={match.close} />
  }

  if (match.phase === 'connecting' || match.phase === 'waiting') {
    return (
      <main className="ttt-waiting-page">
        <div className="ttt-waiting-piece"><i /><i /><i /></div>
        <p>{pendingOptions?.role === 'spectator' ? 'Opening the spectator rail' : 'Looking for a worthy rival'}</p>
        <h1>{match.phase === 'connecting' ? 'Pulling up a chair…' : 'Your table is ready.'}</h1>
        {pendingOptions?.inviteCode && <div className="ttt-waiting-code"><span>Invite code</span><strong>{pendingOptions.inviteCode}</strong></div>}
        <small>{match.error || 'You can leave this tab in the foreground while matchmaking finishes.'}</small>
        <button type="button" onClick={match.close}>Cancel search</button>
      </main>
    )
  }

  return <TicTacToeLobby defaultName={match.name} onJoin={join} />
}

export default TicTacToeGame
