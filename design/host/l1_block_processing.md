# Standardizing L1 Data Processing

## Requirements

1. Standardise the way in which we process l1 blocks to find the relevant data needed for processing on the L2
2. Reduce the processing load on the enclave

## Current Problems
* Multiple redundant loops through L1 block data
* Inconsistent processing patterns
* Unnecessary load on the enclave
* Scattered responsibility for L1 data extraction

## Proposed Solution

Filter logs from blocks using the Management Contract and MessageBus Contract address. Build a map of emitted event types
against the transactions that created them. The events we care about: 

* Initialize secret
* Request Secret
* Secret Response 
* Rollup
* Cross chain messages
* Value transfers
* Enclave granted sequencer 
* Enclave sequencer revoked

```go
const (
    RollupTx L1TxType = iota
    InitialiseSecretTx
    SecretRequestTx
    SecretResponseTx
    CrossChainMessageTx
    CrossChainValueTranserTx
    SequencerAddedTx
    SequencerRevokedTx
    SetImportantContractsTx
)

// ProcessedL1Data is submitted to the enclave by the guardian
type ProcessedL1Data struct {
    BlockHeader *types.Header
    Events      []L1Event
}

// L1Event represents a single event type and its associated transactions
type L1Event struct {
	Type uint8
    Txs  []*L1TxData
}

// L1TxData represents an L1 transaction that are relevant to us
type L1TxData struct {
    Transaction        *types.Transaction
    Receipt            *types.Receipt
    Blobs              []*kzg4844.Blob     // Only populated for blob transactions
    SequencerEnclaveID gethcommon.Address  // Only non-zero when a new enclave is added as a sequencer
    CrossChainMessages CrossChainMessages  // Only populated for xchain messages
    ValueTransfers     ValueTransferEvents // Only populated for xchain transfers
    Proof              []byte              // Some merkle proof TBC
}
```
## Guardian
In the guardian we do all the transaction extraction to look for the event types we care about and then submit a
`ProcessedL1Data` object to the enclave via gRPC in the `SubmitL1Block` function.

`TODO` what proof to submit?

## Enclave side

On the enclave side we handle each of the `processedData.GetEvents[L1TxType]` individually and don't have duplicate loops
through the transactions.

Correct ordering of these event processing is going to be the biggest pain point I suspect.
