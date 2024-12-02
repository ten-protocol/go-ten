package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ten-protocol/go-ten/go/common/privacy"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/gas"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/l2chain"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

// EncryptionManager manages the decryption and encryption of enclave comms.
type EncryptionManager struct {
	chain                l2chain.ObscuroChain
	rpcKeyService        *crypto.RPCKeyService
	storage              storage.Storage
	cacheService         *storage.CacheService
	registry             components.BatchRegistry
	processors           *crosschain.Processors
	mempool              *txpool.TxPool
	gasOracle            gas.Oracle
	blockResolver        storage.BlockResolver
	l1BlockProcessor     components.L1BlockProcessor
	config               *enclaveconfig.EnclaveConfig
	logger               gethlog.Logger
	storageSlotWhitelist *privacy.Whitelist
}

func NewEncryptionManager(storage storage.Storage, cacheService *storage.CacheService, registry components.BatchRegistry, mempool *txpool.TxPool, processors *crosschain.Processors, config *enclaveconfig.EnclaveConfig, oracle gas.Oracle, blockResolver storage.BlockResolver, l1BlockProcessor components.L1BlockProcessor, chain l2chain.ObscuroChain, rpcKeyService *crypto.RPCKeyService, logger gethlog.Logger) *EncryptionManager {
	return &EncryptionManager{
		storage:              storage,
		cacheService:         cacheService,
		registry:             registry,
		processors:           processors,
		chain:                chain,
		config:               config,
		blockResolver:        blockResolver,
		l1BlockProcessor:     l1BlockProcessor,
		gasOracle:            oracle,
		logger:               logger,
		rpcKeyService:        rpcKeyService,
		storageSlotWhitelist: privacy.NewWhitelist(),
		mempool:              mempool,
	}
}

// DecryptBytes decrypts the bytes with the enclave's private key.
func (rpc *EncryptionManager) DecryptBytes(encryptedBytes []byte) ([]byte, error) {
	bytes, err := rpc.rpcKeyService.DecryptRPCRequest(encryptedBytes)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with enclave private key. Cause: %w", err)
	}

	return bytes, nil
}
