package enclavedb

import (
	"fmt"
	"sync"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/jmoiron/sqlx"
)

// PreparedStatementCache provides thread-safe caching of prepared SQL statements
type PreparedStatementCache struct {
	cache  map[string]*sqlx.Stmt
	mutex  sync.RWMutex
	db     *sqlx.DB
	logger gethlog.Logger
}

// NewStatementCache creates a new prepared statement cache
func NewStatementCache(db *sqlx.DB, logger gethlog.Logger) *PreparedStatementCache {
	return &PreparedStatementCache{
		cache:  make(map[string]*sqlx.Stmt),
		db:     db,
		logger: logger,
	}
}

// GetOrPrepare returns a cached prepared statement or creates and caches a new one
func (sc *PreparedStatementCache) GetOrPrepare(query string) (*sqlx.Stmt, error) {
	// First try to get from cache using read lock
	sc.mutex.RLock()
	stmt, found := sc.cache[query]
	sc.mutex.RUnlock()

	if found {
		sc.logger.Debug("Using cached prepared statement", "query", query)
		return stmt, nil
	}

	// Not found, prepare and cache using write lock
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Check again in case another goroutine prepared it while we were waiting
	stmt, found = sc.cache[query]
	if found {
		return stmt, nil
	}

	// Prepare new statement
	var err error
	stmt, err = sc.db.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	// Cache it
	sc.cache[query] = stmt
	sc.logger.Debug("Cached new prepared statement", "query", query)
	return stmt, nil
}

// Clear removes all cached statements, closing them properly
func (sc *PreparedStatementCache) Clear() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	var lastErr error
	for key, stmt := range sc.cache {
		if err := stmt.Close(); err != nil {
			sc.logger.Error("Error closing prepared statement", "query", key, "error", err)
			lastErr = err
		}
	}

	sc.cache = make(map[string]*sqlx.Stmt)
	return lastErr
}
