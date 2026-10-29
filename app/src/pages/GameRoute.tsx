import { useParams } from "react-router-dom";
import { GameLobby } from "../components/room/GameLobby";
import { useRoomFlow } from "../components/room/useRoomFlow";
import { Waiting } from "../components/room/Waiting";
import { getGame, type GameDefinition } from "../games/catalog";
import { ChessTable } from "../games/tables/ChessTable";
import { GridTable } from "../games/tables/GridTable";
import { MahjongTable } from "../games/tables/MahjongTable";
import { PokerTable } from "../games/tables/PokerTable";
import type { SocketPhase } from "../hooks/useGameSocket";
import { useMahjongMatch } from "../hooks/useMahjongMatch";
import { useOnlineMatch, type GridKind } from "../hooks/useOnlineMatch";
import { usePokerMatch } from "../hooks/usePokerMatch";
import { ErrorPage } from "./ErrorPage";

export function GameRoute() {
  const { slug } = useParams();
  const game = getGame(slug);
  if (!game) return <ErrorPage />;
  // key resets all match state when moving between games.
  if (game.slug === "poker") return <PokerRoom key={game.slug} game={game} />;
  if (game.slug === "mahjong") return <MahjongRoom key={game.slug} game={game} />;
  if (game.slug === "chess") return <ChessRoom key={game.slug} game={game} />;
  return <GridRoom key={game.slug} game={game} kind={game.slug as GridKind} />;
}

interface Pending {
  game: GameDefinition;
  phase: SocketPhase;
  error: string;
  flow: ReturnType<typeof useRoomFlow>;
}

/** Lobby or waiting screen until the server sends the first table state. */
function BeforeTable({ game, phase, error, flow }: Pending) {
  if (phase === "connecting" || phase === "waiting") {
    return <Waiting game={game} request={flow.request} connecting={phase === "connecting"} error={error} onCancel={flow.leave} />;
  }
  return <GameLobby game={game} error={error} onJoin={flow.join} />;
}

function GridRoom({ game, kind }: { game: GameDefinition; kind: GridKind }) {
  const match = useOnlineMatch(kind);
  const flow = useRoomFlow(match.connect, match.close, match.resume);
  if (!match.state) return <BeforeTable game={game} phase={match.phase} error={match.error} flow={flow} />;
  return <GridTable game={game} state={match.state} error={match.error} reconnecting={match.reconnecting} onMove={match.move} onChat={match.chat} onResign={match.resign} onRematch={match.rematch} onLeave={flow.leave} />;
}

function ChessRoom({ game }: { game: GameDefinition }) {
  const match = useOnlineMatch<string>("chess");
  const flow = useRoomFlow(match.connect, match.close, match.resume);
  if (!match.state) return <BeforeTable game={game} phase={match.phase} error={match.error} flow={flow} />;
  return (
    <ChessTable
      game={game}
      state={match.state}
      error={match.error}
      reconnecting={match.reconnecting}
      onMove={match.chessMove}
      onChat={match.chat}
      onResign={match.resign}
      onRematch={match.rematch}
      onOfferDraw={match.offerDraw}
      onAcceptDraw={match.acceptDraw}
      onDeclineDraw={match.declineDraw}
      onLeave={flow.leave}
    />
  );
}

function PokerRoom({ game }: { game: GameDefinition }) {
  const match = usePokerMatch();
  const flow = useRoomFlow(match.connect, match.close, match.resume);
  if (!match.state) return <BeforeTable game={game} phase={match.phase} error={match.error} flow={flow} />;
  return <PokerTable game={game} state={match.state} error={match.error} reconnecting={match.reconnecting} onAction={match.action} onDiscard={match.discard} onChat={match.chat} onNext={match.nextHand} onLeave={flow.leave} />;
}

function MahjongRoom({ game }: { game: GameDefinition }) {
  const match = useMahjongMatch();
  const flow = useRoomFlow(match.connect, match.close, match.resume);
  if (!match.state) return <BeforeTable game={game} phase={match.phase} error={match.error} flow={flow} />;
  return <MahjongTable game={game} state={match.state} error={match.error} reconnecting={match.reconnecting} onPair={match.pair} onHint={match.hint} onChat={match.chat} onRematch={match.rematch} onLeave={flow.leave} />;
}
