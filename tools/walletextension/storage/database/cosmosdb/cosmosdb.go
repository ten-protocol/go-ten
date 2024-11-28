package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	dbcommon "github.com/ten-protocol/go-ten/tools/walletextension/storage/database/common"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

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

// userWithETag struct is used to store the user data along with its ETag
type userWithETag struct {
	user dbcommon.GWUserDB
	etag azcore.ETag
}

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
	ctx := context.Background()
	keyString, partitionKey := c.dbKey(userID)

	user := dbcommon.GWUserDB{
		UserId:     userID,
		PrivateKey: privateKey,
		Accounts:   []dbcommon.GWAccountDB{},
	}
	docJSON, err := c.createEncryptedDoc(user, keyString)
	if err != nil {
		return err
	}

	_, err = c.usersContainer.CreateItem(ctx, partitionKey, docJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (c *CosmosDB) DeleteUser(userID []byte) error {
	ctx := context.Background()

	keyString, partitionKey := c.dbKey(userID)
	_, err := c.usersContainer.DeleteItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *CosmosDB) AddSessionKey(userID []byte, key common.GWSessionKey) error {
	ctx := context.Background()

	user, err := c.getUserDB(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.user.SessionKey = &dbcommon.GWSessionKeyDB{
		PrivateKey: crypto.FromECDSA(key.PrivateKey.ExportECDSA()),
		Account: dbcommon.GWAccountDB{
			AccountAddress: key.Account.Address.Bytes(),
			Signature:      key.Account.Signature,
			SignatureType:  int(key.Account.SignatureType),
		},
	}
	return c.updateUser(ctx, user.user)
}

func (c *CosmosDB) ActivateSessionKey(userID []byte, active bool) error {
	ctx := context.Background()

	user, err := c.getUserDB(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.user.ActiveSK = active
	return c.updateUser(ctx, user.user)
}

func (c *CosmosDB) RemoveSessionKey(userID []byte) error {
	ctx := context.Background()

	user, err := c.getUserDB(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.user.SessionKey = nil
	return c.updateUser(ctx, user.user)
}

func (c *CosmosDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	ctx := context.Background()

	user, err := c.getUserDB(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Add the new account
	newAccount := dbcommon.GWAccountDB{
		AccountAddress: accountAddress,
		Signature:      signature,
		SignatureType:  int(signatureType),
	}
	user.user.Accounts = append(user.user.Accounts, newAccount)

	return c.updateUser(ctx, user.user)
}

func (c *CosmosDB) GetUser(userID []byte) (*common.GWUser, error) {
	user, err := c.getUserDB(userID)
	if err != nil {
		return nil, err
	}
	return user.user.ToGWUser()
}

func (c *CosmosDB) getUserDB(userID []byte) (userWithETag, error) {
	keyString, partitionKey := c.dbKey(userID)

	ctx := context.Background()

	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		return userWithETag{}, err
	}

	var doc EncryptedDocument
	err = json.Unmarshal(itemResponse.Value, &doc)
	if err != nil {
		return userWithETag{}, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	data, err := c.encryptor.Decrypt(doc.Data)
	if err != nil {
		return userWithETag{}, fmt.Errorf("failed to decrypt data: %w", err)
	}

	var user dbcommon.GWUserDB
	err = json.Unmarshal(data, &user)
	if err != nil {
		return userWithETag{}, fmt.Errorf("failed to unmarshal user data: %w", err)
	}
	return userWithETag{user: user, etag: itemResponse.ETag}, nil
}

func (c *CosmosDB) updateUser(ctx context.Context, user dbcommon.GWUserDB) error {
	// Attempt to update without retries
	currentUser, err := c.getUserDB(user.UserId)
	if err != nil {
		return fmt.Errorf("failed to get current user state: %w", err)
	}

	keyString, partitionKey := c.dbKey(user.UserId)
	encryptedDoc, err := c.createEncryptedDoc(user, keyString)
	if err != nil {
		return fmt.Errorf("failed to marshal updated document: %w", err)
	}

	options := &azcosmos.ItemOptions{
		IfMatchEtag: &currentUser.etag,
	}

	_, err = c.usersContainer.ReplaceItem(ctx, partitionKey, keyString, encryptedDoc, options)
	if err != nil {
		if strings.Contains(err.Error(), "Precondition Failed") {
			return fmt.Errorf("ETag mismatch: the user document was modified by another process")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (c *CosmosDB) createEncryptedDoc(user dbcommon.GWUserDB, keyString string) ([]byte, error) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	ciphertext, err := c.encryptor.Encrypt(userJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt user data: %w", err)
	}

	// Create an EncryptedDocument struct to store in CosmosDB
	doc := EncryptedDocument{
		ID:   keyString,
		Data: ciphertext,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}
	return docJSON, nil
}

func (c *CosmosDB) dbKey(userID []byte) (string, azcosmos.PartitionKey) {
	key := c.encryptor.HashWithHMAC(userID)
	keyString := hex.EncodeToString(key)
	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	return keyString, partitionKey
}

// GetEncryptionKey returns the encryption key used by the CosmosDB instance
func (c *CosmosDB) GetEncryptionKey() []byte {
	return c.encryptor.GetKey()
}
