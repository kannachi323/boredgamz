import { create } from 'zustand'
import { ServerResponse, GameState, Player, ClientRequest, LobbyRequest } from '../types'
import { convertTime } from '../utils.ts'

interface GameStore {
  gameState: GameState | null
  conn: WebSocket | null
  player: Player
  opponent: Player
  
  setPlayer: (player: Player) => void
  setOpponent: (opponent: Player) => void
  setConnection: (lobbyType: string, player: Player, onMessage : (data: ServerResponse) => void) => WebSocket | null
  handler: (payload: ServerResponse) => void
  sendClientRequest: (socket: WebSocket, data: ClientRequest) => void
  sendLobbyRequest: (socket: WebSocket, data: LobbyRequest) => void

}

export const useGameStore = create<GameStore>((set) => ({
  gameState: null,
  conn: null,
  player: { playerID: '', playerName: '', color: 'black', playerClock: { remaining: convertTime(5, "minutes", "nanoseconds") } },
  opponent: { playerID: '', playerName: '', color: 'black', playerClock: { remaining: convertTime(5, "minutes", "nanoseconds")} },

  setPlayer: (player: Player) => set({ player }),
  setOpponent: (opponent: Player) => set({ opponent }),

  setConnection: (lobbyType, player, onMessage) => {
    const conn = join(lobbyType, player, onMessage)
    set({ conn })
    return conn
  },

  handler: (payload : ServerResponse) => {
    switch (payload.type) {
      case 'update':{
        set({ gameState: payload.data});
        break;
      }
      case 'chat':
        console.log('Chat message:', payload)
        break
        
    }
  },
  sendClientRequest: (socket: WebSocket, req: ClientRequest) => {
    if (socket.readyState !== WebSocket.OPEN) return;
    socket.send(JSON.stringify(req));
  },
  sendLobbyRequest: (socket: WebSocket, req: LobbyRequest) => {
    if (socket.readyState !== WebSocket.OPEN) return;
    socket.send(JSON.stringify(req));
  }

}))


function join(lobbyType: string, player: Player, onMessage : (data: ServerResponse) => void) : WebSocket {
  const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}/join-lobby`);

  socket.onopen = () => {
    console.log(player)
    const lobbyReq : LobbyRequest = {
      lobbyType: lobbyType,
      player: player,
    };
    socket.send(JSON.stringify(lobbyReq));
  };

  socket.onmessage = (event) => {
    const payload = JSON.parse(event.data);
    console.log(payload);
    onMessage(payload);
  }

  socket.onerror = () => {
    //TODO: show popup that shows error status
  };

  socket.onclose = () => {
    //TODO: show popup that signals end of game
  };

  return socket
}

