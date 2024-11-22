import { ethers } from "hardhat";
import { ManagementContract } from "../../../../typechain-types";

async function main() {
    const mgmtContractAddr = process.env.MGMT_CONTRACT_ADDRESS;
    if (!mgmtContractAddr) {
        throw new Error("MGMT_CONTRACT_ADDRESS not set");
    }

    const enclaveIDs = process.env.ENCLAVE_IDS;
    if (!enclaveIDs) {
        throw new Error("ENCLAVE_IDS not set");
    }

    // Get the Management Contract
    const managementContract = await ethers.getContractAt(
        "ManagementContract",
        mgmtContractAddr
    ) as ManagementContract;

    // Parse the comma-separated list of enclave IDs
    const enclaveAddresses = enclaveIDs.split(",");

    // Grant sequencer status to each enclave
    for (const enclaveAddr of enclaveAddresses) {
        console.log(`Granting sequencer status to enclave: ${enclaveAddr}`);
        const tx = await managementContract.GrantSequencerEnclave(enclaveAddr);
        await tx.wait();
        console.log(`Successfully granted sequencer status to: ${enclaveAddr}`);

        // Verify the sequencer status
        const isSequencer = await managementContract.IsSequencerEnclave(enclaveAddr);
        if (!isSequencer) {
            throw new Error(`Failed to verify sequencer status for enclave: ${enclaveAddr}. IsSequencerEnclave returned false`);
        }
        console.log(`Verified sequencer status for enclave: ${enclaveAddr}`);
    }

    // Final verification of all enclaves
    console.log("\nFinal verification of all sequencer permissions:");
    for (const enclaveAddr of enclaveAddresses) {
        const isSequencer = await managementContract.IsSequencerEnclave(enclaveAddr);
        console.log(`Enclave ${enclaveAddr}: IsSequencerEnclave = ${isSequencer}`);
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });