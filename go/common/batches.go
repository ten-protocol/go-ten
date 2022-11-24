package common

// ExtBatch is an encrypted form of batch used when passing the batch around outside of an enclave.
// TODO - #718 - Expand this structure to contain the required fields.
type ExtBatch struct {
	Header          *Header
	TxHashes        []TxHash // The hashes of the transactions included in the batch.
	EncryptedTxBlob EncryptedTransactions
}

// TODO - #718 - Cache hash calculation.
