package connectfour

import (
	"boredgamz/core"
	"boredgamz/db"
	"container/list"
	"log"
	"sync"
	"time"
)

type ConnectFourLobby struct {
    *core.Lobby
    GameType   string
    RedQueue   *list.List
    YellowQueue *list.List
    PlayerMap  map[*core.Player]*ConnectFourLobbySlot
    LobbyType  string
    mu         sync.RWMutex
}

type ConnectFourLobbySlot struct {
    Element *list.Element
    Queue   *list.List
}

func NewConnectFourLobby(maxPlayers int, gameType string, db *db.Database) core.LobbyController {
    lobby := &ConnectFourLobby{
        Lobby: &core.Lobby{
            NumPlayers:  0,
            MaxPlayers:  maxPlayers,
            RoomManager: core.NewRoomManager(),
            DB:          db,
        },
        GameType:   gameType,
        RedQueue:   list.New(),
        YellowQueue: list.New(),
        LobbyType:  gameType,
        PlayerMap:  make(map[*core.Player]*ConnectFourLobbySlot),
    }

    go lobby.MatchPlayers() // background matchmaking loop

    return lobby
}

func (lobby *ConnectFourLobby) AddPlayer(player *core.Player) {
    lobby.mu.Lock()
    defer lobby.mu.Unlock()

    if lobby.NumPlayers >= lobby.MaxPlayers {
        return
    }
    if _, exists := lobby.PlayerMap[player]; exists {
        return
    }

    var elem *list.Element
    var queue *list.List

    switch player.Color {
    case "red": // Player chooses red
        log.Println("Adding player to red queue:", player.PlayerID)
        elem = lobby.RedQueue.PushBack(player)
        queue = lobby.RedQueue
    case "yellow": // Player chooses yellow
        log.Println("Adding player to yellow queue:", player.PlayerID)
        elem = lobby.YellowQueue.PushBack(player)
        queue = lobby.YellowQueue
    default:
        // Invalid color
        return
    }

    lobby.PlayerMap[player] = &ConnectFourLobbySlot{
        Element: elem,
        Queue:   queue,
    }

    lobby.NumPlayers++
}

func (lobby *ConnectFourLobby) RemovePlayer(player *core.Player) {
    lobby.mu.Lock()
    defer lobby.mu.Unlock()

    slot, ok := lobby.PlayerMap[player]
    if !ok {
        return
    }

    slot.Queue.Remove(slot.Element)
    delete(lobby.PlayerMap, player)
    lobby.NumPlayers--
}

func (lobby *ConnectFourLobby) MatchPlayers() {
    for {
        r, y, ok := lobby.tryMatch()
        if ok {
            log.Println("Matched:", r.PlayerID, y.PlayerID)

            room := NewConnectFourRoom(r, y, lobby.LobbyType, lobby.DB)
            if room != nil {
                lobby.RoomManager.RegisterPlayerToRoom(r.PlayerID, room)
                lobby.RoomManager.RegisterPlayerToRoom(y.PlayerID, room)

                r.StartPlayer()
                y.StartPlayer()
                room.Start()
            }
        }

        time.Sleep(100 * time.Millisecond) // prevent spin-loop
    }
}

// PRIVATE
func (lobby *ConnectFourLobby) tryMatch() (*core.Player, *core.Player, bool) {
    lobby.mu.Lock()
    defer lobby.mu.Unlock()

    for lobby.RedQueue.Len() > 0 && lobby.YellowQueue.Len() > 0 {

        r := lobby.RedQueue.Front().Value.(*core.Player)
        y := lobby.YellowQueue.Front().Value.(*core.Player)

        // Drop disconnected players
        if r.Conn == nil {
            lobby.RedQueue.Remove(lobby.RedQueue.Front())
            delete(lobby.PlayerMap, r)
            lobby.NumPlayers--
            continue
        }
        if y.Conn == nil {
            lobby.YellowQueue.Remove(lobby.YellowQueue.Front())
            delete(lobby.PlayerMap, y)
            lobby.NumPlayers--
            continue
        }

        if r.PlayerID == y.PlayerID {
            // Corrupted: same player in both queues
            lobby.RedQueue.Remove(lobby.RedQueue.Front())
            lobby.YellowQueue.Remove(lobby.YellowQueue.Front())
            delete(lobby.PlayerMap, r)
            delete(lobby.PlayerMap, y)
            lobby.NumPlayers -= 2
            continue
        }

        // VALID MATCH
        lobby.RedQueue.Remove(lobby.RedQueue.Front())
        lobby.YellowQueue.Remove(lobby.YellowQueue.Front())

        delete(lobby.PlayerMap, r)
        delete(lobby.PlayerMap, y)
        lobby.NumPlayers -= 2

        return r, y, true
    }

    return nil, nil, false
}
