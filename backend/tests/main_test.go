package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"boredgamz/api"
	"boredgamz/core"
	"boredgamz/core/gomoku"
	"boredgamz/db"
	gomokudb "boredgamz/db/gomoku"
	"boredgamz/middleware"
	"boredgamz/server"
	"boredgamz/utils"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var testDB *db.Database

func StartTestDB(testDB *db.Database) error {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_TEST_URL")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	

	testDB.Pool = pool

	return nil
}	

func CloseTestDB(testDB *db.Database) {
	if testDB.Pool != nil {
		testDB.Pool.Close()
	}
}

func ResetTestDB(testDB *db.Database) error {
    ctx := context.Background()

    _, err := testDB.Pool.Exec(ctx, `
DO $$
DECLARE
    stmt text;
BEGIN
    SELECT 'TRUNCATE TABLE ' || string_agg(format('%I.%I', schemaname, tablename), ', ') || ' RESTART IDENTITY CASCADE;'
    INTO stmt
    FROM pg_tables
    WHERE schemaname = 'public';

    EXECUTE stmt;
END $$;`)
    if err != nil {
        return fmt.Errorf("failed to truncate all tables: %w", err)
    }

    // seed (optional)
    sqlBytes, err := os.ReadFile("../db/sql/data.sql")
    if err != nil {
        return fmt.Errorf("failed to read seed SQL file: %w", err)
    }

    _, err = testDB.Pool.Exec(ctx, string(sqlBytes))
    if err != nil {
        return fmt.Errorf("failed to seed test database: %w", err)
    }

    return nil
}


func CreateTestServer() *server.Server {
	s := &server.Server{
		Router: chi.NewRouter(),
		APIRouter: chi.NewRouter(),
		LobbyManager: core.NewLobbyManager(),
		DB: testDB, //this database is already mounted from main test setup
	}

  s.Router.Mount("/api", s.APIRouter)

  s.MountHandlers()
	return s
}

func CreateTestGomokuGameState() *gomoku.GomokuGameState {
	return &gomoku.GomokuGameState{
		GameID: "04c717b7-1234-4db6-afa1-c92d4afa9f0f",
		Players: []*core.Player{
			{PlayerID: "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0", PlayerName: "Alice", Color: "black"},
			{PlayerID: "88d0cd1e-912c-4d7f-9bc8-f9ef324d3df9", PlayerName: "Bob", Color: "white"},
		},
		Moves: []*gomoku.Move{
			{Row: 0, Col: 1, Color: "black"},
			{Row: 1, Col: 1, Color: "white"},
		},
	}
}

func executeRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.APIRouter.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Failed to load .env:", err)
		os.Exit(1)
	}

	testDB = &db.Database{}
	err = StartTestDB(testDB)
	if err != nil {
		fmt.Println("Failed to start test DB:", err)
		os.Exit(1)
	}
	code := m.Run()

	os.Exit(code)
}

func TestAuthMiddleware(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.ContextKey("userID"))
		require.NotNil(t, userID)
		require.Equal(t, "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0", userID)

		resp := api.AuthResponse{
			Username: "test",
			UserID:   userID.(string),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	token, err := utils.GenerateAccessJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
	require.NoError(t, err)

	handler := middleware.AuthMiddleware(dummyHandler)
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: token,
	})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var actual api.AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	require.NoError(t, err)

	expected := api.AuthResponse{
		Username: "test",
		UserID: "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0",
	}

	require.Equal(t, expected, actual)
}

func TestCheckAuth(t *testing.T) {
	err := ResetTestDB(testDB)
	require.NoError(t, err)


	s := CreateTestServer()
	

	token, err := utils.GenerateAccessJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/check-auth", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name: "access_token",
		Value: token,
	})

	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var actual api.AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	require.NoError(t, err)

	expected := api.AuthResponse{
		Username: "testuser",
		UserID: "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0",
	}

	require.Equal(t, expected, actual)
}

