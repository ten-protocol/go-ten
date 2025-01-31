import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { network } from 'hardhat';

/* 
    This script sets the fee for the message bus to prevent spam.
*/
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 

    const l2Accounts = await l2Network.getNamedAccounts();

    const networkConfig = await l2Network.network.provider.request({
        method: "net_config",
    });    

    const signer = await l2Network.ethers.getSigner(l2Accounts.deployer);
    const fees = await l2Network.ethers.getContractAt(
        'Fees', 
        networkConfig["PublicSystemContracts"]["Fees"], 
        signer
    );

    const owner = await fees.owner();
    console.log(`Owner = ${owner}`);
    console.log(`Signer = ${l2Accounts.deployer}`);

    const tx = await fees.setMessageFee(32*10000);
    const receipt =await tx.wait();

    if (receipt.status != 1) {
        throw new Error("Failed to set message fee");
    }
    console.log(`Fee set at ${receipt.hash}`);
}
export default func;
func.tags = ['SetFees', 'SetFees_deploy'];
func.dependencies = ['ZenBase'];
