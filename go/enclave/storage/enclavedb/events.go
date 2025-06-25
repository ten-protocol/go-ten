package enclavedb

import (
	"context"
	"database/sql"
	"errors"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

func WriteEventType(ctx context.Context, dbTX *sqlx.Tx, et *EventType) (uint64, error) {
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

func ReadEventTypesForContract(ctx context.Context, dbTX *sqlx.Tx, contractId uint64) ([]*EventType, error) {
	rows, err := dbTX.QueryContext(ctx,
		"select id, event_sig, auto_visibility, auto_public, config_public, topic1_can_view, topic2_can_view, topic3_can_view, sender_can_view "+
			"from event_type where contract=?", contractId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventTypes []*EventType
	for rows.Next() {
		et := EventType{
			Contract: &Contract{Id: contractId},
		}
		err := rows.Scan(&et.Id, &et.EventSignature, &et.AutoVisibility, &et.AutoPublic, &et.ConfigPublic, &et.Topic1CanView, &et.Topic2CanView, &et.Topic3CanView, &et.SenderCanView)
		if err != nil {
			return nil, err
		}
		eventTypes = append(eventTypes, &et)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return eventTypes, nil
}

func WriteEventTopic(ctx context.Context, dbTX *sqlx.Tx, topic *gethcommon.Hash, addressId *uint64, eventTypeId uint64) (uint64, error) {
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

func UpdateEventTypeAutoPublic(ctx context.Context, dbTx *sqlx.Tx, etId uint64, isPublic bool) error {
	_, err := dbTx.ExecContext(ctx, "update event_type set auto_public=? where id=?", isPublic, etId)
	return err
}

func ReadEventTopic(ctx context.Context, dbTX *sqlx.Tx, topic []byte, eventTypeId uint64) (*EventTopic, error) {
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

func WriteEventLog(ctx context.Context, dbTX *sqlx.Tx, eventTypeId uint64, userTopics []*EventTopic, data []byte, logIdx uint, execTx uint64) error {
	var t1, t2, t3 *uint64
	if len(userTopics) > 0 && userTopics[0] != nil {
		t1 = &(userTopics[0].Id)
	}
	if len(userTopics) > 1 && userTopics[1] != nil {
		t2 = &(userTopics[1].Id)
	}
	if len(userTopics) > 2 && userTopics[2] != nil {
		t3 = &(userTopics[2].Id)
	}
	_, err := dbTX.ExecContext(ctx, "insert into event_log (event_type, topic1, topic2, topic3, datablob, log_idx, receipt) values (?,?,?,?,?,?,?)",
		eventTypeId, t1, t2, t3, data, logIdx, execTx)
	return err
}

func WriteEoa(ctx context.Context, dbTX *sqlx.Tx, sender gethcommon.Address) (uint64, error) {
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

func ReadEoa(ctx context.Context, dbTx *sqlx.Tx, addr gethcommon.Address) (uint64, error) {
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

func WriteContractConfig(ctx context.Context, dbTX *sqlx.Tx, contractAddress gethcommon.Address, eoaId uint64, cfg *core.ContractVisibilityConfig, txId uint64) (*uint64, error) {
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

func ReadContractByAddress(ctx context.Context, dbTx *sqlx.Tx, addr gethcommon.Address) (*Contract, error) {
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

func ReadEventTypes(ctx context.Context, dbTX *sqlx.Tx, eventSignature gethcommon.Hash) ([]*EventType, error) {
	rows, err := dbTX.QueryContext(ctx,
		"select c.id, c.address, c.auto_visibility, c.transparent, eoa.address, et.id, et.event_sig, et.auto_visibility, et.auto_public, et.config_public, et.topic1_can_view, et.topic2_can_view, et.topic3_can_view, et.sender_can_view "+
			"from event_type et join contract c on et.contract=c.id join externally_owned_account eoa on c.creator=eoa.id "+
			"where et.event_sig=?",
		eventSignature.Bytes(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ets []*EventType
	for rows.Next() {
		var c Contract
		et := EventType{Contract: &c}
		err := rows.Scan(&c.Id, &c.Address, &c.AutoVisibility, &c.Transparent, &c.Creator, &et.Id, &et.EventSignature, &et.AutoVisibility, &et.AutoPublic, &et.ConfigPublic, &et.Topic1CanView, &et.Topic2CanView, &et.Topic3CanView, &et.SenderCanView)
		if err != nil {
			return nil, err
		}
		ets = append(ets, &et)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ets, nil
}