func TestSignUp(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()
	reqBody := api.SignUpRequest{
		Email:    "helloworld@gmail.com",
		Password: "123456789",
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create HTTP POST request with JSON body
	req := httptest.NewRequest("POST", "/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := executeRequest(req, s)

	checkResponseCode(t, http.StatusCreated, rr.Code)

}

func TestLogIn(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()

	reqBody := api.LogInRequest{
		Email: "testuser@gmail.com",
		Password: "test1234",
	}

	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, rr.Code)
}

func TestLogOut(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()

	req := httptest.NewRequest("GET", "/logout", nil)
	
	token, err := utils.GenerateAccessJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: token,
	})

	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, rr.Code)

	// Check if the cookie is cleared
	cookie := rr.Result().Cookies()[0]
	require.Equal(t, "access_token", cookie.Name)
	require.Equal(t, "", cookie.Value)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestJWTRefresh(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()

	req := httptest.NewRequest("GET", "/refresh", nil)

    secret := os.Getenv("JWT_SECRET_KEY")

    expiredAccessClaims := &jwt.RegisteredClaims{
		Subject:   "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
	}

	expiredAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredAccessClaims).SignedString([]byte(secret))
	require.NoError(t, err)

	// Create a valid refresh token (e.g., expires in 7 days)
	refreshClaims := &jwt.RegisteredClaims{
		Subject:   "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}

	validRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: expiredAccessToken,
	})
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: validRefreshToken,
	})


	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, rr.Code)
}
func TestPostGame(t *testing.T) {
    require.NoError(t, ResetTestDB(testDB))
    s := CreateTestServer()

    // Create auth JWT
    token, err := utils.GenerateAccessJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
    require.NoError(t, err)

    // Prepare request body
    gameState := CreateTestGomokuGameState()
    gameStateBytes, _ := json.Marshal(gameState)

    reqBody, _ := json.Marshal(gomoku.GomokuClientRequest{
        Type: "save",
        Data: gameStateBytes,
    })

    // Call correct endpoint with API prefix
    req := httptest.NewRequest("POST", "/gomoku/game", bytes.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")

    // Add auth cookie
    req.AddCookie(&http.Cookie{
        Name:  "access_token",
        Value: token,
    })

    rr := executeRequest(req, s)
    checkResponseCode(t, http.StatusOK, rr.Code)
}

func TestGetGame(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()

	// Insert game directly through DB layer
	gameState := CreateTestGomokuGameState()
	err := gomokudb.InsertGame(s.DB, gameState)
	require.NoError(t, err)

	// Auth token
	token, err := utils.GenerateAccessJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
	require.NoError(t, err)

	url := fmt.Sprintf("/gomoku/game?gameID=%s", gameState.GameID)

	req := httptest.NewRequest("GET", url, nil)
	req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: token,
	})

	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, rr.Code)

	// Decode response
	var returnedState gomoku.GomokuGameState
	err = json.Unmarshal(rr.Body.Bytes(), &returnedState)
	require.NoError(t, err)

	// Assertions
	require.Equal(t, gameState.GameID, returnedState.GameID)
	require.Equal(t, len(gameState.Moves), len(returnedState.Moves))
	require.Equal(t, gameState.Players[0].PlayerID, returnedState.Players[0].PlayerID)
}

func TestGetGameEmpty(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()

	// Auth token
	token, err := utils.GenerateAccessJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
	require.NoError(t, err)

	url := fmt.Sprintf("/gomoku/game?gameID=%s", "935fc971-5363-4817-b8db-bcce7b56809b")

	req := httptest.NewRequest("GET", url, nil)
	req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: token,
	})

	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusNotFound, rr.Code)
}

func TestGetGames(t *testing.T) {
	require.NoError(t, ResetTestDB(testDB))
	s := CreateTestServer()

	// Insert games directly through DB layer
	playerID := "f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0"

	gameState1 := CreateTestGomokuGameState()
	err := gomokudb.InsertGame(s.DB, gameState1)
	require.NoError(t, err)
	
	gameState2 := CreateTestGomokuGameState()
	gameState2.GameID = "bbd217b7-1234-4db6-afa1-c92d4afa9f0f"
	err = gomokudb.InsertGame(s.DB, gameState2)
	require.NoError(t, err)

	// Auth token
	token, err := utils.GenerateAccessJWT(playerID)
	require.NoError(t, err)

	url := fmt.Sprintf("/gomoku/games?playerID=%s", playerID)

	req := httptest.NewRequest("GET", url, nil)
	req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: token,
	})

	rr := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, rr.Code)

	// Decode response
	var returnedStates []*gomoku.GomokuGameState
	err = json.Unmarshal(rr.Body.Bytes(), &returnedStates)
	require.NoError(t, err)

	// Assertions	
	require.Equal(t, 2, len(returnedStates))
}