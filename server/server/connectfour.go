package server

import (
	connectfourapi "boredgamz/api/connectfour"
	connectfourcore "boredgamz/core/connectfour"
	"boredgamz/middleware"
)

func (s *Server) MountConnectFourLobbies() {
	s.LobbyManager.RegisterLobby("connectfour-classic", connectfourcore.NewConnectFourLobby(1000, "connectfour-classic", s.DB))
}

func (s *Server) MountConnectFourHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-connect-four", connectfourapi.JoinConnectFourLobby(s.LobbyManager))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/connectfour/game", connectfourapi.GetConnectFourGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/connectfour/games", connectfourapi.GetConnectFourGames(s.DB))
}
