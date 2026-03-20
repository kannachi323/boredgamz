package connectfour

import (
	"boredgamz/core"
	"boredgamz/core/connectfour/model"
	"boredgamz/db"
	cfdb "boredgamz/db/connectfour"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectFourGameEvent struct {
	PlayerID string
	Data     []byte
}

type ConnectFourTimeoutEvent struct {
	PlayerID string
}

type ConnectFourRoom struct {
	*core.Room
	GameState   *ConnectFourGameState
	RoomManager *core.RoomManager
	finishOnce  sync.Once
	mu          sync.Mutex
}

func NewConnectFourRoom(p1, p2 *core.Player, gameType string, db *db.Database, roomManager *core.RoomManager) *ConnectFourRoom {
	return &ConnectFourRoom{
		Room: &core.Room{
			DB:      db,
			RoomID:  uuid.New().String(),
			Players: []*core.Player{p1, p2},
			Events:  make(chan interface{}, 100),
			Done:    make(chan struct{}),
		},
		GameState:   NewConnectFourGame(gameType, p1, p2),
		RoomManager: roomManager,
	}
}

func (room *ConnectFourRoom) Start() {
	go room.watchIncoming()
	go room.watchDisconnections()
	go room.runClock()
	go room.eventLoop()

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
	room.CloseOnce.Do(func() {
		close(room.Done)
		for _, player := range room.Players {
			if room.RoomManager != nil {
				room.RoomManager.RemovePlayerFromRoom(player.PlayerID)
			}
			player.ClosePlayer()
		}
		log.Println("Connect Four room closed:", room.RoomID)
	})
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
	case <-room.Done:
		return
	case room.Events <- raw:
	default:
	}
}

func (room *ConnectFourRoom) eventLoop() {
	for ev := range room.Events {
		select {
		case <-room.Done:
			return
		default:
		}

		switch e := ev.(type) {
		case ConnectFourTimeoutEvent:
			room.handleTimeout(e)
		case ConnectFourGameEvent:
			room.handleConnectFourEvent(e)
		}

		if room.GameState.Status.Code == "offline" {
			room.handleGameFinished()
			return
		}
	}
}

func (room *ConnectFourRoom) watchIncoming() {
	if len(room.Players) != 2 {
		return
	}
	p1 := room.Players[0]
	p2 := room.Players[1]

	for {
		select {
		case <-room.Done:
			return
		default:
		}

		if p1.Disconnected.Load() && p2.Disconnected.Load() {
			return
		}

		select {
		case <-room.Done:
			return
		case raw, ok := <-p1.Incoming:
			if ok {
				room.HandleEvent(ConnectFourGameEvent{PlayerID: p1.PlayerID, Data: raw})
			}
		case raw, ok := <-p2.Incoming:
			if ok {
				room.HandleEvent(ConnectFourGameEvent{PlayerID: p2.PlayerID, Data: raw})
			}
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

	for {
		select {
		case <-room.Done:
			return
		case <-ticker.C:
		}

		if room.GameState.Status.Code == "offline" {
			return
		}

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

func (room *ConnectFourRoom) handleClientRequest(playerID string, raw []byte) {
	var req ConnectFourClientRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return
	}

	switch req.Type {
	case "move":
		if room.GameState.Turn != playerID {
			return
		}
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
	case "chat":
		var chatData ConnectFourChatData
		if json.Unmarshal(req.Data, &chatData) != nil {
			return
		}

		message := room.buildChatMessage(playerID, &chatData)
		if message == nil {
			return
		}

		room.GameState.Messages = append(room.GameState.Messages, message)

		data, _ := json.Marshal(message)
		resp := &ConnectFourServerResponse{
			Type: "chat",
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
	room.handleClientRequest(ev.PlayerID, ev.Data)
}

func (room *ConnectFourRoom) buildChatMessage(playerID string, chatData *ConnectFourChatData) *ChatMessage {
	if chatData == nil {
		return nil
	}

	content := strings.TrimSpace(chatData.Content)
	if content == "" || len(content) > 500 {
		return nil
	}

	sender := GetPlayerByID(room.GameState, playerID)
	if sender == nil {
		return nil
	}

	return &ChatMessage{
		SenderID:   sender.PlayerID,
		SenderName: sender.PlayerName,
		Content:    content,
		SentAt:     time.Now().UTC().Format(time.RFC3339),
	}
}

func (room *ConnectFourRoom) handleGameFinished() {
	room.finishOnce.Do(func() {
		if room.DB == nil || room.DB.Pool == nil {
			log.Println("Skipping connect four persistence: database is not initialized")
		} else {
			gameRow := room.GameState.ToRow()
			go func() {
				err := cfdb.SaveGameWithMessages(
					room.DB,
					room.GameState.GameID,
					room.GameState.Players[0].PlayerID,
					room.GameState.Players[1].PlayerID,
					gameRow,
					gameRow.Messages,
				)
				if err != nil {
					log.Println("Error saving finished Connect Four game to database:", err)
				} else {
					log.Println("Finished Connect Four game saved:", room.GameState.GameID)
				}
			}()
		}

		data, _ := json.Marshal(room.GameState)
		res := &ConnectFourServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(res)
		room.Broadcast(resBytes)

		go func() {
			time.Sleep(2 * time.Second)
			room.Close()
		}()
	})
}

func (room *ConnectFourRoom) runClock() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-room.Done:
			return
		case <-ticker.C:
		}

		if room.GameState.Status.Code == "offline" {
			return
		}

		room.tickActivePlayerClock()
	}
}

func (room *ConnectFourRoom) tickActivePlayerClock() {
	if room.GameState.Status.Code != "online" {
		return
	}

	currPlayer := GetPlayerByID(room.GameState, room.GameState.Turn)
	if currPlayer == nil || currPlayer.Clock == nil {
		return
	}

	if currPlayer.Clock.Remaining <= 0 {
		currPlayer.Clock.Remaining = 0
		room.HandleEvent(ConnectFourTimeoutEvent{PlayerID: currPlayer.PlayerID})
	} else {
		currPlayer.Clock.Remaining -= time.Second
	}
}

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

	messages := make([]*model.ChatMessage, len(gs.Messages))
	for i, m := range gs.Messages {
		messages[i] = &model.ChatMessage{
			SenderID:   m.SenderID,
			SenderName: m.SenderName,
			Content:    m.Content,
			SentAt:     m.SentAt,
		}
	}

	return &model.ConnectFourGameStateRow{
		GameID:   gs.GameID,
		Players:  players,
		Moves:    moves,
		Messages: messages,
		Result:   gs.Status.Result,
		Winner:   winner,
	}
}

func (room *ConnectFourRoom) ReconnectPlayer(playerID string, conn *websocket.Conn) error {
	room.mu.Lock()
	defer room.mu.Unlock()

	var player *core.Player
	for _, p := range room.Players {
		if p.PlayerID == playerID {
			player = p
			break
		}
	}

	if player == nil {
		return fmt.Errorf("player not found in room")
	}

	player.ReconnectPlayer(conn)

	resData, err := json.Marshal(room.GameState)
	if err != nil {
		return err
	}

	res := &ConnectFourServerResponse{
		Type: "update",
		Data: resData,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	room.Send(player, resBytes)
	return nil
}
