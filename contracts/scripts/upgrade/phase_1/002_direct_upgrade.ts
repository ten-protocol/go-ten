import { ethers } from "hardhat";

/**
 * Direct Upgrade Script for Initial Mainnet Phase
 * 
 * This script allows the multisig to upgrade contracts immediately without any delays.
 * It's designed for the initial mainnet phase when rapid iteration is needed.
 * 
 * WARNING: This bypasses all timelock delays. Only use during initial mainnet
 * phase when the protocol is not public-facing.
 */

interface UpgradeConfig {
    contractName: string;
    proxyAddress: string;
    description: string;
}

interface DirectUpgradeResult {
    contractName: string;
    oldImplementation: string;
    newImplementation: string;
    upgradeTx: any;
    success: boolean;
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
 * Perform direct upgrade (no timelock, immediate execution)
 */
async function performDirectUpgrade(
    config: UpgradeConfig,
    newImplementation: string
): Promise<DirectUpgradeResult> {
    const { contractName, proxyAddress, description } = config;
    
    console.log(`\n=== Direct Upgrade: ${contractName} ===`);
    console.log("Proxy address:", proxyAddress);
    console.log("New implementation:", newImplementation);
    console.log("Description:", description);
    
    try {
        // Get the TransparentUpgradeableProxy contract
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
        const proxy = TransparentUpgradeableProxy.attach(proxyAddress);
        
        // Get current implementation before upgrade
        const oldImplementation = await (proxy as any).implementation();
        console.log("Current implementation:", oldImplementation);
        
        // Verify the multisig is the admin
        const currentAdmin = await (proxy as any).admin();
        const multisigAddress = process.env.MULTISIG_ADDRESS;
        
        if (currentAdmin.toLowerCase() !== multisigAddress?.toLowerCase()) {
            throw new Error(`Proxy admin is not the multisig. Current admin: ${currentAdmin}`);
        }
        
        console.log("Multisig control verified");
        
        // Perform the upgrade immediately
        console.log("Performing immediate upgrade...");
        const upgradeTx = await (proxy as any).upgradeTo(newImplementation);
        await upgradeTx.wait();
        
        console.log("Upgrade completed successfully!");
        console.log("Transaction hash:", upgradeTx.hash);
        
        // Verify the upgrade
        const newImpl = await (proxy as any).implementation();
        if (newImpl.toLowerCase() === newImplementation.toLowerCase()) {
            console.log("Implementation verification successful");
        } else {
            throw new Error("Implementation verification failed");
        }
        
        return {
            contractName,
            oldImplementation,
            newImplementation,
            upgradeTx,
            success: true
        };
        
    } catch (error) {
        console.error(`Failed to upgrade ${contractName}:`, error);
        return {
            contractName,
            oldImplementation: "0x...",
            newImplementation,
            upgradeTx: null,
            success: false
        };
    }
}

/**
 * Batch upgrade multiple contracts
 */
async function batchDirectUpgrade(
    configs: UpgradeConfig[]
): Promise<DirectUpgradeResult[]> {
    console.log("=== Batch Direct Upgrade Process ===");
    console.log(`Upgrading ${configs.length} contracts...\n`);
    
    const results: DirectUpgradeResult[] = [];
    
    for (const config of configs) {
        console.log(`\n--- Processing ${config.contractName} ---`);
        
        // Deploy new implementation
        const newImplementation = await deployNewImplementation(config.contractName);
        
        // Perform upgrade
        const result = await performDirectUpgrade(config, newImplementation);
        results.push(result);
        
        if (result.success) {
            console.log(`${config.contractName} upgrade completed`);
        } else {
            console.log(`${config.contractName} upgrade failed`);
        }
    }
    
    return results;
}

/**
 * Main upgrade function
 */
async function main() {
    console.log("=== Direct Upgrade Script for Initial Mainnet Phase ===\n");
    console.log("WARNING: This script bypasses all timelock delays!");
    console.log("Only use during initial mainnet phase for rapid iteration.\n");
    
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    
    // Configuration validation
    const multisigAddress = process.env.MULTISIG_ADDRESS || "0x...";
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR || "0x...";
    
    if (multisigAddress === "0x..." || networkConfigAddr === "0x...") {
        throw new Error('Please set MULTISIG_ADDRESS and NETWORK_CONFIG_ADDR environment variables');
    }
    
    console.log("Configuration:");
    console.log("- Multisig address:", multisigAddress);
    console.log("- NetworkConfig address:", networkConfigAddr);
    console.log("- Deployer address:", deployer.address);
    
    try {
        // Get addresses from network config
        const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
        const addresses = await networkConfig.addresses();
        
        console.log('\nCurrent proxy addresses:');
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
        
        console.log('\n=== Deploying New Implementations ===');
        const implementations: { [key: string]: string } = {};
        
        // Deploy all implementations first
        for (const config of upgradeConfigs) {
            console.log(`\n--- Deploying ${config.contractName} Implementation ---`);
            const newImplementation = await deployNewImplementation(config.contractName);
            implementations[config.contractName] = newImplementation;
        }
        
        console.log('\n=== SAFE TRANSACTION BUNDLE ===');
        console.log('Copy these transactions to your Safe UI for manual execution:\n');
        
        // Generate Safe transaction bundle
        for (const config of upgradeConfigs) {
            const newImplementation = implementations[config.contractName];
            
            console.log(`--- ${config.contractName} Upgrade Transaction ---`);
            console.log(`Description: ${config.description}`);
            console.log(`Contract Address: ${config.proxyAddress}`);
            console.log(`Method: upgradeTo(address)`);
            console.log(`Implementation Address: ${newImplementation}`);
            console.log(`Value: 0 ETH`);
            console.log('');
            
            // Generate calldata for verification
            const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
            const proxy = TransparentUpgradeableProxy.attach(config.proxyAddress);
            const calldata = proxy.interface.encodeFunctionData("upgradeTo", [newImplementation]);
            console.log(`Calldata: ${calldata}`);
            console.log('---\n');
        }
        
        
        console.log('=== IMPLEMENTATION ADDRESSES ===');
        console.log('Save these for verification:');
        for (const [contractName, address] of Object.entries(implementations)) {
            console.log(`${contractName}: ${address}`);
        }
        console.log('===============================\n');
        
    } catch (error) {
        console.error("Failed to perform direct upgrades:", error);
        throw error;
    }
}

// Run the main function if this script is executed directly
if (require.main === module) {
    main()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}

export {
    performDirectUpgrade,
    batchDirectUpgrade,
    deployNewImplementation
};
