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

	"github.com/gorilla/websocket"
)

func JoinGomokuLobby(lm *core.LobbyManager) http.HandlerFunc {
    return joinGomokuLobby(lm, false)
}

func JoinGomokuBotLobby(lm *core.LobbyManager) http.HandlerFunc {
    return joinGomokuLobby(lm, true)
}

func joinGomokuLobby(lm *core.LobbyManager, botsOnly bool) http.HandlerFunc {
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
            log.Println("Error reading join message:", err)
            _ = conn.Close()
            return
        }

        var req gomoku.GomokuClientRequest
        if err := json.Unmarshal(msg, &req); err != nil {
            log.Println("Invalid join message format:", err)
            writeWSError(conn, "invalid join message format")
            return
        }

        var reqBody gomoku.GomokuLobbyData
        if err := json.Unmarshal(req.Data, &reqBody); err != nil {
            log.Println("Invalid join message data:", err)
            writeWSError(conn, "invalid join message data")
            return
        }
        
        if reqBody.OpeningRule == "" {
            reqBody.OpeningRule = "freestyle"
        }

        if botsOnly {
            reqBody.Mode = "bots"
        } else if reqBody.Mode == "bots" {
            writeWSError(conn, "bots mode must use bot lobby endpoint")
            return
        }

        lobbyController, ok := lm.GetLobby(gomoku.GetGomokuLobbyID(reqBody.Name, reqBody.Mode, reqBody.OpeningRule))
		if !ok {
			log.Println("Lobby not found:", reqBody.Name)
            writeWSError(conn, "lobby not found")
			return
		}

        gomokuLobby, ok := lobbyController.(*gomoku.GomokuLobby)
        if !ok {
            writeWSError(conn, "invalid gomoku lobby")
            return
        }

        player := core.NewPlayer(
            reqBody.PlayerID,
            reqBody.PlayerName, 
            reqBody.PlayerColor,
            core.NewPlayerClock(reqBody.TimeControl),
            conn,
        )
        player.OpeningRule = reqBody.OpeningRule
        player.SwapRuleEnabled = reqBody.SwapRuleEnabled
        player.FirstMoveCenterEnabled = reqBody.FirstMoveCenterEnabled
        player.BotDifficulty = reqBody.BotDifficulty

        gomokuLobby.AddPlayer(player)
    }
}

func ReconnectToGomokuRoom(lm *core.LobbyManager, db *db.Database) http.HandlerFunc {
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
            log.Println("Error reading reconnect:", err)
            _ = conn.Close()
            return
        }

        var req gomoku.GomokuClientRequest
        if err := json.Unmarshal(msg, &req); err != nil {
            log.Println("Invalid reconnect format:", err)
            writeWSError(conn, "invalid reconnect format")
            return
        }

        var reqBody gomoku.GomokuReconnectData
        if err := json.Unmarshal(req.Data, &reqBody); err != nil {
            log.Println("Invalid reconnect data:", err)
            writeWSError(conn, "invalid reconnect data")
            return
        }

    
        normalizedLobbyID := gomoku.NormalizeLobbyID(reqBody.LobbyID)
        lobbyController, ok := lm.GetLobby(normalizedLobbyID)
		if !ok {
            log.Println("Lobby not found:", reqBody.LobbyID, "normalized as:", normalizedLobbyID)
            writeWSError(conn, "lobby not found")
			return
		}

        gomokuLobby, ok := lobbyController.(*gomoku.GomokuLobby)
        if !ok { 
            log.Println("not a valid gomoku lobby")
            writeWSError(conn, "invalid gomoku lobby")
            return 
        }

    
        roomController, ok := gomokuLobby.RoomManager.GetPlayerRoom(reqBody.PlayerID)
        if !ok { 
            log.Println("player not found")
            writeWSError(conn, "player room not found")
            return 
        }

        room, ok := roomController.(*gomoku.GomokuRoom)
        if !ok { 
            log.Println("not a valid gomoku room")
            writeWSError(conn, "invalid gomoku room")
            return 
        }


        log.Println("reconnecting player")
        if err := room.ReconnectPlayer(reqBody.PlayerID, conn); err != nil {
            log.Println("failed to reconnect player:", err)
            writeWSError(conn, "failed to reconnect player")
            return
        }
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