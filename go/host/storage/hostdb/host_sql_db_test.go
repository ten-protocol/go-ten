package hostdb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ten-protocol/go-ten/go/common/storage/database/init/sqlite"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"testing"
)

func createSQLiteDB(t *testing.T) (*sql.DB, error) {
	db, err := sqlite.CreateTemporarySQLiteHostDB("", "mode=memory", testlog.Logger(), "host_init.sql")
	if err != nil {
		t.Fatalf("unable to create temp sql db: %s", err)
	}
	return db, nil
}
