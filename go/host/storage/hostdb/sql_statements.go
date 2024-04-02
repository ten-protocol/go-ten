package hostdb

// SQLStatements struct holds SQL statements for a specific database type
type SQLStatements struct {
	InsertBatch        string
	InsertTransactions string
	SelectTxCount      string
	InsertTxCount      string
	Placeholder        string
}

func SQLiteSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:        "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES (?, ?, ?, ?, ?)",
		InsertTransactions: "REPLACE INTO transactions_host (hash, b_sequence) VALUES (?, ?)",
		SelectTxCount:      "SELECT total FROM transaction_count WHERE id = 1",
		InsertTxCount:      "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET total = EXCLUDED.total",
		Placeholder:        "?",
	}
}

func PostgresSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:        "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES ($1, $2, $3, $4, $5)",
		InsertTransactions: "INSERT INTO transactions_host (hash, b_sequence) VALUES ($1, $2) ON CONFLICT (hash) DO NOTHING",
		SelectTxCount:      "SELECT total FROM transaction_count WHERE id = 1",
		InsertTxCount:      "INSERT INTO transaction_count (id, total) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET total = EXCLUDED.total",
		Placeholder:        "$1",
	}
}
