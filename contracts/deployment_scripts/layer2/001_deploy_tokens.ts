import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {
        deployer,
        hocowner,
        pocowner,
    } = await getNamedAccounts();

    console.log(`HOC = ${hocowner}  POC=${pocowner}`);

    await hre.run('obscuro:wallet-extension:add-key', {address: hocowner});
    await hre.run('obscuro:wallet-extension:add-key', {address: pocowner});

    await deployments.deploy('L2HOCERC20', {
        from: hocowner,
        contract: "WrappedERC20",
        args: [ "HOC", "HOC" ],
        log: true
    });

    await deployments.execute('L2HOCERC20', {
        from: hocowner,
        log: true
    }, "issueFor", hocowner, "1000000000000000000000000000000");

    await deployments.deploy('L2POCERC20', {
        from: pocowner,
        contract: "WrappedERC20",
        args: [ "POC", "POC" ],
        log: true
    });

    await deployments.execute('L2POCERC20', {
        from: pocowner,
        log: true
    }, "issueFor", pocowner, "1000000000000000000000000000000");
};

export default func;
func.tags = ['L2HPERC20', 'L2HPERC20_deploy'];
