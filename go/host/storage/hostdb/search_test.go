package hostdb

import (
	"strings"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

func TestIdentifyInputType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", "hash"},
		{"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", "hash"},
		{"12345", "number"},
		{"0", "number"},
		{"invalid", "unknown"},
		{"", "unknown"},
	}

	for _, test := range tests {
		result := identifyInputType(test.input)
		if result != test.expected {
			t.Errorf("identifyInputType(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestSearchByRollupAndBatchHash(t *testing.T) {
	db, err := CreateSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	metadata := createRollupMetadata()
	rollup := createRollup(batchNumber-10, batchNumber)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup, &common.ExtRollupMetadata{}, &metadata, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	dbtx.Write()

	// Create and store a batch
	txHashes := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne"))}
	batch := CreateBatch(batchNumber, txHashes)
	dbtx, _ = db.NewDBTransaction()
	err = AddBatch(dbtx, db.GetSQLStatement(), &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	// Test search by rollup hash
	rollupHash := rollup.Header.Hash().Hex()
	searchResponse, err := Search(db, rollupHash)
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 1 {
		t.Errorf("expected 1 result, got %d", searchResponse.Total)
	}

	if len(searchResponse.ResultsData) != 1 {
		t.Errorf("expected 1 result in ResultsData, got %d", len(searchResponse.ResultsData))
	}

	result := searchResponse.ResultsData[0]
	if result.Type != "rollup" {
		t.Errorf("expected type 'rollup', got '%s'", result.Type)
	}

	trimmedRollupHash := strings.TrimPrefix(rollupHash, "0x")
	if result.Hash != trimmedRollupHash {
		t.Errorf("expected hash %s, got %s", trimmedRollupHash, result.Hash)
	}

	// Test search by batch hash
	batchHash := batch.Header.Hash().Hex()
	searchResponse, err = Search(db, batchHash)
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 1 {
		t.Errorf("expected 1 result, got %d", searchResponse.Total)
	}

	result = searchResponse.ResultsData[0]
	if result.Type != "batch" {
		t.Errorf("expected type 'batch', got '%s'", result.Type)
	}

	trimmedBatchHash := strings.TrimPrefix(batchHash, "0x")
	if result.Hash != trimmedBatchHash {
		t.Errorf("expected hash %s, got %s", trimmedBatchHash, result.Hash)
	}
}

func TestSearchByTransactionHash(t *testing.T) {
	db, err := CreateSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	// Create and store a batch with a transaction
	txHash := gethcommon.BytesToHash([]byte("magicStringOne"))
	txHashes := []common.L2TxHash{txHash}
	batch := CreateBatch(batchNumber, txHashes)
	dbtx, _ := db.NewDBTransaction()
	err = AddBatch(dbtx, db.GetSQLStatement(), &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	// Test search by transaction hash
	txHashStr := txHash.Hex()
	searchResponse, err := Search(db, txHashStr)
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 1 {
		t.Errorf("expected 1 result, got %d", searchResponse.Total)
	}

	if len(searchResponse.ResultsData) != 1 {
		t.Errorf("expected 1 result in ResultsData, got %d", len(searchResponse.ResultsData))
	}

	result := searchResponse.ResultsData[0]
	if result.Type != "transaction" {
		t.Errorf("expected type 'transaction', got '%s'", result.Type)
	}

	trimmedTxHash := strings.TrimPrefix(txHashStr, "0x")

	if result.Hash != trimmedTxHash {
		t.Errorf("expected hash %s, got %s", trimmedTxHash, result.Hash)
	}

	// Verify the transaction data is in ExtraData
	if result.ExtraData == nil {
		t.Errorf("expected ExtraData to contain transaction information")
	}
}

func TestSearchByNumber(t *testing.T) {
	db, err := CreateSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	// Create and store a batch
	txHashes := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne"))}
	batch := CreateBatchWithDiffHeight(batchNumber, batchNumber+100, txHashes)
	dbtx, _ := db.NewDBTransaction()
	err = AddBatch(dbtx, db.GetSQLStatement(), &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	// search by height
	heightStr := batch.Header.Number.String()
	searchResponse, err := Search(db, heightStr)
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 1 {
		t.Errorf("expected 2 result, got %d", searchResponse.Total)
	}

	result := searchResponse.ResultsData[0]
	if result.Type != "batch" {
		t.Errorf("expected type 'batch', got '%s'", result.Type)
	}

	if result.Height.Cmp(batch.Header.Number) != 0 {
		t.Errorf("expected height %s, got %s", batch.Header.Number.String(), result.Height.String())
	}

	// search by sequence
	seqStr := batch.SeqNo().String()
	searchResponse, err = Search(db, seqStr)
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 1 {
		t.Errorf("expected 1 result, got %d", searchResponse.Total)
	}

	result = searchResponse.ResultsData[0]
	if result.Type != "batch" {
		t.Errorf("expected type 'batch', got '%s'", result.Type)
	}

	if result.Sequence.Cmp(batch.SeqNo()) != 0 {
		t.Errorf("expected sequence %s, got %s", batch.SeqNo().String(), result.Sequence.String())
	}
}

func TestAmbiguousSearch(t *testing.T) {
	db, err := CreateSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	// Create and store a batch with height 100
	txHashes := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne"))}
	batch := CreateBatchWithDiffHeight(batchNumber, 100, txHashes)
	dbtx, _ := db.NewDBTransaction()
	err = AddBatch(dbtx, db.GetSQLStatement(), &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	// Test search with a number that could match both height and sequence
	// In this case, we're searching for "100" which could be both height and sequence
	searchResponse, err := Search(db, "100")
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	// Should find at least 1 result (the batch by height)
	if searchResponse.Total < 1 {
		t.Errorf("expected at least 1 result, got %d", searchResponse.Total)
	}

	// Verify we have a batch result
	hasBatch := false
	for _, result := range searchResponse.ResultsData {
		if result.Type == "batch" {
			hasBatch = true
			break
		}
	}

	if !hasBatch {
		t.Errorf("expected to find batch in search results")
	}
}

func TestSearchEmptyResults(t *testing.T) {
	db, err := CreateSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	// Test search with non-existent hash
	searchResponse, err := Search(db, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 0 {
		t.Errorf("expected 0 results, got %d", searchResponse.Total)
	}

	if len(searchResponse.ResultsData) != 0 {
		t.Errorf("expected 0 results in ResultsData, got %d", len(searchResponse.ResultsData))
	}

	// Test search with non-existent number
	searchResponse, err = Search(db, "999999")
	if err != nil {
		t.Errorf("search failed: %s", err)
	}

	if searchResponse.Total != 0 {
		t.Errorf("expected 0 results, got %d", searchResponse.Total)
	}
}
