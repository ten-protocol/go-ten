import { ethers } from "hardhat";
import { HardhatRuntimeEnvironment } from "hardhat/types";
async function main(hre: HardhatRuntimeEnvironment) {
  // Get network information
  const network = await ethers.provider.getNetwork();
  console.log(`Connected to network: ${network.name} (chainId: ${network.chainId})`);
  
  // Get signers - first is funded, create an unfunded one
  const signers = await ethers.getSigners();
  const fundedWallet = signers[0]!;
  const unfundedWallet = signers[2]!;
  
  // Print wallet addresses
  console.log(`Funded Wallet: ${fundedWallet.address}`);
  console.log(`Unfunded Wallet: ${unfundedWallet.address}`);
  
  // Set gas price for testing
  const gasPrice = ethers.parseUnits("100", "gwei");
  console.log(`Gas price for tests: ${ethers.formatUnits(gasPrice, 'gwei')} gwei`);
  
  // Get balance of funded wallet
  const fundedBalance = await ethers.provider.getBalance(fundedWallet.address);
  console.log(`Funded Wallet Balance: ${ethers.formatEther(fundedBalance)} ETH`);
  
  console.log("\n=== Testing eth_call ===");
  
  // Test 1: eth_call with funded account (basic call to check ETH balance)
  try {
    console.log("\nTest 1: eth_call with funded account, no gas price");
    const result1 = await ethers.provider.call({
      from: fundedWallet.address,
      to: fundedWallet.address,
      value: 0
    });
    console.log(`Result: ${result1}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  // Test 2: eth_call with funded account with gas price
  try {
    console.log("\nTest 2: eth_call with funded account, with gas price");
    const result2 = await ethers.provider.call({
      from: fundedWallet.address,
      to: fundedWallet.address,
      value: 0,
      gasPrice
    });
    console.log(`Result: ${result2}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  // Test 3: eth_call with unfunded account, no gas price
  try {
    console.log("\nTest 3: eth_call with unfunded account, no gas price");
    const result3 = await ethers.provider.call({
      from: unfundedWallet.address,
      to: fundedWallet.address,
      value: 0,
      gasPrice: null
    });
    console.log(`Result: ${result3}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  // Test 4: eth_call with unfunded account, with gas price (should fail)
  try {
    console.log("\nTest 4: eth_call with unfunded account, with gas price");
    const result4 = await ethers.provider.call({
      from: unfundedWallet.address,
      to: fundedWallet.address,
      value: 0,
      gasPrice
    });
    console.log(`Result: ${result4}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  console.log("\n=== Testing eth_estimateGas ===");
  
  // Test 5: Estimate gas - Funded account with no gas price
  try {
    console.log("\nTest 5: estimateGas with funded account, no gas price");
    const gasEstimate = await ethers.provider.estimateGas({
      from: fundedWallet.address,
      to: fundedWallet.address,
      value: ethers.parseEther("0.001"),
      gasPrice: null
    });
    console.log(`Gas estimate: ${gasEstimate.toString()}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  // Test 6: Estimate gas - Funded account with gas price
  try {
    console.log("\nTest 6: estimateGas with funded account, with gas price");
    const gasEstimate = await ethers.provider.estimateGas({
      from: fundedWallet.address,
      to: fundedWallet.address,
      value: ethers.parseEther("0.001"),
      gasPrice
    });
    console.log(`Gas estimate: ${gasEstimate.toString()}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  // Test 7: Estimate gas - Unfunded account with no gas price
  try {
    console.log("\nTest 7: estimateGas with unfunded account, no gas price");
    const gasEstimate = await ethers.provider.estimateGas({
      from: unfundedWallet.address,
      to: fundedWallet.address,
      value: ethers.parseEther("0.001"),
      gasPrice: null
    });
    console.log(`Gas estimate: ${gasEstimate.toString()}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
  
  // Test 8: Estimate gas - Unfunded account with gas price (should fail)
  try {
    console.log("\nTest 8: estimateGas with unfunded account, with gas price");
    const gasEstimate = await ethers.provider.estimateGas({
      from: unfundedWallet.address,
      to: fundedWallet.address,
      value: ethers.parseEther("0.001"),
      gasPrice
    });
    console.log(`Gas estimate: ${gasEstimate.toString()}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }

  // Test 9: Direct RPC call to eth_estimateGas with unfunded account
  try {
    console.log("\nTest 9: Direct RPC call to eth_estimateGas with unfunded account");
    const result = await ethers.provider.send("eth_estimateGas", [
      {
        from: unfundedWallet.address,
        to: fundedWallet.address,
        value: ethers.parseEther("0.001").toString(),
        gasPrice: gasPrice.toString()
      }
    ]);
    console.log(`Result: ${result}`);
  } catch (error) {
    console.log("Error object structure:");
    console.log(JSON.stringify(error, null, 2));
    console.error(`Error: ${error.message}`);
    
    // If there's an error.data field, it might contain the revert reason
    if (error.data) {
      console.error(`Error data: ${error.data}`);
    }
  }

  // Test 10: Direct RPC call to eth_call with unfunded account
  try {
    console.log("\nTest 10: Direct RPC call to eth_call with unfunded account");
    const result = await ethers.provider.send("eth_call", [
      {
        from: unfundedWallet.address,
        to: fundedWallet.address,
        value: ethers.parseEther("0.001").toString(),
        gasPrice: gasPrice.toString()
      },
      "latest"
    ]);
    console.log(`Result: ${result}`);
  } catch (error) {
    console.log("Error object structure:");
    console.log(JSON.stringify(error, null, 2));
    console.error(`Error: ${error.message}`);
    
    // If there's an error.data field, it might contain the revert reason
    if (error.data) {
      console.error(`Error data: ${error.data}`);
    }
  }
}

// Run the script
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
