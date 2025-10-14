import { ethers } from 'hardhat';

console.log('=== Recover Testnet Funds Script ===');

const recoverFunds = async function (): Promise<void> {
    const [deployer] = await ethers.getSigners();
    console.log(`Using signer: ${deployer.address}`);

    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    if (!networkConfigAddr) {
        throw new Error("NETWORK_CONFIG_ADDR environment variable is not set");
    }
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr)
    const addresses = await networkConfig.addresses()
    const bridgeAddress = addresses.l1Bridge
    console.log(`Bridge address: ${bridgeAddress}`);

    // Gnosis sepolia ETH address
    const receiverAddress = process.env.RECEIVER_ADDRESS || "0xea052c9635f1647a8a199c2315b9a66ce7d1e2a7"
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

