package obscuro

import (
	c "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	"io"
	"simulation/common"
	"sync"
	"sync/atomic"
)

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncodedL2Tx []byte

//todo
//type EncryptedTransactionBlob []byte
type EncryptedTransactionBlob []L2Tx

// This datastructure is in plaintext
type Header struct {
	ParentHash  common.L2RootHash
	Agg         common.NodeId
	Nonce       common.Nonce
	L1Proof     common.L1RootHash // the L1 block where the Parent was published
	State       StateRoot
	Withdrawals []Withdrawal
}

type Rollup struct {
	// header
	Header *Header

	hash   atomic.Value
	height atomic.Value
	size   atomic.Value
	//Height       uint32
	//RootHash     common.L2RootHash
	//CreationTime time.Time

	// payload
	Transactions []L2Tx
}

// Transfers and Withdrawals for now
type L2TxType uint64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

type Withdrawal struct {
	Amount  uint64
	Address common.Address
}

// todo - signing
type L2Tx struct {
	Id     common.TxHash
	TxType L2TxType
	Amount uint64
	From   common.Address
	To     common.Address
}

const GenesisHash = "1000000000000000000000000000000000000000000000000000000000000000"

var GenesisRollup = NewRollup(&common.GenesisBlock, nil, 0, []L2Tx{}, []Withdrawal{}, common.GenerateNonce(), "")
var encodedGenesis, _ = GenesisRollup.Encode()
var GenesisTx = common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: encodedGenesis}

func (r Rollup) Proof(l1BlockResolver common.BlockResolver) *common.Block {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}

// ProofHeight - return the height of the L1 proof, or 0 - if the block is not known
//todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r Rollup) ProofHeight(l1BlockResolver common.BlockResolver) int {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		return 0
	}
	return v.Height(l1BlockResolver)
}

func NewRollup(b *common.Block, parent *Rollup, a common.NodeId, txs []L2Tx, withdrawals []Withdrawal, nonce common.Nonce, state StateRoot) Rollup {
	parentHash := c.HexToHash(GenesisHash)
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := Header{
		Agg:         a,
		ParentHash:  parentHash,
		L1Proof:     b.Hash(),
		Nonce:       nonce,
		State:       state,
		Withdrawals: withdrawals,
	}
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return r
}

type extrollup struct {
	Header *Header
	Txs    EncryptedTransactionBlob
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() common.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(common.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.L2RootHash {
	return rlpHash(h)
}

//func (h *Header) Height() uint32 {
//	return rlpHash(h)
//}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h common.L2RootHash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// DecodeRLP decodes the Ethereum
func (b *Rollup) DecodeRLP(s *rlp.Stream) error {
	var eb extrollup
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.Header, b.Transactions = eb.Header, eb.Txs
	b.size.Store(c.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b Rollup) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extrollup{
		Header: b.Header,
		Txs:    b.Transactions,
	})
}

func (b *Rollup) Height(db Db) int {
	if height := b.height.Load(); height != nil {
		return height.(int)
	}
	if b.Hash() == GenesisRollup.Hash() {
		b.height.Store(common.L2GenesisHeight)
		return common.L2GenesisHeight
	}
	v := b.Parent(db).Height(db) + 1
	b.height.Store(v)
	return v
}
