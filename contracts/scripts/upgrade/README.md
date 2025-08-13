# Multisig Upgrade System

This directory contains the complete upgrade system for the Ten protocol, implementing a **phase-based approach** that starts with direct multisig control and transitions to timelock governance.

## Overview

The upgrade system follows this progression:

1. **Phase 1 (Initial Mainnet)**: Direct multisig control for rapid iteration
2. **Phase 2 (Public Production)**: Timelock governance for transparent, secure upgrades

## Architecture

```
Phase 1: Deployer → Multisig → Proxy Admin → Implementation
Phase 2: Multisig → Timelock → Proxy Admin → Implementation
```

## Scripts

### Phase 1: Direct Multisig Control

#### `001_direct_multisig_setup.ts`
Sets up direct multisig control by transferring proxy admin ownership directly to the multisig.

```bash
export MULTISIG_ADDRESS="0x..."
export NETWORK_CONFIG_ADDR="0x..."
npx hardhat run deployment_scripts/upgrade/001_direct_multisig_setup.ts --network mainnet
```

**What it does:**
- Transfers all proxy admin ownership to multisig
- Enables immediate upgrades (no delays)
- Perfect for initial mainnet phase

#### `002_direct_upgrade.ts`
Performs direct upgrades without any timelock delays.

```bash
export MULTISIG_ADDRESS="0x..."
export NETWORK_CONFIG_ADDR="0x..."
npx hardhat run deployment_scripts/upgrade/002_direct_upgrade.ts --network mainnet
```

**What it does:**
- Deploys new implementations
- Upgrades contracts immediately

### Phase 2: Timelock Governance

#### `001_deploy_timelock.ts`
Deploys the TimelockController contract for transparent governance.

```bash
export MULTISIG_ADDRESS="0x..."
npx hardhat run scripts/upgrade/phase_2/001_deploy_timelock.ts --network mainnet
```

**What it does:**
- Deploys OpenZeppelin TimelockController
- Grants proposer/executor roles to multisig
- Sets up 24-hour delay (configurable)


#### `002_transition_to_timelock.ts` (Phase-Based Approach)
Bridges Phase 1 to Phase 2 by transferring control from direct multisig to timelock.

```bash
export TIMELOCK_ADDRESS="0x..."
export MULTISIG_ADDRESS="0x..."
export NETWORK_CONFIG_ADDR="0x..."
npx hardhat run scripts/upgrade/phase_2/003_transition_to_timelock.ts --network mainnet
```

**What it does:**
- Verifies current state
- Transfers proxy admin from multisig to timelock
- Tests timelock functionality

#### `003_upgrade_contracts.ts`
Performs timelock-based upgrades with 24-hour delay protection.

```bash
export TIMELOCK_ADDRESS="0x..."
export MULTISIG_ADDRESS="0x..."
export NETWORK_CONFIG_ADDR="0x..."
npx hardhat run scripts/upgrade/phase_2/004_upgrade_contracts.ts --network mainnet
```

**What it does:**
- Schedules upgrades through timelock
- Enforces 24-hour delay for all upgrades
- Provides execution commands for after delay period

## Deployment Flow

### Initial Mainnet Setup (Phase 1)

1. **Create Gnosis Safe Multisig**
   ```bash
   # Use https://app.safe.global/
   # Recommended: 3-of-5 for quick decision making
   ```

2. **Setup Direct Multisig Control**
   ```bash
   export MULTISIG_ADDRESS="0x..."
   export NETWORK_CONFIG_ADDR="0x..."
   npx hardhat run deployment_scripts/upgrade/001_direct_multisig_setup.ts --network mainnet
   ```

3. **Test Direct Upgrades**
   ```bash
   npx hardhat run deployment_scripts/upgrade/002_direct_upgrade.ts --network mainnet
   ```

### Transition to Timelock (Phase 2)

1. **Deploy Timelock (while keeping direct control)**
   ```bash
   export MULTISIG_ADDRESS="0x..."
   npx hardhat run scripts/upgrade/phase_2/001_deploy_timelock.ts --network mainnet
   ```

2. **Test Timelock functionality**
   ```bash
   # Verify roles and delay mechanisms
   # Test with non-critical operations
   ```

3. **Transition to Timelock Governance**
   ```bash
   export TIMELOCK_ADDRESS="0x..."
   export MULTISIG_ADDRESS="0x..."
   export NETWORK_CONFIG_ADDR="0x..."
   npx hardhat run scripts/upgrade/phase_2/002_transition_to_timelock.ts --network mainnet
   ```

4. **Verify Transition**
   ```bash
   # Confirm all proxies under timelock control
   # Test upgrade process through timelock
   ```

## Environment Variables

### Required for All Scripts
- `NETWORK_CONFIG_ADDR`: Address of the NetworkConfig contract

### Required for Phase 1
- `MULTISIG_ADDRESS`: Address of the Gnosis Safe multisig

### Required for Phase 2
- `MULTISIG_ADDRESS`: Address of the Gnosis Safe multisig
- `TIMELOCK_ADDRESS`: Address of the TimelockController contract

### Optional
- `TIMELOCK_DELAY`: Delay in seconds (default: 24 hours = 86400 seconds)

## Security Considerations

## Emergency Procedures

### Phase 1
- **Immediate pause**: Multisig can pause contracts directly
- **Quick rollback**: Deploy previous implementation and upgrade immediately
- **Team coordination**: All members must be available for critical decisions

### Phase 2
- **Emergency pause**: Multisig can pause through timelock (24-hour delay)
- **Emergency bypass**: Consider shorter timelock for critical fixes
- **Rollback plan**: Keep previous implementations ready

## Monitoring and Alerting

### Phase 1
- Monitor multisig transactions
- Alert on any admin ownership changes
- Track implementation upgrades

### Phase 2
- Monitor timelock operations
- Alert on scheduled upgrades
- Track execution of scheduled operations

## Testing

### Testnet Testing
1. Deploy contracts to testnet
2. Setup multisig governance
3. Test upgrade processes
4. Verify all functionality

### Mainnet Testing
1. Start with Phase 1 (direct multisig)
2. Test upgrade process thoroughly
3. Deploy timelock and test functionality
4. Transition to Phase 2
5. Verify complete governance system

## Troubleshooting

### Common Issues

1. **Proxy admin not deployer**
   - Check current admin ownership
   - Manual intervention may be required

2. **Multisig not authorized**
   - Verify multisig address
   - Check role assignments

3. **Timelock not configured**
   - Verify timelock deployment
   - Check role assignments
   - Verify delay settings

### Recovery Procedures

1. **Failed upgrade**
   - Deploy previous implementation
   - Upgrade to previous version
   - Investigate failure cause

2. **Governance issues**
   - Check multisig composition
   - Verify role assignments
   - Review transaction logs

## References

- [OpenZeppelin Timelock](https://docs.openzeppelin.com/contracts/4.x/api/governance#TimelockController)
- [Gnosis Safe](https://docs.safe.global/)
- [Hardhat Upgrades](https://docs.openzeppelin.com/upgrades-plugins/1.x/)
- [Proxy Patterns](https://docs.openzeppelin.com/upgrades-plugins/1.x/proxies)
