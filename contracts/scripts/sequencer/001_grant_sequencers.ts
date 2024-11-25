import { ethers } from "hardhat";
import { ManagementContract } from "../../typechain-types";

const grantSequencerStatus = async function (mgmtContractAddr: string, enclaveIDsStr: string) {
    const managementContract = await ethers.getContractAt(
        "ManagementContract",
        mgmtContractAddr
    ) as ManagementContract;

    const enclaveAddresses = enclaveIDsStr.split(",");

    for (const enclaveAddr of enclaveAddresses) {
        console.log(`Granting sequencer status to enclave: ${enclaveAddr}`);
        const tx = await managementContract.GrantSequencerEnclave(enclaveAddr);
        await tx.wait();
        console.log(`Successfully granted sequencer status to: ${enclaveAddr}`);

        // check they've been added
        const isSequencer = await managementContract.IsSequencerEnclave(enclaveAddr);
        if (!isSequencer) {
            throw new Error(`Failed to verify sequencer status for enclave: ${enclaveAddr}. IsSequencerEnclave returned false`);
        }
        console.log(`Verified sequencer status for enclave: ${enclaveAddr}`);
    }

    console.log("\nFinal verification of all sequencer permissions:");
    for (const enclaveAddr of enclaveAddresses) {
        const isSequencer = await managementContract.IsSequencerEnclave(enclaveAddr);
        console.log(`Enclave ${enclaveAddr}: IsSequencerEnclave = ${isSequencer}`);
    }
};

const args = process.argv.slice(2);
if (args.length !== 2) {
    throw new Error("Required arguments: <mgmtContractAddr> <enclaveIDs>");
}

const [mgmtContractAddr, enclaveIDs] = args as [string, string];
grantSequencerStatus(mgmtContractAddr, enclaveIDs)
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default grantSequencerStatus;