package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"
)

// responsible for saving event logs
type eventsStorage struct {
	db             enclavedb.EnclaveDB
	cachingService *CacheService
	logger         gethlog.Logger
}

func newEventsStorage(cachingService *CacheService, db enclavedb.EnclaveDB, logger gethlog.Logger) *eventsStorage {
	return &eventsStorage{cachingService: cachingService, db: db, logger: logger}
}

func (es *eventsStorage) storeReceiptAndEventLogs(ctx context.Context, dbTX *sqlx.Tx, batch *common.BatchHeader, txExecResult *core.TxExecResult) error {
	txId, senderId, err := enclavedb.ReadTransactionIdAndSender(ctx, dbTX, txExecResult.Receipt.TxHash)
	if err != nil {
		return fmt.Errorf("could not get transaction id. Cause: %w", err)
	}

	// store the contracts created by this tx
	for createdContract, cfg := range txExecResult.CreatedContracts {
		err := es.storeNewContractWithEventTypeConfigs(ctx, dbTX, createdContract, senderId, cfg, *txId)
		if err != nil {
			return err
		}
	}

	receiptId, err := es.storeReceipt(ctx, dbTX, batch, txExecResult, txId)
	if err != nil {
		return err
	}

	isReceiptPublic := false
	eoaMap := make(map[uint64]bool)
	for _, l := range txExecResult.Receipt.Logs {
		ets, err := es.storeEventLog(ctx, dbTX, receiptId, l)
		if err != nil {
			return fmt.Errorf("could not store log entry %v. Cause: %w", l, err)
		}
		contract, err := es.readContract(ctx, dbTX, l.Address)
		if err != nil {
			return fmt.Errorf("could not read contract address. %s. Cause: %w", l.Address, err)
		}
		eventType := contract.EventType(l.Topics[0])
		if eventType.IsPublic() {
			isReceiptPublic = true
		}
		for _, et := range ets {
			if et != nil && et.RelevantAddressId != nil {
				eoaMap[*et.RelevantAddressId] = true
			}
		}
	}

	eoas := make([]uint64, 0)
	for k := range eoaMap {
		eoas = append(eoas, k)
	}

	err = enclavedb.WriteReceiptViewers(ctx, dbTX, receiptId, isReceiptPublic, eoas)
	if err != nil {
		return err
	}
	return nil
}

func (es *eventsStorage) storeNewContractWithEventTypeConfigs(ctx context.Context, dbTX *sqlx.Tx, contractAddr gethcommon.Address, senderId *uint64, cfg *core.ContractVisibilityConfig, txId uint64) error {
	_, err := enclavedb.WriteContractConfig(ctx, dbTX, contractAddr, *senderId, cfg, txId)
	if err != nil {
		return fmt.Errorf("could not write contract address. cause %w", err)
	}

	c, err := es.readContract(ctx, dbTX, contractAddr)
	if err != nil {
		return err
	}

	// create the event types for the events that were configured
	for eventSig, eventCfg := range cfg.EventConfigs {
		et := enclavedb.EventType{
			Contract:       c,
			EventSignature: eventSig,
			AutoVisibility: eventCfg.AutoConfig,
			ConfigPublic:   eventCfg.Public,
			Topic1CanView:  eventCfg.Topic1CanView,
			Topic2CanView:  eventCfg.Topic2CanView,
			Topic3CanView:  eventCfg.Topic3CanView,
			SenderCanView:  eventCfg.SenderCanView,
		}

		// sanity check
		err = et.Validate()
		if err != nil {
			return err
		}

		_, err = enclavedb.WriteEventType(ctx, dbTX, &et)
		if err != nil {
			return fmt.Errorf("could not write event type. cause %w", err)
		}

		es.cachingService.InvalidateContract(contractAddr)
	}
	return nil
}

func (es *eventsStorage) storeReceipt(ctx context.Context, dbTX *sqlx.Tx, batch *common.BatchHeader, txExecResult *core.TxExecResult, txId *uint64) (uint64, error) {
	execTxId, err := enclavedb.WriteReceipt(ctx, dbTX, batch.SequencerOrderNo.Uint64(), txId, txExecResult.Receipt)
	if err != nil {
		return 0, fmt.Errorf("could not write receipt. Cause: %w", err)
	}
	return execTxId, nil
}

