package gomoku

import (
	"encoding/json"
	"log"
	"playgomoku/backend/game"
	"playgomoku/backend/manager"
	"sync"
	"sync/atomic"
	"time"
)


type GomokuRoom struct {
	manager.Room
}

type GomokuRoomManager struct {
	playerRoomMap map[string]*GomokuRoom
	mu	sync.RWMutex
}

func CreateRoomManager() *GomokuRoomManager {
	newRoomManager := &GomokuRoomManager{
		playerRoomMap: make(map[string]*GomokuRoom),
	}

	return newRoomManager;
}

func (rm *GomokuRoomManager) Broadcast(r *GomokuRoom, res *manager.ServerResponse) error {

	Send(r.Player1, res)
	Send(r.Player2, res)

	return nil
}

func Send(p *game.Player, res *manager.ServerResponse) {
	msg, err := json.Marshal(res)
	if err != nil {
		log.Println("unable to send messages")
		return
	}
	
	if p.Disconnected.Load() { return }

	select {
	case p.Outgoing <- msg:
	default:
	}
}

func (rm *GomokuRoomManager) StartRoom(r *GomokuRoom) {
	r.Player1.StartPlayer()
	r.Player2.StartPlayer()


	rm.StartPlayersListener(r)
	rm.StartEventsListener(r)
	rm.StartTimeoutListener(r)
	rm.StartConnectionListener(r)
}

func (rm *GomokuRoomManager) CloseRoom(r *GomokuRoom) {
	r.closeOnce.Do(func() {
		r.Player1.ClosePlayer()
		r.Player2.ClosePlayer()
	})
}

func (rm *GomokuRoomManager) StartPlayersListener(r *GomokuRoom) {
	go func() {
		for {
			if r.Player1.Disconnected.Load() && r.Player2.Disconnected.Load() { return }
			select {
			case msg, ok := <-r.Player1.Incoming:
				if !ok { continue }
				rm.handleRequest(r, msg)
			case msg, ok := <-r.Player2.Incoming:
				if !ok { continue }
				rm.handleRequest(r, msg)
			default:
				//no incoming messages
			}
		}
	}()
}

func (rm *GomokuRoomManager) StartEventsListener(r *GomokuRoom) {
	go func() {
		for req := range r.Events {
			log.Printf("Room %s received event: %v\n", r.GameState.GameID, req)
			switch (req.Type) {
			case "move":
				game.UpdateGameState(r.GameState, req.Data)
				var res *manager.ServerResponse
				res = &manager.ServerResponse{
					Type: "update",
					Data: r.GameState,
				}
				rm.Broadcast(r, res)
			}
		}
	}()
}

func (rm *GomokuRoomManager) StartTimeoutListener(r *GomokuRoom) {
	go func() {
		for playerID := range r.Timeout {
			if (r.GameState.Status.Code == "offline") { return }
			var res *manager.ServerResponse
			game.UpdateGameStatus(r.GameState, "timeout", playerID)
			res = &manager.ServerResponse{
				Type: "update",
				Data: r.GameState,
			}
			rm.Broadcast(r, res)
		}
	}()
}

func (rm *GomokuRoomManager) StartConnectionListener(r *GomokuRoom) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		var p1Time, p2Time time.Duration
		const maxTime = 10 * time.Second

		for range ticker.C {
			if (r.GameState.Status.Code == "offline") { return }
			if r.Player1.Disconnected.Load() {
				p1Time += 2 * time.Second
				log.Println("Player 1 disconnected for ", p1Time)
				if p1Time >= maxTime {
					select {
					case r.Timeout <- r.Player1.PlayerID:
					default:
					}
					return
				}
			} else {
				p1Time = 0
			}
			if r.Player2.Disconnected.Load() {
				p2Time += 2 * time.Second
				log.Println("Player 2 disconnected for ", p2Time)
				if p2Time >= maxTime {
					select {
					case r.Timeout <- r.Player2.PlayerID:
					default:
					}
					return
				}
			} else {
				p2Time = 0
			}
		}
	}()
}

func (rm *GomokuRoomManager) handleRequest(r *GomokuRoom, msg []byte) {
	var req *manager.ClientRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println("Invalid client message:", err)
		return
	}

	select {
	case r.Events <- req:
	default:
		log.Printf("Room %s event queue full — dropping message\n", r.GameState.GameID)
	}
}

func (rm *GomokuRoomManager) CreateNewRoom(player1 *game.Player, player2 *game.Player, lobbyType string) *GomokuRoom {
	var size int;
	switch lobbyType {
	case "9x9":
		size = 9
	case "15x15":
		size = 15
	case "19x19":
		size = 19
	}
	newRoom := &GomokuRoom{
		Player1: player1,
		Player2: player2,
		GameState: game.CreateGameState(size, player1, player2),
		Events: make(chan *manager.ClientRequest, 50),
		Timeout: make(chan string),
	}

	//IMPORTANT: Link player timeout to room timeout channel
	player1.Clock = &game.PlayerClock{
		Remaining: player1.Clock.Remaining * time.Nanosecond,
		IsActive: atomic.Bool{},
		Timeout: newRoom.Timeout,
	}
	player2.Clock = &game.PlayerClock{
		Remaining: player2.Clock.Remaining * time.Nanosecond,
		IsActive: atomic.Bool{},
		Timeout: newRoom.Timeout,
	}
		
	rm.AddPlayerToRoom(player1, newRoom)
	rm.AddPlayerToRoom(player2, newRoom)
	
	return newRoom
}

func (rm *GomokuRoomManager) AddPlayerToRoom(player *game.Player, room *GomokuRoom) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.playerRoomMap[player.PlayerID] = room
}

func (rm *GomokuRoomManager) RemovePlayerFromRoom(player *game.Player) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.playerRoomMap, player.PlayerID)
}

func (rm *GomokuRoomManager) GetRoom(playerID string) (*GomokuRoom, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, ok := rm.playerRoomMap[playerID]
	
	return room, ok
}

func (rm *GomokuRoomManager) ReconnectPlayer(playerID string, newPlayer *game.Player) bool {
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