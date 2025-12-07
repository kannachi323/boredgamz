import { create } from 'zustand'
import { ServerResponse, GameState, Player, ClientRequest, Move, Board, AnalysisState, GameStateRow, Stone, LobbyRequest, ReconnectRequest } from '../../pages/Games/Gomoku/types.tsx'
import { createPlayer, createLobbyRequest } from './utils.ts'

interface GomokuStore {
  lobbyRequest: LobbyRequest
  gameState: GameState | null
  conn: WebSocket | null
  player: Player
  opponent: Player
  analysis: AnalysisState
  showGameEndModal: boolean
 


  setLobbyRequest: (lobbyRequest : LobbyRequest) => void
  setGameState: (gameState: GameState) => void
  setPlayer: (player: Player) => void
  setOpponent: (opponent: Player) => void
  startAnalysis: () => void
  exitAnalysis: () => void
  setAnalysisIndex: (idx: number) => void
  loadGame: (gameID: string) => Promise<void>
  setConnection: (lobbyRequest: LobbyRequest) => void
  closeConnection: () => void;
  reconnect: (gameID: string, playerID: string) => void
  handler: (payload: ServerResponse) => void
  send: (socket: WebSocket | null, data: ClientRequest) => void
  refreshPlayers: () => void
  buildBoardFromMoves: (size: number, moves: Move[], end: number) => Board | null
  buildGameState: (data: GameStateRow) => GameState | null
}


export const useGomokuStore = create<GomokuStore>((set, get) => ({
  lobbyRequest: createLobbyRequest(),
  gameState: null,
  conn: null,
  player: createPlayer(),
  opponent: createPlayer(),
  analysis: { moves: [], board: null, active: false, index: 0 },
  showGameEndModal: false,

  setLobbyRequest: (lobbyRequest : LobbyRequest) => set({ lobbyRequest }),
  setGameState: (gameState: GameState) => set({ gameState }),
  setPlayer: (player: Player) => set({ player }),
  setOpponent: (opponent: Player) => set({ opponent }),


  startAnalysis: () => {
    const { setAnalysisIndex } = get();
    setAnalysisIndex(-1);
  },

  exitAnalysis: () => {
    const { gameState, buildBoardFromMoves } = get();
    const moves = gameState?.moves || []
    set({
      analysis: {
        moves: moves,
        active: false,
        index: moves.length - 1,
        board: buildBoardFromMoves(gameState?.board?.size || -1, moves, moves.length - 1),
      }
    });
  },

  setAnalysisIndex: (idx: number) => {
    const { gameState, buildBoardFromMoves } = get();
    const moves = gameState?.moves || []
    set({
      analysis: {
        moves: moves,
        active: true,
        index: idx,
        board: buildBoardFromMoves(gameState?.board?.size || -1, moves, idx),
      }
    });
  },

  loadGame: async (gameID: string): Promise<void> => {
    const { buildGameState } = get();
    const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/gomoku/game?gameID=${gameID}`, {
      method: "GET",
      credentials: "include",
    });
    if (res.ok) {
      const data = await res.json();
      const newGameState = buildGameState(data as GameStateRow);
      set({ gameState: newGameState as GameState });
    } else {
      console.error("Failed to fetch game");
    }
  },

  setConnection: (lobbyRequest) => {
    const oldConn = get().conn;
    if (oldConn && oldConn.readyState !== WebSocket.CLOSED) {
        oldConn.close(1000);
    }
    const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}/join-gomoku-lobby`);
    socket.onopen = () => {
        console.log("WebSocket connected");
        socket.send(JSON.stringify(lobbyRequest));
    };
    socket.onmessage = (event) => {
        const payload = JSON.parse(event.data);
        console.log(payload);
        get().handler(payload);
    };
    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
    socket.onclose = (event) => {
      console.log("WebSocket closed:", event.code, event.reason);
      set({ conn: null });
    };
    
    set({ conn: socket });
  },
  closeConnection: () => {
    const conn = get().conn;
    if (!conn) return;
    
    if (conn.readyState === WebSocket.OPEN || conn.readyState === WebSocket.CONNECTING) {
        conn.close(1000, "User cancelled");
    } else {
        set({ conn: null });
    }
  },
  reconnect: (gameID: string, playerID: string) => {
    if (get().conn && get().conn!.readyState !== WebSocket.CLOSED) { return }
    console.log("reconnecting to game:", gameID, "as player:", playerID);
    const reconnectRequest : ReconnectRequest = {
      type: "reconnect",
      data: {
        playerID: playerID,
        gameID: gameID,
      }
    };
    const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}/reconnect-gomoku-room`);
    socket.onopen = () => {
      console.log("WebSocket connected");
      socket.send(JSON.stringify(reconnectRequest));
    };
    socket.onmessage = (event) => {
        const payload = JSON.parse(event.data);
        console.log(payload);
        get().handler(payload);
    };
    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
    socket.onclose = (event) => {
      console.log("WebSocket closed:", event.code, event.reason);
      set({ conn: null });
    };
    
    set({ conn: socket });
    
  },

  handler: (payload : ServerResponse) => {
    switch (payload.type) {
      case 'update':{
        console.log(payload);
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

  refreshPlayers: () => {
    const { gameState, player, setPlayer, setOpponent } = get();
    if (!gameState) return;
    const p1 = gameState.players[0];
    const p2 = gameState.players[1];
    setPlayer(p1.playerID === player.playerID ? p1 : p2);
    setOpponent(p1.playerID === player.playerID ? p2 : p1);
  },

  buildBoardFromMoves: (size: number, moves: Move[], end: number) => {
    if (size == -1) return null;
    const stones: Stone[][] = Array.from({ length: size }, () =>
      Array.from({ length: size }, () => ({ color: null }))
    );

    let numStones = 0;

    for (let i = 0; i <= end && i < moves.length; i++) {
      const m = moves[i];
      stones[m.row][m.col] = { color: m.color };
      numStones++;
    }
    
    return { stones, size, numStones }
  },

  buildGameState: (data: GameStateRow) => {
    const { buildBoardFromMoves } = get();
    
    const newBoard = buildBoardFromMoves(data.boardSize, data.moves, data.moves.length - 1);
    if (!newBoard) { return null }

    console.log(newBoard);

    const newPlayers: Player[] = data.players.map((p) => ({
      playerID: p.playerID,
      playerName: p.playerName,
      color: p.color,
      playerClock: null,
    }));

    const winner: Player | null = data.winner && {
      playerID: data.winner.playerID,
      playerName: data.winner.playerName,
      color: data.winner.color,
      playerClock: null,
    }
    
    const newGameState: GameState = {
      gameID: data.gameID,
      board: newBoard,
      size: data.boardSize,
      players: newPlayers,
      turn: "",
      status: {
        result: data.result,
        code: "offline",
        winner: winner,
      },
      lastMove: data.moves.length > 0 ? data.moves[data.moves.length - 1] : null,
      moves: data.moves,
    };

    console.log(newGameState);

    return newGameState;
  },

}));