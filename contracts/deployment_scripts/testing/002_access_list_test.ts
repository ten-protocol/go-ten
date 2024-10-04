import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    return;
    const l2Network = hre; 
    const {deployer} = await hre.getNamedAccounts();

    const gcb = await l2Network.deployments.deploy("GasConsumerBalance", {
        from: deployer,
        log: true
    })
    
    const signer = await hre.ethers.getSigner(deployer);

    const gasConsumerBalance = await hre.ethers.getContractAt("GasConsumerBalance", gcb.address)
    const gasConsumer = await gasConsumerBalance.connect(signer)


    const tx = await gasConsumer.getFunction("resetOwner").populateTransaction(deployer);
    tx.accessList = [
        { 
            address: gcb.address,
            storageKeys: []
        },
    ];
    const resp = await signer.sendTransaction(tx);
    const receipt = await resp.wait();
    console.log(`Receipt.Status=${receipt.status}`);
};


export default func;
func.tags = ['GasDebug'];
