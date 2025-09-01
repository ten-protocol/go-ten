package webserver

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func routeItems(r *gin.Engine, server *WebServer) {
	// info
	r.GET("/info/obscuro/", server.getConfig)
	r.GET("/info/health/", server.getHealthStatus)

	// batches
	r.GET("/items/batches/", server.getBatchListing)
	r.GET("/items/batch/latest/", server.getLatestBatch)
	r.GET("/items/batch/:hash", server.getBatch)
	r.GET("/items/batch/:hash/transactions", server.getBatchTransactions)
	r.GET("/items/batch/height/:height", server.getBatchByHeight)
	r.GET("/items/batch/seq/:seq", server.getBatchBySeq)

	// rollups
	r.GET("/items/rollups/", server.getRollupListing) // New
	r.GET("/items/rollup/latest/", server.getLatestRollupHeader)
	r.GET("/items/rollup/:hash", server.getRollup)
	r.GET("/items/rollup/:hash/batches", server.getRollupBatches)
	r.GET("/items/rollup/batch/:seq", server.getRollupBySeq)

	// transactions
	r.GET("/items/transactions/", server.getPublicTransactions)
	r.GET("/items/transaction/:hash", server.getTransaction)
	r.GET("/items/transactions/count", server.getTotalTxCount)
	r.GET("/items/blocks/", server.getBlockListing)

	// search
	r.GET("/items/search/", server.search)
}

func (w *WebServer) getHealthStatus(c *gin.Context) {
	healthStatus, err := w.backend.GetTenNodeHealthStatus()

	// TODO: error handling, since this does not easily tell connection errors from health errors
	c.JSON(http.StatusOK, gin.H{"result": healthStatus, "errors": fmt.Sprintf("%s", err)})
}

func (w *WebServer) getLatestBatch(c *gin.Context) {
	batch, err := w.backend.GetLatestBatch()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getLatestBatch request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getLatestRollupHeader(c *gin.Context) {
	rollup, err := w.backend.GetLatestRollupHeader()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getLatestRollupHeader request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": rollup})
}

func (w *WebServer) getBatch(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetBatchByHash(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBatch request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getBatchByHeight(c *gin.Context) {
	heightStr := c.Param("height")

	heightBigInt := new(big.Int)
	heightBigInt.SetString(heightStr, 10)
	batch, err := w.backend.GetBatchByHeight(heightBigInt)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBatchByHeight request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getBatchBySeq(c *gin.Context) {
	seqStr := c.Param("seq")

	seqBigInt := new(big.Int)
	seqBigInt.SetString(seqStr, 10)
	batch, err := w.backend.GetBatchBySeq(seqBigInt)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBatchByHeight request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getRollupBySeq(c *gin.Context) {
	seqNo := c.Param("seq")

	seq, err := strconv.ParseUint(seqNo, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse sequence number: %w", err), w.logger)
		return
	}

	batch, err := w.backend.GetRollupBySeqNo(seq)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getRollupBySeq request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getBatchHeader(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetBatchHeader(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBatchHeader request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getTransaction(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetTransaction(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getTransaction request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getPublicTransactions(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getPublicTransactions offset units %w", err), w.logger)
		return
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getPublicTransactions size units %w", err), w.logger)
		return
	}

	publicTxs, err := w.backend.GetPublicTransactions(offset, size)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getPublicTransactions request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": publicTxs})
}

func (w *WebServer) getTotalTxCount(c *gin.Context) {
	txCount, err := w.backend.GetTotalTransactionCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getTotalTxCount request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": txCount})
}

func (w *WebServer) getBatchListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getBatchListing offset units %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getBatchListing size units %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBatchesListing(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBatchListing request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func (w *WebServer) getRollupListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getRollupListing offset units %w", err), w.logger)
		return
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getRollupListing size units %w", err), w.logger)
		return
	}

	rollupListing, err := w.backend.GetRollupListing(offset, size)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getRollupListing request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": rollupListing})
}

func (w *WebServer) getBlockListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getBlockListing offset units %w", err), w.logger)
		return
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getBlockListing size units %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBlockListing(offset, size)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBlockListing request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func (w *WebServer) getRollup(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	rollup, err := w.backend.GetRollupByHash(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getRollup request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": rollup})
}

func (w *WebServer) getRollupBatches(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)

	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getRollupBatches offset units %w", err), w.logger)
		return
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getRollupBatches size units %w", err), w.logger)
		return
	}

	rollupBatchesListing, err := w.backend.GetRollupBatches(parsedHash, offset, size)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getRollupBatches request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": rollupBatchesListing})
}

func (w *WebServer) getBatchTransactions(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)

	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getBatchTransactions offset units %w", err), w.logger)
		return
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to parse getBatchTransactions size units %w", err), w.logger)
		return
	}

	txListing, err := w.backend.GetBatchTransactions(parsedHash, offset, size)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getBatchTransactions request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": txListing})
}

func (w *WebServer) getConfig(c *gin.Context) {
	config, err := w.backend.GetConfig()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute getConfig request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": config})
}

func (w *WebServer) search(c *gin.Context) {
	query := c.Query("query")

	results, err := w.backend.Search(query)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute search request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": results})
}
