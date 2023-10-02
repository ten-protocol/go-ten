import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();

    const mgmtContractAddress = process.env.MGMT_CONTRACT_ADDRESS!!
    // todo: if we want to support this we need to add the payAcc address param to the RetrieveAllBridgeFunds solidity defn
    const addressToPay = process.env.ACC_TO_PAY!!

    const mgmtContract = (await hre.ethers.getContractFactory('ManagementContract')).attach(mgmtContractAddress)
    const tx = await mgmtContract.RetrieveAllBridgeFunds();
    const receipt = await tx.wait();

    // Check the receipt for success, logs, etc.
    if (receipt.status === 1) {
      console.log("Successfully recovered funds from the bridge.");
    } else {
      console.log("Recovery transaction failed");
    }
};

export default func;
// No dependencies