import { Radio, Users, WifiOff } from 'lucide-react'
import { useMemo, useState } from 'react'
import { getGame } from '../../data/games'
import { useOnlineMatch, type ArenaPlayer } from '../../hooks/useOnlineMatch'
import { gridConfigs, gridWinner, placeGridMove, type GridKind, type GridMark } from '../../lib/gridGames'
import { GameFrame } from './GameFrame'

const labels = { one: 'Player one', two: 'Player two' }

export function GridGame({ kind }: { kind: GridKind }) {
  const config = gridConfigs[kind]
  const game = getGame(kind)!
  const [mode, setMode] = useState<'local' | 'online'>('local')
  const [board, setBoard] = useState<GridMark[]>(() => Array(config.rows * config.columns).fill(''))
  const [turn, setTurn] = useState<Exclude<GridMark, ''>>('one')
  const [winner, setWinner] = useState<GridMark>('')
  const [draw, setDraw] = useState(false)
  const online = useOnlineMatch(kind)
  const playerID = online.playerID
  const onlinePlayer = online.state?.players.find((player) => player?.id === playerID)
  const visibleBoard = mode === 'online' && online.state ? online.state.board : board
  const visibleTurn = mode === 'online' && online.state ? online.state.turn : turn
  const visibleWinner = mode === 'online' && online.state ? online.state.winner ?? '' : winner

  const status = useMemo(() => {
    if (mode === 'online') {
      if (online.phase === 'idle') return 'Online table ready'
      if (online.phase === 'connecting') return 'Connecting…'
      if (online.phase === 'waiting') return 'Finding an opponent…'
      if (online.phase === 'finished') return visibleWinner ? `${labels[visibleWinner]} wins` : 'Draw game'
      return visibleTurn === onlinePlayer?.mark ? 'Your turn' : "Opponent's turn"
    }
    if (winner) return `${labels[winner]} wins`
    if (draw) return 'Draw game'
    return `${labels[turn]}'s turn`
  }, [draw, mode, online.phase, onlinePlayer?.mark, turn, visibleTurn, visibleWinner, winner])

  function reset() {
    setBoard(Array(config.rows * config.columns).fill(''))
    setTurn('one')
    setWinner('')
    setDraw(false)
    online.close()
  }

  function select(position: number) {
    if (mode === 'online') {
      if (!online.state || online.state.status !== 'playing' || visibleTurn !== onlinePlayer?.mark) return
      online.move(kind === 'connect-four' ? position % config.columns : position)
      return
    }
    if (winner || draw) return
    const result = placeGridMove(board, config, kind === 'connect-four' ? position % config.columns : position, turn)
    if (!result) return
    setBoard(result.board)
    const nextWinner = gridWinner(result.board, config, result.position)
    if (nextWinner) setWinner(nextWinner)
    else if (result.board.every(Boolean)) setDraw(true)
    else setTurn(turn === 'one' ? 'two' : 'one')
  }

  const aside = (
    <>
      <div className="mode-tabs" aria-label="Play mode">
        <button className={mode === 'local' ? 'is-active' : ''} onClick={() => { online.close(); setMode('local') }}><Users size={16} /> Local</button>
        <button className={mode === 'online' ? 'is-active' : ''} onClick={() => setMode('online')}><Radio size={16} /> Online</button>
      </div>
      {mode === 'local' ? (
        <div className="player-list">
          <PlayerRow mark="one" name="Player one" active={!winner && !draw && turn === 'one'} />
          <PlayerRow mark="two" name="Player two" active={!winner && !draw && turn === 'two'} />
        </div>
      ) : online.phase === 'idle' ? (
        <form className="join-form" onSubmit={(event) => { event.preventDefault(); const data = new FormData(event.currentTarget); online.connect(String(data.get('name') ?? 'Guest')) }}>
          <label htmlFor={`${kind}-name`}>Your table name</label>
          <input id={`${kind}-name`} name="name" maxLength={24} defaultValue={online.name} autoComplete="nickname" />
          <button className="button button--primary" type="submit">Find opponent</button>
          <p>We will match you with the next available player. Keep this tab open while you wait.</p>
        </form>
      ) : (
        <div className="online-panel">
          <div className="connection-state"><span className={online.phase === 'playing' ? 'is-live' : ''} /> {online.phase === 'waiting' ? 'Searching…' : online.phase}</div>
          {online.state?.players.filter((player): player is ArenaPlayer => Boolean(player)).map((player) => <PlayerRow key={player.id} mark={player.mark} name={player.id === playerID ? `${player.name} (you)` : player.name} active={online.state?.turn === player.mark} disconnected={!player.connected} />)}
          <button className="button button--ghost" type="button" onClick={online.close}><WifiOff size={16} /> Leave table</button>
        </div>
      )}
      {online.error && <p className="game-error" role="alert">{online.error}</p>}
      <div className="rules-card"><p className="eyebrow">How to play</p><p>{kind === 'connect-four' ? 'Choose a column. Your disc falls to the lowest open space.' : `Place your piece on an open space. First to connect ${config.connect} wins.`}</p></div>
    </>
  )

  return (
    <GameFrame game={game} status={status} aside={aside} onReset={reset}>
      {(mode === 'online' && (online.phase === 'idle' || online.phase === 'connecting' || online.phase === 'waiting')) ? (
        <div className="table-waiting"><span className="waiting-orbit" /><p>{online.phase === 'idle' ? 'Choose Online and find an opponent to open the board.' : 'Setting the table…'}</p></div>
      ) : (
        <div className={`grid-board grid-board--${kind}`} style={{ gridTemplateColumns: `repeat(${config.columns}, 1fr)` }} role="grid" aria-label={`${game.name} board`}>
          {visibleBoard.map((mark, index) => (
            <button key={index} className={`board-button ${mark ? `is-${mark}` : ''} ${online.state?.lastMove === index ? 'is-last' : ''}`} type="button" role="gridcell" aria-label={`${Math.floor(index / config.columns) + 1}, ${index % config.columns + 1}${mark ? `: ${labels[mark]}` : ': empty'}`} onClick={() => select(index)} disabled={Boolean(mark) && kind !== 'connect-four'}>
              {mark && <span />}
            </button>
          ))}
        </div>
      )}
    </GameFrame>
  )
}

function PlayerRow({ mark, name, active, disconnected = false }: { mark: Exclude<GridMark, ''>; name: string; active: boolean; disconnected?: boolean }) {
  return <div className={`player-row ${active ? 'is-active' : ''}`}><i className={`piece piece--${mark}`} /><span>{name}</span>{disconnected ? <small>reconnecting</small> : active ? <small>to move</small> : null}</div>
}
