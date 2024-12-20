package storage

import (
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
	"golang.org/x/crypto/acme/autocert"
)

// CertStorage defines the interface for certificate storage
type CertStorage interface {
	autocert.Cache
}

// NewCertStorage creates a new certificate storage instance based on the database type
func NewCertStorage(dbType, dbConnectionURL string, randomKey []byte, encryptionEnabled bool, logger gethlog.Logger) (CertStorage, error) {
	switch dbType {
	case "cosmosDB":
		return cosmosdb.NewCertStorageCosmosDB(dbConnectionURL, randomKey, encryptionEnabled)
	default:
		return autocert.DirCache("/data/certs"), nil
	}
}
