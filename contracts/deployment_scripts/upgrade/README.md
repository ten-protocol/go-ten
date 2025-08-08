# Complete Multisig Setup & Upgrade Process for Timelock

This document provides a step-by-step walkthrough of the complete multisig governance setup and upgrade process for Ten Protocol's L1 contracts.

## Architecture

```
Gnosis Safe (Multisig) → TimelockController → Proxy Admin → Implementation
```

## Step-by-Step Process

### Step 1: Create Gnosis Safe Multisig Wallet

1. **Create New Safe**
   - Click "Create Safe"
   - Choose network (Ethereum mainnet/testnet)
   - Add owners (3-5 team members)
   - Set threshold (e.g., 3-of-5)
   - Review and deploy

3. **Save the Safe Address**
   ```bash
   # Example: 0x1234567890123456789012345678901234567890
   export MULTISIG_ADDRESS="0x1234567890123456789012345678901234567890"
   ```

### Step 2: Deploy Governance Infrastructure

**Run the deployment script:**
```bash
# Set your multisig address
export MULTISIG_ADDRESS="0x1234567890123456789012345678901234567890"

# Deploy governance infrastructure
npx hardhat run deployment_scripts/upgrade/001_deploy_timelock.ts --network mainnet
```

**What this script does:**
1. **Deploys TimelockController** - Creates a time-delayed governance contract
2. **Grants roles** - Gives multisig proposer and executor permissions
3. **Revokes deployer admin** - Removes deployer's admin role for security
4. **Prints environment variables** - Shows what to set for next steps

**Output will show:**
```
=== Setting up complete governance system ===

Deploying TimelockController...
Deployer address: 0x...
Configuration:
- Multisig address: 0x1234567890123456789012345678901234567890
- Delay (seconds): 86400
- Delay (hours): 24

TimelockController deployed successfully!
Timelock address: 0x9876543210987654321098765432109876543210
Proposer role successfully granted to multisig
Executor role successfully granted to multisig
Admin role successfully revoked from deployer

=== Deployment Summary ===
TimelockController: 0x9876543210987654321098765432109876543210
Multisig (proposer/executor): 0x1234567890123456789012345678901234567890
Delay: 86400 seconds ( 24 hours)
==========================

=== Governance Setup Complete ===
Environment variables to set:
export TIMELOCK_ADDRESS="0x9876543210987654321098765432109876543210"
export MULTISIG_ADDRESS="0x1234567890123456789012345678901234567890"
===============================

Next Steps:
1. Transfer proxy admin ownership to Timelock for each proxy:
   - CrossChain proxy admin → Timelock
   - NetworkEnclaveRegistry proxy admin → Timelock
   - DataAvailabilityRegistry proxy admin → Timelock
2. Use the upgrade script to schedule upgrades through timelock
```

### Step 3: Set Environment Variables

After deployment, you'll get addresses like this:
```bash
export TIMELOCK_ADDRESS="0x9876543210987654321098765432109876543210"
export MULTISIG_ADDRESS="0x1234567890123456789012345678901234567890"
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"
```

### Step 4: Transfer Proxy Admin Ownership

**Run the transfer script:**
```bash
# Set environment variables
export TIMELOCK_ADDRESS="0x9876543210987654321098765432109876543210"
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"

# Transfer admin ownership to timelock
npx hardhat run deployment_scripts/upgrade/002_transfer_proxy_admin.ts --network mainnet
```

**What this script does:**
1. **Gets proxy addresses** - Reads current proxy addresses from NetworkConfig
2. **Checks current admin** - Verifies who currently owns each proxy
3. **Transfers ownership** - Moves admin ownership from deployer to timelock
4. **Validates transfers** - Ensures all transfers completed successfully

