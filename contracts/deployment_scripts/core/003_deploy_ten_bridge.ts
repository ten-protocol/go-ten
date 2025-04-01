import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { network } from 'hardhat';

/* 
    This deployment script deploys the Ten Bridge smart contracts on both L1 and L2
    and links them together using the 'setRemoteBridge' call.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {
        deployments, 
        getNamedAccounts
    } = hre;
    
    const { deployer } = await getNamedAccounts();
    const messengerL1 = await deployments.get("CrossChainMessenger");
    const networkConfig = await deployments.get("NetworkConfig");
    const networkConfigAddress = networkConfig.address;

    // Deploy the layer 1 part of the bridge
    const layer1BridgeDeployment = await deployments.deploy('TenBridge', {
        from: deployer,
        log: true,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [ messengerL1.address ]
                }
            }
        }
    });

    // Set L1 bridge address in network config
    const networkConfigContract = (await hre.ethers.getContractFactory('NetworkConfig')).attach(networkConfigAddress);
    const recordL1AddressTx = await networkConfigContract.getFunction("setL1BridgeAddress").populateTransaction(layer1BridgeDeployment.address);
    await deployments.rawTx({
        from: deployer,
        to: networkConfigAddress,
        data: recordL1AddressTx.data,
        log: true,
        waitConfirmations: 1,
    });
    console.log(`L1BridgeAddress=${layer1BridgeDeployment.address}`);
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];
func.dependencies = ['CrossChainMessenger'];
