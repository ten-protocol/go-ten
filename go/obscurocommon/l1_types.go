package obscurocommon

import (
	"math"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// L1TxType - Just two types of relevant L1 transactions: Deposits and Rollups
// this does not actually exist in the real implementation
type L1TxType uint8

const (
	DepositTx L1TxType = iota
	RollupTx
	StoreSecretTx
	RequestSecretTx
)

// For now all the fields are placeholders for arguments sent to the management contract
type L1TxData struct {
	TxType L1TxType

	// if the type is rollup
	// todo -payload
	Rollup EncodedRollup

	Secret      EncryptedSharedEnclaveSecret
	Attestation AttestationReport

	// if the type is deposit
	Amount uint64
	Dest   common.Address
}

type L1Tx = types.Transaction

func NewL1Tx(data L1TxData) *L1Tx {
	enc, err := rlp.EncodeToBytes(data)
	if err != nil {
		panic(err)
	}
	return types.NewTx(&types.LegacyTx{
		Nonce:    rand.Uint64(), //nolint:gosec
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	})
}

func TxData(tx *L1Tx) L1TxData {
	data := L1TxData{}
	err := rlp.DecodeBytes(tx.Data(), &data)
	if err != nil {
		panic(err)
	}
	return data
}

type (
	EncodedL1Tx  []byte
	Transactions types.Transactions
)

// the encoded version of an ExtBlock
type EncodedBlock []byte

var (
	GenesisBlock = core.DefaultGenesisBlock().ToBlock(nil)
	GenesisHash  = GenesisBlock.Hash()
)

func NewBlock(parent *types.Block, nodeID common.Address, txs []*L1Tx) *types.Block {
	blockNum := big.NewInt(0).Add(parent.Number(), big.NewInt(1))

	// todo - joel - fold into a function
	var baseFee *big.Int
	londonBlock := uint64(1)
	if blockNum.Uint64() <= londonBlock {
		baseFee = big.NewInt(1000000000)
	} else {
		x := float64(parent.BaseFee().Uint64()) * 0.875
		baseFee = big.NewInt(int64(math.Ceil(x)))
	}

	header := types.Header{
		ParentHash:  parent.Hash(),
		Coinbase:    nodeID,
		Root:        common.HexToHash("d7f8974fb5ac78d9ac099b9ad5018bedc2ce0a72dad1827a1709da30580f0544"),
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  big.NewInt(0),
		Number:      blockNum,
		GasLimit:    10000,
		GasUsed:     0,
		Time:        0,
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       types.EncodeNonce(0), // 0 is the Beacon nonce.
		BaseFee:     baseFee,
	}

	return types.NewBlock(&header, txs, nil, nil, &trie.StackTrie{})
}

type EncryptedSharedEnclaveSecret []byte

type AttestationReport struct {
	Owner common.Address
	// todo public key
	// hash of code
	// other stuff
}
