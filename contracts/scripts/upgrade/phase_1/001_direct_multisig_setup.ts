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
    console.log("WARNING: This removes all upgrade delay protection!");
    console.log("Only use during initial mainnet phase for rapid iteration.\n");
    
    console.log("Deployer address:", deployer.address);
    
    // Configuration - these must be set as environment variables
    const multisigAddress = process.env.MULTISIG_ADDRESS || "0x...";
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR || "0x...";
    
    if (multisigAddress === "0x...") {
        throw new Error('Please set MULTISIG_ADDRESS environment variable');
    }
    
    if (networkConfigAddr === "0x...") {
        throw new Error('Please set NETWORK_CONFIG_ADDR environment variable');
    }
    
    console.log("Configuration:");
    console.log("- Multisig address:", multisigAddress);
    console.log("- NetworkConfig address:", networkConfigAddr);
    
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
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        console.log("\n=== Transferring and Accepting Contract Ownership ===");
        console.log("This will complete the 2-step ownership transfer process");
        
        for (const contract of contracts) {
            console.log(`\n--- Processing ${contract.name} Contract ---`);
            
            try {
                // Get the contract instance
                const contractInstance = await ethers.getContractAt(contract.name, contract.address);
                
                // Get current owner
                const currentOwner = await (contractInstance as any).owner();
                console.log(`Current owner: ${currentOwner}`);
                
                if (currentOwner.toLowerCase() === multisigAddress.toLowerCase()) {
                    console.log(`${contract.name} ownership already transferred to Multisig`);
                    continue;
                }
                
                if (currentOwner.toLowerCase() !== deployer.address.toLowerCase()) {
                    console.log(`Warning: ${contract.name} owner is not the deployer (${currentOwner})`);
                    console.log("Skipping this contract - manual intervention required");
                    continue;
                }
                
                        // Transfer ownership to multisig
        console.log(`Transferring ownership from ${deployer.address} to ${multisigAddress}...`);
        
        const transferTx = await (contractInstance as any).transferOwnership(multisigAddress);
        await transferTx.wait();
        
        console.log(`${contract.name} ownership transfer initiated successfully!`);
        console.log(`Transaction hash: ${transferTx.hash}`);
        
        // Now accept the ownership transfer as the multisig
        console.log(`Accepting ownership transfer for ${contract.name}...`);
        
        // Switch to multisig signer for accepting ownership
        const multisigSigner = await ethers.getSigner(multisigAddress);
        if (!multisigSigner) {
            console.log(`Warning: Could not get multisig signer for ${contract.name}`);
            continue;
        }
        
        const contractWithMultisig = contractInstance.connect(multisigSigner);
        const acceptTx = await (contractWithMultisig as any).acceptOwnership();
        await acceptTx.wait();
        
        console.log(`${contract.name} ownership accepted successfully!`);
        console.log(`Accept transaction hash: ${acceptTx.hash}`);
        
        // Verify the transfer
        const newOwner = await (contractInstance as any).owner();
        console.log(`New owner: ${newOwner}`);
        
        if (newOwner.toLowerCase() === multisigAddress.toLowerCase()) {
            console.log(`${contract.name} ownership transfer completed successfully!`);
        } else {
            console.log(`Warning: ${contract.name} ownership transfer may have failed`);
        }
                
            } catch (error) {
                console.log(`Error processing ${contract.name}:`, error);
                console.log("Skipping this contract");
            }
        }
        
        console.log("\n=== Direct Multisig Setup Complete ===");
        console.log("All contract ownership transfers completed (2-step process)");
        
        
        return {
            multisigAddress,
            networkConfigAddr,
            contracts: contracts.map(c => ({ name: c.name, address: c.address }))
        };
        
    } catch (error) {
        console.error("Failed to setup direct multisig control:", error);
        throw error;
    }
}

/**
 * Verify multisig control
 */
async function verifyMultisigControl() {
    const multisigAddress = process.env.MULTISIG_ADDRESS || "0x...";
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR || "0x...";
    
    if (multisigAddress === "0x..." || networkConfigAddr === "0x...") {
        throw new Error('Please set MULTISIG_ADDRESS and NETWORK_CONFIG_ADDR environment variables');
    }
    
    console.log("=== Verifying Multisig Control ===");
    
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
                    console.log(`${contract.name}: NOT controlled by Multisig (${currentOwner})`);
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
            console.log("Please check the ownership transfer process");
        }
        
        return allControlled;
        
    } catch (error) {
        console.error("Failed to verify multisig control:", error);
        throw error;
    }
}

/**
 * Main setup function
 */
async function main() {
    console.log("\n=== Setup Starting ===")
    // Setup direct multisig control
    await setupDirectMultisig();
    
    console.log("\n=== Verification ===");
    await verifyMultisigControl();
    
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
    verifyMultisigControl
};
