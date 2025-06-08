package api

import (
	"encoding/json"
	"log"
	"net/http"
	"playgomoku/backend/db"
)

type SignUpRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	// more fields later...
}

func SignUp(db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("Received sign-up request: %s %s", req.Email, req.Password)
		
		err = db.CreateUser(req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}