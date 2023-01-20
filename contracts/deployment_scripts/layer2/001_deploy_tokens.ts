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

    // Create the keys for the hoc and poc owner accounts.
    console.log(`Adding VK for HOC Owner = ${hocowner}`);
    await hre.run('obscuro:wallet-extension:add-key', {address: hocowner});

    console.log(`Adding VK for POC owner = ${pocowner}`);
    await hre.run('obscuro:wallet-extension:add-key', {address: pocowner});

    console.log(`Added keys!`);

    await deployments.deploy('L2HOCERC20', {
        from: hocowner,
        contract: "WrappedERC20",
        args: [ "HOC", "HOC" ],
        log: true
    });

    // Mint the initial supply. This is different to the older smart contracts that had a
    // static supply minted on contract creation
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

    // Mint initial supply for POC too.
    await deployments.execute('L2POCERC20', {
        from: pocowner,
        log: true
    }, "issueFor", pocowner, "1000000000000000000000000000000");
};

export default func;
func.tags = ['L2HPERC20', 'L2HPERC20_deploy'];
