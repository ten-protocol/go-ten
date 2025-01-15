package eth2network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/datagenerator"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

func TestEnsureBinariesAreAvail(t *testing.T) {
	path, err := EnsureBinariesExist()
	assert.Nil(t, err)
	t.Logf("Successfully downloaded files to %s", path)
}

func TestStartPosEth2Network(t *testing.T) {
	binDir, err := EnsureBinariesExist()
	assert.Nil(t, err)

	startPort := integration.TestPorts.TestStartPosEth2NetworkPort
	network := NewPosEth2Network(
		binDir,
		startPort+integration.DefaultGethNetworkPortOffset,
		startPort+integration.DefaultPrysmP2PPortOffset,
		startPort+integration.DefaultGethAUTHPortOffset,
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethHTTPPortOffset,
		startPort+integration.DefaultPrysmRPCPortOffset,
		startPort+integration.DefaultPrysmGatewayPortOffset,
		integration.EthereumChainID,
		3*time.Minute,
	)

	// wait until the merge has happened
	assert.Nil(t, network.Start())

	defer network.Stop() //nolint: errcheck

	// test input configurations
	t.Run("areConfigsUphold", func(t *testing.T) {
		areConfigsUphold(t, startPort, gethcommon.HexToAddress(integration.GethNodeAddress), integration.EthereumChainID)
	})

	// test number of nodes
	t.Run("numberOfNodes", func(t *testing.T) {
		numberOfNodes(t, startPort)
	})

	minerWallet := wallet.NewInMemoryWalletFromConfig(
		integration.GethNodePK,
		integration.EthereumChainID,
		gethlog.New())

	t.Run("txsAreMinted", func(t *testing.T) {
		txsAreMinted(t, startPort, minerWallet)
	})
}

func areConfigsUphold(t *testing.T, startPort int, addr gethcommon.Address, chainID int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultGethHTTPPortOffset)
	conn, err := ethclient.Dial(url)
	assert.Nil(t, err)

	at, err := conn.BalanceAt(context.Background(), addr, nil)
	assert.Nil(t, err)
	assert.True(t, at.Cmp(big.NewInt(1)) == 1)

	id, err := conn.NetworkID(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(chainID), id.Int64())
}

func numberOfNodes(t *testing.T, startPort int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultGethHTTPPortOffset)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		strings.NewReader(`{"jsonrpc": "2.0", "method": "net_peerCount", "params": [], "id": 1}`),
	)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	assert.Nil(t, err)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &res)
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf("0x%x", 0), res["result"])
}

func txsAreMinted(t *testing.T, startPort int, w wallet.Wallet) {
	var err error

	ethClient, err := ethadapter.NewEthClient("127.0.0.1", uint(startPort+integration.DefaultGethWSPortOffset), 30*time.Second, gethlog.New())
	assert.Nil(t, err)

	toAddr := datagenerator.RandomAddress()
	estimatedTx, err := ethadapter.SetTxGasPrice(context.Background(), ethClient, &types.LegacyTx{
		To:    &toAddr,
		Value: big.NewInt(100),
	}, w.Address(), w.GetNonceAndIncrement(), 0)
	assert.Nil(t, err)

	signedTx, err := w.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = ethClient.SendTransaction(signedTx)
	assert.Nil(t, err)

	fmt.Printf("Created Tx: %s", signedTx.Hash().Hex())
	// make sure it's mined into a block within an acceptable time and avail in all nodes
	fmt.Printf("Checking for tx receipt for %s", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < 30*time.Second; time.Sleep(time.Second) {
		receipt, err = ethClient.TransactionReceipt(signedTx.Hash())
		if err == nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			t.Fatal(err)
		}
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds - receipt: %+v", 30*time.Second, receipt)
	}
}
