package walletmock

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type Wallet struct {
	Address common.Address
	// TODO - Store key securely. Geth stores the key encrypted on disk.
	Key keystore.Key
}

func New() Wallet {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		panic(fmt.Sprintf("Could not generate keypair for wallet: %v", err))
	}

	key := keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}

	return Wallet{Address: key.Address, Key: key}
}

// NewEncryptedL2Transfer creates an encrypted L2Tx of type TransferTx.
func NewEncryptedL2Transfer(from common.Address, dest common.Address, amount uint64) common2.EncryptedTx {
	txData := enclave.L2TxData{Type: enclave.TransferTx, From: from, To: dest, Amount: amount}
	return newEncryptedL2Tx(txData)
}

// NewEncryptedL2Withdrawal creates an encrypted L2Tx of type WithdrawalTx.
func NewEncryptedL2Withdrawal(from common.Address, amount uint64) common2.EncryptedTx {
	txData := enclave.L2TxData{Type: enclave.WithdrawalTx, From: from, Amount: amount}
	return newEncryptedL2Tx(txData)
}

// newL2Tx creates an L2Tx, using a random nonce (to avoid hash collisions) and with the L2 data encoded in the
// transaction's data field, then encrypts it.
func newEncryptedL2Tx(data enclave.L2TxData) common2.EncryptedTx {
	// We should probably use a deterministic nonce instead, as in L1.
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	enc, err := rlp.EncodeToBytes(data)
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

	return enclave.EncryptTx(tx)
}
