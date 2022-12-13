import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { ManagementContract } from '../typechain-types/contracts/management';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    const mgmtDeployment = await deployments.get("ManagementContract");

    const contract : ManagementContract = await hre.ethers.getContract("ManagementContract");
    const ManagementContract : ManagementContract = await contract.attach(mgmtDeployment.address);
    const busAddress = await ManagementContract.messageBus();

    await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        args: [ busAddress ],
        log: true,
    });
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
func.dependencies = ['ManagementContract'];
