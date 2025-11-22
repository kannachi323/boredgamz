package server

import (
	"boredgamz/api"
	"boredgamz/core"
	"boredgamz/db"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	Router       *chi.Mux
	APIRouter    *chi.Mux
	LobbyManager *core.LobbyManager
	DB	*db.Database
}


func CreateServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
		APIRouter: chi.NewRouter(),
		LobbyManager: core.NewLobbyManager(),
		DB: &db.Database{},
	}
	
	s.Router.Mount("/api", s.APIRouter)

	s.MountDatabase()
	s.MountHandlers()
	s.MountLobbies()

	return s
}

func (s *Server) MountDatabase() {
	err := s.DB.Start()
	if err != nil {
		fmt.Print("could not mount database: ", err)
	}
}

func (s* Server) MountLobbies() {
	//IMPORTANT: Do not remove this lobby registration
	s.MountGomokuLobbies()
	
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
	s.MountGomokuHandlers()
}


