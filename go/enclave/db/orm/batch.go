package orm

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	obscurosql "github.com/obscuronet/go-obscuro/go/enclave/db/sql"
)

const (
	bodyInsert         = "replace into batch_body values (?,?)"
	bInsert            = "insert into batch values (?,?,?,?,?,?,?,?,?)"
	updateNonCanonical = "update batch set is_canonical=? where hash=?"
	txInsert           = "insert into tx values "
	txInsertValue      = "(?,?,?,?,?,?)"
	selectBatch        = "select b.header, bb.content from batch b join batch_body bb on b.body=bb.hash"
	selectHeader       = "select b.header from batch b"
)

func WriteBatch(dbtx *obscurosql.Batch, batch *core.Batch) error {
	bodyHash := batch.Header.TxHash.Bytes()

	body, err := rlp.EncodeToBytes(batch.Transactions)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	header, err := rlp.EncodeToBytes(batch.Header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}

	dbtx.ExecuteSQL(bodyInsert, bodyHash, body)

	dbtx.ExecuteSQL(bInsert,
		batch.Hash().Bytes(),
		batch.Header.ParentHash.Bytes(),
		batch.Header.SequencerOrderNo.Uint64(),
		batch.Header.Number.Uint64(),
		true,
		string(header),
		bodyHash,
		batch.Header.L1Proof.Bytes(),
		"", // todo
	)

	if len(batch.Transactions) > 0 {
		insert := txInsert + strings.Repeat(txInsertValue+",", len(batch.Transactions))
		args := make([]any, 6*len(batch.Transactions))
		for _, transaction := range batch.Transactions {
			args = append(args, transaction.Hash())
			args = append(args, transaction.Hash())
		}
		dbtx.ExecuteSQL(insert[0:len(insert)-1], args...)
	}

	return nil
}

func MarkNonCanonicalBatch(dbtx *obscurosql.Batch, batchHash gethcommon.Hash) {
	// updateNonCanonical
}

func FindBatchBySeqNo(db *sql.DB, seqNo uint64) (*core.Batch, error) {
	var header string
	var body []byte
	query := selectBatch + " where sequence=?"
	err := db.QueryRow(query, seqNo).Scan(&header, &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	h := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	txs := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(body), txs); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions. Cause: %w", err)
	}

	return &core.Batch{
		Header:       h,
		Transactions: *txs,
	}, nil
}

func FetchBatch(db *sql.DB, hash common.L2BatchHash) (*core.Batch, error) {
	var header string
	var body []byte
	query := selectBatch + " where b.hash=?"
	err := db.QueryRow(query, hash.Bytes()).Scan(&header, &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	h := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	txs := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(body), txs); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions. Cause: %w", err)
	}

	return &core.Batch{
		Header:       h,
		Transactions: *txs,
	}, nil
}

func FetchCanonicalBatchByHeight(db *sql.DB, height uint64) (*core.Batch, error) {
	var header string
	var body []byte
	query := selectBatch + " where b.height=? and is_canonical"
	err := db.QueryRow(query, height).Scan(&header, &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	h := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	txs := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(body), txs); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions. Cause: %w", err)
	}

	return &core.Batch{
		Header:       h,
		Transactions: *txs,
	}, nil
}

func ReadBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	var header string
	query := selectHeader + " where hash=?"
	err := db.QueryRow(query, hash.Bytes()).Scan(&header)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	h := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}

	return h, nil
}

func FetchHeadBatch(db *sql.DB) (*core.Batch, error) {
	var header string
	var body []byte
	query := selectBatch + " where b.height=(select max(b1.height) from batch b1) and is_canonical"
	err := db.QueryRow(query).Scan(&header, &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	h := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	txs := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(body), txs); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions. Cause: %w", err)
	}

	return &core.Batch{
		Header:       h,
		Transactions: *txs,
	}, nil
}

func FetchCurrentSequencerNo(db *sql.DB) (*big.Int, error) {
	var seq int64
	query := "select max(sequence) from batch"
	err := db.QueryRow(query).Scan(&seq)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return big.NewInt(seq), nil
}
