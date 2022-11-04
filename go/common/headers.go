package common

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// Used to hash headers.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// Header is a public / plaintext struct that holds common properties of rollups batches.
// Making changes to this struct will require GRPC + GRPC Converters regen
type Header struct {
	ParentHash  L2RootHash
	Agg         common.Address
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

// Withdrawal - this is the withdrawal instruction that is included in the rollup header.
type Withdrawal struct {
	// Type      uint8 // the type of withdrawal. For now only ERC20. Todo - add this once more ERCs are supported
	Amount    *big.Int
	Recipient common.Address // the user account that will receive the money
	Contract  common.Address // the contract
}

// HeaderWithTxHashes pairs a header with the hashes of the transactions in that rollup.
type HeaderWithTxHashes struct {
	Header   *Header
	TxHashes []TxHash
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding excluding the signature.
func (h *Header) Hash() L2RootHash {
	cp := *h
	cp.R = nil
	cp.S = nil
	hash, err := rlpHash(cp)
	if err != nil {
		panic("err hashing a rollup header")
	}
	return hash
}

// Encodes value, hashes the encoded bytes and returns the hash.
func rlpHash(value interface{}) (common.Hash, error) {
	var hash common.Hash

	sha := hasherPool.Get().(crypto.KeccakState)
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
