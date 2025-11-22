package gomoku

import (
	"boredgamz/core"
	"boredgamz/core/gomoku"
	gomokucore "boredgamz/core/gomoku"
	"boredgamz/db"
	gomokudb "boredgamz/db/gomoku"
	"boredgamz/utils"
	"encoding/json"
	"log"
	"net/http"
)

func JoinGomokuLobby(lm *core.LobbyManager) http.HandlerFunc {
	 return func(w http.ResponseWriter, r *http.Request) {
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

        var clientReq gomokucore.GomokuClientRequest
        if err := json.Unmarshal(msg, &clientReq); err != nil {
            log.Println("Invalid join message format:", err)
            return
        }

        var lobbyData gomokucore.GomokuLobbyData
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

func PostGomokuGame(db* db.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req gomoku.GomokuClientRequest
        
        err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
        if req.Type != "save" {
            http.Error(w, "Invalid request type", http.StatusBadRequest)
        }

        var gameState gomokucore.GomokuGameState
        if err := json.Unmarshal(req.Data, &gameState); err != nil {
            http.Error(w, "Invalid game state data", http.StatusBadRequest)
        }

        err = gomokudb.InsertGame(db, &gameState)
        if err != nil {
            http.Error(w, "Failed to save game", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
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