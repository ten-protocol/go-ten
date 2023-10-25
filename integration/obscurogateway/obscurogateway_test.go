package faucet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/retry"

	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/httputil"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	integrationCommon "github.com/obscuronet/go-obscuro/integration/common"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/tools/walletextension/config"
	"github.com/obscuronet/go-obscuro/tools/walletextension/container"
	"github.com/obscuronet/go-obscuro/tools/walletextension/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "obscurogateway",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/obscurogateway/"
)

func TestObscuroGateway(t *testing.T) {
	startPort := integration.StartPortObscuroGatewayUnitTest
	createObscuroNetwork(t, startPort, 1)

	obscuroGatewayConf := config.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: startPort + integration.DefaultObscuroGatewayHTTPPortOffset,
		WalletExtensionPortWS:   startPort + integration.DefaultObscuroGatewayWSPortOffset,
		NodeRPCHTTPAddress:      fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		NodeRPCWebsocketAddress: fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		LogPath:                 "sys_out",
		VerboseFlag:             false,
		DBType:                  "sqlite",
	}

	obscuroGwContainer := container.NewWalletExtensionContainerFromConfig(obscuroGatewayConf, testlog.Logger())
	go func() {
		err := obscuroGwContainer.Start()
		if err != nil {
			fmt.Printf("error stopping WE - %s", err)
		}
	}()

	// wait for the msg bus contract to be deployed
	time.Sleep(5 * time.Second)

	// make sure the server is ready to receive requests
	httpURL := fmt.Sprintf("http://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortHTTP)
	wsURL := fmt.Sprintf("ws://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortWS)

	// make sure the server is ready to receive requests
	err := waitServerIsReady(httpURL)
	require.NoError(t, err)

	// run the tests against the exis
	for name, test := range map[string]func(*testing.T, string, string){
		//"testAreTxsMinted":            testAreTxsMinted, this breaks the other tests bc, enable once concurency issues are fixed
		"testErrorHandling":           testErrorHandling,
		"testErrorsRevertedArePassed": testErrorsRevertedArePassed,
	} {
		t.Run(name, func(t *testing.T) {
			test(t, httpURL, wsURL)
		})
	}

	// Gracefully shutdown
	err = obscuroGwContainer.Stop()
	assert.NoError(t, err)
}

func TestObscuroGatewaySubscriptionsWithMultipleAccounts(t *testing.T) {
	// t.Skip("Commented it out until more testing is driven from this test")
	startPort := integration.StartPortObscuroGatewayUnitTest
	wallets := createObscuroNetwork(t, startPort, 5)

	obscuroGatewayConf := config.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: startPort + integration.DefaultObscuroGatewayHTTPPortOffset,
		WalletExtensionPortWS:   startPort + integration.DefaultObscuroGatewayWSPortOffset,
		NodeRPCHTTPAddress:      fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		NodeRPCWebsocketAddress: fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		LogPath:                 "sys_out",
		VerboseFlag:             false,
		DBType:                  "sqlite",
	}

	obscuroGwContainer := container.NewWalletExtensionContainerFromConfig(obscuroGatewayConf, testlog.Logger())
	go func() {
		err := obscuroGwContainer.Start()
		if err != nil {
			fmt.Printf("error stopping WE - %s", err)
		}
	}()

	// wait for the msg bus contract to be deployed
	time.Sleep(5 * time.Second)

	// make sure the server is ready to receive requests
	gatewayAddressHTTP := fmt.Sprintf("http://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortHTTP)
	gatewayAddressWS := fmt.Sprintf("ws://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortWS)
	fmt.Println("gatewayAddressHTTP: ", gatewayAddressHTTP)
	fmt.Println("gatewayAddressWS: ", gatewayAddressWS)

	// make sure the server is ready to receive requests
	err := waitServerIsReady(gatewayAddressHTTP)
	require.NoError(t, err)

	// Server is now ready and we can create requests

	// Create users
	user0, err := NewUser([]wallet.Wallet{wallets.L2FaucetWallet}, gatewayAddressHTTP, gatewayAddressWS)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user0.UserID)

	user1, err := NewUser(wallets.SimObsWallets[0:2], gatewayAddressHTTP, gatewayAddressWS)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user1.UserID)

	user2, err := NewUser(wallets.SimObsWallets[2:4], gatewayAddressHTTP, gatewayAddressWS)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user2.UserID)

	// register all the accounts for that user
	err = user0.RegisterAccounts()
	require.NoError(t, err)
	err = user1.RegisterAccounts()
	require.NoError(t, err)
	err = user2.RegisterAccounts()
	require.NoError(t, err)

	// Transfer some funds to user1 and user2 wallets, because they need it to make transactions
	var amountToTransfer int64 = 1_000_000_000_000_000_000
	// Transfer some funds to user1 and user2 wallets, because they need it to make transactions
	_, err = TransferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	_, err = TransferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = TransferETHToAddress(user0.HTTPClient, user0.Wallets[0], user2.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = TransferETHToAddress(user0.HTTPClient, user0.Wallets[0], user2.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)

	// Print balances of all registered accounts to check if all accounts have funds
	err = user0.PrintUserAccountsBalances()
	require.NoError(t, err)
	err = user1.PrintUserAccountsBalances()
	require.NoError(t, err)
	err = user2.PrintUserAccountsBalances()
	require.NoError(t, err)

	// User0 deploys a contract that will later emit events
	bytecode := `60806040523480156200001157600080fd5b506040518060400160405280600381526020017f666f6f00000000000000000000000000000000000000000000000000000000008152506000908162000058919062000320565b506040518060400160405280600381526020017f666f6f0000000000000000000000000000000000000000000000000000000000815250600190816200009f919062000320565b5062000407565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200012857607f821691505b6020821081036200013e576200013d620000e0565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620001a87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000169565b620001b4868362000169565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b600062000201620001fb620001f584620001cc565b620001d6565b620001cc565b9050919050565b6000819050919050565b6200021d83620001e0565b620002356200022c8262000208565b84845462000176565b825550505050565b600090565b6200024c6200023d565b6200025981848462000212565b505050565b5b8181101562000281576200027560008262000242565b6001810190506200025f565b5050565b601f821115620002d0576200029a8162000144565b620002a58462000159565b81016020851015620002b5578190505b620002cd620002c48562000159565b8301826200025e565b50505b505050565b600082821c905092915050565b6000620002f560001984600802620002d5565b1980831691505092915050565b6000620003108383620002e2565b9150826002028217905092915050565b6200032b82620000a6565b67ffffffffffffffff811115620003475762000346620000b1565b5b6200035382546200010f565b6200036082828562000285565b600060209050601f83116001811462000398576000841562000383578287015190505b6200038f858262000302565b865550620003ff565b601f198416620003a88662000144565b60005b82811015620003d257848901518255600182019150602085019450602081019050620003ab565b86831015620003f25784890151620003ee601f891682620002e2565b8355505b6001600288020188555050505b505050505050565b6107ee80620004176000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063368b877214610051578063c2d366581461006d578063c5ced0361461008b578063e21f37ce146100a7575b600080fd5b61006b600480360381019061006691906103e6565b6100c5565b005b610075610126565b60405161008291906104ae565b60405180910390f35b6100a560048036038101906100a091906103e6565b6101b4565b005b6100af6101fe565b6040516100bc91906104ae565b60405180910390f35b80600090816100d491906106e6565b503373ffffffffffffffffffffffffffffffffffffffff167fe31c2ad953ded70272b94617f9181f8cc33755f1b40f4c706660f6ee0dfb634a8260405161011b91906104ae565b60405180910390a250565b60018054610133906104ff565b80601f016020809104026020016040519081016040528092919081815260200182805461015f906104ff565b80156101ac5780601f10610181576101008083540402835291602001916101ac565b820191906000526020600020905b81548152906001019060200180831161018f57829003601f168201915b505050505081565b80600190816101c391906106e6565b507f4fcdf2659dcf2254d2bce07af2baaf0c6ddf6da052dd241b2445a2f6398ae575816040516101f391906104ae565b60405180910390a150565b6000805461020b906104ff565b80601f0160208091040260200160405190810160405280929190818152602001828054610237906104ff565b80156102845780601f1061025957610100808354040283529160200191610284565b820191906000526020600020905b81548152906001019060200180831161026757829003601f168201915b505050505081565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6102f3826102aa565b810181811067ffffffffffffffff82111715610312576103116102bb565b5b80604052505050565b600061032561028c565b905061033182826102ea565b919050565b600067ffffffffffffffff821115610351576103506102bb565b5b61035a826102aa565b9050602081019050919050565b82818337600083830152505050565b600061038961038484610336565b61031b565b9050828152602081018484840111156103a5576103a46102a5565b5b6103b0848285610367565b509392505050565b600082601f8301126103cd576103cc6102a0565b5b81356103dd848260208601610376565b91505092915050565b6000602082840312156103fc576103fb610296565b5b600082013567ffffffffffffffff81111561041a5761041961029b565b5b610426848285016103b8565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561046957808201518184015260208101905061044e565b60008484015250505050565b60006104808261042f565b61048a818561043a565b935061049a81856020860161044b565b6104a3816102aa565b840191505092915050565b600060208201905081810360008301526104c88184610475565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061051757607f821691505b60208210810361052a576105296104d0565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026105927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610555565b61059c8683610555565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b60006105e36105de6105d9846105b4565b6105be565b6105b4565b9050919050565b6000819050919050565b6105fd836105c8565b610611610609826105ea565b848454610562565b825550505050565b600090565b610626610619565b6106318184846105f4565b505050565b5b818110156106555761064a60008261061e565b600181019050610637565b5050565b601f82111561069a5761066b81610530565b61067484610545565b81016020851015610683578190505b61069761068f85610545565b830182610636565b50505b505050565b600082821c905092915050565b60006106bd6000198460080261069f565b1980831691505092915050565b60006106d683836106ac565b9150826002028217905092915050565b6106ef8261042f565b67ffffffffffffffff811115610708576107076102bb565b5b61071282546104ff565b61071d828285610659565b600060209050601f831160018114610750576000841561073e578287015190505b61074885826106ca565b8655506107b0565b601f19841661075e86610530565b60005b8281101561078657848901518255600182019150602085019450602081019050610761565b868310156107a3578489015161079f601f8916826106ac565b8355505b6001600288020188555050505b50505050505056fea264697066735822122076146d8c796917af248ecb981f38094293788d92ad21704dc623fd8412cb9dc964736f6c63430008120033`
	abiString := `
	[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "newMessage",
				"type": "string"
			}
		],
		"name": "Message2Updated",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "newMessage",
				"type": "string"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "MessageUpdatedWithAddress",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "message",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "message2",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "newMessage",
				"type": "string"
			}
		],
		"name": "setMessage",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "newMessage",
				"type": "string"
			}
		],
		"name": "setMessage2",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]
	`

	_, contractAddress, err := DeploySmartContract(user0.HTTPClient, user0.Wallets[0], bytecode)
	require.NoError(t, err)
	fmt.Println("Deployed contract address: ", contractAddress)

	// contract abi
	contractAbi, err := abi.JSON(strings.NewReader(abiString))
	require.NoError(t, err)

	// check if contract was deployed and call one of the implicit getter functions
	// call getter for a message
	resultMessage, err := getStringValueFromSmartContractGetter(contractAddress, contractAbi, "message", user1.HTTPClient)
	require.NoError(t, err)

	// check if the value is the same as hardcoded in smart contract
	hardcodedMessageValue := "foo"
	assert.Equal(t, hardcodedMessageValue, resultMessage)

	// subscribe with all three users for all events in deployed contract
	var user0logs []types.Log
	var user1logs []types.Log
	var user2logs []types.Log
	subscribeToEvents([]gethcommon.Address{contractAddress}, nil, user0.WSClient, &user0logs)
	subscribeToEvents([]gethcommon.Address{contractAddress}, nil, user1.WSClient, &user1logs)
	subscribeToEvents([]gethcommon.Address{contractAddress}, nil, user2.WSClient, &user2logs)

	time.Sleep(time.Second)

	// user1 calls setMessage and setMessage2 on deployed smart contract with the account
	// that was registered as the first in OG
	user1MessageValue := "user1PrivateEvent"
	// interact with smart contract and cause events to be emitted
	_, err = InteractWithSmartContract(user1.HTTPClient, user1.Wallets[0], contractAbi, "setMessage", "user1PrivateEvent", contractAddress)
	require.NoError(t, err)
	_, err = InteractWithSmartContract(user1.HTTPClient, user1.Wallets[0], contractAbi, "setMessage2", "user1PublicEvent", contractAddress)
	require.NoError(t, err)

	// check if value was changed in the smart contract with the interactions above
	resultMessage, err = getStringValueFromSmartContractGetter(contractAddress, contractAbi, "message", user1.HTTPClient)
	require.NoError(t, err)
	assert.Equal(t, user1MessageValue, resultMessage)

	// user2 calls setMessage and setMessage2 on deployed smart contract with the account
	// that was registered as the second in OG
	user2MessageValue := "user2PrivateEvent"
	// interact with smart contract and cause events to be emitted
	_, err = InteractWithSmartContract(user2.HTTPClient, user2.Wallets[1], contractAbi, "setMessage", "user2PrivateEvent", contractAddress)
	require.NoError(t, err)
	_, err = InteractWithSmartContract(user2.HTTPClient, user2.Wallets[1], contractAbi, "setMessage2", "user2PublicEvent", contractAddress)
	require.NoError(t, err)

	// check if value was changed in the smart contract with the interactions above
	resultMessage, err = getStringValueFromSmartContractGetter(contractAddress, contractAbi, "message", user1.HTTPClient)
	require.NoError(t, err)
	assert.Equal(t, user2MessageValue, resultMessage)

	// wait a few seconds to be completely sure all events arrived
	time.Sleep(time.Second * 3)

	// Assert the number of logs received by each client
	// user0 should see two lifecycle events (1 for each interaction with setMessage2)
	assert.Equal(t, 2, len(user0logs))
	// user1 should see three events (two lifecycle events - same as user0) and event with his interaction with setMessage
	assert.Equal(t, 3, len(user1logs))
	// user2 should see three events (two lifecycle events - same as user0) and event with his interaction with setMessage
	assert.Equal(t, 2, len(user2logs))

	// Gracefully shutdown
	err = obscuroGwContainer.Stop()
	assert.NoError(t, err)
}

