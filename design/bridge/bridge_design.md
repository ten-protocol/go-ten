# TEN Cross-Chain Messaging & Bridge

## 1. Introduction

TEN is an L2 that needs to communicate with Ethereum L1. This document describes the cross-chain messaging protocol and reference bridge built on top of it..

The messaging protocol is a general-purpose API. Any developer can use it to build cross-chain applications — bridges, oracles, governance relays, or anything else that requires authenticated data to move between L1 and L2.

The reference bridge is one such application. It locks ERC20 tokens (and native ETH) on L1 and mints wrapped equivalents on L2, and vice versa.

### Goals

- Provide a generic, permissionless messaging API between L1 and L2
- Build the bridge as an ordinary application on top of that API (no special privileges)
- Support both push and pull semantics for message delivery
- Ensure security through Merkle proofs, challenge periods, and replay protection

---

## 2. Architecture Overview

The system has three layers:
1. Transport - publishes and stores messages
2. Relay - verifies and delivers messages
3. Applications (eg Bridge) - dApps built by anyone


Each layer exists on both the L1 and L2, but with different implementations suited to each environment:

| Component | L1 | L2 |
|---|---|---|
| Transport | `MerkleTreeMessageBus` verifies messages via Merkle proofs against rollup state roots | `MessageBus` stores messages submitted by the enclave from L1 block processing |
| Relay | `CrossChainMessenger` relays messages using Merkle proofs | `CrossChainMessenger` relays messages after finality verification |
| App | `TenBridge` holds locked tokens, manages whitelist | `EthereumBridge` mints/burns wrapped tokens |

---

## 3. Core Data Structures

Before diving into contracts, here are the two message types that flow through the system.

### CrossChainMessage

Every cross-chain communication is encoded as a `CrossChainMessage`:

```solidity
struct CrossChainMessage {
    address sender;           // address that published the message
    uint64  sequence;         // auto-incrementing index per sender
    uint64  nonce;            // app-defined deduplication key
    uint32  topic;            // message category (TRANSFER, MANAGEMENT, VALUE)
    bytes   payload;          // arbitrary encoded data
    uint8   consistencyLevel; // block confirmations before finality
}
```

- `sender` + `sequence` uniquely identify a message globally
- `nonce` lets applications group related messages or deduplicate identical payloads
- `topic` provides basic routing and versioning — the bridge uses topics `0` (transfer), `1` (management), and `2` (value)
- `consistencyLevel` controls how many L1 block confirmations must pass before the message is considered final

### ValueTransferMessage

A specialised struct for native currency transfers:

```solidity
struct ValueTransferMessage {
    address sender;
    address receiver;
    uint256 amount;
    uint64  sequence;
}
```

This exists because native value transfers don't go through the messenger relay — they're verified directly via Merkle proofs on L1.

---

## 4. MessageBus (Transport)

The MessageBus is the foundation. It does two things:
1. Publish — accept messages from local contracts and emit them as events
2. Store and verify — accept messages from the other chain and make them queryable

### Publishing a Message

Any contract calls `publishMessage()`:

```solidity
function publishMessage(
    uint64 nonce,
    uint32 topic,
    bytes  payload,
    uint8  consistencyLevel
) external payable returns (uint64 sequence);
```

The MessageBus:
1. Collects the publishing fee (`msg.value`)
2. Assigns a sequence number to the sender
3. Emits `LogMessagePublished(sender, sequence, nonce, topic, payload, consistencyLevel)`

Messages are not stored on-chain at the source. They exist only as events. The enclave (for L1 to L2) or the rollup publisher (for L2 to L1) picks them up from there.

### Storing Incoming Messages (L2 only)

On L2, when the enclave processes an L1 block and finds `LogMessagePublished` events, it creates "synthetic transactions" that call `storeCrossChainMessage` on the L2 MessageBus:

```solidity
function storeCrossChainMessage(
    CrossChainMessage message,
    uint256 finalAfterTimestamp
) external;
```

