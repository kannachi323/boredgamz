package connectfour

import (
	"boredgamz/core"
	"boredgamz/core/connectfour/model"
	"boredgamz/db"
	cfdb "boredgamz/db/connectfour"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type ConnectFourGameEvent struct {
	Type     string
	PlayerID string
	Data     []byte
}

type ConnectFourTimeoutEvent struct {
	PlayerID string
}

type ConnectFourRoom struct {
	*core.Room
	GameState *ConnectFourGameState
}

func NewConnectFourRoom(p1, p2 *core.Player, gameType string, db *db.Database) core.RoomController {
	newRoom := &ConnectFourRoom{
		Room: &core.Room{
			DB:      db,
			RoomID:  uuid.New().String(),
			Players: []*core.Player{p1, p2},
			Events:  make(chan interface{}, 100),
		},
		GameState: NewConnectFourGame(gameType, p1, p2),
	}

	return newRoom
}

////////////////////////////
// ROOM LIFECYCLE METHODS
///////////////////////////
func (room *ConnectFourRoom) Start() {
	go room.watchIncoming()
	go room.watchDisconnections()
	go room.runClock()
	go room.eventLoop()

	// Broadcast initial game state
	resData, _ := json.Marshal(room.GameState)
	res := &ConnectFourServerResponse{
		Type: "update",
		Data: resData,
	}
	resBytes, err := json.Marshal(res)
	if err == nil {
		log.Println("Broadcasting initial Connect Four game state for room:", room.RoomID)
		room.Broadcast(resBytes)
	}
}

func (room *ConnectFourRoom) Close() {
	for _, player := range room.Players {
		player.ClosePlayer()
	}
	close(room.Events)
	log.Println("Connect Four room closed:", room.RoomID)
}

func (room *ConnectFourRoom) Broadcast(res []byte) {
	for _, player := range room.Players {
		if player.Disconnected.Load() {
			continue
		}
		room.Send(player, res)
	}
}

func (room *ConnectFourRoom) Send(p *core.Player, res []byte) {
	if p.Disconnected.Load() {
		return
	}

	select {
	case p.Outgoing <- res:
	default:
	}
}

func (room *ConnectFourRoom) HandleEvent(raw interface{}) {
	select {
	case room.Events <- raw:
	default:
	}
}

////////////////////////////
// EVENT LOOP
///////////////////////////
func (room *ConnectFourRoom) eventLoop() {
	for ev := range room.Events {
		switch e := ev.(type) {
		case []byte:
			room.handleClientRequest(e)
		case ConnectFourTimeoutEvent:
			room.handleTimeout(e)
		case ConnectFourGameEvent:
			room.handleConnectFourEvent(e)
		default:
			// unknown event → ignore
		}

		if room.GameState.Status.Code == "offline" {
			room.handleGameFinished()
			return
		}
	}
}

////////////////////////////
// WATCHERS
///////////////////////////
func (room *ConnectFourRoom) watchIncoming() {
	if len(room.Players) != 2 {
		return
	}
	p1 := room.Players[0]
	p2 := room.Players[1]

	for {
		if p1.Disconnected.Load() && p2.Disconnected.Load() {
			return
		}
		select {
		case raw, ok := <-p1.Incoming:
			if ok {
				room.HandleEvent(raw)
			}
		case raw, ok := <-p2.Incoming:
			if ok {
				room.HandleEvent(raw)
			}
		default:
			// no messages
		}
	}
}

func (room *ConnectFourRoom) watchDisconnections() {
	if len(room.Players) != 2 {
		return
	}
	p1 := room.Players[0]
	p2 := room.Players[1]

	var (
		p1Disc time.Duration
		p2Disc time.Duration
	)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if p1.Disconnected.Load() {
			p1Disc += time.Second
			if p1Disc >= 30*time.Second {
				room.HandleEvent(ConnectFourTimeoutEvent{PlayerID: p1.PlayerID})
				return
			}
		} else {
			p1Disc = 0
		}

		if p2.Disconnected.Load() {
			p2Disc += time.Second
			if p2Disc >= 30*time.Second {
				room.HandleEvent(ConnectFourTimeoutEvent{PlayerID: p2.PlayerID})
				return
			}
		} else {
			p2Disc = 0
		}
	}
}

