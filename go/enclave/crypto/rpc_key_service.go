package crypto

import (
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/status-im/keycard-go/hexutils"
)

const rpcSuffix = 1

// RPCKeyService - manages the "TEN - RPC key" used by clients (like the TEN gateway) to make RPC requests
type RPCKeyService struct {
	privKey *ecies.PrivateKey
}

func NewRPCKeyService(service *SharedSecretService) *RPCKeyService {
	// the key is derived from the shared secret to allow transactions to be broadcast
	// bytes := service.ExtendEntropy([]byte{byte(rpcSuffix)})
	// todo - identify where we have the hardcoded public key - and create the logic to get the pub key
	bytes := hexutils.HexToBytes("81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207")
	ecdsaKey, err := gethcrypto.ToECDSA(bytes)
	if err != nil {
		panic(err)
	}
	return &RPCKeyService{privKey: ecies.ImportECDSA(ecdsaKey)}
}

func (s RPCKeyService) DecryptRPCRequest(bytes []byte) ([]byte, error) {
	return s.privKey.Decrypt(bytes, nil, nil)
}
