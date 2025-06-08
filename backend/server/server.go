package server

import (
	"fmt"
	"playgomoku/backend/api"
	"playgomoku/backend/db"
	"playgomoku/backend/manager"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Router       *chi.Mux
	LobbyManager *manager.LobbyManager
	DB	*db.Database
}


func CreateServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
		LobbyManager: manager.NewLobbyManager(),
		DB: &db.Database{},
	}

	s.MountDatabase()
	s.MountResources()
	s.MountHandlers()


	return s
}

func (s *Server) MountDatabase() {
	err := s.DB.Start()
	
	if err != nil {
		fmt.Print("could not mount database: ", err)
	}

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



	s.Router.Get("/join-lobby", api.JoinLobby(s.LobbyManager))
	s.Router.Post("/signup", api.SignUp(s.DB))
}


