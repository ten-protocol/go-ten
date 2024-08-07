package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"
)

// responsible for saving event logs
type eventsStorage struct {
	cachingService *CacheService
	logger         gethlog.Logger
}

func newEventsStorage(cachingService *CacheService, logger gethlog.Logger) *eventsStorage {
	return &eventsStorage{cachingService: cachingService, logger: logger}
}

func (es *eventsStorage) storeReceiptAndEventLogs(ctx context.Context, dbTX *sql.Tx, batch *common.BatchHeader, receipt *types.Receipt, createdContracts []*gethcommon.Address) error {
	txId, senderId, err := enclavedb.ReadTransactionIdAndSender(ctx, dbTX, receipt.TxHash)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not get transaction id. Cause: %w", err)
	}

	for _, createdContract := range createdContracts {
		_, err = enclavedb.WriteContractAddress(ctx, dbTX, createdContract, *senderId)
		if err != nil {
			return fmt.Errorf("could not write contract address. cause %w", err)
		}
	}

	// Convert the receipt into its storage form and serialize
	// this removes information that can be recreated
	// todo - in a future iteration, this can be slimmed down further because we already store the logs separately
	storageReceipt := (*types.ReceiptForStorage)(receipt)
	receiptBytes, err := rlp.EncodeToBytes(storageReceipt)
	if err != nil {
		return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
	}

	execTxId, err := enclavedb.WriteReceipt(ctx, dbTX, batch.SequencerOrderNo.Uint64(), txId, receiptBytes)
	if err != nil {
		return fmt.Errorf("could not write receipt. Cause: %w", err)
	}

	for _, l := range receipt.Logs {
		err := es.storeEventLog(ctx, dbTX, execTxId, l)
		if err != nil {
			return fmt.Errorf("could not store log entry %v. Cause: %w", l, err)
		}
	}
	return nil
}

func (es *eventsStorage) storeEventLog(ctx context.Context, dbTX *sql.Tx, execTxId uint64, l *types.Log) error {
	topicIds, isLifecycle, err := es.handleUserTopics(ctx, dbTX, l)
	if err != nil {
		return err
	}

	eventTypeId, err := es.handleEventType(ctx, dbTX, l, isLifecycle)
	if err != nil {
		return err
	}

	// normalize data
	data := l.Data
	if len(data) == 0 {
		data = nil
	}
	err = enclavedb.WriteEventLog(ctx, dbTX, eventTypeId, topicIds, data, l.Index, execTxId)
	if err != nil {
		return fmt.Errorf("could not write event log. Cause: %w", err)
	}

	return nil
}

func (es *eventsStorage) handleEventType(ctx context.Context, dbTX *sql.Tx, l *types.Log, isLifecycle bool) (uint64, error) {
	et, err := es.readEventType(ctx, dbTX, l.Address, l.Topics[0])
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return 0, fmt.Errorf("could not read event type. Cause: %w", err)
	}
	if err == nil {
		// in case we determined the current emitted event is not lifecycle, we must update the EventType
		if !isLifecycle && et.isLifecycle {
			err := enclavedb.UpdateEventTypeLifecycle(ctx, dbTX, et.id, isLifecycle)
			if err != nil {
				return 0, fmt.Errorf("could not update the event type. cause: %w", err)
			}
		}
		return et.id, nil
	}

	// the first time an event of this type is emitted we must store it
	contractAddId, err := es.readContractAddress(ctx, dbTX, l.Address)
	if err != nil {
		// the contract was already stored when it was created
		return 0, fmt.Errorf("could not read contract address. %s. Cause: %w", l.Address, err)
	}
	return enclavedb.WriteEventType(ctx, dbTX, contractAddId, l.Topics[0], isLifecycle)
}

