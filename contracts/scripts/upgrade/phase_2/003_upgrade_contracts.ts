import { BaseContract } from 'ethers';
import { ethers } from 'hardhat';
import { upgrades } from 'hardhat';
import { UpgradeOptions } from '@openzeppelin/hardhat-upgrades/dist/utils';
import * as path from 'path';
const hre = require("hardhat");

console.log('=== Multisig Upgrade Script started ===');

interface UpgradeConfig {
    contractName: string;
    proxyAddress: string;
    description: string;
}

interface GovernanceConfig {
    timelockAddress: string;
    multisigAddress: string;
}

/**
 * Deploy new implementation for a contract
 */
async function deployNewImplementation(contractName: string): Promise<string> {
    console.log(`Deploying new ${contractName} implementation...`);
    
    const factory = await ethers.getContractFactory(contractName);
    const implementation = await factory.deploy();
    await implementation.waitForDeployment();
    
    const address = await implementation.getAddress();
    console.log(`${contractName} implementation deployed:`, address);
    return address;
}

/**
 * Prepare upgrade transaction for multisig governance
 * Each TransparentUpgradeableProxy has its own admin who can upgrade it
 */
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
    
    console.log(`\n=== Upgrade Transaction for ${contractName} ===`);
    console.log("Target:", proxyAddress);
    console.log("Value: 0");
    console.log("Calldata:", upgradeCalldata);
    console.log("Description:", description);
    console.log("New Implementation:", newImplementation);
    console.log("==========================================\n");
    
    return {
        target: proxyAddress,
        value: 0,
        calldata: upgradeCalldata,
        description
    };
}

/**
 * Schedule upgrade through timelock
 */
async function scheduleTimelockUpgrade(
    timelockAddress: string,
    upgradeTx: any,
    delay: number = 24 * 60 * 60 // 24 hours
) {
    // Get TimelockController contract using the same approach as existing code
    const TimelockController = await ethers.getContractFactory("TimelockController");
    const timelock = TimelockController.attach(timelockAddress) as any;
    
    console.log("Scheduling upgrade through timelock...");
    
    // Schedule the upgrade using the correct interface
    const scheduleTx = await timelock.schedule(
        upgradeTx.target,
        upgradeTx.value,
        upgradeTx.calldata,
        ethers.ZeroHash, // predecessor
        ethers.ZeroHash, // salt
        delay
    );
    
    console.log("Upgrade scheduled. Transaction hash:", scheduleTx.hash);
    console.log(`Upgrade will be executable after ${delay} seconds (${delay / 3600} hours)`);
    
    return scheduleTx;
}

/**
 * Execute upgrade through timelock (after delay period)
 */
async function executeTimelockUpgrade(timelockAddress: string, upgradeTx: any) {
    // Get TimelockController contract using the same approach as existing code
    const TimelockController = await ethers.getContractFactory("TimelockController");
    const timelock = TimelockController.attach(timelockAddress) as any;
    
    console.log("Executing upgrade...");
    
    const executeTx = await timelock.execute(
        upgradeTx.target,
        upgradeTx.value,
        upgradeTx.calldata,
        ethers.ZeroHash, // predecessor
        ethers.ZeroHash  // salt
    );
    
    console.log("Upgrade executed. Transaction hash:", executeTx.hash);
    return executeTx;
}

/**
 * Generate execution commands for a scheduled upgrade
 */
function generateExecutionCommands(contractName: string, timelockAddress: string, upgradeTx: any) {
    console.log(`\n=== Execution Commands for ${contractName} ===`);
    console.log("# Execute upgrade after delay period");
    console.log(`export TIMELOCK_ADDRESS="${timelockAddress}"`);
    console.log(`export UPGRADE_TARGET="${upgradeTx.target}"`);
    console.log(`export UPGRADE_CALLDATA="${upgradeTx.calldata}"`);
    console.log(`export UPGRADE_VALUE="${upgradeTx.value}"`);
    console.log(`npx hardhat run scripts/upgrade/001_upgrade_contracts.ts --network mainnet`);
    console.log("================================================");
}

/**
 * Upgrade contract using multisig governance
 */
