package server

import (
	"playgomoku/backend/api"
	"playgomoku/backend/middleware"
)

func (s *Server) MountLobbyHandlers() {
	s.APIRouter.With(middleware.AuthMiddleware).Get("/join-lobby", api.JoinLobby(s.LobbyManager))
}