package eth2network

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"
)

func TestStartEth2Network(t *testing.T) {
	binDir, err := EnsureBinariesExist()
	assert.Nil(t, err)

	randomAddr := datagenerator.RandomAddress()
	network := NewEth2Network(binDir, 8545, 9000, 30303, 8551, 2, 1, []string{randomAddr.Hex()})
	// wait until the merge has happened
	assert.Nil(t, network.Start())

	defer network.Stop()

	// test prefunding
	t.Run("isAddressPrefunded", func(t *testing.T) {
		isAddressPrefunded(t, randomAddr)
	})

	// run additional tests
	for testName, test := range map[string]func(t *testing.T, network *Eth2Network){
		"GenesisParamsAreUsed": genesisParamsAreUsed,
	} {
		t.Run(testName, func(t *testing.T) {
			test(t, network)
		})
	}
}

func isAddressPrefunded(t *testing.T, addr gethcommon.Address) {
	conn, err := ethclient.Dial("http://127.0.0.1:8545")
	assert.Nil(t, err)

	at, err := conn.BalanceAt(context.Background(), addr, nil)
	assert.Nil(t, err)

	assert.True(t, at.Cmp(big.NewInt(1)) == 1)
}

func genesisParamsAreUsed(t *testing.T, network *Eth2Network) {
	fmt.Println("Stopping now")
}