async function upgradeContractWithMultisig(
    contractName: string,
    proxyAddress: string,
    description: string,
    governanceConfig: GovernanceConfig,
    scheduleOnly: boolean = true
): Promise<any> {
    console.log(
        `Preparing ${contractName} upgrade for multisig governance (proxy: ${proxyAddress})`
    );
    
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

/**
 * Main upgrade function using multisig governance
 */
const upgradeContractsWithMultisig = async function (): Promise<void> {
    console.log('=== Starting multisig upgrade process ===');
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    
    // Governance configuration - these should be set as environment variables
    const governanceConfig: GovernanceConfig = {
        timelockAddress: process.env.TIMELOCK_ADDRESS || "0x...",
        multisigAddress: process.env.MULTISIG_ADDRESS || "0x..."
    };
    
    // Validate governance config
    if (governanceConfig.timelockAddress === "0x..." || 
        governanceConfig.multisigAddress === "0x...") {
        throw new Error('Please set TIMELOCK_ADDRESS and MULTISIG_ADDRESS environment variables');
    }
    
    console.log('Governance Configuration:');
    console.table(governanceConfig);
    
    // Get addresses from network config
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    if (!networkConfigAddr) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable is not set');
    }

    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const addresses = await networkConfig.addresses();

    console.log('\nCurrent proxy addresses');
    console.table({
        NetworkConfig: networkConfigAddr,
        CrossChain: addresses.crossChain,
        NetworkEnclaveRegistry: addresses.networkEnclaveRegistry,
        DataAvailabilityRegistry: addresses.dataAvailabilityRegistry
    });

    // Define upgrade configurations
    const upgradeConfigs: UpgradeConfig[] = [
        {
            contractName: 'CrossChain',
            proxyAddress: addresses.crossChain,
            description: 'Upgrade CrossChain contract to v2.0.0 with improved security'
        },
        {
            contractName: 'NetworkEnclaveRegistry',
            proxyAddress: addresses.networkEnclaveRegistry,
            description: 'Upgrade NetworkEnclaveRegistry contract to v2.0.0'
        },
        {
            contractName: 'DataAvailabilityRegistry',
            proxyAddress: addresses.dataAvailabilityRegistry,
            description: 'Upgrade DataAvailabilityRegistry contract to v2.0.0'
        }
    ];

    // Perform upgrades
    const upgradeResults = [];
    for (const config of upgradeConfigs) {
        console.log(`\n=== Processing ${config.contractName} upgrade ===`);
        const result = await upgradeContractWithMultisig(
            config.contractName,
            config.proxyAddress,
            config.description,
            governanceConfig,
            true // scheduleOnly = true
        );
        upgradeResults.push(result);
    }

    console.log('\n=== Upgrade Summary ===');
    console.log('All upgrades have been scheduled through timelock.');
    console.log('Next steps:');
    console.log('1. Wait for 24-hour delay period');
    console.log('2. Execute upgrades through multisig');
    console.log('3. Verify upgrades on blockchain');
    
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
    
    // Print cancellation commands
    console.log('\n=== Cancellation Commands (if needed) ===');
    for (let i = 0; i < upgradeResults.length; i++) {
        const result = upgradeResults[i];
        const config = upgradeConfigs[i];
        if (config && result) {
            generateExecutionCommands(config.contractName, governanceConfig.timelockAddress, result.upgradeTx);
        }
    }
    
    console.log('\nAll upgrades scheduled successfully');
}

/**
 * Execute previously scheduled upgrades (run after delay period)
 */
const executeScheduledUpgrades = async function (): Promise<void> {
    console.log('=== Executing scheduled upgrades ===');
    
    const governanceConfig: GovernanceConfig = {
        timelockAddress: process.env.TIMELOCK_ADDRESS || "0x...",
        multisigAddress: process.env.MULTISIG_ADDRESS || "0x..."
    };
    
    // This would be called after the delay period
    // Implementation would depend on how you want to handle the execution
    console.log('Execute upgrades through multisig after delay period');
}

// Export functions for use in other scripts
export {
    upgradeContractWithMultisig,
    scheduleTimelockUpgrade,
    executeTimelockUpgrade,
    prepareUpgradeTransaction,
    executeScheduledUpgrades,
    generateExecutionCommands
};

// Run the main function if this script is executed directly
if (require.main === module) {
    upgradeContractsWithMultisig()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}

export default upgradeContractsWithMultisig;