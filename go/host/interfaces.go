package host

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
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
	// ReceiveRollup processes a rollup received from a peer host.
	ReceiveRollup(r common.EncodedRollup)
	// ReceiveTx processes a transaction received from a peer host.
	ReceiveTx(tx common.EncryptedTx)
	// Stop gracefully stops the host execution.
	Stop()
}

// MockHost extends Host with additional methods that are only used for integration testing.
type MockHost interface {
	Host

	// MockedNewHead receives the notification of new blocks.
	// TODO - Remove this method.
	MockedNewHead(b common.EncodedBlock, p common.EncodedBlock)
	// MockedNewFork receives the notification of a new fork.
	MockedNewFork(b []common.EncodedBlock)
	// ConnectToEthNode connects the Aggregator to the Ethereum node.
	ConnectToEthNode(node ethadapter.EthClient)
}

// P2P is the layer responsible for sending and receiving messages to Obscuro network peers.
type P2P interface {
	StartListening(callback Host)
	StopListening() error
	UpdatePeerList([]string)
	BroadcastRollup(r common.EncodedRollup) error
	BroadcastTx(tx common.EncryptedTx) error
}

type StatsCollector interface {
	// L2Recalc - called when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id gethcommon.Address)
	NewBlock(block *types.Block)
	NewRollup(node gethcommon.Address)
	RollupWithMoreRecentProof()
}