func testAreTxsMinted(t *testing.T, httpURL, wsURL string) { //nolint: unused
	// set up the ogClient
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

	w := wallet.NewInMemoryWalletFromConfig(genesis.TestnetPrefundedPK, integration.ObscuroChainID, testlog.Logger())
	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	// use a standard eth client via the og
	ethStdClient, err := ethclient.Dial(ogClient.HTTP())
	require.NoError(t, err)

	// check the balance
	balance, err := ethStdClient.BalanceAt(context.Background(), w.Address(), nil)
	require.NoError(t, err)
	require.True(t, big.NewInt(0).Cmp(balance) == -1)

	// issue a tx and check it was successfully minted
	txHash := transferRandomAddr(t, ethStdClient, w)
	receipt, err := ethStdClient.TransactionReceipt(context.Background(), txHash)
	assert.NoError(t, err)
	require.True(t, receipt.Status == 1)
}

func testErrorHandling(t *testing.T, httpURL, wsURL string) {
	// set up the ogClient
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

	// register an account
	w := wallet.NewInMemoryWalletFromConfig(genesis.TestnetPrefundedPK, integration.ObscuroChainID, testlog.Logger())
	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	// make requests to geth for comparison

	for _, req := range []string{
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":[],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getgetget","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1,"extra":"extra_field"}`,
		`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "0x1234"]],"id":1}`,
	} {
		// ensure the geth request is issued correctly (should return 200 ok with jsonRPCError)
		_, response, err := httputil.PostDataJSON(ogClient.HTTP(), []byte(req))
		require.NoError(t, err)

		// unmarshall the response to JSONRPCMessage
		jsonRPCError := wecommon.JSONRPCMessage{}
		err = json.Unmarshal(response, &jsonRPCError)
		require.NoError(t, err)

		// repeat the process for the gateway
		_, response, err = httputil.PostDataJSON(fmt.Sprintf("http://localhost:%d", integration.StartPortObscuroGatewayUnitTest), []byte(req))
		require.NoError(t, err)

		// we only care about format
		jsonRPCError = wecommon.JSONRPCMessage{}
		err = json.Unmarshal(response, &jsonRPCError)
		require.NoError(t, err)
	}
}

