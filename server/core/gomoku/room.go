package gomoku

import (
	"boredgamz/core"
	"boredgamz/core/gomoku/model"
	"boredgamz/db"
	gomokudb "boredgamz/db/gomoku"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GomokuGameEvent struct {
	Type     string
	PlayerID string
	Data     []byte
}

type GomokuTimeoutEvent struct {
	PlayerID string
}

type GomokuRoom struct {
	*core.Room
	GameState     *GomokuGameState
	RoomManager   *core.RoomManager
	BotPlayerID   string
	BotDifficulty string
	finishOnce    sync.Once
	mu            sync.Mutex
}

func NewGomokuRoom(
	p1,
	p2 *core.Player,
	name string,
	openingRule string,
	swapRuleEnabled bool,
	firstMoveCenterEnabled bool,
	db *db.Database,
	roomManager *core.RoomManager,
) *GomokuRoom {
	newGomokuRoom := &GomokuRoom{
		Room: &core.Room{
			DB:      db,
			RoomID:  uuid.New().String(),
			Players: []*core.Player{p1, p2},
			Events:  make(chan interface{}, 100),
			Done:    make(chan struct{}),
		},
		GameState:   NewGomokuGame(name, p1, p2, openingRule, swapRuleEnabled, firstMoveCenterEnabled),
		RoomManager: roomManager,
	}

	return newGomokuRoom
}

// //////////////////////////
// ROOM LIFECYCLE METHODS
// /////////////////////////
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

	room.scheduleBotTurn()
}

func (room *GomokuRoom) Close() {
	room.CloseOnce.Do(func() {
		close(room.Done)
		for _, player := range room.Players {
			if room.RoomManager != nil {
				room.RoomManager.RemovePlayerFromRoom(player.PlayerID)
			}
			player.ClosePlayer()
		}
		log.Println("Gomoku room closed:", room.RoomID)
	})
}

func (room *GomokuRoom) Broadcast(res []byte) {
	for _, player := range room.Players {
		if player.Disconnected.Load() {
			continue
		}
		room.Send(player, res)
	}
}

func (room *GomokuRoom) Send(p *core.Player, res []byte) {
	if p.Disconnected.Load() {
		log.Println("player is still disconnected")
		return
	}

	select {
	case p.Outgoing <- res:
	default:
	}
}

func (room *GomokuRoom) HandleEvent(raw interface{}) {
	select {
	case <-room.Done:
		return
	case room.Events <- raw:
	default:
	}
}

// //////////////////////////
// EVENT LOOP
// /////////////////////////
func (room *GomokuRoom) eventLoop() {
	for ev := range room.Events {
		select {
		case <-room.Done:
			return
		default:
		}

		switch e := ev.(type) {
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

// //////////////////////////
// WATCHERS
// /////////////////////////
func (room *GomokuRoom) watchIncoming() {
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
				room.HandleEvent(GomokuGameEvent{PlayerID: p1.PlayerID, Data: raw})
			}
		case raw, ok := <-p2.Incoming:
			if ok {
				room.HandleEvent(GomokuGameEvent{PlayerID: p2.PlayerID, Data: raw})
			}
		}
	}
}

func (room *GomokuRoom) watchDisconnections() {
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
				room.HandleEvent(GomokuTimeoutEvent{PlayerID: p1.PlayerID})
				return
			}
		} else {
			p1Disc = 0
		}

		if p2.Disconnected.Load() {
			p2Disc += time.Second
			if p2Disc >= 30*time.Second {
				room.HandleEvent(GomokuTimeoutEvent{PlayerID: p2.PlayerID})
				return
			}
		} else {
			p2Disc = 0
		}
	}
}

// //////////////////////////
// HANDLERS
// /////////////////////////
func (room *GomokuRoom) handleClientRequest(playerID string, raw []byte) {
	var req GomokuClientRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return
	}
	switch req.Type {
	case "move":
		if room.GameState.Turn != playerID {
			return
		}
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
		room.scheduleBotTurn()
	case "chat":
		var chatData GomokuChatData
		if json.Unmarshal(req.Data, &chatData) != nil {
			return
		}

		message := room.buildChatMessage(playerID, &chatData)
		if message == nil {
			return
		}

		room.GameState.Messages = append(room.GameState.Messages, message)

		data, _ := json.Marshal(message)
		resp := &GomokuServerResponse{
			Type: "chat",
			Data: data,
		}
		resBytes, _ := json.Marshal(resp)
		room.Broadcast(resBytes)
	case "swap":
		if room.GameState.Turn != playerID {
			return
		}
		if err := HandleGomokuSwap(room.GameState); err != nil {
			return
		}

		data, _ := json.Marshal(room.GameState)
		resp := &GomokuServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(resp)
		room.Broadcast(resBytes)
		room.scheduleBotTurn()
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
	room.handleClientRequest(ev.PlayerID, ev.Data)
}

