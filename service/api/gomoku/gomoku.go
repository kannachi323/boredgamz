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
	"time"

	"github.com/gorilla/websocket"
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
        
        lobbyController, ok := lm.GetLobby(gomoku.GetGomokuLobbyID(reqBody.Name, reqBody.Mode))
		if !ok {
			log.Println("Lobby not found:", reqBody.Name)
			return
		}

        gomokuLobby, ok := lobbyController.(*gomoku.GomokuLobby)
        if !ok {
            return
        }

        player := core.NewPlayer(
            reqBody.PlayerID,
            reqBody.PlayerName, 
            reqBody.PlayerColor,
            core.NewPlayerClock(reqBody.TimeControl),
            conn,
        )

        player.StartPlayer()
        gomokuLobby.AddPlayer(player)

        player.Conn.SetCloseHandler(func(code int, text string) error {
            message := websocket.FormatCloseMessage(code, text)
            player.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
            gomokuLobby.RemovePlayer(player)
            return nil
        })
    }
}

func ReconnectToGomokuRoom(lm *core.LobbyManager, db *db.Database) http.HandlerFunc {
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
            log.Println("Error reading reconnect:", err)
            return
        }

        var req gomoku.GomokuClientRequest
        if err := json.Unmarshal(msg, &req); err != nil {
            log.Println("Invalid reconnect format:", err)
            return
        }

        var reqBody gomoku.GomokuReconnectData
        if err := json.Unmarshal(req.Data, &reqBody); err != nil {
            log.Println("Invalid reconnect data:", err)
            return
        }

    
        lobbyController, ok := lm.GetLobby(reqBody.LobbyID)
		if !ok {
			log.Println("Lobby not found:", reqBody.LobbyID)
			return
		}

        gomokuLobby, ok := lobbyController.(*gomoku.GomokuLobby)
        if !ok { 
            log.Println("not a valid gomoku lobby")
            return 
        }

    
        roomController, ok := gomokuLobby.RoomManager.GetPlayerRoom(reqBody.PlayerID)
        if !ok { 
            log.Println("player not found")
            return 
        }

        room, ok := roomController.(*gomoku.GomokuRoom)
        if !ok { 
            log.Println("not a valid gomoku room")
            return 
        }


        log.Println("reconnecting player")
        room.ReconnectPlayer(reqBody.PlayerID, conn)
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