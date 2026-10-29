import { useEffect, useState } from "react";
import type { MahjongOptions } from "../../../hooks/useMahjongMatch";
import { useMahjongMatch } from "../../../hooks/useMahjongMatch";
import { MahjongLobby } from "./MahjongLobby";
import { MahjongTable } from "./MahjongTable";

export function MahjongGame() {
  const match = useMahjongMatch();
  const connectMatch = match.connect;
  const playerName = match.name;
  const [pending, setPending] = useState<MahjongOptions | null>(null);
  useEffect(() => {
    const roomId = new URLSearchParams(location.search).get("watch");
    if (!roomId) return;
    const options: MahjongOptions = { name: playerName, role: "spectator", roomId };
    setPending(options);
    connectMatch(options);
  }, [connectMatch, playerName]);
  function join(options: MahjongOptions) {
    setPending(options);
    match.connect(options);
  }
  if (match.state) {
    return <MahjongTable state={match.state} playerID={match.playerID} error={match.error} onPair={match.pair} onHint={match.hint} onChat={match.chat} onRematch={match.rematch} onLeave={match.close} />;
  }
  if (match.phase === "connecting" || match.phase === "waiting") {
    return (
      <main className="mjh-waiting">
        <div className="mjh-waiting__tiles"><i>東</i><i>中</i><i>發</i></div>
        <small>{pending?.role === "spectator" ? "OPENING THE GARDEN" : "FINDING YOUR TABLE"}</small>
        <h1>Setting the tiles.</h1>
        <p>{pending?.inviteCode ? <>Private garden <b>{pending.inviteCode}</b></> : `Finding another ${modeName(pending?.mode)} player.`}</p>
        {match.error && <p className="mjh-error">{match.error}</p>}
        <button type="button" onClick={match.close}>Leave queue</button>
      </main>
    );
  }
  return <MahjongLobby name={match.name} onJoin={join} />;
}

function modeName(mode?: MahjongOptions["mode"]) {
  return mode === "duel" ? "Mirror Duel" : mode === "zen" ? "Zen Partners" : "Shared Garden";
}

export default MahjongGame;
