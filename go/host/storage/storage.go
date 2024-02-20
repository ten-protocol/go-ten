package storage

import "database/sql"

type HostDatabase struct {
	db *sql.DB
}
