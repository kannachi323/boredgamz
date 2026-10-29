package server

import (
	"boredgamz/api"
	"boredgamz/core/arena"
	"boredgamz/core/mahjong"
	"boredgamz/core/poker"
	"boredgamz/db"
	"boredgamz/utils"
	"fmt"
	"log/slog"
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
	Arena          *arena.Hub
	Mahjong        *mahjong.Hub
	Poker          *poker.Hub
	DB             *db.Database
	StartedAt      time.Time
	allowedOrigins map[string]struct{}
	socketLimiter  *authRateLimiter
}

func CreateServer() (*Server, error) {
	s := &Server{
		Router:         chi.NewRouter(),
		APIRouter:      chi.NewRouter(),
		Arena:          arena.NewHub(),
		Mahjong:        mahjong.NewHub(),
		Poker:          poker.NewHub(),
		DB:             &db.Database{},
		StartedAt:      time.Now(),
		allowedOrigins: parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS")),
		socketLimiter:  newRateLimiter(60, time.Minute, "too many connection attempts"),
	}

	if err := s.MountDatabase(); err != nil {
		return nil, err
	}
	s.MountHandlers()
	s.MountStatic()

	return s, nil
}

func (s *Server) MountDatabase() error {
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) == "" {
		slog.Warn("DATABASE_URL is not set; accounts are disabled and sockets accept guests (development only)")
		return nil
	}
	if err := s.DB.Start(); err != nil {
		return fmt.Errorf("start account database: %w", err)
	}
	return nil
}

func (s *Server) accountsEnabled() bool {
	return s.DB != nil && s.DB.Pool != nil
}

