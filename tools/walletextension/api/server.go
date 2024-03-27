package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed all:static
var staticFiles embed.FS

const (
	staticDir = "static"
)

func StaticFilesHandler(prefix string) http.Handler {
	// Serves the web assets for the management of viewing keys.
	fileSystem, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(fileSystem)))
}
