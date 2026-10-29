package arena

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Room struct {
	id             string
	game           *Game
	players        [2]*Player
	spectators     map[string]*Player
	settings       RoomSettings
	turn           string
	winner         string
	seriesWinner   string
	status         string
	lastMove       int
	round          int
	score          map[string]int
	rematchReady   map[string]bool
	chat           []ChatMessage
	lastChatAt     map[string]time.Time
	turnDeadline   *time.Time
	turnTimer      *time.Timer
	timerVersion   uint64
	clockInitial   time.Duration
	clockIncrement time.Duration
	clockRemaining map[string]time.Duration
	turnStarted    time.Time
	drawOfferedBy  string
	mu             sync.Mutex
	updated        time.Time
	cleanupOnce    sync.Once
	onFinished     func(*Room)
}

func newRoom(id string, game *Game, first, second *Player, settings RoomSettings, onFinished func(*Room)) *Room {
	room := &Room{
		id: id, game: game, players: [2]*Player{first, second}, spectators: make(map[string]*Player),
		settings: settings, turn: "one", status: "playing", lastMove: -1, round: 1,
		score: map[string]int{"one": 0, "two": 0}, rematchReady: make(map[string]bool),
		lastChatAt: make(map[string]time.Time), updated: time.Now(), onFinished: onFinished,
	}
	room.initializeChessClocksLocked()
	room.resetTurnTimerLocked()
	return room
}

func (r *Room) HandleClient(playerID string, message ClientMessage) error {
	switch message.Type {
	case "move":
		if r.game.Kind == Chess {
			return r.HandleChessMove(playerID, message.From, message.To, message.Promotion)
		}
		return r.HandleMove(playerID, message.Position)
	case "chat":
		return r.AddChat(playerID, message.Content)
	case "resign":
		return r.Resign(playerID)
	case "rematch":
		return r.ReadyRematch(playerID)
	case "offer-draw":
		return r.OfferDraw(playerID)
	case "accept-draw":
		return r.AcceptDraw(playerID)
	case "decline-draw":
		return r.DeclineDraw(playerID)
	default:
		return errors.New("unsupported message type")
	}
}

func (r *Room) HandleChessMove(playerID, from, to, promotion string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.status != "playing" {
		return errors.New("game is not active")
	}
	player := r.playerByIDLocked(playerID)
	if player == nil {
		return errors.New("player is not seated in this room")
	}
	if player.Mark != r.turn {
		return errors.New("it is not your turn")
	}
	remaining, clocked := r.liveChessClockLocked(player.Mark)
	if clocked && remaining <= 0 {
		r.clockRemaining[player.Mark] = 0
		r.finishRoundLocked(opposite(player.Mark))
		r.broadcastStateLocked()
		return errors.New("time expired")
	}
	if err := r.game.ApplyChessMove(from, to, promotion); err != nil {
		return err
	}
	if r.drawOfferedBy != "" && r.drawOfferedBy != player.Mark {
		r.drawOfferedBy = ""
	}
	if clocked {
		r.clockRemaining[player.Mark] = remaining + r.clockIncrement
	}
	if winner, finished := r.game.ChessResult(); finished {
		r.finishRoundLocked(winner)
	} else {
		r.turn = opposite(r.turn)
		r.resetTurnTimerLocked()
	}
	r.updated = time.Now()
	r.broadcastStateLocked()
	return nil
}

func (r *Room) OfferDraw(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.game.Kind != Chess || r.status != "playing" {
		return errors.New("draw offers are only available during an active chess game")
	}
	player := r.playerByIDLocked(playerID)
	if player == nil {
		return errors.New("only seated players can offer a draw")
	}
	if r.game.Moves == 0 || player.Mark == r.turn {
		return errors.New("offer a draw after making your move")
	}
	if r.drawOfferedBy != "" {
		return errors.New("a draw offer is already pending")
	}
	r.drawOfferedBy = player.Mark
	r.updated = time.Now()
	r.broadcastStateLocked()
	return nil
}

func (r *Room) AcceptDraw(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	player := r.playerByIDLocked(playerID)
	if player == nil || r.drawOfferedBy == "" || r.drawOfferedBy == player.Mark {
		return errors.New("there is no opponent draw offer to accept")
	}
	r.drawOfferedBy = ""
	r.finishRoundLocked("")
	r.updated = time.Now()
	r.broadcastStateLocked()
	return nil
}

func (r *Room) DeclineDraw(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	player := r.playerByIDLocked(playerID)
	if player == nil || r.drawOfferedBy == "" || r.drawOfferedBy == player.Mark {
		return errors.New("there is no opponent draw offer to decline")
	}
	r.drawOfferedBy = ""
	r.updated = time.Now()
	r.broadcastStateLocked()
	return nil
}

