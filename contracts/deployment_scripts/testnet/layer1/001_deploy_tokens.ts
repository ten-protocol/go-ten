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
    await deployments.deploy('HOCERC20', {
        from: deployer,
        contract: "ConstantSupplyERC20",
        args: [ "HOC", "HOC", "1000000000000000000000000000000" ],
        log: true,
    });

    // Deploy a constant supply (constructor mints) erc20
    await deployments.deploy('POCERC20', {
        from: deployer,
        contract: "ConstantSupplyERC20",
        args: [ "POC", "POC", "1000000000000000000000000000000" ],
        log: true,
    });
};

export default func;
func.tags = ['HPERC20', 'HPERC20_deploy'];
