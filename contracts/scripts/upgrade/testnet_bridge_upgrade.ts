import { BaseContract } from 'ethers';
import { ethers } from 'hardhat';
import { upgrades } from 'hardhat';
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

    // Import the existing proxy using the CURRENT contract type
    // The proxy was deployed as TenBridge, so we need to import it as TenBridge first
    console.log('Importing existing proxy (as TenBridge) into OpenZeppelin tracking system...');
    const currentFactory = await ethers.getContractFactory('TenBridge');
    await upgrades.forceImport(proxyAddress, currentFactory, {
        kind: 'transparent'
    } as UpgradeOptions);

    console.log(`Performing upgrade from TenBridge to ${newContractName}...`);
    const upgraded = await upgrades.upgradeProxy(proxyAddress, newFactory, { 
        kind: 'transparent'
    });

    const address = await upgraded.getAddress();
    const newImpl = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    
    if (newImpl === currentImpl) {
        throw new Error(`Upgrade failed: implementation address unchanged (${currentImpl})`);
    }
    
    console.log(`${newContractName} upgraded successfully:`);
    console.log(`  Old implementation: ${currentImpl}`);
    console.log(`  New implementation: ${newImpl}`);
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

    // get addresses from network config
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    if (!networkConfigAddr) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable is not set');
    }

    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const { l1Bridge } = await networkConfig.addresses();

    console.log('\nCurrent proxy addresses');
    console.table({
        NetworkConfig: networkConfigAddr,
        L1Bridge: l1Bridge
    });

    // Perform upgrades
    await upgradeContract('TenBridgeTestnet', l1Bridge);

    console.log('Upgrade completed successfully');
}

upgradeContracts()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default upgradeContracts