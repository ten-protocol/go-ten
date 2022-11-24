package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"

	"github.com/ethereum/go-ethereum/core/types"
)

type (
	StateRoot = common.Hash
	TxHash    = common.Hash

	// MainNet aliases
	L1Address     = common.Address
	L1RootHash    = common.Hash
	L1Block       = types.Block
	L1Transaction = types.Transaction
	L1Receipt     = types.Receipt
	L1Receipts    = types.Receipts

	// Local Obscuro aliases
	L2RootHash     = common.Hash
	L2Tx           = types.Transaction
	L2Transactions = types.Transactions
	L2Address      = common.Address
	L2Receipt      = types.Receipt
	L2Receipts     = types.Receipts

	CrossChainMessage     = MessageBus.StructsCrossChainMessage
	CrossChainMessages    = []CrossChainMessage
	EncryptedTx           []byte // A single transaction, encoded as a JSON list of transaction binary hexes and encrypted using the enclave's public key
	EncryptedTransactions []byte // A blob of encrypted transactions, as they're stored in the rollup, with the nonce prepended.

	EncryptedParamsGetBalance      []byte // The params for an RPC getBalance request, as a JSON object encrypted with the public key of the enclave.
	EncryptedParamsCall            []byte // As above, but for an RPC call request.
	EncryptedParamsGetTxByHash     []byte // As above, but for an RPC getTransactionByHash request.
	EncryptedParamsGetTxReceipt    []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedParamsLogSubscription []byte // As above, but for an RPC logs subscription request.
	EncryptedParamsSendRawTx       []byte // As above, but for an RPC sendRawTransaction request.
	EncryptedParamsGetTxCount      []byte // As above, but for an RPC getTransactionCount request.
	EncryptedParamsEstimateGas     []byte // As above, but for an RPC estimateGas request.
	EncryptedParamsGetLogs         []byte // As above, but for an RPC getLogs request.

	EncryptedResponseGetBalance   []byte // The response for an RPC getBalance request, as a JSON object encrypted with the viewing key of the user.
	EncryptedResponseCall         []byte // As above, but for an RPC call request.
	EncryptedResponseGetTxReceipt []byte // As above, but for an RPC getTransactionReceipt request.
	EncryptedResponseSendRawTx    []byte // As above, but for an RPC sendRawTransaction request.
	EncryptedResponseGetTxByHash  []byte // As above, but for an RPC getTransactionByHash request.
	EncryptedResponseGetTxCount   []byte // As above, but for an RPC getTransactionCount request.
	EncryptedLogSubscription      []byte // As above, but for a log subscription request.
	EncryptedLogs                 []byte // As above, but for a log subscription response.
	EncryptedResponseEstimateGas  []byte // As above, but for an RPC estimateGas response.
	EncryptedResponseGetLogs      []byte // As above, but for an RPC getLogs request.

	Nonce         = uint64
	EncodedRollup []byte
	EncodedBatch  []byte
)

const (
	L2GenesisHeight = uint64(0)
	L1GenesisHeight = uint64(0)
	// HeightCommittedBlocks is the number of blocks deep a transaction must be to be considered safe from reorganisations.
	HeightCommittedBlocks = 15
)

// AttestationReport represents a signed attestation report from a TEE and some metadata about the source of it to verify it
type AttestationReport struct {
	Report      []byte         // the signed bytes of the report which includes some encrypted identifying data
	PubKey      []byte         // a public key that can be used to send encrypted data back to the TEE securely (should only be used once Report has been verified)
	Owner       common.Address // address identifying the owner of the TEE which signed this report, can also be verified from the encrypted Report data
	HostAddress string         // the IP address on which the host can be contacted by other Obscuro hosts for peer-to-peer communication
}

type (
	EncryptedSharedEnclaveSecret []byte
	EncodedAttestationReport     []byte
)
