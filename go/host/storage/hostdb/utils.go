package hostdb

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"
)

// An arbitrary number to put in the header
const batchNumber = 777

func CreateSQLiteDB(t *testing.T) (HostDB, error) {
	hostDB, err := sqlite.CreateTemporarySQLiteHostDB("", "mode=memory")
	if err != nil {
		t.Fatalf("unable to create temp sql db: %s", err)
	}

	// Create a test logger for the database
	testLogger := log.New(log.HostCmp, 1, log.SysOut)
	return NewHostDB(hostDB, SQLiteSQLStatements(), testLogger)
}

func CreateBatch(batchNum int64, txHashes []common.L2BatchHash) common.ExtBatch {
	header := common.BatchHeader{
		SequencerOrderNo: big.NewInt(batchNum),
		Number:           big.NewInt(batchNum),
		Time:             uint64(time.Now().Unix()),
	}
	batch := common.ExtBatch{
		Header:   &header,
		TxHashes: txHashes,
	}

	return batch
}

func CreateBatchWithDiffHeight(seqNo int64, height int64, txHashes []common.L2BatchHash) common.ExtBatch {
	header := common.BatchHeader{
		SequencerOrderNo: big.NewInt(seqNo),
		Number:           big.NewInt(height),
		Time:             uint64(time.Now().Unix()),
	}
	batch := common.ExtBatch{
		Header:   &header,
		TxHashes: txHashes,
	}

	return batch
}

func bytesToHexString(bytes []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes))
}
