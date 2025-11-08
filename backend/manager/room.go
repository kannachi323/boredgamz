package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"playgomoku/backend/game"
	"sync"
)

type Room struct {
	Player1 *game.Player
	Player2 *game.Player
	Game   *game.GameState
	closeOnce sync.Once
}

type RoomManager struct {
	playerRoomMap map[string]*Room
	mu	sync.RWMutex
}

func CreateRoomManager() *RoomManager {
	newRoomManager := &RoomManager{
		playerRoomMap: make(map[string]*Room),
	}

	return newRoomManager;
}

func (rm *RoomManager) Broadcast(r *Room, res *ServerResponse) error {

	_ = Send(r.Player1, res)
	_ = Send(r.Player2, res)

	return nil
}

func Send(p *game.Player, res *ServerResponse) error {
	msg, err := json.Marshal(res)
	if err != nil {
		return err
	}

	if p.Disconnected.Load() {
		return nil
	}
	select {
	case p.Outgoing <- msg:
		return nil
	default:
		return fmt.Errorf("failed to send message to player %s", p.PlayerID	)
	}
}

func (rm *RoomManager) StartRoom(r *Room) {
	go func() {
		for {
			if r.Player1.Disconnected.Load() && r.Player2.Disconnected.Load() {
				log.Println("Both players disconnected — closing room")
				rm.CloseRoom(r)
				return
			}
			
			select {
			case msg, ok := <-r.Player1.Incoming:
				if !ok { continue }
				rm.handleRequest(r, msg)
			case msg, ok := <-r.Player2.Incoming:
				if !ok { continue }
				rm.handleRequest(r, msg)
			}
		}
	}()
}


func (rm *RoomManager) CloseRoom(r *Room) {
	r.closeOnce.Do(func() {
		//TODO: cleanup room resources
	})
}

func (rm *RoomManager) handleRequest(r *Room, msg []byte) {
	var req ClientRequest
	json.Unmarshal(msg, &req)

	clientGameState := &req.Data

	var res *ServerResponse
	
	switch (req.Type) {
	case "move":
		game.UpdateGameState(r .Game, clientGameState)
		res = &ServerResponse{
			Type: "update",
			Data: r.Game,
		}
	}
	rm.Broadcast(r, res)
}



func (rm *RoomManager) CreateNewRoom(player1 *game.Player, player2 *game.Player, lobbyType string) *Room {
	var size int;
	switch lobbyType {
	case "9x9":
		size = 9
	case "15x15":
		size = 15
	case "19x19":
		size = 19
	}
	
	newRoom := &Room{
		Player1: player1,
		Player2: player2,
		Game: game.CreateGameState(size, player1, player2),
	}

	rm.AddPlayerToRoom(player1, newRoom)
	rm.AddPlayerToRoom(player2, newRoom)
	
	return newRoom
}

func (rm *RoomManager) AddPlayerToRoom(player *game.Player, room *Room) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.playerRoomMap[player.PlayerID] = room
}

func (rm *RoomManager) RemovePlayerFromRoom(player *game.Player) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.playerRoomMap, player.PlayerID)
}

func (rm *RoomManager) GetRoom(playerID string) (*Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, ok := rm.playerRoomMap[playerID]
	
	return room, ok
}

func (rm *RoomManager) ReconnectPlayer(playerID string, newPlayer *game.Player) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.playerRoomMap[playerID]
	if !ok {
		return false
	}

	if room.Player1.PlayerID == playerID {
		room.Player1.Conn = newPlayer.Conn
		room.Player1.Incoming = newPlayer.Incoming
		room.Player1.Outgoing = newPlayer.Outgoing
		room.Player1.Disconnected.Store(false)
		return true
	} else if room.Player2.PlayerID == playerID {
		room.Player2.Conn = newPlayer.Conn
		room.Player2.Incoming = newPlayer.Incoming
		room.Player2.Outgoing = newPlayer.Outgoing
		room.Player2.Disconnected.Store(false)
		return true
	}

	return false
}






