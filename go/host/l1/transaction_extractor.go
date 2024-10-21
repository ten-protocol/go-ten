package l1

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

// TransactionExtractor defines the interface for extracting TEN-relevant transactions from L1 blocks
type TransactionExtractor interface {
	// ExtractRelevantTenTransactions extracts any rollup or important contract txs from an L1 block
	ExtractRelevantTenTransactions(ctx context.Context, block *types.Block, receipts types.Receipts) ([]*common.TxAndReceiptAndBlobs, []*ethadapter.L1RollupTx, []*ethadapter.L1SetImportantContractsTx)
	// FindSecretResponseTx will return the secret response tx from an L1 block
	FindSecretResponseTx(block *types.Block) []*ethadapter.L1RespondSecretTx
}

// TenTransactionExtractor extracts TEN-relevant transactions from L1 blocks
type TenTransactionExtractor struct {
	mgmtContractLib mgmtcontractlib.MgmtContractLib
	blobResolver    BlobResolver
	logger          gethlog.Logger
}

// NewTransactionExtractor creates a new TenTransactionExtractor
func NewTransactionExtractor(mgmtContractLib mgmtcontractlib.MgmtContractLib, blobResolver BlobResolver, logger gethlog.Logger) TenTransactionExtractor {
	return TenTransactionExtractor{
		mgmtContractLib: mgmtContractLib,
		blobResolver:    blobResolver,
		logger:          logger,
	}
}

func (e *TenTransactionExtractor) Start() error {
	// todo (#2495) we should monitor for relevant L1 events instead of scanning every transaction in the block
	return nil
}

func (e *TenTransactionExtractor) Stop() error {
	return nil
}

func (e *TenTransactionExtractor) HealthStatus(_ context.Context) error {
	return nil
}

// ExtractRelevantTenTransactions extracts any transactions from the block that are relevant to TEN
// todo (#2495) we should monitor for relevant L1 events instead of scanning every transaction in the block
func (e *TenTransactionExtractor) ExtractRelevantTenTransactions(ctx context.Context, block *types.Block, receipts types.Receipts) ([]*common.TxAndReceiptAndBlobs, []*ethadapter.L1RollupTx, []*ethadapter.L1SetImportantContractsTx) {
	txWithReceiptsAndBlobs := make([]*common.TxAndReceiptAndBlobs, 0)
	rollupTxs := make([]*ethadapter.L1RollupTx, 0)
	contractAddressTxs := make([]*ethadapter.L1SetImportantContractsTx, 0)

	txs := block.Transactions()
	for i, rec := range receipts {
		if rec.BlockNumber == nil {
			continue // Skip non-relevant transactions
		}

		decodedTx := e.mgmtContractLib.DecodeTx(txs[i])
		var blobs []*kzg4844.Blob
		var err error

		switch typedTx := decodedTx.(type) {
		case *ethadapter.L1SetImportantContractsTx:
			contractAddressTxs = append(contractAddressTxs, typedTx)
		case *ethadapter.L1RollupHashes:
			blobs, err = e.blobResolver.FetchBlobs(ctx, block.Header(), typedTx.BlobHashes)
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					e.logger.Crit("Blobs were not found on beacon chain or archive service", "block", block.Hash(), "error", err)
				} else {
					e.logger.Crit("could not fetch blobs", log.ErrKey, err)
				}
				continue
			}

			encodedRlp, err := ethadapter.DecodeBlobs(blobs)
			if err != nil {
				e.logger.Crit("could not decode blobs.", log.ErrKey, err)
				continue
			}

			rlp := &ethadapter.L1RollupTx{
				Rollup: encodedRlp,
			}
			rollupTxs = append(rollupTxs, rlp)
		}

		// compile the tx, receipt and blobs into a single struct for submission to the enclave
		txWithReceiptsAndBlobs = append(txWithReceiptsAndBlobs, &common.TxAndReceiptAndBlobs{
			Tx:      txs[i],
			Receipt: rec,
			Blobs:   blobs,
		})
	}

	return txWithReceiptsAndBlobs, rollupTxs, contractAddressTxs
}

// FindSecretResponseTx will scan the block for any secret response transactions. This is separate from the above method
// as we do not require the receipts for these transactions.
func (e *TenTransactionExtractor) FindSecretResponseTx(block *types.Block) []*ethadapter.L1RespondSecretTx {
	secretRespTxs := make([]*ethadapter.L1RespondSecretTx, 0)

	for _, tx := range block.Transactions() {
		t := e.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			secretRespTxs = append(secretRespTxs, scrtTx)
			continue
		}
	}
	return secretRespTxs
}