////////////////////////////
// HANDLERS
///////////////////////////
func (room *ConnectFourRoom) handleClientRequest(raw []byte) {
	var req ConnectFourClientRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return
	}

	switch req.Type {
	case "move":
		var moveData ConnectFourMoveData
		if json.Unmarshal(req.Data, &moveData) != nil {
			return
		}

		HandleConnectFourMove(room.GameState, moveData.Column)

		data, _ := json.Marshal(room.GameState)
		resp := &ConnectFourServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(resp)
		room.Broadcast(resBytes)
	}
}

func (room *ConnectFourRoom) handleTimeout(ev ConnectFourTimeoutEvent) {
	UpdateGameStatus(room.GameState, "timeout", ev.PlayerID)
	data, _ := json.Marshal(room.GameState)
	res := &ConnectFourServerResponse{
		Type: "update",
		Data: data,
	}
	resBytes, _ := json.Marshal(res)
	room.Broadcast(resBytes)
}

func (room *ConnectFourRoom) handleConnectFourEvent(ev ConnectFourGameEvent) {
	return
}

func (room *ConnectFourRoom) handleGameFinished() {
	room.CloseOnce.Do(func() {
		go func() {
			err := cfdb.InsertGame(
				room.DB,
				room.GameState.GameID,
				room.GameState.Players[0].PlayerID,
				room.GameState.Players[1].PlayerID,
				room.GameState.ToRow(),
			)
			if err != nil {
				log.Println("Error saving finished Connect Four game to database:", err)
			} else {
				log.Println("Finished Connect Four game saved:", room.GameState.GameID)
			}
		}()

		data, _ := json.Marshal(room.GameState)
		res := &ConnectFourServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(res)
		room.Broadcast(resBytes)
	})
}

////////////////////////////
// CLOCK MANAGEMENT
///////////////////////////
func (room *ConnectFourRoom) runClock() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		room.tickActivePlayerClock()
	}
}

func (room *ConnectFourRoom) tickActivePlayerClock() {
	if room.GameState.Status.Code != "online" {
		return
	}

	currPlayer := GetPlayerByID(room.GameState, room.GameState.Turn)

	if currPlayer.Clock.Remaining <= 0 {
		currPlayer.Clock.Remaining = 0
		room.HandleEvent(ConnectFourTimeoutEvent{PlayerID: currPlayer.PlayerID})
	} else {
		currPlayer.Clock.Remaining -= time.Second
	}
}

////////////////////////////
// DATABASE PERSISTENCE
///////////////////////////
func (gs *ConnectFourGameState) ToRow() *model.ConnectFourGameStateRow {
	moves := make([]*model.Move, len(gs.Moves))
	for i, m := range gs.Moves {
		moves[i] = &model.Move{
			Column: m.Col,
			Row:    m.Row,
			Color:  m.Color,
		}
	}

	var winner *model.Player
	if gs.Status.Winner != nil {
		winner = &model.Player{
			PlayerID:   gs.Status.Winner.PlayerID,
			PlayerName: gs.Status.Winner.PlayerName,
			Color:      gs.Status.Winner.Color,
		}
	}

	var players []*model.Player
	for _, p := range gs.Players {
		players = append(players, &model.Player{
			PlayerID:   p.PlayerID,
			PlayerName: p.PlayerName,
			Color:      p.Color,
		})
	}

	return &model.ConnectFourGameStateRow{
		GameID: gs.GameID,
		Players: players,
		Moves:  moves,
		Result: gs.Status.Result,
		Winner: winner,
	}
}
