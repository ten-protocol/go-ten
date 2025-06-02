package common

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rlp"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
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
	L1Transaction = types.Transaction
	L1Receipt     = types.Receipt
	L1Receipts    = types.Receipts

	// L1 addresses
	NetworkConfigAddress         = common.Address
	CrossChainAddress            = common.Address
	EnclaveRegistryAddress       = common.Address
	DARegistryAddress            = common.Address
	L1MessageBusAddress          = common.Address
	L1BridgeAddress              = common.Address
	L1CrossChainMessengerAddress = common.Address

	// L2 addresses
	L2BridgeAddress                 = common.Address
	L2CrossChainMessengerAddress    = common.Address
	L2MessageBusAddress             = common.Address
	TransactionPostProcessorAddress = common.Address
	SystemContractsUpgraderAddress  = common.Address

	// Local Obscuro aliases
	L2BatchHash              = common.Hash
	L2RollupHash             = common.Hash
	L2TxHash                 = common.Hash
	L2Tx                     = types.Transaction
	L2Transactions           = types.Transactions
	L2Address                = common.Address
	L2Receipt                = types.Receipt
	L2Receipts               = types.Receipts
	SerializedCrossChainTree = []byte

	L2PricedTransaction struct {
		Tx *L2Tx
		// todo - add sender
		PublishingCost *big.Int
		FromSelf       bool
		SystemDeployer bool // Free contract construction
	}
	L2PricedTransactions []*L2PricedTransaction

	CrossChainMessage  MessageBus.StructsCrossChainMessage
	CrossChainMessages []CrossChainMessage
	ValueTransferEvent struct {
		Sender   common.Address
		Receiver common.Address
		Amount   *big.Int
		Sequence uint64
	}
	ValueTransferEvents            = []ValueTransferEvent
	EncryptedRequest               []byte
	EncryptedTx                    []byte // A single transaction, encoded as a JSON list of transaction binary hexes and encrypted using the enclave's public key
	EncryptedTransactions          []byte // A blob of encrypted transactions, as they're stored in the rollup, with the nonce prepended.
	EncryptedParamsLogSubscription []byte
	Nonce                          = uint64
	EncodedRollup                  []byte
	EncodedBatchMsg                []byte
	EncodedBatchRequest            []byte
	EncodedBlobHashes              []byte

	EnclaveID = common.Address

	// RollupSignature represents a signature over a rollup's composite hash
	RollupSignature = []byte

	// CreateRollupResult contains all data returned from creating a rollup
	CreateRollupResult struct {
		Signature RollupSignature // The signature over the composite hash
		Blobs     []*kzg4844.Blob // The blobs containing the rollup data
	}
)

func (c CrossChainMessage) IsValueTransfer(bridgeAuthority common.Address) bool {
	return c.Sender == bridgeAuthority
}

func (c CrossChainMessages) FilterValueTransfers(bridgeAuthority common.Address) CrossChainMessages {
	ret := make(CrossChainMessages, 0)
	for _, msg := range c {
		if msg.IsValueTransfer(bridgeAuthority) {
			ret = append(ret, msg)
		}
	}
	return ret
}

// FailedDecryptErr - when the TEN enclave fails to decrypt an RPC request
var FailedDecryptErr = errors.New("failed to decrypt RPC payload. please use the correct enclave key")

// EncryptedRPCRequest - an encrypted request with extra plaintext metadata
type EncryptedRPCRequest struct {
	Req  EncryptedRequest
	IsTx bool // we can make this an enum if we need to provide more info to the TEN host
}

func (txs L2PricedTransactions) ToTransactions() types.Transactions {
	ret := make(types.Transactions, 0)
	for _, tx := range txs {
		ret = append(ret, tx.Tx)
	}
	return ret
}

const (
	L1GenesisHeight = uint64(0)

	L2GenesisHeight           = uint64(0)
	L2GenesisSeqNo            = uint64(1)
	L2SysContractGenesisSeqNo = uint64(2)

	SyntheticTxGasLimit = params.MaxGasLimit
)

var GethGenesisParentHash = common.Hash{}

// SystemError is the interface for the internal errors generated in the Enclave and consumed by the Host
type SystemError interface {
	Error() string
}

// AttestationReport represents a signed attestation report from a TEE and some metadata about the source of it to verify it
type AttestationReport struct {
	Report      []byte         // the signed bytes of the report which includes some encrypted identifying data
	PubKey      []byte         // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	EnclaveID   common.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data
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
	BlockHeader            *types.Header
	TxsWithReceipts        []*TxAndReceiptAndBlobs
	successfulTransactions *types.Transactions
}

