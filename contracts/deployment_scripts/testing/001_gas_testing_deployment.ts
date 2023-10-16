import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 
    const {deployer} = await hre.getNamedAccounts();

    const gcb = await l2Network.deployments.deploy("GasConsumerBalance", {
        from: deployer,
        log: true
    })
    
    
    const gasConsumerBalance = await hre.ethers.getContractAt("GasConsumerBalance", gcb.address)
    const gasEstimation = await gasConsumerBalance.estimateGas.get_balance({
        from: deployer,
        gasPrice: 2,
    });
    
    await hre.deployments.execute("GasConsumerBalance", {
        from: deployer,
        gasLimit: gasEstimation,
        log: true
    }, "get_balance");
};


export default func;
func.tags = ['GasDebug'];
