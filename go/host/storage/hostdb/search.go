package hostdb

import (
	"database/sql"
	"errors"
	"math/big"
	"strconv"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
)

// Search queries the host DB using the provided string. The input type is determined by length and will attempt to query
// by hash, sequence number or height.
func Search(db HostDB, query string) (*common.SearchResponse, error) {
	inputType := identifyInputType(query)

	var results []*common.SearchResult

	switch inputType {
	case "hash":
		// could be rollup, batch, or transaction hash
		results = append(results, searchByHash(db, query)...)
	case "number":
		// could be batch height or sequence
		results = append(results, searchByNumber(db, query)...)
	default:
		// try all searches
		results = append(results, searchByHash(db, query)...)
		results = append(results, searchByNumber(db, query)...)
	}

	searchResults := make([]common.SearchResult, len(results))
	for i, result := range results {
		searchResults[i] = *result
	}

	return &common.SearchResponse{
		ResultsData: searchResults,
		Total:       uint64(len(results)),
	}, nil
}

func identifyInputType(query string) string {
	query = strings.TrimPrefix(query, "0x")

	if len(query) == 64 {
		return "hash"
	}

	if len(query) == 40 {
		return "address"
	}

	if _, err := strconv.ParseInt(query, 10, 64); err == nil {
		return "number"
	}

	return "unknown"
}

func searchByHash(db HostDB, hash string) []*common.SearchResult {
	var results []*common.SearchResult
	logger := db.Logger()

	trimmedHash := strings.TrimPrefix(hash, "0x")

	rollup, err := GetRollupByHash(db, gethcommon.HexToHash(hash))
	if err == nil {
		results = append(results, &common.SearchResult{
			Type:      "rollup",
			Hash:      trimmedHash,
			Timestamp: rollup.Timestamp,
			ExtraData: map[string]interface{}{
				"rollup": rollup,
			},
		})
		return results // early return if rollup found due to uniqueness of hash lookup
	} else if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("No rollup found for hash during search", "hash", hash, "error", err)
	}

	batch, err := GetPublicBatch(db, gethcommon.HexToHash(hash))
	if err == nil {
		results = append(results, &common.SearchResult{
			Type:      "batch",
			Hash:      trimmedHash,
			Height:    batch.Height,
			Sequence:  batch.SequencerOrderNo,
			Timestamp: batch.Header.Time,
			ExtraData: map[string]interface{}{
				"batch": batch,
			},
		})
		return results // early return if rollup found due to uniqueness of hash lookup
	} else if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("No batch found for hash during search", "hash", hash, "error", err)
	}

	tx, err := GetTransaction(db, gethcommon.HexToHash(hash))
	if err == nil {
		results = append(results, &common.SearchResult{
			Type:      "transaction",
			Hash:      trimmedHash,
			Timestamp: tx.BatchTimestamp,
			ExtraData: map[string]interface{}{
				"transaction": tx,
			},
		})
	} else if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("No tx found for hash during search", "hash", hash, "error", err)
	}

	return results
}

func searchByNumber(db HostDB, number string) []*common.SearchResult {
	var results []*common.SearchResult
	logger := db.Logger()

	num, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		logger.Error("Failed to parse number for search", "number", number, "error", err)
		return results
	}

	batch, err := GetBatchByHeight(db, big.NewInt(num))
	if err == nil {
		results = append(results, &common.SearchResult{
			Type:      "batch",
			Hash:      batch.FullHash.Hex(),
			Height:    batch.Height,
			Sequence:  batch.SequencerOrderNo,
			Timestamp: batch.Header.Time,
			ExtraData: map[string]interface{}{
				"batch": batch,
			},
		})
	} else if errors.Is(err, sql.ErrNoRows) {
		logger.Error("No batch found for height during search", "height", num, "error", err)
	}

	batch, err = GetPublicBatchBySequenceNumber(db, uint64(num))
	if err == nil {
		results = append(results, &common.SearchResult{
			Type:      "batch",
			Hash:      batch.FullHash.Hex(),
			Height:    batch.Height,
			Sequence:  batch.SequencerOrderNo,
			Timestamp: batch.Header.Time,
			ExtraData: map[string]interface{}{
				"batch": batch,
			},
		})
	} else if errors.Is(err, sql.ErrNoRows) {
		logger.Error("No batch found for sequence number during search", "seq", num, "error", err)
	}

	return results
}
