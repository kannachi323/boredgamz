package server

import (
	gomokuapi "boredgamz/api/gomoku"
	gomokucore "boredgamz/core/gomoku"
	"boredgamz/middleware"
)

func (s* Server) MountGomokuLobbies() {
	modes := []string{"casual", "ranked", "bots"}
	boards := []string{"19x19", "13x13", "9x9"}
	rules := []string{"freestyle", "standard", "renju"}

	for _, mode := range modes {
		for _, board := range boards {
			for _, rule := range rules {
				lobbyID := gomokucore.GetGomokuLobbyID(board, mode, rule)
				s.LobbyManager.RegisterLobby(lobbyID, gomokucore.NewGomokuLobby(1000, lobbyID, s.DB))
			}
		}
	}
}

func (s *Server) MountGomokuHandlers() {
	s.APIRouter.Get("/join-gomoku-lobby", gomokuapi.JoinGomokuLobby(s.LobbyManager))
	s.APIRouter.Get("/join-gomoku-bot-lobby", gomokuapi.JoinGomokuBotLobby(s.LobbyManager))
	s.APIRouter.Get("/reconnect-gomoku-room", gomokuapi.ReconnectToGomokuRoom(s.LobbyManager, s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/game", gomokuapi.GetGomokuGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/games", gomokuapi.GetGomokuGames(s.DB))
}