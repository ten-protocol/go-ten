package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	cfgExists = "select exists(select 1 from config where ky = ?)"
	cfgInsert = "insert into config values (?,?)"
	cfgUpdate = "update config set val = ? where ky = ?"
	cfgSelect = "select val from config where ky=?"
)

func InsertOrUpdateConfig(ctx context.Context, dbtx *sqlx.Tx, key string, value any) error {
	var exists bool
	// check if it exists then insert or update - this keeps it agnostic to the type of sql database

	// convert placeholders based on the database driver
	reboundExists := dbtx.Rebind(cfgExists)
	err := dbtx.GetContext(ctx, &exists, reboundExists, key)
	if err != nil {
		return fmt.Errorf("failed to check existence of config key %q: %w", key, err)
	}

	if exists {
		reboundUpdate := dbtx.Rebind(cfgUpdate)
		_, err = dbtx.ExecContext(ctx, reboundUpdate, value, key)
		if err != nil {
			return fmt.Errorf("failed to update config key %q: %w", key, err)
		}
	} else {
		reboundInsert := dbtx.Rebind(cfgInsert)
		_, err = dbtx.ExecContext(ctx, reboundInsert, key, value)
		if err != nil {
			return fmt.Errorf("failed to insert config key %q: %w", key, err)
		}
	}

	return nil
}

func WriteConfig(ctx context.Context, db *sqlx.Tx, key string, value []byte) (sql.Result, error) {
	reboundInsert := db.Rebind(cfgInsert)
	return db.ExecContext(ctx, reboundInsert, key, value)
}

func FetchConfig(ctx context.Context, db *sqlx.DB, key string) ([]byte, error) {
	var res []byte

	reboundQuery := db.Rebind(cfgSelect)
	err := db.QueryRowContext(ctx, reboundQuery, key).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}
