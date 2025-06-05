import { Stone } from "../types";

import { sendData } from "./connection";

export function createGame() {
  return {
    board: Array(0).fill(null), //default board
    size: 9, //default invalid size
    players: {
      p1: {
        playerID: "player1",
        color: "1",
      },
      p2: {
        playerID: "player2",
        color: "2",
      }
    },
    turn: "player1",
    status: "online"
  }
}

export function placeStone(socket: WebSocket, idx: number, setBoard: React.Dispatch<React.SetStateAction<Stone[]>>) {
  console.log('im sending a move');
  sendData(socket, {
    type: "move",
    move: idx,
  });

  setBoard(prev => {
    const newBoard = [...prev];
    newBoard[idx] = { color: 0, colorName: "black"};
    return newBoard;
  });
}
