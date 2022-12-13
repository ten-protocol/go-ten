import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    // TODO: Remove hardcoded L2 message bus address when properly exposed.
    const busAddress = hre.ethers.utils.getAddress("0x526c84529b2b8c11f57d93d3f5537aca3aecef9b")

    console.log(`Beginning deploy of cross chain messenger`);

    await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        args: [ busAddress ],
        log: true,
    });
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];