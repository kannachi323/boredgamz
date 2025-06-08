import { createBrowserRouter } from 'react-router-dom';

import App from './App';

import Play from "./pages/Play";
import Home from "./pages/Home";
import Signup from './pages/Signup';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {index: true, element: <Home />},
      {path: "play", element: <Play />},
    ],
  },
  {
    path: '/signup',
    element: <Signup />,
  }
]);

export default router;
