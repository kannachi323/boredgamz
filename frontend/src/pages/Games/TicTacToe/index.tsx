
import { RouteObject } from "react-router-dom"

export function TicTacToe() {
  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50 text-white">
      <div className="text-3xl font-bold">
        COMING SOON
      </div>
    </div>
  );
}



export default function TicTacToeRoutes() : RouteObject {
  const routes : RouteObject = {
    path: "/games/tictactoe", 
    element: <TicTacToe />,
  }
  return routes
}