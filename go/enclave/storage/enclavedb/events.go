package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	baseReceiptJoin = " from receipt rec " +
		"join batch b on rec.batch=b.sequence " +
		"join tx curr_tx on rec.tx=curr_tx.id " +
		"   join externally_owned_account tx_sender on curr_tx.sender_address=tx_sender.id " +
		"   left join contract tx_contr on curr_tx.to_address=tx_contr.id "

	baseEventJoin = " left join event_log e on e.receipt=rec.id " +
		"left join event_type et on e.event_type=et.id " +
		"	left join contract c on et.contract=c.id " +
		//"		left join tx creator_tx on c.tx=creator_tx.id " +
		"left join event_topic t1 on e.topic1=t1.id and et.id=t1.event_type " +
		"   left join externally_owned_account eoa1 on t1.rel_address=eoa1.id " +
		"left join event_topic t2 on e.topic2=t2.id and et.id=t2.event_type " +
		"   left join externally_owned_account eoa2 on t2.rel_address=eoa2.id " +
		"left join event_topic t3 on e.topic3=t3.id and et.id=t3.event_type " +
		"   left join externally_owned_account eoa3 on t3.rel_address=eoa3.id " +
		"where b.is_canonical=true "
)

func WriteEventType(ctx context.Context, dbTX *sql.Tx, et *EventType) (uint64, error) {
	res, err := dbTX.ExecContext(ctx, "insert into event_type (contract, event_sig, auto_visibility,auto_public, config_public, topic1_can_view, topic2_can_view, topic3_can_view, sender_can_view) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		et.Contract.Id, et.EventSignature.Bytes(), et.AutoVisibility, et.AutoPublic, et.ConfigPublic, et.Topic1CanView, et.Topic2CanView, et.Topic3CanView, et.SenderCanView)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func ReadEventType(ctx context.Context, dbTX *sql.Tx, contract *Contract, eventSignature gethcommon.Hash) (*EventType, error) {
	et := EventType{Contract: contract}
	err := dbTX.QueryRowContext(ctx,
		"select id, event_sig, auto_visibility, auto_public, config_public, topic1_can_view, topic2_can_view, topic3_can_view, sender_can_view from event_type where contract=? and event_sig=?",
		contract.Id, eventSignature.Bytes(),
	).Scan(&et.Id, &et.EventSignature, &et.AutoVisibility, &et.AutoPublic, &et.ConfigPublic, &et.Topic1CanView, &et.Topic2CanView, &et.Topic3CanView, &et.SenderCanView)
	if errors.Is(err, sql.ErrNoRows) {
		// make sure the error is converted to obscuro-wide not found error
		return nil, errutil.ErrNotFound
	}
	return &et, err
}

func WriteEventTopic(ctx context.Context, dbTX *sql.Tx, topic *gethcommon.Hash, addressId *uint64, eventTypeId uint64) (uint64, error) {
	res, err := dbTX.ExecContext(ctx, "insert into event_topic (event_type, topic, rel_address) values (?, ?, ?)", eventTypeId, topic.Bytes(), addressId)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func UpdateEventTypeAutoPublic(ctx context.Context, dbTx *sql.Tx, etId uint64, isPublic bool) error {
	_, err := dbTx.ExecContext(ctx, "update event_type set auto_public=? where id=?", isPublic, etId)
	return err
}

func ReadEventTopic(ctx context.Context, dbTX *sql.Tx, topic []byte, eventTypeId uint64) (*EventTopic, error) {
	var id uint64
	var address *uint64
	err := dbTX.QueryRowContext(ctx,
		"select id, rel_address from event_topic where topic=? and event_type=?", topic, eventTypeId).Scan(&id, &address)
	if errors.Is(err, sql.ErrNoRows) {
		// make sure the error is converted to obscuro-wide not found error
		return nil, errutil.ErrNotFound
	}
	return &EventTopic{Id: id, RelevantAddressId: address}, err
}

func ReadRelevantAddressFromEventTopic(ctx context.Context, dbTX *sql.Tx, id uint64) (*uint64, error) {
	var address *uint64
	err := dbTX.QueryRowContext(ctx,
		"select rel_address from event_topic where id=?", id).Scan(&address)
	if errors.Is(err, sql.ErrNoRows) {
		// make sure the error is converted to obscuro-wide not found error
		return nil, errutil.ErrNotFound
	}
	return address, err
}

func WriteEventLog(ctx context.Context, dbTX *sql.Tx, eventTypeId uint64, userTopics []*uint64, data []byte, logIdx uint, execTx uint64) error {
	_, err := dbTX.ExecContext(ctx, "insert into event_log (event_type, topic1, topic2, topic3, datablob, log_idx, receipt) values (?,?,?,?,?,?,?)",
		eventTypeId, userTopics[0], userTopics[1], userTopics[2], data, logIdx, execTx)
	return err
}

func FilterLogs(ctx context.Context, db *sql.DB, requestingAccount *gethcommon.Address, fromBlock, toBlock *big.Int, batchHash *common.L2BatchHash, addresses []gethcommon.Address, topics [][]gethcommon.Hash) ([]*types.Log, error) {
	queryParams := []any{}
	query := ""

	if batchHash != nil {
		query += " AND b.hash = ? "
		queryParams = append(queryParams, batchHash.Bytes())
	}

	// ignore negative numbers
	if fromBlock != nil && fromBlock.Sign() > 0 {
		query += " AND b.height >= ?"
		queryParams = append(queryParams, fromBlock.Int64())
	}
	if toBlock != nil && toBlock.Sign() > 0 {
		query += " AND b.height <= ?"
		queryParams = append(queryParams, toBlock.Int64())
	}

	if len(addresses) > 0 {
		query += " AND c.address in (" + repeat("?", ",", len(addresses)) + ")"
		for _, address := range addresses {
			queryParams = append(queryParams, address.Bytes())
		}
	}
	if len(topics) > 4 {
		return nil, fmt.Errorf("invalid filter. Too many topics")
	}

	for i := 0; i < len(topics); i++ {
		if len(topics[i]) > 0 {
			valuesIn := "IN (" + repeat("?", ",", len(topics[i])) + ")"
			if i == 0 {
				query += " AND et.event_sig " + valuesIn
			} else {
				query += " AND t" + strconv.Itoa(i) + ".topic " + valuesIn
			}
			for _, hash := range topics[i] {
				queryParams = append(queryParams, hash.Bytes())
			}
		}
	}

	_, logs, err := loadReceiptsAndEventLogs(ctx, db, requestingAccount, query, queryParams, false)
	return logs, err
}

func DebugGetLogs(ctx context.Context, db *sql.DB, fromBlock *big.Int, toBlock *big.Int, address gethcommon.Address, eventSig gethcommon.Hash) ([]*common.DebugLogVisibility, error) {
	var queryParams []any
	query := "select c.transparent, c.auto_visibility, et.config_public, et.topic1_can_view, et.topic2_can_view, et.topic3_can_view, et.sender_can_view, et.auto_visibility, et.auto_public, eoa1.address, eoa2.address, eoa3.address, b.height, curr_tx.hash, curr_tx.idx, b.hash, log_idx " +
		baseReceiptJoin + baseEventJoin

	// ignore negative numbers
	if fromBlock != nil && fromBlock.Sign() > 0 {
		query += " AND b.height >= ?"
		queryParams = append(queryParams, fromBlock.Int64())
	}
	if toBlock != nil && toBlock.Sign() > 0 {
		query += " AND b.height <= ?"
		queryParams = append(queryParams, toBlock.Int64())
	}

	query += " AND c.address = ? "
	queryParams = append(queryParams, address.Bytes())

	query += " AND et.event_sig = ? "
	queryParams = append(queryParams, eventSig.Bytes())

	result := make([]*common.DebugLogVisibility, 0)

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := common.DebugLogVisibility{EventSig: &eventSig, Address: &address}

		var relAddress1, relAddress2, relAddress3 []byte
		err = rows.Scan(
			&l.TransparentContract,
			&l.AutoContract,

			&l.EventConfigPublic,
			&l.Topic1,
			&l.Topic2,
			&l.Topic3,
			&l.Sender,

			&l.AutoVisibility,
			&l.AutoPublic,
			&relAddress1,
			&relAddress2,
			&relAddress3,

			&l.BlockNumber,
			&l.TxHash,
			&l.TxIndex,
			&l.BlockHash,
			&l.Index,
		)
		if err != nil {
			return nil, fmt.Errorf("could not debug load log entry from db: %w", err)
		}

		r1 := relAddress1 != nil
		r2 := relAddress2 != nil
		r3 := relAddress3 != nil
		l.RelAddress1 = &r1
		l.RelAddress2 = &r2
		l.RelAddress3 = &r3

		result = append(result, &l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func loadReceiptList(ctx context.Context, db *sql.DB, requestingAccount *gethcommon.Address, whereCondition string, whereParams []any, orderBy string, orderByParams []any) ([]*core.InternalReceipt, error) {
	if requestingAccount == nil {
		return nil, fmt.Errorf("you have to specify requestingAccount")
	}
	var queryParams []any

	query := "select b.hash, b.height, curr_tx.hash, curr_tx.idx, rec.post_state, rec.status, rec.gas_used, rec.effective_gas_price, rec.created_contract_address, tx_sender.address, tx_contr.address, curr_tx.type "
	query += baseReceiptJoin
	query += " WHERE 1=1 "

	// visibility
	query += " AND tx_sender.address = ? "
	queryParams = append(queryParams, requestingAccount.Bytes())

	query += whereCondition
	queryParams = append(queryParams, whereParams...)

	if len(orderBy) > 0 {
		query += orderBy
		queryParams = append(queryParams, orderByParams...)
	}

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := make([]*core.InternalReceipt, 0)

	empty := true
	for rows.Next() {
		empty = false
		r, err := onRowWithReceipt(rows)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, r)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if empty {
		return nil, errutil.ErrNotFound
	}
	return receipts, nil
}

func onRowWithReceipt(rows *sql.Rows) (*core.InternalReceipt, error) {
	r := core.InternalReceipt{}

	var txIndex *uint
	var blockHash, transactionHash *gethcommon.Hash
	var blockNumber *uint64
	res := []any{&blockHash, &blockNumber, &transactionHash, &txIndex, &r.PostState, &r.Status, &r.GasUsed, &r.EffectiveGasPrice, &r.CreatedContract, &r.From, &r.To, &r.TxType}

	err := rows.Scan(res...)
	if err != nil {
		return nil, fmt.Errorf("could not load receipt from db: %w", err)
	}

	r.BlockHash = *blockHash
	r.BlockNumber = big.NewInt(int64(*blockNumber))
	r.TxHash = *transactionHash
	r.TransactionIndex = *txIndex
	return &r, nil
}

// utility function that knows how to load relevant logs from the database together with a receipt
// returns either receipts with logs, or only logs
// this complexity is necessary to avoid executing multiple queries.
// todo always pass in the actual batch hashes because of reorgs, or make sure to clean up log entries from discarded batches
func loadReceiptsAndEventLogs(ctx context.Context, db *sql.DB, requestingAccount *gethcommon.Address, whereCondition string, whereParams []any, withReceipts bool) ([]*core.InternalReceipt, []*types.Log, error) {
	logsQuery := " et.event_sig, t1.topic, t2.topic, t3.topic, datablob, log_idx, b.hash, b.height, curr_tx.hash, curr_tx.idx, c.address "
	receiptQuery := " rec.post_state, rec.status, rec.gas_used, rec.effective_gas_price, rec.created_contract_address, tx_sender.address, tx_contr.address, curr_tx.type "

	query := "select " + logsQuery
	if withReceipts {
		query += "," + receiptQuery
	}
	query += baseReceiptJoin
	query += baseEventJoin

	var queryParams []any

	if requestingAccount != nil {
		// Add log visibility rules
		logsVisibQuery, logsVisibParams := logsVisibilityQuery(requestingAccount, withReceipts)
		query += logsVisibQuery
		queryParams = append(queryParams, logsVisibParams...)

		// add receipt visibility rules
		if withReceipts {
			receiptsVisibQuery, receiptsVisibParams := receiptsVisibilityQuery(requestingAccount)
			query += receiptsVisibQuery
			queryParams = append(queryParams, receiptsVisibParams...)
		}
	}

	query += whereCondition
	queryParams = append(queryParams, whereParams...)

	if withReceipts && requestingAccount != nil {
		// there is a corner case when a receipt has logs, but none are visible to the requester
		query += " UNION ALL "
		query += " select null, null, null, null, null, null, b.hash, b.height, curr_tx.hash, curr_tx.idx, null, " + receiptQuery
		query += baseReceiptJoin
		query += " where b.is_canonical=true "
		query += " AND tx_sender.address = ? "
		queryParams = append(queryParams, requestingAccount.Bytes())
		query += whereCondition
		queryParams = append(queryParams, whereParams...)
	}

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	receipts := make([]*core.InternalReceipt, 0)
	logList := make([]*types.Log, 0)

	empty := true
	for rows.Next() {
		empty = false
		r, l, err := onRowWithEventLogAndReceipt(rows, withReceipts)
		if err != nil {
			return nil, nil, err
		}
		receipts = append(receipts, r)
		logList = append(logList, l)
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}

	if withReceipts {
		if empty {
			return nil, nil, errutil.ErrNotFound
		}
		// group the logs manually to avoid complicating the db indexes
		result := groupReceiptAndLogs(receipts, logList)
		return result, nil, nil
	}
	return nil, logList, nil
}

func groupReceiptAndLogs(receipts []*core.InternalReceipt, logs []*types.Log) []*core.InternalReceipt {
	recMap := make(map[gethcommon.Hash]*core.InternalReceipt)
	logMap := make(map[gethcommon.Hash][]*types.Log)
	for _, r := range receipts {
		recMap[r.TxHash] = r
	}
	for _, log := range logs {
		logList := logMap[log.TxHash]
		if logList == nil {
			logList = make([]*types.Log, 0)
			logMap[log.TxHash] = logList
		}
		logMap[log.TxHash] = append(logList, log)
	}
	result := make([]*core.InternalReceipt, 0)
	for txHash, receipt := range recMap {
		receipt.Logs = logMap[txHash]
		result = append(result, receipt)
	}
	return result
}

func onRowWithEventLogAndReceipt(rows *sql.Rows, withReceipts bool) (*core.InternalReceipt, *types.Log, error) {
	l := types.Log{
		Topics: make([]gethcommon.Hash, 0),
	}
	r := core.InternalReceipt{}

	var t0, t1, t2, t3 []byte
	var logIndex, txIndex *uint
	var blockHash, transactionHash *gethcommon.Hash
	var address *gethcommon.Address
	var blockNumber *uint64
	res := []any{&t0, &t1, &t2, &t3, &l.Data, &logIndex, &blockHash, &blockNumber, &transactionHash, &txIndex, &address}

	if withReceipts {
		// when loading receipts, add the extra fields
		res = append(res, &r.PostState, &r.Status, &r.GasUsed, &r.EffectiveGasPrice, &r.CreatedContract, &r.From, &r.To, &r.TxType)
	}

	err := rows.Scan(res...)
	if err != nil {
		return nil, nil, fmt.Errorf("could not load log entry from db: %w", err)
	}

	if withReceipts {
		r.BlockHash = *blockHash
		r.BlockNumber = big.NewInt(int64(*blockNumber))
		r.TxHash = *transactionHash
		r.TransactionIndex = *txIndex
	}

	if logIndex != nil {
		l.Index, l.BlockHash, l.BlockNumber, l.TxHash, l.TxIndex = *logIndex, *blockHash, *blockNumber, *transactionHash, *txIndex
		if address != nil {
			l.Address = *address
		}
		for _, topic := range [][]byte{t0, t1, t2, t3} {
			if len(topic) > 0 {
				l.Topics = append(l.Topics, byteArrayToHash(topic))
			}
		}
	}

	if withReceipts {
		return &r, &l, nil
	}
	return nil, &l, nil
}

func receiptsVisibilityQuery(requestingAccount *gethcommon.Address) (string, []any) {
	// the visibility rules for the receipt:
	// - the sender can query
	// - anyone can query if the contract is transparent
	// - anyone who can view an event log should also be able to view the receipt
	query := " AND ( (e.id IS NOT NULL) OR (tx_sender.address = ?) OR (tx_contr.transparent=true) )"
	queryParams := []any{requestingAccount.Bytes()}
	return query, queryParams
}

// this function encodes the event log visibility rules
func logsVisibilityQuery(requestingAccount *gethcommon.Address, withReceipts bool) (string, []any) {
	acc := requestingAccount.Bytes()

	visibParams := make([]any, 0)

	visibQuery := "AND ("

	if withReceipts {
		// this condition only affects queries that return receipts that have no events logs
		visibQuery += " (e.id is NULL) "
	} else {
		visibQuery += " (1=0) "
	}

	// everyone can query config_public events
	visibQuery += " OR (et.config_public=true) "

	// For event logs that have no explicit configuration, an event is visible by all account owners whose addresses are used in any topic
	visibQuery += " OR (et.auto_visibility=true AND (et.auto_public=true OR eoa1.address=? OR eoa2.address=? OR eoa3.address=?)) "
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)

	// Configured events that are not public specify explicitly which event topics are addresses empowered to view that event
	visibQuery += " OR (" +
		"et.auto_visibility=false AND et.config_public=false AND " +
		"  (" +
		"       (et.topic1_can_view AND eoa1.address=?) " +
		"    OR (et.topic2_can_view AND eoa2.address=?) " +
		"    OR (et.topic3_can_view AND eoa3.address=?)" +
		"    OR (et.sender_can_view AND tx_sender.address=?)" +
		"  )" +
		")"
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)

	visibQuery += ") "
	return visibQuery, visibParams
}

func WriteEoa(ctx context.Context, dbTX *sql.Tx, sender gethcommon.Address) (uint64, error) {
	insert := "insert into externally_owned_account (address) values (?)"
	res, err := dbTX.ExecContext(ctx, insert, sender.Bytes())
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func ReadEoa(ctx context.Context, dbTx *sql.Tx, addr gethcommon.Address) (uint64, error) {
	row := dbTx.QueryRowContext(ctx, "select id from externally_owned_account where address = ?", addr.Bytes())

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return 0, errutil.ErrNotFound
		}
		return 0, err
	}

	return id, nil
}

func WriteContractConfig(ctx context.Context, dbTX *sql.Tx, contractAddress gethcommon.Address, eoaId uint64, cfg *core.ContractVisibilityConfig, txId uint64) (*uint64, error) {
	insert := "insert into contract (address, creator, auto_visibility, transparent, tx) values (?,?,?,?,?)"
	res, err := dbTX.ExecContext(ctx, insert, contractAddress.Bytes(), eoaId, cfg.AutoConfig, cfg.Transparent, txId)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	v := uint64(id)
	return &v, nil
}

func ReadContractByAddress(ctx context.Context, dbTx *sql.Tx, addr gethcommon.Address) (*Contract, error) {
	row := dbTx.QueryRowContext(ctx, "select c.id, c.address, c.auto_visibility, c.transparent, eoa.address from contract c join externally_owned_account eoa on c.creator=eoa.id where c.address = ?", addr.Bytes())

	var c Contract
	err := row.Scan(&c.Id, &c.Address, &c.AutoVisibility, &c.Transparent, &c.Creator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func byteArrayToHash(b []byte) gethcommon.Hash {
	result := gethcommon.Hash{}
	result.SetBytes(b)
	return result
}
