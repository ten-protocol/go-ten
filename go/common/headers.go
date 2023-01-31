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
	Agg                common.Address                        // TODO - Can this be removed and replaced with the `Coinbase` field?
	L1Proof            L1RootHash                            // the L1 block used by the enclave to generate the current batch
	R, S               *big.Int                              // signature values
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
	Agg                common.Address                        // TODO - Can this be removed and replaced with the `Coinbase` field?
	L1Proof            L1RootHash                            // the L1 block used by the enclave to generate the current rollup
	R, S               *big.Int                              // signature values
	CrossChainMessages []MessageBus.StructsCrossChainMessage `json:"crossChainMessages"`
	HeadBatchHash      common.Hash                           // The latest batch included in this rollup.

	// The block hash of the latest block that has been scanned for cross chain messages.
	LatestInboundCrossChainHash common.Hash `json:"inboundCrossChainHash"`

	// The block height of the latest block that has been scanned for cross chain messages.
	LatestInboundCrossChainHeight *big.Int `json:"inboundCrossChainHeight"`
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

func (b *BatchHeader) ToRollupHeader(parentRollupHeader *RollupHeader) *RollupHeader {
	header := &RollupHeader{
		// The fields copied directly from the batch header.
		ParentHash:                    b.ParentHash,
		UncleHash:                     b.UncleHash,
		Coinbase:                      b.Coinbase,
		Root:                          b.Root,
		ReceiptHash:                   b.ReceiptHash,
		Bloom:                         b.Bloom,
		Difficulty:                    b.Difficulty,
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
		CrossChainMessages:            b.CrossChainMessages,
		LatestInboundCrossChainHash:   b.LatestInboudCrossChainHash,
		LatestInboundCrossChainHeight: b.LatestInboundCrossChainHeight,

		// The fields that differ from the batch header.
		HeadBatchHash: b.TxHash,
	}

	// We set the parent hash to the genesis hash, or to the parent hash if one is provided.
	// todo - joel - is this the correct handling of the genesis rollup?
	header.ParentHash = L2RootHash{}
	header.Number = big.NewInt(int64(L2GenesisHeight))
	if parentRollupHeader != nil {
		header.ParentHash = parentRollupHeader.Hash()
		header.Number = big.NewInt(0).Add(parentRollupHeader.Number, big.NewInt(1))
	}

	return header
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
