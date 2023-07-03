package host

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/responses"
)

// Host is the half of the Obscuro node that lives outside the enclave.
type Host interface {
	Config() *config.HostConfig
	DB() *db.DB
	EnclaveClient() common.Enclave

	// Start initializes the main loop of the host.
	Start() error
	// SubmitAndBroadcastTx submits an encrypted transaction to the enclave, and broadcasts it to the other hosts on the network.
	SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error)
	// Subscribe feeds logs matching the encrypted log subscription to the matchedLogs channel.
	Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogs chan []byte) error
	// Unsubscribe terminates a log subscription between the host and the enclave.
	Unsubscribe(id rpc.ID)
	// Stop gracefully stops the host execution.
	Stop() error

	// HealthCheck returns the health status of the host + enclave + db
	HealthCheck() (*HealthCheck, error)

	P2PSubscriber // todo (@matt) remove this from host interface and host tests when it's done at a per-enclave level
}

type P2PSubscriber interface {
	// ReceiveTx processes a transaction received from a peer host.
	ReceiveTx(tx common.EncryptedTx)
	// ReceiveBatches receives a set of batches from a peer host.
	ReceiveBatches(batches common.EncodedBatchMsg)
	// ReceiveBatchRequest receives a batch request from a peer host. Used during catch-up.
	ReceiveBatchRequest(batchRequest common.EncodedBatchRequest)
}

// P2P is the layer responsible for sending and receiving messages to Obscuro network peers.
type P2P interface {
	StartListening(callback P2PSubscriber)
	StopListening() error
	UpdatePeerList([]string)
	// SendTxToSequencer sends the encrypted transaction to the sequencer.
	SendTxToSequencer(tx common.EncryptedTx) error
	// BroadcastBatch sends the batch to every other node on the network.
	BroadcastBatch(batchMsg *BatchMsg) error
	// RequestBatchesFromSequencer requests batches from the sequencer.
	RequestBatchesFromSequencer(batchRequest *common.BatchRequest) error
	// SendBatches sends batches to a specific node, in response to a batch request.
	SendBatches(batchMsg *BatchMsg, to string) error

	// Status returns the status of the p2p communications.
	Status() *P2PStatus

	// HealthCheck returns whether the p2p lib is healthy.
	HealthCheck() bool
}

type BlockStream struct {
	Stream <-chan *types.Block // the channel which will receive the consecutive, canonical blocks
	Stop   func()              // function to permanently stop the stream and clean up any associated processes/resources
}

type BatchMsg struct {
	Batches   []*common.ExtBatch // The batches being sent.
	IsCatchUp bool               // Whether these batches are being sent as part of a catch-up request.
}
