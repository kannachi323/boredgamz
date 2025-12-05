package connectfour

import (
	"boredgamz/core"
	"encoding/json"
	"log"
	"net/http"

	cf "boredgamz/core/connectfour"
	"boredgamz/utils"
)

func JoinConnectFourLobby(lm *core.LobbyManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure WebSocket upgrade
		if r.Header.Get("Connection") != "Upgrade" && r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "Expected WebSocket upgrade", http.StatusUpgradeRequired)
			return
		}

		// Upgrade connection
		conn, err := utils.UpgradeConnection(w, r)
		if err != nil {
			http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
			return
		}

		// Read initial join message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading join message:", err)
			return
		}

		// Parse client request
		var req cf.ConnectFourClientRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("Invalid join message format:", err)
			return
		}

		// Parse lobby data
		var reqBody cf.ConnectFourLobbyData
		if err := json.Unmarshal(req.Data, &reqBody); err != nil {
			log.Println("Invalid join message data:", err)
			return
		}

		// Get lobby from manager
		connectFourLobby, ok := lm.GetLobby(reqBody.LobbyType)
		if !ok {
			log.Println("Lobby not found:", reqBody.LobbyType)
			return
		}

		// Create player and add to lobby
		player := core.NewPlayer(
			reqBody.Player.PlayerID,
			reqBody.Player.PlayerName,
			reqBody.Player.Color,
			reqBody.Player.Clock,
			conn,
		)
		connectFourLobby.AddPlayer(player)
	}
}
