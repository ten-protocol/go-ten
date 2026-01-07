# Guardian Queue Architecture Proposal

## Problem Statement
Current guardian has multiple submission paths competing for enclave access, causing:
- Race conditions between lock release and state updates
- TryLock churn with errEnclaveBusy retries
- Complex error recovery with recursive calls
- State desynchronization during slow block processing

## Proposed Solution: Hybrid Queue Architecture

### Design Principles
1. Single consumer for enclave submissions (sequential processing)
2. Multiple producers can enqueue without blocking
3. Smart deduplication prevents redundant work
4. Maintain backward compatibility with existing interfaces

### Core Components

```go
type EnclaveSubmissionQueue struct {
    l1Blocks     chan *L1SubmissionRequest
    l2Batches    chan *L2SubmissionRequest
    transactions chan *TxSubmissionRequest
    
    state        *StateTracker
    enclaveClient EnclaveClient
    logger       Logger
    
    // Deduplication
    pendingBlocks sync.Map // hash -> struct{}
    
    // Metrics
    queueDepth   atomic.Int64
    rejections   atomic.Int64
}

type L1SubmissionRequest struct {
    Block     *types.Header
    IsLatest  bool
    Priority  Priority
    ResultCh  chan<- SubmissionResult // Optional callback
}
```

### Processing Loop

```go
func (q *EnclaveSubmissionQueue) processLoop() {
    for {
        select {
        case req := <-q.l1Blocks:
            q.handleL1Block(req)
            
        case req := <-q.l2Batches:
            q.handleL2Batch(req)
            
        case req := <-q.transactions:
            q.handleTransaction(req)
            
        case <-q.stopCh:
            return
        }
    }
}

func (q *EnclaveSubmissionQueue) handleL1Block(req *L1SubmissionRequest) {
    // Check if already processed
    if _, exists := q.pendingBlocks.LoadAndDelete(req.Block.Hash()); !exists {
        return // duplicate, skip
    }
    
    // Submit to enclave (no lock needed - single consumer)
    resp, err := q.enclaveClient.SubmitL1Block(req.Block)
    
    // Handle response
    if err == nil && resp.RejectError == nil {
        // Success - update state
        q.state.OnProcessedBlock(req.Block.Hash())
        if req.ResultCh != nil {
            req.ResultCh <- SubmissionResult{Success: true}
        }
        return
    }
    
    // Handle "already processed"
    if isAlreadyProcessed(resp) {
        enclaveHead := resp.RejectError.L1Head
        q.state.OnProcessedBlock(enclaveHead)
        q.clearPendingBlocksUpTo(enclaveHead)
        if req.ResultCh != nil {
            req.ResultCh <- SubmissionResult{Skipped: true, NewHead: enclaveHead}
        }
        return
    }
    
    // Handle other errors
    q.logger.Error("Block submission failed", "error", err)
    if req.ResultCh != nil {
        req.ResultCh <- SubmissionResult{Error: err}
    }
}
```

### Producer Interface (Backward Compatible)

```go
// HandleBlock (from L1 stream) becomes a producer
func (g *Guardian) HandleBlock(block *types.Header) {
    if !g.running.Load() {
        return
    }
    
    g.state.OnReceivedBlock(block.Hash())
    
    if !g.state.InSyncWithL1() {
        return // Still use state check to avoid flooding queue
    }
    
    // Non-blocking enqueue with deduplication
    if _, exists := g.submissionQueue.pendingBlocks.LoadOrStore(block.Hash(), struct{}{}); exists {
        return // Already queued
    }
    
    select {
    case g.submissionQueue.l1Blocks <- &L1SubmissionRequest{
        Block:    block,
        IsLatest: true,
        Priority: PriorityLive,
    }:
        // Queued successfully
    default:
        // Queue full - drop and let catchup handle it
        g.submissionQueue.pendingBlocks.Delete(block.Hash())
        g.logger.Warn("Submission queue full, dropping live block", "block", block.Hash())
    }
}

// catchupWithL1 also becomes a producer
func (g *Guardian) catchupWithL1() error {
    for g.running.Load() && g.state.GetStatus() == L1Catchup {
        enclaveHead := g.state.GetEnclaveL1Head()
        l1Block, isLatest, err := g.sl.L1Data().FetchNextBlock(enclaveHead)
        if err != nil {
            return err
        }
        
        // Enqueue with result channel to wait for completion
        resultCh := make(chan SubmissionResult, 1)
        g.submissionQueue.l1Blocks <- &L1SubmissionRequest{
            Block:    l1Block,
            IsLatest: isLatest,
            Priority: PriorityCatchup,
            ResultCh: resultCh,
        }
        
        // Wait for result
        result := <-resultCh
        if result.Error != nil {
            return result.Error
        }
        if result.Skipped {
            // Fast-forward to new head
            g.state.OnProcessedBlock(result.NewHead)
        }
    }
    return nil
}
```

## Migration Path

### Stage 1: Add Queue (No Behavior Change)
- Implement queue infrastructure
- Route all submissions through queue
- Keep existing error handling logic
- Run in parallel with old code for comparison

### Stage 2: Simplify Logic
- Remove TryLock (no longer needed)
- Remove recursive calls (already done)
- Consolidate error handling in single consumer

### Stage 3: Optimize
- Add priority-based processing
- Implement smart deduplication
- Add metrics and monitoring

## Benefits

1. **Eliminates Race Conditions**: Single consumer owns all state updates
2. **Reduces Lock Contention**: No more TryLock/errEnclaveBusy churn
3. **Clearer Error Handling**: All error recovery in one place
4. **Better Observability**: Queue depth metrics show system health
5. **Graceful Degradation**: Queue full = automatic backpressure

## Risks & Mitigations

### Risk: Queue Overflow During Slow Blocks
**Mitigation**: 
- Bounded queue with sensible size (e.g., 100 blocks)
- Drop oldest catchup blocks when full (keep live blocks)
- Alert on sustained high queue depth

### Risk: Memory Leaks in Dedup Map
**Mitigation**:
- Periodic cleanup of old entries
- Size limits with LRU eviction
- Clear on fast-forward (clearPendingBlocksUpTo)

### Risk: Deadlock on Shutdown
**Mitigation**:
- Close queue channels in specific order
- Use context cancellation for all operations
- Timeout on drain operations

## Testing Strategy

1. **Unit Tests**: Queue operations, deduplication, error handling
2. **Integration Tests**: Multi-producer scenarios, state sync
3. **Chaos Tests**: Slow block processing, queue overflow, concurrent submissions
4. **Production Canary**: Deploy to single node first, monitor metrics

## Alternative: Keep Current Architecture

If queue approach seems too risky, current architecture can be improved with:
1. ✅ Remove recursive calls (done)
2. ✅ Use enclave L1Head for sync (done)
3. ✅ Proper lock handling (done)
4. Add retry limits to prevent infinite loops
5. Add circuit breaker for repeated failures
6. Better metrics around submission patterns

These changes are lower risk and might be sufficient.

