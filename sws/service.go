package sws

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func New(
	httpFileSystem http.FileSystem,
	opts ...ConfigOption,
) *Service {
	c := &config{
		httpFileSystem:  httpFileSystem,
		notFoundHandler: http.NotFoundHandler(),
		indexFileName:   "index.html",
	}

	for _, opt := range opts {
		opt(c)
	}

	return &Service{
		config: c,
	}
}

type Service struct {
	config *config
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedFileName := r.URL.Path

	isRequestingDirectory := strings.HasSuffix(requestedFileName, "/")
	if s.config.indexFileName != "" && isRequestingDirectory {
		requestedFileName += s.config.indexFileName
	}

	file, err := s.config.httpFileSystem.Open(requestedFileName)
	if err != nil {
		s.config.notFoundHandler.ServeHTTP(w, r)

		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	// Redirect to add forward slash
	if fileInfo.IsDir() {
		if s.config.indexFileName == "" {
			s.config.notFoundHandler.ServeHTTP(w, r)

			return
		}

		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusTemporaryRedirect)

			return
		}
	}

	s.config.serveFile(w, r, file, fileInfo, http.StatusOK)
}

func (c *config) serveFile(
	w http.ResponseWriter,
	r *http.Request,
	file http.File,
	fileInfo fs.FileInfo,
	statusCode int,
) {
	// Content-Type Handling
	fileExtension := strings.ToLower(filepath.Ext(fileInfo.Name()))
	fileTypeHeader := mime.TypeByExtension(fileExtension)
	w.Header().Set("Content-Type", fileTypeHeader)

	// Cache-Control Handling - Default to no-cache
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	// Request Hook Handling
	for _, requestHook := range c.requestHooks {
		requestHook(w, r, file, fileInfo)
	}

	// Actually write to the response
	w.WriteHeader(statusCode)
	_, _ = io.Copy(w, file)
}
