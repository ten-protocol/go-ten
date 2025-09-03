package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/jmoiron/sqlx"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/merkle"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	selectExtRollup         = "SELECT ext_rollup from rollup_host r join block_host b on r.compression_block=b.id "
	selectLatestExtRollup   = "SELECT ext_rollup FROM rollup_host ORDER BY time_stamp DESC LIMIT 1"
	selectLatestRollupCount = "SELECT id FROM rollup_host ORDER BY id DESC LIMIT 1"
	selectRollupBatches     = "SELECT b.sequence, b.hash, b.height, b.ext_batch FROM rollup_host r JOIN batch_host b ON b.sequence >= r.start_seq AND b.sequence <= r.end_seq AND b.sequence IS NOT NULL"
	selectRollups           = "SELECT rh.id, rh.hash, rh.start_seq, rh.end_seq, rh.time_stamp, rh.ext_rollup, bh.hash FROM rollup_host rh join block_host bh on rh.compression_block=bh.id "

	// SQL statements that need placeholder conversion
	insertRollup             = "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values (?,?,?,?,?,?) RETURNING id"
	insertCrossChainMessage  = "INSERT INTO cross_chain_message_host (message_hash, message_type, rollup_id) values (?,?,?)"
	selectRollupIdByMessage  = "SELECT rollup_id FROM cross_chain_message_host WHERE message_hash = ?"
	selectMessagesByRollupId = "SELECT message_hash, message_type FROM cross_chain_message_host WHERE rollup_id = ?"

	// WHERE clause patterns
	whereRollupHash     = " WHERE r.hash = ?"
	whereBlockHash      = " WHERE b.hash = ?"
	whereRollupHostHash = " WHERE rh.hash = ?"
	whereSeqBetween     = " WHERE ? BETWEEN start_seq AND end_seq"
)

// AddRollup adds a rollup to the DB
func AddRollup(dbtx *dbTransaction, db HostDB, rollup *common.ExtRollup, extMetadata *common.ExtRollupMetadata, metadata *common.PublicRollupMetadata, block *types.Header) error {
	extRollup, err := rlp.EncodeToBytes(rollup)
	if err != nil {
		return fmt.Errorf("could not encode rollup: %w", err)
	}

	var blockId int
	reboundSelectBlockId := db.GetSQLDB().Rebind(selectBlockId)
	err = dbtx.Tx.QueryRow(reboundSelectBlockId, block.Hash().Bytes()).Scan(&blockId)
	if err != nil {
		return fmt.Errorf("could not read block id: %w", err)
	}

	// Use QueryRow instead of Exec to retrieve the id directly.
	var rollupId int64
	reboundInsertRollup := db.GetSQLDB().Rebind(insertRollup)
	err = dbtx.Tx.QueryRow(reboundInsertRollup,
		rollup.Header.Hash().Bytes(),  // hash
		rollup.Header.FirstBatchSeqNo, // first batch sequence
		rollup.Header.LastBatchSeqNo,  // last batch sequence
		metadata.StartTime,            // timestamp
		extRollup,                     // rollup blob
		blockId,                       // l1 block hash
	).Scan(&rollupId)
	if err != nil {
		if IsRowExistsError(err) {
			return errutil.ErrAlreadyExists
		}
		return fmt.Errorf("could not insert rollup: %w", err)
	}

	if len(extMetadata.CrossChainTree) == 0 {
		return nil
	}

	tree, err := merkle.UnmarshalCrossChainTree(extMetadata.CrossChainTree)
	if err != nil {
		return err
	}

	for _, message := range tree {
		reboundInsertCrossChainMessage := db.GetSQLDB().Rebind(insertCrossChainMessage)
		_, err = dbtx.Tx.Exec(reboundInsertCrossChainMessage,
			message[1].(gethcommon.Hash).Bytes(),
			message[0],
			rollupId,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRollupListing returns latest rollups given a pagination.
// For example, offset 1, size 10 will return the latest 11-20 rollups.
func GetRollupListing(db HostDB, pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	orderQuery := " ORDER BY rh.id DESC "

	// Handle pagination with Rebind
	var paginationQuery string
	driverName := db.GetSQLDB().DriverName()
	if sqlx.BindType(driverName) == sqlx.QUESTION {
		paginationQuery = " LIMIT ? OFFSET ?"
	} else {
		paginationQuery = " LIMIT $1 OFFSET $2"
	}

	query := selectRollups + orderQuery + paginationQuery

	rows, err := db.GetSQLDB().Query(query, int64(pagination.Size), int64(pagination.Offset))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %s - %w", query, err)
	}
	defer rows.Close()
	var rollups []common.PublicRollup

	for rows.Next() {
		var id, startSeq, endSeq, timeStamp int
		var hash, extRollup, compressionBlock []byte

		var rollup common.PublicRollup
		err = rows.Scan(&id, &hash, &startSeq, &endSeq, &timeStamp, &extRollup, &compressionBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to scan query %s - %w", query, err)
		}

		extRollupDecoded := new(common.ExtRollup)
		if err := rlp.DecodeBytes(extRollup, extRollupDecoded); err != nil {
			return nil, fmt.Errorf("could not decode rollup header. Cause: %w", err)
		}

		rollup = common.PublicRollup{
			ID:        big.NewInt(int64(id)),
			Hash:      bytesToHexString(hash),
			FirstSeq:  big.NewInt(int64(startSeq)),
			LastSeq:   big.NewInt(int64(endSeq)),
			Timestamp: uint64(timeStamp),
			Header:    extRollupDecoded.Header,
			L1Hash:    bytesToHexString(compressionBlock),
		}
		rollups = append(rollups, rollup)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// TODO @will we will want to cache this value in the future
	totalRollups, err := fetchTotalRollups(db.GetSQLDB())
	if err != nil {
		return nil, fmt.Errorf("could not fetch total rollups. Cause: %w", err)
	}

	return &common.RollupListingResponse{
		RollupsData: rollups,
		Total:       totalRollups.Uint64(),
	}, nil
}

func GetExtRollup(db HostDB, hash gethcommon.Hash) (*common.ExtRollup, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereRollupHash)
	return fetchExtRollup(db.GetSQLDB(), reboundWhereQuery, hash.Bytes())
}

// GetRollupHeader returns the rollup with the given hash.
func GetRollupHeader(db HostDB, hash gethcommon.Hash) (*common.RollupHeader, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereRollupHash)
	return fetchRollupHeader(db.GetSQLDB(), reboundWhereQuery, hash.Bytes())
}

// GetRollupHeaderByBlock returns the rollup for the given block
func GetRollupHeaderByBlock(db HostDB, blockHash gethcommon.Hash) (*common.RollupHeader, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereBlockHash)
	return fetchRollupHeader(db.GetSQLDB(), reboundWhereQuery, blockHash)
}

