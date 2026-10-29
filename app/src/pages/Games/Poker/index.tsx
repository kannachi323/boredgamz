import { useEffect, useState } from "react";
import type { PokerOptions } from "../../../hooks/usePokerMatch";
import { usePokerMatch } from "../../../hooks/usePokerMatch";
import { PokerLobby } from "./PokerLobby";
import { PokerTable } from "./PokerTable";

export function PokerGame() {
  const match = usePokerMatch();
  const connectMatch = match.connect;
  const playerName = match.name;
  const [pending, setPending] = useState<PokerOptions | null>(null);
  useEffect(() => {
    const roomId = new URLSearchParams(location.search).get("watch");
    if (!roomId) return;
    const options: PokerOptions = {
      name: playerName,
      role: "spectator",
      roomId,
    };
    setPending(options);
    connectMatch(options);
  }, [connectMatch, playerName]);
  function join(options: PokerOptions) {
    setPending(options);
    match.connect(options);
  }
  if (match.state)
    return (
      <PokerTable
        state={match.state}
        playerID={match.playerID}
        error={match.error}
        onAction={match.action}
        onDiscard={match.discard}
        onChat={match.chat}
        onNext={match.nextHand}
        onLeave={match.close}
      />
    );
  if (match.phase === "connecting" || match.phase === "waiting")
    return (
      <main className="pkr-waiting">
        <div>
          <i>♠</i>
          <i>♥</i>
          <i>♦</i>
          <i>♣</i>
        </div>
        <small>
          {pending?.role === "spectator"
            ? "OPENING THE RAIL"
            : "CARD ROOM MATCHMAKING"}
        </small>
        <h1>Holding your seat.</h1>
        <p>
          {pending?.inviteCode ? (
            <>
              Private code <b>{pending.inviteCode}</b>
            </>
          ) : (
            `Finding another ${pending?.mode === "five-card-draw" ? "draw" : "Hold’em"} player.`
          )}
        </p>
        {match.error && <p className="pkr-error">{match.error}</p>}
        <button type="button" onClick={match.close}>
          Leave queue
        </button>
      </main>
    );
  return <PokerLobby name={match.name} onJoin={join} />;
}

export default PokerGame;
