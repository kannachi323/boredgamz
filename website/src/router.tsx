import { createBrowserRouter } from 'react-router-dom';

import App from './App';

//PAGES 
import Games from './pages/Games';
import Home from './pages/Home';
import Signup from './pages/Signup';
import Login from './pages/Login';
import Community from './pages/Community';
import About from './pages/About';
import Help from './pages/Help';
import Donate from './pages/Donate';

//GAME ROUTES
import GomokuRoutes from './pages/Games/Gomoku';
import TicTacToeRoutes from './pages/Games/TicTacToe';
import ConnectFourRoutes from './pages/Games/ConnectFour';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {index: true, element: <Home />},
      {path: '/games', element: <Games />},
      {path: "/help", element: <Help />},
      {path: "/donate", element: <Donate />},
      {path: "/community", element: <Community />},
      {path: '/about', element: <About />},
      {path: '/signup', element: <Signup />},
      {path: '/login', element: <Login />},

      GomokuRoutes(),
      TicTacToeRoutes(),
      ConnectFourRoutes(),
    ],
  },
]);

export default router;
