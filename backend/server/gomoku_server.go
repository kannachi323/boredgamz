package server

import (
	"boredgamz/api"
	"boredgamz/core/gomoku"
	"boredgamz/middleware"
)

func (s* Server) MountGomokuLobbies() {
	s.Lobbycore.RegisterLobby("19x19", gomoku.NewGomokuLobby(1000, "19x19"))
	s.Lobbycore.RegisterLobby("15x15", gomoku.NewGomokuLobby(1000, "15x15"))
	s.Lobbycore.RegisterLobby("9x9", gomoku.NewGomokuLobby(1000, "9x9"))
}

func (s *Server) MountGomokuHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-gomoku-lobby", api.JoinGomokuLobby(s.Lobbycore))
}