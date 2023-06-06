package storage

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

func TestStoringMultipleKeysPerUser(t *testing.T) {
	userID := "abc"
	wallet1 := wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	wallet2 := wallet.NewInMemoryWalletFromConfig(
		"111114725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b21111",
		777,
		gethlog.New())

	vk1, _ := rpc.GenerateAndSignViewingKey(wallet1)
	vk2, _ := rpc.GenerateAndSignViewingKey(wallet2)

	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	err = myStorage.SaveUserVK(userID, vk1, "")
	if err != nil {
		t.Fatal(err)
	}
	err = myStorage.SaveUserVK(userID, vk2, "")
	if err != nil {
		t.Fatal(err)
	}
	result, err := myStorage.GetUserVKs(userID)
	if err != nil {
		t.Errorf("error getting user VKs")
	}

	// Check if vk1 is in the result
	foundVK1 := false
	for _, vk := range result {
		if reflect.DeepEqual(vk, vk1) {
			foundVK1 = true
			break
		}
	}
	if !foundVK1 {
		t.Errorf("vk1 is not found in the result")
	}

	// Check if vk2 is in the result
	foundVK2 := false
	for _, vk := range result {
		if reflect.DeepEqual(vk, vk2) {
			foundVK2 = true
			break
		}
	}
	if !foundVK2 {
		t.Errorf("vk2 is not found in the result")
	}
}

func TestMultipleUsersStoringKeys(t *testing.T) {
	userID1 := "user1"
	wallet1 := wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())

	userID2 := "user2"
	wallet2 := wallet.NewInMemoryWalletFromConfig(
		"111114725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b21111",
		777,
		gethlog.New())

	vk1, _ := rpc.GenerateAndSignViewingKey(wallet1)
	vk2, _ := rpc.GenerateAndSignViewingKey(wallet2)

	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	err = myStorage.SaveUserVK(userID1, vk1, "")
	if err != nil {
		t.Fatal(err)
	}

	err = myStorage.SaveUserVK(userID2, vk2, "")
	if err != nil {
		t.Fatal(err)
	}

	// userId1 should get only the item that belongs to him and not an item that belongs to userId2
	result, err := myStorage.GetUserVKs(userID1)
	if err != nil {
		t.Errorf("error getting user VKs")
	}

	// Check if vk1 is in the result
	foundVK1 := false
	for _, vk := range result {
		if reflect.DeepEqual(vk, vk1) {
			foundVK1 = true
			break
		}
	}
	if !foundVK1 {
		t.Errorf("vk1 is not found in the result and it should be in")
	}

	// Check if vk2 is not in the result
	foundVK2 := false
	for _, vk := range result {
		if reflect.DeepEqual(vk, vk2) {
			foundVK2 = true
			break
		}
	}
	if foundVK2 {
		t.Errorf("vk2 is not found in the result for the wrong user")
	}

	// userId2 should get only the item that belongs to him and not an item that belongs to userId1
	result, err = myStorage.GetUserVKs(userID2)
	if err != nil {
		t.Errorf("error getting user VKs")
	}

	// Check if vk1 is not in the result
	foundVK1 = false
	for _, vk := range result {
		if reflect.DeepEqual(vk, vk1) {
			foundVK1 = true
			break
		}
	}
	if foundVK1 {
		t.Errorf("vk1 is not found in the result for the wrong user")
	}

	// Check if vk2 is not in the result
	foundVK2 = false
	for _, vk := range result {
		if reflect.DeepEqual(vk, vk2) {
			foundVK2 = true
			break
		}
	}
	if !foundVK2 {
		t.Errorf("vk2 is not found in the result for the wrong user")
	}
}

func TestAddAndGetMessage(t *testing.T) {
	userID1 := "user1"
	wallet1 := wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	vk1, _ := rpc.GenerateAndSignViewingKey(wallet1)
	userAccount := vk1.Account.Bytes()
	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	originalMessage := "mymessage"
	err = myStorage.SaveUserVK(userID1, vk1, originalMessage)
	if err != nil {
		t.Fatal(err)
	}

	// Test initial retrieval of message and signature
	message, signature, err := myStorage.GetMessageAndSignature(userID1, userAccount)
	if err != nil {
		t.Fatalf("Error getting user VKs: %v", err)
	}

	if message != originalMessage {
		t.Errorf("Original message and retrieved message are not the same")
	}

	if signature != "" {
		t.Errorf("No signature was added, but the retrieved signature is not empty")
	}

	originalSignature := "mysignature"
	err = myStorage.AddSignature(userID1, userAccount, originalSignature)
	if err != nil {
		t.Errorf("Error adding signature: %v", err)
	}

	// Test retrieval after adding the signature
	message, signature, err = myStorage.GetMessageAndSignature(userID1, userAccount)
	if err != nil {
		t.Errorf("Error getting message and signature: %v", err)
	}

	if message != originalMessage {
		t.Errorf("Original message and retrieved message are not the same")
	}

	if signature != originalSignature {
		t.Errorf("Signature from the database and original signature are not the same")
	}
}

func TestGetUnauthenticatedUserPrivateKey(t *testing.T) {
	// create new storage
	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	// generate new keys
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	vk := &rpc.ViewingKey{
		Account:    nil,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  nil,
		SignedKey:  nil,
	}

	// store viewing key and calculate userID
	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	userID := crypto.Keccak256Hash(viewingPublicKeyBytes)
	err = myStorage.SaveUserVK(userID.Hex(), vk, "")
	if err != nil {
		t.Errorf("Unable to store UserVK")
	}

	// get private key from the database
	storedPrivateKeyBytes, err := myStorage.GetUnauthenticatedUserPrivateKey(userID.Hex())
	if err != nil {
		t.Errorf("Unable to get private key for unauthenticated user: %s", userID.Hex())
	}

	// check if we have the same private key
	if !bytes.Equal(crypto.FromECDSA(viewingKeyPrivate), storedPrivateKeyBytes) {
		t.Errorf("Stored viewing key is not the same as original")
	}
}

func TestStoreAuthenticatedDataForUserAndGetViewingKey(t *testing.T) {
	const userID = "user"
	wallet1 := wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	vk, _ := rpc.GenerateAndSignViewingKey(wallet1)
	accountAddress := vk.Account.Bytes()
	vk.Account = nil
	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	// save viewing key (private key) - simulate call to /join endpoint
	err = myStorage.SaveUserVK(userID, vk, "")
	if err != nil {
		t.Fatal(err)
	}

	// store data for that user
	myMessage := "mymessage"
	mySignature := "mysignature"
	privateKeyBytes := crypto.FromECDSA(vk.PrivateKey.ExportECDSA())
	err = myStorage.StoreAuthenticatedDataForUser(userID, privateKeyBytes, accountAddress, myMessage, mySignature)
	if err != nil {
		t.Fatal(err)
	}

	// get the data back from the database
	storedPrivKey, storedMessage, storedSignedMessage, err := myStorage.GetDataForUserAndAddress(userID, accountAddress)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(storedPrivKey, privateKeyBytes) {
		t.Fatal("stored private key not the same as original")
	}

	if storedMessage != myMessage {
		t.Fatal("stored message not the same as original")
	}

	if storedSignedMessage != mySignature {
		t.Fatal("store signedMessage not the same as original")
	}
}
