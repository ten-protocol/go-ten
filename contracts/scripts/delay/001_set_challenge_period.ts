import { ethers } from "hardhat";
import {RollupContract} from "../../typechain-types";

const setChallengePeriod = async function (rollupContractAddress: string, challengPeriod: number) {
    const rollupContract = await ethers.getContractAt(
        "RollupContract",
        rollupContractAddress
    ) as RollupContract;


    console.log(`Setting challenge period to: ${challengPeriod}`);
    const tx = await rollupContract.setChallengePeriod(BigInt(challengPeriod));
    await tx.wait();
    console.log(`Successfully set challenge period to: ${challengPeriod}`);
    
    const rollupChallengePeriod = await rollupContract.getChallengePeriod();
    if (BigInt(challengPeriod) !== rollupChallengePeriod) {
        throw new Error(`Failed to set the challenge period to: ${challengPeriod}. Returned value is: ${rollupChallengePeriod}`);
    }
}

const crossChainContractAddress = process.env.ROLLUP_CONTRACT_ADDR;
const challengePeriod = process.env.L1_CHALLENGE_PERIOD ?
    Number(process.env.L1_CHALLENGE_PERIOD) : 0;


if (!crossChainContractAddress) {
    console.error("Missing required environment variables: ROLLUP_CONTRACT_ADDR");
    process.exit(1);
}

setChallengePeriod(crossChainContractAddress, challengePeriod)
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default setChallengePeriod;