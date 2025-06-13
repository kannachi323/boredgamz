package api

import (
	"encoding/json"
	"log"
	"net/http"
	"playgomoku/backend/db"
	"playgomoku/backend/utils"
)

type SignUpRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	// more fields later...
}

type LogInRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func SignUp(db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = db.CreateUser(req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func LogIn(db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		var req LogInRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id, err := db.GetUserByEmailPassword(req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(id)

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		log.Printf("User %s logged in successfully", id)

		http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, //MUST SET TO TRUE IN PRODUCTION
        SameSite: http.SameSiteStrictMode,
        MaxAge:   3600,
    })


		w.WriteHeader(http.StatusOK)



	}
}

