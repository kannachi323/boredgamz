package server

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestSPAHandlerServesAssetsAndRouteFallback(t *testing.T) {
	files := fstest.MapFS{
		"index.html":             &fstest.MapFile{Data: []byte("<main>app</main>")},
		"assets/app-deadbeef.js": &fstest.MapFile{Data: []byte("console.log('ok')")},
	}
	fileServer := http.FileServer(http.FS(fs.FS(files)))
	handler := spaHandler(files, fileServer)

	asset := httptest.NewRecorder()
	handler.ServeHTTP(asset, httptest.NewRequest(http.MethodGet, "/assets/app-deadbeef.js", nil))
	require.Equal(t, http.StatusOK, asset.Code)
	require.Contains(t, asset.Body.String(), "console.log")
	require.Equal(t, "public, max-age=31536000, immutable", asset.Header().Get("Cache-Control"))

	route := httptest.NewRecorder()
	handler.ServeHTTP(route, httptest.NewRequest(http.MethodGet, "/games/chess", nil))
	require.Equal(t, http.StatusOK, route.Code)
	require.Equal(t, "<main>app</main>", route.Body.String())
	require.Equal(t, "no-cache", route.Header().Get("Cache-Control"))
}

func TestSPAHandlerRejectsWrites(t *testing.T) {
	files := fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte("app")}}
	handler := spaHandler(files, http.FileServer(http.FS(files)))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/games", nil))
	require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}
