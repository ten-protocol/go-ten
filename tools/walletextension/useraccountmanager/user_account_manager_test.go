package useraccountmanager

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

func TestAddingAndGettingUserAccountManagers(t *testing.T) {
	unauthedClient, _ := rpc.NewNetworkClient("ws://test")
	userAccountManager := NewUserAccountManager(unauthedClient, log.New())
	userID1 := "user1"
	userID2 := "user2"

	// Test adding and getting account manager for userID1
	userAccountManager.AddUserAccountManager(userID1)
	accManager1, err := userAccountManager.GetUserAccountManager(userID1)
	if err != nil {
		t.Fatal(err)
	}
	// We should get error if we try to get Account manager for User2
	_, err = userAccountManager.GetUserAccountManager(userID2)
	fmt.Println(err)
	if err == nil {
		t.Fatal("expecting error when trying to get AccountManager for user that doesn't exist.")
	}

	// After trying to add new AccountManager for the same user we should get the same instance (not overriding old one)
	userAccountManager.AddUserAccountManager(userID1)
	accManager1New, err := userAccountManager.GetUserAccountManager(userID1)
	if err != nil {
		t.Fatal(err)
	}

	if accManager1 != accManager1New {
		t.Fatal("AccountManagers are not the same after adding new account manager for the same userID")
	}

	// We get a new instance of AccountManager when we add it for a new user
	userAccountManager.AddUserAccountManager(userID2)
	accManager2, err := userAccountManager.GetUserAccountManager(userID2)
	if err != nil {
		t.Fatal(err)
	}

	if accManager1 == accManager2 {
		t.Fatal("AccountManagers are the same for two different users")
	}
}
