import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import 'process';

/*
    This script deploys the CrossChainMessenger contract on the L1. 
    It depends on knowing the address of the message bus from the management contract predeployment.
    This is passed using the environment variables.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const deployments = hre.deployments;

    const { deployer } = await hre.getNamedAccounts();

    const networkConfigDeployment = await deployments.get("NetworkConfig");
    const networkConfigAddress = networkConfigDeployment.address;

    const messageBusAddress = await deployments.read("NetworkConfig", {}, "messageBusContractAddress");

    // Setup the cross chain messenger and point it to the message bus from the management contract to be used for validation
    const crossChainDeployment = await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        log: true,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [ messageBusAddress ]
                }
            }
        }
    });

    console.log("Setting L1 Cross chain messenger")
    // get network config contract and write the cross chain messenger address to it
    const networkConfigContract = (await hre.ethers.getContractFactory('NetworkConfig')).attach(networkConfigAddress)
    const tx = await  networkConfigContract.getFunction("setL1CrossChainMessengerAddress").populateTransaction(crossChainDeployment.address);
    const receipt = await hre.deployments.rawTx({
        from: deployer,
        to: networkConfigAddress,
        data: tx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt.events?.length === 0) {
        console.log(`Failed to set L1CrossChainMessenger=${crossChainDeployment.address} on network config contract.`);
    } else {
        console.log(`L1CrossChainMessenger=${crossChainDeployment.address}`);
    }
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
func.dependencies = ['NetworkConfig'];
