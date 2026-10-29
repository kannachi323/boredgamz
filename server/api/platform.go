package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type GameCatalogItem struct {
	Slug        string   `json:"slug"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Players     string   `json:"players"`
	Modes       []string `json:"modes"`
	Difficulty  string   `json:"difficulty"`
}

var GameCatalog = []GameCatalogItem{
	{Slug: "tictactoe", Name: "Tic-Tac-Toe", Description: "Three compact live rulesets with series play.", Players: "2", Modes: []string{"classic", "four-square", "misere"}, Difficulty: "Easy"},
	{Slug: "connect-four", Name: "Connect Four", Description: "Three gravity-board formats with live clocks.", Players: "2", Modes: []string{"classic", "quick", "connect-five"}, Difficulty: "Easy"},
	{Slug: "gomoku", Name: "Gomoku", Description: "Freestyle, exact-five, and compact club boards.", Players: "2", Modes: []string{"freestyle", "pro", "compact"}, Difficulty: "Tactical"},
	{Slug: "chess", Name: "Chess", Description: "Server-verified chess with real clocks and notation.", Players: "2", Modes: []string{"blitz", "rapid", "standard"}, Difficulty: "Strategic"},
	{Slug: "poker", Name: "Poker", Description: "Heads-up private-card poker with authoritative wagering.", Players: "2", Modes: []string{"holdem", "five-card-draw"}, Difficulty: "Strategic"},
	{Slug: "mahjong", Name: "Mahjong Solitaire", Description: "Clear a shared tile garden or race mirrored layouts.", Players: "2", Modes: []string{"garden", "duel", "zen"}, Difficulty: "Relaxed"},
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func Health(startedAt time.Time, databaseReady func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{
			"status":        "ok",
			"service":       "boredgamz",
			"uptimeSeconds": int(time.Since(startedAt).Seconds()),
			"database":      databaseReady(),
		})
	}
}

// Ready reports 503 while a dependency (the database) is unreachable so a
// load balancer stops routing new players here.
func Ready(check func(context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := check(ctx); err != nil {
			WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "service": "boredgamz"})
			return
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ready", "service": "boredgamz"})
	}
}

func Catalog() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{"games": GameCatalog})
	}
}
