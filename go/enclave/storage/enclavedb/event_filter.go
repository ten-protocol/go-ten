package enclavedb

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	NR_TOPICS       = 3
	baseReceiptJoin = " from receipt rec " +
		"join batch b on rec.batch=b.sequence " +
		"join tx curr_tx on rec.tx=curr_tx.id " +
		"   join externally_owned_account tx_sender on curr_tx.sender_address=tx_sender.id " +
		"   left join contract tx_contr on curr_tx.contract=tx_contr.id "

	baseReceiptJoinWithViewer = " from receipt rec " +
		"left join receipt_viewer rv on rec.id=rv.receipt " +
		"join batch b on rec.batch=b.sequence " +
		"join tx curr_tx on rec.tx=curr_tx.id " +
		"   join externally_owned_account tx_sender on curr_tx.sender_address=tx_sender.id " +
		"   left join contract tx_contr on curr_tx.contract=tx_contr.id "

	personalTxCondition       = "(tx_sender.id = ? OR rv.eoa = ? OR curr_tx.to_eoa = ?)"
	personalTxConditionPublic = "(tx_sender.id = ? OR rv.eoa = ? OR curr_tx.to_eoa = ? OR rec.public)"

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

// FilterLogs - the event types will contain all contracts
func FilterLogs(ctx context.Context, stmtCache *PreparedStatementCache, requestingAccountId *uint64, fromBlock, toBlock *big.Int, batchHash *common.L2BatchHash, eventTypes []*EventType, topics FilterTopics) ([]*types.Log, error) {
	if len(eventTypes) == 0 {
		return nil, fmt.Errorf("should not happen. At least one event type must be specified")
	}
	var queries []string
	var params []any
	for _, eventType := range eventTypes {
		q, p := loadEventLogsQuery(requestingAccountId, fromBlock, toBlock, batchHash, eventType, topics)
		queries = append(queries, q)
		params = append(params, p...)
	}
	query := strings.Join(queries, " UNION ALL ")
	logs, err := executeEventLogsQuery(ctx, stmtCache, query, params)
	if err != nil {
		return nil, err
	}

	// the database returns an unsorted list of event logs.
	// we have to perform the sorting programmatically
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].Index < logs[j].Index
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})

	// remove duplicates
	logs = slices.CompactFunc(logs, func(a, b *types.Log) bool {
		return a.BlockNumber == b.BlockNumber && a.Index == b.Index && a.TxHash == b.TxHash
	})
	return logs, nil
}

func loadEventLogsQuery(requestingAccountId *uint64, fromBlock, toBlock *big.Int, batchHash *common.L2BatchHash, eventType *EventType, topics FilterTopics) (string, []any) {
	query := "select et.event_sig, t1.topic, t2.topic, t3.topic, datablob, log_idx, b.hash, b.height, curr_tx.hash, curr_tx.idx, c.address " +
		"from batch b " +
		"join receipt rec on rec.batch=b.sequence " +
		"	join tx curr_tx on rec.tx=curr_tx.id " +
		"   	join externally_owned_account tx_sender on curr_tx.sender_address=tx_sender.id " +
		" 	join event_log e on e.receipt=rec.id " +
		" 		join event_type et on e.event_type=et.id " +
		"	 		join contract c on et.contract=c.id "

	for i := 1; i <= NR_TOPICS; i++ {
		joinType := "left"
		/*
			// todo - consider this more carefully
			// When we know we have to search for topics on a position, we can use a straight join.
			// If there is nothing to join to, the query will end early, which is correct.
			if !topics.HasTopicsOnPos(i) {
				joinType = "left"
			}
		*/
		query += fmt.Sprintf(" %s join event_topic t%d on e.topic%d=t%d.id ", joinType, i, i, i)
		if eventType.IsTopicRelevant(i) || eventType.AutoVisibility {
			// we join with `externally_owned_account` only if we need to filter visibility on this topic
			query += fmt.Sprintf("   %s join externally_owned_account eoa%d on t%d.rel_address=eoa%d.id ", joinType, i, i, i)
		}
	}

	query += " WHERE b.is_canonical=true "

	var queryParams []any
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

	for i := 1; i <= NR_TOPICS; i++ {
		if topics.HasTopicsOnPos(i) {
			topicsRow := topics.TopicsOnPos(i)
			valuesIn := "IN (" + repeat("?", ",", len(topicsRow)) + ")"
			query += fmt.Sprintf(" AND t%d.topic %s", i, valuesIn)
			for _, topicHash := range topicsRow {
				queryParams = append(queryParams, topicHash.Bytes())
			}
		}
	}

	query += " AND et.id = ? "
	queryParams = append(queryParams, eventType.Id)

	// for non-public events add visibility rules
	if !eventType.IsPublic() {
		query += " AND ( (1=0) "

		if eventType.AutoVisibility {
			query += " OR eoa1.id=? OR eoa2.id=? OR eoa3.id=? "
			queryParams = append(queryParams, requestingAccountId)
			queryParams = append(queryParams, requestingAccountId)
			queryParams = append(queryParams, requestingAccountId)
		}

		if eventType.SenderCanView != nil && *eventType.SenderCanView {
			query += " OR tx_sender.id=? "
			queryParams = append(queryParams, requestingAccountId)
		}

		for i := 1; i <= NR_TOPICS; i++ {
			if eventType.IsTopicRelevant(i) {
				query += fmt.Sprintf(" OR eoa%d.id=? ", i)
				queryParams = append(queryParams, requestingAccountId)
			}
		}
		query += ") "
	}

	return query, queryParams
}

