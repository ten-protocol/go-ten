package storage

import (
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
)

func TestAddBatchWithDuplicateSequenceNumberStorageImpl(t *testing.T) {
	db, _ := hostdb.CreateSQLiteDB(t)
	logger := testlog.Logger()
	s := NewStorage(db, logger)
	batch1 := hostdb.CreateBatch(2, []common.L2TxHash{})
	err := s.AddBatch(&batch1)
	if err != nil {
		t.Errorf("could not store first batch. Cause: %s", err)
	}

	batch2 := hostdb.CreateBatch(2, []common.L2TxHash{gethcommon.BytesToHash([]byte("different"))})
	err = s.AddBatch(&batch2)
	// verify constraint error is returned
	if err == nil {
		t.Errorf("expected error but got: success")
	}
}