The function is gated by `ownerOrSelf` — synthetic transactions are sent from a masked version of the MessageBus address, which the protocol treats as a trusted internal caller. User transactions to this function are blocked at the protocol level, even if sent from an address with the correct role.

The message is stored in a nested mapping (`sender to topic to message[]`) and its hash is recorded with a finality timestamp. The message cannot be consumed until `block.timestamp >= finalAfterTimestamp`.

### Verifying a Message

```solidity
function verifyMessageFinalized(CrossChainMessage message) external view returns (bool);
```

Returns `true` if the message has been stored **and** its challenge period has elapsed.

### L1 Extension: MerkleTreeMessageBus

On L1, the MessageBus is extended with Merkle proof verification. Instead of storing every L2 to L1 message individually, rollup state roots are registered and messages are verified against them:

```solidity
function addStateRoot(bytes32 stateRoot, uint256 activationTime) external;
function verifyMessageInclusion(CrossChainMessage msg, bytes32[] proof, bytes32 root) external view returns (bool);
function verifyValueTransferInclusion(ValueTransferMessage msg, bytes32[] proof, bytes32 root) external view returns (bool);
```

**Leaf construction for messages:**
```
leaf = keccak256(abi.encode("m", keccak256(abi.encode(message))))
```

**Leaf construction for value transfers:**
```
leaf = keccak256(abi.encode("v", keccak256(abi.encode(message))))
```

State roots have an activation time — they can't be used for verification until `block.timestamp >= activationTime`. This is the challenge period. State roots can also be disabled if a rollup is successfully challenged.

---

## 5. Cross Chain Messenger 

The messenger sits between the MessageBus and application contracts. It takes a verified message and executes it — calling the target contract with the encoded calldata.

### Cross Chain Function Calls

When a dApp wants to send a cross-chain function call, it encodes a `CrossChainCall`:

```solidity
struct CrossChainCall {
    address target;  // contract to call on destination chain
    bytes   data;    // encoded function call
    uint256 gas;     // gas allocation (currently unused)
}
```

This struct is packed into the `payload` field of a `CrossChainMessage`.

### Relaying a Message

Anyone can trigger message relay — it's permissionless. The messenger exposes two paths:

**Direct verification** (L2, where messages are stored by the enclave):
```solidity
function relayMessage(CrossChainMessage message) external;
```

**Merkle proof verification** (L1, where messages are verified against state roots):
```solidity
function relayMessageWithProof(
    CrossChainMessage message,
    bytes32[] proof,
    bytes32 root
) external;
```

In both cases the messenger:
1. Verifies the message 
2. Marks it as consumed (prevents replay)
3. Decodes the `CrossChainCall` from the payload
4. Sets `crossChainSender` and `crossChainTarget` temporarily
5. Calls `target.call(data)`
6. Clears the sender/target state

### Replay Protection

A `mapping(bytes32 => bool) messageConsumed` ensures each message can only be relayed once. The key is the hash of the full `CrossChainMessage`.

---

## 6. Developer API (CrossChainEnabledTEN)

Application contracts inherit `CrossChainEnabledTEN` to get cross-chain capabilities. This abstract contract provides:

### Sending Messages

```solidity
function queueMessage(
    address remoteTarget,
    bytes memory payload,
    uint32 topic,
    uint256 gas,
    uint8 consistencyLevel
) internal;
```

This encodes a `CrossChainCall` and publishes it through the MessageBus. The bridge uses this to send function calls like `receiveAssets(...)` to its counterpart on the other chain.

For raw data (not function calls):
```solidity
function publishRawMessage(bytes memory data, uint32 topic, uint8 consistencyLevel) internal;
```

### Receiving Messages

The `onlyCrossChainSender` modifier verifies that the current call originated from a cross-chain relay and came from the expected sender:

```solidity
modifier onlyCrossChainSender(address expectedSender) {
    require(_isCrossChain());
    require(_crossChainSender() == expectedSender);
    require(_crossChainTarget() == address(this));
    _;
}
```

This is how the bridge ensures that only its counterpart on the other chain can call `receiveAssets`.

---

