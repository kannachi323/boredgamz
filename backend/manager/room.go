package manager

import "sync"

type RoomController interface {
	//Room lifecycle methods
	Start()
	Close()
	Broadcast(res []byte)
	Send(p *Player, res []byte)
	HandleEvent(req []byte)
}

type Room struct {
	RoomID 	  string
	Players	 	[]*Player
	Events    chan []byte
	Timeout   chan []byte
	GameID    string
	CloseOnce sync.Once
}

type RoomManager struct { 
	PlayerRoomMap map[string]RoomController
	mu	sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		PlayerRoomMap: make(map[string]RoomController),
	}
}

func (rm *RoomManager) RegisterPlayerToRoom(playerID string, room RoomController) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.PlayerRoomMap[playerID] = room
}

func (rm *RoomManager) RemovePlayerFromRoom(playerID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.PlayerRoomMap, playerID)
}

func (rm *RoomManager) GetPlayerRoom(playerID string) (RoomController, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, ok := rm.PlayerRoomMap[playerID]
	return room, ok
}