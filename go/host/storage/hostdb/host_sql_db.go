package hostdb

import (
	"database/sql"
	"fmt"
)

// BeginTx handles beginning a transaction and committing or rolling it back based on the error
func BeginTx(db *sql.DB, txFunc func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return txFunc(tx)
}
