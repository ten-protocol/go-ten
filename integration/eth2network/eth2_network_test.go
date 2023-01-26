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
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	_testBasePort = 18545
	_numTestNodes = 1
)

func TestStartEth2Network(t *testing.T) {
	binDir, err := EnsureBinariesExist()
	assert.Nil(t, err)

	chainID := int(datagenerator.RandomUInt64())
	randomWallet := datagenerator.RandomWallet(int64(chainID))

	network := NewEth2Network(
		binDir,
		_testBasePort,
		_testBasePort+100,
		_testBasePort+200,
		_testBasePort+300,
		chainID,
		_numTestNodes,
		1,
		[]string{randomWallet.Address().Hex()},
	)
	// wait until the merge has happened
	assert.Nil(t, network.Start())

	defer network.Stop()

	// test input configurations
	t.Run("areConfigsUphold", func(t *testing.T) {
		areConfigsUphold(t, randomWallet.Address(), chainID)
	})

	// test number of nodes
	t.Run("numberOfNodes", func(t *testing.T) {
		numberOfNodes(t)
	})

	// txs are minted
	t.Run("txsAreMinted", func(t *testing.T) {
		txsAreMinted(t, randomWallet)
	})
}

func areConfigsUphold(t *testing.T, addr gethcommon.Address, chainID int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", _testBasePort)
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
	for i := 0; i < _numTestNodes; i++ {
		url := fmt.Sprintf("http://127.0.0.1:%d", _testBasePort+i)

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
		assert.Equal(t, res["result"], fmt.Sprintf("0x%x", _numTestNodes-1))
	}
}

func txsAreMinted(t *testing.T, w wallet.Wallet) {
	ethClient, err := ethadapter.NewEthClient("127.0.0.1", _testBasePort+100, 30*time.Second, common.L2Address{}, gethlog.New())
	assert.Nil(t, err)

	nonce, err := ethClient.Nonce(w.Address())
	assert.Nil(t, err)
	w.SetNonce(nonce)

	toAddr := datagenerator.RandomAddress()
	estimatedTx, err := ethClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce: w.GetNonceAndIncrement(),
		To:    &toAddr,
		Value: big.NewInt(100),
	}, w.Address())
	assert.Nil(t, err)

	signedTx, err := w.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = ethClient.SendTransaction(signedTx)
	assert.Nil(t, err)

	// make sure it's mined into a block within an acceptable time
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