func testErrorsRevertedArePassed(t *testing.T, httpURL, wsURL string) {
	// set up the ogClient
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

	w := wallet.NewInMemoryWalletFromConfig(genesis.TestnetPrefundedPK, integration.ObscuroChainID, testlog.Logger())
	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	// use a standard eth client via the og
	ethStdClient, err := ethclient.Dial(ogClient.HTTP())
	require.NoError(t, err)

	// check the balance
	balance, err := ethStdClient.BalanceAt(context.Background(), w.Address(), nil)
	require.NoError(t, err)
	require.True(t, big.NewInt(0).Cmp(balance) == -1)

	// deploy errors contract
	deployTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(errorsContractBytecode),
	}

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = ethStdClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	receipt, err := integrationCommon.AwaitReceiptEth(context.Background(), ethStdClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	pack, _ := errorsContractABI.Pack("force_require")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, err.Error(), "execution reverted: Forced require")

	// convert error to WE error
	errBytes, err := json.Marshal(err)
	require.NoError(t, err)
	weError := wecommon.JSONError{}
	err = json.Unmarshal(errBytes, &weError)
	require.NoError(t, err)
	require.Equal(t, weError.Message, "execution reverted: Forced require")
	require.Equal(t, weError.Data, "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e466f726365642072657175697265000000000000000000000000000000000000")
	require.Equal(t, weError.Code, 3)

	pack, _ = errorsContractABI.Pack("force_revert")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, err.Error(), "execution reverted: Forced revert")

	pack, _ = errorsContractABI.Pack("force_assert")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, err.Error(), "execution reverted")
}