func (es *eventsStorage) storeEventLog(ctx context.Context, dbTX *sqlx.Tx, receiptId uint64, l *types.Log) ([]*enclavedb.EventTopic, error) {
	contract, err := es.readContract(ctx, dbTX, l.Address)
	if err != nil {
		// the contract should already have been stored when it was created
		return nil, fmt.Errorf("could not read contract address. %s. Cause: %w", l.Address, err)
	}

	eventSig := l.Topics[0]
	eventType := contract.EventType(eventSig)
	if eventType == nil {
		// this is the first type an event of this type is emitted, so we must store it
		eventType, err = es.storeAutoConfigEventType(ctx, dbTX, contract, l)
		if err != nil {
			return nil, fmt.Errorf("could not write event type. cause %w", err)
		}
	}

	topics, err := es.storeTopics(ctx, dbTX, eventType, l)
	if err != nil {
		return nil, fmt.Errorf("could not store topics. cause: %w", err)
	}

	// normalize data
	data := l.Data
	if len(data) == 0 {
		data = nil
	}
	err = enclavedb.WriteEventLog(ctx, dbTX, eventType.Id, topics, data, l.Index, receiptId)
	if err != nil {
		return nil, fmt.Errorf("could not write event log. Cause: %w", err)
	}

	// event types that were not configured explicitly can be "Public events" as well.
	// based on the topics, this logic determines whether the event type has any relevant addresses
	// this is called only the first time an event is emitted
	err = es.setAutoVisibilityWhenEventFirstEmitted(ctx, dbTX, eventType, topics)
	if err != nil {
		return nil, fmt.Errorf("could not update the auto visibility. Cause: %w", err)
	}

	return topics, nil
}

func (es *eventsStorage) setAutoVisibilityWhenEventFirstEmitted(ctx context.Context, dbTX *sqlx.Tx, eventType *enclavedb.EventType, topics []*enclavedb.EventTopic) error {
	if !eventType.ConfigPublic && eventType.AutoVisibility && eventType.AutoPublic == nil {
		isPublic := true
		for _, topic := range topics {
			if topic != nil {
				if topic.RelevantAddressId != nil {
					isPublic = false
					break
				}
			}
		}
		// for private events with autovisibility, the first time we need to determine whether they are public
		err := enclavedb.UpdateEventTypeAutoPublic(ctx, dbTX, eventType.Id, isPublic)
		if err != nil {
			return fmt.Errorf("could not update event type. cause: %w", err)
		}
		es.cachingService.InvalidateContract(eventType.Contract.Address)
	}
	return nil
}

// stores an event type the first time it is emitted
// since it wasn't saved on contract deployment, it means that there is no explicit configuration for it
func (es *eventsStorage) storeAutoConfigEventType(ctx context.Context, dbTX *sqlx.Tx, contract *enclavedb.Contract, l *types.Log) (*enclavedb.EventType, error) {
	eventType := enclavedb.EventType{
		Contract:       contract,
		EventSignature: l.Topics[0],
		ConfigPublic:   contract.IsTransparent(),
	}

	// event types that are not public - will have the default rules
	if !eventType.ConfigPublic {
		eventType.AutoVisibility = true
	}

	// sanity check
	err := eventType.Validate()
	if err != nil {
		return nil, err
	}

	id, err := enclavedb.WriteEventType(ctx, dbTX, &eventType)
	if err != nil {
		return nil, fmt.Errorf("could not write event type. cause: %w", err)
	}
	es.cachingService.InvalidateContract(contract.Address)
	eventType.Id = id
	return &eventType, nil
}

func (es *eventsStorage) storeTopics(ctx context.Context, dbTX *sqlx.Tx, eventType *enclavedb.EventType, l *types.Log) ([]*enclavedb.EventTopic, error) {
	eventTopics := make([]*enclavedb.EventTopic, 3)
	// iterate the topics containing user values
	// reuse them if already inserted
	// if not, discover if there is a relevant externally owned address
	for i := 1; i < len(l.Topics); i++ {
		topic := l.Topics[i]
		// first check if there is an entry already for this topic
		eventTopic, err := es.findTopic(ctx, dbTX, topic.Bytes(), eventType.Id)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not read the event topic. Cause: %w", err)
		}
		if errors.Is(err, errutil.ErrNotFound) {
			// if no entry was found
			eventTopic, err = es.storeTopic(ctx, dbTX, eventType, i, topic)
			if err != nil {
				return nil, fmt.Errorf("could not store the event topic. Cause: %w", err)
			}
		}
		eventTopics[i-1] = eventTopic
	}
	return eventTopics, nil
}

