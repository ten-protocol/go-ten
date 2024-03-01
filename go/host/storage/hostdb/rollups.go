package hostdb

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	selectRollupHeader = "SELECT header from rollup"
	selectRollups      = "SELECT id, hash, start_seq, end_seq, time_stamp, header, compression_block FROM rollup ORDER BY id DESC LIMIT ? OFFSET ?"
	insertRollup       = "INSERT INTO rollup (hash, start_seq, end_seq, time_stamp, header, compression_block) values (?,?,?,?,?,?)"
)

// AddRollupHeader adds a rollup to the DB
func AddRollupHeader(db *sql.DB, rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error {
	// Check if the Header is already stored
	_, err := GetRollupHeader(db, rollup.Header.Hash())
	if err == nil {
		// The rollup is already stored, so we return early.
		return errutil.ErrAlreadyExists
	}

	rollupHeader := rollup.Header
	header, err := rlp.EncodeToBytes(rollupHeader)
	if err != nil {
		return fmt.Errorf("could not encode batch header: %w", err)
	}
	_, err = db.Exec(insertRollup,
		truncTo16(rollup.Header.Hash()),      // short hash
		metadata.FirstBatchSequence.Uint64(), // first batch sequence
		rollupHeader.LastBatchSeqNo,          // last batch sequence
		metadata.StartTime,                   // timestamp
		header,                               // header blob
		block.Hash(),                         // l1 block hash
	)

	if err != nil {
		return fmt.Errorf("could not store rollup in db: %w", err)
	}

	return nil
}

// GetRollupListing returns latest rollups given a pagination.
// For example, offset 1, size 10 will return the latest 11-20 rollups.
func GetRollupListing(db *sql.DB, pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	rows, err := db.Query(selectRollups, pagination.Size, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rollups []common.PublicRollup

	for rows.Next() {
		var id, startSeq, endSeq, timeStamp int
		var hash, headerBlob, compressionBlock []byte

		var rollup common.PublicRollup
		err = rows.Scan(&id, &hash, &startSeq, &endSeq, &timeStamp, &headerBlob, &compressionBlock)
		if err != nil {
			return nil, err
		}

		header := new(common.RollupHeader)
		if err := rlp.DecodeBytes(headerBlob, header); err != nil {
			return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
		}

		rollup = common.PublicRollup{
			ID:        big.NewInt(int64(id)),
			Hash:      hash,
			FirstSeq:  big.NewInt(int64(startSeq)),
			LastSeq:   big.NewInt(int64(endSeq)),
			Timestamp: uint64(timeStamp),
			Header:    header,
			L1Hash:    compressionBlock,
		}
		rollups = append(rollups, rollup)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &common.RollupListingResponse{
		Rollups: rollups,
		Total:   uint64(len(rollups)),
	}, nil
}

// GetRollupHeader returns the rollup with the given hash.
func GetRollupHeader(db *sql.DB, hash gethcommon.Hash) (*common.RollupHeader, error) {

	return fetchRollupHeader(db, " where r.hash=?", truncTo16(hash))
}

// GetRollupHeaderByBlock returns the rollup for the given block
func GetRollupHeaderByBlock(db *sql.DB, blockHash gethcommon.Hash) (*common.RollupHeader, error) {
	return fetchRollupHeader(db, " where r.compression_block=?", truncTo16(blockHash))
}

func fetchRollupHeader(db *sql.DB, whereQuery string, args ...any) (*common.RollupHeader, error) {
	var headerBlob []byte

	query := selectRollupHeader + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&headerBlob)
	} else {
		err = db.QueryRow(query).Scan(&headerBlob)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch rollup header by hash: %w", err)
	}
	var header common.RollupHeader
	err = rlp.DecodeBytes(headerBlob, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup header: %w", err)
	}

	return &header, nil
}
