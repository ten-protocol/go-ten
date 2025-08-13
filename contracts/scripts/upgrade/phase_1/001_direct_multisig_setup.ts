import { ethers } from "hardhat";

/**
 * Direct Multisig Setup for Initial Mainnet Phase
 * 
 * This script transfers proxy admin ownership directly to a Gnosis Safe multisig
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
        
        console.log("\nCurrent proxy addresses:");
        console.table({
            NetworkConfig: networkConfigAddr,
            CrossChain: addresses.crossChain,
            NetworkEnclaveRegistry: addresses.networkEnclaveRegistry,
            DataAvailabilityRegistry: addresses.dataAvailabilityRegistry
        });
        
        // Get the TransparentUpgradeableProxy contract factory
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
        
        // List of proxies to transfer admin ownership
        const proxies = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        console.log("\n=== Transferring Proxy Admin Ownership to Multisig ===");
        console.log("This will enable immediate upgrades by the multisig (no delays)");
        
        for (const proxy of proxies) {
            console.log(`\n--- Processing ${proxy.name} Proxy ---`);
            
            // Get the proxy contract
            const proxyContract = TransparentUpgradeableProxy.attach(proxy.address);
            
            // Get current admin
            const currentAdmin = await (proxyContract as any).admin();
            console.log(`Current admin: ${currentAdmin}`);
            
            if (currentAdmin.toLowerCase() === multisigAddress.toLowerCase()) {
                console.log(`${proxy.name} proxy admin already transferred to Multisig`);
                continue;
            }
            
            if (currentAdmin.toLowerCase() !== deployer.address.toLowerCase()) {
                console.log(`Warning: ${proxy.name} proxy admin is not the deployer (${currentAdmin})`);
                console.log("Skipping this proxy - manual intervention required");
                continue;
            }
            
            // Transfer admin ownership to multisig
            console.log(`Transferring admin ownership from ${deployer.address} to ${multisigAddress}...`);
            
            const transferTx = await (proxyContract as any).changeAdmin(multisigAddress);
            await transferTx.wait();
            
            console.log(`${proxy.name} proxy admin ownership transferred successfully!`);
            console.log(`Transaction hash: ${transferTx.hash}`);
        }
        
        console.log("\n=== Direct Multisig Setup Complete ===");
        console.log("All proxy admin ownership transferred to Multisig");
        
        
        return {
            multisigAddress,
            networkConfigAddr,
            proxies: proxies.map(p => ({ name: p.name, address: p.address }))
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
        
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
        
        const proxies = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        let allControlled = true;
        
        for (const proxy of proxies) {
            const proxyContract = TransparentUpgradeableProxy.attach(proxy.address);
            const currentAdmin = await (proxyContract as any).admin();
            
            if (currentAdmin.toLowerCase() === multisigAddress.toLowerCase()) {
                console.log(`${proxy.name}: Controlled by Multisig`);
            } else {
                console.log(`${proxy.name}: NOT controlled by Multisig (${currentAdmin})`);
                allControlled = false;
            }
        }
        
        if (allControlled) {
            console.log("\nAll proxies are under Multisig control!");
            console.log("Direct upgrades are now possible (no delays)");
        } else {
            console.log("\nSome proxies are not under Multisig control");
            console.log("Please complete the transfer process");
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
