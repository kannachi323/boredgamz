package main

import (
	"net/http"
	"playgomoku/backend/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
    s := server.CreateServer()
	
	defer s.DB.Stop()
	
	http.ListenAndServe(":3000", s.Router)
}
