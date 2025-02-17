package manualtests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/tools/walletextension/lib"
)

func TestSubscribeToOG(t *testing.T) {
	t.Skip("skip manual tests")

	// Using http
	ogHTTPAddress := "https://dev-testnet.ten.xyz:443"
	ogWSAddress := "wss://dev-testnet.ten.xyz:81"
	// ogWSAddress := "ws://51.132.131.47:81"

	ogClient := lib.NewTenGatewayLibrary(ogHTTPAddress, ogWSAddress)

	// join the network
	err := ogClient.Join()
	require.NoError(t, err)
	fmt.Println(ogClient.UserID())

	// register an account
	err = ogClient.RegisterAccount(l2Wallet.PrivateKey(), l2Wallet.Address())
	require.NoError(t, err)
	fmt.Println("Registered account: ", l2Wallet.Address().Hex())

	// Using WS ->

	// Connect to WebSocket server using the standard geth client
	client, err := ethclient.Dial(ogWSAddress)
	require.NoError(t, err)

	// Create a simple request
	at, err := client.BalanceAt(context.Background(), l2Wallet.Address(), nil)
	require.NoError(t, err)

	fmt.Println("Balance for account ", l2Wallet.Address().Hex(), " - ", at.String())

	// Create a subscription
	query := ethereum.FilterQuery{
		Addresses: []gethcommon.Address{l2Wallet.Address()},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	// Listen for events from the contract
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			// Process the contract event
			// This is just a simple example printing the block number; you'll want to decode and handle the logs according to your contract's ABI
			log.Printf("Received log in block number: %v", vLog.BlockNumber)
		}
	}
}
