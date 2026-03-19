import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";
import { GomokuLobbyOptions } from "./GomokuLobbyOptions";
import { GomokuLobbyBoards } from "./GomokuLobbyBoards"
import { GomokuLobbyModes } from "@Gomoku/features/Lobby/GomokuLobbyModes";
import { GomokuLobbyWaiting } from "./GomokuLobbyWaiting";
import { GomokuScrollTooltip } from "./GomokuScrollToolTip";
import { GomokuLobbyRules } from "./GomokuLobbyRules";
import { useAuthStore } from "@/stores/useAuthStore";

export function GomokuLobby() {
  const { user } = useAuthStore();
  const {
    gameState,
    setConnection,
    lobbyRequest,
    setLobbyRequest,
    closeConnection,
    connectionStatus,
    connectionError,
    retryConnection,
  } = useGomokuStore();
  const navigate = useNavigate();

  const showWaitingModal =
    connectionStatus === "connecting" ||
    connectionStatus === "queued" ||
    connectionStatus === "reconnecting" ||
    connectionStatus === "error";

  useEffect(() => {
    if (!user?.id || !user?.username) return;
    if (lobbyRequest.data.playerID === user.id && lobbyRequest.data.playerName === user.username) {
      return;
    }

    setLobbyRequest({...lobbyRequest, data: {...lobbyRequest.data,
      playerName: user.username,
      playerID: user.id, 
    }})
  }, [user?.id, user?.username, lobbyRequest, setLobbyRequest])

  useEffect(() => {
    if (gameState?.gameID && gameState?.status.code === "online") {
      const timer = window.setTimeout(() => {
        navigate(`/games/gomoku/gomoku-${lobbyRequest.data.mode}-${lobbyRequest.data.name}-${lobbyRequest.data.openingRule}`);
      }, 2000);
      return () => clearTimeout(timer);
    }
  }, [gameState, navigate, lobbyRequest.data.mode, lobbyRequest.data.name, lobbyRequest.data.openingRule]);

  function handlePlayNow() {
    setConnection(lobbyRequest);
  }

  function handlePlayCancel() {
    closeConnection()
  }

  return (
    <div className="relative h-[92vh] w-full flex flex-col bg-[#1b1918]">
      <div className="flex-1 overflow-y-auto py-8">
        <div className="flex flex-col justify-start items-center gap-10 pb-8 px-6 md:px-10">
          <h1 className="text-6xl text-[#C3B299] font-bold text-center">Gomoku</h1>

          {/* Game Options */}
          <section className="flex flex-col items-center gap-1">
            <p className="text-lg text-[#C3B299] font-bold mb-1">Game</p>
            <div className="bg-[#433d3a] flex flex-row items-center justify-evenly p-3 rounded-xl gap-5">
              <GomokuLobbyOptions />
            </div>
          </section>

          {/* Mode */}
          <section className="flex flex-col items-center gap-1">
            <p className="text-lg text-[#C3B299] font-bold">Mode</p>
            <div className="bg-[#433d3a] p-3 rounded-xl flex flex-row justify-evenly gap-5">
              <GomokuLobbyModes/>
            </div>
          </section>

          {/* Rules */}
          <section className="flex flex-col items-center gap-1">
            <p className="text-lg text-[#C3B299] font-bold">Rules</p>
            <div className="bg-[#433d3a] p-3 rounded-xl flex flex-col justify-evenly gap-4">
              <GomokuLobbyRules />
            </div>
          </section>


          {/* Board */}
          <section className="flex flex-col items-center gap-1">
            <p className="text-lg text-[#C3B299] font-bold mb-1">Board</p>
            <div className="bg-[#433d3a] flex flex-row items-center justify-center p-3 rounded-xl gap-5">
              <GomokuLobbyBoards />
            </div>
          </section>


          {/* tooltip */}
          <div className="fixed bottom-24 right-4 text-[#C3B299] rounded-lg shadow-lg text-sm 
            font-semibold flex items-center gap-3 cursor-pointer">
            <GomokuScrollTooltip />
          </div>
        </div>
      </div>

      <div className="sticky bottom-0 w-full bg-[#1b1918]/95 backdrop-blur py-4 flex justify-center z-20">
        {/* Start / Play Button */}
        <button
          className="px-10 py-3 bg-[#C3B299] text-[#433d3a] font-bold rounded-lg hover:bg-[#d7c9b8] transition disabled:opacity-60 disabled:cursor-not-allowed"
          onClick={() => handlePlayNow()}
          disabled={showWaitingModal && connectionStatus !== "error"}
        >
          Play Now
        </button>
      </div>

      {showWaitingModal && (
        <GomokuLobbyWaiting
          onCancel={() => handlePlayCancel()}
          onRetry={() => retryConnection()}
          status={connectionStatus}
          mode={lobbyRequest.data.mode}
          errorMessage={connectionError}
        />
      )}
    </div>
  );
}
