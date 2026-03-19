package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Custom claims with guest flag
type CustomClaims struct {
	IsGuest bool `json:"is_guest"`
	jwt.RegisteredClaims
}

func GenerateAccessJWT(userID string) (string, error) {
    secret := os.Getenv("JWT_SECRET_KEY")

    
    claims := &jwt.RegisteredClaims{
        Subject:   userID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    return token.SignedString([]byte(secret))
}

func GenerateGuestJWT() (string, string, error) {
    secret := os.Getenv("JWT_SECRET_KEY")

    randomBytes := make([]byte, 6)
    if _, err := rand.Read(randomBytes); err != nil {
        return "", "", err
    }

    guestID := "guest_" + hex.EncodeToString(randomBytes)

    claims := &CustomClaims{
        IsGuest: true,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   guestID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    signedToken, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", "", err
    }

    return guestID, signedToken, nil
}

func GenerateRefreshJWT(userID string) (string, error) {
    secret := os.Getenv("JWT_SECRET_KEY")

    claims := &jwt.RegisteredClaims{
        Subject: userID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(secret))
}

func VerifyJWT(tokenStr string) (string, error) {
    secret := os.Getenv("JWT_SECRET_KEY")


    token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil || !token.Valid {
        return "", errors.New("invalid token")
    }

    claims := token.Claims.(*jwt.RegisteredClaims)
    return claims.Subject, nil
}

func VerifyGuestOrUserJWT(tokenStr string) (string, bool, error) {
    secret := os.Getenv("JWT_SECRET_KEY")
    
    // Try custom claims first (for guest tokens)
    tokenCustom, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err == nil && tokenCustom.Valid {
        claims := tokenCustom.Claims.(*CustomClaims)
        return claims.Subject, claims.IsGuest, nil
    }
    
    // Fall back to standard claims (for user tokens)
    tokenStandard, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err != nil || !tokenStandard.Valid {
        return "", false, errors.New("invalid token")
    }
    
    claims := tokenStandard.Claims.(*jwt.RegisteredClaims)
    return claims.Subject, false, nil
}