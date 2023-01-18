package gethnetwork

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartEth2Network(t *testing.T) {
	network := NewEth2Network(2)
	assert.NotNil(t, network.Start())

	defer network.Stop()

	for testName, test := range map[string]func(t *testing.T, network *Eth2Network){
		"GenesisParamsAreUsed": genesisParamsAreUsed,
	} {
		t.Run(testName, func(t *testing.T) {
			test(t, network)
		})
	}
}

func genesisParamsAreUsed(t *testing.T, network *Eth2Network) {

}
