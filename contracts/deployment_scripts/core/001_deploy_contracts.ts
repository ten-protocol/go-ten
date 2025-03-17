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
                    args: []
                }
            }
        },
        log: true,
    });

    const merkleMessageBusAddress = await deployments.read('CrossChain', 'merkleMessageBus');
    console.log(`MerkleMessageBus deployed to: ${merkleMessageBusAddress}`);

    const networkEnclaveRegistryDeployment = await deployments.deploy('NetworkEnclaveRegistry', {
        from: deployer,
        execute: {
            init: {
                methodName: "initialize",
                args: []
            }
        },
        log: true,
    });

    const rollupContractDeployment = await deployments.deploy('RollupContract', {
        from: deployer,
        execute: {
            init: {
                methodName: "initialize",
                args: [
                    merkleMessageBusAddress,
                    networkEnclaveRegistryDeployment.address
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
                        rollupContract: rollupContractDeployment.address
                    }, deployer]
                }
            }
        },
        log: true,
    });

    console.log(`NetworkConfig deployed to: ${networkConfigDeployment.address}`);
    console.log(`MerkleMessageBus deployed to: ${merkleMessageBusAddress}`);
    console.log(`CrossChain deployed to: ${crossChainDeployment.address}`);
    console.log(`NetworkEnclaveRegistry deployed to: ${networkEnclaveRegistryDeployment.address}`);
    console.log(`RollupContract deployed to: ${rollupContractDeployment.address}`);
};

export default func;
func.tags = ['NetworkConfig', 'NetworkConfig_deploy'];
// No dependencies