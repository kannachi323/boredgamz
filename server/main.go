package main

import (
	"boredgamz/server"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	s := server.CreateServer()

	defer s.DB.Stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		panic(err)
	}
}
