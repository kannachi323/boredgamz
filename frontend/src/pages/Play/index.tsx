import { useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";


import { useAuthStore } from "../../stores/useAuthStore";
import { PlayQuick } from "./PlayQuick";
import { PlayRanked } from "./PlayRanked";
import { PlayCustom } from "./PlayCustom";



export default function Play() {
  const { isAuthenticated, user, checkAuth } = useAuthStore();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated || !user) {
      checkAuth(() => navigate('/login'));
    }
  })

  const { mode } = useParams();
  
  switch (mode) {
    case "quick": return <PlayQuick />;
    case "ranked": return <PlayRanked />;
    case "custom": return <PlayCustom />;
    default: return null;
  }
}
