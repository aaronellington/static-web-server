package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aaronellington/environment-go/environment"
	"github.com/aaronellington/static-web-server/sws"
)

type Config struct {
	Port           int    `env:"PORT"`
	Host           string `env:"HOST"`
	FileSystemPath string `env:"FILE_SYSTEM_PATH"`
	NotFoundFile   string `env:"NOT_FOUND_FILE"`
	NotFoundStatus int    `env:"NOT_FOUND_STATUS"`
}

func main() {
	env := environment.New(true)

	config := Config{
		Port:           2828,
		Host:           "0.0.0.0",
		FileSystemPath: "./public",
		NotFoundFile:   "404.html",
		NotFoundStatus: http.StatusNotFound,
	}

	if err := env.Decode(&config); err != nil {
		log.Fatal(err)
	}

	s := sws.New(
		http.FS(os.DirFS(config.FileSystemPath)),
		sws.SetHashedFilenameCachePolicy(),
		sws.SetNotFoundFile(config.NotFoundFile, config.NotFoundStatus),
	)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("HTTP Request: %s %s", r.Method, r.URL.Path)
			s.ServeHTTP(w, r)
		}),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
