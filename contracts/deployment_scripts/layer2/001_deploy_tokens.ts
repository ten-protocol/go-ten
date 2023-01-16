import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    await deployments.deploy('L2HOCERC20', {
        from: deployer,
        contract: "WrappedERC20",
        args: [ "HOC", "HOC", "1000000000000000000000000000000" ],
        log: true,
    });

    await deployments.deploy('L2POCERC20', {
        from: deployer,
        contract: "WrappedERC20",
        args: [ "POC", "POC" ],
        log: true,
    });

    await deployments.execute('L2HOCERC20', {
        from: deployer
    }, "issueFor", deployer, "1000000000000000000000000000000");


    await deployments.execute('L2POCERC20', {
        from: deployer
    }, "issueFor", deployer, "1000000000000000000000000000000");
};

export default func;
func.tags = ['L2HPERC20', 'L2HPERC20_deploy'];
func.dependencies = ['ObscuroBridge'];
