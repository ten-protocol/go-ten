package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/ten-protocol/go-ten/tools/walletextension/encryption"
	"golang.org/x/crypto/acme/autocert"
)

const (
	CERT_CONTAINER_NAME = "certificates"
)

// CertStorageCosmosDB implements autocert.Cache interface using CosmosDB
type CertStorageCosmosDB struct {
	client         *azcosmos.Client
	certsContainer *azcosmos.ContainerClient
	encryptor      encryption.Encryptor
}

// EncryptedCertDocument represents the structure of a certificate document in CosmosDB
type EncryptedCertDocument struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

// NewCertStorageCosmosDB creates a new CosmosDB-based certificate storage
func NewCertStorageCosmosDB(connectionString string, encryptionKey []byte) (*CertStorageCosmosDB, error) {
	encryptor, err := encryption.NewEncryptor(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %w", err)
	}

	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	// Ensure database exists
	ctx := context.Background()
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Create container for certificates
	certsContainer, err := client.NewContainer(DATABASE_NAME, CERT_CONTAINER_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificates container: %w", err)
	}

	return &CertStorageCosmosDB{
		client:         client,
		certsContainer: certsContainer,
		encryptor:      *encryptor,
	}, nil
}

// Get retrieves a certificate data for the given key
func (c *CertStorageCosmosDB) Get(ctx context.Context, key string) ([]byte, error) {
	keyString, partitionKey := c.dbKey([]byte(key))

	itemResponse, err := c.certsContainer.ReadItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}

	var doc EncryptedCertDocument
	err = json.Unmarshal(itemResponse.Value, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return c.encryptor.Decrypt(doc.Data)
}

// Put stores certificate data with the given key
func (c *CertStorageCosmosDB) Put(ctx context.Context, key string, data []byte) error {
	keyString, partitionKey := c.dbKey([]byte(key))

	encryptedData, err := c.encryptor.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt certificate data: %w", err)
	}

	doc := EncryptedCertDocument{
		ID:   keyString,
		Data: encryptedData,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	_, err = c.certsContainer.UpsertItem(ctx, partitionKey, docJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to upsert certificate: %w", err)
	}

	return nil
}

// Delete removes certificate data for the given key
func (c *CertStorageCosmosDB) Delete(ctx context.Context, key string) error {
	keyString, partitionKey := c.dbKey([]byte(key))

	_, err := c.certsContainer.DeleteItem(ctx, partitionKey, keyString, nil)
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}

	return nil
}

// dbKey generates a consistent key for CosmosDB storage
func (c *CertStorageCosmosDB) dbKey(key []byte) (string, azcosmos.PartitionKey) {
	keyString := string(key)
	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	return keyString, partitionKey
}
