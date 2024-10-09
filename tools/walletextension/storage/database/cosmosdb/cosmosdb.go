package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	common "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const TEMP_PARTITION_KEY = "p"

// CosmosDB represents a CosmosDB storage implementation
type CosmosDB struct {
	client                *azcosmos.Client
	usersContainer        *azcosmos.ContainerClient
	accountsContainer     *azcosmos.ContainerClient
	transactionsContainer *azcosmos.ContainerClient
}

type Account struct {
	ID             string `json:"id"`     // Required by CosmosDB
	UserId         string `json:"userId"` // Partition Key
	AccountAddress string `json:"accountAddress"`
	Signature      string `json:"signature"`
	SignatureType  int    `json:"signatureType"`
}

type User struct {
	ID           string `json:"id"` // Required by CosmosDB
	UserId       string `json:"userId"`
	PartitionKey string `json:"partitionKey"`
	PrivateKey   string `json:"privateKey"`
}

type Transaction struct {
	ID           string `json:"id"` // Required by CosmosDB
	UserId       string `json:"userId"`
	PartitionKey string `json:"partitionKey"`
	RawTx        string `json:"rawTx"`
}

// NewCosmosDB creates a new CosmosDB storage instance
func NewCosmosDB(connectionString, dbPath string) (*CosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		fmt.Println("Failed to create CosmosDB client:", err)
		return nil, err
	}

	// Hardcoded database and container names
	databaseName := dbPath
	usersContainerName := "users"
	accountsContainerName := "accounts"
	transactionsContainerName := "transactions"
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

	transactionsContainer, err := client.NewContainer(databaseName, transactionsContainerName)
	if err != nil {
		fmt.Println("Failed to create transactions container:", err)
		return nil, fmt.Errorf("failed to create transactions container: %w", err)
	}

	return &CosmosDB{
		client:                client,
		usersContainer:        usersContainer,
		accountsContainer:     accountsContainer,
		transactionsContainer: transactionsContainer,
	}, nil
}

// Add more methods as required by your storage interface

