package sws

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedFileName := r.URL.Path

	isRequestingDirectory := strings.HasSuffix(requestedFileName, "/")
	if isRequestingDirectory {
		requestedFileName += s.indexFileName
	}

	file, err := s.httpFileSystem.Open(requestedFileName)
	if err != nil {
		s.notFoundHandler.ServeHTTP(w, r)

		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	// Redirect to add forward slash
	if fileInfo.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusTemporaryRedirect)

			return
		}
	}

	s.serveFile(w, r, file, fileInfo, http.StatusOK)
}

func (s *Service) serveFile(
	w http.ResponseWriter,
	r *http.Request,
	file http.File,
	fileInfo fs.FileInfo,
	statusCode int,
) {
	fileExtension := filepath.Ext(fileInfo.Name())
	fileTypeHeader := mime.TypeByExtension(fileExtension)
	w.Header().Set("Content-Type", fileTypeHeader)

	for _, requestHook := range s.requestHooks {
		requestHook(w, r, file, fileInfo)
	}

	w.WriteHeader(statusCode)

	_, _ = io.Copy(w, file)
}
