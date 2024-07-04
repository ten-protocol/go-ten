package hostdb

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"
)

// An arbitrary number to put in the header
const batchNumber = 777

func createSQLiteDB(t *testing.T) (HostDB, error) {
	hostDB, err := sqlite.CreateTemporarySQLiteHostDB("", "mode=memory")
	if err != nil {
		t.Fatalf("unable to create temp sql db: %s", err)
	}
	return NewHostDB(hostDB, SQLiteSQLStatements())
}

func bytesToHexString(bytes []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes))
}
