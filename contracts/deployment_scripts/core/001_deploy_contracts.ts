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

    // Deploy MerkleTreeMessageBus first
    const merkleMessageBusDeployment = await deployments.deploy('MerkleTreeMessageBus', {
        from: deployer,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [deployer, deployer] // initialOwner and withdrawalManager
                }
            }
        },
        log: true,
    });

    // Deploy CrossChain with MessageBus address
    const crossChainDeployment = await deployments.deploy('CrossChain', {
        from: deployer,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [deployer, merkleMessageBusDeployment.address]
                }
            }
        },
        log: true,
    });

    const merkleMessageBusAddress = await deployments.read('CrossChain', 'merkleMessageBus');

    const networkEnclaveRegistryDeployment = await deployments.deploy('NetworkEnclaveRegistry', {
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

    const daRegistryDeployment = await deployments.deploy('DataAvailabilityRegistry', {
        from: deployer,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
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
                        dataAvailabilityRegistry: daRegistryDeployment.address
                    }, deployer]
                }
            }
        },
        log: true,
    });

    console.log(`NetworkConfig= ${networkConfigDeployment.address}`); // line[0] in docker container
    console.log(`CrossChain= ${crossChainDeployment.address}`);
    console.log(`MerkleMessageBus= ${merkleMessageBusAddress}`);
    console.log(`NetworkEnclaveRegistry= ${networkEnclaveRegistryDeployment.address}`);
    console.log(`DataAvailabilityRegistry= ${daRegistryDeployment.address}`);
    console.log(`L1Start= ${networkConfigDeployment.receipt!!.blockHash}`);

    // Grant the DataAvailabilityRegistry the stateRootManager permission on MerkleMessageBus so it can publish rollups
    const merkleMessageBusContract = await hre.ethers.getContractAt('MerkleTreeMessageBus', merkleMessageBusAddress);
    const tx = await merkleMessageBusContract.addStateRootManager(daRegistryDeployment.address);
    const receipt = await tx.wait();
    if (receipt!.status !== 1) {
        throw new Error('Failed to add DataAvailabilityRegistry as stateRootManager to MerkleMessageBus');
    }
};

export default func;
func.tags = ['NetworkConfig', 'NetworkConfig_deploy'];
// No dependencies