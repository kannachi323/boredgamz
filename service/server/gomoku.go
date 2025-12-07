package server

import (
	gomokuapi "boredgamz/api/gomoku"
	gomokucore "boredgamz/core/gomoku"
	"boredgamz/middleware"
)

func (s* Server) MountGomokuLobbies() {
	s.LobbyManager.RegisterLobby("19x19", gomokucore.NewGomokuLobby(1000, "19x19", s.DB))
	s.LobbyManager.RegisterLobby("15x15", gomokucore.NewGomokuLobby(1000, "13x13", s.DB))
	s.LobbyManager.RegisterLobby("9x9", gomokucore.NewGomokuLobby(1000, "9x9", s.DB))
}

func (s *Server) MountGomokuHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-gomoku-lobby", gomokuapi.JoinGomokuLobby(s.LobbyManager))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/reconnect-gomoku-room", gomokuapi.ReconnectToGomokuRoom(s.LobbyManager, s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/game", gomokuapi.GetGomokuGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/games", gomokuapi.GetGomokuGames(s.DB))
}