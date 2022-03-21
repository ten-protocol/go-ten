package nodecommon

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/hashing"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/ethereum/go-ethereum/common"
)

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncryptedTx []byte

type EncryptedTransactions []EncryptedTx

// Header is a public / plaintext struct that holds common properties of the Rollup
type Header struct {
	ParentHash  obscurocommon.L2RootHash
	Agg         common.Address
	Nonce       obscurocommon.Nonce
	L1Proof     obscurocommon.L1RootHash // the L1 block where the Parent was published
	State       StateRoot
	Withdrawals []Withdrawal
}

type Withdrawal struct {
	Amount  uint64
	Address common.Address
}

// ExtRollup Data structure that is used to communicate between the enclave and the outside world
type ExtRollup struct {
	Header *Header
	Txs    EncryptedTransactions
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() obscurocommon.L2RootHash {
	hash, err := hashing.RLPHash(h)
	if err != nil {
		// todo - log / surface these
		fmt.Printf("err hashing the l2roothash")
	}
	return hash
}
