package manager

import (
	"encoding/json"
	"errors"
	"log"
	"playgomoku/backend/game"
	"sync"
)

type Room struct {
	Player1 *game.Player
	Player2 *game.Player
	Game   *game.GameState
	Quit	chan struct{}
}

type RoomManager struct {
	playerRoomMap map[string]*Room
	mu	sync.RWMutex
}

type ServerMessage struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

func CreateRoomManager() *RoomManager {
	newRoomManager := &RoomManager{
		playerRoomMap: make(map[string]*Room),
	}

	return newRoomManager;
}

func (r *Room) Broadcast(msg ServerMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	send := func(p *game.Player) error {
		select {
		case p.Outgoing <- b:
			return nil
		default:
			return errors.New("player outgoing channel full or blocked")
		}
	}

	if err1 := send(r.Player1); err1 != nil {
		return err1
	}
	if err2 := send(r.Player2); err2 != nil {
		return err2
	}

	return nil
}

func (r *Room) Start() {
	for {
		select {
		case msg, ok := <-r.Player1.Incoming:
			if !ok {
				r.Close()
				return
			}
			r.handleMessage(r.Player1, msg)
		
		case msg, ok := <-r.Player2.Incoming:
			if !ok {
				r.Close()
				return
			}
			r.handleMessage(r.Player2, msg)
		}
	}
}

func (r *Room) Close() {
	close(r.Quit)
	//TODO: add other clean up methods
}

func (r *Room) handleMessage(player *game.Player, msg []byte) {
	log.Printf("received message from player %s: %s", player.ID, string(msg))

	//TODO: update game state, broadcast updated state to both players
	newMsg := &ServerMessage{
		
	}
	r.Broadcast()
}


func (rm *RoomManager) CreateNewRoom(player1 *game.Player, player2 *game.Player, lobbyType string) *Room {
	var size int;
	switch lobbyType {
	case "9x9":
		size = 8
	case "15x15":
		size = 15
	case "19x19":
		size = 19
	}
	
	newRoom := &Room{
		Player1: player1,
		Player2: player2,
		Game: game.CreateGameState(size, player1, player2),
		Quit: make(chan struct{}),
	}

	//add the players to the playerRoomMap here
	rm.AddPlayerToRoom(player1, newRoom)
	rm.AddPlayerToRoom(player2, newRoom)

	go newRoom.Start()
	
	return newRoom
}

func (rm *RoomManager) AddPlayerToRoom(player *game.Player, room *Room) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	rm.playerRoomMap[player.ID] = room
}

func (rm* RoomManager) RemovePlayerFromRoom(player *game.Player) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.playerRoomMap, player.ID)
}

func (rm *RoomManager) GetRoom(playerID string) (*Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, ok := rm.playerRoomMap[playerID]
	
	return room, ok
}






