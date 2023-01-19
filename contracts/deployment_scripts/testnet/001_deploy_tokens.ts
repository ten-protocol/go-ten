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
        contract: "ConstantSupplyERC20",
        args: [ "HOC", "HOC", "1000000000000000000000000000000" ],
        log: true,
    });

    await deployments.deploy('POCERC20', {
        from: deployer,
        contract: "ConstantSupplyERC20",
        args: [ "POC", "POC", "1000000000000000000000000000000" ],
        log: true,
    });
};

export default func;
func.tags = ['HPERC20', 'HPERC20_deploy'];
func.dependencies = ['ObscuroBridge'];
