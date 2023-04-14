package tracers

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// DebugLogs are the logs returned when using the DebugGetLogs endpoint
type DebugLogs struct {
	RelAddress1    gethcommon.Hash `json:"relAddress1"`
	RelAddress2    gethcommon.Hash `json:"relAddress2"`
	RelAddress3    gethcommon.Hash `json:"relAddress3"`
	RelAddress4    gethcommon.Hash `json:"relAddress4"`
	LifecycleEvent bool            `json:"lifecycleEvent"`

	gethtypes.Log
}

// MarshalJSON marshals as JSON.
// this holds a copy of the gethtypes.Log log marshaller
func (l DebugLogs) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RelAddress1    gethcommon.Hash    `json:"relAddress1"`
		RelAddress2    gethcommon.Hash    `json:"relAddress2"`
		RelAddress3    gethcommon.Hash    `json:"relAddress3"`
		RelAddress4    gethcommon.Hash    `json:"relAddress4"`
		LifecycleEvent bool               `json:"lifecycleEvent"`
		Address        gethcommon.Address `json:"address" gencodec:"required"`
		Topics         []gethcommon.Hash  `json:"topics" gencodec:"required"`
		Data           hexutil.Bytes      `json:"data" gencodec:"required"`
		BlockNumber    hexutil.Uint64     `json:"blockNumber"`
		TxHash         gethcommon.Hash    `json:"transactionHash" gencodec:"required"`
		TxIndex        hexutil.Uint       `json:"transactionIndex"`
		BlockHash      gethcommon.Hash    `json:"blockHash"`
		Index          hexutil.Uint       `json:"logIndex"`
		Removed        bool               `json:"removed"`
	}{
		l.RelAddress1,
		l.RelAddress2,
		l.RelAddress3,
		l.RelAddress4,
		l.LifecycleEvent,
		l.Address,
		l.Topics,
		l.Data,
		hexutil.Uint64(l.BlockNumber),
		l.TxHash,
		hexutil.Uint(l.TxIndex),
		l.BlockHash,
		hexutil.Uint(l.Index),
		l.Removed,
	})
}
