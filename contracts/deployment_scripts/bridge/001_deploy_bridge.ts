import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { network } from 'hardhat';


/* 
    This deployment script deploys the Obscuro Bridge smart contracts on both L1 and L2
    and links them together using the 'setRemoteBridge' call.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    // Make sure the companion network exists
    if (!hre.companionNetworks || !hre.companionNetworks.layer1) {
        throw new Error("Companion network 'layer1' is not configured. Please check your hardhat config.");
    }

    const { deployer } = await hre.companionNetworks.layer1.getNamedAccounts();
    if (!deployer) {
        throw new Error("Deployer account not found in named accounts");
    }

    console.log("Retrieving network configuration...");
    const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
    console.log("Network config:", JSON.stringify(networkConfig, null, 2));

    // Check if L1Bridge exists
    if (!networkConfig || !networkConfig.L1Bridge) {
        throw new Error("L1Bridge not found in network config. Check that the bridge is deployed on L1.");
    }

    const l1BridgeAddress = networkConfig.L1Bridge;
    const l2BridgeAddress = networkConfig.L2Bridge;
    console.log(`Using L1 bridge address: ${l1BridgeAddress}`);
    console.log(`Using L2 bridge address: ${l2BridgeAddress}`);

    // Check if L2BridgeAddress exists
    if (!l2BridgeAddress) {
        throw new Error("L2Bridge not found in network config. Check that the bridge is deployed on L2.");
    }

    // Get the contract factory and attach to the bridge address
    console.log("Attaching L1 bridge to TenBridge contract...");
    const bridgeContract = (await hre.ethers.getContractFactory('TenBridge')).attach(l1BridgeAddress);
    
    // Create transaction to set remote bridge
    console.log("Creating setRemoteBridge transaction...");
    const tx = await bridgeContract.getFunction("setRemoteBridge").populateTransaction(l2BridgeAddress);
    
    // Execute transaction via layer1 network
    console.log("Sending transaction to set remote bridge...");
    const receipt = await hre.companionNetworks.layer1.deployments.rawTx({
        from: deployer,
        to: l1BridgeAddress,
        data: tx.data,
        log: true,
        waitConfirmations: 1,
    });
    
    if (!receipt.events || receipt.events.length === 0) {
        console.log(`Warning: No events emitted. This might indicate a failure when setting L2Bridge=${l2BridgeAddress} on bridge contract.`);
    } else {
        console.log(`Successfully set L2BridgeAddress=${l2BridgeAddress} on bridge contract`);
    }
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];