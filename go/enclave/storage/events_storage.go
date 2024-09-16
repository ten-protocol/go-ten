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

func (es *eventsStorage) storeReceiptAndEventLogs(ctx context.Context, dbTX *sql.Tx, batch *common.BatchHeader, txExecResult *core.TxExecResult) error {
	txId, senderId, err := enclavedb.ReadTransactionIdAndSender(ctx, dbTX, txExecResult.Receipt.TxHash)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not get transaction id. Cause: %w", err)
	}

	// store the contracts created by this tx
	for createdContract, cfg := range txExecResult.CreatedContracts {
		_, err := es.storeNewContract(ctx, dbTX, createdContract, senderId, cfg)
		if err != nil {
			return err
		}

		c, err := es.readContract(ctx, dbTX, createdContract)
		if err != nil {
			return err
		}

		// create the event types for the events that were configured
		for eventSig, eventCfg := range cfg.EventConfigs {
			_, err = enclavedb.WriteEventType(ctx, dbTX, &enclavedb.EventType{
				Contract:       c,
				EventSignature: eventSig,
				AutoVisibility: eventCfg.AutoConfig,
				Public:         eventCfg.Public,
				Topic1CanView:  eventCfg.Topic1CanView,
				Topic2CanView:  eventCfg.Topic2CanView,
				Topic3CanView:  eventCfg.Topic3CanView,
				SenderCanView:  eventCfg.SenderCanView,
			})
			if err != nil {
				return fmt.Errorf("could not write event type. cause %w", err)
			}
		}
	}

	receiptId, err := es.storeReceipt(ctx, dbTX, batch, txExecResult, txId)
	if err != nil {
		return err
	}

	for _, l := range txExecResult.Receipt.Logs {
		err := es.storeEventLog(ctx, dbTX, receiptId, l)
		if err != nil {
			return fmt.Errorf("could not store log entry %v. Cause: %w", l, err)
		}
	}

	return nil
}

// Convert the receipt into its storage form and serialize
// this removes information that can be recreated
// todo - in a future iteration, this can be slimmed down further because we already store the logs separately
func (es *eventsStorage) storeReceipt(ctx context.Context, dbTX *sql.Tx, batch *common.BatchHeader, txExecResult *core.TxExecResult, txId *uint64) (uint64, error) {
	storageReceipt := (*types.ReceiptForStorage)(txExecResult.Receipt)
	receiptBytes, err := rlp.EncodeToBytes(storageReceipt)
	if err != nil {
		return 0, fmt.Errorf("failed to encode block receipts. Cause: %w", err)
	}

	execTxId, err := enclavedb.WriteReceipt(ctx, dbTX, batch.SequencerOrderNo.Uint64(), txId, receiptBytes)
	if err != nil {
		return 0, fmt.Errorf("could not write receipt. Cause: %w", err)
	}
	return execTxId, nil
}

func (es *eventsStorage) storeNewContract(ctx context.Context, dbTX *sql.Tx, createdContract gethcommon.Address, senderId *uint64, cfg *core.ContractVisibilityConfig) (*uint64, error) {
	ctrId, err := enclavedb.WriteContractConfig(ctx, dbTX, createdContract, *senderId, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not write contract address. cause %w", err)
	}
	return ctrId, nil
}

func (es *eventsStorage) storeEventLog(ctx context.Context, dbTX *sql.Tx, execTxId uint64, l *types.Log) error {
	eventSig := l.Topics[0]

	contract, err := es.readContract(ctx, dbTX, l.Address)
	if err != nil {
		// the contract should already have been stored when it was created
		return fmt.Errorf("could not read contract address. %s. Cause: %w", l.Address, err)
	}

	eventType, err := es.readEventType(ctx, dbTX, l.Address, eventSig)
	if errors.Is(err, errutil.ErrNotFound) {
		// this is the first type an event of this type is emitted, so we must store it
		eventType, err = es.storeAutoConfigEventType(ctx, dbTX, contract, l)
		if err != nil {
			return fmt.Errorf("could not write event type. cause %w", err)
		}
	} else if err != nil {
		// unexpected event type
		return fmt.Errorf("could not read event type. Cause: %w", err)
	}

	topicIds, err := es.storeTopics(ctx, dbTX, eventType, l)
	if err != nil {
		return fmt.Errorf("could not store topics. cause: %w", err)
	}

	// normalize data
	data := l.Data
	if len(data) == 0 {
		data = nil
	}
	err = enclavedb.WriteEventLog(ctx, dbTX, eventType.Id, topicIds, data, l.Index, execTxId)
	if err != nil {
		return fmt.Errorf("could not write event log. Cause: %w", err)
	}

	return nil
}

// handles the visibility config detection
func (es *eventsStorage) storeAutoConfigEventType(ctx context.Context, dbTX *sql.Tx, contract *enclavedb.Contract, l *types.Log) (*enclavedb.EventType, error) {
	eventType := enclavedb.EventType{
		Contract:       contract,
		EventSignature: l.Topics[0],
		Public:         contract.IsTransparent(),
	}

	// event types that are not public - will have the default rules
	if !eventType.Public {
		eventType.AutoVisibility = true
	}

	id, err := enclavedb.WriteEventType(ctx, dbTX, &eventType)
	if err != nil {
		return nil, fmt.Errorf("could not write event type. cause: %w", err)
	}
	eventType.Id = id
	return &eventType, nil
}

