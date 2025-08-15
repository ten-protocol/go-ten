import { ethers } from "hardhat";

/**
 * Transition Script: Direct Multisig → Timelock Governance
 * 
 * This script transitions the protocol from direct multisig control (Option 1)
 * to timelock-based governance (Option 2) after the initial mainnet phase.
 * 
 * IMPORTANT: Only run this after:
 * 1. Initial mainnet phase is stable (2+ weeks without critical issues)
 * 2. TimelockController has been deployed and tested
 * 3. Team is ready for transparent governance
 */

interface TransitionConfig {
    timelockAddress: string;
    multisigAddress: string;
    networkConfigAddr: string;
}

/**
 * Verify current state before transition
 */
async function verifyPreTransitionState(config: TransitionConfig): Promise<boolean> {
    console.log("=== Verifying Pre-Transition State ===");
    
    try {
        const networkConfig = await ethers.getContractAt('NetworkConfig', config.networkConfigAddr);
        const addresses = await networkConfig.addresses();
        
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
        
        const proxies = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        console.log("Checking current proxy admin ownership...");
        
        let allReady = true;
        
        for (const proxy of proxies) {
            const proxyContract = TransparentUpgradeableProxy.attach(proxy.address);
            const currentAdmin = await (proxyContract as any).admin();
            
            if (currentAdmin.toLowerCase() === config.multisigAddress.toLowerCase()) {
                console.log(`${proxy.name}: Ready for transition (currently multisig)`);
            } else if (currentAdmin.toLowerCase() === config.timelockAddress.toLowerCase()) {
                console.log(`${proxy.name}: Already under timelock control`);
            } else {
                console.log(`${proxy.name}: Unexpected admin (${currentAdmin})`);
                allReady = false;
            }
        }
        
        if (allReady) {
            console.log("\nAll proxies are ready for transition to timelock");
        } else {
            console.log("\nSome proxies are not ready for transition");
        }
        
        return allReady;
        
    } catch (error) {
        console.error("Failed to verify pre-transition state:", error);
        return false;
    }
}

/**
 * Verify timelock is properly configured
 */
async function verifyTimelockConfiguration(config: TransitionConfig): Promise<boolean> {
    console.log("\n=== Verifying Timelock Configuration ===");
    
    try {
        const TimelockController = await ethers.getContractFactory("TimelockController");
        const timelock = TimelockController.attach(config.timelockAddress) as any;
        
        // Check if multisig has proposer role
        const proposerRole = await timelock.PROPOSER_ROLE();
        const hasProposerRole = await timelock.hasRole(proposerRole, config.multisigAddress);
        
        if (hasProposerRole) {
            console.log("Multisig has proposer role");
        } else {
            console.log("Multisig missing proposer role");
            return false;
        }
        
        // Check if multisig has executor role
        const executorRole = await timelock.EXECUTOR_ROLE();
        const hasExecutorRole = await timelock.hasRole(executorRole, config.multisigAddress);
        
        if (hasExecutorRole) {
            console.log("Multisig has executor role");
        } else {
            console.log("Multisig missing executor role");
            return false;
        }
        
        // Check delay
        const delay = await timelock.getMinDelay();
        console.log(`Timelock delay: ${delay} seconds (${delay / 3600} hours)`);
        
        console.log("\nTimelock is properly configured");
        return true;
        
    } catch (error) {
        console.error("Failed to verify timelock configuration:", error);
        return false;
    }
}

/**
 * Perform the transition from direct multisig to timelock
 */
async function performTransition(config: TransitionConfig): Promise<boolean> {
    console.log("\n=== Performing Transition to Timelock ===");
    
    try {
        const networkConfig = await ethers.getContractAt('NetworkConfig', config.networkConfigAddr);
        const addresses = await networkConfig.addresses();
        
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
        
        const proxies = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        let allTransferred = true;
        
        for (const proxy of proxies) {
            console.log(`\n--- Transitioning ${proxy.name} ---`);
            
            const proxyContract = TransparentUpgradeableProxy.attach(proxy.address);
            const currentAdmin = await (proxyContract as any).admin();
            
            if (currentAdmin.toLowerCase() === config.timelockAddress.toLowerCase()) {
                console.log(`${proxy.name} already under timelock control, skipping`);
                continue;
            }
            
            if (currentAdmin.toLowerCase() !== config.multisigAddress.toLowerCase()) {
                console.log(`Warning: ${proxy.name} admin is not multisig (${currentAdmin})`);
                console.log("Skipping this proxy - manual intervention required");
                allTransferred = false;
                continue;
            }
            
            // Transfer admin ownership from multisig to timelock
            console.log(`Transferring admin from multisig to timelock...`);
            
            const transferTx = await (proxyContract as any).changeAdmin(config.timelockAddress);
            await transferTx.wait();
            
            console.log(`${proxy.name} admin transferred to timelock`);
            console.log(`Transaction hash: ${transferTx.hash}`);
        }
        
        return allTransferred;
        
    } catch (error) {
        console.error("Failed to perform transition:", error);
        return false;
    }
}

