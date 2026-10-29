package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessAndRefreshTokensAreNotInterchangeable(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret-that-is-longer-than-thirty-two-bytes")
	access, err := GenerateAccessJWT("player-1")
	require.NoError(t, err)
	refresh, err := GenerateRefreshJWT("player-1")
	require.NoError(t, err)
	userID, err := VerifyJWT(access)
	require.NoError(t, err)
	require.Equal(t, "player-1", userID)
	userID, err = VerifyRefreshJWT(refresh)
	require.NoError(t, err)
	require.Equal(t, "player-1", userID)
	_, err = VerifyRefreshJWT(access)
	require.Error(t, err)
	_, err = VerifyJWT(refresh)
	require.Error(t, err)
}

func TestJWTRequiresStrongConfiguredSecret(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "short")
	_, err := GenerateAccessJWT("player-1")
	require.EqualError(t, err, "JWT_SECRET_KEY must contain at least 32 bytes")
}
