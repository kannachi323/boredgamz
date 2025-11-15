import { createBrowserRouter } from 'react-router-dom';

import App from './App';

//PAGES 
import Games from './pages/Games';
import Home from "./pages/Home";
import Signup from './pages/Signup';
import Login from './pages/Login';
import Community from './pages/Community';

//GAME ROUTES
import GomokuRoutes from './pages/Games/Gomoku';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {index: true, element: <Home />},
      {path: "/games", element: <Games />},
      {path: "/community", element: <Community />},
      {path: '/signup', element: <Signup />},
      {path: '/login', element: <Login />},

      GomokuRoutes(),
    ],
  },
]);

export default router;
