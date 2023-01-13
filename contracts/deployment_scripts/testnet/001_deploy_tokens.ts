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

    await deployments.deploy('POCERC20', {
        from: deployer,
        contract: "WrappedERC20",
        args: [ "POC", "POC" ],
        log: true,
    });


    await deployments.execute('HOCERC20', {
        from: deployer,
        log: true,
    }, "issueFor", deployer, "1000000000000000000000000000000");

    await deployments.execute('POCERC20', {
        from: deployer,
        log: true, 
    }, "issueFor", deployer, "1000000000000000000000000000000");
};

export default func;
func.tags = ['HPERC20', 'HPERC20_deploy'];
func.dependencies = ['ObscuroBridge'];
