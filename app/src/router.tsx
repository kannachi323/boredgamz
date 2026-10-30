import { createBrowserRouter, Navigate } from "react-router-dom";
import App from "./App";
import { ErrorPage } from "./pages/ErrorPage";
import { GameRoute } from "./pages/GameRoute";
import Home from "./pages/Home";
import { AuthPage } from "./pages/AuthPage";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Home /> },
      { path: "games", element: <Navigate to="/" replace /> },
      { path: "games/:slug", element: <GameRoute /> },
      { path: "login", element: <AuthPage mode="login" /> },
      { path: "signup", element: <AuthPage mode="signup" /> },
      { path: "*", element: <ErrorPage /> },
    ],
  },
]);

export default router;
