package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeCounts(r *gin.Engine, server *WebServer) {
	r.GET("/count/contracts/", server.getTotalContractCount)
	r.GET("/count/contracts/historical", server.getHistoricalContractCount)
	r.GET("/count/transactions/", server.getTotalTransactionCount)
	r.GET("/count/transactions/historical/", server.getHistoricalTransactionCount)
}

func (w *WebServer) getTotalContractCount(c *gin.Context) {
	count, err := w.backend.GetTotalContractCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (w *WebServer) getTotalTransactionCount(c *gin.Context) {
	count, err := w.backend.GetTotalTransactionCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (w *WebServer) getHistoricalTransactionCount(c *gin.Context) {
	count, err := w.backend.GetHistoricalTransactionCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"historicalCount": count})
}

func (w *WebServer) getHistoricalContractCount(c *gin.Context) {
	count, err := w.backend.GetHistoricalContractCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"historicalCount": count})
}
