import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { network } from 'hardhat';


/* 
    This deployment script deploys the TEN Bridge smart contracts on both L1 and L2
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
    const l2CrossChainMessengerAddress = networkConfig.L2CrossChainMessenger;
    const networkConfigAddress = networkConfig.NetworkConfig;
    
    console.log(`Using L1 bridge address: ${l1BridgeAddress}`);
    console.log(`Using L2 bridge address: ${l2BridgeAddress}`);
    console.log(`Using L2 CrossChainMessenger address: ${l2CrossChainMessengerAddress}`);
    console.log(`Using NetworkConfig address: ${networkConfigAddress}`);

    // Check if L2BridgeAddress exists
    if (!l2BridgeAddress) {
        throw new Error("L2Bridge not found in network config. Check that the bridge is deployed on L2.");
    }
    
    if (!l2CrossChainMessengerAddress) {
        throw new Error("L2CrossChainMessenger not found in network config. Check that it is deployed on L2.");
    }
    
    if (!networkConfigAddress) {
        throw new Error("NetworkConfig address not found in network config.");
    }

    // Step 1: Set remote bridge on L1 TenBridge
    console.log("Attaching L1 bridge to TenBridge contract...");
    const bridgeContract = (await hre.ethers.getContractFactory('TenBridge')).attach(l1BridgeAddress);
    
    console.log("Creating setRemoteBridge transaction...");
    const setRemoteBridgeTx = await bridgeContract.getFunction("setRemoteBridge").populateTransaction(l2BridgeAddress);

    console.log("Sending transaction to set remote bridge...");
    try {
        const receipt = await hre.companionNetworks.layer1.deployments.rawTx({
            from: deployer,
            to: l1BridgeAddress,
            data: setRemoteBridgeTx.data,
            log: true,
            waitConfirmations: 1,
        });

        if (receipt.status !== 1) {
            throw new Error("setRemoteBridge transaction failed");
        }
        console.log("Remote bridge set successfully");
    } catch (err: any) {
        const msg = err?.error?.message || err?.message || String(err);
        if (msg.includes("Remote bridge address already set")) {
            console.log("Remote bridge already set; skipping transaction.");
        } else {
            throw err;
        }
    }

    // Step 2: Set L2 Bridge address on L1 NetworkConfig
    console.log("\nSetting L2 Bridge address on L1 NetworkConfig...");
    const networkConfigContract = (await hre.ethers.getContractFactory('NetworkConfig')).attach(networkConfigAddress);
    
    try {
        const setL2BridgeTx = await networkConfigContract.getFunction("setL2BridgeAddress").populateTransaction(l2BridgeAddress);
        
        const receipt = await hre.companionNetworks.layer1.deployments.rawTx({
            from: deployer,
            to: networkConfigAddress,
            data: setL2BridgeTx.data,
            log: true,
            waitConfirmations: 1,
        });

        if (receipt.status !== 1) {
            throw new Error("setL2BridgeAddress transaction failed");
        }
        console.log("L2 Bridge address set successfully on NetworkConfig");
    } catch (err: any) {
        const msg = err?.error?.message || err?.message || String(err);
        if (msg.includes("already set")) {
            console.log("L2 Bridge address already set on NetworkConfig; skipping.");
        } else {
            console.error("Failed to set L2 Bridge address:", err);
            throw err;
        }
    }

    // Step 3: Set L2 CrossChainMessenger address on L1 NetworkConfig
    console.log("\nSetting L2 CrossChainMessenger address on L1 NetworkConfig...");
    
    try {
        const setL2MessengerTx = await networkConfigContract.getFunction("setL2CrossChainMessengerAddress").populateTransaction(l2CrossChainMessengerAddress);
        
        const receipt = await hre.companionNetworks.layer1.deployments.rawTx({
            from: deployer,
            to: networkConfigAddress,
            data: setL2MessengerTx.data,
            log: true,
            waitConfirmations: 1,
        });

        if (receipt.status !== 1) {
            throw new Error("setL2CrossChainMessengerAddress transaction failed");
        }
        console.log("L2 CrossChainMessenger address set successfully on NetworkConfig");
    } catch (err: any) {
        const msg = err?.error?.message || err?.message || String(err);
        if (msg.includes("already set")) {
            console.log("L2 CrossChainMessenger address already set on NetworkConfig; skipping.");
        } else {
            console.error("Failed to set L2 CrossChainMessenger address:", err);
            throw err;
        }
    }
    
    console.log("\n=== Bridge deployment and configuration completed successfully ===");
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];