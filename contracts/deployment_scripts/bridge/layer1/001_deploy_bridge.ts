import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    const messengerDeployment = await deployments.get("CrossChainMessenger");

    await deployments.deploy('ObscuroBridge', {
        from: deployer,
        args: [ messengerDeployment.address ],
        log: true,
    });
};

export default func;
func.tags = ['ObscuroBridge', 'ObscuroBridge_deploy'];
func.dependencies = ['CrossChainMessenger'];
