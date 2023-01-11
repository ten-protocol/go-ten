import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    await deployments.deploy('HOCERC20', {
        from: deployer,
        contract: "WrappedERC20",
        args: [ "HOC", "HOC" ],
        log: true,
    });
};

export default func;
func.tags = ['HOCERC20', 'HOCERC20_deploy'];
func.dependencies = ['ObscuroBridge'];
