# Rollback Scripts

This directory contains scripts for managing emergency situations and role transfers in the TEN protocol contracts.

## Scripts Overview

### 1. Emergency Pause Script (`001_emergency_pause.ts`)

**Purpose**: Immediately pause all contracts if a malicious upgrade is detected.

**Features**:
- Bypasses timelock for immediate response
- Pauses all core contracts (CrossChain, NetworkEnclaveRegistry, DataAvailabilityRegistry)
- Requires `EMERGENCY_PAUSER_ROLE` on contracts
- Can be used for both pause and unpause operations

**Usage**:
```bash
# Pause all contracts
npx hardhat run scripts/rollback/001_emergency_pause.ts

# Unpause all contracts
ACTION=unpause npx hardhat run scripts/rollback/001_emergency_pause.ts
```

**Environment Variables**:
- `EMERGENCY_PAUSER_ADDRESS`: Address with emergency pauser role
- `NETWORK_CONFIG_ADDR`: Address of NetworkConfig contract

### 2. Transfer Unpauser Roles Script (`002_transfer_unpauser_roles.ts`)

**Purpose**: Transfer `UNPAUSER_ROLE` from deployer to multisig wallet for all contracts.

**Features**:
- Transfers unpauser role for all contracts using `PausableWithRoles`
- Verifies role transfers were successful
- Includes verification mode to check current role assignments
- Supports all core contracts and bridge contracts

**Usage**:
```bash
# Transfer unpauser roles to multisig
npx hardhat run scripts/rollback/002_transfer_unpauser_roles.ts

# Verify current role assignments
ACTION=verify npx hardhat run scripts/rollback/002_transfer_unpauser_roles.ts
```

**Environment Variables**:
- `MULTISIG_ADDRESS`: Address of the multisig wallet
- `NETWORK_CONFIG_ADDR`: Address of NetworkConfig contract

## Contract Roles Overview

### PausableWithRoles System

All contracts now use a unified role-based pause system:

- **`PAUSER_ROLE`**: Can pause contracts (typically deployer for quick response)
- **`UNPAUSER_ROLE`**: Can unpause contracts (typically multisig for controlled recovery)
- **`DEFAULT_ADMIN_ROLE`**: Can grant/revoke roles and transfer unpauser role

### Contracts Using PausableWithRoles

1. **CrossChain** - Cross-chain value transfers and message verification
2. **NetworkEnclaveRegistry** - Network enclave management
3. **DataAvailabilityRegistry** - Rollup data availability
4. **MessageBus** - Cross-chain messaging infrastructure
4. **MerkleMessageBus** - Cross-chain messaging infrastructure
5. **TenBridge** - L1 bridge contract
6. **EthereumBridge** - L2 bridge contract

## Workflow

### 1. Initial Setup
- Deployer has both `PAUSER_ROLE` and `UNPAUSER_ROLE` on all contracts
- Deployer can pause and unpause contracts

### 2. Role Transfer
- Run `002_transfer_unpauser_roles.ts` to transfer `UNPAUSER_ROLE` to multisig
- Deployer retains `PAUSER_ROLE` for emergency response
- Multisig gets `UNPAUSER_ROLE` for controlled recovery

### 3. Emergency Response
- **Quick Pause**: Deployer uses `PAUSER_ROLE` to pause contracts immediately
- **Controlled Recovery**: Multisig uses `UNPAUSER_ROLE` to unpause when ready

### 4. Verification
- Use verification mode to check current role assignments
- Ensure multisig has `UNPAUSER_ROLE` on all contracts
- Verify deployer still has `PAUSER_ROLE` for emergencies

## Security Considerations

- **Emergency Pause**: Only use when malicious activity is detected
- **Role Management**: Regularly verify role assignments
- **Multisig Security**: Ensure multisig wallet is properly secured
- **Access Control**: Monitor who has access to pause/unpause functions

## Example Environment Setup

```bash
# .env file
MULTISIG_ADDRESS=0x1234567890123456789012345678901234567890
NETWORK_CONFIG_ADDR=0x0987654321098765432109876543210987654321
EMERGENCY_PAUSER_ADDRESS=0x1111111111111111111111111111111111111111
```

## Testing

Before running on mainnet:

1. **Test on local network** with hardhat
2. **Test on testnet** with small amounts
3. **Verify role transfers** work correctly
4. **Test pause/unpause** functionality
5. **Verify multisig** can unpause contracts

## Troubleshooting

### Common Issues

1. **"Contract does not have transferUnpauserRoleToMultisig function"**
   - Contract may not be using `PausableWithRoles`
   - Check contract inheritance

2. **"Deployer does not have UNPAUSER_ROLE"**
   - Role may have already been transferred
   - Check current role assignments

3. **"Transaction failed"**
   - Insufficient gas
   - Network congestion
   - Contract may be paused

### Verification Commands

```bash
# Check current roles
ACTION=verify npx hardhat run scripts/rollback/002_transfer_unpauser_roles.ts

# Check specific contract roles
npx hardhat console
> const contract = await ethers.getContractAt("ContractName", "address")
> const role = ethers.keccak256(ethers.toUtf8Bytes("UNPAUSER_ROLE"))
> await contract.hasRole(role, "address")
```