func (room *GomokuRoom) buildChatMessage(playerID string, chatData *GomokuChatData) *ChatMessage {
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

func (room *GomokuRoom) scheduleBotTurn() {
	if room.BotPlayerID == "" || room.GameState.Status.Code != "online" {
		return
	}
	if room.GameState.Turn != room.BotPlayerID {
		return
	}

	go func() {
		time.Sleep(450 * time.Millisecond)

		room.mu.Lock()
		defer room.mu.Unlock()

		if room.GameState.Status.Code != "online" || room.GameState.Turn != room.BotPlayerID {
			return
		}

		botPlayer := GetPlayerByID(room.GameState, room.BotPlayerID)
		if botPlayer == nil {
			return
		}

		move := room.pickBotMove(botPlayer.Color)
		if move == nil {
			return
		}

		HandleGomokuMove(room.GameState, move)

		data, _ := json.Marshal(room.GameState)
		resp := &GomokuServerResponse{
			Type: "update",
			Data: data,
		}
		resBytes, _ := json.Marshal(resp)
		room.Broadcast(resBytes)
	}()
}

func (room *GomokuRoom) pickBotMove(botColor string) *Move {
	size := room.GameState.Board.Size
	empty := room.emptyMoves()
	if len(empty) == 0 {
		return nil
	}

	if room.BotDifficulty == "advanced" {
		if winning := room.findWinningMove(botColor); winning != nil {
			return winning
		}
		opponent := "white"
		if botColor == "white" {
			opponent = "black"
		}
		if block := room.findWinningMove(opponent); block != nil {
			return &Move{Row: block.Row, Col: block.Col, Color: botColor}
		}
	}

	if room.BotDifficulty == "intermediate" || room.BotDifficulty == "advanced" {
		if room.GameState.LastMove != nil {
			near := room.neighboringMoves(room.GameState.LastMove.Row, room.GameState.LastMove.Col, botColor)
			if len(near) > 0 {
				return near[rand.Intn(len(near))]
			}
		}

		center := size / 2
		centerMove := &Move{Row: center, Col: center, Color: botColor}
		if IsValidMove(room.GameState.Board, centerMove) {
			return centerMove
		}
	}

	r := empty[rand.Intn(len(empty))]
	return &Move{Row: r.Row, Col: r.Col, Color: botColor}
}

func (room *GomokuRoom) emptyMoves() []*Move {
	var moves []*Move
	for r := 0; r < room.GameState.Board.Size; r++ {
		for c := 0; c < room.GameState.Board.Size; c++ {
			if room.GameState.Board.Stones[r][c].Color == "" {
				moves = append(moves, &Move{Row: r, Col: c})
			}
		}
	}
	return moves
}

func (room *GomokuRoom) neighboringMoves(row, col int, color string) []*Move {
	var result []*Move
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}

			r := row + dr
			c := col + dc
			move := &Move{Row: r, Col: c, Color: color}
			if IsValidMove(room.GameState.Board, move) {
				result = append(result, move)
			}
		}
	}
	return result
}

func (room *GomokuRoom) findWinningMove(color string) *Move {
	for r := 0; r < room.GameState.Board.Size; r++ {
		for c := 0; c < room.GameState.Board.Size; c++ {
			if room.GameState.Board.Stones[r][c].Color != "" {
				continue
			}

			move := &Move{Row: r, Col: c, Color: color}
			room.GameState.Board.Stones[r][c].Color = color
			isWin := IsGomoku(room.GameState.Board.Stones, move)
			room.GameState.Board.Stones[r][c].Color = ""

			if isWin {
				return move
			}
		}
	}

	return nil
}

func (room *GomokuRoom) handleGameFinished() {
	room.finishOnce.Do(func() {
		if room.DB == nil || room.DB.Pool == nil {
			log.Println("Skipping gomoku game persistence: database is not initialized")
		} else {
			gameRow := room.GameState.ToRow()

			// persist the game by saving to database
			go func() {
				err := gomokudb.SaveGameWithMessages(
					room.DB,
					room.GameState.GameID,
					room.GameState.Players[0].PlayerID,
					room.GameState.Players[1].PlayerID,
					gameRow,
					gameRow.Messages,
				)
				if err != nil {
					log.Println("Error saving finished gomoku game to database:", err)
				} else {
					log.Println("Finished gomoku game saved to database:", room.GameState.GameID)
				}
			}()
		}

		// Send final game state ONE MORE TIME just to be safe
		data, _ := json.Marshal(room.GameState)
		res := &GomokuServerResponse{
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

// //////////////////////////
// CLOCK MANAGEMENT
// /////////////////////////
func (room *GomokuRoom) runClock() {
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

func (room *GomokuRoom) tickActivePlayerClock() {
	if room.GameState.Status.Code != "online" {
		return
	}

	currPlayer := GetPlayerByID(room.GameState, room.GameState.Turn)

	if currPlayer.Clock.Remaining <= 0 {
		currPlayer.Clock.Remaining = 0
		room.HandleEvent(GomokuTimeoutEvent{PlayerID: currPlayer.PlayerID})
	} else {
		currPlayer.Clock.Remaining -= time.Second
	}
}

// //////////////////////////
// DATABASE MODEL PERSISTENCE
// /////////////////////////
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

	return &model.GomokuGameStateRow{
		GameID:    gs.GameID,
		BoardSize: gs.Board.Size,
		Players:   players,
		Moves:     moves,
		Messages:  messages,
		Result:    gs.Status.Result,
		Winner:    winner,
	}
}

// //////////////////////////
// CONNECTION MANAGEMENT
// /////////////////////////
func (room *GomokuRoom) ReconnectPlayer(playerID string, conn *websocket.Conn) error {
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

	log.Println("Reconnecting player:", playerID, "to room:", room.RoomID)
	player.ReconnectPlayer(conn)

	resData, err := json.Marshal(room.GameState)
	if err != nil {
		return err
	}

	res := &GomokuServerResponse{
		Type: "update",
		Data: resData,
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Send game state to reconnected player
	room.Send(player, resBytes)

	log.Println("Game state replayed to:", playerID)
	return nil
}
