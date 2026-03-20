import { create } from 'zustand'
import { ServerResponse, GameState, Player, ClientRequest, Move, Board, AnalysisState, GameStateRow, Stone, LobbyRequest, ReconnectRequest, ChatMessage } from '../../pages/Games/Gomoku/types.tsx'
import { createPlayer, createLobbyRequest } from './utils.ts'

type GomokuConnectionStatus = "idle" | "connecting" | "queued" | "reconnecting" | "connected" | "error"

type ConnectionIntent =
  | { type: "lobby"; payload: LobbyRequest }
  | { type: "reconnect"; payload: ReconnectRequest }

interface GomokuStore {
  lobbyRequest: LobbyRequest
  gameState: GameState | null
  conn: WebSocket | null
  connectionStatus: GomokuConnectionStatus
  connectionError: string | null
  lastConnectionIntent: ConnectionIntent | null
  player: Player
  opponent: Player
  messages: ChatMessage[]
  analysis: AnalysisState
  showGameEndModal: boolean
 


  setLobbyRequest: (lobbyRequest : LobbyRequest) => void
  setGameState: (gameState: GameState) => void
  setPlayer: (player: Player) => void
  setOpponent: (opponent: Player) => void
  setMessages: (messages: ChatMessage[]) => void
  startAnalysis: () => void
  exitAnalysis: () => void
  setAnalysisIndex: (idx: number) => void
  loadGame: (gameID: string) => Promise<void>
  setConnection: (lobbyRequest: LobbyRequest) => void
  closeConnection: () => void;
  reconnect: (lobbyID: string, playerID: string) => void
  retryConnection: () => void
  clearConnectionError: () => void
  markConnectionError: (message: string) => void
  handler: (payload: ServerResponse) => void
  send: (data: ClientRequest) => void
  refreshPlayers: () => void
  buildBoardFromMoves: (size: number, moves: Move[], end: number) => Board | null
  buildGameState: (data: GameStateRow) => GameState | null
}


