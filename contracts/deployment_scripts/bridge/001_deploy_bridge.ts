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
    var networkConfigAddress = process.env.NETWORK_CONFIG_ADDR!!
    if (networkConfigAddress === undefined) {
        const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
        networkConfigAddress = networkConfig.NetworkConfig;
        console.log(`Fallback read of network config contract address = ${networkConfigAddress}`);
    }

    // L2 address of a prefunded deployer account to be used in smart contracts
    const accountsL2 = await getNamedAccounts();
    const accountsL1 = await hre.companionNetworks.layer1.getNamedAccounts();

    // L1 Cross Chain Messenger Deployment.
    const messengerL1 = await hre.companionNetworks.layer1.deployments.get("CrossChainMessenger");

    // We deploy the layer 1 part of the bridge.
    const layer1BridgeDeployment = await hre.companionNetworks.layer1.deployments.deploy('TenBridge', {
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

    // get management contract and write the L1 bridge address to it
    const networkConfigContract = (await hre.ethers.getContractFactory('NetworkConfig')).attach(networkConfigAddress)
    const recordL1AddressTx = await networkConfigContract.getFunction("setL1BridgeAddress").populateTransaction(layer1BridgeDeployment.address);
    const receipt = await hre.companionNetworks.layer1.deployments.rawTx({
        from: accountsL1.deployer,
        to: networkConfigAddress,
        data: recordL1AddressTx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt.events?.length === 0) {
        console.log(`Failed to set L1BridgeAddress=${layer1BridgeDeployment.address} on network config contract.`);
    } else {
        console.log(`L1BridgeAddress=${layer1BridgeDeployment.address}`);
    }

    // We get the Cross chain messenger deployment on the layer 2 network.
    const messengerL2 = await deployments.get("CrossChainMessenger");

    // Deploy the layer 2 part of the bridge and instruct it to use the address of the L2 cross chain messenger to enable functionality
    // and be subordinate of the L1 TenBridge
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

    await hre.companionNetworks.layer1.deployments.execute("TenBridge", {
        from: accountsL1.deployer, 
        log: true,
    }, "setRemoteBridge", layer2BridgeDeployment.address);

    const recordL2AddressTx = await networkConfigContract.getFunction("setL2BridgeAddress").populateTransaction(layer2BridgeDeployment.address);
    const receipt2 = await hre.companionNetworks.layer1.deployments.rawTx({
        from: accountsL1.deployer,
        to: networkConfigAddress,
        data: recordL2AddressTx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt2.events?.length === 0) {
        console.log(`Failed to set L2BridgeAddress=${layer2BridgeDeployment.address} on network config contract.`);
    } else {
        console.log(`L2BridgeAddress=${layer2BridgeDeployment.address}`);
    }
    console.log(` Bridge deployed with from L1 address=${accountsL1.deployer} L2 Address=${accountsL2.deployer}`); 
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];

// This should only be deployed after the L2 CrossChainMessenger
func.dependencies = ['CrossChainMessenger'];
