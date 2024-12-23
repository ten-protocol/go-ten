package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	dbcommon "github.com/ten-protocol/go-ten/tools/walletextension/storage/database/common"

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

// maxRetries defines how many times we'll retry in case of ETag conflicts
const maxRetries = 5

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

// updateUser updates the user's data in Cosmos DB using optimistic concurrency (ETag). If another
// process updates the document concurrently, the ETag check fails and this function retries until
// it either succeeds or the maximum number of retries is reached.
func (c *CosmosDB) updateUser(ctx context.Context, newUser dbcommon.GWUserDB) error {
	// Retry loop for handling concurrency conflicts
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Retrieves the latest version of the user from the database, including its current ETag
		currentUser, err := c.getUserDB(newUser.UserId)
		if err != nil {
			return fmt.Errorf("failed to get current user state: %w", err)
		}

		// Creates an encrypted document using the new user data
		keyString, partitionKey := c.dbKey(newUser.UserId)
		encryptedDoc, err := c.createEncryptedDoc(newUser, keyString)
		if err != nil {
			return fmt.Errorf("failed to marshal updated document: %w", err)
		}

		// Attempts to replace the document in Cosmos DB, enforcing the ETag check to detect concurrent modifications
		options := &azcosmos.ItemOptions{
			IfMatchEtag: &currentUser.etag,
		}
		_, err = c.usersContainer.ReplaceItem(ctx, partitionKey, keyString, encryptedDoc, options)
		if err != nil {
			// If a concurrent update has changed the ETag, retry by reading the updated document again
			if strings.Contains(err.Error(), "Precondition Failed") {
				continue
			}
			// Any other error is returned directly
			return fmt.Errorf("failed to update user: %w", err)
		}

		// Update succeeded without any concurrency conflict
		return nil
	}

	// If all attempts fail due to repeated concurrency conflicts, return an error
	return fmt.Errorf("updateUser failed due to concurrency conflicts after %d retries", maxRetries)
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
