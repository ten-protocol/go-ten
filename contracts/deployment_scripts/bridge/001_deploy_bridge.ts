import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


/* 
    This deployment script deploys the Obscuro Bridge smart contracts on both L1 and L2
    and links them together using the 'setRemoteBridge' call.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    // L2 address of a prefunded deployer account to be used in smart contracts
    const accountsL2 = await getNamedAccounts();
    const accountsL1 = await hre.companionNetworks.layer1.getNamedAccounts();

    // L1 Cross Chain Messenger Deployment.
    const messengerL1 = await hre.companionNetworks.layer1.deployments.get("CrossChainMessenger");

    // We deploy the layer 1 part of the bridge.
    const layer1BridgeDeployment = await hre.companionNetworks.layer1.deployments.deploy('ObscuroBridge', {
        from: accountsL1.deployer,
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

    // We get the Cross chain messenger deployment on the layer 2 network.
    const messengerL2 = await deployments.get("CrossChainMessenger");

    // Deploy the layer 2 part of the bridge and instruct it to use the address of the L2 cross chain messenger to enable functionality
    // and be subordinate of the L1 ObscuroBridge
    const layer2BridgeDeployment = await deployments.deploy('EthereumBridge', {
        from: accountsL2.deployer,
        log: true,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [ messengerL2.address, layer1BridgeDeployment.address ]
                }
            }
        }
    });

    await hre.companionNetworks.layer1.deployments.execute("ObscuroBridge", {
        from: accountsL1.deployer, 
        log: true,
    }, "setRemoteBridge", layer2BridgeDeployment.address);

    console.log(` Bridge deployed with from L1 address=${accountsL1.deployer} L2 Address=${accountsL2.deployer}`);
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];

// This should only be deployed after the L2 CrossChainMessenger
func.dependencies = ['CrossChainMessenger'];
