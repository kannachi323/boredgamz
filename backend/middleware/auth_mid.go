package middleware

import (
	"context"
	"net/http"
	"playgomoku/backend/utils"
)

type contextKey string

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

        ctx := context.WithValue(r.Context(), contextKey("userID"), userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
