import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    const layer1BridgeDeployment = await hre.companionNetworks.layer1.deployments.get("ObscuroBridge");
    const messengerDeployment = await deployments.get("CrossChainMessenger");

    await deployments.deploy('ObscuroL2Bridge', {
        from: deployer,
        args: [ messengerDeployment.address, layer1BridgeDeployment.address ],
        log: true,
    });
};

export default func;
func.tags = ['ObscuroL2Bridge', 'ObscuroL2Bridge_deploy'];
func.dependencies = ['CrossChainMessenger'];
