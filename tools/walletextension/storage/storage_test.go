package storage

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"testing"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

var tests = map[string]func(storage Storage, t *testing.T){
	"testAddAndGetUser":     testAddAndGetUser,
	"testAddAndGetAccounts": testAddAndGetAccounts,
	"testDeleteUser":        testDeleteUser,
	"testGetAllUsers":       testGetAllUsers,
	"testStoringNewTx":      testStoringNewTx,
}

func TestGatewayStorage(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// storage, err := New("mariaDB", "obscurouser:password@tcp(127.0.0.1:3306)/ogdb", "") allows to run tests against a local instance of MariaDB
			//storage, err := New("sqlite", "", "")
			storage, err := New("cosmosDB", "", "dev-testnet-gateway")
			require.NoError(t, err)

			test(storage, t)
		})
	}
}

func testAddAndGetUser(storage Storage, t *testing.T) {
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

	err = storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	returnedPrivateKey, err := storage.GetUserPrivateKey(userID)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(returnedPrivateKey, privateKey) {
		t.Errorf("privateKey mismatch: got %v, want %v", returnedPrivateKey, privateKey)
	}
}

func testAddAndGetAccounts(storage Storage, t *testing.T) {
	userID := make([]byte, 20)
	rand.Read(userID)
	privateKey := make([]byte, 32)
	rand.Read(privateKey)
	accountAddress1 := make([]byte, 20)
	rand.Read(accountAddress1)
	signature1 := make([]byte, 65)
	rand.Read(signature1)

	err := storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.AddAccount(userID, accountAddress1, signature1, viewingkey.EIP712Signature)
	if err != nil {
		t.Fatal(err)
	}

	accountAddress2 := make([]byte, 20)
	rand.Read(accountAddress2)
	signature2 := make([]byte, 65)
	rand.Read(signature2)

	err = storage.AddAccount(userID, accountAddress2, signature2, viewingkey.EIP712Signature)
	if err != nil {
		t.Fatal(err)
	}

	accounts, err := storage.GetAccounts(userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

	foundAccount1 := false
	foundAccount2 := false

	for _, account := range accounts {
		if bytes.Equal(account.AccountAddress, accountAddress1) && bytes.Equal(account.Signature, signature1) {
			foundAccount1 = true
		}
		if bytes.Equal(account.AccountAddress, accountAddress2) && bytes.Equal(account.Signature, signature2) {
			foundAccount2 = true
		}
	}

	if !foundAccount1 {
		t.Errorf("Account 1 was not found in the result")
	}

	if !foundAccount2 {
		t.Errorf("Account 2 was not found in the result")
	}
}

func testDeleteUser(storage Storage, t *testing.T) {
	userID := make([]byte, 20)
	rand.Read(userID)
	privateKey := make([]byte, 32)
	rand.Read(privateKey)

	err := storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.DeleteUser(userID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = storage.GetUserPrivateKey(userID)
	fmt.Println("err:", err)
	if err == nil || !errors.Is(err, errutil.ErrNotFound) {
		t.Fatal("Expected error when getting deleted user, but got none")
	}
}

func testGetAllUsers(storage Storage, t *testing.T) {
	initialUsers, err := storage.GetAllUsers()
	if err != nil {
		t.Fatal(err)
	}

	userID := make([]byte, 20)
	rand.Read(userID)
	privateKey := make([]byte, 32)
	rand.Read(privateKey)

	err = storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	afterInsertUsers, err := storage.GetAllUsers()
	if err != nil {
		t.Fatal(err)
	}

	if len(afterInsertUsers) != len(initialUsers)+1 {
		t.Errorf("Expected user count to increase by 1. Got %d initially and %d after insert", len(initialUsers), len(afterInsertUsers))
	}
}

func testStoringNewTx(storage Storage, t *testing.T) {
	userID := []byte("userID")
	rawTransaction := "0x0123456789"

	err := storage.StoreTransaction(rawTransaction, userID)
	if err != nil {
		t.Fatal(err)
	}
}
