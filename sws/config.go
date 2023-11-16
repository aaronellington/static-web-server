package sws

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type ConfigOption func(c *config)

type config struct {
	httpFileSystem  http.FileSystem
	notFoundHandler http.Handler
	indexFileName   string
	requestHooks    []requestHook
}

type requestHook func(
	w http.ResponseWriter,
	r *http.Request,
	file http.File,
	fileInfo fs.FileInfo,
)

func SetNotFoundHandler(notFoundHandler http.Handler) ConfigOption {
	return func(c *config) {
		c.notFoundHandler = notFoundHandler
	}
}

func SetNotFoundFile(fileName string, statusCode int) ConfigOption {
	return func(c *config) {
		SetNotFoundHandler(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					file, err := c.httpFileSystem.Open(fileName)
					if err != nil {
						http.NotFound(w, r)

						return
					}
					defer file.Close()

					fileInfo, err := file.Stat()
					if err != nil {
						http.NotFound(w, r)

						return
					}

					c.serveFile(w, r, file, fileInfo, statusCode)
				},
			),
		)(c)
	}
}

func SetIndexFileName(indexFileName string) ConfigOption {
	return func(c *config) {
		c.indexFileName = indexFileName
	}
}

func SetCSP(csp ContentSecurityPolicy) ConfigOption {
	return func(c *config) {
		cspString := csp.String()
		if cspString == "" {
			return
		}

		c.requestHooks = append(c.requestHooks, func(w http.ResponseWriter, r *http.Request, file http.File, fileInfo fs.FileInfo) {
			fileExtension := strings.ToLower(filepath.Ext(fileInfo.Name()))

			if !strings.HasPrefix(fileExtension, ".htm") {
				return
			}

			if cspString := cspString; cspString == "" {
				return
			}

			w.Header().Set("Content-Security-Policy", cspString)
		})
	}
}

func SetHashedFilenameCachePolicy() ConfigOption {
	return func(c *config) {
		c.requestHooks = append(
			c.requestHooks,
			func(w http.ResponseWriter, r *http.Request, file http.File, fileInfo fs.FileInfo) {
				if strings.HasPrefix(r.URL.Path, "/assets/") {
					w.Header().Set("Cache-Control", "public,max-age=31536000,immutable")
					w.Header().Del("Pragma")
					w.Header().Del("Expires")
				}
			},
		)
	}
}
