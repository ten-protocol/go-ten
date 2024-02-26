package hostdb

import (
	"database/sql"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/storage/database/init/sqlite"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"testing"
)

const truncHash = 16

// An arbitrary number to put in the header
const batchNumber = 777

func truncTo16(hash gethcommon.Hash) []byte {
	return truncLastTo16(hash.Bytes())
}

func truncLastTo16(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	start := len(bytes) - truncHash
	if start < 0 {
		start = 0
	}
	b := bytes[start:]
	c := make([]byte, truncHash)
	copy(c, b)
	return c
}

func createSQLiteDB(t *testing.T) (*sql.DB, error) {
	db, err := sqlite.CreateTemporarySQLiteHostDB("", "mode=memory", testlog.Logger(), "host_init.sql")
	if err != nil {
		t.Fatalf("unable to create temp sql db: %s", err)
	}
	return db, nil
}
