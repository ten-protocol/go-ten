import { ethers } from "hardhat";

/**
 * Transfer unpauser role to multisig for all contracts that implement PausableWithRoles
 */
async function transferUnpauserRolesOnAllContracts() {
    const [deployer] = await ethers.getSigners()

    if (!deployer) {
        throw new Error("No deployer signer found")
    }
    
    const multisigAddress = process.env.MULTISIG_ADDR || "0x...";
    const networkConfigAddress = process.env.NETWORK_CONFIG_ADDR || "0x...";
    // this isnt stored in NetworkConfig so we need to pass it in
    // const merkleMessageBusAddress = process.env.MERKLE_MESSAGE_BUS_ADDR || "0x...";
    
    if (multisigAddress === "0x...") {
        throw new Error('Please set MULTISIG_ADDR environment variable');
    }
    
    if (networkConfigAddress === "0x...") {
        throw new Error('Please set NETWORK_CONFIG_ADDR environment variable');
    }
    
    // if (merkleMessageBusAddress === "0x...") {
    //   throw new Error("Please set MERKLE_MESSAGE_BUS_ADDR environment variable")
    // }
    
    console.log("Configuration:");
    console.log("- Multisig address:", multisigAddress);
    console.log("- Network config address:", networkConfigAddress);
    // console.log("- Merkle message bus address:", merkleMessageBusAddress);

    const networkConfig = await ethers.getContractAt("NetworkConfig", networkConfigAddress);
    const networkConfigData = await networkConfig.addresses();
    
    // contracts that implement PausableWithRoles
    console.log("\nFound contracts:");
    // console.log("- CrossChain:", networkConfigData.crossChain);
    console.log("- MessageBus:", networkConfigData.messageBus);
    console.log("- NetworkEnclaveRegistry:", networkConfigData.networkEnclaveRegistry);
    console.log("- DataAvailabilityRegistry:", networkConfigData.dataAvailabilityRegistry);
    console.log("- TenBridge:", networkConfigData.l1Bridge)
    console.log("- EthereumBridge:", networkConfigData.l2Bridge)
    // console.log("- MerkleTreeMessageBus:", merkleMessageBusAddress)

    const pausableContracts = [
      // { name: "CrossChain", address: networkConfigData.crossChain },
      { name: "MessageBus", address: networkConfigData.messageBus },
      {
        name: "NetworkEnclaveRegistry",
        address: networkConfigData.networkEnclaveRegistry,
      },
      {
        name: "DataAvailabilityRegistry",
        address: networkConfigData.dataAvailabilityRegistry,
      },
      // { name: "MerkleTreeMessageBus", address: merkleMessageBusAddress },
      { name: "MessageBus", address: networkConfigData.messageBus },
      { name: "TenBridge", address: networkConfigData.l2Bridge },
      { name: "EthereumBridge", address: networkConfigData.l1Bridge },
    ]
    
    console.log(`\nTransferring UNPAUSER_ROLE from ${deployer.address} to ${multisigAddress} on all contracts...`);
    
    for (const contract of pausableContracts) {
        try {
            console.log(`\n--- Processing ${contract.name} at ${contract.address} ---`);
            
            const pausableContract = await ethers.getContractAt(contract.name, contract.address) as any;
            
            const tx = await pausableContract.transferUnpauserRoleToMultisig(multisigAddress)
            const receipt = await tx.wait()
            if (receipt!.status !== 1) {
                throw new Error(`Failed to grant UNPAUSER_ROLE to ${multisigAddress} on ${contract.name}`)
            }
            
        } catch (error) {
            console.error(`Failed to process ${contract.name}:`, error);
        }
    }
    
    console.log(`\n=== Role transfer process completed ===`);
}


async function main() {
    console.log("\n=== Starting transferUnpauserRolesOnAllContracts ===")
    await transferUnpauserRolesOnAllContracts();
    console.log("\n=== Transfer Complete ===");
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

export { transferUnpauserRolesOnAllContracts };