func executeEventLogsQuery(ctx context.Context, stmtCache *PreparedStatementCache, query string, queryParams []any) ([]*types.Log, error) {
	stmt, err := stmtCache.GetOrPrepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, queryParams...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	logList := make([]*types.Log, 0)

	for rows.Next() {
		_, l, err := onRowWithEventLogAndReceipt(rows, false)
		if err != nil {
			return nil, err
		}
		logList = append(logList, l)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return logList, nil
}

func DebugGetLogs(ctx context.Context, db *sqlx.DB, fromBlock *big.Int, toBlock *big.Int, address gethcommon.Address, eventSig gethcommon.Hash) ([]*common.DebugLogVisibility, error) {
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

func loadPersonalTxs(ctx context.Context, db *sqlx.DB, requestingAccountId *uint64, showPublic bool, whereCondition string, whereParams []any, orderBy string, orderByParams []any) ([]*core.InternalReceipt, error) {
	if requestingAccountId == nil {
		return nil, fmt.Errorf("you have to specify requestingAccount")
	}
	var queryParams []any

	query := "select b.hash, b.height, curr_tx.hash, curr_tx.idx, rec.post_state, rec.status, rec.gas_used, rec.effective_gas_price, rec.created_contract_address, tx_sender.address, tx_contr.address, curr_tx.type "
	query += baseReceiptJoinWithViewer

	// visibility
	query += " WHERE "

	if showPublic {
		query += personalTxConditionPublic
	} else {
		query += personalTxCondition
	}

	queryParams = append(queryParams, *requestingAccountId)
	queryParams = append(queryParams, *requestingAccountId)
	queryParams = append(queryParams, *requestingAccountId)

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

	r.Logs = make([]*types.Log, 0)
	return &r, nil
}

// utility function that knows how to load relevant logs from the database together with a receipt
// returns either receipts with logs, or only logs
// this complexity is necessary to avoid executing multiple queries.
// todo always pass in the actual batch hashes because of reorgs, or make sure to clean up log entries from discarded batches
func loadReceiptsAndEventLogs(ctx context.Context, stmtCache *PreparedStatementCache, requestingAccountId *uint64, whereCondition string, whereParams []any, withReceipts bool) ([]*core.InternalReceipt, []*types.Log, error) {
	logsQuery := " et.event_sig, t1.topic, t2.topic, t3.topic, datablob, log_idx, b.hash, b.height, curr_tx.hash, curr_tx.idx, c.address "
	receiptQuery := " rec.post_state, rec.status, rec.gas_used, rec.effective_gas_price, rec.created_contract_address, tx_sender.address, tx_contr.address, curr_tx.type "

	query := "select " + logsQuery
	if withReceipts {
		query += "," + receiptQuery
	}
	query += baseReceiptJoin
	query += baseEventJoin

	var queryParams []any

	if requestingAccountId != nil {
		// Add log visibility rules
		logsVisibQuery, logsVisibParams := logsVisibilityQuery(*requestingAccountId, withReceipts)
		query += logsVisibQuery
		queryParams = append(queryParams, logsVisibParams...)

		// add receipt visibility rules
		if withReceipts {
			receiptsVisibQuery, receiptsVisibParams := receiptsVisibilityQuery(*requestingAccountId)
			query += receiptsVisibQuery
			queryParams = append(queryParams, receiptsVisibParams...)
		}
	}

	query += whereCondition
	queryParams = append(queryParams, whereParams...)

	if withReceipts && requestingAccountId != nil {
		// there is a corner case when a receipt has logs, but none are visible to the requester
		query += " UNION ALL "
		query += " select null, null, null, null, null, null, b.hash, b.height, curr_tx.hash, curr_tx.idx, null, " + receiptQuery
		query += baseReceiptJoin
		query += " where b.is_canonical=true "
		query += " AND tx_sender.id = ? "
		queryParams = append(queryParams, requestingAccountId)
		query += whereCondition
		queryParams = append(queryParams, whereParams...)
	}

	stmt, err := stmtCache.GetOrPrepare(query)
	if err != nil {
		return nil, nil, fmt.Errorf("could not prepare query: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, queryParams...)
	if err != nil {
		return nil, nil, fmt.Errorf("could not execute query: %w", err)
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

func receiptsVisibilityQuery(requestingAccountId uint64) (string, []any) {
	// the visibility rules for the receipt:
	// - the sender can query
	// - anyone can query if the contract is transparent
	// - anyone who can view an event log should also be able to view the receipt
	query := " AND ( (e.id IS NOT NULL) OR (tx_sender.id = ?) OR (tx_contr.transparent=true) )"
	queryParams := []any{requestingAccountId}
	return query, queryParams
}

// this function encodes the event log visibility rules
func logsVisibilityQuery(requestingAccountId uint64, withReceipts bool) (string, []any) {
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
	visibQuery += " OR (et.auto_visibility=true AND (et.auto_public=true OR eoa1.id=? OR eoa2.id=? OR eoa3.id=?)) "
	visibParams = append(visibParams, requestingAccountId)
	visibParams = append(visibParams, requestingAccountId)
	visibParams = append(visibParams, requestingAccountId)

	// Configured events that are not public specify explicitly which event topics are addresses empowered to view that event
	visibQuery += " OR (" +
		"et.auto_visibility=false AND et.config_public=false AND " +
		"  (" +
		"       (et.topic1_can_view AND eoa1.id=?) " +
		"    OR (et.topic2_can_view AND eoa2.id=?) " +
		"    OR (et.topic3_can_view AND eoa3.id=?)" +
		"    OR (et.sender_can_view AND tx_sender.id=?)" +
		"  )" +
		")"
	visibParams = append(visibParams, requestingAccountId)
	visibParams = append(visibParams, requestingAccountId)
	visibParams = append(visibParams, requestingAccountId)
	visibParams = append(visibParams, requestingAccountId)

	visibQuery += ") "
	return visibQuery, visibParams
}
