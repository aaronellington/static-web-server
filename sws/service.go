package sws

import (
	"io/fs"
	"net/http"
)

func New(
	httpFileSystem http.FileSystem,
	opts ...configOption,
) *Service {
	service := &Service{
		httpFileSystem:  httpFileSystem,
		notFoundHandler: http.NotFoundHandler(),
		indexFileName:   "index.html",
	}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

type RequestHook func(
	w http.ResponseWriter,
	r *http.Request,
	file http.File,
	fileInfo fs.FileInfo,
)

type Service struct {
	httpFileSystem  http.FileSystem
	notFoundHandler http.Handler
	indexFileName   string
	requestHooks    []RequestHook
}
