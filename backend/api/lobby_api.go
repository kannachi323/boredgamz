package api

import (
	"encoding/json"
	"log"
	"net/http"
	"playgomoku/backend/game"
	"playgomoku/backend/manager"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //remove this in production
	},
}

func JoinLobby(lm *manager.LobbyManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Connection") != "Upgrade" && r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "Expected WebSocket upgrade", http.StatusUpgradeRequired)
			return
		}

        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
            return
        }

        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading join message:", err)
            return
        }

        var reqBody manager.LobbyRequest
        if err := json.Unmarshal(msg, &reqBody); err != nil {
            log.Println("Invalid join message format:", err)
            return
        }

        lobbyType := reqBody.LobbyType
        lobby, _:= lm.GetLobby(lobbyType)

        player := &game.Player{
            PlayerID:      reqBody.Player.PlayerID,
            PlayerName:    reqBody.Player.PlayerName,
            Color:         reqBody.Player.Color,
            Clock:      reqBody.Player.Clock,
            Conn:     conn,
            Incoming: make(chan []byte, 10),
            Outgoing: make(chan []byte, 10),
        }

        lm.AddPlayerToQueue(lobby, player)
        players, success := lm.MatchPlayers(lobby)
        
        if success {
            rm := lobby.RoomManager

		    room := rm.CreateNewRoom(players[0], players[1], lobby.LobbyType)

            rm.Broadcast(room, &manager.ServerResponse{
                Type: "update",
                Data: room.GameState,
            })
            rm.StartRoom(room)
        }
    }
}