func (r *Room) HandleMove(playerID string, position int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.status != "playing" {
		return errors.New("game is not active")
	}
	player := r.playerByIDLocked(playerID)
	if player == nil {
		return errors.New("player is not seated in this room")
	}
	if player.Mark != r.turn {
		return errors.New("it is not your turn")
	}
	actual, err := r.game.ApplyMove(position, player.Mark)
	if err != nil {
		return err
	}
	r.lastMove = actual
	if lineWinner := r.game.Winner(actual); lineWinner != "" {
		winner := lineWinner
		if r.game.Misere {
			winner = opposite(lineWinner)
		}
		r.finishRoundLocked(winner)
	} else if r.game.IsDraw() {
		r.finishRoundLocked("")
	} else {
		r.turn = opposite(r.turn)
		r.resetTurnTimerLocked()
	}
	r.updated = time.Now()
	r.broadcastStateLocked()
	return nil
}

func (r *Room) Resign(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.status != "playing" {
		return errors.New("game is not active")
	}
	player := r.playerByIDLocked(playerID)
	if player == nil {
		return errors.New("only seated players can resign")
	}
	r.finishRoundLocked(opposite(player.Mark))
	r.updated = time.Now()
	r.broadcastStateLocked()
	return nil
}

func (r *Room) ReadyRematch(playerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	player := r.playerByIDLocked(playerID)
	if player == nil {
		return errors.New("only seated players can request a rematch")
	}
	if r.status != "round-finished" && r.status != "finished" {
		return errors.New("finish the current round first")
	}
	r.rematchReady[player.Mark] = true
	if r.rematchReady["one"] && r.rematchReady["two"] {
		if r.status == "finished" {
			r.score = map[string]int{"one": 0, "two": 0}
			r.round = 0
			r.seriesWinner = ""
		}
		r.startNextRoundLocked()
	}
	r.broadcastStateLocked()
	return nil
}

func (r *Room) finishRoundLocked(winner string) {
	r.stopTurnTimerLocked()
	r.winner = winner
	if winner != "" {
		r.score[winner]++
	}
	winsNeeded := r.settings.BestOf/2 + 1
	if winner != "" && r.score[winner] >= winsNeeded {
		r.status, r.seriesWinner = "finished", winner
		r.scheduleCleanupLocked()
	} else {
		r.status = "round-finished"
	}
}

func (r *Room) startNextRoundLocked() {
	game, _ := NewGameWithMode(r.game.Kind, r.settings.Mode)
	r.game = game
	r.round++
	r.winner, r.status, r.lastMove = "", "playing", -1
	r.drawOfferedBy = ""
	r.rematchReady = make(map[string]bool)
	if r.game.Kind == Chess {
		r.turn = "one"
	} else if r.round%2 == 0 {
		r.turn = "two"
	} else {
		r.turn = "one"
	}
	r.initializeChessClocksLocked()
	r.resetTurnTimerLocked()
}

func (r *Room) SendState(player *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sendLocked(player, ServerMessage{Type: "state", State: r.snapshotForLocked(player)})
}

func (r *Room) Snapshot() *Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.snapshotLocked()
}

func (r *Room) snapshotForLocked(viewer *Player) *Snapshot {
	snapshot := r.snapshotLocked()
	if viewer != nil {
		snapshot.ViewerMark = viewer.Mark
	}
	return snapshot
}

func (r *Room) snapshotLocked() *Snapshot {
	var players [2]*PlayerView
	for index, player := range r.players {
		if player != nil {
			players[index] = player.view()
		}
	}
	ready := make([]string, 0, len(r.rematchReady))
	for mark, value := range r.rematchReady {
		if value {
			ready = append(ready, mark)
		}
	}
	return &Snapshot{
		RoomID: r.id, Kind: r.game.Kind, Status: r.status, Rows: r.game.Rows, Columns: r.game.Columns,
		Board: append([]string(nil), r.game.Board...), Turn: r.turn, Winner: r.winner,
		SeriesWinner: r.seriesWinner, Players: players, LastMove: r.lastMove, Settings: r.settings,
		Round: r.round, Score: map[string]int{"one": r.score["one"], "two": r.score["two"]},
		RematchReady: ready, Chat: append([]ChatMessage{}, r.chat...), Spectators: len(r.spectators),
		TurnDeadline:  cloneTime(r.turnDeadline),
		Chess:         r.game.ChessState(),
		Clocks:        r.clockSnapshotLocked(),
		DrawOfferedBy: r.drawOfferedBy,
		GridMoves:     append([]GridMove{}, r.game.gridMoves...),
	}
}