func transferRandomAddr(t *testing.T, client *ethclient.Client, w wallet.Wallet) common.TxHash { //nolint: unused
	ctx := context.Background()
	toAddr := datagenerator.RandomAddress()
	nonce, err := client.NonceAt(ctx, w.Address(), nil)
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	}
	assert.Nil(t, err)

	fmt.Println("Transferring from:", w.Address(), " to:", toAddr)

	signedTx, err := w.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = client.SendTransaction(ctx, signedTx)
	assert.Nil(t, err)

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client, signedTx.Hash(), time.Minute)
	assert.NoError(t, err)

	fmt.Println("Successfully minted the transaction - ", signedTx.Hash())
	return signedTx.Hash()
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int, nrSimWallets int) *params.SimWallets {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(nrSimWallets, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 1 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        startPort,
		WithPrefunding:   true,
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}

	return wallets
}

func waitServerIsReady(serverAddr string) error {
	for now := time.Now(); time.Since(now) < 30*time.Second; time.Sleep(500 * time.Millisecond) {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/health/", serverAddr))
		if err != nil {
			// give it time to boot up
			if strings.Contains(err.Error(), "connection") {
				continue
			}
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}
	}
	return fmt.Errorf("timed out before server was ready")
}

func ComputeContractAddress(sender gethcommon.Address, nonce uint64) (gethcommon.Address, error) {
	// RLP encode the byte array of the sender's address and nonce
	encoded, err := rlp.EncodeToBytes([]interface{}{sender, nonce})
	if err != nil {
		return gethcommon.Address{}, err
	}
	// Compute the Keccak-256 hash of the RLP encoded byte array
	hash := crypto.Keccak256(encoded)

	// The contract address is the last 20 bytes of this hash
	return gethcommon.BytesToAddress(hash[len(hash)-20:]), nil
}

