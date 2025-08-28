import { ethers } from "hardhat";
const hre = require("hardhat");

/**
 * Direct Upgrade Script for Initial Mainnet Phase
 *
 * This script allows the multisig to upgrade contracts immediately without any delays.
 */

interface UpgradeConfig {
    contractName: string;
    proxyAddress: string;
    description: string;
}

/**
 * Verify that the multisig wallet controls all contracts
 */
async function verifyMultisigOwnership(
    multisigAddress: string,
    networkConfigAddr: string,
    proxyAdminAddr: string
): Promise<boolean> {
    console.log("Checking contract ownership...");

    try {
        const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
        const addresses = await networkConfig.addresses();

        const contracts = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry },
            { name: "ProxyAdmin", address: proxyAdminAddr }
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

        // Additional verification: Check that ProxyAdmin is the admin of all proxy contracts
        console.log("\n=== Verifying Proxy Admin Control ===");
        try {
            const proxyAdmin = await ethers.getContractAt("ProxyAdmin", proxyAdminAddr);
            
            // Check each proxy contract to see if ProxyAdmin is its admin
            const proxyContracts = [
                { name: "CrossChain", address: addresses.crossChain },
                { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
                { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
            ];

            for (const proxy of proxyContracts) {
                try {
                    // Check if this proxy is owned by the ProxyAdmin
                    // We can do this by checking the proxy's admin storage slot directly
                    const adminSlot = "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103";
                    const adminOfProxy = await ethers.provider.getStorage(proxy.address, adminSlot);
                    
                    if (adminOfProxy !== "0x0000000000000000000000000000000000000000000000000000000000000000") {
                        const adminAddress = "0x" + adminOfProxy.slice(26); // Remove padding
                        console.log(`${proxy.name} proxy admin: ${adminAddress}`);
                        
                        if (adminAddress.toLowerCase() === proxyAdminAddr.toLowerCase()) {
                            console.log(`${proxy.name}: ProxyAdmin IS the admin`);
                        } else {
                            console.log(`${proxy.name}: ProxyAdmin is NOT the admin (Current: ${adminAddress})`);
                            allControlled = false;
                        }
                    } else {
                        console.log(`${proxy.name}: No proxy admin found`);
                        allControlled = false;
                    }
                } catch (error) {
                    console.log(`Error checking ${proxy.name} proxy admin:`, error);
                    allControlled = false;
                }
            }
        } catch (error) {
            console.log(`Error checking ProxyAdmin contract:`, error);
            allControlled = false;
        }

        if (allControlled) {
            console.log("\nAll contracts are under Multisig control!");
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
    const multisigAddress = process.env.MULTISIG_ADDR || "0x...";
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR || "0x...";
    const proxyAdminAddr = process.env.PROXY_ADMIN_ADDR || "0x...";

    if (multisigAddress === "0x..." || networkConfigAddr === "0x..." || proxyAdminAddr === "0x...") {
        throw new Error('Please set MULTISIG_ADDR, NETWORK_CONFIG_ADDR, and PROXY_ADMIN_ADDR environment variables');
    }

    console.log("Configuration:");
    console.log("- Multisig address:", multisigAddress);
    console.log("- NetworkConfig address:", networkConfigAddr);
    console.log("- ProxyAdmin address:", proxyAdminAddr);
    console.log("- Deployer address:", deployer.address);

    // Verify multisig ownership before proceeding
    console.log("\n=== Verifying Multisig Ownership ===");
    const ownershipVerified = await verifyMultisigOwnership(multisigAddress, networkConfigAddr, proxyAdminAddr);

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

        // Define upgrade configurations for all contracts
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

        console.log('\n=== DEPLOYING NEW IMPLEMENTATIONS ===');
        
        // Deploy all new implementations first
        const implementations: { [key: string]: string } = {};
        
        for (const config of upgradeConfigs) {
            console.log(`\n--- Deploying ${config.contractName} Implementation ---`);
            try {
                const newImplementation = await deployNewImplementation(config.contractName);
                implementations[config.contractName] = newImplementation;
                
                // Verify the new implementation is valid
                console.log(`Verifying new implementation at ${newImplementation}...`);
                const newImplContract = await ethers.getContractAt(config.contractName, newImplementation);
                
                // Try to call a basic function to ensure it's working
                try {
                    const owner = await (newImplContract as any).owner();
                    console.log(`New implementation verified - owner(): ${owner}`);
                } catch (error) {
                    console.log(`Warning: New implementation may have issues:`, (error as Error).message);
                }
                
            } catch (error) {
                console.error(`Failed to deploy ${config.contractName}:`, error);
                throw error;
            }
        }
        
        console.log('\n=== UPGRADE INSTRUCTIONS ===');
        console.log('For each contract, create a new transaction in Gnosis Safe:');
        console.log('='.repeat(80));
        
        for (const config of upgradeConfigs) {
            const newImplementation = implementations[config.contractName];
            if (!newImplementation) {
                console.error(`No implementation found for ${config.contractName}`);
                continue;
            }
            
            console.log(`\n--- ${config.contractName} Upgrade ---`);
            console.log(`1. Create a new transaction to ${proxyAdminAddr}`);
            console.log(`2. Select the upgrade function`);
            console.log(`3. Set proxy address to ${config.proxyAddress}`);
            console.log(`4. Set implementation address to ${newImplementation}`);
            console.log(`5. Execute the transaction`);
        }
        
        console.log('\n' + '='.repeat(80));
        
        console.log('\n=== IMPLEMENTATION ADDRESSES ===');
        console.log('Save these for verification:');
        for (const [contractName, address] of Object.entries(implementations)) {
            console.log(`${contractName}: ${address}`);
        }

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