/**
 * Verify post-transition state
 */
async function verifyPostTransitionState(config: TransitionConfig): Promise<boolean> {
    console.log("\n=== Verifying Post-Transition State ===");
    
    try {
        const networkConfig = await ethers.getContractAt('NetworkConfig', config.networkConfigAddr);
        const addresses = await networkConfig.addresses();
        
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");
        
        const proxies = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];
        
        let allUnderTimelock = true;
        
        for (const proxy of proxies) {
            const proxyContract = TransparentUpgradeableProxy.attach(proxy.address);
            const currentAdmin = await (proxyContract as any).admin();
            
            if (currentAdmin.toLowerCase() === config.timelockAddress.toLowerCase()) {
                console.log(`${proxy.name}: Under timelock control`);
            } else {
                console.log(`${proxy.name}: NOT under timelock control (${currentAdmin})`);
                allUnderTimelock = false;
            }
        }
        
        if (allUnderTimelock) {
            console.log("\nAll proxies are now under timelock control!");
            console.log("Your protocol has successfully transitioned to transparent governance");
        } else {
            console.log("\nSome proxies are not under timelock control");
            console.log("Please investigate and complete the transition");
        }
        
        return allUnderTimelock;
        
    } catch (error) {
        console.error("Failed to verify post-transition state:", error);
        return false;
    }
}

/**
 * Test timelock functionality
 */
async function testTimelockFunctionality(config: TransitionConfig): Promise<boolean> {
    console.log("\n=== Testing Timelock Functionality ===");
    
    try {
        const TimelockController = await ethers.getContractFactory("TimelockController");
        const timelock = TimelockController.attach(config.timelockAddress) as any;
        
        // Get current delay
        const delay = await timelock.getMinDelay();
        console.log(`Timelock delay: ${delay} seconds (${delay / 3600} hours)`);
        
        // Test a simple operation (like checking roles)
        const adminRole = await timelock.DEFAULT_ADMIN_ROLE();
        const hasAdminRole = await timelock.hasRole(adminRole, config.multisigAddress);
        
        console.log(`Multisig has admin role: ${hasAdminRole}`);
        
        console.log("Timelock functionality test completed");
        return true;
        
    } catch (error) {
        console.error("Failed to test timelock functionality:", error);
        return false;
    }
}

/**
 * Main transition function
 */
async function main() {
    console.log("=== Transition Script: Direct Multisig → Timelock Governance ===\n");
    
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    
    // Configuration validation
    const config: TransitionConfig = {
        timelockAddress: process.env.TIMELOCK_ADDRESS || "0x...",
        multisigAddress: process.env.MULTISIG_ADDRESS || "0x...",
        networkConfigAddr: process.env.NETWORK_CONFIG_ADDR || "0x..."
    };
    
    if (config.timelockAddress === "0x..." || 
        config.multisigAddress === "0x..." || 
        config.networkConfigAddr === "0x...") {
        throw new Error('Please set TIMELOCK_ADDRESS, MULTISIG_ADDRESS, and NETWORK_CONFIG_ADDR environment variables');
    }
    
    console.log("Configuration:");
    console.log("- Timelock address:", config.timelockAddress);
    console.log("- Multisig address:", config.multisigAddress);
    console.log("- NetworkConfig address:", config.networkConfigAddr);
    console.log("- Deployer address:", deployer.address);
    
    try {
        // Step 1: Verify pre-transition state
        const preTransitionReady = await verifyPreTransitionState(config);
        if (!preTransitionReady) {
            throw new Error("Pre-transition state verification failed");
        }
        
        // Step 2: Verify timelock configuration
        const timelockReady = await verifyTimelockConfiguration(config);
        if (!timelockReady) {
            throw new Error("Timelock configuration verification failed");
        }
        
        // Step 3: Perform the transition
        const transitionSuccess = await performTransition(config);
        if (!transitionSuccess) {
            throw new Error("Transition failed for some proxies");
        }
        
        // Step 4: Verify post-transition state
        const postTransitionSuccess = await verifyPostTransitionState(config);
        if (!postTransitionSuccess) {
            throw new Error("Post-transition state verification failed");
        }
        
        // Step 5: Test timelock functionality
        await testTimelockFunctionality(config);
        
        console.log("\n=== Transition Complete ===");
    } catch (error) {
        console.error("Transition failed:", error);
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
    verifyPreTransitionState,
    verifyTimelockConfiguration,
    performTransition,
    verifyPostTransitionState,
    testTimelockFunctionality
};
