import { ethers } from "hardhat";
import { Contract, ContractFactory } from "ethers";
import { verifyContract } from "../../utils/verify";

/**
 * Emergency Rollback Script
 * 
 * This script handles emergency rollback procedures after detecting malicious upgrades.
 * It follows the industry-standard approach of using normal timelock processes for security.
 * 
 * IMPORTANT: This script does NOT bypass the timelock - it uses the normal upgrade process
 * to ensure security and prevent rushed decisions during emergencies.
 * 
 * Usage:
 * export NETWORK_CONFIG_ADDR="0x..."
 * export EMERGENCY_PAUSER_ADDRESS="0x..."
 * export TIMELOCK_ADDRESS="0x..."
 * npx hardhat run scripts/emergency/002_emergency_rollback.ts --network mainnet
 */

interface ContractConfig {
  name: string;
  address: string;
  implementation?: string;
}

interface RollbackResult {
  contractName: string;
  currentImpl: string;
  newImpl: string;
  rollbackTx: string;
  success: boolean;
  error?: string;
}

async function main() {
  console.log("Starting Emergency Rollback Procedure");
  console.log("========================================");

  // Load environment variables
  const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
  const emergencyPauserAddr = process.env.EMERGENCY_PAUSER_ADDRESS;
  const timelockAddr = process.env.TIMELOCK_ADDRESS;

  if (!networkConfigAddr || !emergencyPauserAddr || !timelockAddr) {
    throw new Error(
      "Missing required environment variables:\n" +
      "- NETWORK_CONFIG_ADDR: Address of the network configuration contract\n" +
      "- EMERGENCY_PAUSER_ADDRESS: Address with emergency pause permissions\n" +
      "- TIMELOCK_ADDRESS: Address of the timelock contract"
    );
  }

  console.log(`Network Config: ${networkConfigAddr}`);
  console.log(`Emergency Pauser: ${emergencyPauserAddr}`);
  console.log(`Timelock: ${timelockAddr}`);

  // Get signer
  const [deployer] = await ethers.getSigners();
  console.log(`Deployer: ${deployer.address}`);

  // Check if deployer has sufficient balance
  const balance = await deployer.getBalance();
  console.log(`Deployer Balance: ${ethers.utils.formatEther(balance)} ETH`);

  if (balance.lt(ethers.utils.parseEther("0.1"))) {
    throw new Error("Insufficient balance for deployment. Need at least 0.1 ETH");
  }

  // Define contracts to rollback
  const contractsToRollback: ContractConfig[] = [
    {
      name: "CrossChain",
      address: networkConfigAddr, // This will be updated with actual address
    },
    {
      name: "NetworkEnclaveRegistry",
      address: networkConfigAddr, // This will be updated with actual address
    },
    {
      name: "DataAvailabilityRegistry",
      address: networkConfigAddr, // This will be updated with actual address
    },
  ];

  console.log("\nContracts to Rollback:");
  contractsToRollback.forEach((contract, index) => {
    console.log(`  ${index + 1}. ${contract.name} at ${contract.address}`);
  });

  // Step 1: Deploy Secure Implementations
  console.log("\nStep 1: Deploying Secure Implementations");
  console.log("=============================================");

  const secureImplementations: { [key: string]: string } = {};

  for (const contractConfig of contractsToRollback) {
    try {
      console.log(`\nDeploying secure implementation for ${contractConfig.name}...`);
      
      // Deploy the secure implementation
      const secureImpl = await deploySecureImplementation(contractConfig.name);
      secureImplementations[contractConfig.name] = secureImpl.address;
      
      console.log(`${contractConfig.name} secure implementation deployed at: ${secureImpl.address}`);
      
      // Verify contract on Etherscan (if not on local network)
      if (process.env.ETHERSCAN_API_KEY && process.env.HARDHAT_NETWORK !== "localhost") {
        console.log(`Verifying ${contractConfig.name} implementation on Etherscan...`);
        try {
          await verifyContract(secureImpl.address, []);
          console.log(`${contractConfig.name} verified on Etherscan`);
        } catch (verifyError) {
          console.log(`Failed to verify ${contractConfig.name}: ${verifyError}`);
        }
      }
      
    } catch (error) {
      console.error(`Failed to deploy secure implementation for ${contractConfig.name}:`, error);
      throw error;
    }
  }

  // Step 2: Execute Rollbacks Through Timelock
  console.log("\nStep 2: Executing Rollbacks Through Timelock");
  console.log("===============================================");
  console.log("IMPORTANT: Rollbacks will be scheduled with 24-hour delay for security");
  console.log("   This follows industry standards and cannot be bypassed");

  const rollbackResults: RollbackResult[] = [];

  for (const contractConfig of contractsToRollback) {
    try {
      console.log(`\nRolling back ${contractConfig.name}...`);
      
      const secureImplAddress = secureImplementations[contractConfig.name];
      if (!secureImplAddress) {
        throw new Error(`No secure implementation found for ${contractConfig.name}`);
      }

      // Get current implementation
      const currentImpl = await getCurrentImplementation(contractConfig.address);
      console.log(`   Current implementation: ${currentImpl}`);
      console.log(`   New secure implementation: ${secureImplAddress}`);

      // Execute rollback through timelock
      const rollbackResult = await executeRollback(
        contractConfig.address,
        secureImplAddress,
        timelockAddr,
        deployer
      );

      rollbackResults.push({
        contractName: contractConfig.name,
        currentImpl,
        newImpl: secureImplAddress,
        rollbackTx: rollbackResult.transactionHash,
        success: true,
      });

      console.log(`${contractConfig.name} rollback scheduled successfully`);
      console.log(`   Transaction: ${rollbackResult.transactionHash}`);
      console.log(`   Rollback will execute after 24-hour timelock delay`);

    } catch (error) {
      console.error(`Failed to rollback ${contractConfig.name}:`, error);
      
      rollbackResults.push({
        contractName: contractConfig.name,
        currentImpl: "Unknown",
        newImpl: secureImplementations[contractConfig.name] || "Unknown",
        rollbackTx: "Failed",
        success: false,
        error: error.message,
      });
    }
  }

  // Step 3: Summary Report
  console.log("\nStep 3: Rollback Summary Report");
  console.log("===================================");

  const successfulRollbacks = rollbackResults.filter(r => r.success);
  const failedRollbacks = rollbackResults.filter(r => !r.success);

  console.log(`\nSuccessful Rollbacks: ${successfulRollbacks.length}/${rollbackResults.length}`);
  successfulRollbacks.forEach(result => {
    console.log(`   ${result.contractName}: ${result.currentImpl} → ${result.newImpl}`);
    console.log(`   TX: ${result.rollbackTx}`);
  });

  if (failedRollbacks.length > 0) {
    console.log(`\nFailed Rollbacks: ${failedRollbacks.length}/${rollbackResults.length}`);
    failedRollbacks.forEach(result => {
      console.log(`   ${result.contractName}: ${result.error}`);
    });
  }

  // Step 4: Next Steps Instructions
  console.log("\nNext Steps");
  console.log("=============");
  console.log("1. Wait for 24-hour timelock delay to complete");
  console.log("2. Monitor rollback execution on blockchain");
  console.log("3. Verify all contracts are using secure implementations");
  console.log("4. Unpause contracts using emergency unpause script");
  console.log("5. Test system functionality");
  console.log("6. Document incident and lessons learned");

  console.log("\nSecurity Note:");
  console.log("   The 24-hour timelock delay is a security feature that prevents");
  console.log("   rushed decisions during emergencies. This follows industry");
  console.log("   standards used by Uniswap, Compound, Aave, and other major protocols.");

  console.log("\nEmergency Rollback Procedure Complete!");
}