func (es *eventsStorage) storeTopics(ctx context.Context, dbTX *sql.Tx, eventType *enclavedb.EventType, l *types.Log) ([]*uint64, error) {
	topicIds := make([]*uint64, 3)
	// iterate the topics containing user values
	// reuse them if already inserted
	// if not, discover if there is a relevant externally owned address
	for i := 1; i < len(l.Topics); i++ {
		topic := l.Topics[i]
		// first check if there is an entry already for this topic
		eventTopicId, _, err := es.findEventTopic(ctx, dbTX, topic.Bytes())
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not read the event topic. Cause: %w", err)
		}
		if errors.Is(err, errutil.ErrNotFound) {
			// if no entry was found
			eventTopicId, err = es.storeEventTopic(ctx, dbTX, eventType, i, topic)
			if err != nil {
				return nil, fmt.Errorf("could not read the event topic. Cause: %w", err)
			}
		}
		topicIds[i-1] = &eventTopicId
	}
	return topicIds, nil
}

// this function contains visibility logic
func (es *eventsStorage) storeEventTopic(ctx context.Context, dbTX *sql.Tx, eventType *enclavedb.EventType, i int, topic gethcommon.Hash) (uint64, error) {
	relevantAddress, err := es.visibililty(ctx, dbTX, eventType, i, topic)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return 0, fmt.Errorf("could not determine visibility rules. cause: %w", err)
	}

	var relAddressId *uint64
	if relevantAddress != nil {
		var err error
		relAddressId, err = es.readEOA(ctx, dbTX, *relevantAddress)
		if err != nil {
			return 0, err
		}
	}
	eventTopicId, err := enclavedb.WriteEventTopic(ctx, dbTX, &topic, relAddressId)
	if err != nil {
		return 0, fmt.Errorf("could not write event topic. Cause: %w", err)
	}
	return eventTopicId, nil
}

func (es *eventsStorage) visibililty(ctx context.Context, dbTX *sql.Tx, eventType *enclavedb.EventType, i int, topic gethcommon.Hash) (*gethcommon.Address, error) {
	var relevantAddress *gethcommon.Address
	switch {
	case eventType.AutoVisibility:
		var err error
		// if there is no configuration, we have to autodetect the address
		relevantAddress, err = es.autoDetectRelevantAddress(ctx, dbTX, topic)
		if err != nil {
			return nil, err
		}

		// when autodetecting, we assume that any address that is not a contract is an EOA
		_, err = es.readEOA(ctx, dbTX, *relevantAddress)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}
		if errors.Is(err, errutil.ErrNotFound) {
			_, err := enclavedb.WriteEoa(ctx, dbTX, *relevantAddress)
			if err != nil {
				return nil, err
			}
		}

	case eventType.IsPublic():
		// for public events, there is no relevant address
		relevantAddress = nil

	case eventType.IsTopicRelevant(i):
		relevantAddress = common.ExtractPotentialAddress(topic)
		if relevantAddress == nil {
			return nil, fmt.Errorf("invalid configuration. expected address in topic %d : %s", i, topic.String())
		}

	default:
		es.logger.Crit("impossible case. Should not get here")
	}

	return relevantAddress, nil
}

// Of the log's topics, returns those that are (potentially) user addresses. A topic is considered a user address if:
//   - It has at least 12 leading zero bytes (since addresses are 20 bytes long, while hashes are 32) and at most 22 leading zero bytes
//   - It is not a smart contract address
func (es *eventsStorage) autoDetectRelevantAddress(ctx context.Context, dbTX *sql.Tx, topic gethcommon.Hash) (*gethcommon.Address, error) {
	potentialAddr := common.ExtractPotentialAddress(topic)
	if potentialAddr == nil {
		return nil, errutil.ErrNotFound
	}

	// first check whether there is already an entry in the EOA table
	_, err := es.readEOA(ctx, dbTX, *potentialAddr)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return potentialAddr, nil
	}

	// if the address is a contract then it's clearly not an EOA
	_, err = es.readContract(ctx, dbTX, *potentialAddr)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, errutil.ErrNotFound
	}

	return potentialAddr, nil
}

func (es *eventsStorage) readEventType(ctx context.Context, dbTX *sql.Tx, contractAddress gethcommon.Address, eventSignature gethcommon.Hash) (*enclavedb.EventType, error) {
	defer es.logDuration("ReadEventType", measure.NewStopwatch())

	return es.cachingService.ReadEventType(ctx, contractAddress, eventSignature, func(v any) (*enclavedb.EventType, error) {
		contract, err := es.readContract(ctx, dbTX, contractAddress)
		if err != nil {
			return nil, err
		}
		return enclavedb.ReadEventType(ctx, dbTX, contract, eventSignature)
	})
}

func (es *eventsStorage) readContract(ctx context.Context, dbTX *sql.Tx, addr gethcommon.Address) (*enclavedb.Contract, error) {
	defer es.logDuration("readContract", measure.NewStopwatch())
	return es.cachingService.ReadContractAddr(ctx, addr, func(v any) (*enclavedb.Contract, error) {
		return enclavedb.ReadContractByAddress(ctx, dbTX, addr)
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
