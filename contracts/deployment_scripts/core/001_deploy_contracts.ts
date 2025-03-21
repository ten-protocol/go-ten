import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/*
    This deployment script instantiates the network contracts and stores them in the deployed NetworkConfig contract.
*/
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;
    const {deployer} = await getNamedAccounts();
    // Deploy CrossChain with MessageBus address
    const crossChainDeployment = await deployments.deploy('CrossChain', {
        from: deployer,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [deployer]
                }
            }
        },
        log: true,
    });

    const merkleMessageBusAddress = await deployments.read('CrossChain', 'merkleMessageBus');

    const networkEnclaveRegistryDeployment = await deployments.deploy('NetworkEnclaveRegistry', {
        from: deployer,
        execute: {
            init: {
                methodName: "initialize",
                args: [deployer]
            }
        },
        log: true,
    });

    const daRegistryDeployment = await deployments.deploy('DataAvailabilityRegistry', {
        from: deployer,
        execute: {
            init: {
                methodName: "initialize",
                args: [
                    merkleMessageBusAddress,
                    networkEnclaveRegistryDeployment.address,
                    deployer
                ]
            }
        },
        log: true,
    });

    // Then deploy NetworkConfig with all addresses
    const networkConfigDeployment = await deployments.deploy('NetworkConfig', {
        from: deployer,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [{
                        crossChain: crossChainDeployment.address,
                        messageBus: merkleMessageBusAddress,
                        networkEnclaveRegistry: networkEnclaveRegistryDeployment.address,
                        rollupContract: daRegistryDeployment.address
                    }, deployer]
                }
            }
        },
        log: true,
    });

    console.log(`NetworkConfig= ${networkConfigDeployment.address}`);
    console.log(`CrossChain= ${crossChainDeployment.address}`);
    console.log(`MerkleMessageBus= ${merkleMessageBusAddress}`);
    console.log(`NetworkEnclaveRegistry= ${networkEnclaveRegistryDeployment.address}`);
    console.log(`DataAvailabilityRegistry= ${daRegistryDeployment.address}`);
    console.log(`L1Start= ${networkConfigDeployment.receipt!!.blockHash}`);
};

export default func;
func.tags = ['NetworkConfig', 'NetworkConfig_deploy'];
// No dependencies