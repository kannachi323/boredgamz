package api

import (
	"boredgamz/core"
	"boredgamz/core/gomoku"
	"boredgamz/utils"
	"encoding/json"
	"log"
	"net/http"
)

func JoinGomokuLobby(lm *core.LobbyManager) http.HandlerFunc {
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

        var req gomoku.GomokuClientRequest
        if err := json.Unmarshal(msg, &req); err != nil {
            log.Println("Invalid join message format:", err)
            return
        }

        var reqBody gomoku.GomokuLobbyData
        if err := json.Unmarshal(req.Data, &reqBody); err != nil {
            log.Println("Invalid join message data:", err)
            return
        }
       
        gomokuLobby, ok := lm.GetLobby(reqBody.LobbyType)
		if !ok {
			log.Println("Lobby not found:", reqBody.LobbyType)
			return
		}
        player := core.NewPlayer(
            reqBody.Player.PlayerID,
            reqBody.Player.PlayerName, 
            reqBody.Player.Color,
            reqBody.Player.Clock,
            conn,
        )
        gomokuLobby.AddPlayer(player)
    }
}