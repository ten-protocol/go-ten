package storage

import (
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"reflect"
	"testing"

	gethlog "github.com/ethereum/go-ethereum/log"
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

	vk1, _ := viewingkey.GenerateAndSignViewingKey(wallet1)
	vk2, _ := viewingkey.GenerateAndSignViewingKey(wallet2)

	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	err = myStorage.SaveUserVK(userID, vk1)
	if err != nil {
		t.Fatal(err)
	}
	err = myStorage.SaveUserVK(userID, vk2)
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

	vk1, _ := viewingkey.GenerateAndSignViewingKey(wallet1)
	vk2, _ := viewingkey.GenerateAndSignViewingKey(wallet2)

	myStorage, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	err = myStorage.SaveUserVK(userID1, vk1)
	if err != nil {
		t.Fatal(err)
	}

	err = myStorage.SaveUserVK(userID2, vk2)
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
