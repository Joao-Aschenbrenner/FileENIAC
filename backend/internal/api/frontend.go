// SPDX-License-Identifier: MIT
package api

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ENIACSystems/FileENIAC/backend/webui"
)

var frontendFS fs.FS

func init() {
	if os.Getenv("ENIAC_NO_FRONTEND") != "1" {
		frontendFS = webui.FS()
	}
}

func (s *Server) ServeFrontend() {
	if frontendFS == nil {
		return
	}
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if strings.HasPrefix(p, "/api") {
			http.NotFound(w, r)
			return
		}
		name := strings.TrimPrefix(p, "/")
		if name == "" {
			name = "index.html"
		}
		f, err := frontendFS.Open(name)
		if err != nil {
			name = "index.html"
			f, err = frontendFS.Open(name)
			if err != nil {
				http.NotFound(w, r)
				return
			}
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			http.NotFound(w, r)
			return
		}
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(f); err != nil {
			http.NotFound(w, r)
			return
		}
		ct := mimeType(path.Ext(name))
		if ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		http.ServeContent(w, r, name, stat.ModTime(), bytes.NewReader(buf.Bytes()))
	})
}

func mimeType(ext string) string {
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".json":
		return "application/json"
	case ".woff2":
		return "font/woff2"
	}
	return ""
}