**Output will show:**
```
Transferring proxy admin ownership to Timelock...
Deployer address: 0x...
Configuration:
- Timelock address: 0x9876543210987654321098765432109876543210
- NetworkConfig address: 0x1111111111111111111111111111111111111111

Current proxy addresses:
┌─────────────────────┬──────────────────────────────────────┐
│ NetworkConfig       │ 0x1111111111111111111111111111111111 │
│ CrossChain          │ 0x2222222222222222222222222222222222 │
│ NetworkEnclaveRegistry│ 0x33333333333333333333333333333333 │
│ DataAvailabilityRegistry│ 0x44444444444444444444444444444444 │
└─────────────────────┴──────────────────────────────────────┘

=== Transferring CrossChain proxy admin ownership ===
Current admin: 0x...
Transferring admin ownership from 0x... to 0x9876543210987654321098765432109876543210...
CrossChain proxy admin ownership transferred successfully!
Transaction hash: 0x...

=== Transferring NetworkEnclaveRegistry proxy admin ownership ===
Current admin: 0x...
Transferring admin ownership from 0x... to 0x9876543210987654321098765432109876543210...
NetworkEnclaveRegistry proxy admin ownership transferred successfully!
Transaction hash: 0x...

=== Transferring DataAvailabilityRegistry proxy admin ownership ===
Current admin: 0x...
Transferring admin ownership from 0x... to 0x9876543210987654321098765432109876543210...
DataAvailabilityRegistry proxy admin ownership transferred successfully!
Transaction hash: 0x...

=== Transfer Summary ===
All proxy admin ownership transfers completed
Timelock now controls all proxy upgrades
==========================
```

### Step 5: Understanding the Upgrade Script

The upgrade script (`scripts/upgrade/001_upgrade_contracts.ts`) orchestrates the entire upgrade process. Here's what each part does:

#### Script Structure & Functions

##### 1. Interfaces
```typescript
interface UpgradeConfig {
    contractName: string;    // Name of contract to upgrade (e.g., "CrossChain")
    proxyAddress: string;    // Proxy address of the contract
    description: string;     // Human-readable description of the upgrade
}

interface GovernanceConfig {
    timelockAddress: string;     // TimelockController contract address
    multisigAddress: string;     // Gnosis Safe multisig address
}
```

##### 2. `deployNewImplementation(contractName)`
```typescript
async function deployNewImplementation(contractName: string): Promise<string> {
    console.log(`Deploying new ${contractName} implementation...`);
    
    const factory = await ethers.getContractFactory(contractName);
    const implementation = await factory.deploy();
    await implementation.waitForDeployment();
    
    const address = await implementation.getAddress();
    console.log(`${contractName} implementation deployed:`, address);
    return address;
}
```
**What it does:**
- Compiles and deploys a new implementation of the specified contract
- Returns the address of the new implementation
- This is the new version that will replace the old one

##### 3. `prepareUpgradeTransaction(config, newImplementation)`
```typescript
async function prepareUpgradeTransaction(
    config: UpgradeConfig,
    newImplementation: string
) {
    const { contractName, proxyAddress, description } = config;
    
    // Get the TransparentUpgradeableProxy contract
    const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
    const proxy = TransparentUpgradeableProxy.attach(proxyAddress);
    
    // Prepare upgrade calldata - the admin calls upgradeTo on the proxy directly
    const upgradeCalldata = proxy.interface.encodeFunctionData(
        "upgradeTo",
        [newImplementation]
    );
    
    return {
        target: proxyAddress,
        value: 0,
        calldata: upgradeCalldata,
        description
    };
}
```
**What it does:**
- Creates the transaction data needed to upgrade a proxy
- Encodes the `upgradeTo` function call for the proxy
- Returns a transaction object ready for timelock scheduling

##### 4. `scheduleTimelockUpgrade(timelockAddress, upgradeTx, delay)`
```typescript
async function scheduleTimelockUpgrade(
    timelockAddress: string,
    upgradeTx: any,
    delay: number = 24 * 60 * 60 // 24 hours
) {
    const TimelockController = await ethers.getContractFactory("TimelockController");
    const timelock = TimelockController.attach(timelockAddress) as any;
    
    console.log("Scheduling upgrade through timelock...");
    
    const scheduleTx = await timelock.schedule(
        upgradeTx.target,      // Proxy address
        upgradeTx.value,       // 0 (no ETH sent)
        upgradeTx.calldata,    // Encoded upgrade transaction
        ethers.ZeroHash,       // predecessor (none)
        ethers.ZeroHash,       // salt (none)
        delay                  // 24 hours
    );
    
    return scheduleTx;
}
```
**What it does:**
- Schedules the upgrade transaction through the timelock
- Creates a 24-hour delay before the upgrade can be executed
- Returns the scheduling transaction

