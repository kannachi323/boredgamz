package server

import (
	"boredgamz/api"
	"boredgamz/core"
	"boredgamz/core/arena"
	"boredgamz/core/mahjong"
	"boredgamz/core/poker"
	"boredgamz/db"
	"boredgamz/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

type Server struct {
	Router         *chi.Mux
	APIRouter      *chi.Mux
	LobbyManager   *core.LobbyManager
	Arena          *arena.Hub
	Mahjong        *mahjong.Hub
	Poker          *poker.Hub
	DB             *db.Database
	StartedAt      time.Time
	allowedOrigins map[string]struct{}
}

func CreateServer() (*Server, error) {
	s := &Server{
		Router:         chi.NewRouter(),
		LobbyManager:   core.NewLobbyManager(),
		Arena:          arena.NewHub(),
		Mahjong:        mahjong.NewHub(),
		Poker:          poker.NewHub(),
		DB:             &db.Database{},
		StartedAt:      time.Now(),
		allowedOrigins: parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS")),
	}
	s.APIRouter = chi.NewRouter()

	if err := s.MountDatabase(); err != nil {
		return nil, err
	}
	s.MountHandlers()
	s.MountStatic()

	return s, nil
}

func (s *Server) MountDatabase() error {
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) == "" {
		log.Print("DATABASE_URL is not set; account persistence is disabled")
		return nil
	}
	if err := s.DB.Start(); err != nil {
		return fmt.Errorf("start account database: %w", err)
	}
	return nil
}

func (s *Server) MountLobbies() {
	//IMPORTANT: Do not remove this lobby registration
	s.MountGomokuLobbies()
	s.MountConnectFourLobbies()

}

func (s *Server) MountHandlers() {
	s.Router.Use(chimiddleware.RequestID)
	s.Router.Use(chimiddleware.RealIP)
	s.Router.Use(chimiddleware.Recoverer)
	s.Router.Use(securityHeaders)

	s.APIRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOriginPatterns(s.allowedOrigins),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.APIRouter.Get("/hello", api.HelloWorld())
	s.APIRouter.Get("/v1/catalog", api.Catalog())
	s.APIRouter.Get("/v1/status", s.RuntimeStatus())
	s.APIRouter.Get("/v1/arena/{game}", s.ArenaSocket())
	s.APIRouter.Get("/v1/mahjong", s.MahjongSocket())
	s.APIRouter.Get("/v1/poker", s.PokerSocket())
	s.Router.Get("/healthz", api.Health(s.StartedAt, func() bool { return s.DB != nil && s.DB.Pool != nil }))
	s.Router.Get("/readyz", api.Ready())

	if s.DB != nil && s.DB.Pool != nil {
		s.MountAuthHandlers()
	}
	s.Router.Mount("/api", s.APIRouter)
}

func (s *Server) RuntimeStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		api.WriteJSON(w, http.StatusOK, map[string]any{
			"service":       "boredgamz",
			"uptimeSeconds": int(time.Since(s.StartedAt).Seconds()),
			"arena":         s.Arena.Stats(),
			"poker":         s.Poker.Stats(),
			"mahjong":       s.Mahjong.Stats(),
		})
	}
}

func (s *Server) MahjongSocket() http.HandlerFunc {
	upgrader := websocket.Upgrader{HandshakeTimeout: 5 * time.Second, ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: s.originAllowed}
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := strings.TrimSpace(r.URL.Query().Get("playerId"))
		if playerID == "" || len(playerID) > 80 {
			http.Error(w, "invalid playerId", http.StatusBadRequest)
			return
		}
		playerID, playerName := s.socketIdentity(r, playerID, r.URL.Query().Get("name"), r.URL.Query().Get("role"))
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		options := mahjong.JoinOptions{PlayerID: playerID, Name: playerName, Mode: r.URL.Query().Get("mode"), InviteCode: r.URL.Query().Get("inviteCode"), Role: r.URL.Query().Get("role"), RoomID: r.URL.Query().Get("roomId")}
		if err := mahjong.Serve(s.Mahjong, options, conn); err != nil {
			log.Printf("mahjong connection ended: %v", err)
		}
	}
}

