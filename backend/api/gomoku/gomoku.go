package gomoku

import (
	"boredgamz/core"
	"boredgamz/core/gomoku"
	"boredgamz/db"
	gomokudb "boredgamz/db/gomoku"
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


func GetGomokuGame(db *db.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        gameID := r.URL.Query().Get("gameID")
        if gameID == "" {
            http.Error(w, "missing gameID", http.StatusBadRequest)
            return
        }

        game, err := gomokudb.GetGameByID(db, gameID)
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

func GetGomokuGames(db *db.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        playerID := r.URL.Query().Get("playerID")
        if playerID == "" {
            http.Error(w, "missing playerID", http.StatusBadRequest)
            return
        }

        games, err := gomokudb.GetGamesByPlayerID(db, playerID)
        if err != nil {
            http.Error(w, "internal error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(games)
    }
}