##### 5. `executeTimelockUpgrade(timelockAddress, upgradeTx)`
```typescript
async function executeTimelockUpgrade(timelockAddress: string, upgradeTx: any) {
    const TimelockController = await ethers.getContractFactory("TimelockController");
    const timelock = TimelockController.attach(timelockAddress) as any;
    
    console.log("Executing upgrade...");
    
    const executeTx = await timelock.execute(
        upgradeTx.target,
        upgradeTx.value,
        upgradeTx.calldata,
        ethers.ZeroHash,
        ethers.ZeroHash
    );
    
    return executeTx;
}
```
**What it does:**
- Executes the previously scheduled upgrade (after 24-hour delay)
- Actually performs the upgrade on the blockchain
- Returns the execution transaction

##### 6. `upgradeContractWithMultisig(contractName, proxyAddress, description, governanceConfig, scheduleOnly)`
```typescript
async function upgradeContractWithMultisig(
    contractName: string,
    proxyAddress: string,
    description: string,
    governanceConfig: GovernanceConfig,
    scheduleOnly: boolean = true
): Promise<any> {
    // Get current implementation address
    const currentImpl = await hre.upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log(`Current implementation address: ${currentImpl}`);
    
    // Deploy new implementation
    const newImplementation = await deployNewImplementation(contractName);
    
    // Prepare upgrade transaction
    const upgradeTx = await prepareUpgradeTransaction(
        { contractName, proxyAddress, description },
        newImplementation
    );
    
    if (scheduleOnly) {
        // Schedule the upgrade (this would be done by multisig)
        const scheduleTx = await scheduleTimelockUpgrade(governanceConfig.timelockAddress, upgradeTx);
        return { scheduleTx, upgradeTx, newImplementation };
    } else {
        // Execute the upgrade (after delay period)
        const executeTx = await executeTimelockUpgrade(governanceConfig.timelockAddress, upgradeTx);
        return { executeTx, upgradeTx, newImplementation };
    }
}
```
**What it does:**
- Orchestrates the entire upgrade process for a single contract
- Deploys new implementation → Prepares transaction → Schedules/Executes
- Returns all relevant transaction data

##### 7. `upgradeContractsWithMultisig()` (Main Function)
```typescript
const upgradeContractsWithMultisig = async function (): Promise<void> {
    // Validate environment variables
    const governanceConfig: GovernanceConfig = {
        timelockAddress: process.env.TIMELOCK_ADDRESS || "0x...",
        multisigAddress: process.env.MULTISIG_ADDRESS || "0x..."
    };
    
    // Get current proxy addresses from NetworkConfig
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const addresses = await networkConfig.addresses();
    
    // Define upgrade configurations
    const upgradeConfigs: UpgradeConfig[] = [
        {
            contractName: 'CrossChain',
            proxyAddress: addresses.crossChain,
            description: 'Upgrade CrossChain contract to v2.0.0 with improved security'
        },
        // ... more contracts
    ];
    
    // Perform upgrades
    const upgradeResults = [];
    for (const config of upgradeConfigs) {
        const result = await upgradeContractWithMultisig(
            config.contractName,
            config.proxyAddress,
            config.description,
            governanceConfig,
            true // scheduleOnly = true
        );
        upgradeResults.push(result);
    }
    
    // Print execution commands for later use
    console.log('\n=== Execution Commands (run after delay) ===');
    for (let i = 0; i < upgradeResults.length; i++) {
        const result = upgradeResults[i];
        const config = upgradeConfigs[i];
        if (config && result) {
            console.log(`${config.contractName}:`);
            console.log(`  Timelock: ${governanceConfig.timelockAddress}`);
            console.log(`  Target: ${result.upgradeTx.target}`);
            console.log(`  Calldata: ${result.upgradeTx.calldata}`);
        }
    }
}
```
**What it does:**
- Main orchestrator for the entire upgrade process
- Validates environment variables
- Gets current contract addresses from NetworkConfig
- Processes all contracts in the upgrade configuration
- Prints execution commands for later use

