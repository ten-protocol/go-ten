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
	"github.com/ten-protocol/go-ten/go/common/tracers"
)

const (
	baseEventsJoin = "from event_log e " +
		"join receipt rec on e.receipt=rec.id" +
		"	join tx on rec.tx=tx.id " +
		"	join batch b on rec.batch=b.sequence " +
		"join event_type et on e.event_type=et.id " +
		"	join contract c on et.contract=c.id " +
		"left join event_topic t1 on e.topic1=t1.id and et.id=t1.event_type " +
		"   left join externally_owned_account eoa1 on t1.rel_address=eoa1.id " +
		"left join event_topic t2 on e.topic2=t2.id and et.id=t2.event_type " +
		"   left join externally_owned_account eoa2 on t2.rel_address=eoa2.id " +
		"left join event_topic t3 on e.topic3=t3.id and et.id=t1.event_type " +
		"   left join externally_owned_account eoa3 on t3.rel_address=eoa3.id " +
		"where b.is_canonical=true "
)

func WriteEventType(ctx context.Context, dbTX *sql.Tx, et *EventType) (uint64, error) {
	res, err := dbTX.ExecContext(ctx, "insert into event_type (contract, event_sig, auto_visibility,auto_public, public, topic1_can_view, topic2_can_view, topic3_can_view, sender_can_view) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		et.Contract.Id, et.EventSignature.Bytes(), et.AutoVisibility, et.AutoPublic, et.Public, et.Topic1CanView, et.Topic2CanView, et.Topic3CanView, et.SenderCanView)
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
	var et EventType = EventType{Contract: contract}
	err := dbTX.QueryRowContext(ctx,
		"select id, event_sig, auto_visibility, auto_public, public, topic1_can_view, topic2_can_view, topic3_can_view, sender_can_view from event_type where contract=? and event_sig=?",
		contract.Id, eventSignature.Bytes(),
	).Scan(&et.Id, &et.EventSignature, &et.AutoVisibility, &et.AutoPublic, &et.Public, &et.Topic1CanView, &et.Topic2CanView, &et.Topic3CanView, &et.SenderCanView)
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

func ReadEventTopic(ctx context.Context, dbTX *sql.Tx, topic []byte, eventTypeId uint64) (uint64, *uint64, error) {
	var id uint64
	var address *uint64
	err := dbTX.QueryRowContext(ctx,
		"select id, rel_address from event_topic where topic=? and event_type=?", topic, eventTypeId).Scan(&id, &address)
	if errors.Is(err, sql.ErrNoRows) {
		// make sure the error is converted to obscuro-wide not found error
		return 0, nil, errutil.ErrNotFound
	}
	return id, address, err
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

func FilterLogs(
	ctx context.Context,
	db *sql.DB,
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	batchHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
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

	return loadLogs(ctx, db, requestingAccount, query, queryParams)
}

func DebugGetLogs(ctx context.Context, db *sql.DB, txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	var queryParams []any

	query := "select eoa1.address, eoa2.address, eoa3.address, et.public, et.event_sig, t1.topic, t2.topic, t3.topic, datablob, b.hash, b.height, tx.hash, tx.idx, log_idx, c.address " +
		baseEventsJoin +
		" AND tx.hash = ? "

	queryParams = append(queryParams, txHash.Bytes())

	result := make([]*tracers.DebugLogs, 0)

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := tracers.DebugLogs{
			Log: types.Log{
				Topics: []gethcommon.Hash{},
			},
			LifecycleEvent: false,
		}

		var t0, t1, t2, t3 sql.NullString
		var relAddress1, relAddress2, relAddress3 []byte
		err = rows.Scan(
			&relAddress1,
			&relAddress2,
			&relAddress3,
			&l.LifecycleEvent,
			&t0, &t1, &t2, &t3,
			&l.Data,
			&l.BlockHash,
			&l.BlockNumber,
			&l.TxHash,
			&l.TxIndex,
			&l.Index,
			&l.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}

		for _, topic := range []sql.NullString{t0, t1, t2, t3} {
			if topic.Valid {
				l.Topics = append(l.Topics, stringToHash(topic))
			}
		}

		l.RelAddress1 = bytesToAddress(relAddress1)
		l.RelAddress2 = bytesToAddress(relAddress2)
		l.RelAddress3 = bytesToAddress(relAddress3)

		result = append(result, &l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func bytesToAddress(b []byte) *gethcommon.Address {
	if b != nil {
		addr := gethcommon.BytesToAddress(b)
		return &addr
	}
	return nil
}

// utility function that knows how to load relevant logs from the database
// todo always pass in the actual batch hashes because of reorgs, or make sure to clean up log entries from discarded batches
func loadLogs(ctx context.Context, db *sql.DB, requestingAccount *gethcommon.Address, whereCondition string, whereParams []any) ([]*types.Log, error) {
	if requestingAccount == nil { // todo - only restrict to lifecycle events if requesting==nil
		return nil, fmt.Errorf("logs can only be requested for an account")
	}

	result := make([]*types.Log, 0)
	query := "select et.event_sig, t1.topic, t2.topic, t3.topic, datablob, b.hash, b.height, tx.hash, tx.idx, log_idx, c.address" + " " + baseEventsJoin
	var queryParams []any

	// Add visibility rules
	visibQuery, visibParams := visibilityQuery(requestingAccount)

	query += visibQuery
	queryParams = append(queryParams, visibParams...)

	query += whereCondition
	queryParams = append(queryParams, whereParams...)

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := types.Log{
			Topics: make([]gethcommon.Hash, 0),
		}
		var t0, t1, t2, t3 []byte
		err = rows.Scan(&t0, &t1, &t2, &t3, &l.Data, &l.BlockHash, &l.BlockNumber, &l.TxHash, &l.TxIndex, &l.Index, &l.Address)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}

		for _, topic := range [][]byte{t0, t1, t2, t3} {
			if len(topic) > 0 {
				l.Topics = append(l.Topics, byteArrayToHash(topic))
			}
		}

		result = append(result, &l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

// this function encodes the event log visibility rules
func visibilityQuery(requestingAccount *gethcommon.Address) (string, []any) {
	acc := requestingAccount.Bytes()

	visibQuery := "AND ("
	visibParams := make([]any, 0)

	// everyone can query public events
	visibQuery += " et.public=true "

	// For event logs that have no explicit configuration, an event is visible be all account owners whose addresses are used in any topic
	visibQuery += " OR (et.auto_visibility=true AND (et.auto_public=true OR (eoa1.address=? OR eoa2.address=? OR eoa3.address=?))) "
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)
	visibParams = append(visibParams, acc)

	// Configured events that are not public specify explicitly which event topics are addresses empowered to view that event
	visibQuery += " OR (" +
		"et.auto_visibility=false AND et.public=false AND " +
		"  (" +
		"       (et.topic1_can_view AND eoa1.address=?) " +
		"    OR (et.topic2_can_view AND eoa2.address=?) " +
		"    OR (et.topic3_can_view AND eoa3.address=?)" +
		"  )" +
		")"
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

func WriteContractConfig(ctx context.Context, dbTX *sql.Tx, contractAddress gethcommon.Address, eoaId uint64, cfg *core.ContractVisibilityConfig) (*uint64, error) {
	insert := "insert into contract (address, creator, auto_visibility, transparent) values (?,?,?,?)"
	res, err := dbTX.ExecContext(ctx, insert, contractAddress.Bytes(), eoaId, cfg.AutoConfig, cfg.Transparent)
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
	row := dbTx.QueryRowContext(ctx, "select id, address, auto_visibility, transparent from contract where address = ?", addr.Bytes())

	var c Contract
	err := row.Scan(&c.Id, &c.Address, &c.AutoVisibility, &c.Transparent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func ReadContractCreator(ctx context.Context, db *sql.DB, address gethcommon.Address) (*gethcommon.Address, error) {
	row := db.QueryRowContext(ctx, "select eoa.address from contract c join externally_owned_account eoa on c.creator=eoa.id  where c.address = ?", address.Bytes())

	var eoaAddress gethcommon.Address
	err := row.Scan(&eoaAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return &eoaAddress, nil
}

func stringToHash(ns sql.NullString) gethcommon.Hash {
	value, err := ns.Value()
	if err != nil {
		return [32]byte{}
	}
	s, ok := value.(string)
	if !ok {
		return [32]byte{}
	}
	result := gethcommon.Hash{}
	result.SetBytes([]byte(s))
	return result
}

func byteArrayToHash(b []byte) gethcommon.Hash {
	result := gethcommon.Hash{}
	result.SetBytes(b)
	return result
}
