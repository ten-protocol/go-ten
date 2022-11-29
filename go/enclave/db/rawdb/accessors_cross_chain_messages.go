package rawdb

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

func StoreL1Messages(db ethdb.KeyValueWriter, blockHash gethcommon.Hash, messages common.CrossChainMessages, logger gethlog.Logger) bool {
	data, err := rlp.EncodeToBytes(messages)
	if err != nil {
		logger.Crit("Failed to encode the synthetic transactions...", log.ErrKey, err)
		return false
	}

	if err := db.Put(crossChainMessagesKey(blockHash), data); err != nil {
		logger.Crit("Failed to store the synthetic transactions...", log.ErrKey, err)
		return false
	}
	return true
}

func ReadL1Messages(db ethdb.KeyValueReader, blockHash gethcommon.Hash, logger gethlog.Logger) common.CrossChainMessages {
	var messages common.CrossChainMessages

	data, err := db.Get(crossChainMessagesKey(blockHash))
	if err != nil {
		logger.Info("Could not read key from db. ", log.ErrKey, err)
		return messages
	}

	err = rlp.DecodeBytes(data, &messages)
	if err != nil {
		logger.Info("Could not parse synthetic transactions from db.", log.ErrKey, err)
	}
	return messages
}
