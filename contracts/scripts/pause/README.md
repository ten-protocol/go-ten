# Pause Scripts

This directory contains scripts for managing the pause functionality of contracts that implement `PausableWithRoles`.

## Scripts

### 1. `001_transfer_unpauser_role.ts`

**Purpose**: Transfers the `UNPAUSER_ROLE` from the deployer to a specified multisig address for all contracts that implement `PausableWithRoles`.

**Usage**:
```bash
npx hardhat run scripts/pause/001_transfer_unpauser_role.ts --network layer1
```

**Required Environment Variables**:
- `MULTISIG_ADDR`: The multisig wallet address to transfer the unpauser role to
- `NETWORK_CONFIG_ADDR`: The network config contract address
- `MERKLE_TREE_MESSAGE_BUS_ADDR`: The merkle message bus contract address

**What it does**:
1. Connects to all contracts that implement `PausableWithRoles`
2. Transfers the `UNPAUSER_ROLE` from the deployer to the specified multisig address
3. Reports the success/failure of each role transfer

### 2. `002_pause_all_contracts.ts`

**Purpose**: Pauses all contracts that implement `PausableWithRoles` and then verifies they are successfully paused.

**Usage**:
```bash
npx hardhat run scripts/pause/002_pause_all_contracts.ts --network layer1
```

**Required Environment Variables**:
- `NETWORK_CONFIG_ADDR`: The network config contract address
- `MERKLE_TREE_MESSAGE_BUS_ADDR`: The merkle message bus contract address

**What it does**:
1. Connects to all contracts that implement `PausableWithRoles`
2. Pauses each contract using the deployer account (which should have `PAUSER_ROLE`)
3. Verifies that each contract is successfully paused
4. Reports the status of each contract

## Contracts Affected

Both scripts work with the following contracts that implement `PausableWithRoles`:

- **MessageBus**: Cross-chain message handling contract
- **NetworkEnclaveRegistry**: Registry for network enclaves
- **DataAvailabilityRegistry**: Registry for data availability
- **TenBridge**: Bridge contract for token transfers
- **CrossChainMessenger**: Cross-chain messaging contract
- **MerkleTreeMessageBus**: Merkle tree-based message bus

## Workflow

The typical workflow for managing contract pauses is:

1. **Setup Phase**: Deploy contracts with deployer having both `PAUSER_ROLE` and `UNPAUSER_ROLE`
2. **Role Transfer Phase**: Use `001_transfer_unpauser_role.ts` to transfer `UNPAUSER_ROLE` to multisig
3. **Pause Phase**: Use `002_pause_all_contracts.ts` to pause all contracts when needed
4. **Recovery Phase**: Use multisig to unpause contracts when safe to do so

## Security Considerations

- The deployer retains `PAUSER_ROLE` for quick response to emergencies
- Only the multisig wallet can unpause contracts after the role transfer
- Pausing contracts stops all critical functionality
- Use with extreme caution on mainnet environments

## Go Launchers

For automated deployment and management, Go launchers are available in the `testnet/launcher/` directory:

- `transferunpauserroles/`: Launcher for transferring unpauser roles
- `pauseallcontracts/`: Launcher for pausing all contracts

These launchers can be used in CI/CD pipelines and automated workflows.
