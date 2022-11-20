package common

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"golang.org/x/crypto/sha3"
)

// Used to hash headers.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// Header is a public / plaintext struct that holds common properties of rollups..
// Making changes to this struct will require GRPC + GRPC Converters regen
type Header struct {
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
	Agg                           common.Address // TODO - Can this be removed and replaced with the `Coinbase` field?
	L1Proof                       L1RootHash     // the L1 block used by the enclave to generate the current rollup
	R, S                          *big.Int       // signature values
	Withdrawals                   []Withdrawal
	CrossChainMessages            []MessageBus.StructsCrossChainMessage
	LatestInboudCrossChainHash    common.Hash
	LatestInboundCrossChainHeight *big.Int
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
