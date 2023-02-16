import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/*
    This deployment script instantiates the Management Contract and additionally reads and prints
    out the message bus address to be used in CI/CD!
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    // The deployer prefunded address to be used to deploy the management contract
    const {deployer} = await getNamedAccounts();
    // The compiled contract artifact.
    const contractArtifact = await hre.artifacts.readArtifact("ManagementContract");

    // Deploying the management contract
    await deployments.deploy('ManagementContract', {
        from: deployer,
        contract: contractArtifact,
        args: [],
        log: true,
    });
    const busAddress = await deployments.read('ManagementContract', 'messageBus');

    // This is required in CI/CD - look at testnet-deploy-contracts.sh for more information.
    // depends on grep -e MessageBusAddress and a positional cut of the address
    console.log(`MessageBusAddress= ${busAddress}`);
};

export default func;
func.tags = ['ManagementContract', 'ManagementContract_deploy'];
// No dependencies