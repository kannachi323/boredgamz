package server

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Server) MountStatic() {
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}
	info, err := os.Stat(filepath.Join(staticDir, "index.html"))
	if err != nil || info.IsDir() {
		return
	}

	root := os.DirFS(staticDir)
	files := http.FileServer(http.FS(root))
	s.Router.Handle("/*", spaHandler(root, files))
}

func spaHandler(root fs.FS, files http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		path := strings.TrimPrefix(filepath.Clean(r.URL.Path), string(filepath.Separator))
		if path != "." {
			if info, err := fs.Stat(root, path); err == nil && !info.IsDir() {
				if strings.HasPrefix(path, "assets/") {
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				}
				files.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("Cache-Control", "no-cache")
		r.URL.Path = "/"
		files.ServeHTTP(w, r)
	})
}
