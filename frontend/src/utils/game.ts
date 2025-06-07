import { sendData } from "./connection";
import { GameState, Move } from "../types";

export function makeMove(socket: WebSocket | null, gameState: GameState, move: Move) {
  if (!socket || !move) return;
  gameState.lastMove = move;

  sendData(socket, {
    type: "move",
    data: gameState,
  });
}
