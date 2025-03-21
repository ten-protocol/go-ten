import { ethers } from "hardhat";
import {DataAvailabilityRegistry} from "../../typechain-types";

const setChallengePeriod = async function (daRegistryAddress: string, challengPeriod: number) {
    const daRegistryContract = await ethers.getContractAt(
        "DataAvailabilityRegistry",
        daRegistryAddress
    ) as DataAvailabilityRegistry;


    console.log(`Setting challenge period to: ${challengPeriod}`);
    const tx = await daRegistryContract.setChallengePeriod(BigInt(challengPeriod));
    await tx.wait();
    console.log(`Successfully set challenge period to: ${challengPeriod}`);

    const rollupChallengePeriod = await daRegistryContract.getChallengePeriod();
    if (BigInt(challengPeriod) !== rollupChallengePeriod) {
        throw new Error(`Failed to set the challenge period to: ${challengPeriod}. Returned value is: ${rollupChallengePeriod}`);
    }
}

const daRegistryAddress = process.env.DA_REGISTRY_ADDR;
const challengePeriod = process.env.L1_CHALLENGE_PERIOD ?
    Number(process.env.L1_CHALLENGE_PERIOD) : 0;


if (!daRegistryAddress) {
    console.error("Missing required environment variables: DA_REGISTRY_ADDR");
    process.exit(1);
}

setChallengePeriod(daRegistryAddress, challengePeriod)
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default setChallengePeriod;