func (es *eventsStorage) handleUserTopics(ctx context.Context, dbTX *sql.Tx, l *types.Log) ([]*uint64, bool, error) {
	topicIds := make([]*uint64, 3)
	// iterate the topics containing user values
	// reuse them if already inserted
	// if not, discover if there is a relevant externally owned address
	isLifecycle := true
	for i := 1; i < len(l.Topics); i++ {
		topic := l.Topics[i]
		// first check if there is an entry already for this topic
		eventTopicId, relAddressId, err := es.findEventTopic(ctx, dbTX, topic.Bytes())
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, false, fmt.Errorf("could not read the event topic. Cause: %w", err)
		}
		if errors.Is(err, errutil.ErrNotFound) {
			// check whether the topic is an EOA
			relAddressId, err = es.findRelevantAddress(ctx, dbTX, topic)
			if err != nil && !errors.Is(err, errutil.ErrNotFound) {
				return nil, false, fmt.Errorf("could not read relevant address. Cause %w", err)
			}
			eventTopicId, err = enclavedb.WriteEventTopic(ctx, dbTX, &topic, relAddressId)
			if err != nil {
				return nil, false, fmt.Errorf("could not write event topic. Cause: %w", err)
			}
		}

		if relAddressId != nil {
			isLifecycle = false
		}
		topicIds[i-1] = &eventTopicId
	}
	return topicIds, isLifecycle, nil
}

// Of the log's topics, returns those that are (potentially) user addresses. A topic is considered a user address if:
//   - It has at least 12 leading zero bytes (since addresses are 20 bytes long, while hashes are 32) and at most 22 leading zero bytes
//   - It is not a smart contract address
func (es *eventsStorage) findRelevantAddress(ctx context.Context, dbTX *sql.Tx, topic gethcommon.Hash) (*uint64, error) {
	potentialAddr := common.ExtractPotentialAddress(topic)
	if potentialAddr == nil {
		return nil, errutil.ErrNotFound
	}

	// first check whether there is already an entry in the EOA table
	eoaID, err := es.readEOA(ctx, dbTX, *potentialAddr)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return eoaID, nil
	}

	// if the address is a contract then it's clearly not an EOA
	_, err = es.readContractAddress(ctx, dbTX, *potentialAddr)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, errutil.ErrNotFound
	}

	// when we reach this point, the value looks like an address, but we haven't yet seen it
	// for the first iteration, we'll just assume it's an EOA
	// we can make this smarter by passing in more information about the event
	id, err := enclavedb.WriteEoa(ctx, dbTX, *potentialAddr)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (es *eventsStorage) readEventType(ctx context.Context, dbTX *sql.Tx, contractAddress gethcommon.Address, eventSignature gethcommon.Hash) (*EventType, error) {
	defer es.logDuration("ReadEventType", measure.NewStopwatch())

	return es.cachingService.ReadEventType(ctx, contractAddress, eventSignature, func(v any) (*EventType, error) {
		contractAddrId, err := enclavedb.ReadContractAddress(ctx, dbTX, contractAddress)
		if err != nil {
			return nil, err
		}
		id, isLifecycle, err := enclavedb.ReadEventType(ctx, dbTX, *contractAddrId, eventSignature)
		if err != nil {
			return nil, err
		}
		return &EventType{
			id:          id,
			isLifecycle: isLifecycle,
		}, nil
	})
}

func (es *eventsStorage) readContractAddress(ctx context.Context, dbTX *sql.Tx, addr gethcommon.Address) (*uint64, error) {
	defer es.logDuration("readContractAddress", measure.NewStopwatch())
	return es.cachingService.ReadContractAddr(ctx, addr, func(v any) (*uint64, error) {
		return enclavedb.ReadContractAddress(ctx, dbTX, addr)
	})
}

func (es *eventsStorage) findEventTopic(ctx context.Context, dbTX *sql.Tx, topic []byte) (uint64, *uint64, error) {
	defer es.logDuration("findEventTopic", measure.NewStopwatch())
	return enclavedb.ReadEventTopic(ctx, dbTX, topic)
}

func (es *eventsStorage) readEOA(ctx context.Context, dbTX *sql.Tx, addr gethcommon.Address) (*uint64, error) {
	defer es.logDuration("ReadEOA", measure.NewStopwatch())
	return es.cachingService.ReadEOA(ctx, addr, func(v any) (*uint64, error) {
		id, err := enclavedb.ReadEoa(ctx, dbTX, addr)
		if err != nil {
			return nil, err
		}
		return &id, nil
	})
}

func (es *eventsStorage) logDuration(method string, stopWatch *measure.Stopwatch) {
	core.LogMethodDuration(es.logger, stopWatch, fmt.Sprintf("Storage::%s completed", method))
}
