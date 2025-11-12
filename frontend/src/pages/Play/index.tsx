import { useEffect } from "react";
import { useParams } from "react-router-dom";


import { useAuthStore } from "../../stores/useAuthStore";
import { PlayQuick } from "./PlayQuick";
import { PlayRanked } from "./PlayRanked";
import { PlayCustom } from "./PlayCustom";
import { useGameStore } from "../../stores/useGameStore";
import { LoginRedirectModal } from "../Login/LoginRedirectModal";



export default function Play() {
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

  const { mode } = useParams();

  if (!isAuthenticated) return <LoginRedirectModal />
  
  switch (mode) {
    case "quick": return <PlayQuick />;
    case "ranked": return <PlayRanked />;
    case "custom": return <PlayCustom />;
    default: return null;
  }
}
