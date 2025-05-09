import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import {TenBridge} from "../../../typechain-types";
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();

    const bridgeContractAddress = process.env.BRIDGE_CONTRACT_ADDRESS!!
    const bridgeContract = (await hre.ethers.getContractFactory('TenBridge')).attach(bridgeContractAddress) as TenBridge;
    const tx = await bridgeContract.retrieveAllFunds();
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