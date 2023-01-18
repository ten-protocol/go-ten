import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    // Read the message bus address from the management contract deployment.
    const messageBusAddress : string = await deployments.read("ManagementContract", {}, "messageBus");

    // Setup the cross chain messenger and point it to the message bus from the management contract to be used for validation
    await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        args: [ messageBusAddress ],
        log: true,
    });
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
func.dependencies = ['ManagementContract', 'HPERC20']; //TODO: Remove HPERC20, this is only to have matching addresses.