func (s *Server) PokerSocket() http.HandlerFunc {
	upgrader := websocket.Upgrader{HandshakeTimeout: 5 * time.Second, ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: s.originAllowed}
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := strings.TrimSpace(r.URL.Query().Get("playerId"))
		if playerID == "" || len(playerID) > 80 {
			http.Error(w, "invalid playerId", http.StatusBadRequest)
			return
		}
		playerID, playerName := s.socketIdentity(r, playerID, r.URL.Query().Get("name"), r.URL.Query().Get("role"))
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		options := poker.JoinOptions{PlayerID: playerID, Name: playerName, Mode: r.URL.Query().Get("mode"), InviteCode: r.URL.Query().Get("inviteCode"), Role: r.URL.Query().Get("role"), RoomID: r.URL.Query().Get("roomId")}
		if err := poker.Serve(s.Poker, options, conn); err != nil {
			log.Printf("poker connection ended: %v", err)
		}
	}
}

func (s *Server) ArenaSocket() http.HandlerFunc {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      s.originAllowed,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		kind := arena.GameKind(chi.URLParam(r, "game"))
		if !arena.ValidKind(kind) {
			http.Error(w, "unsupported game", http.StatusNotFound)
			return
		}
		playerID := strings.TrimSpace(r.URL.Query().Get("playerId"))
		playerName := r.URL.Query().Get("name")
		if playerID == "" || len(playerID) > 80 {
			http.Error(w, "invalid playerId", http.StatusBadRequest)
			return
		}
		playerID, playerName = s.socketIdentity(r, playerID, playerName, r.URL.Query().Get("role"))
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		bestOf, _ := strconv.Atoi(r.URL.Query().Get("bestOf"))
		options := arena.JoinOptions{
			Kind: kind, PlayerID: playerID, PlayerName: playerName,
			Mode: r.URL.Query().Get("mode"), BestOf: bestOf,
			TimeControl: r.URL.Query().Get("timeControl"), InviteCode: r.URL.Query().Get("inviteCode"),
			Role: r.URL.Query().Get("role"), RoomID: r.URL.Query().Get("roomId"),
		}
		if err := arena.ServePlayer(s.Arena, options, conn); err != nil {
			log.Printf("arena connection ended: %v", err)
		}
	}
}

func (s *Server) socketIdentity(r *http.Request, fallbackID, fallbackName, role string) (string, string) {
	if s.DB == nil || s.DB.Pool == nil {
		return fallbackID, fallbackName
	}
	cookie, err := r.Cookie("access_token")
	if err != nil || cookie.Value == "" {
		return fallbackID, fallbackName
	}
	userID, err := utils.VerifyJWT(cookie.Value)
	if err != nil {
		return fallbackID, fallbackName
	}
	user, err := s.DB.GetUserByID(userID)
	if err != nil {
		return fallbackID, fallbackName
	}
	tabID := fallbackID
	parts := strings.Split(fallbackID, ":")
	if role == "spectator" && len(parts) > 1 && parts[len(parts)-1] == "spectator" {
		parts = parts[:len(parts)-1]
	}
	if len(parts) > 1 {
		tabID = parts[len(parts)-1]
	}
	if len(tabID) > 12 {
		tabID = tabID[len(tabID)-12:]
	}
	identity := user.ID + ":" + tabID
	if role == "spectator" {
		identity += ":spectator"
	}
	return identity, user.Name
}

func parseAllowedOrigins(value string) map[string]struct{} {
	origins := map[string]struct{}{}
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(strings.TrimRight(item, "/"))
		if item != "" {
			origins[item] = struct{}{}
		}
	}
	if len(origins) == 0 && os.Getenv("APP_ENV") != "production" {
		origins["http://localhost:5173"] = struct{}{}
		origins["http://127.0.0.1:5173"] = struct{}{}
		origins["http://[::1]:5173"] = struct{}{}
		origins["http://localhost:5174"] = struct{}{}
		origins["http://127.0.0.1:5174"] = struct{}{}
		origins["http://localhost:3000"] = struct{}{}
		origins["http://127.0.0.1:3000"] = struct{}{}
	}
	return origins
}

func allowedOriginPatterns(origins map[string]struct{}) []string {
	patterns := make([]string, 0, len(origins))
	for origin := range origins {
		patterns = append(patterns, origin)
	}
	return patterns
}

func (s *Server) originAllowed(r *http.Request) bool {
	origin := strings.TrimRight(r.Header.Get("Origin"), "/")
	if origin == "" {
		return os.Getenv("APP_ENV") != "production"
	}
	parsed, err := url.Parse(origin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return false
	}
	_, ok := s.allowedOrigins[origin]
	return ok
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self' ws: wss:; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'")
		if os.Getenv("APP_ENV") == "production" {
			w.Header().Set("Strict-Transport-Security", "max-age=15552000")
		}
		next.ServeHTTP(w, r)
	})
}
