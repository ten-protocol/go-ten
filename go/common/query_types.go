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

type TransactionListingResponse struct {
	TransactionsData []PublicTransaction
	Total            uint64
}

type BatchListingResponse struct {
	BatchesData []PublicBatch
	Total       uint64
}

type BlockListingResponse struct {
	BlocksData []PublicBlock
	Total      uint64
}

type PublicTransaction struct {
	TransactionHash TxHash
	BatchHeight     *big.Int
	Finality        FinalityType
}

type PublicBatch struct {
	BatchHeader
	TxHashes []TxHash `json:"txHashes"`
}

type PublicBlock struct {
	BlockHeader types.Header `json:"blockHeader"`
	RollupHash  common.Hash  `json:"rollupHash"`
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

type ObscuroNetworkInfo struct {
	ManagementContractAddress common.Address
	L1StartHash               common.Hash
	SequencerID               common.Address
	MessageBusAddress         common.Address
	ImportantContracts        map[string]common.Address // map of contract name to address
}
