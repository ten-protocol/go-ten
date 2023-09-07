package common

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
)

type (
	StateRoot = common.Hash
	TxHash    = common.Hash

	// EncryptedSubscriptionLogs - Alias for the event subscription updates going
	// out of the enclave.
	EncryptedSubscriptionLogs = map[rpc.ID][]byte

	// StreamL2UpdatesResponse - the struct encoded for each response message
	// when streaming batches out of the enclave.
	// The properties inside need to be encrypted according to the privacy rules.
	StreamL2UpdatesResponse struct {
		Batch *ExtBatch
		Logs  EncryptedSubscriptionLogs
	}

	// MainNet aliases
	L1Address     = common.Address
	L1BlockHash   = common.Hash
	L1Block       = types.Block
	L1Transaction = types.Transaction
	L1Receipt     = types.Receipt
	L1Receipts    = types.Receipts

	// Local Obscuro aliases
	L2BatchHash    = common.Hash
	L2RollupHash   = common.Hash
	L2TxHash       = common.Hash
	L2Tx           = types.Transaction
	L2Transactions = types.Transactions
	L2Address      = common.Address
	L2Receipt      = types.Receipt
	L2Receipts     = types.Receipts

	CrossChainMessage  = MessageBus.StructsCrossChainMessage
	CrossChainMessages = []CrossChainMessage
	ValueTransferEvent struct {
		Sender   common.Address
		Receiver common.Address
		Amount   *big.Int
	}
	ValueTransferEvents   = []ValueTransferEvent
	EncryptedTx           []byte // A single transaction, encoded as a JSON list of transaction binary hexes and encrypted using the enclave's public key
	EncryptedTransactions []byte // A blob of encrypted transactions, as they're stored in the rollup, with the nonce prepended.

	EncryptedParamsGetBalance      []byte // The params for an RPC getBalance request, as a JSON object encrypted with the public key of the enclave.
	EncryptedParamsCall            []byte // As above, but for an RPC call request.
	EncryptedParamsGetTxByHash     []byte // As above, but for an RPC getTransactionByHash request.
	EncryptedParamsGetTxReceipt    []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedParamsLogSubscription []byte // As above, but for an RPC logs subscription request.
	EncryptedParamsSendRawTx       []byte // As above, but for an RPC sendRawTransaction request.
	EncryptedParamsGetTxCount      []byte // As above, but for an RPC getTransactionCount request.
	EncryptedParamsEstimateGas     []byte // As above, but for an RPC estimateGas request.
	EncryptedParamsGetLogs         []byte // As above, but for an RPC getLogs request.
	EncryptedParamsGetStorageAt    []byte

	Nonce               = uint64
	EncodedRollup       []byte
	EncodedBatchMsg     []byte
	EncodedBatchRequest []byte
)

const (
	L2GenesisHeight = uint64(0)
	L1GenesisHeight = uint64(0)
	L2GenesisSeqNo  = uint64(1)
	// HeightCommittedBlocks is the number of blocks deep a transaction must be to be considered safe from reorganisations.
	HeightCommittedBlocks = 15
)

// SystemError is the interface for the internal errors generated in the Enclave and consumed by the Host
type SystemError interface {
	Error() string
}

// AttestationReport represents a signed attestation report from a TEE and some metadata about the source of it to verify it
type AttestationReport struct {
	Report      []byte         // the signed bytes of the report which includes some encrypted identifying data
	PubKey      []byte         // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	Owner       common.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data
	HostAddress string         // the IP address on which the host can be contacted by other Obscuro hosts for peer-to-peer communication
}

type (
	EncryptedSharedEnclaveSecret []byte
	EncodedAttestationReport     []byte
)

// BlockAndReceipts - a structure that contains a fuller view of a block. It allows iterating over the
// successful transactions, using the receipts. The receipts are bundled in the host node and thus verification
// is performed over their correctness.
// To work properly, all of the receipts are required, due to rlp encoding pruning some of the information.
// The receipts must also be in the correct order.
type BlockAndReceipts struct {
	Block                  *types.Block
	ReceiptsMap            map[int]*types.Receipt // sparse map with obscuro-relevant receipts in it
	Receipts               *types.Receipts
	successfulTransactions *types.Transactions
}

// ParseBlockAndReceipts - will create a container struct that has preprocessed the receipts
// and verified if they indeed match the receipt root hash in the block.
func ParseBlockAndReceipts(block *L1Block, receipts *L1Receipts) (*BlockAndReceipts, error) {
	if len(block.Transactions()) != len(*receipts) {
		// the receipts list is currently a *sparse* list of relevant receipts, it needs to have the same length as the
		// transactions list even though some of the entries may be nil
		return nil, fmt.Errorf("transactions and receipts are not the same length")
	}

	br := BlockAndReceipts{
		Block:                  block,
		Receipts:               receipts,
		ReceiptsMap:            make(map[int]*types.Receipt, len(block.Transactions())),
		successfulTransactions: nil,
	}

	for idx, receipt := range *receipts {
		br.ReceiptsMap[idx] = receipt
	}

	return &br, nil
}

// SuccessfulTransactions - returns slice containing only the transactions that have receipts with successful status.
func (br *BlockAndReceipts) SuccessfulTransactions() *types.Transactions {
	if br.successfulTransactions != nil {
		return br.successfulTransactions
	}

	txs := br.Block.Transactions()
	st := make(types.Transactions, 0)

	for idx, tx := range txs {
		receipt, ok := br.ReceiptsMap[idx]
		if ok && receipt.Status == types.ReceiptStatusSuccessful {
			st = append(st, tx)
		}
	}

	br.successfulTransactions = &st
	return br.successfulTransactions
}

// ChainFork - represents the result of walking the chain when processing a fork
type ChainFork struct {
	NewCanonical *types.Block
	OldCanonical *types.Block

	CommonAncestor   *types.Block
	CanonicalPath    []L1BlockHash
	NonCanonicalPath []L1BlockHash
}

func (cf *ChainFork) IsFork() bool {
	return len(cf.NonCanonicalPath) > 0
}

func (cf *ChainFork) String() string {
	return fmt.Sprintf("ChainFork{NewCanonical: %s, OldCanonical: %s, CommonAncestor: %s, CanonicalPath: %s, NonCanonicalPath: %s}",
		cf.NewCanonical.Hash(), cf.OldCanonical.Hash(), cf.CommonAncestor.Hash(), cf.CanonicalPath, cf.NonCanonicalPath)
}
