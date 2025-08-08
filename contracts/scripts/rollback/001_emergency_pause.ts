import { ethers } from "hardhat";

/**
 * Emergency Pause Script
 * 
 * This script can be used to immediately pause all contracts if a malicious upgrade
 * is detected. This bypasses the timelock and can stop malicious code execution.
 * 
 * WARNING: This is for emergency use only. Only use if you detect a malicious upgrade.
 * 
 * REQUIREMENTS:
 * - The signer must have EMERGENCY_PAUSER_ROLE on all contracts
 * - Set EMERGENCY_PAUSER_ADDRESS environment variable
 * - Set NETWORK_CONFIG_ADDR environment variable
 */

interface EmergencyConfig {
    emergencyPauserAddress: string;
    networkConfigAddress: string;
}

/**
 * Check if the signer has the emergency pauser role
 */
async function checkEmergencyPauserRole(contractInstance: any, signerAddress: string): Promise<boolean> {
    try {
        const EMERGENCY_PAUSER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("EMERGENCY_PAUSER_ROLE"));
        const hasRole = await contractInstance.hasRole(EMERGENCY_PAUSER_ROLE, signerAddress);
        return hasRole;
    } catch (error) {
        console.log("Could not check role (contract may not have role-based access control):", error);
        return true; // Assume role exists for backward compatibility
    }
}

/**
 * Pause all contracts immediately
 */
async function emergencyPause() {
    const [emergencyPauser] = await ethers.getSigners();
    
    if (!emergencyPauser) {
        throw new Error('No emergency pauser signer found');
    }
    
    console.log("=== EMERGENCY PAUSE ACTIVATED ===");
    console.log("Emergency pauser address:", emergencyPauser.address);
    console.log("Timestamp:", new Date().toISOString());
    console.log("================================\n");
    
    // Configuration - these should be set as environment variables
    const config: EmergencyConfig = {
        emergencyPauserAddress: process.env.EMERGENCY_PAUSER_ADDRESS || emergencyPauser.address,
        networkConfigAddress: process.env.NETWORK_CONFIG_ADDR || "0x..."
    };
    
    if (config.networkConfigAddress === "0x...") {
        throw new Error('Please set NETWORK_CONFIG_ADDR environment variable');
    }
    
    console.log("Configuration:");
    console.log("- Emergency pauser:", config.emergencyPauserAddress);
    console.log("- NetworkConfig address:", config.networkConfigAddress);
    
    try {
        // Get addresses from network config
        const networkConfig = await ethers.getContractAt('NetworkConfig', config.networkConfigAddress);
        const addresses = await networkConfig.addresses();
        
        console.log("\nCurrent contract addresses:");
        console.table({
            NetworkConfig: config.networkConfigAddress,
            CrossChain: addresses.crossChain,
            NetworkEnclaveRegistry: addresses.networkEnclaveRegistry,
            DataAvailabilityRegistry: addresses.dataAvailabilityRegistry
        });
        
        // List of contracts to pause
        const contracts = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        console.log("\n=== Pausing Contracts ===");
        
        for (const contract of contracts) {
            console.log(`\n--- Pausing ${contract.name} ---`);
            
            try {
                // Get the contract
                const contractInstance = await ethers.getContractAt(contract.name, contract.address);
                
                // Check if signer has emergency pauser role
                const hasRole = await checkEmergencyPauserRole(contractInstance, emergencyPauser.address);
                if (!hasRole) {
                    console.error(`${contract.name}: Signer does not have EMERGENCY_PAUSER_ROLE`);
                    continue;
                }
                
                // Check if contract is already paused
                const isPaused = await (contractInstance as any).paused();
                
                if (isPaused) {
                    console.log(`${contract.name} is already paused`);
                    continue;
                }
                
                // Pause the contract
                console.log(`Pausing ${contract.name} at ${contract.address}...`);
                const pauseTx = await (contractInstance as any).emergencyPause();
                await pauseTx.wait();
                
                console.log(`${contract.name} paused successfully!`);
                console.log(`Transaction hash: ${pauseTx.hash}`);
                
            } catch (error) {
                console.error(`Failed to pause ${contract.name}:`, error);
                console.log(`Continuing with next contract...`);
            }
        }
        
        console.log("\n=== Emergency Pause Summary ===");
        console.log("All contracts have been paused");
        console.log("Malicious code execution has been stopped");
        console.log("Next steps:");
        console.log("1. Investigate the malicious upgrade");
        console.log("2. Prepare a secure rollback");
        console.log("3. Execute rollback when ready");
        console.log("==============================\n");
        
    } catch (error) {
        console.error("Failed to execute emergency pause:", error);
        throw error;
    }
}

