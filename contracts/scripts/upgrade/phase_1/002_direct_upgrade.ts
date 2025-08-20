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
    success: boolean;
}

interface SafeTransaction {
    to: string;
    value: string;
    data: string;
    operation: number;
    safeTxGas: string;
    baseGas: string;
    gasPrice: string;
    gasToken: string;
    refundReceiver: string;
    nonce: number;
}

interface SafeTransactionBundle {
    version: string;
    chainId: number;
    createdAt: string;
    meta: {
        name: string;
        description: string;
        txBuilderVersion: string;
        createdFromSafeAddress: string;
        createdFromOwnerAddress: string;
        checksums: { [key: string]: string };
    };
    transactions: Array<{
        to: string;
        value: string;
        data: string;
        contractMethod: {
            inputs: Array<{
                name: string;
                type: string;
                internalType: string;
            }>;
            name: string;
            payable: boolean;
        };
        contractInputsValues: { [key: string]: string };
    }>;
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
 * Generate Safe transaction JSON for a contract upgrade
 */
function generateSafeTransaction(
    proxyAddress: string,
    newImplementation: string,
    contractName: string
): SafeTransaction {
    // Get the TransparentUpgradeableProxy ABI for the upgradeTo method
    const upgradeToData = new ethers.Interface([
        "function upgradeTo(address newImplementation)"
    ]).encodeFunctionData("upgradeTo", [newImplementation]);
    
    return {
        to: proxyAddress,
        value: "0",
        data: upgradeToData,
        operation: 0, // 0 = call, 1 = delegatecall
        safeTxGas: "0", // safe estimate
        baseGas: "0", // safe estimate
        gasPrice: "0", //safe estimate
        gasToken: ethers.ZeroAddress,
        refundReceiver: ethers.ZeroAddress,
        nonce: 0 // set by Safe
    };
}

/**
 * Generate Safe transaction bundle JSON for batch upgrade
 */
function generateSafeTransactionBundle(
    transactions: Array<{ proxyAddress: string; newImplementation: string; contractName: string }>,
    multisigAddress: string,
    chainId: number
): SafeTransactionBundle {
    const now = new Date().toISOString();
    
    return {
        version: "1.0",
        chainId: chainId,
        createdAt: now,
        meta: {
            name: "Contract Upgrade Bundle",
            description: "Batch upgrade of TEN protocol contracts",
            txBuilderVersion: "1.0.0",
            createdFromSafeAddress: multisigAddress,
            createdFromOwnerAddress: multisigAddress,
            checksums: {}
        },
        transactions: transactions.map(({ proxyAddress, newImplementation, contractName }) => ({
            to: proxyAddress,
            value: "0",
            data: new ethers.Interface([
                "function upgradeTo(address newImplementation)"
            ]).encodeFunctionData("upgradeTo", [newImplementation]),
            contractMethod: {
                inputs: [
                    {
                        name: "newImplementation",
                        type: "address",
                        internalType: "address"
                    }
                ],
                name: "upgradeTo",
                payable: false
            },
            contractInputsValues: {
                newImplementation: newImplementation
            }
        }))
    };
}

/**
 * Verify that the multisig wallet controls all contracts
 */
async function verifyMultisigOwnership(
    multisigAddress: string,
    networkConfigAddr: string
): Promise<boolean> {
    console.log("Checking contract ownership...");
    
    try {
        const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
        const addresses = await networkConfig.addresses();
        
        const contracts = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        let allControlled = true;
        
        for (const contract of contracts) {
            try {
                const contractInstance = await ethers.getContractAt(contract.name, contract.address);
                const currentOwner = await (contractInstance as any).owner();
                
                if (currentOwner.toLowerCase() === multisigAddress.toLowerCase()) {
                    console.log(`${contract.name}: Controlled by Multisig`);
                } else {
                    console.log(`${contract.name}: NOT controlled by Multisig (Current: ${currentOwner})`);
                    allControlled = false;
                }
            } catch (error) {
                console.log(`Error checking ${contract.name}:`, error);
                allControlled = false;
            }
        }
        
        if (allControlled) {
            console.log("\nAll contracts are under Multisig control!");
            console.log("Direct upgrades are now possible (no delays)");
        } else {
            console.log("\nSome contracts are not under Multisig control");
            console.log("Complete the ownership transfer process first");
        }
        
        return allControlled;
        
    } catch (error) {
        console.error("Failed to verify multisig ownership:", error);
        return false;
    }
}

/**
 * Print JSON to console for easy copying
 */
async function printJsonToConsole(filename: string, data: any): Promise<void> {
    console.log(`\n=== ${filename.toUpperCase()} ===`);
    console.log('Copy this JSON:');
    console.log('='.repeat(60));
    console.log(JSON.stringify(data, null, 2));
    console.log('='.repeat(60));
    console.log(`End of ${filename}\n`);
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
            success: true
        };
        
    } catch (error) {
        console.error(`Failed to upgrade ${contractName}:`, error);
        return {
            contractName,
            oldImplementation: "0x...",
            newImplementation,
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
    
    // Verify multisig ownership before proceeding
    console.log("\n=== Verifying Multisig Ownership ===");
    const ownershipVerified = await verifyMultisigOwnership(multisigAddress, networkConfigAddr);
    
            if (!ownershipVerified) {
            console.error("\nOwnership verification failed!");
            console.error("The multisig wallet does not control all contracts.");
            console.error("Please complete the ownership transfer process first:");
            console.error("1. Run the multisig setup script (001_direct_multisig_setup.ts)");
            console.error("2. Accept ownership in your Gnosis Safe UI");
            console.error("3. Run this script again to verify ownership");
            process.exit(1);
        }
        
        console.log("Ownership verification passed! Proceeding with upgrades...\n");
    
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
        
        console.log('\n=== GENERATING SAFE TRANSACTION FILES ===');
        
        // Get chain ID for the transaction bundle
        const chainId = (await ethers.provider.getNetwork()).chainId;
        console.log(`Chain ID: ${chainId}`);
        
        // Generate individual transaction files for each contract
        for (const config of upgradeConfigs) {
            const newImplementation = implementations[config.contractName];
            if (!newImplementation) {
                console.error(`No implementation found for ${config.contractName}`);
                continue;
            }
            
            console.log(`\n--- Generating ${config.contractName} Transaction ---`);
            
            // Generate Safe transaction
            const safeTx = generateSafeTransaction(
                config.proxyAddress,
                newImplementation,
                config.contractName
            );
            
            // Print individual transaction JSON to console
            const filename = `${config.contractName}_upgrade_tx.json`;
            await printJsonToConsole(filename, safeTx);
            
            // Print transaction details
            console.log(`Contract: ${config.contractName}`);
            console.log(`Proxy Address: ${config.proxyAddress}`);
            console.log(`New Implementation: ${newImplementation}`);
            console.log(`Calldata: ${safeTx.data}`);
            console.log(`File: ${filename}`);
        }
        
        // Generate batch transaction bundle
        console.log('\n--- Generating Batch Transaction Bundle ---');
        const transactionData = upgradeConfigs
            .map(config => {
                const newImplementation = implementations[config.contractName];
                if (!newImplementation) return null;
                return {
                    proxyAddress: config.proxyAddress,
                    newImplementation: newImplementation,
                    contractName: config.contractName
                };
            })
            .filter((item): item is NonNullable<typeof item> => item !== null);
        
        if (!multisigAddress) {
            throw new Error('MULTISIG_ADDRESS environment variable is required');
        }
        
        const batchBundle = generateSafeTransactionBundle(
            transactionData,
            multisigAddress,
            Number(chainId)
        );
        
        // Print batch bundle JSON to console
        await printJsonToConsole('batch_upgrade_bundle.json', batchBundle);
        
        console.log('\n=== IMPLEMENTATION ADDRESSES ===');
        console.log('Save these for verification:');
        for (const [contractName, address] of Object.entries(implementations)) {
            console.log(`${contractName}: ${address}`);
        }
        
        console.log('\n=== USAGE INSTRUCTIONS ===');
        console.log('1. Individual transactions: Use the individual JSON files for single upgrades');
        console.log('2. Batch upgrade: Use batch_upgrade_bundle.json for all upgrades at once');
        console.log('3. In Gnosis Safe: Go to Apps > Transaction Builder > Import JSON');
        console.log('4. Drag and drop the JSON file(s) and execute');
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
