package server

import (
	gomokuapi "boredgamz/api/gomoku"
	gomokucore "boredgamz/core/gomoku"
	"boredgamz/middleware"
)



func (s* Server) MountConnectFourLobbies() {
	s.LobbyManager.RegisterLobby("19x19", gomokucore.NewGomokuLobby(1000, "19x19", s.DB))
	s.LobbyManager.RegisterLobby("15x15", gomokucore.NewGomokuLobby(1000, "15x15", s.DB))
	s.LobbyManager.RegisterLobby("9x9", gomokucore.NewGomokuLobby(1000, "9x9", s.DB))
}

func (s *Server) MountConnectFourHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-connect-four", gomokuapi.JoinGomokuLobby(s.LobbyManager))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/connectfour/game", gomokuapi.GetGomokuGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/connectfour/games", gomokuapi.GetGomokuGames(s.DB))
}