import { BaseContract } from 'ethers';
import { ethers } from 'hardhat';
import { upgrades } from 'hardhat';
import { UpgradeOptions } from '@openzeppelin/hardhat-upgrades/dist/utils';
import * as path from 'path';
const hre = require("hardhat");

console.log('=== Script started ===');

export async function upgradeContract(
    upgraderAddress: string,
    contractName: string,
    proxyAddress: string
): Promise<BaseContract> {
    console.log(
        `Upgrading proxy ${proxyAddress} to new implementation of ${contractName} (sent from ${upgraderAddress})`
    );
    // hardhat will compile the contract if it's not already compiled
    const factory = await ethers.getContractFactory(contractName);

    // get the current implementation address
    const currentImpl = await hre.upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log(`Current implementation address: ${currentImpl}`);

    // Force import the existing proxy with its current implementation
    await upgrades.forceImport(proxyAddress, factory, {
        kind: 'transparent',
        unsafeAllow: ['constructor'],
        implementation: currentImpl
    } as UpgradeOptions);

    const upgraded = await upgrades.upgradeProxy(proxyAddress, factory, {
        kind: 'transparent',
        unsafeAllow: ['constructor']
    } as UpgradeOptions);

    const address = await upgraded.getAddress();
    console.log(`${contractName} upgraded — new implementation at ${address}`);
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
    const addresses = await networkConfig.addresses();

    console.log('\nCurrent proxy addresses');
    console.table({
        NetworkConfig: networkConfigAddr,
        CrossChain: addresses.crossChain,
        NetworkEnclaveRegistry: addresses.networkEnclaveRegistry,
        DataAvailabilityRegistry: addresses.dataAvailabilityRegistry
    });

    // Perform upgrades
    await upgradeContract(upgrader, 'CrossChain', addresses.crossChain);
    await upgradeContract(upgrader, 'NetworkEnclaveRegistry', addresses.networkEnclaveRegistry);
    await upgradeContract(upgrader, 'DataAvailabilityRegistry', addresses.dataAvailabilityRegistry);
    // await upgradeContract(upgrader, 'NetworkConfig', networkConfigAddr);
    console.log('All upgrades completed successfully');
}

upgradeContracts()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default upgradeContracts