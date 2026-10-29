import { createBrowserRouter } from 'react-router-dom'
import App from './App'
import { ErrorPage } from './pages/ErrorPage'
import Games from './pages/Games'
import Home from './pages/Home'
import { ChessGame } from './pages/Games/Chess'
import { ConnectFourGame } from './pages/Games/ConnectFour'
import { GomokuGamePage } from './pages/Games/Gomoku/GomokuGamePage'
import { MahjongGame } from './pages/Games/Mahjong'
import { PokerGame } from './pages/Games/Poker'
import { TicTacToeGame } from './pages/Games/TicTacToe'
import Login from './pages/Login'
import SignUp from './pages/Signup'

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Home /> },
      { path: 'games', element: <Games /> },
      { path: 'games/tictactoe', element: <TicTacToeGame /> },
      { path: 'games/connect-four', element: <ConnectFourGame /> },
      { path: 'games/gomoku', element: <GomokuGamePage /> },
      { path: 'games/chess', element: <ChessGame /> },
      { path: 'games/poker', element: <PokerGame /> },
      { path: 'games/mahjong', element: <MahjongGame /> },
      { path: 'login', element: <Login /> },
      { path: 'signup', element: <SignUp /> },
    ],
  },
])

export default router
