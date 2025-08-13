-- Migration 004: Add audit tables for tracking database changes
CREATE TABLE IF NOT EXISTS audit_log
(
    id              SERIAL PRIMARY KEY,
    table_name      VARCHAR(100) NOT NULL,
    operation       VARCHAR(20) NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    record_id       INT,
    old_values      JSONB,
    new_values      JSONB,
    changed_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    changed_by      VARCHAR(100) DEFAULT 'system'
);

CREATE INDEX IF NOT EXISTS IDX_AUDIT_LOG_TABLE ON audit_log (table_name);
CREATE INDEX IF NOT EXISTS IDX_AUDIT_LOG_OPERATION ON audit_log (operation);
CREATE INDEX IF NOT EXISTS IDX_AUDIT_LOG_TIMESTAMP ON audit_log (changed_at);

-- Add a comment to the rollup_host table
COMMENT ON TABLE rollup_host IS 'Stores rollup information with batch sequence ranges';
COMMENT ON COLUMN rollup_host.start_seq IS 'First batch sequence number in this rollup';
COMMENT ON COLUMN rollup_host.end_seq IS 'Last batch sequence number in this rollup';