func (c *CosmosDB) AddUser(userID []byte, privateKey []byte) error {
	ctx := context.Background()

	user := User{
		ID:           hex.EncodeToString(userID),
		UserId:       hex.EncodeToString(userID),
		PartitionKey: TEMP_PARTITION_KEY,
		PrivateKey:   hex.EncodeToString(privateKey),
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	partitionKey := azcosmos.NewPartitionKeyString(user.PartitionKey)
	_, err = c.usersContainer.CreateItem(ctx, partitionKey, userJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (c *CosmosDB) DeleteUser(userID []byte) error {
	ctx := context.Background()

	// Convert userID to hex string to match the stored ID
	id := hex.EncodeToString(userID)

	// Create the partition key using the userID
	partitionKey := azcosmos.NewPartitionKeyString("p")

	// Attempt to delete the user item from the users container
	_, err := c.usersContainer.DeleteItem(ctx, partitionKey, id, nil)
	if err != nil {
		// TODO check if user not found or call failed for other reasons
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (c *CosmosDB) GetUserPrivateKey(userID []byte) ([]byte, error) {
	ctx := context.Background()

	// Convert userID to hex string to match the stored ID
	id := hex.EncodeToString(userID)

	// Create the partition key using the userID
	partitionKey := azcosmos.NewPartitionKeyString("p")

	// Create the query
	query := fmt.Sprintf("SELECT c.privateKey FROM c WHERE c.userId = '%s'", id)

	// Create a query pager
	queryPager := c.usersContainer.NewQueryItemsPager(query, partitionKey, nil)

	var privateKey string

	// Iterate through the pages
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to query user: %w", err)
		}

		// We expect only one result, so we'll break after processing the first item
		for _, item := range queryResponse.Items {
			var result struct {
				PrivateKey string `json:"privateKey"`
			}
			if err := json.Unmarshal(item, &result); err != nil {
				return nil, fmt.Errorf("failed to unmarshal query result: %w", err)
			}
			privateKey = result.PrivateKey
			break
		}
	}

	if privateKey == "" {
		return nil, errutil.ErrNotFound
	}

	// Decode the private key from hex string back to bytes
	decodedPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	return decodedPrivateKey, nil
}

func (c *CosmosDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	ctx := context.Background()

	// Convert byte slices to hex strings
	userIDHex := hex.EncodeToString(userID)
	accountAddressHex := hex.EncodeToString(accountAddress)
	signatureHex := hex.EncodeToString(signature)
	signatureTypeInt := int(signatureType) // Assuming 'SignatureType' is based on int

	// Create the Account document
	account := Account{
		ID:             uuid.New().String(), // Generate a new UUID as a string
		UserId:         userIDHex,
		AccountAddress: accountAddressHex,
		Signature:      signatureHex,
		SignatureType:  signatureTypeInt,
	}

	// Serialize the account struct to JSON
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %w", err)
	}

	// Use 'userId' as the partition key
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	// Create the item in the accounts container
	_, err = c.accountsContainer.CreateItem(ctx, partitionKey, accountJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to add account: %w", err)
	}
	return nil
}

func (c *CosmosDB) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	ctx := context.Background()
	userIDHex := hex.EncodeToString(userID)
	partitionKey := azcosmos.NewPartitionKeyString(userIDHex)

	// Build the SQL query
	query := fmt.Sprintf("SELECT * FROM c WHERE c.userId = '%s'", userIDHex)

	// Now create the pager using NewQueryItemsPager
	queryPager := c.accountsContainer.NewQueryItemsPager(query, partitionKey, nil)

	var accounts []common.AccountDB

	// Iterate through the pages
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to query user: %w", err)
		}

		// We expect only one result, so we'll break after processing the first item
		for _, item := range queryResponse.Items {
			var result Account
			if err := json.Unmarshal(item, &result); err != nil {
				return nil, fmt.Errorf("failed to unmarshal query result: %w", err)
			}

			accountAddressBytes, err := hex.DecodeString(result.AccountAddress)
			if err != nil {
				return nil, fmt.Errorf("failed to decode accountAddress: %w", err)
			}

			signatureBytes, err := hex.DecodeString(result.Signature)
			if err != nil {
				return nil, fmt.Errorf("failed to decode signature: %w", err)
			}

			account := common.AccountDB{
				AccountAddress: accountAddressBytes,
				Signature:      signatureBytes,
				SignatureType:  result.SignatureType,
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

func (c *CosmosDB) GetAllUsers() ([]common.UserDB, error) {
	fmt.Println("GetAllUsers")
	ctx := context.Background()

	// Build the SQL query to select all users
	query := "SELECT * FROM c"

	partitionKey := azcosmos.NewPartitionKeyString("p")

	// Create the pager using NewQueryItemsPager
	queryPager := c.usersContainer.NewQueryItemsPager(query, partitionKey, nil)

	var users []common.UserDB

	// Iterate through the pages
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to query users: %w", err)
		}

		for _, item := range queryResponse.Items {
			var result struct {
				ID           string `json:"id"`
				PartitionKey string `json:"partitionKey"`
				UserID       string `json:"userId"`
				PrivateKey   string `json:"privateKey"`
			}
			if err := json.Unmarshal(item, &result); err != nil {
				return nil, fmt.Errorf("failed to unmarshal query result: %w", err)
			}

			userIDBytes, err := hex.DecodeString(result.UserID)
			if err != nil {
				return nil, fmt.Errorf("failed to decode userID: %w", err)
			}

			privateKeyBytes, err := hex.DecodeString(result.PrivateKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decode privateKey: %w", err)
			}

			user := common.UserDB{
				UserID:     userIDBytes,
				PrivateKey: privateKeyBytes,
			}

			users = append(users, user)
		}
	}

	return users, nil
}

func (c *CosmosDB) StoreTransaction(rawTx string, userID []byte) error {
	ctx := context.Background()

	transaction := Transaction{
		ID:           uuid.New().String(),
		UserId:       hex.EncodeToString(userID),
		PartitionKey: TEMP_PARTITION_KEY,
		RawTx:        rawTx,
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	partitionKey := azcosmos.NewPartitionKeyString(transaction.PartitionKey)

	_, err = c.transactionsContainer.CreateItem(ctx, partitionKey, transactionJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to store transaction: %w", err)
	}

	return nil
}
