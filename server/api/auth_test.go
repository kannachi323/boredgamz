package api

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthJSONRejectsUnknownTrailingAndOversizedInput(t *testing.T) {
	for _, body := range []string{
		`{"email":"a@example.com","password":"long-enough-password","admin":true}`,
		`{"email":"a@example.com","password":"long-enough-password"}{}`,
		`{"email":"` + strings.Repeat("a", maxAuthBody) + `"}`,
	} {
		request := httptest.NewRequest("POST", "/api/signup", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		var target SignUpRequest
		require.Error(t, decodeAuthJSON(recorder, request, &target))
	}
}

func TestAuthCookiesAreSecureInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	recorder := httptest.NewRecorder()
	setAuthCookie(recorder, "access_token", "token", 60)
	cookies := recorder.Result().Cookies()
	require.Len(t, cookies, 1)
	require.True(t, cookies[0].Secure)
	require.True(t, cookies[0].HttpOnly)
	require.Equal(t, 60, cookies[0].MaxAge)
}

func TestEmailValidationRejectsDisplayNamesAndMalformedInput(t *testing.T) {
	require.True(t, validEmail("player@example.com"))
	require.False(t, validEmail("Player <player@example.com>"))
	require.False(t, validEmail("not-an-email"))
}

func TestUsernameValidation(t *testing.T) {
	require.True(t, validUsername("Board Sage"))
	require.False(t, validUsername("A"))
	require.False(t, validUsername(strings.Repeat("x", 25)))
	require.False(t, validUsername("bad\nname"))
}

func TestAuthJSONRequiresJSONContentType(t *testing.T) {
	body := `{"email":"a@example.com","password":"long-enough-password"}`
	request := httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	request.Header.Set("Content-Type", "text/plain")
	var target LogInRequest
	require.Error(t, decodeAuthJSON(httptest.NewRecorder(), request, &target))

	request = httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	require.NoError(t, decodeAuthJSON(httptest.NewRecorder(), request, &target))
}

func TestPasswordLimitsRespectBcrypt(t *testing.T) {
	require.False(t, validPassword("short"))
	require.True(t, validPassword("twelve-chars"))
	require.True(t, validPassword(strings.Repeat("x", 72)))
	require.False(t, validPassword(strings.Repeat("x", 73)))
	// 24 three-byte runes = 72 bytes is allowed; 25 is not.
	require.True(t, validPassword(strings.Repeat("密", 24)))
	require.False(t, validPassword(strings.Repeat("密", 25)))
}
