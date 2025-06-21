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

	"playgomoku/backend/api"
	"playgomoku/backend/db"
	"playgomoku/backend/manager"
	"playgomoku/backend/middleware"
	"playgomoku/backend/server"
	"playgomoku/backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var testDB *db.Database

func StartTestDB(testDB *db.Database) error {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL_TEST")

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
	_, err := testDB.Pool.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("failed to truncate test database: %w", err)
	}
	
	sqlBytes, err := os.ReadFile("../db/sql/data.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	sql := string(sqlBytes)

	_, err = testDB.Pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to seed test database: %w", err)
	}
	return nil
}

func CreateTestServer() *server.Server {
	s := &server.Server{
		Router: chi.NewRouter(),
		LobbyManager: manager.NewLobbyManager(),
		DB: testDB, //this database is already mounted from main test setup
	}
	s.MountHandlers()

	return s
}

func executeRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

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
	defer CloseTestDB(testDB)

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

	token, err := utils.GenerateJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
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
	

	token, err := utils.GenerateJWT("f06d11d2-e147-45b7-aa29-c2aa5d8e9cc0")
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

	
	