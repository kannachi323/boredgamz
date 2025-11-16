package server

import (
	"playgomoku/backend/api"
	"playgomoku/backend/middleware"
)

func (s *Server) MountAuthHandlers() {
	s.APIRouter.Post("/signup", api.SignUp(s.DB))
	s.APIRouter.Post("/login", api.LogIn(s.DB))
	s.APIRouter.With(middleware.AuthMiddleware).Get("/logout", api.LogOut())
	s.APIRouter.With(middleware.AuthMiddleware).Get("/check-auth", api.CheckAuth(s.DB))
	s.APIRouter.Get("/refresh", api.RefreshAuth(s.DB))
}