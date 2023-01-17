import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();
    const contractArtifact = await hre.artifacts.readArtifact("ManagementContract");

    await deployments.deploy('ManagementContract', {
        from: deployer,
        contract: contractArtifact,
        args: [],
        log: true,
    });
};

export default func;
func.tags = ['ManagementContract', 'ManagementContract_deploy'];
