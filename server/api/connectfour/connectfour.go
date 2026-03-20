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
		if !websocket.IsWebSocketUpgrade(r) {
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
			log.Println("Error reading connect four join message:", err)
			_ = conn.Close()
			return
		}

		var req cf.ConnectFourClientRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			writeWSError(conn, "invalid join message format")
			return
		}

		var reqBody cf.ConnectFourLobbyData
		if err := json.Unmarshal(req.Data, &reqBody); err != nil {
			writeWSError(conn, "invalid join message data")
			return
		}

		lobbyController, ok := lm.GetLobby(cf.LobbyID)
		if !ok {
			writeWSError(conn, "connect four lobby not found")
			return
		}

		player := core.NewPlayer(
			reqBody.PlayerID,
			reqBody.PlayerName,
			reqBody.PlayerColor,
			core.NewPlayerClock(reqBody.TimeControl),
			conn,
		)

		lobbyController.AddPlayer(player)
	}
}

func ReconnectToConnectFourRoom(lm *core.LobbyManager, db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) {
			http.Error(w, "Expected WebSocket upgrade", http.StatusUpgradeRequired)
			return
		}

		_ = db

		conn, err := utils.UpgradeConnection(w, r)
		if err != nil {
			http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
			return
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading reconnect message:", err)
			_ = conn.Close()
			return
		}

		var req cf.ConnectFourClientRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			writeWSError(conn, "invalid reconnect format")
			return
		}

		var reqBody cf.ConnectFourReconnectData
		if err := json.Unmarshal(req.Data, &reqBody); err != nil {
			writeWSError(conn, "invalid reconnect data")
			return
		}

		lobbyController, ok := lm.GetLobby(cf.LobbyID)
		if !ok {
			writeWSError(conn, "connect four lobby not found")
			return
		}

		lobby, ok := lobbyController.(*cf.ConnectFourLobby)
		if !ok {
			writeWSError(conn, "invalid connect four lobby")
			return
		}

		roomController, ok := lobby.RoomManager.GetPlayerRoom(reqBody.PlayerID)
		if !ok {
			writeWSError(conn, "player room not found")
			return
		}

		room, ok := roomController.(*cf.ConnectFourRoom)
		if !ok {
			writeWSError(conn, "invalid connect four room")
			return
		}

		if err := room.ReconnectPlayer(reqBody.PlayerID, conn); err != nil {
			writeWSError(conn, "failed to reconnect player")
			return
		}
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
		_ = json.NewEncoder(w).Encode(game)
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
		_ = json.NewEncoder(w).Encode(games)
	}
}

func writeWSError(conn *websocket.Conn, message string) {
	res := map[string]any{
		"type": "error",
		"data": map[string]string{"message": message},
	}

	resBytes, err := json.Marshal(res)
	if err == nil {
		_ = conn.WriteMessage(websocket.TextMessage, resBytes)
	}

	_ = conn.Close()
}
