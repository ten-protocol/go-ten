import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    await deployments.deploy('POCERC20', {
        from: deployer,
        contract: "WrappedERC20",
        args: [ "POC", "POC" ],
        log: true,
    });
};

export default func;
func.tags = ['POCERC20', 'POCERC20_deploy'];
func.dependencies = ['ObscuroBridge'];
