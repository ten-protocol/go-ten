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
	}

	return results
}

func searchByNumber(db HostDB, number string) []*common.SearchResult {
	var results []*common.SearchResult
	num, _ := strconv.ParseInt(number, 10, 64)

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
