import { create } from 'zustand'
import { ServerResponse, GameState, Player, ClientRequest, Move, Board, AnalysisState } from '../pages/Games/Gomoku/GomokuTypes.tsx'
import { convertTime } from '../utils.ts'

interface GomokuStore {
  gameState: GameState | null
  conn: WebSocket | null
  player: Player
  opponent: Player
  analysis: AnalysisState
  startAnalysis: () => void
  exitAnalysis: () => void
  setAnalysisIndex: (idx: number) => void
  
  setPlayer: (player: Player) => void
  setOpponent: (opponent: Player) => void
  setConnection: (lobbyType: string, player: Player, onMessage : (data: ServerResponse) => void) => void
  handler: (payload: ServerResponse) => void
  send: (socket: WebSocket | null, data: ClientRequest) => void

}


export const useGomokuStore = create<GomokuStore>((set, get) => ({
  gameState: null,
  conn: null,
  player: { playerID: '', playerName: '', color: 'black', playerClock: { remaining: convertTime(5, "minutes", "nanoseconds") } },
  opponent: { playerID: '', playerName: '', color: 'black', playerClock: { remaining: convertTime(5, "minutes", "nanoseconds")} },
  analysis: { moves: [], board: null, active: false, index: 0 },

  startAnalysis: () => {
    const {gameState } = get();
    const moves = gameState?.moves || []
    set({
      analysis: {
        moves: moves,
        active: true,
        index: 0,
        board: buildBoardFromMoves(moves, 0),
      }
    });
  },

  exitAnalysis: () => {
    const { gameState } = get();
    const moves = gameState?.moves || []
    set({
      analysis: {
        moves: moves,
        active: false,
        index: moves.length - 1,
        board: buildBoardFromMoves(moves, moves.length - 1),
      }
    });
  },

  setAnalysisIndex: (idx: number) => {
    const { gameState } = get();
    const moves = gameState?.moves || []
    set({
      analysis: {
        moves: moves,
        active: true,
        index: idx,
        board: buildBoardFromMoves(moves, idx),
      }
    });
  },

  setPlayer: (player: Player) => set({ player }),
  setOpponent: (opponent: Player) => set({ opponent }),

  setConnection: (lobbyType, player, onMessage) => {
    const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}/join-gomoku-lobby`);

    socket.onopen = () => {
      socket.send(JSON.stringify({
        type: "lobby",
        data: {
          lobbyType: lobbyType,
          player: player,
        }
      }));
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
    
    set({ conn: socket })
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
  send: (socket: WebSocket | null, req: ClientRequest) => {
    if (!socket || socket.readyState !== WebSocket.OPEN) return;
    socket.send(JSON.stringify(req));
  },
}))

export function buildBoardFromMoves(
  moves: Move[],
  until: number,
  size: number = 9
): Board {

  // Initialize an empty board
  const stones: Board["stones"] = Array.from({ length: size }, () =>
    Array.from({ length: size }, () => ({ color: null }))
  );

  let numStones = 0;

  // Apply moves sequentially
  for (let i = 0; i <= until && i < moves.length; i++) {
    const m = moves[i];

    stones[m.row][m.col] = { color: m.color };
    numStones++;
  }

  return {
    stones,
    size,
    numStones
  };
}