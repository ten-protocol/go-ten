import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    const messageBusAddress : string = await deployments.read("ManagementContract", {}, "messageBus");

    await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        args: [ messageBusAddress ],
        log: true,
    });
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
func.dependencies = ['ManagementContract', 'HPERC20']; //TODO: Remove HPERC20, this is only to have matching addresses.