// GetLatestRollup returns the latest rollup ordered by timestamp
func GetLatestRollup(db HostDB) (*common.RollupHeader, error) {
	extRollup, err := fetchHeadRollup(db.GetSQLDB())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch head rollup: %w", err)
	}
	return extRollup.Header, nil
}

func GetRollupByHash(db HostDB, rollupHash gethcommon.Hash) (*common.PublicRollup, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereRollupHostHash)
	return fetchPublicRollup(db.GetSQLDB(), reboundWhereQuery, rollupHash.Bytes())
}

func GetRollupBySeqNo(db HostDB, seqNo uint64) (*common.PublicRollup, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereSeqBetween)
	return fetchPublicRollup(db.GetSQLDB(), reboundWhereQuery, seqNo)
}

func GetCrossChainMessagesTree(db HostDB, messageHash gethcommon.Hash) ([][]interface{}, error) {
	// First get the rollupID for this message hash
	var rollupID int64
	reboundMessageQuery := db.GetSQLDB().Rebind(selectRollupIdByMessage)
	err := db.GetSQLDB().QueryRow(reboundMessageQuery, messageHash.Bytes()).Scan(&rollupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch rollup ID for message: %w", err)
	}

	// Get all messages with the same rollupID
	reboundMessagesQuery := db.GetSQLDB().Rebind(selectMessagesByRollupId)
	rows, err := db.GetSQLDB().Query(reboundMessagesQuery, rollupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cross chain messages: %w", err)
	}
	defer rows.Close()

	var messages [][]interface{}
	for rows.Next() {
		var messageHash []byte
		var messageType string
		if err := rows.Scan(&messageHash, &messageType); err != nil {
			return nil, fmt.Errorf("failed to scan cross chain message row: %w", err)
		}
		messages = append(messages, []interface{}{messageType, messageHash})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func GetRollupBatches(db HostDB, rollupHash gethcommon.Hash, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereRollupHash)
	orderQuery := " ORDER BY b.sequence DESC "

	// TODO @will quick fix to unblock main
	var paginationQuery string
	driverName := db.GetSQLDB().DriverName()
	if sqlx.BindType(driverName) == sqlx.QUESTION {
		paginationQuery = " LIMIT ? OFFSET ?"
	} else {
		// PostgreSQL uses $1, $2, $3,
		paginationQuery = " LIMIT $2 OFFSET $3"
	}
	query := selectRollupBatches + reboundWhereQuery + orderQuery + paginationQuery

	countQuery := "SELECT COUNT(*) FROM rollup_host r JOIN batch_host b ON b.sequence BETWEEN r.start_seq AND r.end_seq" + reboundWhereQuery
	var total uint64
	err := db.GetSQLDB().QueryRow(countQuery, rollupHash.Bytes()).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	rows, err := db.GetSQLDB().Query(query, rollupHash.Bytes(), int64(pagination.Size), int64(pagination.Offset))
	if err != nil {
		return nil, fmt.Errorf("query execution for select rollup batches failed: %w", err)
	}
	defer rows.Close()

	var batches []common.PublicBatch
	for rows.Next() {
		var (
			sequenceInt64 int
			fullHash      gethcommon.Hash
			heightInt64   int
			extBatch      []byte
		)
		err := rows.Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch rollup batches: %w", err)
		}
		var b common.ExtBatch
		err = rlp.DecodeBytes(extBatch, &b)
		if err != nil {
			return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
		}

		batch := common.PublicBatch{
			SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
			FullHash:         fullHash,
			Height:           new(big.Int).SetInt64(int64(heightInt64)),
			Header:           b.Header,
			EncryptedTxBlob:  b.EncryptedTxBlob,
			TxHashes:         b.TxHashes,
		}
		batches = append(batches, batch)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &common.BatchListingResponse{
		BatchesData: batches,
		Total:       total,
	}, nil
}

func fetchRollupHeader(db *sqlx.DB, whereQuery string, args ...any) (*common.RollupHeader, error) {
	rollup, err := fetchExtRollup(db, whereQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ext rollup - %w", err)
	}
	return rollup.Header, nil
}

func fetchExtRollup(db *sqlx.DB, whereQuery string, args ...any) (*common.ExtRollup, error) {
	var rollupBlob []byte
	query := selectExtRollup + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&rollupBlob)
	} else {
		err = db.QueryRow(query).Scan(&rollupBlob)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	var rollup common.ExtRollup
	err = rlp.DecodeBytes(rollupBlob, &rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup: %w", err)
	}

	return &rollup, nil
}

