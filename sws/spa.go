package sws

import (
	"net/http"
)

func NewSPA(
	httpFileSystem http.FileSystem,
	indexFileName string,
	csp ContentSecurityPolicy,
) *Service {
	return New(
		httpFileSystem,
		SetHashedFilenameCachePolicy(),
		SetCSP(csp),
		SetIndexFileName(""),
		SetNotFoundFile(indexFileName, http.StatusOK),
	)
}
