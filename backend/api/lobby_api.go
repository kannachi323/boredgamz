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
			log.Println("Ignoring non-WebSocket request")
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
        player := &game.Player{
            PlayerID:       reqBody.Player.PlayerID,
            Color:         reqBody.Player.Color,
            Conn:     conn,
            Incoming: make(chan []byte, 10),
            Outgoing: make(chan []byte, 10),
        }

        lobby, ok := lm.GetLobby(lobbyType)
        if !ok {
            // create a new lobby if it doesn't exist
            log.Print("brooo")
            return
        }

        manager.AddPlayerToQueue(lobby, player)
        
        room, accept := lobby.MatchPlayers() //game state will already be created here

        if !accept {
            // TODO: player continues waiting until timeout
            log.Print("waiting for more players...")
        } else {
            // start the game room
            log.Print(room.Game.Players[0].PlayerID)
            go room.Start()
            room.Broadcast(&manager.ServerResponse{
                Type: "new",
                Data: room.Game,
            })
            log.Print("bro i got here yesss")
        }
    }
}