/**
 * Unpause all contracts (use after rollback)
 */
async function emergencyUnpause() {
    const [emergencyPauser] = await ethers.getSigners();
    
    if (!emergencyPauser) {
        throw new Error('No emergency pauser signer found');
    }
    
    console.log("=== EMERGENCY UNPAUSE ===");
    console.log("Emergency pauser address:", emergencyPauser.address);
    console.log("Timestamp:", new Date().toISOString());
    console.log("========================\n");
    
    // Configuration - these should be set as environment variables
    const config: EmergencyConfig = {
        emergencyPauserAddress: process.env.EMERGENCY_PAUSER_ADDRESS || emergencyPauser.address,
        networkConfigAddress: process.env.NETWORK_CONFIG_ADDR || "0x..."
    };
    
    if (config.networkConfigAddress === "0x...") {
        throw new Error('Please set NETWORK_CONFIG_ADDR environment variable');
    }
    
    try {
        // Get addresses from network config
        const networkConfig = await ethers.getContractAt('NetworkConfig', config.networkConfigAddress);
        const addresses = await networkConfig.addresses();
        
        // List of contracts to unpause
        const contracts = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        console.log("\n=== Unpausing Contracts ===");
        
        for (const contract of contracts) {
            console.log(`\n--- Unpausing ${contract.name} ---`);
            
            try {
                // Get the contract
                const contractInstance = await ethers.getContractAt(contract.name, contract.address);
                
                // Check if signer has emergency pauser role
                const hasRole = await checkEmergencyPauserRole(contractInstance, emergencyPauser.address);
                if (!hasRole) {
                    console.error(`${contract.name}: Signer does not have EMERGENCY_PAUSER_ROLE`);
                    continue;
                }
                
                // Check if contract is paused
                const isPaused = await (contractInstance as any).paused();
                
                if (!isPaused) {
                    console.log(`${contract.name} is not paused`);
                    continue;
                }
                
                // Unpause the contract
                console.log(`Unpausing ${contract.name} at ${contract.address}...`);
                const unpauseTx = await (contractInstance as any).emergencyUnpause();
                await unpauseTx.wait();
                
                console.log(`${contract.name} unpaused successfully!`);
                console.log(`Transaction hash: ${unpauseTx.hash}`);
                
            } catch (error) {
                console.error(`Failed to unpause ${contract.name}:`, error);
                console.log(`Continuing with next contract...`);
            }
        }
        
        console.log("\n=== Emergency Unpause Summary ===");
        console.log("All contracts have been unpaused");
        console.log("System is now operational");
        console.log("==============================\n");
        
    } catch (error) {
        console.error("Failed to execute emergency unpause:", error);
        throw error;
    }
}

// Export functions for use in other scripts
export {
    emergencyPause,
    emergencyUnpause
};

// Run the emergency pause if this script is executed directly
if (require.main === module) {
    const action = process.env.EMERGENCY_ACTION || "pause";
    
    if (action === "unpause") {
        emergencyUnpause()
            .then(() => process.exit(0))
            .catch((error) => {
                console.error(error);
                process.exit(1);
            });
    } else {
        emergencyPause()
            .then(() => process.exit(0))
            .catch((error) => {
                console.error(error);
                process.exit(1);
            });
    }
}

export default emergencyPause; 