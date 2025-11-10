package static

import (
	"embed"
	"io"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

// Handler serves index.html at "/" and static assets at "/assets/*", forbids directory listing
func EmbedStaticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var filePath string
		cleanPath := path.Clean(r.URL.Path)
		if cleanPath == "/" || cleanPath == "/entries" || cleanPath == "/sets" {
			filePath = "dist/index.html"
		} else {
			filePath = "dist" + cleanPath
		}

		// Запретить листинг директорий
		if strings.HasSuffix(filePath, "/") {
			http.NotFound(w, r)
			return
		}

		file, err := distFS.Open(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		ext := path.Ext(filePath)
		switch ext {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		readSeeker, ok := file.(io.ReadSeeker)
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.ServeContent(w, r, filePath, info.ModTime(), readSeeker)
	})
}
