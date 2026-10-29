package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtIssuer = "boredgamz"
const jwtAudience = "boredgamz-web"

type tokenClaims struct {
	TokenUse string `json:"token_use"`
	jwt.RegisteredClaims
}

func GenerateAccessJWT(userID string) (string, error) {
	return generateJWT(userID, "access", time.Hour)
}

func GenerateRefreshJWT(userID string) (string, error) {
	return generateJWT(userID, "refresh", 30*24*time.Hour)
}

func generateJWT(userID, tokenUse string, duration time.Duration) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	claims := &tokenClaims{
		TokenUse: tokenUse,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    jwtIssuer,
			Audience:  jwt.ClaimStrings{jwtAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func VerifyJWT(token string) (string, error) {
	return verifyJWT(token, "access")
}

func VerifyRefreshJWT(token string) (string, error) {
	return verifyJWT(token, "refresh")
}

func verifyJWT(tokenString, expectedUse string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected token signing method")
			}
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(jwtIssuer),
		jwt.WithAudience(jwtAudience),
	)
	if err != nil || !token.Valid || claims.Subject == "" || claims.TokenUse != expectedUse {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func jwtSecret() ([]byte, error) {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET_KEY must contain at least 32 bytes")
	}
	return secret, nil
}
