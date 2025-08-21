import { ethers } from "hardhat";

/**
 * Direct Multisig Setup for Initial Mainnet Phase
 * 
 * This script transfers contract ownership directly to a Gnosis Safe multisig
 * for immediate upgrade control during the initial mainnet release phase.
 * 
 * WARNING: This removes all delay protection. Only use during initial mainnet
 * phase when rapid iteration is needed and the protocol is not public-facing.
 */

async function setupDirectMultisig() {
    const [deployer] = await ethers.getSigners();
    
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    
    console.log("=== Setting up Direct Multisig Control ===");
    console.log("Deployer address:", deployer.address);
    
    // Configuration - these must be set as environment variables
    const multisigAddr = process.env.MULTISIG_ADDR || "0x...";
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR || "0x...";
    const proxyAdminAddr = process.env.PROXY_ADMIN_ADDR || "0x...";

    if (multisigAddr === "0x...") {
        throw new Error('Please set MULTISIG_ADDR environment variable');
    }
    
    if (networkConfigAddr === "0x...") {
        throw new Error('Please set NETWORK_CONFIG_ADDR environment variable');
    }

    if (networkConfigAddr === "0x...") {
        throw new Error('Please set PROXY_ADMIN_ADDR environment variable');
    }

    console.log("Configuration:");
    console.log("- Multisig address:", multisigAddr);
    console.log("- NetworkConfig address:", networkConfigAddr);
    console.log("- ProxyAdmin address:", proxyAdminAddr);

    try {
        // Get addresses from network config
        const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
        const addresses = await networkConfig.addresses();
        
        console.log("\nCurrent contract addresses:");
        console.table({
            NetworkConfig: networkConfigAddr,
            CrossChain: addresses.crossChain,
            NetworkEnclaveRegistry: addresses.networkEnclaveRegistry,
            DataAvailabilityRegistry: addresses.dataAvailabilityRegistry
        });
        
        // List of contracts to transfer ownership
        const contracts = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry },
            { name: "ProxyAdmin", address: proxyAdminAddr }
        ];
        
        console.log("\n=== Transferring Contract Ownership ===");
        console.log("This will initiate the ownership transfer process");
        console.log("Note: You must complete the transfer by accepting ownership in your Gnosis Safe");
        
        for (const contract of contracts) {
            console.log(`\n--- Processing ${contract.name} Contract ---`);
            
            try {
                // Get the contract instance
                const contractInstance = await ethers.getContractAt(contract.name, contract.address);
                
                // Get current owner
                const currentOwner = await (contractInstance as any).owner();
                console.log(`Current owner: ${currentOwner}`);
                
                if (currentOwner.toLowerCase() === multisigAddr.toLowerCase()) {
                    console.log(`${contract.name} ownership already transferred to Multisig`);
                    continue;
                }
                
                if (currentOwner.toLowerCase() !== deployer.address.toLowerCase()) {
                    console.log(`Warning: ${contract.name} owner is not the deployer (${currentOwner})`);
                    console.log("Skipping this contract - manual intervention required");
                    continue;
                }
                
        // Transfer ownership to multisig
        console.log(`Transferring ownership from ${deployer.address} to ${multisigAddr}...`);
        
        const transferTx = await (contractInstance as any).transferOwnership(multisigAddr);
        await transferTx.wait();
        
        console.log(`${contract.name} ownership transfer initiated successfully!`);
        console.log(`Transaction hash: ${transferTx.hash}`);
        console.log(`Ownership transfer initiated for ${contract.name}`);
        console.log(`Pending: You must now accept ownership through your Gnosis Safe`);
        
        // Generate Safe transaction for acceptOwnership
        const acceptOwnershipData = "0x79ba5097"; // acceptOwnership() function selector
        console.log(`Safe Transaction Details for ${contract.name}:`);
        console.log(`   To: ${contract.address}`);
        console.log(`   Value: 0 ETH`);
        console.log(`   Data: ${acceptOwnershipData}`);
        console.log(`   Function: acceptOwnership()`);
        console.log(`   Status: PENDING ACCEPTANCE`);
                
            } catch (error) {
                console.log(`Error processing ${contract.name}:`, error);
                console.log("Skipping this contract");
            }
        }
        
        console.log("\n=== Direct Multisig Setup Complete ===");
        console.log("All ownership transfers have been initiated!");
        console.log("PENDING: You must now accept ownership through your Gnosis Safe");
        
        console.log("\n=== NEXT STEPS IN GNOSIS SAFE UI ===");
        console.log("1. Go to your Gnosis Safe dashboard");
        console.log("2. Click 'New Transaction'");
        console.log("3. For each contract, create a transaction with:");
        console.log("   - To: [Contract Address] (shown above)");
        console.log("   - Value: 0 ETH");
        console.log("   - Data: 0x79ba5097");
        console.log("4. Sign and execute each transaction");
        console.log("5. Run the upgrade script (002_direct_upgrade.ts) to verify ownership");
        
        console.log("\n=== CONTRACT STATUS ===");
        for (const contract of contracts) {
            console.log(`${contract.name}: PENDING ACCEPTANCE`);
            console.log(`   Address: ${contract.address}`);
            console.log(`   Safe Transaction Data: 0x79ba5097`);
        }
        
        return {
            multisigAddress: multisigAddr,
            networkConfigAddr,
            contracts: contracts.map(c => ({ name: c.name, address: c.address }))
        };
        
    } catch (error) {
        console.error("Failed to setup direct multisig control:", error);
        throw error;
    }
}

/**
 * Main setup function
 */
async function main() {
    console.log("\n=== Setup Starting ===")
            // Setup direct multisig control
        const result = await setupDirectMultisig();
        
    console.log("\n=== Setup Complete ===");
}

// Run the setup if this script is executed directly
if (require.main === module) {
    main()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}



export {
    setupDirectMultisig,
};
