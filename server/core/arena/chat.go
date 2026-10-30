package arena

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func (r *Room) AddChat(playerID, content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	player := r.participantByIDLocked(playerID)
	if player == nil {
		return errors.New("player is not in this room")
	}
	content = sanitizeChat(content)
	if content == "" {
		return errors.New("message is empty")
	}
	if time.Since(r.lastChatAt[playerID]) < 750*time.Millisecond {
		return errors.New("messages are being sent too quickly")
	}
	r.lastChatAt[playerID] = time.Now()
	r.chat = append(r.chat, ChatMessage{ID: uuid.NewString(), SenderID: player.ID, SenderName: player.Name, Content: content, CreatedAt: time.Now().UTC()})
	if len(r.chat) > maxChatHistory {
		r.chat = append([]ChatMessage(nil), r.chat[len(r.chat)-maxChatHistory:]...)
	}
	r.broadcastStateLocked()
	return nil
}

func sanitizeChat(content string) string {
	content = strings.TrimSpace(content)
	runes := []rune(content)
	if len(runes) > maxChatLength {
		runes = runes[:maxChatLength]
	}
	return strings.Map(func(value rune) rune {
		if unicode.IsControl(value) && value != '\n' && value != '\t' {
			return -1
		}
		return value
	}, string(runes))
}
