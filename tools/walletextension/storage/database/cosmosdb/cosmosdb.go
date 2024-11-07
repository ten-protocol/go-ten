package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	dbcommon "github.com/ten-protocol/go-ten/tools/walletextension/storage/database/common"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/encryption"
)

/*
This is a CosmosDB implementation of the Storage interface.

We need to make sure we have a CosmosDB account and a database created before using this.

Quick summary of the CosmosDB setup:
- Create a CosmosDB account (Azure Cosmos DB for NoSQL)
- account name should follow the format: <network_type>-gateway-cosmosdb
- use Serverless capacity mode for testnets
- go to "Data Explorer" in the CosmosDB account and create new database named "gatewayDB"
- inside the database create a container named "users" with partition key of "/id"
- to get your connection string go to settings -> keys -> primary connection string

*/

// CosmosDB struct represents the CosmosDB storage implementation
type CosmosDB struct {
	client         *azcosmos.Client
	usersContainer *azcosmos.ContainerClient
	encryptor      encryption.Encryptor
}

// EncryptedDocument struct is used to store encrypted user data in CosmosDB
// We use this structure to add an extra layer of security by encrypting the actual user data
// The 'ID' field is used as the document ID and partition key in CosmosDB
// The 'Data' field contains the base64-encoded encrypted user data
type EncryptedDocument struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

// Constants for the CosmosDB database and container names
const (
	DATABASE_NAME        = "gatewayDB"
	USERS_CONTAINER_NAME = "users"
)

func NewCosmosDB(connectionString string, encryptionKey []byte) (*CosmosDB, error) {
	// Create encryptor
	encryptor, err := encryption.NewEncryptor(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %w", err)
	}

	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	// Create database if it doesn't exist
	ctx := context.Background()
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Create container client for users container
	usersContainer, err := client.NewContainer(DATABASE_NAME, USERS_CONTAINER_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to create users container: %w", err)
	}

	return &CosmosDB{
		client:         client,
		usersContainer: usersContainer,
		encryptor:      *encryptor,
	}, nil
}

func (c *CosmosDB) AddUser(userID []byte, privateKey []byte) error {
	user := dbcommon.GWUserDB{
		UserId:     userID,
		PrivateKey: privateKey,
		Accounts:   []dbcommon.GWAccountDB{},
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	ciphertext, err := c.encryptor.Encrypt(userJSON)
	if err != nil {
		return fmt.Errorf("failed to encrypt user data: %w", err)
	}

	key := c.encryptor.HashWithHMAC(userID)
	keyString := hex.EncodeToString(key)

	// Create an EncryptedDocument struct to store in CosmosDB
	doc := EncryptedDocument{
		ID:   keyString,
		Data: ciphertext,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	ctx := context.Background()
	_, err = c.usersContainer.CreateItem(ctx, partitionKey, docJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (c *CosmosDB) DeleteUser(userID []byte) error {
	key := c.encryptor.HashWithHMAC(userID)
	keyString := hex.EncodeToString(key)
	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	ctx := context.Background()

	_, err := c.usersContainer.DeleteItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *CosmosDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	key := c.encryptor.HashWithHMAC(userID)
	keyString := hex.EncodeToString(key)
	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	ctx := context.Background()

	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	var doc EncryptedDocument
	err = json.Unmarshal(itemResponse.Value, &doc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal document: %w", err)
	}

	data, err := c.encryptor.Decrypt(doc.Data)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	var user dbcommon.GWUserDB
	err = json.Unmarshal(data, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Add the new account
	newAccount := dbcommon.GWAccountDB{
		AccountAddress: accountAddress,
		Signature:      signature,
		SignatureType:  int(signatureType),
	}
	user.Accounts = append(user.Accounts, newAccount)

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling updated user: %w", err)
	}

	ciphertext, err := c.encryptor.Encrypt(userJSON)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated user data: %w", err)
	}

	// Update the document
	doc.Data = ciphertext

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal updated document: %w", err)
	}

	// Replace the item in the container
	_, err = c.usersContainer.ReplaceItem(ctx, partitionKey, keyString, docJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to update user with new account: %w", err)
	}
	return nil
}

func (c *CosmosDB) GetUser(userID []byte) (*common.GWUser, error) {
	key := c.encryptor.HashWithHMAC(userID)
	keyString := hex.EncodeToString(key)
	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	ctx := context.Background()

	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		return nil, errutil.ErrNotFound
	}

	var doc EncryptedDocument
	err = json.Unmarshal(itemResponse.Value, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	data, err := c.encryptor.Decrypt(doc.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	var user dbcommon.GWUserDB
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}
	return user.ToGWUser(), nil
}

// GetEncryptionKey returns the encryption key used by the CosmosDB instance
func (c *CosmosDB) GetEncryptionKey() []byte {
	return c.encryptor.GetKey()
}
