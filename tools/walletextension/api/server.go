package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
)

//go:embed static
//go:embed staticOG
var staticFiles embed.FS

const staticDir = "static"

// Server is a wrapper for the http server
type Server struct {
	server *http.Server
}

// Start starts the server in its own goroutine and returns an error chan where errors can be monitored
func (s *Server) Start() chan error {
	errChan := make(chan error)
	go func() {
		// start the server and serve any errors over the channel
		errChan <- s.server.ListenAndServe()
	}()
	return errChan
}

// Stop synchronously stops the server
func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

// NewHTTPServer returns the HTTP server for the WE
func NewHTTPServer(address string, routes []Route) *Server {
	return &Server{
		server: createHTTPServer(address, routes),
	}
}

// NewWSServer returns the WS server for the WE
func NewWSServer(address string, routes []Route) *Server {
	return &Server{
		server: createWSServer(address, routes),
	}
}

func createHTTPServer(address string, routes []Route) *http.Server {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	for _, route := range routes {
		serveMux.HandleFunc(route.Name, route.Func)
	}

	// Serves the web assets for the management of viewing keys.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMux.Handle(common.PathViewingKeys, http.StripPrefix(common.PathViewingKeys, http.FileServer(http.FS(noPrefixStaticFiles))))

	noPrefixStaticFilesOG, err := fs.Sub(staticFiles, "staticOG")
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMux.Handle(common.PathObscuroGateway, http.StripPrefix(common.PathObscuroGateway, http.FileServer(http.FS(noPrefixStaticFilesOG))))

	// Creates the actual http server with a ReadHeaderTimeout to avoid Potential Slowloris Attack
	server := &http.Server{Addr: address, Handler: serveMux, ReadHeaderTimeout: common.ReaderHeadTimeout}
	return server
}

func createWSServer(address string, routes []Route) *http.Server {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	for _, route := range routes {
		serveMux.HandleFunc(route.Name, route.Func)
	}

	return &http.Server{Addr: address, Handler: serveMux, ReadHeaderTimeout: 10 * time.Second}
}
