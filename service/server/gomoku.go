package server

import (
	gomokuapi "boredgamz/api/gomoku"
	gomokucore "boredgamz/core/gomoku"
	"boredgamz/middleware"
)

func (s* Server) MountGomokuLobbies() {
	s.LobbyManager.RegisterLobby("gomoku-casual-19x19", gomokucore.NewGomokuLobby(1000, "gomoku-casual-19x19", s.DB))
	s.LobbyManager.RegisterLobby("gomoku-casual-13x13", gomokucore.NewGomokuLobby(1000, "gomoku-casual-13x13", s.DB))
	s.LobbyManager.RegisterLobby("gomoku-casual-9x9", gomokucore.NewGomokuLobby(1000, "gomoku-casual-9x9", s.DB))
	s.LobbyManager.RegisterLobby("gomoku-ranked-19x19", gomokucore.NewGomokuLobby(1000, "gomoku-ranked-19x19", s.DB))
	s.LobbyManager.RegisterLobby("gomoku-ranked-13x13", gomokucore.NewGomokuLobby(1000, "gomoku-ranked-13x13", s.DB))
	s.LobbyManager.RegisterLobby("gomoku-ranked-9x9", gomokucore.NewGomokuLobby(1000, "gomoku-ranked-9x9", s.DB))
}

func (s *Server) MountGomokuHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-gomoku-lobby", gomokuapi.JoinGomokuLobby(s.LobbyManager))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/reconnect-gomoku-room", gomokuapi.ReconnectToGomokuRoom(s.LobbyManager, s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/game", gomokuapi.GetGomokuGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/games", gomokuapi.GetGomokuGames(s.DB))
}