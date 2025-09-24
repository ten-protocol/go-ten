import { ethers } from 'hardhat';
import { upgrades } from 'hardhat';

console.log('=== TenBridgeTestnet Deployment Script Started ===');

export async function deployTenBridgeTestnet(
    messengerAddress: string,
    ownerAddress: string
): Promise<string> {
    console.log(`Deploying TenBridgeTestnet with:`);
    console.log(`  Messenger: ${messengerAddress}`);
    console.log(`  Owner: ${ownerAddress}`);
    
    // Get the contract factory
    const TenBridgeTestnetFactory = await ethers.getContractFactory('TenBridgeTestnet');
    
    // Deploy the implementation contract
    console.log('Deploying TenBridgeTestnet implementation...');
    const implementation = await TenBridgeTestnetFactory.deploy();
    await implementation.waitForDeployment();
    const implementationAddress = await implementation.getAddress();
    console.log(`Implementation deployed at: ${implementationAddress}`);
    
    // Deploy the proxy
    console.log('Deploying proxy...');
    const proxy = await upgrades.deployProxy(TenBridgeTestnetFactory, [messengerAddress, ownerAddress], {
        kind: 'transparent',
        initializer: 'initialize'
    });
    await proxy.waitForDeployment();
    const proxyAddress = await proxy.getAddress();
    
    console.log(`TenBridgeTestnet deployed successfully:`);
    console.log(`  Implementation: ${implementationAddress}`);
    console.log(`  Proxy: ${proxyAddress}`);
    
    return proxyAddress;
}

const deployContract = async function (): Promise<void> {
    console.log('=== Starting TenBridgeTestnet deployment ===');
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    const deployerAddress = deployer.address;
    console.log(`Using signer: ${deployerAddress}`);

    // Get parameters from command line or use defaults
    const messengerAddress = process.argv[2];
    const ownerAddress = process.argv[3] || deployerAddress;

    if (!messengerAddress) {
        console.log('Usage: npx hardhat run scripts/deploy/deploy_ten_bridge_testnet.ts --network <network> <messengerAddress> [ownerAddress]');
        console.log('Example: npx hardhat run scripts/deploy/deploy_ten_bridge_testnet.ts --network localGeth 0x1234... 0x5678...');
        throw new Error('Messenger address is required as first parameter');
    }

    // Deploy the contract
    const contractAddress = await deployTenBridgeTestnet(messengerAddress, ownerAddress);
    
    console.log('\n=== Deployment Summary ===');
    console.log(`TenBridgeTestnet deployed at: ${contractAddress}`);
    console.log(`Owner: ${ownerAddress}`);
    console.log(`Messenger: ${messengerAddress}`);
    
    // Optional: Update NetworkConfig if NETWORK_CONFIG_ADDR is provided
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    if (networkConfigAddr) {
        console.log('\n=== Updating NetworkConfig ===');
        try {
            const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
            const tx = await networkConfig.setL1BridgeAddress(contractAddress);
            await tx.wait();
            console.log(`NetworkConfig updated with new bridge address: ${contractAddress}`);
        } catch (error) {
            console.warn('Failed to update NetworkConfig:', error);
        }
    }
    
    console.log('Deployment completed successfully!');
}

deployContract()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default deployContract

