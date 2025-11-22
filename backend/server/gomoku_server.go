package server

import (
	gomokuapi "boredgamz/api/gomoku"
	gomokucore "boredgamz/core/gomoku"
	"boredgamz/middleware"
)

func (s* Server) MountGomokuLobbies() {
	s.LobbyManager.RegisterLobby("19x19", gomokucore.NewGomokuLobby(1000, "19x19"))
	s.LobbyManager.RegisterLobby("15x15", gomokucore.NewGomokuLobby(1000, "15x15"))
	s.LobbyManager.RegisterLobby("9x9", gomokucore.NewGomokuLobby(1000, "9x9"))
}

func (s *Server) MountGomokuHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-gomoku-lobby", gomokuapi.JoinGomokuLobby(s.LobbyManager))
	s.APIRouter.With(middleware.AuthMiddleware).Post("/gomoku/game", gomokuapi.PostGomokuGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/game", gomokuapi.GetGomokuGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/gomoku/games", gomokuapi.GetGomokuGames(s.DB))
}