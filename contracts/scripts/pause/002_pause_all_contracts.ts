import { ethers } from "hardhat";

/**
 * Pause or unpause all contracts that implement PausableWithRoles and verify their status
 */
async function pauseOrUnpauseAllContracts() {
    const [deployer] = await ethers.getSigners()

    if (!deployer) {
        throw new Error("No deployer signer found")
    }
    
    const networkConfigAddress = process.env.NETWORK_CONFIG_ADDR || "0x...";
    const merkleMessageBusAddress = process.env.MERKLE_TREE_MESSAGE_BUS_ADDR || "0x...";
    const action = process.env.ACTION || "PAUSE";
    
    if (networkConfigAddress === "0x...") {
        throw new Error('Please set NETWORK_CONFIG_ADDR environment variable');
    }
    
    if (merkleMessageBusAddress === "0x...") {
      throw new Error("Please set MERKLE_TREE_MESSAGE_BUS_ADDR environment variable")
    }
    
    if (action !== "PAUSE" && action !== "UNPAUSE") {
        throw new Error('ACTION environment variable must be either "PAUSE" or "UNPAUSE"');
    }
    
    console.log("Configuration:");
    console.log("- Deployer address:", deployer.address);
    console.log("- Network config address:", networkConfigAddress);
    console.log("- Merkle message bus address:", merkleMessageBusAddress);
    console.log("- Action:", action);

    const networkConfig = await ethers.getContractAt("NetworkConfig", networkConfigAddress);
    const networkConfigData = await networkConfig.addresses();
    
    // contracts that implement PausableWithRoles
    console.log("\nFound contracts:");
    console.log("- MessageBus:", networkConfigData.messageBus);
    console.log("- NetworkEnclaveRegistry:", networkConfigData.networkEnclaveRegistry);
    console.log("- DataAvailabilityRegistry:", networkConfigData.dataAvailabilityRegistry);
    console.log("- TenBridge:", networkConfigData.l1Bridge)
    console.log("- CrossChainMessenger:", networkConfigData.l1CrossChainMessenger)
    console.log("- MerkleTreeMessageBus:", merkleMessageBusAddress)

    const pausableContracts = [
      { name: "MessageBus", address: networkConfigData.messageBus },
      {
        name: "NetworkEnclaveRegistry",
        address: networkConfigData.networkEnclaveRegistry,
      },
      {
        name: "DataAvailabilityRegistry",
        address: networkConfigData.dataAvailabilityRegistry,
      },
      { name: "TenBridge", address: networkConfigData.l1Bridge },
      { name: "CrossChainMessenger", address: networkConfigData.l1CrossChainMessenger },
      { name: "MerkleTreeMessageBus", address: merkleMessageBusAddress },
    ]
    
    console.log(`\n${action === "PAUSE" ? "Pausing" : "Unpausing"} all contracts using deployer account ${deployer.address}...`);
    
    // First, pause or unpause all contracts
    for (const contract of pausableContracts) {
        try {
            console.log(`\n--- ${action === "PAUSE" ? "Pausing" : "Unpausing"} ${contract.name} at ${contract.address} ---`);
            
            const pausableContract = await ethers.getContractAt(contract.name, contract.address) as any;
            
            // Check current pause status
            const isPaused = await pausableContract.paused();
            
            if (action === "PAUSE") {
                if (isPaused) {
                    console.log(`${contract.name} is already paused, skipping...`);
                    continue;
                }
                
                const tx = await pausableContract.pause()
                const receipt = await tx.wait()
                if (receipt!.status !== 1) {
                    throw new Error(`Failed to pause ${contract.name}`)
                }
                
                console.log(`Successfully paused ${contract.name}`);
            } else { // UNPAUSE
                if (!isPaused) {
                    console.log(`${contract.name} is already unpaused, skipping...`);
                    continue;
                }
                
                const tx = await pausableContract.unpause()
                const receipt = await tx.wait()
                if (receipt!.status !== 1) {
                    throw new Error(`Failed to unpause ${contract.name}`)
                }
                
                console.log(`Successfully unpaused ${contract.name}`);
            }
            
        } catch (error) {
            console.error(`Failed to ${action.toLowerCase()} ${contract.name}:`, error);
        }
    }
    
    console.log(`\n=== ${action} process completed ===`);
    
    // Now verify all contracts have the expected status
    console.log(`\nVerifying all contracts are ${action === "PAUSE" ? "paused" : "unpaused"}...`);
    
    for (const contract of pausableContracts) {
        try {
            console.log(`\n--- Verifying ${contract.name} at ${contract.address} ---`);
            
            const pausableContract = await ethers.getContractAt(contract.name, contract.address) as any;
            
            const isPaused = await pausableContract.paused();
            const expectedStatus = action === "PAUSE";
            
            if (isPaused === expectedStatus) {
                console.log(`${contract.name} is successfully ${action === "PAUSE" ? "paused" : "unpaused"}`);
            } else {
                console.log(`${contract.name} is ${isPaused ? "paused" : "unpaused"} but should be ${expectedStatus ? "paused" : "unpaused"} - this is an error!`);
            }
            
        } catch (error) {
            console.error(`Failed to verify ${action.toLowerCase()} status for ${contract.name}:`, error);
        }
    }
    
    console.log(`\n=== Verification process completed ===`);
}

async function main() {
    await pauseOrUnpauseAllContracts();
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

export { pauseOrUnpauseAllContracts };