## 7. Message Flows

There are two types of cross-chain communication, and they work differently:

- **Message-based** (`queueMessage`) — encodes a function call as a `CrossChainCall`. Requires explicit relay on the destination chain.
- **Value transfers** (`publishRawMessage` with `Topics.VALUE`) — encodes a `ValueTransfer` struct. On L2, the enclave processes these automatically as native balance increases. No relay needed.

### L1 to L2 (Message-based)

Used by: ERC20 deposits, token whitelisting, WETH notifications.

```
  Ethereum L1                                TEN L2

  1. User/admin calls dApp (e.g. sendERC20)
  2. dApp calls queueMessage()
     to MessageBus.publishMessage()
  3. LogMessagePublished event emitted
                    │
                    │  (automatic)
                    │  Enclave processes L1 block,
                    │  extracts LogMessagePublished events,
                    │  creates synthetic transactions
                    ▼
                                             4. MessageBus.storeCrossChainMessage()
                                                (message now stored on L2)
                                             5. Poll MessageBus.verifyMessageFinalized()
                                                until it returns true
                                             6. Call CrossChainMessenger.relayMessage(message)
                                             7. Messenger decodes CrossChainCall,
                                                calls target contract
```

Step 4 is automatic (the enclave creates synthetic transactions). Steps 5-6 require someone to act — typically the user's frontend extracts the `CrossChainMessage` from the L1 `LogMessagePublished` event logs, polls for finality on L2, then submits the relay transaction.

### L1 to L2 (Native value transfer)

Used by `TenBridge.sendNative()` (native ETH deposits).

```
  Ethereum L1                                TEN L2


  1. User calls TenBridge.sendNative{value}()
  2. publishRawMessage(ValueTransfer, Topics.VALUE)
     to MessageBus.publishMessage()
  3. LogMessagePublished event emitted
                    │
                    │  
                    │ 
                    │  
                    │  
                    ▼
                                             4. Receiver's balance updated
```

Native value transfers bypass the relay mechanism entirely. The enclave recognises `Topics.VALUE` messages with a `ValueTransfer` payload and credits the receiver's L2 balance directly via state modification.

### L2 to L1

Used by ERC20 withdrawals, native withdrawals, WETH withdrawals.

```
  TEN L2                                     Ethereum L1
  ──────                                     ──────────

  1. User calls dApp (e.g. sendERC20)
  2. dApp calls queueMessage()
     to MessageBus.publishMessage()
  3. LogMessagePublished event emitted
  4. Message included in rollup state
                    │
                    │  
                    │ 
                    │  
                    │  
                    │  
                    ▼
                                             5. Challenge period elapses,
                                                state root activates
                                             6. Call CrossChainMessenger
                                                .relayMessageWithProof(
                                                    message, proof, root)
                                             7. MerkleTreeMessageBus verifies
                                                proof against state root
                                             8. Messenger decodes CrossChainCall,
                                                calls target contract
```

On L1, messages from L2 are never stored individually. The rollup's state root is registered in the `MerkleTreeMessageBus`, and messages are verified on-demand using Merkle proofs. This is gas-efficient — only messages that are actually claimed incur verification cost. The caller must provide the Merkle proof and the state root to use.

---

## 8. Reference Bridge

The bridge is built entirely on top of the messaging API. It consists of two contracts:

| Contract | Chain | Role |
|---|---|---|
| `TenBridge` | L1 | Holds locked tokens, manages whitelist |
| `EthereumBridge` | L2 | Mints/burns wrapped tokens |


### Token Whitelisting

Before a token can be bridged, it must be whitelisted on L1. This triggers wrapped token creation on L2.

**On L1 (admin action):**
1. Admin calls `TenBridge.whitelistToken(tokenAddress, name, symbol)`
2. Bridge grants `ERC20_TOKEN_ROLE` to the token address
3. Bridge calls `queueMessage()` to encodes `onCreateTokenCommand(tokenAddress, name, symbol)` as a `CrossChainCall` to `MessageBus.publishMessage()` emits `LogMessagePublished`

