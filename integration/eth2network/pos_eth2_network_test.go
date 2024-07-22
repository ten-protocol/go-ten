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
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/datagenerator"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	_startPort = integration.StartPortEth2NetworkTests
)

func TestEnsureBinariesAreAvail(t *testing.T) {
	path, err := EnsureBinariesExist()
	assert.Nil(t, err)
	t.Logf("Successfully downloaded files to %s", path)
}

func TestStartPosEth2Network(t *testing.T) {
	binDir, err := EnsureBinariesExist()
	assert.Nil(t, err)

	network := NewPosEth2Network(
		binDir,
		_startPort+integration.DefaultGethNetworkPortOffset,
		_startPort+integration.DefaultPrysmP2PPortOffset,
		_startPort+integration.DefaultGethAUTHPortOffset, // RPC
		_startPort+integration.DefaultGethWSPortOffset,
		_startPort+integration.DefaultGethHTTPPortOffset,
		_startPort+integration.DefaultPrysmRPCPortOffset, // RPC
		integration.EthereumChainID,
		6*time.Minute,
	)

	// wait until the merge has happened
	assert.Nil(t, network.Start())

	defer network.Stop() //nolint: errcheck

	// test input configurations
	t.Run("areConfigsUphold", func(t *testing.T) {
		areConfigsUphold(t, gethcommon.HexToAddress(integration.GethNodeAddress), integration.EthereumChainID)
	})

	// test number of nodes
	t.Run("numberOfNodes", func(t *testing.T) {
		numberOfNodes(t)
	})

	minerWallet := wallet.NewInMemoryWalletFromConfig(
		integration.GethNodePK,
		integration.EthereumChainID,
		gethlog.New())

	t.Run("txsAreMinted", func(t *testing.T) {
		txsAreMinted(t, minerWallet)
	})
}

func areConfigsUphold(t *testing.T, addr gethcommon.Address, chainID int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", _startPort+integration.DefaultGethHTTPPortOffset)
	conn, err := ethclient.Dial(url)
	assert.Nil(t, err)

	at, err := conn.BalanceAt(context.Background(), addr, nil)
	assert.Nil(t, err)
	assert.True(t, at.Cmp(big.NewInt(1)) == 1)

	id, err := conn.NetworkID(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(chainID), id.Int64())
}

func numberOfNodes(t *testing.T) {
	url := fmt.Sprintf("http://127.0.0.1:%d", _startPort+integration.DefaultGethHTTPPortOffset)

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

func txsAreMinted(t *testing.T, w wallet.Wallet) {
	var err error

	ethClient, err := ethadapter.NewEthClient("127.0.0.1", uint(_startPort+integration.DefaultGethWSPortOffset), 30*time.Second, common.L2Address{}, gethlog.New())
	assert.Nil(t, err)

	toAddr := datagenerator.RandomAddress()
	estimatedTx, err := ethClient.PrepareTransactionToSend(context.Background(), &types.LegacyTx{
		To:    &toAddr,
		Value: big.NewInt(100),
	}, w.Address())
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