func fetchHeadRollup(db *sqlx.DB) (*common.ExtRollup, error) {
	var extRollup []byte
	err := db.QueryRow(selectLatestExtRollup).Scan(&extRollup)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	var rollup common.ExtRollup
	err = rlp.DecodeBytes(extRollup, &rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup: %w", err)
	}

	return &rollup, nil
}

func fetchTotalRollups(db *sqlx.DB) (*big.Int, error) {
	var total int
	err := db.QueryRow(selectLatestRollupCount).Scan(&total)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("failed to fetch rollup latest rollup ID: %w", err)
	}

	bigTotal := big.NewInt(int64(total))
	return bigTotal, nil
}

func fetchPublicRollup(db *sqlx.DB, whereQuery string, args ...any) (*common.PublicRollup, error) {
	query := selectRollups + whereQuery
	var rollup common.PublicRollup
	var hash, extRollup, compressionblock []byte
	var id, firstSeq, lastSeq, timestamp int

	err := db.QueryRow(query, args...).Scan(
		&id,
		&hash,
		&firstSeq,
		&lastSeq,
		&timestamp,
		&extRollup,
		&compressionblock,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	rollup.ID = big.NewInt(int64(id))
	rollup.Hash = bytesToHexString(hash)
	rollup.FirstSeq = big.NewInt(int64(firstSeq))
	rollup.LastSeq = big.NewInt(int64(lastSeq))
	rollup.Timestamp = uint64(timestamp)
	rollup.L1Hash = bytesToHexString(compressionblock)

	extRollupDecoded := new(common.ExtRollup)
	if err := rlp.DecodeBytes(extRollup, extRollupDecoded); err != nil {
		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
	}

	rollup.Header = extRollupDecoded.Header
	return &rollup, nil
}
