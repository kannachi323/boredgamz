package manager

import (
	"playgomoku/backend/game"
	"sync"
)

type RoomController interface {
	Start()
	Close()
	GetRoomID() string
	GetPlayerByID() string
	GetGameState() *game.GameState
	GetPlayerChannels(playerID string) (incoming chan []byte, ouutgoing chan []byte)
	ReconnectPlayer(newPlayer *game.Player) bool
	RemovePlayer(playerID string) bool
}

type Room struct {
	RoomController

	RoomID 	  string
	Player1   *game.Player
    Player2   *game.Player
    Events    chan *ClientRequest
    Timeout   chan string
    GameID    string
    GameState *game.GameState
	mu 	  sync.RWMutex
}

type RoomManager struct { 
	PlayerRoomMap map[string]*Room
}

func (rm *RoomManager) StartRoom(r *Room) {
	r.Start() 
}

func (rm *RoomManager) CloseRoom(r Room) {
	r.Close()
}