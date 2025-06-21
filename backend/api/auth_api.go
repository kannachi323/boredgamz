package api

import (
	"encoding/json"
	"log"
	"net/http"
	"playgomoku/backend/db"
	"playgomoku/backend/utils"
	"strings"
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

//auth always returns a user-like object with username and id
type AuthResponse struct {
	Username string `json:"username"`
	UserID string `json:"id"`
}

func CheckAuth(db *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		userID, err := utils.VerifyJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		email, err := db.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&AuthResponse{
			Username: strings.Split(email, "@")[0],
			UserID: userID,
		})
	}

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

		w.WriteHeader(http.StatusCreated)
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
        SameSite: http.SameSiteLaxMode,
        MaxAge:   3600,
    })
		w.WriteHeader(http.StatusOK)
	}
}

