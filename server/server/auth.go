package server

import (
	"boredgamz/api"
	"boredgamz/middleware"
	"time"
)

func (s *Server) MountAuthHandlers() {
	limiter := newAuthRateLimiter(10, time.Minute)
	s.APIRouter.With(limiter.Middleware).Post("/v1/auth/signup", api.SignUp(s.DB))
	s.APIRouter.With(limiter.Middleware).Post("/v1/auth/login", api.LogIn(s.DB))
	s.APIRouter.Post("/v1/auth/logout", api.LogOut())
	s.APIRouter.With(middleware.AuthMiddleware).Get("/v1/auth/session", api.CheckAuth(s.DB))
	s.APIRouter.With(limiter.Middleware).Post("/v1/auth/refresh", api.RefreshAuth(s.DB))
}
