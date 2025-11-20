package api

import (
	"boredgamz/core"
	"boredgamz/core/gomoku"
	"boredgamz/utils"
	"encoding/json"
	"log"
	"net/http"
)

func JoinGomokuLobby(lm *core.Lobbycore) http.HandlerFunc {
	 return func(w http.ResponseWriter, r *http.Request) {
        log.Println("New join gomoku lobby request")
        if r.Header.Get("Connection") != "Upgrade" && r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "Expected WebSocket upgrade", http.StatusUpgradeRequired)
			return
		}

        conn, err := utils.UpgradeConnection(w, r)
        if err != nil {
            http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
            return
        }

        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading join message:", err)
            return
        }

        var clientReq gomoku.GomokuClientRequest
        if err := json.Unmarshal(msg, &clientReq); err != nil {
            log.Println("Invalid join message format:", err)
            return
        }

        var lobbyData gomoku.GomokuLobbyData
        if err := json.Unmarshal(clientReq.Data, &lobbyData); err != nil {
            log.Println("Invalid lobby request data:", err)
            return
        }

        gomokuLobby, ok := lm.GetLobby(lobbyData.LobbyType)
		if !ok {
			log.Println("Lobby not found:", lobbyData.LobbyType)
			return
		}
        player := core.NewPlayer(
            lobbyData.Player.PlayerID,
            lobbyData.Player.PlayerName, 
            lobbyData.Player.Color,
            lobbyData.Player.Clock,
            conn,
        )
        gomokuLobby.AddPlayer(player)
		players, matched := gomokuLobby.MatchPlayers()
		
        if !matched { return }

        p1 := players[0]
		p2 := players[1]
        room := gomoku.NewGomokuRoom(p1, p2, lobbyData.LobbyType)
     

        p1.StartPlayer()
        p2.StartPlayer()
        room.Start()

        data, _ := json.Marshal(room.GameState)
        room.Broadcast(&gomoku.GomokuServerResponse{
            Type: "update",
            Data: data,
        })
       
    }
}