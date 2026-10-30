package server

import (
	"boredgamz/api"
	"time"
)

func (s *Server) MountAuthHandlers() {
	credentials := newAuthRateLimiter(10, time.Minute)
	// Refresh runs on every page load and periodically per tab, so it gets a
	// looser budget than password attempts.
	refresh := newRateLimiter(60, time.Minute, "too many session refreshes")
	s.APIRouter.With(credentials.Middleware).Post("/v1/auth/signup", api.SignUp(s.DB))
	s.APIRouter.With(credentials.Middleware).Post("/v1/auth/login", api.LogIn(s.DB))
	s.APIRouter.Post("/v1/auth/logout", api.LogOut())
	s.APIRouter.With(refresh.Middleware).Get("/v1/auth/session", api.Session(s.DB))
	s.APIRouter.With(refresh.Middleware).Post("/v1/auth/refresh", api.RefreshAuth(s.DB))
}