func (s *Server) MountHandlers() {
	s.Router.Use(chimiddleware.RequestID)
	if os.Getenv("TRUST_PROXY") == "true" {
		// Only honour X-Forwarded-For / X-Real-IP behind a proxy that sets them;
		// otherwise clients could pick their own rate-limit key.
		s.Router.Use(chimiddleware.RealIP)
	}
	s.Router.Use(requestLogger)
	s.Router.Use(chimiddleware.Recoverer)
	s.Router.Use(securityHeaders)

	s.APIRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOriginPatterns(s.allowedOrigins),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	s.APIRouter.Use(s.requireTrustedOrigin)

	s.APIRouter.Get("/v1/catalog", api.Catalog())
	s.APIRouter.Get("/v1/status", s.RuntimeStatus())
	s.APIRouter.With(s.socketLimiter.Middleware).Get("/v1/arena/{game}", s.ArenaSocket())
	s.APIRouter.With(s.socketLimiter.Middleware).Get("/v1/mahjong", s.MahjongSocket())
	s.APIRouter.With(s.socketLimiter.Middleware).Get("/v1/poker", s.PokerSocket())
	s.Router.Get("/healthz", api.Health(s.StartedAt, s.accountsEnabled))
	s.Router.Get("/readyz", api.Ready(s.DB.Ready))

	if s.accountsEnabled() {
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

// socketParams is the validated part of a socket handshake shared by every game.
type socketParams struct {
	PlayerID, Name, Role, Mode, InviteCode, RoomID string
}

// prepareSocket authenticates and validates a handshake. It writes the HTTP
// error itself and returns ok=false when the upgrade must not happen.
func (s *Server) prepareSocket(w http.ResponseWriter, r *http.Request) (socketParams, bool) {
	query := r.URL.Query()
	params := socketParams{
		PlayerID:   strings.TrimSpace(query.Get("playerId")),
		Name:       query.Get("name"),
		Role:       query.Get("role"),
		Mode:       boundedParam(query.Get("mode"), 32),
		InviteCode: boundedParam(query.Get("inviteCode"), 32),
		RoomID:     boundedParam(query.Get("roomId"), 64),
	}
	if params.PlayerID == "" || len(params.PlayerID) > 80 {
		http.Error(w, "invalid playerId", http.StatusBadRequest)
		return params, false
	}
	if params.Role != "" && params.Role != "player" && params.Role != "spectator" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return params, false
	}
	id, name, ok := s.socketIdentity(r, params.PlayerID, params.Name, params.Role)
	if !ok {
		http.Error(w, "sign in to join a table", http.StatusUnauthorized)
		return params, false
	}
	params.PlayerID, params.Name = id, name
	return params, true
}

func boundedParam(value string, limit int) string {
	value = strings.TrimSpace(value)
	if len(value) > limit {
		return value[:limit]
	}
	return value
}

func (s *Server) upgrader() websocket.Upgrader {
	return websocket.Upgrader{HandshakeTimeout: 5 * time.Second, ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: s.originAllowed}
}

func (s *Server) MahjongSocket() http.HandlerFunc {
	upgrader := s.upgrader()
	return func(w http.ResponseWriter, r *http.Request) {
		params, ok := s.prepareSocket(w, r)
		if !ok {
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		options := mahjong.JoinOptions{PlayerID: params.PlayerID, Name: params.Name, Mode: params.Mode, InviteCode: params.InviteCode, Role: params.Role, RoomID: params.RoomID}
		if err := mahjong.Serve(s.Mahjong, options, conn); err != nil {
			slog.Debug("mahjong connection ended", "err", err)
		}
	}
}

func (s *Server) PokerSocket() http.HandlerFunc {
	upgrader := s.upgrader()
	return func(w http.ResponseWriter, r *http.Request) {
		params, ok := s.prepareSocket(w, r)
		if !ok {
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		options := poker.JoinOptions{PlayerID: params.PlayerID, Name: params.Name, Mode: params.Mode, InviteCode: params.InviteCode, Role: params.Role, RoomID: params.RoomID}
		if err := poker.Serve(s.Poker, options, conn); err != nil {
			slog.Debug("poker connection ended", "err", err)
		}
	}
}

func (s *Server) ArenaSocket() http.HandlerFunc {
	upgrader := s.upgrader()
	return func(w http.ResponseWriter, r *http.Request) {
		kind := arena.GameKind(chi.URLParam(r, "game"))
		if !arena.ValidKind(kind) {
			http.Error(w, "unsupported game", http.StatusNotFound)
			return
		}
		params, ok := s.prepareSocket(w, r)
		if !ok {
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		bestOf, _ := strconv.Atoi(r.URL.Query().Get("bestOf"))
		options := arena.JoinOptions{
			Kind: kind, PlayerID: params.PlayerID, PlayerName: params.Name,
			Mode: params.Mode, BestOf: bestOf,
			TimeControl: boundedParam(r.URL.Query().Get("timeControl"), 16), InviteCode: params.InviteCode,
			Role: params.Role, RoomID: params.RoomID,
		}
		if err := arena.ServePlayer(s.Arena, options, conn); err != nil {
			slog.Debug("arena connection ended", "err", err)
		}
	}
}

// socketIdentity resolves who is connecting. With accounts enabled (always in
// production) a valid access cookie is required and the display name comes
// from the account, never the query string. Without a database (local
// development and tests only) the client-supplied guest ID is used.
func (s *Server) socketIdentity(r *http.Request, fallbackID, fallbackName, role string) (string, string, bool) {
	if !s.accountsEnabled() {
		return fallbackID, fallbackName, true
	}
	cookie, err := r.Cookie("access_token")
	if err != nil || cookie.Value == "" {
		return "", "", false
	}
	userID, err := utils.VerifyJWT(cookie.Value)
	if err != nil {
		return "", "", false
	}
	user, err := s.DB.GetUserByID(userID)
	if err != nil {
		return "", "", false
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
	return identity, user.Name, true
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
		for _, host := range []string{"localhost", "127.0.0.1"} {
			for _, port := range []string{"3000", "5201", "8201"} {
				origins["http://"+host+":"+port] = struct{}{}
			}
		}
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
	return s.isAllowedOrigin(origin)
}

func (s *Server) isAllowedOrigin(origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return false
	}
	_, ok := s.allowedOrigins[origin]
	return ok
}

// requireTrustedOrigin is the CSRF defence for cookie-authenticated writes.
// Session cookies are SameSite=Lax already; on top of that every state-changing
// request must carry an Origin (or Referer) from ALLOWED_ORIGINS.
func (s *Server) requireTrustedOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		origin := strings.TrimRight(r.Header.Get("Origin"), "/")
		if origin == "" {
			if referer, err := url.Parse(r.Header.Get("Referer")); err == nil && referer.Host != "" {
				origin = referer.Scheme + "://" + referer.Host
			}
		}
		if origin == "" && os.Getenv("APP_ENV") != "production" {
			next.ServeHTTP(w, r)
			return
		}
		if !s.isAllowedOrigin(origin) {
			http.Error(w, "cross-site request rejected", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		path := r.URL.Path // the SPA handler rewrites r.URL.Path
		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		status := ww.Status()
		if status == 0 {
			status = http.StatusSwitchingProtocols
		}
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		}
		slog.LogAttrs(r.Context(), level, "request",
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
			slog.String("request_id", chimiddleware.GetReqID(r.Context())),
			slog.String("remote", clientAddress(r.RemoteAddr)),
		)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src 'self' ws: wss:; img-src 'self' data:; font-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		if os.Getenv("APP_ENV") == "production" {
			w.Header().Set("Strict-Transport-Security", "max-age=15552000")
		}
		next.ServeHTTP(w, r)
	})
}
