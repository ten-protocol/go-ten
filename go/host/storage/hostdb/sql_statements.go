package hostdb

// SQLStatements struct holds SQL statements for a specific database type
type SQLStatements struct {
	InsertBatch         string
	InsertTransactions  string
	InsertTxCount       string
	InsertRollup        string
	InsertBlock         string
	SelectRollups       string
	SelectBlocks        string
	SelectRollupBatches string
	Placeholder         string
}

func SQLiteSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:         "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES (?, ?, ?, ?, ?)",
		InsertTransactions:  "REPLACE INTO transactions_host (hash, b_sequence) VALUES (?, ?)",
		InsertTxCount:       "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET total = EXCLUDED.total",
		InsertRollup:        "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values (?,?,?,?,?,?)",
		InsertBlock:         "REPLACE INTO block_host (hash, header, rollup_hash) values (?,?,?)",
		SelectRollups:       "SELECT id, hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block FROM rollup_host ORDER BY id DESC LIMIT ? OFFSET ?",
		SelectBlocks:        "SELECT id, hash, header, rollup_hash FROM block_host ORDER BY id DESC LIMIT ? OFFSET ?",
		SelectRollupBatches: "SELECT b.sequence, b.hash, b.full_hash, b.height, b.ext_batch FROM rollup_host r JOIN batch_host b ON r.start_seq <= b.sequence AND r.end_seq >= b.sequence WHERE r.hash = ?",
		Placeholder:         "?",
	}
}

func PostgresSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:         "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES ($1, $2, $3, $4, $5)",
		InsertTransactions:  "INSERT INTO transactions_host (hash, b_sequence) VALUES ($1, $2) ON CONFLICT (hash) DO NOTHING",
		InsertTxCount:       "INSERT INTO transaction_count (id, total) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET total = EXCLUDED.total",
		InsertRollup:        "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values ($1, $2, $3, $4, $5, $6)",
		InsertBlock:         "INSERT INTO block_host (hash, header, rollup_hash) VALUES ($1, $2, $3) ON CONFLICT (hash) DO NOTHING",
		SelectRollups:       "SELECT id, hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block FROM rollup_host ORDER BY id DESC LIMIT $1 OFFSET $2",
		SelectBlocks:        "SELECT id, hash, header, rollup_hash FROM block_host ORDER BY id DESC LIMIT $1 OFFSET $2",
		SelectRollupBatches: "SELECT b.sequence, b.hash, b.full_hash, b.height, (SELECT COUNT(*) FROM transactions t WHERE t.batch_hash = b.hash) AS tx_count, b.ext_batch FROM rollup_host JOIN batch_host b ON r.start_seq <= b.sequence AND r.end_seq >= b.sequence WHERE r.hash = $1",
		Placeholder:         "$1",
	}
}
