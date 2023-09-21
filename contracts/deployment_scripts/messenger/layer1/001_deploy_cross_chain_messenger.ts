import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import 'process';

/*
    This script deploys the CrossChainMessenger contract on the L1. 
    It depends on knowing the address of the message bus from the management contract predeployment.
    This is passed using the environment variables.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const deployments = hre.companionNetworks.layer1.deployments;

    const { deployer } = await hre.companionNetworks.layer1.getNamedAccounts();

    // Read the message bus address from the management contract deployment.
    const messageBusAddress : string = process.env.MESSAGE_BUS_ADDRESS || "0xa1fdA5f6Df55a326f5f4300F3A716317f0f03110"
    console.log(`Message Bus address ${messageBusAddress}`);

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
