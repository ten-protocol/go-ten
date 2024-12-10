import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script instantiates the L1 side of the HOC and POC tokens.
    It is equivalent to what the old contract deployer was doing, except for
    address prefunding.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre.companionNetworks.layer1;

    const {deployer} = await getNamedAccounts();

    // Deploy a constant supply (constructor mints) erc20
    await deployments.execute('TenBridge', {
        from: deployer
    }, 'promoteToAdmin', '0xE09a37ABc1A63441404007019E5BC7517bE2c43f');

    await deployments.execute('TenBridge', {
        from: deployer
    }, 'promoteToAdmin', '0xeC3f9B38a3B30AdC9fB3dF3a0D8f50127E6c2C8f');
};

export default func;
func.tags = ['BRIDGEADMIN', 'BRIDGEADMIN_deploy'];
func.dependencies = ['TenBridge']
