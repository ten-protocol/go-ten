import { BaseContract } from 'ethers';
import { ethers } from 'hardhat';
import { upgrades } from 'hardhat';
import { UpgradeOptions } from '@openzeppelin/hardhat-upgrades/dist/utils';

console.log('=== Script started ===');

export async function upgradeContract(
    upgraderAddress: string,
    contractName: string,
    proxyAddress: string
): Promise<BaseContract> {
    console.log(
        `Upgrading proxy ${proxyAddress} to new implementation of ${contractName} (sent from ${upgraderAddress})`
    );
    // Assumes the contract is already compiled, otherwise ensure `npx hardhat compile` is run first
    const factory = await ethers.getContractFactory(contractName);
    
    // get the current implementation address
    const currentImpl = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log(`Current implementation address: ${currentImpl}`);

    // Force import the existing proxy (let OpenZeppelin detect current implementation)
    await upgrades.forceImport(proxyAddress, factory, {
        kind: 'transparent',
        unsafeAllow: ['constructor']
    } as UpgradeOptions);

    const upgraded = await upgrades.upgradeProxy(proxyAddress, factory, { 
        kind: 'transparent',
        unsafeAllow: ['constructor']    
    });

    const address = await upgraded.getAddress();
    const newImpl = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    
    if (newImpl === currentImpl) {
        throw new Error(`Upgrade failed: implementation address unchanged (${currentImpl})`);
    }
    
    console.log(`${contractName} upgraded successfully:`);
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
    await upgradeContract(upgrader, 'TenBridgeTestnet', l1Bridge);

    console.log('Upgrade completed successfully');
}

upgradeContracts()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default upgradeContracts