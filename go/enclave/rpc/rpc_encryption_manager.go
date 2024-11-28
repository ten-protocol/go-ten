package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ten-protocol/go-ten/go/common/privacy"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/gas"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/l2chain"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// EncryptionManager manages the decryption and encryption of enclave comms.
type EncryptionManager struct {
	chain                  l2chain.ObscuroChain
	enclavePrivateKeyECIES *ecies.PrivateKey
	storage                storage.Storage
	cacheService           *storage.CacheService
	registry               components.BatchRegistry
	processors             *crosschain.Processors
	mempool                *txpool.TxPool
	gasOracle              gas.Oracle
	blockResolver          storage.BlockResolver
	l1BlockProcessor       components.L1BlockProcessor
	config                 *enclaveconfig.EnclaveConfig
	logger                 gethlog.Logger
	storageSlotWhitelist   *privacy.Whitelist
}

func NewEncryptionManager(enclavePrivateKeyECIES *ecies.PrivateKey, storage storage.Storage, cacheService *storage.CacheService, registry components.BatchRegistry, mempool *txpool.TxPool, processors *crosschain.Processors, config *enclaveconfig.EnclaveConfig, oracle gas.Oracle, blockResolver storage.BlockResolver, l1BlockProcessor components.L1BlockProcessor, chain l2chain.ObscuroChain, logger gethlog.Logger) *EncryptionManager {
	return &EncryptionManager{
		storage:                storage,
		cacheService:           cacheService,
		registry:               registry,
		processors:             processors,
		chain:                  chain,
		config:                 config,
		blockResolver:          blockResolver,
		l1BlockProcessor:       l1BlockProcessor,
		gasOracle:              oracle,
		logger:                 logger,
		enclavePrivateKeyECIES: enclavePrivateKeyECIES,
		storageSlotWhitelist:   privacy.NewWhitelist(),
		mempool:                mempool,
	}
}

// DecryptBytes decrypts the bytes with the enclave's private key.
func (rpc *EncryptionManager) DecryptBytes(encryptedBytes []byte) ([]byte, error) {
	bytes, err := rpc.enclavePrivateKeyECIES.Decrypt(encryptedBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with enclave private key. Cause: %w", err)
	}

	return bytes, nil
}
