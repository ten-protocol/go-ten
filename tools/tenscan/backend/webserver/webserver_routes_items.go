package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func routeItems(r *gin.Engine, server *WebServer) {
	r.GET("/items/batch/latest/", server.getLatestBatch)
	r.GET("/items/batch/:hash", server.getBatch)
	r.GET("/items/rollup/latest/", server.getLatestRollupHeader)
	r.GET("/items/rollups/", server.getRollupListing) // New
	r.GET("/items/batches/", server.getBatchListingDeprecated)
	r.GET("/items/new/batches/", server.getBatchListingNew)
	r.GET("/items/blocks/", server.getBlockListing) // Deprecated
	r.GET("/items/transactions/", server.getPublicTransactions)
	r.GET("/info/obscuro/", server.getConfig)
	r.POST("/info/health/", server.getHealthStatus)
}

func (w *WebServer) getHealthStatus(c *gin.Context) {
	healthStatus, err := w.backend.GetTenNodeHealthStatus()

	// TODO: error handling, since this does not easily tell connection errors from health errors
	c.JSON(http.StatusOK, gin.H{"result": healthStatus, "errors": fmt.Sprintf("%s", err)})
}

func (w *WebServer) getLatestBatch(c *gin.Context) {
	batch, err := w.backend.GetLatestBatch()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getLatestRollupHeader(c *gin.Context) {
	rollup, err := w.backend.GetLatestRollupHeader()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": rollup})
}

func (w *WebServer) getBatch(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetBatchByHash(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getBatchHeader(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetBatchHeader(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getTransaction(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetTransaction(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getPublicTransactions(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	publicTxs, err := w.backend.GetPublicTransactions(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": publicTxs})
}

func (w *WebServer) getBatchListingNew(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBatchesListing(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func (w *WebServer) getBatchListingDeprecated(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBatchesListingDeprecated(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func (w *WebServer) getRollupListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	rollupListing, err := w.backend.GetRollupListing(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": rollupListing})
}

func (w *WebServer) getBlockListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBlockListing(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func (w *WebServer) getConfig(c *gin.Context) {
	config, err := w.backend.GetConfig()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": config})
}
