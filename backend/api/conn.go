package api

import (
	"log"
	"net/http"
	"playgomoku/backend/game"
	"playgomoku/backend/manager"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //remove this in production
	},
}

func JoinLobby(lm *manager.LobbyManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        lobbyType := r.URL.Query().Get("type")
        if lobbyType == "" {
            http.Error(w, "missing lobby type...", http.StatusBadRequest)
            return
        }

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

        player := &game.Player{
            ID:       uuid.NewString(),
            Conn:     conn,
            Incoming: make(chan []byte, 10),
            Outgoing: make(chan []byte, 10),
        }

        lobby, ok := lm.GetLobby(lobbyType)
        if !ok {
            http.Error(w, "cannot connect to lobby", http.StatusNotFound)
            return
        }

        manager.AddPlayerToQueue(lobby, player)
        room, accept := lobby.MatchPlayers()

        if !accept {
            // TODO: player continues waiting until timeout
            log.Print("waiting for more players...")
            return
        }

        room.Start()
        
    }
}

