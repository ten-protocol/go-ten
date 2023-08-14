package common

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
)

type PrivateQueryResponse struct {
	Receipts types.Receipts
	Total    uint64
}

type PublicTxListingResponse struct {
	PublicTxData []PublicTxData
	Total        uint64
}

type BatchListingResponse struct {
	BatchData []PublicBatchListing
	Total     uint64
}

type BlockListingResponse struct {
	BlockData []PublicBlockListing
	Total     uint64
}

type PublicTxData struct {
	TransactionHash TxHash
	BatchHeight     *big.Int
	Finality        FinalityType
}

type PublicBatchListing struct {
	BatchHeader
}

type PublicBlockListing struct {
	BlockHeader types.Header `json:"blockHeader"`
	RollupHash  common.Hash  `json:"rollupHash"`
}

// MarshalJSON custom marshals the RollupHeader into a json
func (p *PublicBlockListing) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		RollupHash string       `json:"rollupHash"`
		Header     types.Header `json:"blockHeader"`
	}{
		RollupHash: p.RollupHash.Hex(),
		Header:     p.BlockHeader,
	})
}

func (p *PublicBlockListing) UnmarshalJSON(data []byte) error {
	// Create an anonymous structure that matches the expected JSON structure
	type Alias struct {
		RollupHash string       `json:"rollupHash"`
		Header     types.Header `json:"blockHeader"`
	}

	// Temporary variable to hold the unmarshalled data
	var aux Alias

	// Unmarshal the JSON data into the temporary structure
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Transfer the data from the temporary structure to the actual TestEntity struct
	p.RollupHash = common.HexToHash(aux.RollupHash)
	p.BlockHeader = aux.Header

	return nil
}

type FinalityType string

const (
	MempoolPending FinalityType = "Pending"
	BatchFinal     FinalityType = "Final"
)

type QueryPagination struct {
	Offset uint64
	Size   uint
}

func (p *QueryPagination) UnmarshalJSON(data []byte) error {
	// Use a temporary struct to avoid infinite unmarshalling loop
	type Temp struct {
		Size   uint `json:"size"`
		Offset uint64
	}

	var temp Temp
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Size < 1 || temp.Size > 100 {
		return fmt.Errorf("size must be between 1 and 100")
	}

	p.Size = temp.Size
	p.Offset = temp.Offset
	return nil
}

type PrivateCustomQueryListTransactions struct {
	Address    common.Address  `json:"address"`
	Pagination QueryPagination `json:"pagination"`
}
