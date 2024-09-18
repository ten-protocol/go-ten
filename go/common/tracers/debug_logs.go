package tracers

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// DebugLogs are the logs returned when using the DebugGetLogs endpoint
type DebugLogs struct {
	RelAddress1    *gethcommon.Address `json:"relAddress1"`
	RelAddress2    *gethcommon.Address `json:"relAddress2"`
	RelAddress3    *gethcommon.Address `json:"relAddress3"`
	LifecycleEvent bool                `json:"lifecycleEvent"`

	gethtypes.Log
}

// MarshalJSON marshals as JSON.
// this holds a copy of the gethtypes.Log log marshaller
func (l DebugLogs) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address        string              `json:"address" gencodec:"required"`
		Topics         []gethcommon.Hash   `json:"topics" gencodec:"required"`
		Data           hexutil.Bytes       `json:"data" gencodec:"required"`
		BlockNumber    uint64              `json:"blockNumber"`
		TxHash         gethcommon.Hash     `json:"transactionHash" gencodec:"required"`
		TxIndex        uint                `json:"transactionIndex"`
		BlockHash      gethcommon.Hash     `json:"blockHash"`
		Index          uint                `json:"logIndex"`
		Removed        bool                `json:"removed"`
		LifecycleEvent bool                `json:"lifecycleEvent"`
		RelAddress1    *gethcommon.Address `json:"relAddress1"`
		RelAddress2    *gethcommon.Address `json:"relAddress2"`
		RelAddress3    *gethcommon.Address `json:"relAddress3"`
	}{
		l.Address.Hex(),
		l.Topics,
		l.Data,
		l.BlockNumber,
		l.TxHash,
		l.TxIndex,
		l.BlockHash,
		l.Index,
		l.Removed,
		l.LifecycleEvent,
		l.RelAddress1,
		l.RelAddress2,
		l.RelAddress3,
	})
}
