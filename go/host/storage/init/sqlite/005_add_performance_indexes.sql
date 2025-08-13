-- Migration 005: Add performance optimization indexes
-- Add composite indexes for better query performance

-- Index for rollup queries by timestamp and sequence range
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_TIMESTAMP_SEQ ON rollup_host (time_stamp, start_seq, end_seq);

-- Index for batch queries by height and sequence
CREATE INDEX IF NOT EXISTS IDX_BATCH_HEIGHT_SEQ ON batch_host (height, sequence);

-- Index for transaction queries by batch sequence
CREATE INDEX IF NOT EXISTS IDX_TX_BATCH_SEQ ON transaction_host (b_sequence, hash);

-- Add a function to get database statistics
CREATE OR REPLACE FUNCTION get_db_stats()
RETURNS TABLE (
    table_name VARCHAR(100),
    row_count BIGINT,
    table_size TEXT,
    index_size TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        schemaname||'.'||tablename as table_name,
        n_tup_ins as row_count,
        pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as table_size,
        pg_size_pretty(pg_indexes_size(schemaname||'.'||tablename)) as index_size
    FROM pg_stat_user_tables 
    WHERE schemaname = 'public'
    ORDER BY tablename;
END;
$$ LANGUAGE plpgsql;
