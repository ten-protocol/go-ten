package nodecommon

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/obscuro-playground/go/hashing"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/ethereum/go-ethereum/common"
)

type (
	StateRoot             = common.Hash
	L2Tx                  = types.Transaction
	EncryptedTx           []byte // A single transaction encrypted using the enclave's public key
	EncryptedTransactions []byte // A blob of encrypted transactions, as they're stored in the rollup.

	EncryptedParamsGetBalance     []byte // The params for an RPC getBalance request, as a JSON object encrypted with the public key of the enclave.
	EncryptedParamsCall           []byte // As above, but for an RPC call request.
	EncryptedParamsGetTxReceipt   []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedResponseGetBalance   []byte // The response for an RPC getBalance request, as a JSON object encrypted with the viewing key of the user.
	EncryptedResponseCall         []byte // As above, but for an RPC call request.
	EncryptedResponseGetTxReceipt []byte // As above, but for an RPC getTransactionReceipt request.
)

// Header is a public / plaintext struct that holds common properties of the Rollup
// Making changes to this struct will require GRPC + GRPC Converters regen
type Header struct {
	ParentHash  obscurocommon.L2RootHash
	Agg         common.Address
	Nonce       obscurocommon.Nonce
	L1Proof     obscurocommon.L1RootHash // the L1 block where the Parent was published
	Root        StateRoot
	TxHash      common.Hash // todo - include the synthetic deposits
	Number      *big.Int    // the rollup height
	Withdrawals []Withdrawal
	Bloom       types.Bloom
	ReceiptHash common.Hash
	Extra       []byte
}

type Withdrawal struct {
	Amount  uint64
	Address common.Address
}

// ExtRollup is used for communication between the enclave and the outside world.
type ExtRollup struct {
	Header          *Header
	EncryptedTxBlob EncryptedTransactions
}

// Rollup extends ExtRollup with additional fields.
// This parallels the Block/extblock split in Go Ethereum.
type Rollup struct {
	Header *Header

	hash atomic.Value
	size atomic.Value //nolint

	Transactions EncryptedTransactions
}

func (er ExtRollup) ToRollup() *Rollup {
	return &Rollup{
		Header:       er.Header,
		Transactions: er.EncryptedTxBlob,
	}
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() obscurocommon.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(obscurocommon.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)

	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() obscurocommon.L2RootHash {
	hash, err := hashing.RLPHash(h)
	if err != nil {
		log.Error("err hashing the l2roothash")
	}
	return hash
}