### Step 6: Running the Upgrade Process

#### Phase 1: Schedule Upgrades
```bash
# Run the upgrade script to schedule upgrades
npx hardhat run scripts/upgrade/001_upgrade_contracts.ts --network mainnet
```

**Output will show:**
```
=== Starting multisig upgrade process ===
Governance Configuration:
┌─────────────────┬──────────────────────────────────────┐
│ timelockAddress │ 0x9876543210987654321098765432109876 │
│ multisigAddress │ 0x1234567890123456789012345678901234 │
└─────────────────┴──────────────────────────────────────┘

Current proxy addresses
┌─────────────────────┬──────────────────────────────────────┐
│ NetworkConfig       │ 0x1111111111111111111111111111111111 │
│ CrossChain          │ 0x2222222222222222222222222222222222 │
│ NetworkEnclaveRegistry│ 0x33333333333333333333333333333333 │
│ DataAvailabilityRegistry│ 0x44444444444444444444444444444444 │
└─────────────────────┴──────────────────────────────────────┘

=== Processing CrossChain upgrade ===
Deploying new CrossChain implementation...
CrossChain implementation deployed: 0x5555555555555555555555555555555555555555

=== Upgrade Transaction for CrossChain ===
Target: 0x2222222222222222222222222222222222222222
Value: 0
Calldata: 0x3659cfe60000000000000000000000005555555555555555555555555555555555555555
Description: Upgrade CrossChain contract to v2.0.0 with improved security
New Implementation: 0x5555555555555555555555555555555555555555
==========================================

Scheduling upgrade through timelock...
Upgrade scheduled. Transaction hash: 0x6666666666666666666666666666666666666666666666666666666666666666
Upgrade will be executable after 86400 seconds (24 hours)

=== Upgrade Summary ===
All upgrades have been scheduled through timelock.
Next steps:
1. Wait for 24-hour delay period
2. Execute upgrades through multisig
3. Verify upgrades on blockchain

=== Execution Commands (run after delay) ===
CrossChain:
  Timelock: 0x9876543210987654321098765432109876543210
  Target: 0x2222222222222222222222222222222222222222
  Calldata: 0x3659cfe60000000000000000000000005555555555555555555555555555555555555555
```

#### Phase 2: Execute Upgrades (After 24 Hours)

**Option A: Through Gnosis Safe UI**
1. Go to your Gnosis Safe
2. Create new transaction
3. Target: `0x9876543210987654321098765432109876543210` (Timelock)
4. Value: `0`
5. Data: `0x3659cfe60000000000000000000000005555555555555555555555555555555555555555`

**Option B: Through Script (if you have multisig keys)**
```bash
# After 24 hours, execute the upgrades
npx hardhat run scripts/upgrade/001_upgrade_contracts.ts --network mainnet
```

### Step 7: Verification

After execution, verify the upgrades:
```bash
# Check new implementation addresses
npx hardhat console --network mainnet
> const networkConfig = await ethers.getContractAt('NetworkConfig', '0x1111111111111111111111111111111111111111')
> const addresses = await networkConfig.addresses()
> console.log('CrossChain proxy:', addresses.crossChain)
> const impl = await hre.upgrades.erc1967.getImplementationAddress(addresses.crossChain)
> console.log('New implementation:', impl)
```

## Key Benefits of This Process

1. **Security**: 24-hour delay prevents rushed upgrades
2. **Multisig**: Multiple approvals required
3. **Transparency**: All upgrades are visible and scheduled
4. **Rollback**: Previous implementations are preserved
5. **Automation**: Script handles the entire process

## Script Files

- `deployment_scripts/upgrade/001_deploy_timelock.ts` - Deploy timelock script
- `deployment_scripts/upgrade/002_transfer_proxy_admin.ts` - Transfer proxy admin ownership
- `scripts/upgrade/001_upgrade_contracts.ts` - Upgrade workflow
- `scripts/emergency/001_emergency_pause.ts` - Emergency pause all contracts

## Documentation

- `contracts/docs/ROLLBACK_PROCEDURES_COMPLETE.md` - Complete rollback procedures and emergency response documentation

This gives you a **production-ready, secure upgrade system** that follows industry best practices! 