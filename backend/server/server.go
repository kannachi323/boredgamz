package server

import (
	"fmt"
	"playgomoku/backend/api"
	"playgomoku/backend/db"
	"playgomoku/backend/manager"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	Router       *chi.Mux
	APIRouter    *chi.Mux
	LobbyManager *manager.LobbyManager
	DB	*db.Database
}


func CreateServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
		LobbyManager: manager.NewLobbyManager(),
		DB: &db.Database{},
	}
	s.Router.Route("/api", func(r chi.Router) {
		s.APIRouter = r.(*chi.Mux)
	})

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
	s.APIRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	
	//DO NOT REMOVE THIS
	s.APIRouter.Get("/hello", api.HelloWorld())
	
	s.MountAuthHandlers()
	s.MountLobbyHandlers()

	
}


