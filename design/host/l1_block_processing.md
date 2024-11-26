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

Loop through transactions to find events we care about in the guardian and submit these with the block to the enclave:

* Rollup
* Secret request
* Cross chain messages
* GrantSequencer events

```go
// L1EventType identifies different types of L1 transactions we care about
type L1EventType int

const (
    RollupEvent L1EventType = iota
    SecretRequestEvent
    CrossChainMessageEvent
    SequencerAddedEvent
)

// ProcessedL1Data is submitted to the enclave by the guardian 
type ProcessedL1Data struct {
    BlockHeader *types.Header
    Events      map[L1EventType][]*L1EventData
}

// L1Event represents a processed L1 transaction that's relevant to us
type L1EventData struct {
    Type        L1EventType
    Transaction *types.Transaction
    Receipt     *types.Receipt
    BlockHeader *types.Header
    Blobs       []*kzg4844.Blob          // Only populated for blob transactions
}
```
## Guardian 
In the guardian we do all the transaction filtering to look for the event types we care about and then submit a
`ProcessedL1Data` object to the enclave via gRPC in the `SubmitL1Block` function. 

`TODO` what proof to submit?

```go
// This can be added to the guardian, or we include in one of the existing guardian services
type L1Processor struct {
    mgmtContractLib mgmtcontractlib.MgmtContractLib
    blobResolver    BlobResolver
    logger          gethlog.Logger
}

func (p *L1Processor) ProcessBlock(br *common.BlockAndReceipts) (*common.ProcessedL1Data, error) {
    processed := &common.ProcessedL1Data{
        BlockHeader: br.BlockHeader,
        Events:     make(map[L1EventType][]*common.L1EventData),
    }

    // Extract blobs and hashes once
    blobs, blobHashes, err := p.extractBlobsAndHashes(br)
    if err != nil {
        return nil, fmt.Errorf("failed to extract blobs: %w", err)
    }

    // Single pass through transactions
    for _, tx := range *br.RelevantTransactions() {
        decodedTx := p.mgmtContractLib.DecodeTx(tx)
        if decodedTx == nil {
            continue
        }

        // Find receipt for this transaction
        receipt := br.GetReceipt(tx.Hash())
        
        switch t := decodedTx.(type) {
        case *ethadapter.L1RollupHashes:
            // Verify blob hashes match for rollups
            if err := verifyBlobHashes(t, blobHashes); err != nil {
                p.logger.Warn("Invalid rollup blob hashes", "tx", tx.Hash(), "error", err)
                continue
            }
            
            processed.Events[RollupEvent] = append(processed.Events[RollupEvent], &common.L1EventData{
                TxHash:      tx.Hash(),
                Transaction: tx,
                Receipt:     receipt,
                Blobs:       blobs,
            })

        case *ethadapter.L1RequestSecretTx:
            processed.Events[SecretRequestEvent] = append(processed.Events[SecretRequestEvent], &common.L1EventData{
                TxHash:      tx.Hash(),
                Transaction: tx,
                Receipt:     receipt,
            })

        case *ethadapter.L1InitializeSecretTx:
            processed.Events[SecretRequestEvent] = append(processed.Events[SecretRequestEvent], &common.L1EventData{
                TxHash:      tx.Hash(),
                Transaction: tx,
                Receipt:     receipt,
            })

        // Add other event types...
        }
    }

    return processed, nil
}
```

## Enclave side 

On the enclave side we handle each of the `processedData.Events[L1EventType]` individually and don't have duplicate loops
through the transactions. 

Correct ordering of these event processing is going to be the biggest pain point I suspect. 

```go
func (e *enclaveAdminService) ingestL1Block(ctx context.Context, processedData *common.ProcessedL1Data) (*components.BlockIngestionType, error) {
    
    // Process block first to ensure it's valid and get ingestion type
    ingestion, err := e.l1BlockProcessor.Process(ctx, processedData.BlockHeader)
    if err != nil {
        if errors.Is(err, errutil.ErrBlockAncestorNotFound) || errors.Is(err, errutil.ErrBlockAlreadyProcessed) {
            e.logger.Debug("Did not ingest block", log.ErrKey, err, log.BlockHashKey, processedData.BlockHeader.Hash())
        } else {
            e.logger.Warn("Failed ingesting block", log.ErrKey, err, log.BlockHashKey, processedData.BlockHeader.Hash())
        }
        return nil, err
    }

    // Process each event type in order
    var secretResponses []*common.ProducedSecretResponse

    // rollups
    if rollupEvents, exists := processedData.Events[RollupEvent]; exists {
        for _, event := range rollupEvents {
            if err := e.rollupConsumer.ProcessRollup(ctx, event); err != nil {
                if !errors.Is(err, components.ErrDuplicateRollup) {
                    e.logger.Error("Failed processing rollup", log.ErrKey, err)
                    // Continue processing other events even if one rollup fails
                }
            }
        }
    }

    // secret requests
    if secretEvents, exists := processedData.Events[SecretRequestEvent]; exists {
        for _, event := range secretEvents {
            resp, err := e.sharedSecretProcessor.ProcessSecretRequest(ctx, event)
            if err != nil {
                e.logger.Error("Failed to process secret request", "tx", event.TxHash, "error", err)
                continue
            }
            if resp != nil {
                secretResponses = append(secretResponses, resp)
            }
        }
    }
    
    // cross chain messages & add sequencer etc

    // Handle any L1 reorg/fork
    if ingestion.IsFork() {
        e.registry.OnL1Reorg(ingestion)
        if err := e.service.OnL1Fork(ctx, ingestion.ChainFork); err != nil {
            return nil, err
        }
    }

    // Add secret responses to ingestion result
    ingestion.SecretResponses = secretResponses

    return ingestion, nil
}
```