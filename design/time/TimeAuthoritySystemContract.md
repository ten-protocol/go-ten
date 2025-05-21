# L2 Sequencer Based Time Authority System Contract

In standard EVM environments, time is determined by `block.timestamp` (alias `block.time`), which represents 
the timestamp of the block being processed. This creates limitations when precise timing is needed between blocks.

For example, if block N has timestamp 100 and block N+1 has timestamp 120, there is no way to reference
or trigger actions at timestamp 110, even though logically events may need to occur at that time.
This forces developers to stretch timing-sensitive logic to align with block boundaries.

The Time Authority System Contract provides a solution by allowing the L2 sequencer to maintain its own
time authority. This enables:

1. More granular time tracking between L2 blocks
2. Deterministic time progression independent of block timestamps 
3. The ability to reference and trigger actions at precise moments between blocks
4. Better support for time-sensitive applications like auctions, gaming, and DeFi


## Transaction Time Authority

In addition to block-level timing, the Time Authority also tracks the precise arrival time of each transaction in the mempool. This "transaction time" represents the canonical moment when a transaction is considered to have occurred, regardless of when it is ultimately included in a block.

### Transaction Arrival Time

When a transaction is received by the sequencer's RPC endpoint, it is immediately timestamped with the current Time Authority time. This timestamp becomes an immutable part of the transaction's metadata and represents its official "arrival time" in the system.

### Expanded Mempool

To support transaction timing, the mempool is expanded to maintain additional metadata for each transaction:
- Timestamp of receival in the RPC. 
This of course means there is a delay between transaction being sent and being received, and the action taking place at a different time, but this is inevitable flaw of all online systems that rely on time.


### Utilizing the transaction metadata

The Time Authority system contract is predeployed at genesis and plays a central role in managing transaction timing:

1. During batch preparation:
   - For each transaction targeting the Time Authority contract
   - The sequencer retrieves the transaction's arrival timestamp from its metadata
   - A synthetic transaction is created to store the tx.hash -> timestamp mapping on-chain
   - This synthetic transaction is included directly in the batch transactions
   - Unlike other synthetic transactions, these timing records are part of the main batch

2. During batch processing:
   - When a transaction targets the Time Authority contract
   - The Time Authority first stores the transaction's timestamp in a transient storage variable
   - The transaction is then forwarded to the target contract via:
     ```solidity
     targetContract.call{value: tx.value}(calldata)
     ```
   - The called contract can access its transaction's timestamp by calling:
     ```solidity
     timeAuthority.getTransactionTimestamp()
     ```

This mechanism ensures that contracts can deterministically access the canonical arrival time of their triggering transaction, enabling precise time-based logic that operates at transaction-level granularity rather than being constrained to block boundaries.


## Contract Interface

The Time Authority system contract exposes the following interface:

```solidity
interface ITimeAuthority {
    /**
     * @notice Relays a transaction with its timestamp to the target contract
     * @param target The address of the contract to call
     * @param data The calldata to send to the target contract
     */
    function relayTimestampedTransaction(address target, bytes calldata data) external payable;

    /**
     * @notice Sets timestamps for a batch of transactions
     * @param txHashes Array of transaction hashes
     * @param timestamps Array of corresponding timestamps
     * @dev This function is called once per batch at the beginning and includes only
     * information to be consumed within the same batch. Old records can be safely deleted
     * as the next call is guaranteed to be at a point where the data is no longer needed.
     */
    function setTimestampsForTransactions(bytes32[] calldata txHashes, uint256[] calldata timestamps) external;

    /**
     * @notice Gets the timestamp of the current transaction
     * @return The timestamp when the transaction was received by the sequencer
     */
    function getTransactionTimestamp() external view returns (uint256);
}
```

The `setTimestampsForTransactions` function is called once per batch at the beginning and includes only information to be consumed within the same batch. Old records can be safely deleted as the next call is guaranteed to be at a point where the data is no longer needed.






