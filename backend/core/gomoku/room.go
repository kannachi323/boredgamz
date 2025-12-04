package gomoku

import (
	"boredgamz/core"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type GomokuRoom struct {
	*core.Room
	GameState *GomokuGameState
}

func NewGomokuRoom(p1, p2 *core.Player, gomokuType string) core.RoomController{
	newGomokuRoom := &GomokuRoom{
		Room: &core.Room{
			RoomID: uuid.New().String(),
			Players: []*core.Player{p1, p2},
			Events: make(chan []byte, 100),
			Timeout: make(chan []byte, 100),
			GameID: uuid.New().String(),
		},
		GameState: NewGomokuGame(gomokuType, p1, p2),
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

	//make initial broadcast of game state
	resData, _ := json.Marshal(room.GameState)
	res := &GomokuServerResponse{
		Type: "update",
		Data: resData,
	}
	resBytes, err := json.Marshal(res)
	if err == nil {
		log.Println("Broadcasting initial game state for room:", room.RoomID)
		room.Broadcast(resBytes)
	}
}


func (room *GomokuRoom) Close() {
	room.CloseOnce.Do(func() {
		for _, player := range room.Players {
			player.ClosePlayer()
		}
	})
}

func (room *GomokuRoom) Broadcast(res []byte)  {
	for _, player := range room.Players {
		if player.Disconnected.Load() { continue }
		room.Send(player, res)
	}
}

func (room *GomokuRoom) Send(p *core.Player, res []byte) {
	if p.Disconnected.Load() { return }

	select {
	case p.Outgoing <- res:
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
			var moveData GomokuMoveData
			err := json.Unmarshal(req.Data, &moveData)
			log.Println("RAW DATA:", string(req.Data))

			log.Println(moveData)
			if err != nil { continue }
			HandleGomokuMove(room.GameState, moveData.Move.Row, moveData.Move.Col, moveData.Move.Color)
			
			resData, _ := json.Marshal(room.GameState)
			var res *GomokuServerResponse
			res = &GomokuServerResponse{
				Type: "update",
				Data: resData,
			}
			resBytes, err := json.Marshal(res)
			if err == nil {
				room.Broadcast(resBytes)
			}
			
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
		data, _ := json.Marshal(room.GameState)
		res = &GomokuServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, err := json.Marshal(res)
		if err == nil {
			room.Broadcast(resBytes)
		}
	}
}