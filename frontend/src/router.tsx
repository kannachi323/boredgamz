import { createBrowserRouter } from 'react-router-dom';

import App from './App';

import Play from "./pages/Play";
import Game from './features/Game';




import Home from "./pages/Home";
import Signup from './pages/Signup';
import Login from './pages/Login';


const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {index: true, element: <Home />},
      {path: "/play/:mode", element: <Play />,
        children: [
          { path: ":gameID", element: <Game /> },
        ]
      },
      {path: '/signup', element: <Signup />},
      {path: '/login', element: <Login />},
    ],
  },
  
]);

export default router;
