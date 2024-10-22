package cosmosdb

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

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

type CosmosDB struct {
	client         *azcosmos.Client
	usersContainer *azcosmos.ContainerClient
	encryptor      *encryption.Encryptor
}

type EncryptedDocument struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

const (
	DATABASE_NAME        = "gatewayDB"
	USERS_CONTAINER_NAME = "users"
	PARTITION_KEY        = "/id"
)

func NewCosmosDB(connectionString string, encryptionKey []byte) (*CosmosDB, error) {
	// create encryptor
	encryptor, err := encryption.NewEncryptor(encryptionKey)
	if err != nil {
		return nil, err
	}

	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
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
		encryptor:      encryptor,
	}, nil
}

func (c *CosmosDB) AddUser(userID []byte, privateKey []byte) error {
	user := common.GWUserDB{
		ID:         hex.EncodeToString(userID),
		UserId:     userID,
		PrivateKey: privateKey,
		Accounts:   []common.GWAccountDB{},
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	var encryptedData string
	if c.encryptor != nil {
		ciphertext, err := c.encryptor.Encrypt(userJSON)
		if err != nil {
			return err
		}
		encryptedData = base64.StdEncoding.EncodeToString(ciphertext)
	} else {
		encryptedData = base64.StdEncoding.EncodeToString(userJSON)
	}

	// Hash the userID to use as the key
	key := userID
	if c.encryptor != nil {
		key = c.encryptor.HashWithHMAC(userID)
	}
	keyString := hex.EncodeToString(key)

	// Create the encrypted document
	doc := EncryptedDocument{
		ID:   keyString,
		Data: encryptedData,
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
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
	key := userID
	if c.encryptor != nil {
		key = c.encryptor.HashWithHMAC(userID)
	}
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
	key := userID
	if c.encryptor != nil {
		key = c.encryptor.HashWithHMAC(userID)
	}
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

	encryptedData, err := base64.StdEncoding.DecodeString(doc.Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64 data: %w", err)
	}

	data := encryptedData
	if c.encryptor != nil {
		data, err = c.encryptor.Decrypt(encryptedData)
		if err != nil {
			return fmt.Errorf("failed to decrypt data: %w", err)
		}
	}

	var user common.GWUserDB
	err = json.Unmarshal(data, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Add the new account
	newAccount := common.GWAccountDB{
		AccountAddress: accountAddress,
		Signature:      signature,
		SignatureType:  int(signatureType),
	}
	user.Accounts = append(user.Accounts, newAccount)

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling updated user: %w", err)
	}

	var encryptedDataStr string
	if c.encryptor != nil {
		ciphertext, err := c.encryptor.Encrypt(userJSON)
		if err != nil {
			return fmt.Errorf("failed to encrypt updated user data: %w", err)
		}
		encryptedDataStr = base64.StdEncoding.EncodeToString(ciphertext)
	} else {
		encryptedDataStr = base64.StdEncoding.EncodeToString(userJSON)
	}

	// Update the document
	doc.Data = encryptedDataStr

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

func (c *CosmosDB) GetUser(userID []byte) (common.GWUserDB, error) {
	key := userID
	if c.encryptor != nil {
		key = c.encryptor.HashWithHMAC(userID)
	}
	keyString := hex.EncodeToString(key)
	partitionKey := azcosmos.NewPartitionKeyString(keyString)
	ctx := context.Background()

	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, keyString, nil)
	if err != nil {
		return common.GWUserDB{}, errutil.ErrNotFound
	}

	var doc EncryptedDocument
	err = json.Unmarshal(itemResponse.Value, &doc)
	if err != nil {
		return common.GWUserDB{}, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	encryptedData, err := base64.StdEncoding.DecodeString(doc.Data)
	if err != nil {
		return common.GWUserDB{}, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	data := encryptedData
	if c.encryptor != nil {
		data, err = c.encryptor.Decrypt(encryptedData)
		if err != nil {
			return common.GWUserDB{}, fmt.Errorf("failed to decrypt data: %w", err)
		}
	}

	var user common.GWUserDB
	err = json.Unmarshal(data, &user)
	if err != nil {
		return common.GWUserDB{}, fmt.Errorf("failed to unmarshal user data: %w", err)
	}
	return user, nil
}
