import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { network } from 'hardhat';

/* 
    This deployment script deploys the TenBridgeTestnet contract.
    TenBridgeTestnet is a testnet-specific version of TenBridge that includes
    a function to recover native funds from the bridge contract.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {
        deployments, 
        getNamedAccounts
    } = hre;
    
    const { deployer } = await getNamedAccounts();
    
    // Get NetworkConfig address from environment variable or command line
    const networkConfigAddress = process.env.NETWORK_CONFIG_ADDR;
    if (!networkConfigAddress) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable or network config address parameter is required');
    }

    console.log(`Using NetworkConfig address: ${networkConfigAddress}`);

    // Get the NetworkConfig contract and retrieve addresses
    const networkConfigContract = await hre.ethers.getContractAt('NetworkConfig', networkConfigAddress);
    const addresses = await networkConfigContract.addresses();
    
    console.log(`Deploying TenBridgeTestnet with:`);
    console.log(`  Deployer: ${deployer}`);
    console.log(`  Messenger: ${addresses.l1CrossChainMessenger}`);
    console.log(`  NetworkConfig: ${networkConfigAddress}`);

    // Deploy the TenBridgeTestnet contract
    const tenBridgeTestnetDeployment = await deployments.deploy(
      "TenBridgeTestnet",
      {
        from: deployer,
        log: true,
        proxy: {
          proxyContract: "OpenZeppelinTransparentProxy",
          execute: {
            init: {
              methodName: "initialize",
              args: [addresses.l1CrossChainMessenger, deployer],
            },
          },
        },
      }
    )

    console.log(`TenBridgeTestnet deployed at: ${tenBridgeTestnetDeployment.address}`);

    console.log('Updating NetworkConfig with new bridge address...');
    await deployments.rawTx({
        from: deployer,
        to: networkConfigAddress,
        data: recordL1AddressTx.data,
        log: true,
        waitConfirmations: 1,
    });
    
    console.log(`L1BridgeTestnetAddress=${tenBridgeTestnetDeployment.address}`);
    
    // Verify the deployment
    const deployedContract = await hre.ethers.getContractAt('TenBridgeTestnet', tenBridgeTestnetDeployment.address);
    const remoteBridgeAddress = await deployedContract.remoteBridgeAddress();
    const hasAdminRole = await deployedContract.hasRole(await deployedContract.ADMIN_ROLE(), deployer);
    
    console.log('\n=== Deployment Verification ===');
    console.log(`Contract Address: ${tenBridgeTestnetDeployment.address}`);
    console.log(`Remote Bridge Address: ${remoteBridgeAddress}`);
    console.log(`Deployer has Admin Role: ${hasAdminRole}`);
    console.log('TenBridgeTestnet deployment completed successfully!');
};

export default func;
func.tags = ['TenBridgeTestnet', 'TenBridgeTestnet_deploy'];
func.dependencies = ['CrossChainMessenger'];
