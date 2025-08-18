import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/*
    This deployment script verifies all deployed contracts on Etherscan
    Only runs if ETHERSCAN_API_KEY is provided
*/
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    // Skip verification if no API key is provided
    if (!process.env.ETHERSCAN_API_KEY) {
        console.log("Skipping contract verification - ETHERSCAN_API_KEY not provided");
        return;
    }

    const { deployments, getNamedAccounts } = hre;
    const { deployer } = await getNamedAccounts();

    console.log("Starting contract verification on Etherscan...");

    try {
        // Get all deployed contracts
        const networkConfigDeployment = await deployments.get('NetworkConfig');
        const crossChainDeployment = await deployments.get('CrossChain');
        const merkleMessageBusDeployment = await deployments.get('MerkleTreeMessageBus');
        const networkEnclaveRegistryDeployment = await deployments.get('NetworkEnclaveRegistry');
        const daRegistryDeployment = await deployments.get('DataAvailabilityRegistry');
        const feesDeployment = await deployments.get('Fees');

        // Get sequencer host address for verification
        const sequencerHostAddress = process.env.SEQUENCER_HOST_ADDRESS;
        if (!sequencerHostAddress) {
            console.error("SEQUENCER_HOST_ADDRESS environment variable is not set.");
            return;
        }

        // Verify Fees contract
        console.log("Verifying Fees contract...");
        await hre.run("verify:verify", {
            address: feesDeployment.address,
            constructorArguments: []
        });

        // Verify MerkleTreeMessageBus contract
        console.log("Verifying MerkleTreeMessageBus contract...");
        await hre.run("verify:verify", {
            address: merkleMessageBusDeployment.address,
            constructorArguments: []
        });

        // Verify CrossChain contract
        console.log("Verifying CrossChain contract...");
        await hre.run("verify:verify", {
            address: crossChainDeployment.address,
            constructorArguments: []
        });

        // Verify NetworkEnclaveRegistry contract
        console.log("Verifying NetworkEnclaveRegistry contract...");
        await hre.run("verify:verify", {
            address: networkEnclaveRegistryDeployment.address,
            constructorArguments: []
        });

        // Verify DataAvailabilityRegistry contract
        console.log("Verifying DataAvailabilityRegistry contract...");
        await hre.run("verify:verify", {
            address: daRegistryDeployment.address,
            constructorArguments: []
        });

        // Verify NetworkConfig contract
        console.log("Verifying NetworkConfig contract...");
        await hre.run("verify:verify", {
            address: networkConfigDeployment.address,
            constructorArguments: []
        });

        console.log("✅ All contracts verified successfully on Etherscan");

    } catch (error) {
        console.error("❌ Contract verification failed:", error);
        // Don't fail the deployment if verification fails
        console.log("Deployment continues despite verification failure");
    }
};

export default func;
func.tags = ['NetworkConfig', 'Verify'];
func.dependencies = ['NetworkConfig_deploy']; // Run after main deployment