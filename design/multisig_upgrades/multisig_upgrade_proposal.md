# Multisig Upgrade Processess

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

1. **Transfer proxy admin ownership directly to multisig**
   ```bash
   # Transfer each proxy's admin to multisig address
   await proxy.changeAdmin(safeAddress);
   ```

### Upgrade Process

1. **Deploy new implementation**
2. **Multisig directly calls proxy.upgradeTo()**

### Pros
- Simple setup
- No delay for upgrades

### Cons
- No time delay protection
- Risk of rushed upgrades

---

## Option 2: Timelock (Recommended)

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
   export MULTISIG_ADDRESS="0x..."
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

## Recommended Implementation: Option 2

The reasons for this:

1. **Security**: The 24-hour delay provides protection against rushed upgrades
2. **Industry Standard**: Used by major DeFi protocols and proven governance model
3. **Transparency**: All upgrades are visible and scheduled

### Security Considerations

1. **Multisig Composition**
   - Use 3-of-5 or 4-of-7 for good security
   - Distribute keys geographically
   - Use hardware wallets where possible

2. **Timelock Delay**
   - 24 hours is standard
   - Can be adjusted based on risk tolerance
   - Consider emergency bypass mechanisms

3. **Upgrade Validation**
   - Have rollback plan ready

### Emergency Procedures

1. **Emergency Pause**
   - Add pause functionality to contracts
   - Multisig can pause immediately (no timelock)

2. **Emergency Upgrade**
   - Consider shorter timelock for critical fixes
   - Or direct multisig upgrade for emergencies

3. **Rollback Plan**
   - Keep previous implementation addresses
   - Have rollback transaction ready

---

## Next Steps

1. Choose governance model
2. Set up Gnosis Safe with team members
3. Deploy Timelock and transfer proxy admin ownership
4. Test upgrade process on testnet
5. Document emergency procedures
6. Deploy to mainnet

## References

- [OpenZeppelin Timelock](https://docs.openzeppelin.com/contracts/4.x/api/governance#TimelockController)
- [Gnosis Safe](https://docs.safe.global/)
- [Uniswap Governance](https://docs.uniswap.org/protocol/concepts/governance)
- [Compound Governance](https://docs.compound.finance/governance/) 