drop index receipt on tendb.receipt_viewer;
CREATE INDEX receipt ON tendb.receipt_viewer (eoa, receipt);