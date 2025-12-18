package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	getQry = `select sdb.val from %s sdb where sdb.ky = ?`
	// `replace` will perform insert or replace if existing and this syntax works for both sqlite and edgeless db
	putQryBatchSqlite = `replace into %s (ky, val) values`
	putQryBatchEdb1   = `INSERT INTO %s (ky, val) VALUES `
	putQryValues      = `(?,?)`
	putQryBatchEdb2   = ` ON DUPLICATE KEY UPDATE val=VALUES(val)`
	delQry            = `delete from %s where ky = ?`
	// todo - how is the performance of this? probably extraordinarily slow
	searchQry   = `select ky, val from %s sdb where substring(sdb.ky, 1, ?) = ? and sdb.ky >= ? order by sdb.ky asc`
	dbChunkSize = 32 * 1024 // 32 KB chunks
)

var stateIDPrefix = []byte("L")

// routes the table based on the key length and prefix
// mirrors the prefixes used by go-ethereum
func getTable(key []byte) string {
	switch {
	case len(key) <= 32:
		return "statedb32"
	case len(key) == 33:
		switch key[0] {
		case rawdb.TrieNodeAccountPrefix[0]:
			return "statedb33_trie_node_account"
		case rawdb.TrieNodeStoragePrefix[0]:
			return "statedb33_trie_node_storage"
		case stateIDPrefix[0]:
			return "statedb33_state_id"
		case rawdb.SnapshotAccountPrefix[0]:
			return "statedb33_snapshot_account"
		case rawdb.SnapshotStoragePrefix[0]:
			return "statedb33_snapshot_storage"
		case rawdb.CodePrefix[0]:
			return "statedb33_code"
		default:
			return "statedb33"
		}
	case len(key) == 34:
		return "statedb34"
	case len(key) <= 65:
		return "statedb65"
	default:
		// it will fail here
		panic(fmt.Sprintf("key too long: %d", len(key)))
		return "non-existent-table"
	}
}

func has(ctx context.Context, db *sqlx.DB, key []byte) (bool, error) {
	var dummy []byte
	err := db.QueryRowContext(ctx, fmt.Sprintf(getQry, getTable(key)), key).Scan(&dummy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getJournal(ctx context.Context, db *sqlx.DB) ([]byte, error) {
	q := "select val from triedb_journal order by id asc"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []byte
	rowCount := 0

	for rows.Next() {
		var val []byte
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		result = append(result, val...)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rowCount == 0 {
		// No rows found
		return nil, errutil.ErrNotFound
	}

	return result, nil
}

// the journal can be quite large, so we split it into chunks and insert them one by one
// because edglessdb fails silently when the data is too large
func putJournal(ctx context.Context, db *sqlx.DB, value []byte) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction - %w", err)
	}
	defer tx.Rollback()

	// Truncate the journal table
	_, err = tx.ExecContext(ctx, "DELETE FROM triedb_journal")
	if err != nil {
		return fmt.Errorf("failed to truncate journal table - %w", err)
	}

	// Split value into chunks and insert
	totalLen := len(value)
	numChunks := (totalLen + dbChunkSize - 1) / dbChunkSize // ceiling division

	for i := 0; i < numChunks; i++ {
		start := i * dbChunkSize
		end := start + dbChunkSize
		if end > totalLen {
			end = totalLen
		}
		chunk := value[start:end]

		// Insert chunk with auto-incrementing id (id column should be AUTO_INCREMENT)
		_, err = tx.ExecContext(ctx, "INSERT INTO triedb_journal (val) VALUES (?)", chunk)
		if err != nil {
			return fmt.Errorf("failed to insert journal chunk %d - %w", i, err)
		}
	}

	return tx.Commit()
}

func get(ctx context.Context, db *sqlx.DB, key []byte) ([]byte, error) {
	var res []byte
	q := fmt.Sprintf(getQry, getTable(key))
	err := db.QueryRowxContext(ctx, q, key).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

// sanity check that we don't try to insert large values in the db and get unexpected errors later
func valTooLarge(val []byte) bool {
	return len(val) > dbChunkSize
}

func put(ctx context.Context, db *sqlx.DB, key []byte, value []byte) error {
	if valTooLarge(value) {
		panic(fmt.Sprintf("value too large: %d", len(value)))
		// return fmt.Errorf("value too large")
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = putKeyValues(ctx, tx, [][]byte{key}, [][]byte{value})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func putKeyValues(ctx context.Context, tx *sqlx.Tx, keys [][]byte, vals [][]byte) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("invalid command. should not happen")
	}

	for _, val := range vals {
		if valTooLarge(val) {
			panic(fmt.Sprintf("value too large: %d", len(val)))
			// return fmt.Errorf("value too large")
		}
	}

	// Group keys and values by table name, using getTable for routing.
	groupedKeys := make(map[string][][]byte)
	groupedVals := make(map[string][][]byte)

	for i, key := range keys {
		tableName := getTable(key)
		groupedKeys[tableName] = append(groupedKeys[tableName], key)
		groupedVals[tableName] = append(groupedVals[tableName], vals[i])
	}

	// Insert into each table we have accumulated keys for.
	for table, tKeys := range groupedKeys {
		tVals := groupedVals[table]
		if err := insertIntoTable(ctx, tx, table, tKeys, tVals); err != nil {
			return err
		}
	}

	return nil
}

func insertIntoTable(ctx context.Context, tx *sqlx.Tx, table string, keys [][]byte, vals [][]byte) error {
	if len(keys) == 0 {
		return nil
	}
	var update string
	if isMysql(tx.DriverName()) {
		update = fmt.Sprintf(putQryBatchEdb1, table) + repeat(putQryValues, ",", len(keys)) + putQryBatchEdb2
	} else {
		update = fmt.Sprintf(putQryBatchSqlite, table) + repeat(putQryValues, ",", len(keys))
	}
	values := make([]any, 0)
	for i := range keys {
		values = append(values, keys[i], vals[i])
	}
	_, err := tx.ExecContext(ctx, update, values...)
	if err != nil {
		// for some unknown reason, the mysql-panic driver doesn't intercept this error
		// until we figure out the reason, we'll panic here to bounce the server
		if errors.Is(err, mysql.ErrInvalidConn) {
			panic("Invalid connection")
		}
		return fmt.Errorf("failed to exec k/v transaction statement table=%s. kv=%v, err=%w", table, values, err)
	}
	return nil
}

func deleteKey(ctx context.Context, db *sqlx.DB, key []byte) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(delQry, getTable(key)), key)
	return err
}

func deleteKeys(ctx context.Context, db *sqlx.Tx, keys [][]byte) error {
	for _, del := range keys {
		_, err := db.ExecContext(ctx, fmt.Sprintf(delQry, getTable(del)), del)
		if err != nil {
			return err
		}
	}
	return nil
}

func newIterator(ctx context.Context, db *sqlx.DB, prefix []byte, start []byte) ethdb.Iterator {
	// Avoid mutating `prefix` backing array.
	pr := prefix
	st := make([]byte, 0, len(prefix)+len(start))
	st = append(st, prefix...)
	st = append(st, start...)

	rows, err := db.QueryContext(ctx, fmt.Sprintf(searchQry, getTable(st)), len(pr), pr, st)
	if err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}

	if err = rows.Err(); err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}
	return &iterator{
		rows: rows,
	}
}

func isMysql(driverName string) bool {
	return strings.Index(driverName, "mysql") == 0
}