/**
 * Deploy a secure implementation of the specified contract
 */
async function deploySecureImplementation(contractName: string): Promise<Contract> {
  let contractFactory: ContractFactory;
  
  switch (contractName) {
    case "CrossChain":
      contractFactory = await ethers.getContractFactory("CrossChain");
      break;
    case "NetworkEnclaveRegistry":
      contractFactory = await ethers.getContractFactory("NetworkEnclaveRegistry");
      break;
    case "DataAvailabilityRegistry":
      contractFactory = await ethers.getContractFactory("DataAvailabilityRegistry");
      break;
    default:
      throw new Error(`Unknown contract type: ${contractName}`);
  }

  // Deploy with constructor arguments if needed
  // Note: These will need to be updated based on actual contract constructors
  const secureImpl = await contractFactory.deploy();
  await secureImpl.deployed();
  
  return secureImpl;
}

/**
 * Get the current implementation address of a proxy contract
 */
async function getCurrentImplementation(proxyAddress: string): Promise<string> {
  try {
    // Try to get implementation from OpenZeppelin proxy pattern
    const proxy = new ethers.Contract(
      proxyAddress,
      ["function implementation() view returns (address)"],
      ethers.provider
    );
    
    return await proxy.implementation();
  } catch (error) {
    try {
      // Try alternative method for different proxy patterns
      const proxy = new ethers.Contract(
        proxyAddress,
        ["function getImplementation() view returns (address)"],
        ethers.provider
      );
      
      return await proxy.getImplementation();
    } catch (altError) {
      console.log(`Could not determine current implementation for ${proxyAddress}`);
      return "Unknown";
    }
  }
}

/**
 * Execute rollback through timelock
 */
async function executeRollback(
  proxyAddress: string,
  newImplementation: string,
  timelockAddress: string,
  deployer: any
): Promise<any> {
  // This is a simplified example - actual implementation will depend on your timelock contract
  const timelock = new ethers.Contract(
    timelockAddress,
    [
      "function schedule(address target, uint256 value, bytes calldata data, bytes32 predecessor, bytes32 salt, uint256 delay) external",
      "function execute(address target, uint256 value, bytes calldata data, bytes32 predecessor, bytes32 salt) external payable"
    ],
    deployer
  );

  // Prepare upgrade data
  const upgradeData = ethers.utils.defaultAbiCoder.encode(
    ["address"],
    [newImplementation]
  );

  // Schedule the rollback
  const salt = ethers.utils.randomBytes(32);
  const delay = 24 * 60 * 60; // 24 hours in seconds
  
  const scheduleTx = await timelock.schedule(
    proxyAddress,
    0,
    upgradeData,
    ethers.constants.HashZero,
    salt,
    delay
  );
  
  await scheduleTx.wait();
  
  return scheduleTx;
}

// Error handling
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("Emergency rollback failed:", error);
    process.exit(1);
  });