func (r *Room) resetTurnTimerLocked() {
	r.stopTurnTimerLocked()
	duration := timeControlDuration(r.settings.TimeControl)
	if r.game.Kind == Chess && r.clockRemaining != nil {
		duration = r.clockRemaining[r.turn]
	}
	if duration == 0 || r.status != "playing" {
		return
	}
	deadline := time.Now().UTC().Add(duration)
	r.turnStarted = time.Now()
	r.turnDeadline = &deadline
	r.timerVersion++
	version, timedMark := r.timerVersion, r.turn
	r.turnTimer = time.AfterFunc(duration, func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.status != "playing" || r.timerVersion != version || r.turn != timedMark {
			return
		}
		if r.clockRemaining != nil {
			r.clockRemaining[timedMark] = 0
		}
		r.finishRoundLocked(opposite(timedMark))
		r.updated = time.Now()
		r.broadcastStateLocked()
	})
}

func (r *Room) stopTurnTimerLocked() {
	if r.turnTimer != nil {
		r.turnTimer.Stop()
		r.turnTimer = nil
	}
	r.turnDeadline = nil
	r.turnStarted = time.Time{}
	r.timerVersion++
}

func (r *Room) initializeChessClocksLocked() {
	if r.game.Kind != Chess {
		return
	}
	initial, increment := chessTimeControl(r.settings.TimeControl)
	r.clockInitial, r.clockIncrement = initial, increment
	if initial == 0 {
		r.clockRemaining = nil
		return
	}
	r.clockRemaining = map[string]time.Duration{"one": initial, "two": initial}
}

func (r *Room) liveChessClockLocked(mark string) (time.Duration, bool) {
	if r.clockRemaining == nil || r.turnStarted.IsZero() {
		return 0, false
	}
	return r.clockRemaining[mark] - time.Since(r.turnStarted), true
}

func (r *Room) clockSnapshotLocked() map[string]int64 {
	if r.clockRemaining == nil {
		return nil
	}
	one, two := r.clockRemaining["one"], r.clockRemaining["two"]
	if r.status == "playing" && !r.turnStarted.IsZero() {
		elapsed := time.Since(r.turnStarted)
		if r.turn == "one" {
			one = max(0, one-elapsed)
		} else {
			two = max(0, two-elapsed)
		}
	}
	return map[string]int64{"one": one.Milliseconds(), "two": two.Milliseconds()}
}

func chessTimeControl(control string) (time.Duration, time.Duration) {
	switch control {
	case "3m+2s":
		return 3 * time.Minute, 2 * time.Second
	case "10m":
		return 10 * time.Minute, 0
	case "15m+10s":
		return 15 * time.Minute, 10 * time.Second
	default:
		return 0, 0
	}
}

func timeControlDuration(control string) time.Duration {
	if !strings.HasSuffix(control, "s") {
		return 0
	}
	seconds, err := strconv.Atoi(strings.TrimSuffix(control, "s"))
	if err != nil || seconds <= 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func (r *Room) broadcastStateLocked() {
	for _, player := range r.players {
		r.sendLocked(player, ServerMessage{Type: "state", State: r.snapshotForLocked(player)})
	}
	for _, spectator := range r.spectators {
		r.sendLocked(spectator, ServerMessage{Type: "state", State: r.snapshotForLocked(spectator)})
	}
}

func (r *Room) sendLocked(player *Player, message ServerMessage) {
	if player == nil || !player.connected.Load() {
		return
	}
	select {
	case player.sender() <- message:
	default:
	}
}

func (r *Room) playerByIDLocked(id string) *Player {
	for _, player := range r.players {
		if player != nil && player.ID == id {
			return player
		}
	}
	return nil
}

func (r *Room) participantByIDLocked(id string) *Player {
	if player := r.playerByIDLocked(id); player != nil {
		return player
	}
	return r.spectators[id]
}

func (r *Room) allParticipants() []*Player {
	r.mu.Lock()
	defer r.mu.Unlock()
	players := make([]*Player, 0, 2+len(r.spectators))
	players = append(players, r.players[:]...)
	for _, spectator := range r.spectators {
		players = append(players, spectator)
	}
	return players
}

func (r *Room) scheduleCleanupLocked() {
	if r.onFinished != nil {
		r.onFinished(r)
	}
}

func opposite(mark string) string {
	if mark == "one" {
		return "two"
	}
	return "one"
}
