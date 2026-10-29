package arena

import (
	"sync"
	"sync/atomic"
)

type Player struct {
	ID   string
	Name string
	Mark string

	options JoinOptions // requested settings; the host's win on private tables

	mu         sync.RWMutex
	send       chan ServerMessage
	connected  atomic.Bool
	left       atomic.Bool // pressed Leave; won't be back
	generation atomic.Uint64
}

func newPlayer(id, name string) *Player {
	player := &Player{ID: id, Name: SanitizeName(name)}
	player.replaceSender()
	return player
}

func (p *Player) replaceSender() {
	p.mu.Lock()
	p.send = make(chan ServerMessage, 32)
	p.mu.Unlock()
	p.connected.Store(true)
	p.left.Store(false)
	p.generation.Add(1)
}

func (p *Player) sender() chan ServerMessage {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.send
}

func (p *Player) view() *PlayerView {
	return &PlayerView{ID: p.ID, Name: p.Name, Mark: p.Mark, Connected: p.connected.Load(), Left: p.left.Load()}
}
