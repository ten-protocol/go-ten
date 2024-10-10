package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	common "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// CosmosDB represents a CosmosDB storage implementation
type CosmosDB struct {
	client            *azcosmos.Client
	usersContainer    *azcosmos.ContainerClient
	accountsContainer *azcosmos.ContainerClient
}

// NewCosmosDB creates a new CosmosDB storage instance
func NewCosmosDB(connectionString string) (*CosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		fmt.Println("Failed to create CosmosDB client:", err)
		return nil, err
	}

	// Hardcoded database and container names
	databaseName := "test-gateway"
	usersContainerName := "users"
	accountsContainerName := "accounts"

	// Create database if it doesn't exist
	ctx := context.Background()
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: databaseName}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		fmt.Println("Failed to create database:", err)
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	usersContainer, err := client.NewContainer(databaseName, usersContainerName)
	if err != nil {
		fmt.Println("Failed to create users container:", err)
		return nil, fmt.Errorf("failed to create users container: %w", err)
	}

	accountsContainer, err := client.NewContainer(databaseName, accountsContainerName)
	if err != nil {
		fmt.Println("Failed to create accounts container:", err)
		return nil, fmt.Errorf("failed to create accounts container: %w", err)
	}

	// fmt.Println("CosmosDB connection string:", connectionString)
	// fmt.Println("CosmosDB client:", client)
	return &CosmosDB{
		client:            client,
		usersContainer:    usersContainer,
		accountsContainer: accountsContainer,
	}, nil
}

// Add more methods as required by your storage interface

func (c *CosmosDB) AddUser(userID []byte, privateKey []byte) error {
	ctx := context.Background()

	user := struct {
		ID         string `json:"id"`
		UserID     string `json:"userId"`
		PrivateKey string `json:"privateKey"`
	}{
		ID:         hex.EncodeToString(userID),
		UserID:     hex.EncodeToString(userID),
		PrivateKey: hex.EncodeToString(privateKey),
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	userIDString := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDString)
	_, err = c.usersContainer.CreateItem(ctx, partitionKey, userJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	fmt.Println("User added successfully")

	return nil
}

func (c *CosmosDB) DeleteUser(userID []byte) error {
	// Implementation needed
	return nil
}

func (c *CosmosDB) GetUserPrivateKey(userID []byte) ([]byte, error) {
	// Implementation needed
	return nil, nil
}

func (c *CosmosDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	// Implementation needed
	return nil
}

func (c *CosmosDB) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	// Implementation needed
	return nil, nil
}

func (c *CosmosDB) GetAllUsers() ([]common.UserDB, error) {
	// Implementation needed
	return nil, nil
}

func (c *CosmosDB) StoreTransaction(rawTx string, userID []byte) error {
	// Implementation needed
	return nil
}
