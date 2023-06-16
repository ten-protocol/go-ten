package storage

import (
	"bytes"
	"testing"
)

func TestAddAndGetUser(t *testing.T) {
	storage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	userID := []byte("userID")
	privateKey := []byte("privateKey")

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

func TestAddAndGetAccounts(t *testing.T) {
	storage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	userID := []byte("userID")
	privateKey := []byte("privateKey")
	accountAddress1 := []byte("accountAddress1")
	signature1 := []byte("signature1")

	err = storage.AddUser(userID, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.AddAccount(userID, accountAddress1, signature1)
	if err != nil {
		t.Fatal(err)
	}

	accountAddress2 := []byte("accountAddress2")
	signature2 := []byte("signature2")

	err = storage.AddAccount(userID, accountAddress2, signature2)
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
