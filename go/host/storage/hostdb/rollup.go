package hostdb

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	selectExtRollup     = "SELECT ext_rollup from rollup_host r"
	selectLatestRollup  = "SELECT ext_rollup FROM rollup_host ORDER BY time_stamp DESC LIMIT 1"
	selectRollupBatches = "SELECT b.sequence, b.hash, b.full_hash, b.height, b.ext_batch FROM rollup_host r JOIN batch_host b ON r.start_seq <= b.sequence AND r.end_seq >= b.sequence"
	selectRollups       = "SELECT id, hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block FROM rollup_host ORDER BY id DESC "
)

// AddRollup adds a rollup to the DB
func AddRollup(dbtx *dbTransaction, statements *SQLStatements, rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error {
	extRollup, err := rlp.EncodeToBytes(rollup)
	if err != nil {
		return fmt.Errorf("could not encode rollup: %w", err)
	}

	_, err = dbtx.tx.Exec(statements.InsertRollup,
		truncTo16(rollup.Header.Hash()),      // short hash
		metadata.FirstBatchSequence.Uint64(), // first batch sequence
		rollup.Header.LastBatchSeqNo,         // last batch sequence
		metadata.StartTime,                   // timestamp
		extRollup,                            // rollup blob
		block.Hash(),                         // l1 block hash
	)
	if err != nil {
		return fmt.Errorf("could not insert rollup. Cause: %w", err)
	}
	return nil
}

// GetRollupListing returns latest rollups given a pagination.
// For example, offset 1, size 10 will return the latest 11-20 rollups.
func GetRollupListing(db HostDB, pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	query := selectRollups + db.GetSQLStatement().Pagination

	rows, err := db.GetSQLDB().Query(query, pagination.Size, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rollups []common.PublicRollup

	for rows.Next() {
		var id, startSeq, endSeq, timeStamp int
		var hash, extRollup, compressionBlock []byte

		var rollup common.PublicRollup
		err = rows.Scan(&id, &hash, &startSeq, &endSeq, &timeStamp, &extRollup, &compressionBlock)
		if err != nil {
			return nil, err
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

	return &common.RollupListingResponse{
		RollupsData: rollups,
		Total:       uint64(len(rollups)),
	}, nil
}

func GetExtRollup(db HostDB, hash gethcommon.Hash) (*common.ExtRollup, error) {
	whereQuery := " WHERE r.hash=" + db.GetSQLStatement().Placeholder
	return fetchExtRollup(db.GetSQLDB(), whereQuery, truncTo16(hash))
}

// GetRollupHeader returns the rollup with the given hash.
func GetRollupHeader(db HostDB, hash gethcommon.Hash) (*common.RollupHeader, error) {
	whereQuery := " WHERE r.hash=" + db.GetSQLStatement().Placeholder
	return fetchRollupHeader(db.GetSQLDB(), whereQuery, truncTo16(hash))
}

// GetRollupHeaderByBlock returns the rollup for the given block
func GetRollupHeaderByBlock(db HostDB, blockHash gethcommon.Hash) (*common.RollupHeader, error) {
	whereQuery := " WHERE r.compression_block=" + db.GetSQLStatement().Placeholder
	return fetchRollupHeader(db.GetSQLDB(), whereQuery, blockHash)
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
	whereQuery := " WHERE hash=" + db.GetSQLStatement().Placeholder
	return fetchPublicRollup(db.GetSQLDB(), whereQuery, truncTo16(rollupHash))
}

func GetRollupBySeqNo(db HostDB, seqNo uint64) (*common.PublicRollup, error) {
	whereQuery := " WHERE " + db.GetSQLStatement().Placeholder + " BETWEEN start_seq AND end_seq"
	return fetchPublicRollup(db.GetSQLDB(), whereQuery, seqNo)
}

func GetRollupBatches(db HostDB, rollupHash gethcommon.Hash) (*common.BatchListingResponse, error) {
	whereQuery := " WHERE r.hash=" + db.GetSQLStatement().Placeholder
	orderQuery := " ORDER BY b.height DESC"
	query := selectRollupBatches + whereQuery + orderQuery
	rows, err := db.GetSQLDB().Query(query, truncTo16(rollupHash))
	if err != nil {
		return nil, fmt.Errorf("query execution for select rollup batches failed: %w", err)
	}
	defer rows.Close()

	var batches []common.PublicBatch
	for rows.Next() {
		var (
			sequenceInt64 int
			hash          []byte
			fullHash      gethcommon.Hash
			heightInt64   int
			extBatch      []byte
		)
		err := rows.Scan(&sequenceInt64, &hash, &fullHash, &heightInt64, &extBatch)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errutil.ErrNotFound
			}
			return nil, fmt.Errorf("failed to fetch rollup batches: %w", err)
		}
		var b common.ExtBatch
		err = rlp.DecodeBytes(extBatch, &b)
		if err != nil {
			return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
		}

		batch := common.PublicBatch{
			SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
			Hash:             bytesToHexString(hash),
			FullHash:         fullHash,
			Height:           new(big.Int).SetInt64(int64(heightInt64)),
			TxCount:          new(big.Int).SetInt64(int64(len(b.TxHashes))),
			Header:           b.Header,
			EncryptedTxBlob:  b.EncryptedTxBlob,
		}
		batches = append(batches, batch)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &common.BatchListingResponse{
		BatchesData: batches,
		Total:       uint64(len(batches)),
	}, nil
}

func fetchRollupHeader(db *sql.DB, whereQuery string, args ...any) (*common.RollupHeader, error) {
	rollup, err := fetchExtRollup(db, whereQuery, args...)
	if err != nil {
		return nil, err
	}
	return rollup.Header, nil
}

func fetchExtRollup(db *sql.DB, whereQuery string, args ...any) (*common.ExtRollup, error) {
	var rollupBlob []byte
	query := selectExtRollup + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&rollupBlob)
	} else {
		err = db.QueryRow(query).Scan(&rollupBlob)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	var rollup common.ExtRollup
	err = rlp.DecodeBytes(rollupBlob, &rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup: %w", err)
	}

	return &rollup, nil
}

func fetchHeadRollup(db *sql.DB) (*common.ExtRollup, error) {
	var extRollup []byte
	err := db.QueryRow(selectLatestRollup).Scan(&extRollup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	var rollup common.ExtRollup
	err = rlp.DecodeBytes(extRollup, &rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup: %w", err)
	}

	return &rollup, nil
}

func fetchPublicRollup(db *sql.DB, whereQuery string, args ...any) (*common.PublicRollup, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
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