// this function contains visibility logic
func (es *eventsStorage) storeTopic(ctx context.Context, dbTX *sqlx.Tx, eventType *enclavedb.EventType, topicNo int, topic gethcommon.Hash) (*enclavedb.EventTopic, error) {
	relevantAddress, err := es.determineRelevantAddressForTopic(ctx, dbTX, eventType, topicNo, topic)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not determine visibility rules. cause: %w", err)
	}

	var relAddressId *uint64
	if relevantAddress != nil {
		var err error
		relAddressId, err = es.readEOA(ctx, dbTX, *relevantAddress)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}
		if relAddressId == nil {
			es.logger.Debug("EOA not found when saving topic", "topic", topic.Hex())
		}
	}
	eventTopicId, err := enclavedb.WriteEventTopic(ctx, dbTX, &topic, relAddressId, eventType.Id)
	if err != nil {
		return nil, fmt.Errorf("could not write event topic. Cause: %w", err)
	}
	return &enclavedb.EventTopic{Id: eventTopicId, RelevantAddressId: relAddressId}, nil
}

// based on the eventType configuration, this function determines the address that can view events logs containing this topic
func (es *eventsStorage) determineRelevantAddressForTopic(ctx context.Context, dbTX *sqlx.Tx, eventType *enclavedb.EventType, topicNumber int, topic gethcommon.Hash) (*gethcommon.Address, error) {
	var relevantAddress *gethcommon.Address
	switch {
	case eventType.AutoVisibility:
		extractedAddr := common.ExtractPotentialAddress(topic)
		if extractedAddr == nil {
			break
		}

		// first check whether there is already an entry in the EOA table
		_, err := es.readEOA(ctx, dbTX, *extractedAddr)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}
		if err == nil {
			relevantAddress = extractedAddr
			break
		}

		// if the address is a contract then it's clearly not an EOA
		_, err = es.readContract(ctx, dbTX, *extractedAddr)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}
		if err == nil {
			// the extracted address is a contract
			relevantAddress = nil
			break
		}

		// save the extracted address to the EOA table
		relevantAddress = extractedAddr
		_, err = enclavedb.WriteEoa(ctx, dbTX, *relevantAddress)
		if err != nil {
			return nil, err
		}

	case eventType.IsPublicConfig():
		// for public events, there is no relevant address
		relevantAddress = nil

	case eventType.IsTopicRelevant(topicNumber):
		relevantAddress = common.ExtractPotentialAddress(topic)
		// it is possible for contracts to emit events without an actual address.
		// for example. ERC20.mint emits a transfer event from a "0" address
		if relevantAddress == nil {
			es.logger.Debug(fmt.Sprintf("invalid configuration. expected address in topic %d : %s", topicNumber, topic.String()))
			return nil, errutil.ErrNotFound
		}

	case !eventType.IsTopicRelevant(topicNumber):
		// no need to do anything because this topic was not configured to be an address
		relevantAddress = nil

	default:
		es.logger.Crit("impossible case. Should not get here")
	}

	return relevantAddress, nil
}

func (es *eventsStorage) readContract(ctx context.Context, dbTX *sqlx.Tx, addr gethcommon.Address) (*enclavedb.Contract, error) {
	defer es.logDuration("readContract", measure.NewStopwatch())
	return es.cachingService.ReadContractAddr(ctx, addr, func() (*enclavedb.Contract, error) {
		c, err := enclavedb.ReadContractByAddress(ctx, dbTX, addr)
		if err != nil {
			return nil, err
		}
		ets, err := enclavedb.ReadEventTypesForContract(ctx, dbTX, c.Id)
		if err != nil {
			return nil, err
		}
		c.SetEventTypes(ets)
		return c, nil
	})
}

func (es *eventsStorage) ReadContract(ctx context.Context, addr gethcommon.Address) (*enclavedb.Contract, error) {
	defer es.logDuration("readContract", measure.NewStopwatch())
	dbtx, err := es.db.GetSQLDB().Beginx()
	if err != nil {
		return nil, err
	}
	defer dbtx.Rollback()
	return es.readContract(ctx, dbtx, addr)
}

func (es *eventsStorage) findTopic(ctx context.Context, dbTX *sqlx.Tx, topic []byte, eventTypeId uint64) (*enclavedb.EventTopic, error) {
	defer es.logDuration("findTopic", measure.NewStopwatch())
	return es.cachingService.ReadEventTopic(ctx, topic, eventTypeId, func() (*enclavedb.EventTopic, error) {
		return enclavedb.ReadEventTopic(ctx, dbTX, topic, eventTypeId)
	})
}

func (es *eventsStorage) readEOA(ctx context.Context, dbTX *sqlx.Tx, addr gethcommon.Address) (*uint64, error) {
	defer es.logDuration("ReadEOA", measure.NewStopwatch())
	return es.cachingService.ReadEOA(ctx, addr, func() (*uint64, error) {
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
