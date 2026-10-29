package connectfour

import (
	"boredgamz/core"
	cf "boredgamz/core/connectfour"
	"boredgamz/db"
	cfdb "boredgamz/db/connectfour"
	"boredgamz/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func JoinConnectFourLobby(lm *core.LobbyManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure WebSocket upgrade
		if !websocket.IsWebSocketUpgrade(r) {
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
		player.StartPlayer()
		connectFourLobby.AddPlayer(player)

		player.Conn.SetCloseHandler(func(code int, text string) error {
			player.ClosePlayer()
			connectFourLobby.RemovePlayer(player)
			return nil
		})
	}
}

func GetConnectFourGame(db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.URL.Query().Get("gameID")
		if gameID == "" {
			http.Error(w, "missing gameID", http.StatusBadRequest)
			return
		}

		game, err := cfdb.GetGameByID(db, gameID)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if game == nil {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func GetConnectFourGames(db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := r.URL.Query().Get("playerID")
		if playerID == "" {
			http.Error(w, "missing playerID", http.StatusBadRequest)
			return
		}

		games, err := cfdb.GetGamesByPlayerID(db, playerID)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}
