import { useEffect } from "react";
import { Outlet } from "react-router-dom";
import { AppLayout } from "./Layout";
import { useAuthStore } from "./stores/useAuthStore";

function App() {
  const initialize = useAuthStore((state) => state.initialize);
  const status = useAuthStore((state) => state.status);
  useEffect(() => {
    void initialize();
  }, [initialize]);
  return (
    <AppLayout>
      {status === "checking" ? <div className="loading" role="status" aria-label="Loading"><span className="spinner" /></div> : <Outlet />}
    </AppLayout>
  );
}

export default App;
