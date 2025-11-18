package api

import (
	"boredgamz/manager"
	"boredgamz/manager/gomoku"
	"boredgamz/utils"
	"encoding/json"
	"log"
	"net/http"
)

func JoinGomokuLobby(lm *manager.LobbyManager) http.HandlerFunc {
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

        var reqBody gomoku.GomokuLobbyRequest
        if err := json.Unmarshal(msg, &reqBody); err != nil {
            log.Println("Invalid join message format:", err)
            return
        }
        log.Println(lm.Lobbies)
        gomokuLobby, ok := lm.GetLobby(reqBody.LobbyType)
		if !ok {
			log.Println("Lobby not found:", reqBody.LobbyType)
			return
		}
        player := manager.NewPlayer(
            reqBody.Player.PlayerID,
            reqBody.Player.PlayerName, 
            reqBody.Player.Color,
            reqBody.Player.Clock,
            conn,
        )
        gomokuLobby.AddPlayer(player)
		players, matched := gomokuLobby.MatchPlayers()
		
        if !matched { return }

        p1 := players[0]
		p2 := players[1]


        room := gomoku.NewGomokuRoom(p1, p2, reqBody.LobbyType)
        room.Start()
        
        room.Broadcast(&gomoku.GomokuServerResponse{
            Type: "update",
            Data: room.GameState,
        })
       
    }
}