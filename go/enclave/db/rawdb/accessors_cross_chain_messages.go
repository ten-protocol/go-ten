package rawdb

import (
	"errors"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

func StoreL1Messages(db ethdb.KeyValueWriter, blockHash gethcommon.Hash, messages common.CrossChainMessages, logger gethlog.Logger) error {
	data, err := rlp.EncodeToBytes(messages)
	if err != nil {
		logger.Crit("Failed to encode the synthetic transactions...", log.ErrKey, err)
		return err
	}

	if err := db.Put(crossChainMessagesKey(blockHash), data); err != nil {
		logger.Crit("Failed to store the synthetic transactions...", log.ErrKey, err)
		return err
	}
	return nil
}

func GetL1Messages(db ethdb.KeyValueReader, blockHash gethcommon.Hash, logger gethlog.Logger) (common.CrossChainMessages, error) {
	var messages common.CrossChainMessages

	data, err := db.Get(crossChainMessagesKey(blockHash))
	if err != nil {
		logger.Trace("Could not read key from db. ", log.ErrKey, err)
		// It is expected that not every block will have messages, thus do not surface it.
		if errors.Is(err, errutil.ErrNotFound) {
			return messages, nil
		}
		return nil, err
	}

	err = rlp.DecodeBytes(data, &messages)
	if err != nil {
		logger.Error("Could not parse synthetic transactions from db.", log.ErrKey, err)
		return nil, err
	}
	return messages, nil
}
