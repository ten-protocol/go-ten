package rawdb

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// This type is required because Geth serialises its logs in a reduced form to minimise storage space. For now, it is
// more straightforward for us to serialise all the fields by converting the logs to this type.
type logForStorage struct {
	Address     gethcommon.Address
	Topics      []gethcommon.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      gethcommon.Hash
	TxIndex     uint
	BlockHash   gethcommon.Hash
	Index       uint
}

func toLogForStorage(fullFatLog *types.Log) *logForStorage {
	return &logForStorage{
		Address:     fullFatLog.Address,
		Topics:      fullFatLog.Topics,
		Data:        fullFatLog.Data,
		BlockNumber: fullFatLog.BlockNumber,
		TxHash:      fullFatLog.TxHash,
		TxIndex:     fullFatLog.TxIndex,
		BlockHash:   fullFatLog.BlockHash,
		Index:       fullFatLog.Index,
	}
}

func (l logForStorage) toLog() *types.Log {
	return &types.Log{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		BlockNumber: l.BlockNumber,
		TxHash:      l.TxHash,
		TxIndex:     l.TxIndex,
		BlockHash:   l.BlockHash,
		Index:       l.Index,
	}
}
