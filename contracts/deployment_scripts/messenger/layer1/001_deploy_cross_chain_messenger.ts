import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import 'process';

/*
    This script deploys the CrossChainMessenger contract on the L1. 
    It depends on knowing the address of the message bus from the management contract predeployment.
    This is passed using the environment variables.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const deployments = hre.companionNetworks.layer1.deployments;

    const { deployer } = await hre.companionNetworks.layer1.getNamedAccounts();

    // Use the contract addresses from the management contract deployment.
    const mgmtContractAddress = process.env.MGMT_CONTRACT_ADDRESS!!
    const messageBusAddress : string = process.env.MESSAGE_BUS_ADDRESS!!
    console.log(`Management Contract address ${mgmtContractAddress}`);
    console.log(`Message Bus address ${messageBusAddress}`);

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


    console.log(`Management Contract address ${mgmtContractAddress}`);
    // get management contract and write the cross chain messenger address to it
    const mgmtContract = (await hre.ethers.getContractFactory('ManagementContract')).attach(mgmtContractAddress)
    const tx = await  mgmtContract.getFunction("SetImportantContractAddress").populateTransaction("L1CrossChainMessenger", crossChainDeployment.address);
    const receipt = await hre.companionNetworks.layer1.deployments.rawTx({
        from: deployer,
        to: mgmtContractAddress,
        data: tx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt.events?.length === 0) {
        console.log(`Failed to set L1CrossChainMessenger=${crossChainDeployment.address} on management contract.`);
    } else {
        console.log(`L1CrossChainMessenger=${crossChainDeployment.address}`);
    }
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
func.dependencies = ['ManagementContract', 'HPERC20', 'GasPrefunding']; //TODO: Remove HPERC20, this is only to have matching addresses.
