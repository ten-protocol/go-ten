package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
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
	enc, _ := rlp.EncodeToBytes(data)
	return types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	})
}

func TxData(tx *L1Tx) L1TxData {
	data := L1TxData{}
	rlp.DecodeBytes(tx.Data(), data)
	return data
}

type (
	EncodedL1Tx  []byte
	Transactions types.Transactions
)

type Header = types.Header

type Block = types.Block

// the encoded version of an ExtBlock
type EncodedBlock []byte

var GenesisHash = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")

func NewBlock(parent *Block, nonce uint64, nodeID common.Address, txs []*L1Tx) *Block {
	parentHash := GenesisHash
	if parent != nil {
		parentHash = parent.Hash()
	}

	header := Header{
		ParentHash:  parentHash,
		Coinbase:    nodeID,
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  nil,
		Number:      nil,
		GasLimit:    0,
		GasUsed:     0,
		Time:        0,
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}

	return types.NewBlock(&header, txs, nil, nil, nil)

	//header := Header{
	//	Nonce:      nonce,
	//	Miner:      nodeID,
	//	ParentHash: parentHash,
	//}
	//
	//return Block{
	//	Header:       &header,
	//	Transactions: txs,
	//}
}

var GenesisBlock = NewBlock(nil, 0, common.HexToAddress("0x0"), []*L1Tx{})

type EncryptedSharedEnclaveSecret []byte

type AttestationReport struct {
	Owner NodeID
	// todo public key
	// hash of code
	// other stuff
}
