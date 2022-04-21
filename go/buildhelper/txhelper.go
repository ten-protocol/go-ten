package buildhelper

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// sends a l1 tx that calls the contract and makes a deposit
func FundWallet(ethNode obscurocommon.L1Node, address common.Address) {
	fundAmt := 5000
	txData := obscurocommon.L1TxData{
		TxType: obscurocommon.DepositTx,
		Amount: uint64(fundAmt),
		Dest:   address,
	}
	tx := obscurocommon.NewL1Tx(txData)
	t, _ := obscurocommon.EncodeTx(tx)
	ethNode.BroadcastTx(t)
	fmt.Printf("Funded Address %s with %d \n", address, fundAmt)
}

func Transfer(obsNode host.Node, from common.Address, to common.Address, amount int) {
	txData := enclave.L2TxData{Type: enclave.TransferTx, From: from, To: to, Amount: uint64(amount)}

	// We should probably use a deterministic nonce instead, as in the L1.
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	enc, err := rlp.EncodeToBytes(txData)
	if err != nil {
		// TODO - Surface this error properly.
		panic(err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce.Uint64(),
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	})
	signedTx := wallet_mock.SignTx(tx, Addr1PK())
	encryptedTx := enclave.EncryptTx(signedTx)

	obsNode.ReceiveTx(encryptedTx)
	fmt.Printf("Transferred %d from: %s -> %s\n", amount, from, to)
}

func CheckBalance(obsNode host.Node, addr common.Address) {
	balance := obsNode.RPCBalance(addr)
	fmt.Printf("balance of %s: %d\n", addr.Hex(), balance)
}
