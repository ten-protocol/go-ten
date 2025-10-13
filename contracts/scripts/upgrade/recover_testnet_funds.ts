import { ethers } from 'hardhat';

console.log('=== Recover Testnet Funds Script ===');

const recoverFunds = async function (): Promise<void> {
    const [deployer] = await ethers.getSigners();
    console.log(`Using signer: ${deployer.address}`);

    // Get the bridge address from environment or fetch from NetworkConfig
    let bridgeAddress = process.env.BRIDGE_ADDRESS
  
    if (!bridgeAddress) {
        console.log("BRIDGE_ADDRESS not provided, fetching from NetworkConfig...")
        const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
        if (!networkConfigAddr) {
            throw new Error("NETWORK_CONFIG_ADDR environment variable is not set");
        }
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr)
    const addresses = await networkConfig.addresses()
    bridgeAddress = addresses.l1Bridge
  }

    // Receiver address - defaults to the deployer if not specified
    const receiverAddress = process.env.RECEIVER_ADDRESS || deployer.address;
    console.log(`Receiver address: ${receiverAddress}`);

    // Get the contract
    const bridge = await ethers.getContractAt('TenBridgeTestnet', bridgeAddress);

    // Check current balance
    const currentBalance = await ethers.provider.getBalance(bridgeAddress);
    console.log(`\nCurrent bridge balance: ${ethers.formatEther(currentBalance)} ETH`);

    if (currentBalance === 0n) {
        console.log('Bridge has no funds to recover.');
        return;
    }

    // Call recoverTestnetFunds
    console.log(`\nCalling recoverTestnetFunds to send funds to ${receiverAddress}...`);
    const tx = await bridge.recoverTestnetFunds(receiverAddress);
    console.log(`Transaction hash: ${tx.hash}`);
    
    console.log('Waiting for confirmation...');
    const receipt = await tx.wait();
    console.log(`Transaction confirmed in block ${receipt.blockNumber}`);

    // Check new balance
    const newBalance = await ethers.provider.getBalance(bridgeAddress);
    console.log(`\nNew bridge balance: ${ethers.formatEther(newBalance)} ETH`);
    console.log(`Recovered: ${ethers.formatEther(currentBalance - newBalance)} ETH`);
    
    console.log('\n=== Recovery completed successfully ===');
}

recoverFunds()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error('Recovery failed:', error);
        process.exit(1);
    });

export default recoverFunds;

