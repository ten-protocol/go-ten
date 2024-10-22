package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
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
}

const (
	DATABASE_NAME        = "gatewayDB"
	USERS_CONTAINER_NAME = "users"
	PARTITION_KEY        = "/id"
)

// GWUser is the user struct for the gateway
// both ID and UserID are the same for now, but we will use different values with encryption
type GWUser struct {
	ID         string      `json:"id"` // Required by CosmosDB
	UserId     []byte      `json:"userId"`
	PrivateKey []byte      `json:"privateKey"`
	Accounts   []GWAccount `json:"accounts"` // List of Accounts
}

type GWAccount struct {
	AccountAddress []byte `json:"accountAddress"`
	Signature      []byte `json:"signature"`
	SignatureType  int    `json:"signatureType"`
}

func NewCosmosDB(connectionString string) (*CosmosDB, error) {
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
	}, nil
}

func (c *CosmosDB) AddUser(userID []byte, privateKey []byte) error {
	user := GWUser{
		ID:         hex.EncodeToString(userID),
		UserId:     userID,
		PrivateKey: privateKey,
		Accounts:   []GWAccount{},
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// add to cosmosdb
	partitionKey := azcosmos.NewPartitionKeyString(user.ID)

	ctx := context.Background()
	_, err = c.usersContainer.CreateItem(ctx, partitionKey, userJSON, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *CosmosDB) DeleteUser(userID []byte) error {
	// Convert userID to hex string for use as partition key
	userIDHex := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	ctx := context.Background()

	// Delete the item from the container
	_, err := c.usersContainer.DeleteItem(ctx, partitionKey, userIDHex, nil)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *CosmosDB) GetUserPrivateKey(userID []byte) ([]byte, error) {
	// Convert userID to hex string for use as partition key
	userIDHex := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	ctx := context.Background()

	// Read the item from the container
	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, userIDHex, nil)
	if err != nil {
		return nil, errutil.ErrNotFound
	}

	var user GWUser
	err = json.Unmarshal(itemResponse.Value, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user.PrivateKey, nil
}

func (c *CosmosDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	// Convert userID to hex string for use as partition key
	userIDHex := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	ctx := context.Background()

	// Read the existing user
	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, userIDHex, nil)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	var user GWUser
	err = json.Unmarshal(itemResponse.Value, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Create new account
	newAccount := GWAccount{
		AccountAddress: accountAddress,
		Signature:      signature,
		SignatureType:  int(signatureType),
	}

	// Add new account to user's accounts
	user.Accounts = append(user.Accounts, newAccount)

	// Marshal updated user back to JSON
	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling updated user: %w", err)
	}

	// Update the item in the container
	_, err = c.usersContainer.ReplaceItem(ctx, partitionKey, userIDHex, updatedUserJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to update user with new account: %w", err)
	}
	return nil
}

func (c *CosmosDB) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	// Convert userID to hex string for use as partition key
	userIDHex := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	ctx := context.Background()

	// Read the existing user
	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, userIDHex, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user GWUser
	err = json.Unmarshal(itemResponse.Value, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// Convert GWAccount to common.AccountDB
	accounts := make([]common.AccountDB, len(user.Accounts))
	for i, acc := range user.Accounts {
		accounts[i] = common.AccountDB{
			AccountAddress: acc.AccountAddress,
			Signature:      acc.Signature,
			SignatureType:  acc.SignatureType,
		}
	}

	return accounts, nil
}

func (c *CosmosDB) GetUser(userID []byte) (common.UserDB, error) {
	// Convert userID to hex string for use as partition key
	userIDHex := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	ctx := context.Background()

	// Read the existing user
	itemResponse, err := c.usersContainer.ReadItem(ctx, partitionKey, userIDHex, nil)
	if err != nil {
		return common.UserDB{}, fmt.Errorf("failed to get user: %w", err)
	}

	var user GWUser
	err = json.Unmarshal(itemResponse.Value, &user)
	if err != nil {
		return common.UserDB{}, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	// TODO: @ziga - check if I can use user directly instead of GWUser/UserDB since they are mostly the same...
	// Convert GWUser to common.UserDB
	return common.UserDB{
		UserID:     userID,
		PrivateKey: user.PrivateKey,
	}, nil
}
