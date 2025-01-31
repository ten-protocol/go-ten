package hostdb

import "strconv"

// SQLStatements struct holds SQL statements for a specific database type
type SQLStatements struct {
	InsertBatch             string
	InsertTransactions      string
	UpdateTxCount           string
	InsertRollup            string
	InsertCrossChainMessage string
	InsertBlock             string
	Pagination              string
	Placeholder             string
}

func (s SQLStatements) GetPlaceHolder(pos int) string {
	if s.Placeholder == "?" {
		return s.Placeholder
	}
	return "$" + strconv.Itoa(pos)
}

func SQLiteSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:             "INSERT INTO batch_host (sequence, hash, height, ext_batch) VALUES (?, ?, ?, ?)",
		InsertTransactions:      "INSERT INTO transaction_host (hash, b_sequence) VALUES ",
		UpdateTxCount:           "UPDATE transaction_count SET total=? WHERE id=1",
		InsertRollup:            "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values (?,?,?,?,?,?)",
		InsertBlock:             "INSERT INTO block_host (hash, header) values (?,?)",
		InsertCrossChainMessage: "INSERT INTO cross_chain_message_host (message_hash, message_type, rollup_id) values (?,?,?)",
		Pagination:              "LIMIT ? OFFSET ?",
		Placeholder:             "?",
	}
}

func PostgresSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:             "INSERT INTO batch_host (sequence, hash, height, ext_batch) VALUES ($1, $2, $3, $4)",
		InsertTransactions:      "INSERT INTO transaction_host (hash, b_sequence) VALUES ",
		UpdateTxCount:           "UPDATE transaction_count SET total=$1 WHERE id=1",
		InsertRollup:            "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values ($1, $2, $3, $4, $5, $6)",
		InsertBlock:             "INSERT INTO block_host (hash, header) VALUES ($1, $2)",
		InsertCrossChainMessage: "INSERT INTO cross_chain_message_host (message_hash, message_type, rollup_id) values ($1, $2, $3)",
		Pagination:              "LIMIT $1 OFFSET $2",
		Placeholder:             "$1",
	}
}
