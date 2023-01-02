package common

import (
	"bytes"
	"encoding/json"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiptJSON struct {
	Type              hexutil.Uint64      `json:"type,omitempty"`
	PostState         hexutil.Bytes       `json:"root"`
	Status            hexutil.Uint64      `json:"status"`
	CumulativeGasUsed hexutil.Uint64      `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             types.Bloom         `json:"logsBloom"         gencodec:"required"`
	Logs              []*types.Log        `json:"logs"              gencodec:"required"`
	TxHash            gethcommon.Hash     `json:"transactionHash" gencodec:"required"`
	ContractAddress   *gethcommon.Address `json:"contractAddress"`
	GasUsed           hexutil.Uint64      `json:"gasUsed" gencodec:"required"`
	BlockHash         gethcommon.Hash     `json:"blockHash,omitempty"`
	BlockNumber       *hexutil.Big        `json:"blockNumber,omitempty"`
	TransactionIndex  hexutil.Uint        `json:"transactionIndex"`
}

// MarshalReceipt - serializes a compliant receipt and returns it.
// Cloned from the standard geth marshal function, but fixed issue with
// the returned contract address being returned as zero address instead of nil.
func MarshalReceipt(r *types.Receipt) ([]byte, error) {
	var enc ReceiptJSON
	enc.Type = hexutil.Uint64(r.Type)
	enc.PostState = r.PostState
	enc.Status = hexutil.Uint64(r.Status)
	enc.CumulativeGasUsed = hexutil.Uint64(r.CumulativeGasUsed)
	enc.Bloom = r.Bloom
	enc.Logs = r.Logs
	enc.TxHash = r.TxHash
	// Contract address should only be set for receipts of transactions that initialize a contract.
	// Otherwise return nil to match receipts returned by other testnets.
	if !bytes.Equal(r.ContractAddress.Bytes(), gethcommon.BigToAddress(gethcommon.Big0).Bytes()) {
		enc.ContractAddress = &r.ContractAddress
	}
	enc.GasUsed = hexutil.Uint64(r.GasUsed)
	enc.BlockHash = r.BlockHash
	enc.BlockNumber = (*hexutil.Big)(r.BlockNumber)
	enc.TransactionIndex = hexutil.Uint(r.TransactionIndex)
	return json.Marshal(&enc)
}
