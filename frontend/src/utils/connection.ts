import {ServerMessage, Player} from '../types'


export function createConnection(player: Player, onMessage : (data: ServerMessage) => void) {
  const socket = new WebSocket(`ws://localhost:3000/join-lobby?type=9x9`);

  socket.onopen = () => {
    //TODO: show a popup that starts the game

  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log(data)
    onMessage(data)
  }

  socket.onerror = () => {
    //TODO: show popup that shows error status
  };

  socket.onclose = () => {
    //TODO: show popup that signals end of game
  };
  
  return socket;
}

export function sendData(socket: WebSocket, data: object) {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(data));
  } else {
    console.error("WebSocket is not open. ReadyState:", socket.readyState);
  }
}

