package server

import (
	"playgomoku/backend/api"
	"playgomoku/backend/manager"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Router       *chi.Mux
	LobbyManager *manager.LobbyManager
	//can add db later
}


func CreateServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
		LobbyManager: manager.NewLobbyManager(),
	}

	return s
}

func (s *Server) MountResources() {
	
	// create all the lobbies
	s.LobbyManager.CreateLobby(100, "9x9")

}

func (s *Server) MountHandlers() {
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	
	//all middeware stuff
	s.Router.Use(middleware.Logger)

	//all other handlers
	s.Router.Post("/new-game-state", api.NewGameState)


	s.Router.Get("/join-lobby", api.JoinLobby(s.LobbyManager))
}


