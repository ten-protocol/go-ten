import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { network } from 'hardhat';

/* 
    This script deploys the ZenTestnet contract on the l2 and whitelists it.
*/

// Define type for the network config
interface NetworkConfig {
    NetworkConfig: string;
    EnclaveRegistry: string;
    DataAvailabilityRegistry: string;
    CrossChain: string;
    L1MessageBus: string;
    L2MessageBus: string;
    TransactionsPostProcessor: string;
    L1Bridge: string;
    L2Bridge: string;
    L1CrossChainMessenger: string;
    L2CrossChainMessenger: string;
    L1StartHash: string;
    PublicSystemContracts: {
        CrossChainMessenger: string;
        EthereumBridge: string;
        Fees: string;
        MessageBus: string;
        PublicCallbacks: string;
        TransactionsPostProcessor: string;
    };
    AdditionalContracts: any;
}

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 

    const l2Accounts = await l2Network.getNamedAccounts();
    
    // Ensure deployer exists
    if (!l2Accounts.deployer) {
        throw new Error("Deployer account not found in named accounts");
    }

    // Type cast the network config
    const networkConfig = await l2Network.network.provider.request({
        method: "net_config",
    }) as NetworkConfig;    

    // Verify the TransactionsPostProcessor address exists
    if (!networkConfig.TransactionsPostProcessor) {
        throw new Error("TransactionsPostProcessor address not found in network config");
    }

    console.log(`Trying to deploy ZenBase using: ${l2Accounts.deployer}`);

    const zenTestnet = await l2Network.deployments.deploy("ZenTestnet", {
        from: l2Accounts.deployer,
        log: true,
        args: [],
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [networkConfig.TransactionsPostProcessor]
                }
            }
        }
    });
    console.log(`ZenBase deployed at ${zenTestnet.address}`);

    const signer = await l2Network.ethers.getSigner(l2Accounts.deployer);
    const transactionPostProcessor = await l2Network.ethers.getContractAt(
        'TransactionPostProcessor', 
        networkConfig.TransactionsPostProcessor,
        signer
    );
    
    const tx = await transactionPostProcessor.addOnBlockEndCallback(zenTestnet.address);
    const receipt = await tx.wait();
    if (!receipt || receipt.status !== 1) {
        throw new Error("Failed to register Zen token as a system callback");
    }
    console.log(`Callback added at ${receipt.hash}`);
}
export default func;
func.tags = ['ZenBase', 'ZenBase_deploy'];
func.dependencies = ['EthereumBridge'];
