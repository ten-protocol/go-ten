package backend

import (
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/obsclient"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Backend struct {
	obsClient *obsclient.ObsClient
}

func NewBackend(obsClient *obsclient.ObsClient) *Backend {
	return &Backend{
		obsClient: obsClient,
	}
}

func (b *Backend) GetLatestBatch() (*common.BatchHeader, error) {
	return b.obsClient.GetLatestBatch()
}

func (b *Backend) GetTenNodeHealthStatus() (bool, error) {
	health, err := b.obsClient.Health()
	if err != nil {
		return false, err
	}
	return health.OverallHealth, nil
}

func (b *Backend) GetLatestRollup() (*common.RollupHeader, error) {
	return &common.RollupHeader{}, nil
}

func (b *Backend) GetNodeCount() (int, error) {
	// return b.obsClient.ActiveNodeCount()
	return 0, nil
}

func (b *Backend) GetTotalContractCount() (int, error) {
	return b.obsClient.GetTotalContractCount()
}

func (b *Backend) GetTotalTransactionCount() (int, error) {
	return b.obsClient.GetTotalTransactionCount()
}

func (b *Backend) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return b.obsClient.GetLatestRollupHeader()
}

func (b *Backend) GetBatchByHash(hash gethcommon.Hash) (*common.ExtBatch, error) {
	return b.obsClient.GetBatchByHash(hash)
}

func (b *Backend) GetBatchByHeight(height *big.Int) (*common.PublicBatch, error) {
	return b.obsClient.GetBatchByHeight(height)
}

func (b *Backend) GetBatchBySeq(seq *big.Int) (*common.PublicBatch, error) {
	return b.obsClient.GetBatchBySeq(seq)
}

func (b *Backend) GetRollupBySeqNo(seqNo uint64) (*common.PublicRollup, error) {
	return b.obsClient.GetRollupBySeqNo(seqNo)
}

func (b *Backend) GetBatchHeader(hash gethcommon.Hash) (*common.BatchHeader, error) {
	return b.obsClient.GetBatchHeaderByHash(hash)
}

func (b *Backend) GetTransaction(hash gethcommon.Hash) (*common.PublicTransaction, error) {
	return b.obsClient.GetTransaction(hash)
}

func (b *Backend) GetPublicTransactions(offset uint64, size uint64) (*common.TransactionListingResponse, error) {
	return b.obsClient.GetPublicTxListing(&common.QueryPagination{
		Offset: offset,
		Size:   uint(size),
	})
}

func (b *Backend) GetBatchesListing(offset uint64, size uint64) (*common.BatchListingResponse, error) {
	return b.obsClient.GetBatchesListing(&common.QueryPagination{
		Offset: offset,
		Size:   uint(size),
	})
}

func (b *Backend) GetBlockListing(offset uint64, size uint64) (*common.BlockListingResponse, error) {
	return b.obsClient.GetBlockListing(&common.QueryPagination{
		Offset: offset,
		Size:   uint(size),
	})
}

func (b *Backend) GetRollupListing(offset uint64, size uint64) (*common.RollupListingResponse, error) {
	return b.obsClient.GetRollupListing(&common.QueryPagination{
		Offset: offset,
		Size:   uint(size),
	})
}

func (b *Backend) GetRollupByHash(hash gethcommon.Hash) (*common.PublicRollup, error) {
	return b.obsClient.GetRollupByHash(hash)
}

func (b *Backend) GetRollupBatches(hash gethcommon.Hash, offset uint64, size uint64) (*common.BatchListingResponse, error) {
	return b.obsClient.GetRollupBatches(hash, &common.QueryPagination{
		Offset: offset,
		Size:   uint(size),
	})
}

func (b *Backend) GetBatchTransactions(hash gethcommon.Hash, offset uint64, size uint64) (*common.TransactionListingResponse, error) {
	return b.obsClient.GetBatchTransactions(hash, &common.QueryPagination{
		Offset: offset,
		Size:   uint(size),
	})
}

/*func (b *Backend) DecryptTxBlob(payload string) ([]*common.L2Tx, error) {
	encryptedTxBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("could not decode encrypted transaction blob from Base64. Cause: %w", err)
	}

	key := gethcommon.Hex2Bytes(crypto.RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not initialise AES cipher for enclave rollup key. Cause: %w", err)
	}
	transactionCipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not initialise wrapper for AES cipher for enclave rollup key. Cause: %w", err)
	}

	// The nonce is prepended to the ciphertext.
	nonce := encryptedTxBytes[0:crypto.GCMNonceLength]
	ciphertext := encryptedTxBytes[crypto.GCMNonceLength:]
	compressedTxs, err := transactionCipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt encrypted L2 transactions. Cause: %w", err)
	}

	compressionService := compression.NewBrotliDataCompressionService()

	encodedTxs, err := compressionService.Decompress(compressedTxs)
	if err != nil {
		return nil, fmt.Errorf("could not decompress L2 transactions. Cause: %w", err)
	}

	var cleartextTxs []*common.L2Tx
	if err = rlp.DecodeBytes(encodedTxs, &cleartextTxs); err != nil {
		return nil, fmt.Errorf("could not decode encoded L2 transactions. Cause: %w", err)
	}

	return cleartextTxs, nil
}*/

func (b *Backend) GetConfig() (*common.TenNetworkInfo, error) {
	return b.obsClient.GetConfig()
}

func (b *Backend) Search(query string) (*common.SearchResponse, error) {
	return b.obsClient.Search(query)
}