**Relay to L2 (user/relayer action):**
4. Extract the `CrossChainMessage` from the L1 tx's `LogMessagePublished` event
5. Poll L2 `MessageBus.verifyMessageFinalized(message)` until it returns `true`
6. Call L2 `CrossChainMessenger.relayMessage(message)`
7. Messenger decodes the `CrossChainCall` and calls `EthereumBridge.onCreateTokenCommand()`
8. `EthereumBridge` deploys a new `WrappedERC20` and stores the bidirectional token mapping

### Deposit: L1 to L2 (ERC20)

**On L1 (user action):**
1. User calls `TenBridge.sendERC20(asset, amount, receiver)`
2. Bridge verifies token is whitelisted and not suspended
3. Bridge calls `safeTransferFrom(user, bridge, amount)` — locks tokens
4. Bridge calls `queueMessage()` to encodes `receiveAssets(asset, amount, receiver)` to `MessageBus.publishMessage()`

**Relay to L2 (user/relayer action):**
5. Extract `CrossChainMessage` from L1 tx logs
6. Poll L2 `MessageBus.verifyMessageFinalized(message)` until `true`
7. Call L2 `CrossChainMessenger.relayMessage(message)`
8. Messenger calls `EthereumBridge.receiveAssets(asset, amount, receiver)`
9. Bridge looks up the L2 wrapped token and calls `WrappedERC20.issueFor(receiver, amount)`

### Withdrawal: L2 to L1 (ERC20)

**On L2 (user action):**
1. User calls `EthereumBridge.sendERC20(wrappedToken, amount, receiver)` with `msg.value >= publishFee`
2. Bridge looks up L1 token address for the wrapped token
3. Bridge calls `WrappedERC20.burnFor(user, amount)` — burns wrapped tokens
4. Bridge calls `queueMessage()` to encodes `receiveAssets(l1Token, amount, receiver)` to `MessageBus.publishMessage()`

**Relay to L1 (user/relayer action, after challenge period):**
5. Wait for the rollup containing this message to be published to L1 and the state root to activate
6. Call L1 `CrossChainMessenger.relayMessageWithProof(message, proof, root)`
7. Messenger verifies Merkle proof against the state root
8. Messenger calls `TenBridge.receiveAssets(l1Token, amount, receiver)`
9. Bridge calls `safeTransfer(receiver, amount)` — releases locked tokens

### Native ETH: L1 to L2

**On L1 (user action):**
1. User calls `TenBridge.sendNative{value: amount}(receiver)`
2. Bridge calls `publishRawMessage()` with `ValueTransfer(amount, receiver)` on `Topics.VALUE`

**Automatic (no relay needed):**
3. Enclave processes the value transfer and directly increases the receiver's native L2 balance

### Native: L2 to L1

**On L2 (user action):**
1. User calls `EthereumBridge.sendNative{value: amount + fee}(receiver)`
2. Bridge deducts the publish fee
3. Bridge calls `queueMessage()` to encodes `receiveAssets(address(0), amount, receiver)` on `Topics.VALUE`

**Relay to L1 (after challenge period):**
4. Call L1 `CrossChainMessenger.relayMessageWithProof(message, proof, root)`
5. Messenger calls `TenBridge.receiveAssets(address(0), amount, receiver)`
6. Bridge sends native ETH to receiver

### WETH: L1 to L2

WETH deposits involve two operations — a native value transfer (automatic) and a message (needs relay):

**On L1 (user action):**
1. User calls `TenBridge.sendERC20(weth, amount, receiver)`
2. Bridge unwraps WETH to native ETH via `IWETH.withdraw(amount)`
3. Bridge calls `this.sendNative{value: amount}(remoteBridgeAddress)` — publishes a `ValueTransfer` that increases the L2 bridge's native balance
4. Bridge calls `queueMessage()` to encodes `receiveNativeWrapped(receiver, amount)` on `Topics.TRANSFER`

