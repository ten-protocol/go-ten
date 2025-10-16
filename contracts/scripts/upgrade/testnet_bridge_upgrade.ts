import { BaseContract } from 'ethers';
import { ethers, upgrades } from 'hardhat';
import { UpgradeOptions } from '@openzeppelin/hardhat-upgrades/dist/utils';

console.log('=== Script started ===');

export async function upgradeContract(
    newContractName: string,
    proxyAddress: string
): Promise<BaseContract> {
    console.log(`Upgrading proxy ${proxyAddress} to new implementation of ${newContractName}`);
    
    // Assumes the contract is already compiled, otherwise ensure `npx hardhat compile` is run first
    const newFactory = await ethers.getContractFactory(newContractName);
    
    // get the current implementation address
    const currentImpl = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log(`Current implementation address: ${currentImpl}`);

    // Import proxy and upgrade to new implementation
    console.log('Importing existing proxy into OpenZeppelin tracking system...');
    const currentFactory = await ethers.getContractFactory('TenBridge');
    await upgrades.forceImport(proxyAddress, currentFactory, {
        kind: 'transparent',
        unsafeAllow: ['constructor'],
        implementation: currentImpl
    } as UpgradeOptions);

    console.log(`Deploying new implementation for ${newContractName}...`);
    const newImpl = await upgrades.deployImplementation(newFactory, {
        kind: 'transparent',
        unsafeAllow: ['constructor']
    });
    console.log(`New implementation deployed at: ${newImpl}`);
    
    console.log(`Performing upgrade to ${newContractName}...`);
    const upgraded = await upgrades.upgradeProxy(proxyAddress, newFactory, { 
        kind: 'transparent',
        unsafeAllow: ['constructor']
    });

    const address = await upgraded.getAddress();
    const newImplFinal = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    
    if (newImplFinal === currentImpl) {
        console.log(`Warning: Implementation address unchanged. This may be expected if no code changes were made.`);
        console.log(`Current implementation: ${currentImpl}`);
        console.log(`Deployed implementation: ${newImpl}`);
    }
    
    console.log(`${newContractName} upgraded successfully:`);
    console.log(`  Old implementation: ${currentImpl}`);
    console.log(`  New implementation: ${newImplFinal}`);
    console.log(`  Proxy address: ${address}`);
    return upgraded;
}

const upgradeContracts = async function (): Promise<void> {
    console.log('=== Starting upgrade process ===');
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    const upgrader = deployer.address;
    console.log(`Using signer: ${upgrader}`);

    // Get parameters from command line or use defaults
    const newContractName = process.argv[2] || 'TenBridgeTestnet';

    // get addresses from network config
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    if (!networkConfigAddr) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable is not set');
    }

    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const { l1Bridge } = await networkConfig.addresses();

    // Use provided proxy address or default to l1Bridge
    const targetProxyAddress = l1Bridge;

    console.log('\nCurrent proxy addresses');
    console.table({
        NetworkConfig: networkConfigAddr,
        L1Bridge: l1Bridge,
        TargetProxy: targetProxyAddress
    });

    // Perform upgrades
    await upgradeContract(newContractName, targetProxyAddress);

    console.log('Upgrade completed successfully');
}

upgradeContracts()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default upgradeContracts