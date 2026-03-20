package server

import (
	connectfourapi "boredgamz/api/connectfour"
	connectfourcore "boredgamz/core/connectfour"
	"boredgamz/middleware"
)

func (s *Server) MountConnectFourLobbies() {
	s.LobbyManager.RegisterLobby(connectfourcore.LobbyID, connectfourcore.NewConnectFourLobby(1000, s.DB))
}

func (s *Server) MountConnectFourHandlers() {
	s.APIRouter.Get("/join-connectfour-lobby", connectfourapi.JoinConnectFourLobby(s.LobbyManager))
	s.APIRouter.Get("/reconnect-connectfour-room", connectfourapi.ReconnectToConnectFourRoom(s.LobbyManager, s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/connectfour/game", connectfourapi.GetConnectFourGame(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/connectfour/games", connectfourapi.GetConnectFourGames(s.DB))
}
