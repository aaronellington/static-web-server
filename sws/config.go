package sws

import (
	"io/fs"
	"net/http"
	"regexp"
	"strings"
)

type configOption func(s *Service)

func EnableSPAMode(defaultFile string) configOption {
	return func(s *Service) {
		SetNotFoundFile(defaultFile, http.StatusOK)(s)
		SetHashedFilenameCachePolicy()(s)
	}
}

func SetHashedFilenameCachePolicy() configOption {
	return func(s *Service) {
		hashChecker := regexp.MustCompile(`(?m)\.[0-9a-z]{8,}\.`)

		s.requestHooks = append(s.requestHooks, func(w http.ResponseWriter, r *http.Request, file http.File, fileInfo fs.FileInfo) {
			fileName := strings.ToLower(fileInfo.Name())

			if hashChecker.MatchString(fileName) {
				w.Header().Set("Cache-Control", "public,max-age=31536000,immutable")
			} else {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			}
		})
	}
}

func SetNotFoundFile(notFoundFile string, statusCode int) configOption {
	return func(s *Service) {
		s.notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file, err := s.httpFileSystem.Open(notFoundFile)
			if err != nil {
				http.NotFound(w, r)
				return
			}

			fileStats, err := file.Stat()
			if err != nil {
				http.NotFound(w, r)
				return
			}

			s.serveFile(w, r, file, fileStats, statusCode)
		})
	}
}
