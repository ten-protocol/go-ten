import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script adds admin addresses to the TenBridge contract.
    These admins will have permission to perform administrative actions on the bridge.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    // Deploy a constant supply (constructor mints) erc20
    await deployments.execute('TenBridge', {
        from: deployer
    }, 'promoteToAdmin', '0xE09a37ABc1A63441404007019E5BC7517bE2c43f');

    await deployments.execute('TenBridge', {
        from: deployer
    }, 'promoteToAdmin', '0xeC3f9B38a3B30AdC9fB3dF3a0D8f50127E6c2C8f');

    await deployments.execute('TenBridge', {
        from: deployer
    }, 'promoteToAdmin', '0x7e6866FCA4A210913b668b46d0E618e437F5734b');
};

export default func;
func.tags = ['BRIDGEADMIN', 'BRIDGEADMIN_deploy'];
func.dependencies = ['TenBridge']
