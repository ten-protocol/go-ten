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
    await deployments.execute('ObscuroBridge', {
        from: deployer
    }, 'promoteToAdmin', '0xDEe530E22045939e6f6a0A593F829e35A140D3F1')
};

export default func;
func.tags = ['BRIDGEADMIN', 'BRIDGEADMIN_deploy'];
func.dependencies = ['ObscuroBridge']
