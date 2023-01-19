package common

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
	"golang.org/x/crypto/sha3"
)

// Used to hash headers.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// BatchHeader is a public / plaintext struct that holds common properties of batches.
// Making changes to this struct will require GRPC + GRPC Converters regen
type BatchHeader struct {
	// The fields present in Geth's `types/Header` struct.
	ParentHash  L2RootHash
	UncleHash   common.Hash    `json:"sha3Uncles"`
	Coinbase    common.Address `json:"miner"`
	Root        StateRoot      `json:"stateRoot"`
	TxHash      common.Hash    `json:"transactionsRoot"` // todo - include the synthetic deposits
	ReceiptHash common.Hash    `json:"receiptsRoot"`
	Bloom       types.Bloom    `json:"logsBloom"`
	Difficulty  *big.Int
	Number      *big.Int
	GasLimit    uint64
	GasUsed     uint64
	Time        uint64      `json:"timestamp"`
	Extra       []byte      `json:"extraData"`
	MixDigest   common.Hash `json:"mixHash"`
	Nonce       types.BlockNonce
	BaseFee     *big.Int

	// The custom Obscuro fields.
	Agg     common.Address // TODO - Can this be removed and replaced with the `Coinbase` field?
	L1Proof L1RootHash     // the L1 block used by the enclave to generate the current batch
	R, S    *big.Int       // signature values
	// TODO: mark as deprecated Withdrawals are now contained within cross chain messages.
	Withdrawals        []Withdrawal
	CrossChainMessages []MessageBus.StructsCrossChainMessage `json:"crossChainMessages"`

	// The block hash of the latest block that has been scanned for cross chain messages.
	LatestInboudCrossChainHash common.Hash `json:"inboundCrossChainHash"`

	// The block height of the latest block that has been scanned for cross chain messages.
	LatestInboundCrossChainHeight *big.Int `json:"inboundCrossChainHeight"`
}

// RollupHeader is a public / plaintext struct that holds common properties of rollups.
// Making changes to this struct will require GRPC + GRPC Converters regen
type RollupHeader struct {
	// The fields present in Geth's `types/Header` struct.
	ParentHash  L2RootHash
	UncleHash   common.Hash    `json:"sha3Uncles"`
	Coinbase    common.Address `json:"miner"`
	Root        StateRoot      `json:"stateRoot"`
	TxHash      common.Hash    `json:"transactionsRoot"` // todo - include the synthetic deposits
	ReceiptHash common.Hash    `json:"receiptsRoot"`
	Bloom       types.Bloom    `json:"logsBloom"`
	Difficulty  *big.Int
	Number      *big.Int
	GasLimit    uint64
	GasUsed     uint64
	Time        uint64      `json:"timestamp"`
	Extra       []byte      `json:"extraData"`
	MixDigest   common.Hash `json:"mixHash"`
	Nonce       types.BlockNonce
	BaseFee     *big.Int

	// The custom Obscuro fields.
	Agg     common.Address // TODO - Can this be removed and replaced with the `Coinbase` field?
	L1Proof L1RootHash     // the L1 block used by the enclave to generate the current rollup
	R, S    *big.Int       // signature values
	// TODO: mark as deprecated Withdrawals are now contained within cross chain messages.
	Withdrawals        []Withdrawal
	CrossChainMessages []MessageBus.StructsCrossChainMessage `json:"crossChainMessages"`

	// The block hash of the latest block that has been scanned for cross chain messages.
	LatestInboundCrossChainHash common.Hash `json:"inboundCrossChainHash"`

	// The block height of the latest block that has been scanned for cross chain messages.
	LatestInboundCrossChainHeight *big.Int `json:"inboundCrossChainHeight"`
}

// Withdrawal - this is the withdrawal instruction that is included in the rollup header.
type Withdrawal struct {
	// Type      uint8 // the type of withdrawal. For now only ERC20. Todo - add this once more ERCs are supported
	Amount    *big.Int
	Recipient common.Address // the user account that will receive the money
	Contract  common.Address // the contract
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding excluding the signature.
func (b *BatchHeader) Hash() L2RootHash {
	cp := *b
	cp.R = nil
	cp.S = nil
	hash, err := rlpHash(cp)
	if err != nil {
		panic("err hashing batch header")
	}
	return hash
}

func (b *BatchHeader) ToRollupHeader() *RollupHeader {
	return &RollupHeader{
		ParentHash:                    b.ParentHash,
		UncleHash:                     b.UncleHash,
		Coinbase:                      b.Coinbase,
		Root:                          b.Root,
		TxHash:                        b.TxHash,
		ReceiptHash:                   b.ReceiptHash,
		Bloom:                         b.Bloom,
		Difficulty:                    b.Difficulty,
		Number:                        b.Number,
		GasLimit:                      b.GasLimit,
		GasUsed:                       b.GasUsed,
		Time:                          b.Time,
		Extra:                         b.Extra,
		MixDigest:                     b.MixDigest,
		Nonce:                         b.Nonce,
		BaseFee:                       b.BaseFee,
		Agg:                           b.Agg,
		L1Proof:                       b.L1Proof,
		R:                             b.R,
		S:                             b.S,
		Withdrawals:                   b.Withdrawals,
		CrossChainMessages:            b.CrossChainMessages,
		LatestInboundCrossChainHash:   b.LatestInboudCrossChainHash,
		LatestInboundCrossChainHeight: b.LatestInboundCrossChainHeight,
	}
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding excluding the signature.
func (r *RollupHeader) Hash() L2RootHash {
	cp := *r
	cp.R = nil
	cp.S = nil
	hash, err := rlpHash(cp)
	if err != nil {
		panic("err hashing rollup header")
	}
	return hash
}

func (r *RollupHeader) ToBatchHeader() *BatchHeader {
	return &BatchHeader{
		ParentHash:                    r.ParentHash,
		UncleHash:                     r.UncleHash,
		Coinbase:                      r.Coinbase,
		Root:                          r.Root,
		TxHash:                        r.TxHash,
		ReceiptHash:                   r.ReceiptHash,
		Bloom:                         r.Bloom,
		Difficulty:                    r.Difficulty,
		Number:                        r.Number,
		GasLimit:                      r.GasLimit,
		GasUsed:                       r.GasUsed,
		Time:                          r.Time,
		Extra:                         r.Extra,
		MixDigest:                     r.MixDigest,
		Nonce:                         r.Nonce,
		BaseFee:                       r.BaseFee,
		Agg:                           r.Agg,
		L1Proof:                       r.L1Proof,
		R:                             r.R,
		S:                             r.S,
		Withdrawals:                   r.Withdrawals,
		CrossChainMessages:            r.CrossChainMessages,
		LatestInboudCrossChainHash:    r.LatestInboundCrossChainHash,
		LatestInboundCrossChainHeight: r.LatestInboundCrossChainHeight,
	}
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
