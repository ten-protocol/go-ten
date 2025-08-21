# Multisig Upgrade Processes

This document outlines two different approaches for implementing multisig governance for upgradeable L1/ L2 contracts in the Ten protocol.

## Current Architecture

- **TransparentUpgradeableProxy** - Enables upgradeability with built-in admin
- **Individual Proxy Admins** - Each proxy has its own admin who can upgrade it

## Option 1: Direct Multisig Ownership

**Used by:** Many smaller L2 protocols

### Architecture
```
Gnosis Safe (Multisig) → Proxy Admin → Implementation
```

### Setup Process

1. **Setup multisig wallet**
   ```bash
   # Use Gnosis Safe UI to create multisig
   # https://app.safe.global/
   # Recommended: 3-of-5 or 4-of-7 for initial phase
   ```

2. **Transfer proxy admin ownership directly to multisig**
   ```bash
   export MULTISIG_ADDR="0x..."
   export NETWORK_CONFIG_ADDR="0x..."
   npx hardhat run deployment_scripts/upgrade/001_direct_multisig_setup.ts --network mainnet
   ```

3. **Upgrade contracts**
   ```bash
   export MULTISIG_ADDR="0x..."
   export NETWORK_CONFIG_ADDR="0x..."
   npx hardhat run deployment_scripts/upgrade/002_direct_upgrade.ts --network mainnet
   ```

### Upgrade Process

1. **Deploy new implementation**
2. **Multisig directly calls proxy.upgradeTo()**


### Implementation Details

#### Phase 1A: Quick Setup Script
```bash
# 001_direct_multisig_setup.ts
# Transfers all proxy admin ownership directly to multisig
# No timelock, immediate control
```

#### Phase 1B: Direct Upgrade Script
```bash
# 002_direct_upgrade.ts
# Multisig can upgrade immediately without delays
# Perfect for critical fixes during initial release
```

---

## Option 2: Timelock

**Used by:** Uniswap, Compound, Aave, most major DeFi protocols

### Architecture
```
Gnosis Safe (Multisig) → Timelock → Proxy Admin → Implementation
```

### Setup Process

1. **Setup Gnosis Safe Multisig wallet**
   ```bash
   # Use Gnosis Safe UI to create multisig
   # https://app.safe.global/
   ```

2. **Deploy Timelock Controller**
   ```bash
   export MULTISIG_ADDR="0x..."
   npx hardhat run deployment_scripts/upgrade/001_deploy_timelock.ts --network mainnet
   ```

3. **Transfer proxy admin ownership to Timelock**
   ```bash
   export TIMELOCK_ADDRESS="0x..."
   export NETWORK_CONFIG_ADDR="0x..."
   npx hardhat run deployment_scripts/upgrade/002_transfer_proxy_admin.ts --network mainnet
   ```

### Upgrade Process

1. **Deploy new implementation**
2. **Multisig proposes upgrade via Timelock**
3. **24-hour delay period**
4. **Multisig executes upgrade**

### Pros
- Industry standard
- Time delay prevents rushed upgrades
- Transparent governance

### Cons
- More complex setup
- 24-hour delay for all upgrades (can be made configurable via env var)

---

## Transition Strategy: Option 1 → Option 2

1. **Deploy Timelock (while keeping direct control)**
   ```bash
   # Deploy timelock but don't transfer control yet
   npx hardhat run deployment_scripts/upgrade/001_deploy_timelock.ts --network mainnet
   ```

2. **Test Timelock functionality**
   ```bash
   # Test with non-critical operations
   # Verify delay mechanisms work correctly
   ```

3. **Transfer control to Timelock**
   ```bash
   # Transfer proxy admin ownership from multisig to timelock
   npx hardhat run deployment_scripts/upgrade/003_transition_to_timelock.ts --network mainnet
   ```

4. **Verify transition**
   ```bash
   # Confirm all proxies now controlled by timelock
   # Test upgrade process through timelock
   ```

### Transition Script: 003_transition_to_timelock.ts
```typescript
// This script transfers control from direct multisig to timelock
// It's the bridge between Option 1 and Option 2
// Only run after Option 1 has proven stable
```
---

## Security Considerations

1. **Multisig Composition**
   - Use 3-of-5 for quick decision making
   - Distribute keys geographically
   - Use hardware wallets where possible

2. **Upgrade Validation**
   - Have rollback plan ready
   - Test all upgrades on testnet first
   - Coordinate team availability for critical upgrades


2. **Timelock Delay**
   - 24 hours is standard
   - Can be adjusted based on risk tolerance

1. **Emergency Pause**
   - Add pause functionality to contracts
   - Multisig can pause immediately (no timelock)


3. **Rollback Plan**
   - Keep previous implementation addresses
   - Have rollback transaction ready

---

## References

- [OpenZeppelin Timelock](https://docs.openzeppelin.com/contracts/4.x/api/governance#TimelockController)
- [Gnosis Safe](https://docs.safe.global/)
- [Uniswap Governance](https://docs.uniswap.org/protocol/concepts/governance)
- [Compound Governance](https://docs.compound.finance/governance/) 