# ✅ Gomoku Server Development TODO

This checklist tracks core progress on building and polishing the backend server for the Gomoku multiplayer game.

---

## 🧱 Core Functionality

- [ ] Switch `chan []byte` to typed channels (`chan ClientRequest`, `chan ServerResponse`)
- [ ] Only accept `lastMove` + `color` from client, not full `GameState`
- [ ] Validate moves server-side (bounds, turn order, already occupied)
- [ ] Track `CurrentTurn` in `GameState`
- [ ] Update game state only on server (AddStoneToBoard, increment turn)
- [ ] Detect win condition (`IsGomoku`) after each move
- [ ] Detect draw (board full with no winner)
- [ ] Broadcast game state to both players after each move
- [ ] Broadcast `"end"` message when game is won or drawn

---

## 🧼 Safety and Clean-Up

- [ ] Handle invalid input with proper error messages
- [ ] Close WebSocket connection and channels cleanly
- [ ] Defer `conn.Close()` in all connection goroutines
- [ ] Use `sync.Once` to ensure `Room.Close()` is only called once
- [ ] Add `recover()` around goroutines to prevent panics from crashing

---

## 🧪 Debugging & Logging

- [ ] Log all valid moves, players, and turn state
- [ ] Avoid logging structs with `sync/atomic` or `chan` directly
- [ ] Log when rooms are created, started, and closed

---

## 🚀 Nice-to-Haves (Future)

- [ ] Add `"waiting"` message when one player is waiting in a room
- [ ] Add reconnection support via `RoomManager.ReconnectPlayer`
- [ ] Support spectator mode
- [ ] Add per-move timers (e.g. 30 seconds per turn)

---

## 🗂 Project Structure

- [ ] Move shared types (e.g. `ClientRequest`, `Move`) into a common `models` or `types` package
- [ ] Separate `StartGame()` from `Room.Start()` loop

---

Feel free to check items off as you complete them ✅
