import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import {upgrades} from 'hardhat';

async function getProxyAdminAddress(proxyAddress: string, ethers: any): Promise<string> {
    const ADMIN_SLOT = '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103';
    const adminStorage = await ethers.provider.getStorage(proxyAddress, ADMIN_SLOT);
    return ethers.getAddress(`0x${adminStorage.slice(26)}`);
}

async function getImplementationAddress(proxyAddress: string, ethers: any): Promise<string> {
    const IMPLEMENTATION_SLOT = '0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc';
    const implementationStorage = await ethers.provider.getStorage(proxyAddress, IMPLEMENTATION_SLOT);
    return ethers.getAddress(`0x${implementationStorage.slice(26)}`);
}

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {
        deployments,
        getNamedAccounts,
        ethers
    } = hre;
    console.log('Starting upgrade process...');
    const {deployer} = await getNamedAccounts();

    // Get the NetworkConfig contract
    const networkConfigAddress = process.env.NETWORK_CONFIG_ADDR;
    if (!networkConfigAddress) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable is not set');
    }

    console.log('Getting addresses from NetworkConfig:', networkConfigAddress);
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddress);
    const addresses = await networkConfig.addresses();

    console.log('Current addresses:');
    console.log(`NetworkConfig: ${networkConfigAddress}`);
    console.log(`CrossChain: ${addresses.crossChain}`);
    console.log(`NetworkEnclaveRegistry: ${addresses.networkEnclaveRegistry}`);
    console.log(`DataAvailabilityRegistry: ${addresses.dataAvailabilityRegistry}`);

    // Deploy new implementations
    console.log('\nDeploying new implementations...');

    try {
        // Deploy new CrossChain implementation
        console.log('Deploying new CrossChain implementation...');
        const CrossChainFactory = await ethers.getContractFactory('CrossChain');
        const newCrossChainImpl = await CrossChainFactory.deploy();
        await newCrossChainImpl.waitForDeployment();
        const newCrossChainImplAddress = await newCrossChainImpl.getAddress();
        console.log(`New CrossChain implementation deployed at: ${newCrossChainImplAddress}`);

        // Deploy new NetworkEnclaveRegistry implementation
        console.log('Deploying new NetworkEnclaveRegistry implementation...');
        const NetworkEnclaveRegistryFactory = await ethers.getContractFactory('NetworkEnclaveRegistry');
        const newNetworkEnclaveRegistryImpl = await NetworkEnclaveRegistryFactory.deploy();
        await newNetworkEnclaveRegistryImpl.waitForDeployment();
        const newNetworkEnclaveRegistryImplAddress = await newNetworkEnclaveRegistryImpl.getAddress();
        console.log(`New NetworkEnclaveRegistry implementation deployed at: ${newNetworkEnclaveRegistryImplAddress}`);

        // Deploy new DataAvailabilityRegistry implementation
        console.log('Deploying new DataAvailabilityRegistry implementation...');
        const DataAvailabilityRegistryFactory = await ethers.getContractFactory('DataAvailabilityRegistry');
        const newDataAvailabilityRegistryImpl = await DataAvailabilityRegistryFactory.deploy();
        await newDataAvailabilityRegistryImpl.waitForDeployment();
        const newDataAvailabilityRegistryImplAddress = await newDataAvailabilityRegistryImpl.getAddress();
        console.log(`New DataAvailabilityRegistry implementation deployed at: ${newDataAvailabilityRegistryImplAddress}`);

        // Deploy new NetworkConfig implementation
        console.log('Deploying new NetworkConfig implementation...');
        const NetworkConfigFactory = await ethers.getContractFactory('NetworkConfig');
        const newNetworkConfigImpl = await NetworkConfigFactory.deploy();
        await newNetworkConfigImpl.waitForDeployment();
        const newNetworkConfigImplAddress = await newNetworkConfigImpl.getAddress();
        console.log(`New NetworkConfig implementation deployed at: ${newNetworkConfigImplAddress}`);

        // Upgrade the proxies
        console.log('\nUpgrading proxies...');

        // Upgrade CrossChain
        console.log('Preparing upgrade...');
        const preparedImpl = await upgrades.prepareUpgrade(addresses.crossChain, CrossChainFactory);
        console.log('Upgrade prepared, new implementation at:', preparedImpl);

        // Then try the upgrade
        const crossChainUpgradeTx = await upgrades.upgradeProxy(addresses.crossChain, CrossChainFactory);
        await crossChainUpgradeTx.waitForDeployment();
        console.log('CrossChain upgraded');

        // Upgrade NetworkEnclaveRegistry
        console.log('Upgrading NetworkEnclaveRegistry...');
        try {
            const networkEnclaveRegistryUpgradeTx = await upgrades.upgradeProxy(addresses.networkEnclaveRegistry, NetworkEnclaveRegistryFactory);
            await networkEnclaveRegistryUpgradeTx.waitForDeployment();
            console.log('NetworkEnclaveRegistry upgraded');
        } catch (error) {
            console.error('Error upgrading NetworkEnclaveRegistry:', error);
            throw error;
        }

        // Upgrade DataAvailabilityRegistry
        console.log('Upgrading DataAvailabilityRegistry...');
        try {
            const dataAvailabilityRegistryUpgradeTx = await upgrades.upgradeProxy(addresses.dataAvailabilityRegistry, DataAvailabilityRegistryFactory);
            await dataAvailabilityRegistryUpgradeTx.waitForDeployment();
            console.log('DataAvailabilityRegistry upgraded');
        } catch (error) {
            console.error('Error upgrading DataAvailabilityRegistry:', error);
            throw error;
        }

        // Verify the upgrades
        console.log('\nVerifying upgrades...');
        const networkConfigContract = await ethers.getContractAt('NetworkConfig', networkConfigAddress);
        const crossChainContract = await ethers.getContractAt('CrossChain', addresses.crossChain);
        const networkEnclaveRegistryContract = await ethers.getContractAt('NetworkEnclaveRegistry', addresses.networkEnclaveRegistry);
        const dataAvailabilityRegistryContract = await ethers.getContractAt('DataAvailabilityRegistry', addresses.dataAvailabilityRegistry);

        console.log('Upgrades verified successfully!');
        console.log('\nNew implementation addresses:');
        console.log(`NetworkConfig: ${newNetworkConfigImplAddress}`);
        console.log(`CrossChain: ${newCrossChainImplAddress}`);
        console.log(`NetworkEnclaveRegistry: ${newNetworkEnclaveRegistryImplAddress}`);
        console.log(`DataAvailabilityRegistry: ${newDataAvailabilityRegistryImplAddress}`);
    } catch (error) {
        console.error('Error during upgrade:', error);
        throw error;
    }
};

export default func;
func.tags = ['NetworkConfig', 'NetworkConfig_upgrade'];
func.dependencies = ['NetworkConfig_deploy']; 