package obscurocommon

import (
	"math/big"
	"math/rand"

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
	Attestation EncodedAttestationReport

	// if the type is deposit
	Amount uint64
	Dest   common.Address
}

type L1Tx = types.Transaction

func NewL1Tx(data L1TxData) (*L1Tx, error) {
	enc, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, err
	}
	return types.NewTx(&types.LegacyTx{
		Nonce:    rand.Uint64(), //nolint:gosec
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	}), nil
}

func TxData(tx *L1Tx) (*L1TxData, error) {
	data := L1TxData{}
	err := rlp.DecodeBytes(tx.Data(), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

type (
	EncodedL1Tx  []byte
	Transactions types.Transactions
)

// the encoded version of an ExtBlock
type EncodedBlock []byte

var (
	GenesisHash  = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	GenesisBlock = NewBlock(nil, common.HexToAddress("0x0"), []*L1Tx{})
)

func NewBlock(parent *types.Block, nodeID common.Address, txs []*L1Tx) *types.Block {
	parentHash := GenesisHash
	height := L1GenesisHeight
	if parent != nil {
		parentHash = parent.Hash()
		height = parent.NumberU64() + 1
	}

	header := types.Header{
		ParentHash:  parentHash,
		UncleHash:   common.Hash{},
		Coinbase:    nodeID,
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  big.NewInt(0),
		Number:      big.NewInt(int64(height)),
		GasLimit:    0,
		GasUsed:     0,
		Time:        0,
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}

	return types.NewBlock(&header, txs, nil, nil, &trie.StackTrie{})
}

type EncryptedSharedEnclaveSecret []byte

type EncodedAttestationReport []byte

// AttestationReport represents a signed attestation report from a TEE and some metadata about the source of it to verify it
type AttestationReport struct {
	Report []byte         // the signed bytes of the report which includes some encrypted identifying data
	PubKey []byte         // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	Owner  common.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data
}

type AttestationVerification struct {
	ReportData []byte // the data embedded in the report at the time it was produced (up to 64bytes)
}