**Automatic + relay:**
5. Enclave credits the L2 bridge contract's native balance (automatic)
6. User/relayer polls finality, then calls L2 `CrossChainMessenger.relayMessage(message)`
7. Messenger calls `EthereumBridge.receiveNativeWrapped(receiver, amount)`
8. Bridge wraps native ETH to WETH via `IWETH.deposit()` and transfers WETH to receiver

### WETH: L2 to L1

**On L2 (user action):**
1. User calls `EthereumBridge.sendERC20(localWETH, amount, receiver)` with `msg.value >= 2 * publishFee`
2. Bridge collects WETH from user via `safeTransferFrom`, then unwraps to native ETH
3. Bridge calls `this.sendNative{value: amount + fee}(remoteBridgeAddress)` — queues `receiveAssets(address(0), amount, remoteBridge)` on `Topics.VALUE`
4. Bridge calls `queueMessage()` to encodes `receiveNativeWrapped(receiver, amount)` on `Topics.TRANSFER`

**Relay to L1 (after challenge period):**
5. Relay the value transfer message — `TenBridge.receiveAssets(address(0), ...)` sends native ETH to the bridge itself
6. Relay the WETH notification — `TenBridge.receiveNativeWrapped(receiver, amount)` wraps ETH to WETH and transfers to receiver

---

## 9. Security Model

Security is layered, matching the architecture.

### MessageBus Security

| Threat | Mitigation |
|---|---|
| **Fake messages from L2** | L1 verifies messages against rollup state roots via Merkle proofs. Fake messages would require forging a state root, which requires compromising the rollup. |
| **Premature message consumption** | Challenge period enforced via `activationTime` on state roots and `finalAfterTimestamp` on stored messages. |
| **L1 block reorganisation** | Rollups are bound to specific L1 block hashes. If the L1 reorgs, the rollup referencing the old block becomes invalid, and any messages it contained are rejected. |
| **Fake messages to L2** | Only synthetic transactions from the enclave can call `storeCrossChainMessage`. User transactions to the MessageBus are blocked at the protocol level. |

### Messenger Security

| Threat | Mitigation |
|---|---|
| **Message replay** | `messageConsumed` mapping prevents any message from being relayed twice. |
| **Reentrancy** | `ReentrancyGuard` on all relay functions. |
| **Sender spoofing** | `onlyCrossChainSender` modifier checks both the original sender and intended target. |

### Bridge Security

| Threat | Mitigation |
|---|---|
| **Unauthorised minting** | Only cross-chain messages from the paired bridge contract can trigger `receiveAssets`. |
| **Double withdrawal** | Handled by messenger-level replay protection. Each message can only be consumed once. |
| **Malicious tokens** | Whitelist controlled by admin (eventually DAO). Only approved tokens can be bridged. |
| **Reentrancy on withdrawals** | `ReentrancyGuardTransient` on `TenBridge`. |


### Contract Dependency Chain

```
TenBridge → CrossChainMessenger → MerkleTreeMessageBus ← DataAvailabilityRegistry
                                                        ← CrossChain
```

- `TenBridge` inherits `CrossChainEnabledTEN`, which is configured with the `CrossChainMessenger` address
- `CrossChainMessenger` is initialized with the `MerkleTreeMessageBus` address for proof verification
- `DataAvailabilityRegistry` publishes state roots to `MerkleTreeMessageBus`
- `CrossChain` manages value transfer withdrawals via `MerkleTreeMessageBus`

---


## 12. Contract Reference

### Messaging Contracts

| Contract | Purpose |
|---|---|
| `IMessageBus` | MessageBus interface |
| `MessageBus` | Base implementation — publish, store, verify |
| `MerkleTreeMessageBus` | L1 extension with Merkle proof and state root verification |
| `ICrossChainMessenger` | Messenger interface |
| `CrossChainMessenger` | Relay layer — consumes and executes messages |
| `CrossChainEnabledTEN` | Abstract base for dApps — provides `queueMessage`, `onlyCrossChainSender` |

### Bridge Contracts

