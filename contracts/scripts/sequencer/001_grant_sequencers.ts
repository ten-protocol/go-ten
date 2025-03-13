import { ethers } from "hardhat";
import {NetworkEnclaveRegistry} from "../../typechain-types";

const grantSequencerStatus = async function (enclaveRegistryAddress: string, enclaveIDsStr: string) {
    const enclaveRegistryContract = await ethers.getContractAt(
        "NetworkEnclaveRegistry",
        enclaveRegistryAddress
    ) as NetworkEnclaveRegistry;

    const enclaveAddresses = enclaveIDsStr.split(",");

    for (const enclaveAddr of enclaveAddresses) {
        console.log(`Granting sequencer status to enclave: ${enclaveAddr}`);
        const tx = await enclaveRegistryContract.grantSequencerEnclave(enclaveAddr);
        await tx.wait();
        console.log(`Successfully granted sequencer status to: ${enclaveAddr}`);

        // check they've been added
        const isSequencer = await enclaveRegistryContract.isSequencer(enclaveAddr);
        if (!isSequencer) {
            throw new Error(`Failed to verify sequencer status for enclave: ${enclaveAddr}. IsSequencerEnclave returned false`);
        }
        console.log(`Verified sequencer status for enclave: ${enclaveAddr}`);
    }

    console.log("\nFinal verification of all sequencer permissions:");
    for (const enclaveAddr of enclaveAddresses) {
        const isSequencer = await enclaveRegistryContract.isSequencer(enclaveAddr);
        console.log(`Enclave ${enclaveAddr}: IsSequencerEnclave = ${isSequencer}`);
    }
};

const enclaveRegistryAddress = process.env.ENCLAVE_REGISTRY_ADDRESS;
const enclaveIDs = process.env.ENCLAVE_IDS;

if (!enclaveRegistryAddress || !enclaveIDs) {
    console.error("Missing required environment variables: ENCLAVE_REGISTRY_ADDRESS and ENCLAVE_IDS.");
    process.exit(1);
}

grantSequencerStatus(enclaveRegistryAddress, enclaveIDs)
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default grantSequencerStatus;