package main

import (
	"boredgamz/server"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
  s := server.CreateServer()
	
	defer s.DB.Stop()
	
	if err := http.ListenAndServe(":3000", s.Router); err != nil {
    panic(err)
	}
}
