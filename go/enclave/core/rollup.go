package core

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Rollup Data structure only for the internal use of the enclave since transactions are in clear
// Making changes to this struct will require GRPC + GRPC Converters regen
type Rollup struct {
	Header *common.Header

	hash atomic.Value
	// size   atomic.Value

	Transactions []*common.L2Tx
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() common.L2RootHash {
	// Temporarily disabling the caching of the hash because it's causing bugs.
	// Transforming a Rollup to an ExtRollup and then back to a Rollup will generate a different hash if caching is enabled.
	// Todo - re-enable
	//if hash := r.hash.Load(); hash != nil {
	//	return hash.(common.L2RootHash)
	//}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

func (r *Rollup) NumberU64() uint64 { return r.Header.Number.Uint64() }
func (r *Rollup) Number() *big.Int  { return new(big.Int).Set(r.Header.Number) }

func (r *Rollup) ToExtRollup(transactionBlobCrypto crypto.TransactionBlobCrypto) common.ExtRollup {
	txHashes := make([]gethcommon.Hash, len(r.Transactions))
	for idx, tx := range r.Transactions {
		txHashes[idx] = tx.Hash()
	}

	return common.ExtRollup{
		Header:          r.Header,
		TxHashes:        txHashes,
		EncryptedTxBlob: transactionBlobCrypto.Encrypt(r.Transactions),
	}
}

func ToEnclaveRollup(encryptedRollup *common.ExtRollupWithHash, transactionBlobCrypto crypto.TransactionBlobCrypto) *Rollup {
	return &Rollup{
		Header:       encryptedRollup.Header,
		Transactions: transactionBlobCrypto.Decrypt(encryptedRollup.EncryptedTxBlob),
	}
}

func EmptyRollup(agg gethcommon.Address, parent *common.Header, blkHash gethcommon.Hash, nonce common.Nonce) *Rollup {
	h := common.Header{
		Agg:        agg,
		ParentHash: parent.Hash(),
		L1Proof:    blkHash,
		Number:     big.NewInt(int64(parent.Number.Uint64() + 1)),
		// TODO - Consider how this time should align with the time of the L1 block used as proof.
		Time: uint64(time.Now().Unix()),
		// generate true randomness inside the enclave.
		// note that this randomness will be published in the header of the rollup.
		// the randomness exposed to smart contract is combining this with the shared secret.
		MixDigest: gethcommon.BytesToHash(crypto.GeneratePublicRandomness()),
	}
	r := Rollup{
		Header: &h,
	}
	return &r
}

// NewRollup - produces a new rollup. only used for genesis. todo - review
func NewRollup(blkHash gethcommon.Hash, parent *Rollup, height uint64, a gethcommon.Address, txs []*common.L2Tx, withdrawals []common.Withdrawal, nonce common.Nonce, state common.StateRoot) *Rollup {
	parentHash := common.GenesisHash
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := common.Header{
		Agg:         a,
		ParentHash:  parentHash,
		L1Proof:     blkHash,
		Root:        state,
		TxHash:      types.EmptyRootHash,
		Number:      big.NewInt(int64(height)),
		Withdrawals: withdrawals,
		ReceiptHash: types.EmptyRootHash,
		Time:        uint64(time.Now().Unix()),
	}
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return &r
}
