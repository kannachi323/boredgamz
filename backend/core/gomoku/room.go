package gomoku

import (
	"boredgamz/core"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GomokuRoom struct {
	*core.Room
	GameState *GomokuGameState
}

func NewGomokuRoom(p1, p2 *core.Player, gomokuType string) *GomokuRoom {
	newGomokuRoom := &GomokuRoom{
		Room: &core.Room{
			RoomID: uuid.New().String(),
			Players: []*core.Player{p1, p2},
			Events: make(chan []byte, 100),
			Timeout: make(chan []byte, 100),
			GameID: uuid.New().String(),
		},
		GameState: NewGomokuGameState(gomokuType, p1, p2),
	}

	//IMPORTANT: Link player timeout to room timeout channel
	p1.Clock.Timeout = newGomokuRoom.Timeout
	p2.Clock.Timeout = newGomokuRoom.Timeout

	return newGomokuRoom
}

func (room *GomokuRoom) Start() {
	go room.startEventListener()
	go room.startTimeoutListener()
	go room.startPlayersListener()
	go room.startConnectionListener()

}


func (room *GomokuRoom) Close() {
	room.CloseOnce.Do(func() {
		for _, player := range room.Players {
			player.ClosePlayer()
		}
	})
}

func (room *GomokuRoom) Broadcast(res *GomokuServerResponse)  {
	for _, player := range room.Players {
		if player.Disconnected.Load() { continue }
		room.Send(player, res)
	}
}

func (room *GomokuRoom) Send(p *core.Player, res *GomokuServerResponse) {
	msg, err := json.Marshal(res)
	if err != nil { return }
	if p.Disconnected.Load() { return }

	select {
	case p.Outgoing <- msg:
	default:
	}
}

func (room *GomokuRoom) HandleEvent(raw []byte) {
	select {
	case room.Events <- raw:
	default:
	}
}



func (room *GomokuRoom) startEventListener() {
	for raw := range room.Events {
		var req GomokuClientRequest
		err := json.Unmarshal(raw, &req)
		if err != nil { continue }

		switch (req.Type) {
		case "move":
			UpdateGameState(room.GameState, req.Data)
			var res *GomokuServerResponse
			res = &GomokuServerResponse{
				Type: "update",
				Data: room.GameState,
			}
			room.Broadcast(res)
		}
	}
}

func (room *GomokuRoom) startPlayersListener() {
	if (len(room.Players) != 2) { return }
	
	p1 := room.Players[0]
	p2 := room.Players[1]

	for {
		if p1.Disconnected.Load() && p2.Disconnected.Load() { return }
		select {
		case raw, ok := <-p1.Incoming:
			if !ok { continue }
			room.HandleEvent(raw)

		case raw, ok := <-p2.Incoming:
			if !ok { continue }
			room.HandleEvent(raw)
		default:
			//no incoming messages
		}
	}
}


func (room *GomokuRoom) startConnectionListener() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var p1Time, p2Time time.Duration
	const maxTime = 10 * time.Second

	p1 := room.Players[0]
	p2 := room.Players[1]

	for range ticker.C {
		if (room.GameState.Status.Code == "offline") { return }
		if p1.Disconnected.Load() {
			p1Time += 2 * time.Second
			if p1Time >= maxTime {
				select {
				case room.Timeout <- []byte(p1.PlayerID):
				default:
				}
				return
			}
		} else {
			p1Time = 0
		}
		if p2.Disconnected.Load() {
			p2Time += 2 * time.Second
			if p2Time >= maxTime {
				select {
				case room.Timeout <- []byte(p2.PlayerID):
				default:
				}
				return
			}
		} else {
			p2Time = 0
		}
	}
}

func (room *GomokuRoom) startTimeoutListener() {
	for playerID := range room.Timeout {
		if (room.GameState.Status.Code == "offline") { return }
		var res *GomokuServerResponse
		UpdateGameStatus(room.GameState, "timeout", string(playerID))
		res = &GomokuServerResponse{
			Type: "update",
			Data: room.GameState,
		}
		room.Broadcast(res)
	}
}