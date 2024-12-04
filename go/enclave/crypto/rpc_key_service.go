package crypto

import (
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
)

const rpcSuffix = 1

// RPCKeyService - manages the "TEN - RPC key" used by clients (like the TEN gateway) to make RPC requests
type RPCKeyService struct {
	privKey             *ecies.PrivateKey
	sharedSecretService *SharedSecretService
	logger              gethlog.Logger
}

func NewRPCKeyService(sharedSecretService *SharedSecretService, logger gethlog.Logger) *RPCKeyService {
	s := &RPCKeyService{
		sharedSecretService: sharedSecretService,
		logger:              logger,
	}
	if sharedSecretService.IsInitialised() {
		err := s.Initialise()
		if err != nil {
			logger.Crit("Failed to initialise RPC key service ", log.ErrKey, err)
			return nil
		}
	}
	return s
}

// Initialise - called when the shared secret is available
func (s *RPCKeyService) Initialise() error {
	// the key is derived from the shared secret to allow transactions to be broadcast
	bytes := s.sharedSecretService.ExtendEntropy([]byte{byte(rpcSuffix)})
	ecdsaKey, err := gethcrypto.ToECDSA(bytes)
	if err != nil {
		return err
	}
	s.privKey = ecies.ImportECDSA(ecdsaKey)
	return nil
}

func (s *RPCKeyService) DecryptRPCRequest(bytes []byte) ([]byte, error) {
	return s.privKey.Decrypt(bytes, nil, nil)
}

func (s *RPCKeyService) PublicKey() []byte {
	return gethcrypto.CompressPubkey(s.privKey.PublicKey.ExportECDSA())
}
