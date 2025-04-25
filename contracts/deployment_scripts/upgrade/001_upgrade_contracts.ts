import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {
        deployments,
        getNamedAccounts,
        ethers
    } = hre;
    const {deployer} = await getNamedAccounts();

    // Get the current proxy addresses
    const networkConfig = await deployments.get('NetworkConfig');
    const crossChain = await deployments.get('CrossChain');
    const networkEnclaveRegistry = await deployments.get('NetworkEnclaveRegistry');
    const dataAvailabilityRegistry = await deployments.get('DataAvailabilityRegistry');

    console.log('Current addresses:');
    console.log(`NetworkConfig: ${networkConfig.address}`);
    console.log(`CrossChain: ${crossChain.address}`);
    console.log(`NetworkEnclaveRegistry: ${networkEnclaveRegistry.address}`);
    console.log(`DataAvailabilityRegistry: ${dataAvailabilityRegistry.address}`);

    // Deploy new implementations
    console.log('\nDeploying new implementations...');

    // Deploy new CrossChain implementation
    const CrossChainFactory = await ethers.getContractFactory('CrossChain');
    const newCrossChainImpl = await CrossChainFactory.deploy();
    console.log(`New CrossChain implementation deployed at: ${newCrossChainImpl.target}`);

    // Deploy new NetworkEnclaveRegistry implementation
    const NetworkEnclaveRegistryFactory = await ethers.getContractFactory('NetworkEnclaveRegistry');
    const newNetworkEnclaveRegistryImpl = await NetworkEnclaveRegistryFactory.deploy();
    console.log(`New NetworkEnclaveRegistry implementation deployed at: ${newNetworkEnclaveRegistryImpl.target}`);

    // Deploy new DataAvailabilityRegistry implementation
    const DataAvailabilityRegistryFactory = await ethers.getContractFactory('DataAvailabilityRegistry');
    const newDataAvailabilityRegistryImpl = await DataAvailabilityRegistryFactory.deploy();
    console.log(`New DataAvailabilityRegistry implementation deployed at: ${newDataAvailabilityRegistryImpl.target}`);

    // Deploy new NetworkConfig implementation
    const NetworkConfigFactory = await ethers.getContractFactory('NetworkConfig');
    const newNetworkConfigImpl = await NetworkConfigFactory.deploy();
    console.log(`New NetworkConfig implementation deployed at: ${newNetworkConfigImpl.target}`);

    // Get the ProxyAdmin contract
    const proxyAdmin = await ethers.getContractAt('ProxyAdmin', await deployments.read('NetworkConfig', 'getProxyAdmin'));

    // Upgrade the proxies
    console.log('\nUpgrading proxies...');

    // Upgrade CrossChain
    const crossChainInitData = CrossChainFactory.interface.encodeFunctionData('initialize', [deployer]);
    await proxyAdmin.upgradeAndCall(crossChain.address, newCrossChainImpl.target, crossChainInitData);
    console.log('CrossChain upgraded');

    // Upgrade NetworkEnclaveRegistry
    const networkEnclaveRegistryInitData = NetworkEnclaveRegistryFactory.interface.encodeFunctionData('initialize', [deployer]);
    await proxyAdmin.upgradeAndCall(networkEnclaveRegistry.address, newNetworkEnclaveRegistryImpl.target, networkEnclaveRegistryInitData);
    console.log('NetworkEnclaveRegistry upgraded');

    // Upgrade DataAvailabilityRegistry
    const merkleMessageBus = await deployments.read('CrossChain', 'merkleMessageBus');
    const dataAvailabilityRegistryInitData = DataAvailabilityRegistryFactory.interface.encodeFunctionData('initialize', 
        [merkleMessageBus, newNetworkEnclaveRegistryImpl.target, deployer]);
    await proxyAdmin.upgradeAndCall(dataAvailabilityRegistry.address, newDataAvailabilityRegistryImpl.target, dataAvailabilityRegistryInitData);
    console.log('DataAvailabilityRegistry upgraded');

    // Upgrade NetworkConfig
    const networkConfigInitData = NetworkConfigFactory.interface.encodeFunctionData('initialize', [{
        crossChain: crossChain.address,
        messageBus: merkleMessageBus,
        networkEnclaveRegistry: networkEnclaveRegistry.address,
        dataAvailabilityRegistry: dataAvailabilityRegistry.address
    }, deployer]);
    await proxyAdmin.upgradeAndCall(networkConfig.address, newNetworkConfigImpl.target, networkConfigInitData);
    console.log('NetworkConfig upgraded');

    // Verify the upgrades
    console.log('\nVerifying upgrades...');
    const networkConfigContract = await ethers.getContractAt('NetworkConfig', networkConfig.address);
    const crossChainContract = await ethers.getContractAt('CrossChain', crossChain.address);
    const networkEnclaveRegistryContract = await ethers.getContractAt('NetworkEnclaveRegistry', networkEnclaveRegistry.address);
    const dataAvailabilityRegistryContract = await ethers.getContractAt('DataAvailabilityRegistry', dataAvailabilityRegistry.address);

    console.log('Upgrades verified successfully!');
    console.log('\nNew implementation addresses:');
    console.log(`NetworkConfig: ${newNetworkConfigImpl.target}`);
    console.log(`CrossChain: ${newCrossChainImpl.target}`);
    console.log(`NetworkEnclaveRegistry: ${newNetworkEnclaveRegistryImpl.target}`);
    console.log(`DataAvailabilityRegistry: ${newDataAvailabilityRegistryImpl.target}`);
};

export default func;
func.tags = ['NetworkConfig', 'NetworkConfig_upgrade'];
func.dependencies = ['NetworkConfig_deploy']; 