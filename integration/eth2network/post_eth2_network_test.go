package eth2network

import (
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/integration"
	"testing"
	"time"
)

const (
	startPort = integration.StartPortEth2NetworkTests
)

func TestStartPostEth2Network(t *testing.T) {
	binDir, err := EnsureBinariesExist()
	assert.Nil(t, err)

	network := NewPosEth2Network(
		binDir,
		startPort+integration.DefaultGethAUTHPortOffset, //RPC
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethNetworkPortOffset, //HTTP
		startPort+integration.DefaultPrysmP2PPortOffset,    //RPC
		2*time.Minute, //RPC
	)
	// wait until the merge has happened
	assert.Nil(t, network.Start())

	defer network.Stop() //nolint: errcheck

}
