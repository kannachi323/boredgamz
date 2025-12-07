package gomoku

import (
	"boredgamz/core"
	"boredgamz/core/gomoku/model"
	"boredgamz/db"
	gomokudb "boredgamz/db/gomoku"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type GomokuGameEvent struct {
	Type string
	PlayerID string
	Data []byte
}

type GomokuTimeoutEvent struct {
	PlayerID string
}

type GomokuRoom struct {
	*core.Room
	GameState *GomokuGameState
}

func NewGomokuRoom(p1, p2 *core.Player, gomokuType string, db *db.Database) core.RoomController{
	newGomokuRoom := &GomokuRoom{
		Room: &core.Room{
			DB: db,
			RoomID: uuid.New().String(),
			Players: []*core.Player{p1, p2},
			Events: make(chan interface{}, 100),
		},
		GameState: NewGomokuGame(gomokuType, p1, p2),
	}

	return newGomokuRoom
}


////////////////////////////
//ROOM LIFECYCLE METHODS
///////////////////////////
func (room *GomokuRoom) Start() {
	go room.watchIncoming()
	go room.watchDisconnections()
	go room.runClock()
	go room.eventLoop()

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
	
	for _, player := range room.Players {
		player.ClosePlayer()
	}
	close(room.Events)
	log.Println("Gomoku room closed:", room.RoomID)

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

func (room *GomokuRoom) HandleEvent(raw interface{}) {
	select {
	case room.Events <- raw:
	default:
	}
}

////////////////////////////
//EVENT LOOP
///////////////////////////
func (room *GomokuRoom) eventLoop() {
	for ev := range room.Events {
		switch e := ev.(type) {
		case []byte:
			room.handleClientRequest(e)
		case GomokuTimeoutEvent:
			room.handleTimeout(e)
		case GomokuGameEvent:
			room.handleGomokuEvent(e)
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
//WATCHERS
///////////////////////////
func (room *GomokuRoom) watchIncoming() {
	if len(room.Players) != 2 { return }
	p1 := room.Players[0]
	p2 := room.Players[1]

	for {
		if p1.Disconnected.Load() && p2.Disconnected.Load() { return }
		select {
		case raw, ok := <-p1.Incoming:
			if ok { room.HandleEvent(raw) }
		case raw, ok := <-p2.Incoming:
			if ok { room.HandleEvent(raw) }
		default:
			// no messages
		}
	}
}

func (room *GomokuRoom) watchDisconnections() {
	if len(room.Players) != 2 { return }
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
			if p1Disc >= 30 * time.Second {
				room.HandleEvent(GomokuTimeoutEvent{PlayerID: p1.PlayerID})
				return
			}
		} else {
			p1Disc = 0
		}

		if p2.Disconnected.Load() {
			p2Disc += time.Second
			if p2Disc >= 30 * time.Second {
				room.HandleEvent(GomokuTimeoutEvent{PlayerID: p2.PlayerID})
				return
			}
		} else {
			p2Disc = 0
		}
	}
}



////////////////////////////
//HANDLERS
///////////////////////////
func (room *GomokuRoom) handleClientRequest(raw []byte) {
	var req GomokuClientRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return
	}
	switch req.Type {
	case "move":
		var moveData GomokuMoveData
		if json.Unmarshal(req.Data, &moveData) != nil {
			return
		}
		HandleGomokuMove(room.GameState, &moveData.Move)

		data, _ := json.Marshal(room.GameState)
		resp := &GomokuServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(resp)
		room.Broadcast(resBytes)
	}
}


func (room *GomokuRoom) handleTimeout(ev GomokuTimeoutEvent) {
	UpdateGameStatus(room.GameState, "timeout", ev.PlayerID)
	data, _ := json.Marshal(room.GameState)
	res := &GomokuServerResponse{
		Type: "update",
		Data: data,
	}
	resBytes, _ := json.Marshal(res)
	room.Broadcast(resBytes)
}

func (room *GomokuRoom) handleGomokuEvent(ev GomokuGameEvent) {
	return
}

func (room *GomokuRoom) handleGameFinished() {
	room.CloseOnce.Do(func() {
		//persist the game by saving to database
		go func() {
			err := gomokudb.InsertGame(
				room.DB,
				room.GameState.GameID,
				room.GameState.Players[0].PlayerID,
				room.GameState.Players[1].PlayerID,
				room.GameState.ToRow(),
			)
			if err != nil {
				log.Println("Error saving finished gomoku game to database:", err)
			} else {
				log.Println("Finished gomoku game saved to database:", room.GameState.GameID)
			}
		}()

		// Send final game state ONE MORE TIME just to be safe
		data, _ := json.Marshal(room.GameState)
		res := &GomokuServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(res)
		room.Broadcast(resBytes)
	})
}

////////////////////////////
//CLOCK MANAGEMENT
///////////////////////////
func (room *GomokuRoom) runClock() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		room.tickActivePlayerClock()
	}
}

func (room *GomokuRoom) tickActivePlayerClock() {
	if room.GameState.Status.Code != "online" { return }

	currPlayer := GetPlayerByID(room.GameState, room.GameState.Turn)


    if currPlayer.Clock.Remaining <= 0 {
        currPlayer.Clock.Remaining = 0
        room.HandleEvent(GomokuTimeoutEvent{PlayerID: currPlayer.PlayerID})
    } else {
		currPlayer.Clock.Remaining -= time.Second
	}
}

////////////////////////////
//DATABASE MODEL PERSISTENCE
///////////////////////////
func (gs *GomokuGameState) ToRow() *model.GomokuGameStateRow {
    moves := make([]*model.Move, len(gs.Moves))
    for i, m := range gs.Moves {
        moves[i] = &model.Move{
            Row:   m.Row,
            Col:   m.Col,
            Color: m.Color,
        }
    }


	var winner *model.Player
	if gs.Status.Winner != nil {
		winner = &model.Player{
			PlayerID: gs.Status.Winner.PlayerID,
			PlayerName: gs.Status.Winner.PlayerName,
			Color: gs.Status.Winner.Color,
		}
	}

	var players []*model.Player
	for _, p := range gs.Players {
		players = append(players, &model.Player{
			PlayerID: p.PlayerID,
			PlayerName: p.PlayerName,
			Color: p.Color,
		})
	}
	
    return &model.GomokuGameStateRow{
		GameID: gs.GameID,
		BoardSize: gs.Board.Size,
		Players: players,
        Moves:  moves,
        Result: gs.Status.Result,
        Winner: winner,
    }
}
