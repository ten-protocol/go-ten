package manualtests

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSubscribeToOG(t *testing.T) {
	//if os.Getenv(_IDEFlag) == "" {
	//	t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	//}

	// Using http
	//ogHTTPAddress := "dev-testnet.obscu.ro:443"
	ogWSAddress := "dev-testnet.obscu.ro:81"
	//ogWSAddress := "51.132.131.47:81"

	//// join the network
	//statusCode, userID, err := fasthttp.Get(nil, fmt.Sprintf("https://%s/v1/join/", ogHTTPAddress))
	//require.NoError(t, err) // dialing to the given TCP address timed out
	//fmt.Println(statusCode)
	//fmt.Println(userID)
	//
	//// sign the message
	//messagePayload := signMessage(string(userID))
	//
	//// register an account
	//var regAccountResp []byte
	//regAccountResp, err = registerAccount(ogHTTPAddress, string(userID), messagePayload)
	//require.NoError(t, err)
	//fmt.Println(regAccountResp)

	// Using WS ->

	for i := 0; i < 50; i++ {
		ethclient.Dial("ws://" + ogWSAddress)
		time.Sleep(100 * time.Millisecond)
	}

	// Connect to WebSocket server using the standard geth client
	client, err := ethclient.Dial("ws://" + ogWSAddress)
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

func registerAccount(baseAddress, userID, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"https://"+baseAddress+"/authenticate/?u="+userID,
		strings.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

//	{
//	 "signature": "0xc784adea83ed3ec60528f4747418c85abe553b35a47fd2c95425de654bb9d0d40ede24aec182e6a2ec65c0c7c6aedab7823f21a9b9f7ff5db3a77a9f90dc97b41c",
//	 "message": "Register e097c4a10d4285d13b377985834b4c57e069b5856cc6c2cd4a038f62da4bc459 for 0x06ed49a32fcc5094abee51a4ffd46dd23b62a191"
//	}
func signMessage(userID string) string {
	pk := l2Wallet.PrivateKey()
	address := l2Wallet.Address()
	hexAddress := address.Hex()

	message := fmt.Sprintf("Register %s for %s", userID, strings.ToLower(hexAddress))
	prefixedMessage := fmt.Sprintf(common.PersonalSignMessagePrefix, len(message), message)

	messageHash := crypto.Keccak256([]byte(prefixedMessage))
	sig, err := crypto.Sign(messageHash, pk)
	if err != nil {
		log.Fatalf("Failed to sign message: %v", err)
	}
	return fmt.Sprintf("{\"signature\": \"%s\", \"message\": \"%s\"}", hex.EncodeToString(sig), message)
}
