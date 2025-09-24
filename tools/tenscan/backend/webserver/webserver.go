package webserver

import (
	"context"
	"net/http"
	"time"

	"github.com/ten-protocol/go-ten/tools/tenscan/backend"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type WebServer struct {
	engine      *gin.Engine
	backend     *backend.Backend
	bindAddress string
	logger      log.Logger
	server      *http.Server
}

func New(backend *backend.Backend, bindAddress string, logger log.Logger) *WebServer {
	r := gin.New()
	r.RedirectTrailingSlash = false
	gin.SetMode(gin.ReleaseMode)

	// todo this should be reviewed as anyone can access the api right now
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	r.Use(cors.New(config))

	server := &WebServer{
		engine:      r,
		backend:     backend,
		bindAddress: bindAddress,
		logger:      logger,
	}

	// routes
	routeItems(r, server)
	routeCounts(r, server)

	// todo group/format these into items, counts, actions
	r.GET("/health/", server.health)
	r.GET("/batchHeader/:hash", server.getBatchHeader)
	r.GET("/tx/:hash", server.getTransaction)
	// no longer available because the key is not public any more
	// r.POST("/actions/decryptTxBlob/", server.decryptTxBlob)

	return server
}

func (w *WebServer) Start() error {
	w.server = &http.Server{
		Addr:              w.bindAddress,
		Handler:           w.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling
	go func() {
		w.logger.Info("Starting server on " + w.bindAddress)
		println("Starting server on " + w.bindAddress)
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// todo don't panic
			panic(err)
		}
	}()

	return nil
}

func (w *WebServer) Stop() error {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return w.server.Shutdown(ctx)
}

func (w *WebServer) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"healthy": true})
}

/*func (w *WebServer) decryptTxBlob(c *gin.Context) {
	// Read the payload as a string
	payloadBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read payload"})
		return
	}

	payload := PostData{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	result, err := w.backend.DecryptTxBlob(payload.StrData)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}*/

type PostData struct {
	StrData string `json:"strData"`
}

func errorHandler(c *gin.Context, err error, logger log.Logger) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
	logger.Error(err.Error())
}
