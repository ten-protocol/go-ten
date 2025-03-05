package ethereummock

import (
	"bytes"
	"encoding/gob"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
)

var (
	NetworkConfigAddr      = datagenerator.RandomAddress()
	MessageBusAddr         = datagenerator.RandomAddress()
	DepositTxAddr          = datagenerator.RandomAddress()
	RollupTxAddr           = datagenerator.RandomAddress()
	StoreSecretTxAddr      = datagenerator.RandomAddress()
	RequestSecretTxAddr    = datagenerator.RandomAddress()
	InitializeSecretTxAddr = datagenerator.RandomAddress()
	GrantSeqTxAddr         = datagenerator.RandomAddress()
	CrossChainAddr         = datagenerator.RandomAddress()
)

func DecodeTx(tx *types.Transaction) common.L1TenTransaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify the L1 operation (rollup/storesecret/requestsecret)
	// the mock implementation does not process contracts
	// so this is a way that we can differentiate different contract calls
	var t common.L1TenTransaction
	switch tx.To().Hex() {
	case StoreSecretTxAddr.Hex():
		t = &common.L1RespondSecretTx{}
	case DepositTxAddr.Hex():
		t = &common.L1DepositTx{}
	case RequestSecretTxAddr.Hex():
		t = &common.L1RequestSecretTx{}
	case InitializeSecretTxAddr.Hex():
		t = &common.L1InitializeSecretTx{}
	case GrantSeqTxAddr.Hex():
		// this tx is empty and entirely mocked, no need to decode
		return &common.L1PermissionSeqTx{}
	default:
		panic("unexpected type")
	}

	// decode to interface implementation
	if err := dec.Decode(t); err != nil {
		panic(err)
	}

	return t
}

func EncodeTx(tx common.L1TenTransaction, opType gethcommon.Address) types.TxData {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	// the mock implementation does not process contract calls
	// this uses the To address to distinguish between different contract calls / different l1 transactions
	return &types.LegacyTx{
		Data: buf.Bytes(),
		To:   &opType,
	}
}
