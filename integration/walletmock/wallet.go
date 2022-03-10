package walletmock

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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

// NewL2Transfer creates an enclave.L2Tx of type enclave.TransferTx
func NewL2Transfer(from common.Address, dest common.Address, amount uint64) *enclave.L2Tx {
	txData := enclave.L2TxData{Type: enclave.TransferTx, From: from, To: dest, Amount: amount}
	return newL2Tx(txData)
}

// NewL2Withdrawal creates an enclave.L2Tx of type enclave.WithdrawalTx
func NewL2Withdrawal(from common.Address, amount uint64) *enclave.L2Tx {
	txData := enclave.L2TxData{Type: enclave.WithdrawalTx, From: from, Amount: amount}
	return newL2Tx(txData)
}

// newL2Tx creates an enclave.L2Tx.
//
// A random nonce is used to avoid hash collisions. The enclave.L2TxData is encoded and stored in the transaction's
// data field.
func newL2Tx(data enclave.L2TxData) *enclave.L2Tx {
	// We should probably use a deterministic nonce instead, as in the L1.
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	enc, err := rlp.EncodeToBytes(data)
	if err != nil {
		// TODO - Surface this error properly.
		panic(err)
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce.Uint64(),
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	})
}

// SignTx returns a copy of the enclave.L2Tx signed with the provided ecdsa.PrivateKey
func SignTx(tx *enclave.L2Tx, privateKey *ecdsa.PrivateKey) *enclave.L2Tx {
	signer := types.NewLondonSigner(big.NewInt(enclave.ChainID))
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		panic(fmt.Errorf("could not sign transaction: %w", err))
	}
	return signedTx
}