func TransferETHToAddress(client *ethclient.Client, wallet wallet.Wallet, toAddress gethcommon.Address, amount int64) (*types.Receipt, error) {
	transferTx1 := types.LegacyTx{
		Nonce:    wallet.GetNonceAndIncrement(),
		To:       &toAddress,
		Value:    big.NewInt(amount),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     nil,
	}
	signedTx, err := wallet.SignTransaction(&transferTx1)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return AwaitTransactionReceipt(context.Background(), client, signedTx.Hash(), 2*time.Second)
}

func DeploySmartContract(client *ethclient.Client, wallet wallet.Wallet, bytecode string) (*types.Receipt, gethcommon.Address, error) {
	contractNonce := wallet.GetNonceAndIncrement()
	contractTx := types.LegacyTx{
		Nonce:    contractNonce,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(bytecode),
	}

	signedTx, err := wallet.SignTransaction(&contractTx)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}

	// await for the transaction to be included in a block
	txReceiptContractCreation, err := AwaitTransactionReceipt(context.Background(), client, signedTx.Hash(), 2*time.Second)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}

	// get contract address
	contractAddress, err := ComputeContractAddress(wallet.Address(), contractNonce)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}

	return txReceiptContractCreation, contractAddress, nil
}

func InteractWithSmartContract(client *ethclient.Client, wallet wallet.Wallet, contractAbi abi.ABI, methodName string, methodParam string, contractAddress gethcommon.Address) (*types.Receipt, error) {
	contractInteractionData, err := contractAbi.Pack(methodName, methodParam)
	if err != nil {
		return nil, err
	}

	interactionTx := types.LegacyTx{
		Nonce:    wallet.GetNonceAndIncrement(),
		To:       &contractAddress,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     contractInteractionData,
	}
	signedTx, err := wallet.SignTransaction(&interactionTx)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	txReceipt, err := AwaitTransactionReceipt(context.Background(), client, signedTx.Hash(), 2*time.Second)
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}