// ParseBlockAndReceipts - will create a container struct that has preprocessed the receipts
// and verified if they indeed match the receipt root hash in the block.
func ParseBlockAndReceipts(block *types.Header, receiptsAndBlobs []*TxAndReceiptAndBlobs) (*BlockAndReceipts, error) {
	br := BlockAndReceipts{
		BlockHeader:     block,
		TxsWithReceipts: receiptsAndBlobs,
	}

	return &br, nil
}

func (br *BlockAndReceipts) Receipts() L1Receipts {
	rec := make(L1Receipts, 0)
	for _, txsWithReceipt := range br.TxsWithReceipts {
		rec = append(rec, txsWithReceipt.Receipt)
	}
	return rec
}

// RelevantTransactions - returns slice containing only the transactions that have receipts with successful status.
func (br *BlockAndReceipts) RelevantTransactions() *types.Transactions {
	if br.successfulTransactions != nil {
		return br.successfulTransactions
	}

	st := make(types.Transactions, 0)
	for _, tx := range br.TxsWithReceipts {
		if tx.Receipt.Status == types.ReceiptStatusSuccessful {
			st = append(st, tx.Tx)
		}
	}
	br.successfulTransactions = &st
	return br.successfulTransactions
}

// ChainFork - represents the result of walking the chain when processing a fork
type ChainFork struct {
	NewCanonical *types.Header
	OldCanonical *types.Header

	CommonAncestor   *types.Header
	CanonicalPath    []L1BlockHash
	NonCanonicalPath []L1BlockHash
}

func (cf *ChainFork) IsFork() bool {
	return len(cf.NonCanonicalPath) > 0
}

func (cf *ChainFork) String() string {
	if cf == nil {
		return ""
	}
	return fmt.Sprintf("ChainFork{NewCanonical: %s, OldCanonical: %s, CommonAncestor: %s, CanonicalPath: %s, NonCanonicalPath: %s}",
		cf.NewCanonical.Hash(), cf.OldCanonical.Hash(), cf.CommonAncestor.Hash(), cf.CanonicalPath, cf.NonCanonicalPath)
}

func MaskedSender(address L2Address) L2Address {
	return common.BigToAddress(big.NewInt(0).Sub(address.Big(), big.NewInt(1)))
}

type SystemContractAddresses map[string]*gethcommon.Address

func (s *SystemContractAddresses) ToString() string {
	var str string
	for name, addr := range *s {
		str += fmt.Sprintf("%s: %s; ", name, addr.Hex())
	}
	return str
}

// MaxNegativeTxTimeDeltaMs - to avoid negative numbers in the timestamp delta (block.time - tx.time), we adjust by 10s so that it's impossible to have negative values
// This represents the period until which transactions will come in *after* the sequencer started building a new batch.
// this constant is used by the mempool as well.
const MaxNegativeTxTimeDeltaMs = 10 * 1000

// TxWithTimestamp - RLP serializes a transaction together with the timestamp delta from the block time
type TxWithTimestamp struct {
	Tx          *L2Tx
	TimeDeltaMs *big.Int
}

// createTxWithTimestamp - delta=blockTimeMs - txTimeMs + MaxNegativeTxTimeDeltaMs
func createTxWithTimestamp(tx *L2Tx, blockTimeMs uint64) *TxWithTimestamp {
	return &TxWithTimestamp{
		Tx:          tx,
		TimeDeltaMs: big.NewInt((int64(blockTimeMs) + MaxNegativeTxTimeDeltaMs) - tx.Time().UnixMilli()),
	}
}

// l2Tx - txTimeMs=blockTimeMs - delta + MaxNegativeTxTimeDeltaMs
func (t *TxWithTimestamp) l2Tx(blockTimeMs uint64) *L2Tx {
	t.Tx.SetTime(time.UnixMilli((int64(blockTimeMs) - MaxNegativeTxTimeDeltaMs) - t.TimeDeltaMs.Int64()))
	return t.Tx
}

type TxsWithTimeStamp []*TxWithTimestamp

func (txs TxsWithTimeStamp) Txs(blockTime uint64) []*L2Tx {
	txsOnly := make([]*L2Tx, len(txs))
	blockTimeMs := blockTime * 1000
	for i, tx := range txs {
		txsOnly[i] = tx.l2Tx(blockTimeMs)
	}
	return txsOnly
}

func CreateTxsAndTimeStamp(tx []*L2Tx, blockTime uint64) *TxsWithTimeStamp {
	txs := make(TxsWithTimeStamp, len(tx))
	blockTimeMs := blockTime * 1000
	for i, t := range tx {
		txs[i] = createTxWithTimestamp(t, blockTimeMs)
	}
	return &txs
}

func (txs TxsWithTimeStamp) EncodeIndex(i int, w *bytes.Buffer) {
	rlp.Encode(w, txs[i])
}

func (txs TxsWithTimeStamp) Len() int {
	return len(txs)
}
