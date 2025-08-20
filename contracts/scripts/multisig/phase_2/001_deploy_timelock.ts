import { ethers } from "hardhat";

/**
 * Deploy TimelockController for multisig governance
 * 
 * This script deploys the OpenZeppelin TimelockController contract
 * which is used for time-delayed governance actions
 */

async function deployTimelockController() {
    const [deployer] = await ethers.getSigners();
    
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    
    console.log("Deploying TimelockController...");
    console.log("Deployer address:", deployer.address);
    
    // Configuration - these should be set as environment variables
    const multisigAddress = process.env.MULTISIG_ADDRESS || "0x...";
    const delay = process.env.TIMELOCK_DELAY ? parseInt(process.env.TIMELOCK_DELAY) : 24 * 60 * 60; // 24 hours default
    
    if (multisigAddress === "0x...") {
        throw new Error('Please set MULTISIG_ADDRESS environment variable');
    }
    
    console.log("Configuration:");
    console.log("- Multisig address:", multisigAddress);
    console.log("- Delay (seconds):", delay);
    console.log("- Delay (hours):", delay / 3600);
    
    try {
        // Deploy TimelockController
        const TimelockController = await ethers.getContractFactory("TimelockController");
        
        const timelock = await TimelockController.deploy(
            delay,                    // minDelay: minimum delay for operations
            [multisigAddress],        // proposers: addresses that can propose operations
            [multisigAddress],        // executors: addresses that can execute operations
            deployer.address          // admin: address that can grant/revoke proposer/executor roles
        );
        
        await timelock.waitForDeployment();
        const timelockAddress = await timelock.getAddress();
        
        console.log("\nTimelockController deployed successfully!");
        console.log("Timelock address:", timelockAddress);
        
        // Grant proposer role to multisig
        const proposerRole = await (timelock as any).PROPOSER_ROLE();
        await (timelock as any).grantRole(proposerRole, multisigAddress);
        console.log("Proposer role successfully granted to multisig");
        
        // Grant executor role to multisig
        const executorRole = await (timelock as any).EXECUTOR_ROLE();
        await (timelock as any).grantRole(executorRole, multisigAddress);
        console.log("Executor role successfully granted to multisig");
        
        // Revoke admin role from deployer (optional - for security)
        const adminRole = await (timelock as any).DEFAULT_ADMIN_ROLE();
        await (timelock as any).revokeRole(adminRole, deployer.address);
        console.log("Admin role successfully revoked from deployer");
        
        console.log("\n=== Deployment Summary ===");
        console.log("TimelockController:", timelockAddress);
        console.log("Multisig (proposer/executor):", multisigAddress);
        console.log("Delay:", delay, "seconds (", delay / 3600, "hours)");
        console.log("==========================\n");
        
        return {
            timelockAddress,
            multisigAddress,
            delay
        };
        
    } catch (error) {
        console.error("Failed to deploy TimelockController:", error);
        throw error;
    }
}

/**
 * Setup complete governance system
 */
async function setupGovernance() {
    console.log("=== Setting up complete governance system ===\n");
    
    // Deploy TimelockController
    const { timelockAddress, multisigAddress, delay } = await deployTimelockController();
    
    console.log("\n=== Governance Setup Complete ===");
    console.log("Environment variables to set:");
    console.log(`export TIMELOCK_ADDRESS="${timelockAddress}"`);
    console.log(`export MULTISIG_ADDRESS="${multisigAddress}"`);
    console.log("===============================\n");
    
    return {
        timelockAddress,
        multisigAddress,
        delay
    };
}

// Run the setup if this script is executed directly
if (require.main === module) {
    setupGovernance()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}

export {
    deployTimelockController,
    setupGovernance
}; 