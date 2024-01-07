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
	r.GET("/items/rollup/latest/", server.getLatestRollupHeader)
	r.GET("/items/batch/:hash", server.getBatch)
	r.GET("/items/transactions/", server.getPublicTransactions)
	r.GET("/items/batches/", server.getBatchListing)
	r.GET("/items/blocks/", server.getBlockListing)
	r.GET("/info/obscuro/", server.getConfig)
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
	block, err := w.backend.GetLatestRollupHeader()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": block})
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

func (w *WebServer) getBatchListing(c *gin.Context) {
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
