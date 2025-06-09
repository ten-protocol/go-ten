package hostdb

import (
	"math/big"
	"strconv"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
)

// todo (@will) add pagination if results are large
func Search(db HostDB, query string) (*common.SearchResponse, error) {
	// First try to identify what type of input we have
	inputType := identifyInputType(query)

	var results []*common.SearchResult

	switch inputType {
	case "hash":
		// Could be rollup, batch, or transaction hash
		results = append(results, searchByHash(db, query)...)
	case "number":
		// Could be batch height or sequence
		results = append(results, searchByNumber(db, query)...)
	default:
		// Try all searches
		results = append(results, searchByHash(db, query)...)
		results = append(results, searchByNumber(db, query)...)
	}

	// Convert pointers to values for the response
	searchResults := make([]common.SearchResult, len(results))
	for i, result := range results {
		searchResults[i] = *result
	}

	return &common.SearchResponse{
		TransactionsData: searchResults,
		Total:            uint64(len(results)),
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
	
	// Trim 0x prefix for consistency
	trimmedHash := strings.TrimPrefix(hash, "0x")

	// Search rollups
	rollup, err := GetRollupByHash(db, gethcommon.HexToHash(hash))
	if err == nil {
		println("ROLLUP found with hash: ", hash)
		results = append(results, &common.SearchResult{
			Type:      "rollup",
			Hash:      trimmedHash,
			Timestamp: rollup.Timestamp,
			ExtraData: map[string]interface{}{
				"rollup": rollup,
			},
		})
	}

	// Search batches
	batch, err := GetPublicBatch(db, gethcommon.HexToHash(hash))
	if err == nil {
		println("BATCH found with hash: ", hash)
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
	}

	// Search transactions
	tx, err := GetTransaction(db, gethcommon.HexToHash(hash))
	if err == nil {
		println("TX found with hash: ", hash)
		results = append(results, &common.SearchResult{
			Type:      "transaction",
			Hash:      trimmedHash,
			Timestamp: tx.BatchTimestamp,
			ExtraData: map[string]interface{}{
				"transaction": tx,
			},
		})
	}

	return results
}

func searchByNumber(db HostDB, number string) []*common.SearchResult {
	var results []*common.SearchResult
	num, _ := strconv.ParseInt(number, 10, 64)

	// Try as batch height
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
	}

	// Try as batch sequence
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
	}

	return results
}
