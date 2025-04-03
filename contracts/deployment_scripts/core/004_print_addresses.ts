import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/*
    This deployment script instantiates the network contracts and stores them in the deployed NetworkConfig contract.
*/
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {
        deployments,
    } = hre;
   
    const networkConfigDeployment = await deployments.get("NetworkConfig");
    const crossChainDeployment = await deployments.get("CrossChainMessenger");
    const merkleMessageBusAddress = await deployments.read("NetworkConfig", {}, "messageBusContractAddress");
    const networkEnclaveRegistryDeployment = await deployments.get("NetworkEnclaveRegistry");
    const daRegistryDeployment = await deployments.get("DataAvailabilityRegistry");
    const bridgeDeployment = await deployments.get("TenBridge");

    console.log(`NetworkConfig= ${networkConfigDeployment.address}`); // line[0] in docker container
    console.log(`CrossChain= ${crossChainDeployment.address}`);
    console.log(`MerkleMessageBus= ${merkleMessageBusAddress}`);
    console.log(`NetworkEnclaveRegistry= ${networkEnclaveRegistryDeployment.address}`);
    console.log(`DataAvailabilityRegistry= ${daRegistryDeployment.address}`);
    console.log(`L1Start= ${networkConfigDeployment.receipt!!.blockHash}`);
    console.log(`L1Bridge= ${bridgeDeployment.address}`);
};

export default func;
func.tags = ['Address Printer', 'AddressPrinter_deploy'];
func.dependencies = ['NetworkConfig', 'TenBridge'];