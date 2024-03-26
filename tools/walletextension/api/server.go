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

func StaticFilesHandler() http.Handler {
	// Serves the web assets for the management of viewing keys.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	return http.FileServer(http.FS(noPrefixStaticFiles))
}
