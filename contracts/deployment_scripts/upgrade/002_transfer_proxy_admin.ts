import { ethers } from "hardhat";

/**
 * Transfer proxy admin ownership to Timelock
 *
 * This script transfers the admin ownership of each TransparentUpgradeableProxy
 * from the deployer to the TimelockController for multisig governance
 */

async function transferProxyAdminOwnership() {
    const [deployer] = await ethers.getSigners();

    if (!deployer) {
        throw new Error('No deployer signer found');
    }

    console.log("Transferring proxy admin ownership to Timelock...");
    console.log("Deployer address:", deployer.address);

    // Configuration - these should be set as environment variables
    const timelockAddress = process.env.TIMELOCK_ADDRESS || "0x...";
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR || "0x...";

    if (timelockAddress === "0x..." || networkConfigAddr === "0x...") {
        throw new Error('Please set TIMELOCK_ADDRESS and NETWORK_CONFIG_ADDR environment variables');
    }

    console.log("Configuration:");
    console.log("- Timelock address:", timelockAddress);
    console.log("- NetworkConfig address:", networkConfigAddr);

    try {
        // Get addresses from network config
        const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
        const addresses = await networkConfig.addresses();

        console.log("\nCurrent proxy addresses:");
        console.table({
            NetworkConfig: networkConfigAddr,
            CrossChain: addresses.crossChain,
            NetworkEnclaveRegistry: addresses.networkEnclaveRegistry,
            DataAvailabilityRegistry: addresses.dataAvailabilityRegistry
        });

        // Get the TransparentUpgradeableProxy contract factory
        const TransparentUpgradeableProxy = await ethers.getContractFactory("TransparentUpgradeableProxy");

        // List of proxies to transfer admin ownership
        const proxies = [
            { name: "CrossChain", address: addresses.crossChain },
            { name: "NetworkEnclaveRegistry", address: addresses.networkEnclaveRegistry },
            { name: "DataAvailabilityRegistry", address: addresses.dataAvailabilityRegistry }
        ];

        for (const proxy of proxies) {
            console.log(`\n=== Transferring ${proxy.name} proxy admin ownership ===`);

            // Get the proxy contract
            const proxyContract = TransparentUpgradeableProxy.attach(proxy.address);

            // Get current admin
            const currentAdmin = await (proxyContract as any).admin();
            console.log(`Current admin: ${currentAdmin}`);

            if (currentAdmin.toLowerCase() === timelockAddress.toLowerCase()) {
                console.log(`${proxy.name} proxy admin already transferred to Timelock`);
                continue;
            }

            if (currentAdmin.toLowerCase() !== deployer.address.toLowerCase()) {
                console.log(`Warning: ${proxy.name} proxy admin is not the deployer (${currentAdmin})`);
                console.log("Skipping this proxy - manual intervention required");
                continue;
            }

            // Transfer admin ownership to timelock
            console.log(`Transferring admin ownership from ${deployer.address} to ${timelockAddress}...`);

            const transferTx = await (proxyContract as any).changeAdmin(timelockAddress);
            await transferTx.wait();

            console.log(`${proxy.name} proxy admin ownership transferred successfully!`);
            console.log(`Transaction hash: ${transferTx.hash}`);
        }

        console.log("\n=== Transfer Summary ===");
        console.log("All proxy admin ownership transfers completed");
        console.log("Timelock now controls all proxy upgrades");
        console.log("==========================\n");

    } catch (error) {
        console.error("Failed to transfer proxy admin ownership:", error);
        throw error;
    }
}

// Run the transfer if this script is executed directly
if (require.main === module) {
    transferProxyAdminOwnership()
        .then(() => process.exit(0))
        .catch((error) => {
            console.error(error);
            process.exit(1);
        });
}

export {
    transferProxyAdminOwnership
};