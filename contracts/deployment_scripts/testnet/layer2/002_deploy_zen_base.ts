import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { network } from 'hardhat';

/* 
    This script deploys the ZenTestnet contract on the l2 and whitelists it.
*/


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 

    const l2Accounts = await l2Network.getNamedAccounts();

    const networkConfig = await l2Network.network.provider.request({
        method: "net_config",
    });    

    console.log(`net-cfg: ${JSON.stringify(networkConfig, null, " ")}`);

    const zenTestnet = await l2Network.deployments.deploy("ZenTestnet", {
        from: l2Accounts.deployer,
        log: true,
        args: [],
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [networkConfig["TransactionPostProcessorAddress"]]
                }
            }
        }
    });
    console.log(`ZenBase deployed at ${zenTestnet.address}`);

    const signer = await l2Network.ethers.getSigner(l2Accounts.deployer);
    const transactionPostProcessor = await l2Network.ethers.getContractAt(
        'TransactionPostProcessor', 
        networkConfig["TransactionPostProcessorAddress"], 
        signer
    );

    // TODO: add callback with the security epic when we add the EOA config and all the rules for access
    // to system contracts
    /*
    const receipt = await transactionPostProcessor.addOnBlockEndCallback(zenTestnet.address);
    if (receipt.status !== 1) {
        throw new Error("Failed to register Zen token as a system callback");
    }
    console.log(`Callback added at ${receipt.transactionHash}`); */
}
export default func;
func.tags = ['ZenBase', 'ZenBase_deploy'];
func.dependencies = ['EthereumBridge'];
