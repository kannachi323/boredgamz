package api

import (
	"boredgamz/db"
	"boredgamz/utils"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"unicode"
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

// Session reports who is signed in. It never answers 401: a missing or
// expired access cookie falls back to the refresh cookie (rotating the access
// cookie), and a signed-out visitor gets 204 so page loads stay error-free.
func Session(database *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noStore(w)
		userID := ""
		if cookie, err := r.Cookie("access_token"); err == nil && cookie.Value != "" {
			userID, _ = utils.VerifyJWT(cookie.Value)
		}
		if userID == "" {
			if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
				if id, err := utils.VerifyRefreshJWT(cookie.Value); err == nil {
					if accessToken, err := utils.GenerateAccessJWT(id); err == nil {
						setAuthCookie(w, "access_token", accessToken, 3600)
						userID = id
					}
				}
			}
		}
		if userID == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		user, err := database.GetUserByID(userID)
		if err != nil {
			setAuthCookie(w, "access_token", "", -1)
			setAuthCookie(w, "refresh_token", "", -1)
			w.WriteHeader(http.StatusNoContent)
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
		if !validUsername(request.Username) || !validEmail(request.Email) || !validPassword(request.Password) {
			http.Error(w, "use a 2–24 character display name, a valid email, and a password of 12–72 characters", http.StatusBadRequest)
			return
		}
		user, err := database.CreateUser(request.Username, request.Email, request.Password)
		if errors.Is(err, db.ErrEmailTaken) {
			http.Error(w, "account could not be created", http.StatusConflict)
			return
		}
		if err != nil {
			slog.Error("signup failed", "err", err)
			http.Error(w, "account could not be created", http.StatusInternalServerError)
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
		if !validEmail(request.Email) || len(request.Password) > maxPasswordBytes {
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
	// Requiring JSON keeps HTML forms (which cannot set this type) out.
	if mediaType := strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]); !strings.EqualFold(mediaType, "application/json") {
		return errors.New("request body must be JSON")
	}
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

// bcrypt only reads the first 72 bytes and x/crypto rejects longer input, so
// the limit is enforced up front instead of failing inside CreateUser.
const minPasswordRunes, maxPasswordBytes = 12, 72

func validPassword(value string) bool {
	return len([]rune(value)) >= minPasswordRunes && len(value) <= maxPasswordBytes
}

func validUsername(value string) bool {
	runes := []rune(value)
	if len(runes) < 2 || len(runes) > 24 {
		return false
	}
	for _, char := range runes {
		if unicode.IsControl(char) || unicode.Is(unicode.Cf, char) {
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
