import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    // L2 address of a prefunded deployer account to be used in smart contracts
    const {deployer} = await getNamedAccounts();

    // We get the Obscuro Bridge deployment on the layer 1 network.
    const layer1BridgeDeployment = await hre.companionNetworks.layer1.deployments.get("ObscuroBridge");

    // We get the Cross chain messenger deployment on the layer 2 network.
    const messengerDeployment = await deployments.get("CrossChainMessenger");

    // Deploy the Ethereum Bridge and instruct it to use the address of the L2 cross chain messenger to enable functionality
    // and be subordinate of the L1 ObscuroBridge
    await deployments.deploy('EthereumBridge', {
        from: deployer,
        args: [ messengerDeployment.address, layer1BridgeDeployment.address ],
        log: true,
    });
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];

// This should only be deployed after the L2 CrossChainMessenger
func.dependencies = ['CrossChainMessenger'];
