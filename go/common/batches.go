package common

// ExtBatch is an encrypted form of batch used when passing the batch around outside of an enclave.
type ExtBatch struct {
	EncryptedTxBlob EncryptedTransactions
}
