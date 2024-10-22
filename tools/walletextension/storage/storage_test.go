package storage

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

var tests = map[string]func(storage Storage, t *testing.T){
	"testAddAndGetUser":     testAddAndGetUser,
	"testAddAndGetAccounts": testAddAndGetAccounts,
	"testDeleteUser":        testDeleteUser,
	"testGetUser":           testGetUser,
}

func TestGatewayStorage(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := New("sqlite", "", "")
			// storage, err := New("cosmosDB", "<cosmosdb-connection-string>", "")
			require.NoError(t, err)

			test(storage, t)
		})
	}
}

func testAddAndGetUser(storage Storage, t *testing.T) {
	// Generate random user ID and private key
	userID := make([]byte, 20)
	_, err := rand.Read(userID)
	if err != nil {
		t.Fatal(err)
	}
	privateKey := make([]byte, 32)
	_, err = rand.Read(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Add user to storage
	err = storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve user's private key from storage
	user, err := storage.GetUser(userID)
	if err != nil {
		t.Fatal(err)
	}

	// Check if retrieved private key matches the original
	if !bytes.Equal(user.PrivateKey, privateKey) {
		t.Errorf("privateKey mismatch: got %v, want %v", user.PrivateKey, privateKey)
	}
}

func testAddAndGetAccounts(storage Storage, t *testing.T) {
	// Generate random user ID, private key, and account details
	userID := make([]byte, 20)
	rand.Read(userID)
	privateKey := make([]byte, 32)
	rand.Read(privateKey)
	accountAddress1 := make([]byte, 20)
	rand.Read(accountAddress1)
	signature1 := make([]byte, 65)
	rand.Read(signature1)

	// Add a new user to the storage
	err := storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Add the first account for the user
	err = storage.AddAccount(userID, accountAddress1, signature1, viewingkey.EIP712Signature)
	if err != nil {
		t.Fatal(err)
	}

	// Generate details for a second account
	accountAddress2 := make([]byte, 20)
	rand.Read(accountAddress2)
	signature2 := make([]byte, 65)
	rand.Read(signature2)

	// Add the second account for the user
	err = storage.AddAccount(userID, accountAddress2, signature2, viewingkey.EIP712Signature)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve all accounts for the user
	accounts, err := storage.GetAccounts(userID)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the correct number of accounts were retrieved
	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

	// Flags to check if both accounts are found
	foundAccount1 := false
	foundAccount2 := false

	// Iterate through retrieved accounts and check if they match the added accounts
	for _, account := range accounts {
		if bytes.Equal(account.AccountAddress, accountAddress1) && bytes.Equal(account.Signature, signature1) {
			foundAccount1 = true
		}
		if bytes.Equal(account.AccountAddress, accountAddress2) && bytes.Equal(account.Signature, signature2) {
			foundAccount2 = true
		}
	}

	// Verify that both accounts were found
	if !foundAccount1 {
		t.Errorf("Account 1 was not found in the result")
	}

	if !foundAccount2 {
		t.Errorf("Account 2 was not found in the result")
	}
}

func testDeleteUser(storage Storage, t *testing.T) {
	// Generate random user ID and private key
	userID := make([]byte, 20)
	rand.Read(userID)
	privateKey := make([]byte, 32)
	rand.Read(privateKey)

	// Add user to storage
	err := storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the user
	err = storage.DeleteUser(userID)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to retrieve the deleted user's private key
	// This should fail with a "not found" error
	_, err = storage.GetUser(userID)
	if err == nil || !errors.Is(err, errutil.ErrNotFound) {
		t.Fatal("Expected 'not found' error when getting deleted user, but got none or different error")
	}
}

func testGetUser(storage Storage, t *testing.T) {
	// Generate random user ID and private key
	userID := make([]byte, 20)
	rand.Read(userID)
	privateKey := make([]byte, 32)
	rand.Read(privateKey)

	// Add user to storage
	err := storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}

	// Get user from storage
	user, err := storage.GetUser(userID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Check if retrieved user matches the added user
	if !bytes.Equal(user.UserId, userID) {
		t.Errorf("Retrieved user ID does not match. Expected %x, got %x", userID, user.UserId)
	}

	if !bytes.Equal(user.PrivateKey, privateKey) {
		t.Errorf("Retrieved private key does not match. Expected %x, got %x", privateKey, user.PrivateKey)
	}

	// Try to get a non-existent user
	nonExistentUserID := make([]byte, 20)
	rand.Read(nonExistentUserID)
	_, err = storage.GetUser(nonExistentUserID)
	if err == nil {
		t.Error("Expected error when getting non-existent user, but got none")
	}
}
