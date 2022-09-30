package common

import (
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/eth/filters"

	"github.com/ethereum/go-ethereum/common"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/core/types"
)

type (
	StateRoot             = common.Hash
	L1RootHash            = common.Hash
	L2RootHash            = common.Hash
	TxHash                = common.Hash
	L2Tx                  = types.Transaction
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

	EncryptedResponseGetBalance   []byte // The response for an RPC getBalance request, as a JSON object encrypted with the viewing key of the user.
	EncryptedResponseCall         []byte // As above, but for an RPC call request.
	EncryptedResponseGetTxReceipt []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedResponseSendRawTx    []byte // As above, but for an RPC sendRawTransaction request.
	EncryptedResponseGetTxByHash  []byte // As above, but for an RPC getTransactionByHash request.
	EncryptedResponseGetTxCount   []byte // As above, but for an RPC getTransactionCount request.
	EncryptedLogSubscription      []byte // As above, but for a log subscription request.
	EncryptedLogs                 []byte // As above, but for a log subscription response.
	EncryptedResponseEstimateGas  []byte // As above, but for an RPC estimateGas response.

	Nonce         = uint64
	EncodedRollup []byte
)

const (
	L2GenesisHeight = uint64(0)
	L1GenesisHeight = uint64(0)
	// HeightCommittedBlocks is the number of blocks deep a transaction must be to be considered safe from reorganisations.
	HeightCommittedBlocks = 15
)

// Header is a public / plaintext struct that holds common properties of the Rollup
// Making changes to this struct will require GRPC + GRPC Converters regen
type Header struct {
	ParentHash  L2RootHash
	Agg         common.Address
	RollupNonce Nonce            // RollupNonce holds the lottery rollup nonce
	Nonce       types.BlockNonce // Nonce ensure compatibility with ethereum
	L1Proof     L1RootHash       // the L1 block where the Parent was published
	Root        StateRoot
	TxHash      common.Hash // todo - include the synthetic deposits
	Number      *big.Int    // the rollup height
	Bloom       types.Bloom
	ReceiptHash common.Hash
	Extra       []byte
	R, S        *big.Int // signature values
	Withdrawals []Withdrawal

	// Specification fields - not used for now but are expected to be available
	UncleHash  common.Hash    `json:"sha3Uncles"`
	Coinbase   common.Address `json:"miner"      `
	Difficulty *big.Int       `json:"difficulty" `
	GasLimit   uint64         `json:"gasLimit"  `
	GasUsed    uint64         `json:"gasUsed"    `
	Time       uint64         `json:"timestamp"   `
	MixDigest  common.Hash    `json:"mixHash"`
	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *big.Int `json:"baseFeePerGas"`
}

// HeaderWithTxHashes represents a rollup header and the hashes of the transactions in the rollup.
type HeaderWithTxHashes struct {
	Header   *Header
	TxHashes []TxHash
}

// Withdrawal - this is the withdrawal instruction that is included in the rollup header
type Withdrawal struct {
	// Type      uint8 // the type of withdrawal. For now only ERC20. Todo - add this once more ERCs are supported
	Amount    *big.Int
	Recipient common.Address // the user account that will receive the money
	Contract  common.Address // the contract
}

// ExtRollup is used for communication between the enclave and the outside world.
type ExtRollup struct {
	Header          *Header
	TxHashes        []TxHash // The hashes of the transactions included in the rollup
	EncryptedTxBlob EncryptedTransactions
}

// EncryptedRollup extends ExtRollup with additional fields.
// This parallels the Block/extblock split in Go Ethereum.
type EncryptedRollup struct {
	Header       *Header
	TxHashes     []TxHash // The hashes of the transactions included in the rollup
	Transactions EncryptedTransactions
	hash         atomic.Value
}

// AttestationReport represents a signed attestation report from a TEE and some metadata about the source of it to verify it
type AttestationReport struct {
	Report      []byte         // the signed bytes of the report which includes some encrypted identifying data
	PubKey      []byte         // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	Owner       common.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data
	HostAddress string         // the IP address on which the host can be contacted by other Obscuro hosts for peer-to-peer communication
}

func (er ExtRollup) ToRollup() *EncryptedRollup {
	return &EncryptedRollup{
		Header:       er.Header,
		TxHashes:     er.TxHashes,
		Transactions: er.EncryptedTxBlob,
	}
}

type (
	EncryptedSharedEnclaveSecret []byte
	EncodedAttestationReport     []byte
)

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// RLPHash encodes value, hashes the encoded bytes and returns the hash.
func RLPHash(value interface{}) (common.Hash, error) {
	var hash common.Hash

	sha := hasherPool.Get().(gethcrypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()

	err := rlp.Encode(sha, value)
	if err != nil {
		return hash, fmt.Errorf("unable to encode Value. %w", err)
	}

	_, err = sha.Read(hash[:])
	if err != nil {
		return hash, fmt.Errorf("unable to read encoded value. %w", err)
	}

	return hash, nil
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *EncryptedRollup) Hash() L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)

	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding excluding the signature.
func (h *Header) Hash() L2RootHash {
	cp := *h
	cp.R = nil
	cp.S = nil
	hash, err := RLPHash(cp)
	if err != nil {
		panic("err hashing a rollup header")
	}
	return hash
}

// LogSubscription is an authenticated subscription to logs.
type LogSubscription struct {
	// The account the events relate to.
	Account *common.Address
	// A signature over the account address using a private viewing key. Prevents attackers from subscribing to
	// (encrypted) logs for other accounts to see the pattern of logs.
	// TODO - This does not protect against replay attacks, where someone resends an intercepted subscription request.
	Signature *[]byte
	// A subscriber-defined filter to apply to the stream of logs.
	Filter *filters.FilterCriteria
}

// IDAndEncLog pairs an encrypted log with the ID of the subscription that generated it.
type IDAndEncLog struct {
	ID     rpc.ID
	EncLog []byte
}

// IDAndLog pairs a encrypted log with the ID of the subscription that generated it.
type IDAndLog struct {
	ID  rpc.ID
	Log *types.Log
}
