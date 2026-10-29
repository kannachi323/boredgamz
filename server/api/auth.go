package api

import (
	"boredgamz/db"
	"boredgamz/middleware"
	"boredgamz/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
)

const maxAuthBody = 8 << 10

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
}

func RefreshAuth(database *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noStore(w)
		cookie, err := r.Cookie("refresh_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "missing refresh token", http.StatusUnauthorized)
			return
		}
		userID, err := utils.VerifyRefreshJWT(cookie.Value)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		if _, err := database.GetUserByID(userID); err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		accessToken, err := utils.GenerateAccessJWT(userID)
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}
		setAuthCookie(w, "access_token", accessToken, 3600)
		w.WriteHeader(http.StatusNoContent)
	}
}

func CheckAuth(database *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noStore(w)
		userID, ok := r.Context().Value(middleware.ContextKey("userID")).(string)
		if !ok || userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := database.GetUserByID(userID)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		WriteJSON(w, http.StatusOK, &AuthResponse{Username: user.Name, UserID: user.ID})
	}
}

func SignUp(database *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noStore(w)
		var request SignUpRequest
		if err := decodeAuthJSON(w, r, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		request.Email = strings.ToLower(strings.TrimSpace(request.Email))
		request.Username = strings.TrimSpace(request.Username)
		if !validUsername(request.Username) || !validEmail(request.Email) || len(request.Password) < 12 || len(request.Password) > 128 {
			http.Error(w, "use a 2–24 character display name, valid email, and a password between 12 and 128 characters", http.StatusBadRequest)
			return
		}
		user, err := database.CreateUser(request.Username, request.Email, request.Password)
		if err != nil {
			http.Error(w, "account could not be created", http.StatusConflict)
			return
		}
		if err := issueAuthCookies(w, user.ID); err != nil {
			http.Error(w, "account was created but sign-in failed", http.StatusInternalServerError)
			return
		}
		WriteJSON(w, http.StatusCreated, &AuthResponse{Username: user.Name, UserID: user.ID})
	}
}

func LogIn(database *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noStore(w)
		var request LogInRequest
		if err := decodeAuthJSON(w, r, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		request.Email = strings.ToLower(strings.TrimSpace(request.Email))
		if !validEmail(request.Email) || len(request.Password) > 128 {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		userID, err := database.GetUserByEmailPassword(request.Email, request.Password)
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		if err := issueAuthCookies(w, userID); err != nil {
			http.Error(w, "failed to generate authentication tokens", http.StatusInternalServerError)
			return
		}
		user, err := database.GetUserByID(userID)
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		WriteJSON(w, http.StatusOK, &AuthResponse{Username: user.Name, UserID: user.ID})
	}
}

func LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		noStore(w)
		setAuthCookie(w, "access_token", "", -1)
		setAuthCookie(w, "refresh_token", "", -1)
		w.WriteHeader(http.StatusNoContent)
	}
}

func decodeAuthJSON(w http.ResponseWriter, r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxAuthBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("invalid request body")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain exactly one JSON object")
	}
	return nil
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && len(value) <= 254
}

func validUsername(value string) bool {
	runes := []rune(value)
	if len(runes) < 2 || len(runes) > 24 {
		return false
	}
	for _, char := range runes {
		if char < 32 || char == 127 {
			return false
		}
	}
	return true
}

func issueAuthCookies(w http.ResponseWriter, userID string) error {
	accessToken, accessErr := utils.GenerateAccessJWT(userID)
	refreshToken, refreshErr := utils.GenerateRefreshJWT(userID)
	if accessErr != nil || refreshErr != nil {
		return errors.New("failed to generate tokens")
	}
	setAuthCookie(w, "access_token", accessToken, 3600)
	setAuthCookie(w, "refresh_token", refreshToken, 30*24*3600)
	return nil
}

func setAuthCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name: name, Value: value, Path: "/", MaxAge: maxAge,
		HttpOnly: true, Secure: os.Getenv("APP_ENV") == "production", SameSite: http.SameSiteLaxMode,
	})
}

func noStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}
