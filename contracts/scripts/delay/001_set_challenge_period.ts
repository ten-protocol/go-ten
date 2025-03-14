import { ethers } from "hardhat";
import {CrossChain} from "../../typechain-types";

const setChallengePeriod = async function (crossChainContractAddress: string, challengPeriod: number) {
    const crossChainContract = await ethers.getContractAt(
        "CrossChain",
        crossChainContractAddress
    ) as CrossChain;


    console.log(`Setting challenge period to: ${challengPeriod}`);
    const tx = await crossChainContract.setChallengePeriod(BigInt(challengPeriod));
    await tx.wait();
    console.log(`Successfully set challenge period to: ${challengPeriod}`);
    
    const mgmtContractChallengePeriod = await crossChainContract.getChallengePeriod();
    if (BigInt(challengPeriod) !== mgmtContractChallengePeriod) {
        throw new Error(`Failed to set the challenge period to: ${challengPeriod}. Returned value is: ${mgmtContractChallengePeriod}`);
    }
}

const crossChainContractAddress = process.env.CROSS_CHAIN_ADDRESS;
const challengePeriod = process.env.L1_CHALLENGE_PERIOD ?
    Number(process.env.L1_CHALLENGE_PERIOD) : 0;


if (!crossChainContractAddress) {
    console.error("Missing required environment variables: CROSS_CHAIN_ADDRESS");
    process.exit(1);
}

setChallengePeriod(crossChainContractAddress, challengePeriod)
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default setChallengePeriod;