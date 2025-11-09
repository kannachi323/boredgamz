import { useGameStore } from "../../stores/useGameStore"
import { ClientRequest,   } from "../../types"

export function sendMove(row: number, col: number) {
  const { sendClientRequest, gameState, player, conn} = useGameStore.getState();
  if (!player || !gameState || !conn) return;

  const clientReq : ClientRequest = {
    type: "move",
    data: {
      ...gameState,
      lastMove: {
        r: row,
        c: col,
        color: player.color
      }
    }
  }
  
  sendClientRequest(conn, clientReq);
}