import {ServerResponse, ClientRequest, Player } from '../types'

export function createConnection(lobbyType: string, player: Player, onMessage : (data: ServerResponse) => void) {
  const socket = new WebSocket(`ws://localhost:3000/join-lobby`);

  socket.onopen = () => {
    //TODO: show a popup that starts the game
    socket.send(JSON.stringify({
      type: "join",
      lobbyType: lobbyType,
      player: player,
    }));
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

export function sendData(socket: WebSocket, data: ClientRequest) {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(data));
  } else {
    console.error("WebSocket is not open. ReadyState:", socket.readyState);
  }
}

