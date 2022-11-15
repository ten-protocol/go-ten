package host

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

// Host is the half of the Obscuro node that lives outside the enclave.
type Host interface {
	Config() *config.HostConfig
	DB() *db.DB
	EnclaveClient() common.Enclave

	// Start initializes the main loop of the host.
	Start()
	// SubmitAndBroadcastTx submits an encrypted transaction to the enclave, and broadcasts it to the other hosts on the network.
	SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (common.EncryptedResponseSendRawTx, error)
	// ReceiveTx processes a transaction received from a peer host.
	ReceiveTx(tx common.EncryptedTx)
	// Subscribe feeds logs matching the encrypted log subscription to the matchedLogs channel.
	Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogs chan []byte) error
	// Unsubscribe terminates a log subscription between the host and the enclave.
	Unsubscribe(id rpc.ID)
	// Stop gracefully stops the host execution.
	Stop()

	// HealthCheck returns the health status of the host + enclave + db
	HealthCheck() (bool, error)
}

// MockHost extends Host with additional methods that are only used for integration testing.
// todo - remove this interface
type MockHost interface {
	Host

	// MockedNewHead receives the notification of new blocks.
	// TODO - Remove this method.
	MockedNewHead(b common.EncodedBlock, p common.EncodedBlock)
	// MockedNewFork receives the notification of a new fork.
	MockedNewFork(b []common.EncodedBlock)
}

// P2P is the layer responsible for sending and receiving messages to Obscuro network peers.
type P2P interface {
	StartListening(callback Host)
	StopListening() error
	UpdatePeerList([]string)
	BroadcastTx(tx common.EncryptedTx) error
}

// ReconnectingBlockProvider interface allows host to monitor and await L1 blocks.
//
// The stream channels provide the blocks the way the enclave expects to be fed (consecutive canonical blocks)
//
// ReconnectingBlockProvider handles:
//
//   - reconnecting to the source, it will recover if it can and continue streaming from where it left off
//
//   - forks: block provider only sends blocks that are *currently* canonical. If there was a fork then it will replay
//     from the block after the fork. For example:
//
//     12a --> 13a --> 14a -->
//     \-> 13b --> 14b --> 15b
//     If block provider had just published 14a and then discovered the 'b' fork is canonical, it would next publish 13b, 14b, 15b.
type ReconnectingBlockProvider interface {
	StartStreamingFromHeight(height *big.Int) (<-chan *types.Block, error)
	StartStreamingFromHash(latestHash gethcommon.Hash) (<-chan *types.Block, error)
	Stop()
	IsLive(hash gethcommon.Hash) bool // returns true if hash is of the latest known L1 head block
}