| Contract | Purpose |
|---|---|
| `IBridge` | Bridge interface |
| `TenBridge` | L1 bridge — locks tokens, manages whitelist |
| `EthereumBridge` | L2 bridge — mints/burns wrapped tokens |
| `WrappedERC20` | Wrapped token with mint/burn access control |
| `TenERC20` | Privacy-aware ERC20 base for TEN |

### L1 Management Contracts

| Contract | Purpose |
|---|---|
| `NetworkConfig` | Registry of all deployed contract addresses |
| `CrossChain` | Manages cross-chain value transfer withdrawals and bundle verification |
| `DataAvailabilityRegistry` | Publishes rollup state roots to `MerkleTreeMessageBus` |
| `NetworkEnclaveRegistry` | Manages enclave registration |
| `Fees` | Configurable fee parameters for message publishing |


# Addition stuff(not sure if we want to include)

### L1 Reorganisation Handling

This is the most subtle security concern. Consider:

1. User deposits tokens on L1
2. Message published, enclave stores it on L2, user relays it, tokens minted on L2
3. User withdraws back to L1
4. **L1 reorgs** — the original deposit transaction disappears

The system handles this because rollups are **bound to specific L1 block hashes**. If the L1 reorgs:
- The block hash changes for that block number
- The rollup that processed the now-missing deposit references the old (invalid) block hash
- The L1 rollup contract rejects the rollup
- The enclave regenerates from the new canonical chain

The `consistencyLevel` parameter provides additional protection — setting it to `N` means the message won't be processed until `N` blocks have confirmed the original transaction, exponentially reducing reorg probability.

---

## 10. Fees

Publishing a message on L2 has a direct cost: the sequencer must pay L1 gas to include it in the rollup. The fee model channels this cost to the user.

```solidity
function getPublishFee() external view returns (uint256);
```

The fee is collected as `msg.value` when calling `publishMessage`. It covers:
- Fixed cost for the message metadata
- Variable cost proportional to payload size

Fee parameters are managed by a separate `Fees` contract, configurable by the DAO.

For the bridge specifically, native value transfers deduct the publishing fee from the sent amount:
```
actual_bridged = msg.value - publishFee
```

---

## Deployment

All L1 contracts are deployed behind `OpenZeppelinTransparentProxy` for upgradeability. L2 contracts (`MessageBus`, `CrossChainMessenger`, `EthereumBridge`) are deployed as **system contracts** during L2 network creation — they are not deployed via the L1 deployment scripts.

### L1 Deployment Order

The L1 deployment happens across three scripts in `deployment_scripts/core/`:

**Step 1** Core infrastructure:
1. `Fees` — fee configuration
2. `MerkleTreeMessageBus` — initialized with `(deployer, deployer, feesAddress)`
3. `CrossChain` — initialized with `(deployer, merkleMessageBusAddress)`, manages value transfer withdrawals
4. `NetworkEnclaveRegistry` — enclave registration
5. `DataAvailabilityRegistry` — initialized with `(merkleMessageBusAddress, networkEnclaveRegistryAddress, deployer)`
6. `NetworkConfig` — registry of all deployed addresses
7. Role grants:
   - `DataAvailabilityRegistry` → `STATE_ROOT_MANAGER_ROLE` on `MerkleTreeMessageBus` (so it can publish rollup state roots)
   - `CrossChain` → `WITHDRAWAL_MANAGER_ROLE` on `MerkleTreeMessageBus`

**Step 2** Relay layer:
1. `CrossChainMessenger` — initialized with `messageBusAddress` (read from `NetworkConfig`)
2. Address recorded in `NetworkConfig` via `setL1CrossChainMessengerAddress`

**Step 3** Bridge:
1. `TenBridge` — initialized with `(crossChainMessengerAddress, deployer)`
2. Address recorded in `NetworkConfig` via `setL1BridgeAddress`

### L1 to L2 Linking

After both L1 and L2 contracts exist, `bridge/001_deploy_bridge.ts` links them:
1. Reads L1 and L2 bridge addresses from network config
2. Calls `TenBridge.setRemoteBridge(l2BridgeAddress)` on L1