export const useGomokuStore = create<GomokuStore>((set, get) => ({
  lobbyRequest: createLobbyRequest(),
  gameState: null,
  conn: null,
  connectionStatus: "idle",
  connectionError: null,
  lastConnectionIntent: null,
  player: createPlayer(),
  opponent: createPlayer(),
  messages: [],
  analysis: { moves: [], board: null, active: false, index: 0 },
  showGameEndModal: false,

  setLobbyRequest: (lobbyRequest : LobbyRequest) => set({ lobbyRequest }),
  setGameState: (gameState: GameState) => set({ gameState }),
  setPlayer: (player: Player) => set({ player }),
  setOpponent: (opponent: Player) => set({ opponent }),
  setMessages: (messages: ChatMessage[]) => set({ messages }),
  clearConnectionError: () => set({ connectionError: null }),
  markConnectionError: (message: string) => set({
    connectionStatus: "error",
    connectionError: message,
  }),


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
      set({
        gameState: newGameState as GameState,
        messages: newGameState?.messages || [],
      });
    } else {
      console.error("Failed to fetch game");
    }
  },

  setConnection: (lobbyRequest) => {
    const oldConn = get().conn;
    if (oldConn && oldConn.readyState !== WebSocket.CLOSED) {
        oldConn.close(1000, "Switching session");
    }

    const intent: ConnectionIntent = { type: "lobby", payload: lobbyRequest };
    const joinPath = lobbyRequest.data.mode === "bots"
      ? "/join-gomoku-bot-lobby"
      : "/join-gomoku-lobby";
    const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}${joinPath}`);

  set({
    conn: socket,
    connectionStatus: "connecting",
    connectionError: null,
    lastConnectionIntent: intent,
    messages: [],
  });

    socket.onopen = () => {
      if (get().conn !== socket) return;
      console.log("WebSocket connected");
      set({ connectionStatus: "queued", connectionError: null });
        socket.send(JSON.stringify(lobbyRequest));
    };
    socket.onmessage = (event) => {
        if (get().conn !== socket) return;
        const payload = JSON.parse(event.data);
        console.log(payload);
        get().handler(payload);
    };
    socket.onerror = (error) => {
        if (get().conn !== socket) return;
        console.error("WebSocket error:", error);
        set({
          connectionStatus: "error",
          connectionError: "Could not connect to lobby. Please retry.",
        });
    };
    socket.onclose = (event) => {
      if (get().conn !== socket) return;
      console.log("WebSocket closed:", event.code, event.reason);
      const { gameState } = get();
      const isFinished = gameState?.status?.code === "offline";
      const isUserCancelled =
        event.code === 1000 &&
        (event.reason === "User cancelled" || event.reason === "Switching session");

      set((state) => ({
        conn: null,
        connectionStatus: isFinished
          ? (state.connectionStatus === "connected" ? "connected" : "idle")
          : isUserCancelled
            ? "idle"
            : "error",
        connectionError: isFinished || isUserCancelled
          ? null
          : state.connectionError || event.reason || "Connection closed unexpectedly. Please retry.",
      }));
    };
  },

  closeConnection: () => {
    const conn = get().conn;
    if (!conn) return;
    
    if (conn.readyState === WebSocket.OPEN || conn.readyState === WebSocket.CONNECTING) {
        conn.close(1000, "User cancelled");
    } else {
        set({ conn: null, connectionStatus: "idle", connectionError: null });
    }

    set({ connectionStatus: "idle", connectionError: null });
  },

  reconnect: (lobbyID: string, playerID: string) => {
    if (get().conn && get().conn!.readyState !== WebSocket.CLOSED) { return }
    console.log("reconnecting to lobby:", lobbyID, "as player:", playerID);
    const reconnectRequest : ReconnectRequest = {
      type: "reconnect",
      data: {
        playerID: playerID,
        lobbyID: lobbyID,
      }
    };

    const intent: ConnectionIntent = { type: "reconnect", payload: reconnectRequest };
    const socket = new WebSocket(`${import.meta.env.VITE_WEBSOCKET_ROOT}/reconnect-gomoku-room`);

    set({
      conn: socket,
      connectionStatus: "reconnecting",
      connectionError: null,
      lastConnectionIntent: intent,
      messages: [],
    });

    socket.onopen = () => {
      if (get().conn !== socket) return;
      console.log("WebSocket connected");
      socket.send(JSON.stringify(reconnectRequest));
    };
    socket.onmessage = (event) => {
        if (get().conn !== socket) return;
        const payload = JSON.parse(event.data);
        console.log(payload);
        get().handler(payload);
    };
    socket.onerror = (error) => {
        if (get().conn !== socket) return;
        console.error("WebSocket error:", error);
        set({
          connectionStatus: "error",
          connectionError: "Reconnect failed. Please retry.",
        });
    };
    socket.onclose = (event) => {
      if (get().conn !== socket) return;
      console.log("WebSocket closed:", event.code, event.reason);
      const { gameState } = get();
      const isFinished = gameState?.status?.code === "offline";
      const isUserCancelled =
        event.code === 1000 &&
        (event.reason === "User cancelled" || event.reason === "Switching session");

      set({
        conn: null,
        connectionStatus: isFinished || isUserCancelled ? "idle" : "error",
        connectionError: isFinished || isUserCancelled
          ? null
          : get().connectionError || event.reason || "Reconnect closed unexpectedly. Please retry.",
      });
    };

  },

  retryConnection: () => {
    const { lastConnectionIntent, setConnection, reconnect } = get();
    if (!lastConnectionIntent) {
      set({ connectionStatus: "idle", connectionError: null });
      return;
    }

    if (lastConnectionIntent.type === "lobby") {
      setConnection(lastConnectionIntent.payload);
      return;
    }

    reconnect(lastConnectionIntent.payload.data.lobbyID, lastConnectionIntent.payload.data.playerID);
  },

  handler: (payload : ServerResponse) => {
    switch (payload.type) {
      case 'update':{
        console.log(payload);
        set({
          gameState: payload.data as GameState,
          messages: (payload.data as GameState).messages || [],
          connectionStatus: "connected",
          connectionError: null,
        });
        get().refreshPlayers();
        break;
      }

      case 'error': {
        const message = typeof payload.data === 'object' && payload.data !== null && 'message' in payload.data
          ? String((payload.data as { message: unknown }).message)
          : 'Server reported a connection error. Please retry.';

        set({
          connectionStatus: 'error',
          connectionError: message,
        });
        break;
      }

      case 'chat':
        set((state) => {
          const nextMessages = [...state.messages, payload.data as ChatMessage];
          return {
            messages: nextMessages,
            gameState: state.gameState
              ? { ...state.gameState, messages: nextMessages }
              : state.gameState,
          };
        })
        break
        
    }
  },

  send: (req: ClientRequest) => {
    const socket = get().conn;
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      set({
        connectionStatus: "error",
        connectionError: "Connection is not open. Please retry.",
      });
      return;
    }
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
      openingRule: "freestyle",
      swapRuleEnabled: false,
      firstMoveCenterEnabled: false,
      players: newPlayers,
      turn: "",
      status: {
        result: data.result,
        code: "offline",
        winner: winner,
      },
      lastMove: data.moves.length > 0 ? data.moves[data.moves.length - 1] : null,
      moves: data.moves,
      messages: data.messages || [],
    };

    console.log(newGameState);

    return newGameState;
  },

}));
