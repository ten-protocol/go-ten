package walletmock

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type Wallet struct {
	Address common.Address
	// TODO - Store key securely. Geth stores the key encrypted on disk.
	Key   keystore.Key
	nonce uint64
}

func (w *Wallet) readNonce(cl *obscuroclient.Client) uint64 {
	var result uint64
	err := (*cl).Call(&result, obscuroclient.RPCNonce, w.Address)
	if err != nil {
		panic(err)
	}
	if result == 0 {
		return 0
	}
	return result
}

func (w *Wallet) NextNonce(cl *obscuroclient.Client) uint64 {
	// only returns the nonce when the previous transaction was recorded
	for {
		result := w.readNonce(cl)
		if result == w.nonce {
			atomic.AddUint64(&w.nonce, 1)
			return result
		}
		time.Sleep(time.Millisecond)
	}
}

func New(privateKeyECDSA *ecdsa.PrivateKey) *Wallet {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	key := keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}

	return &Wallet{Address: key.Address, Key: key}
}

// SignTx returns a copy of the enclave.L2Tx signed with the provided ecdsa.PrivateKey
func SignTx(tx *nodecommon.L2Tx, privateKey *ecdsa.PrivateKey) *nodecommon.L2Tx {
	signer := types.NewLondonSigner(big.NewInt(evm.ChainID))
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		panic(fmt.Errorf("could not sign transaction: %w", err))
	}
	return signedTx
}
