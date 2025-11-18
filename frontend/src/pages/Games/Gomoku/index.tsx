import { useEffect } from "react";
import { Outlet, RouteObject } from "react-router-dom";


import { useAuthStore } from "../../../stores/useAuthStore";
import { useGameStore } from "../../../stores/useGomokuStore";
import { LoginRedirectModal } from "../../Login/LoginRedirectModal";
import { GomokuLobby } from "./GomokuLobby";
import GomokuGame from "./GomokuGame";

function Gomoku() {
  const { checkAuth, isAuthenticated } = useAuthStore();
  const { setPlayer, player } = useGameStore();

  useEffect(() => {
    const check = async () => {
      const success = await checkAuth(() => {});
      if (!success) return;

      const user = useAuthStore.getState().user;
      if (!user) return;

      setPlayer({...player, playerID: user.id, playerName: user.username});
    };
    check()
  }, [])

  if (!isAuthenticated) return <LoginRedirectModal />

  return <Outlet />

}

export default function GomokuRoutes() : RouteObject {
  const routes : RouteObject = {
    path: "/games/gomoku", 
    element: <Gomoku />,
    children: [
      { index: true, element: <GomokuLobby />},
      { path: ":gameID", element:  <GomokuGame />},
    ]
  }
  return routes
}