func getStringValueFromSmartContractGetter(contractAddress gethcommon.Address, contractAbi abi.ABI, method string, client *ethclient.Client) (string, error) {
	contract := bind.NewBoundContract(contractAddress, contractAbi, client, client, client)
	var resultMessage string
	callOpts := &bind.CallOpts{}
	results := make([]interface{}, 1)
	results[0] = &resultMessage
	err := contract.Call(callOpts, &results, method)
	if err != nil {
		return "", err
	}

	return resultMessage, nil
}

func AwaitTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash gethcommon.Hash, timeout time.Duration) (*types.Receipt, error) {
	timeoutStrategy := retry.NewTimeoutStrategy(timeout, time.Second)
	return AwaitTransactionReceiptWithRetryStrategy(ctx, client, txHash, timeoutStrategy)
}

func AwaitTransactionReceiptWithRetryStrategy(ctx context.Context, client *ethclient.Client, txHash gethcommon.Hash, retryStrategy retry.Strategy) (*types.Receipt, error) {
	retryStrategy.Reset()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled before receipt was received")
		case <-time.After(retryStrategy.NextRetryInterval()):
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
			if retryStrategy.Done() {
				return nil, fmt.Errorf("receipt not found - %s - %w", retryStrategy.Summary(), err)
			}
		}
	}
}

func subscribeToEvents(addresses []gethcommon.Address, topics [][]gethcommon.Hash, client *ethclient.Client, logs *[]types.Log) {
	// Make a subscription
	filterQuery := ethereum.FilterQuery{
		Addresses: addresses,
		FromBlock: big.NewInt(0), // todo (@ziga) - without those we get errors - fix that and make them configurable
		ToBlock:   big.NewInt(10000),
		Topics:    topics,
	}
	logsCh := make(chan types.Log)

	subscription, err := client.SubscribeFilterLogs(context.Background(), filterQuery, logsCh)
	if err != nil {
		fmt.Printf("Failed to subscribe to filter logs: %v\n", err)
	}
	// todo (@ziga) - unsubscribe when it is fixed...
	// defer subscription.Unsubscribe() // cleanup

	// Listen for logs in a goroutine
	go func() {
		for {
			select {
			case err := <-subscription.Err():
				fmt.Printf("Error from logs subscription: %v\n", err)
				return
			case log := <-logsCh:
				// append logs to be visible from the main thread
				*logs = append(*logs, log)
			}
		}
	}()
}
