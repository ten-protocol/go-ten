import { ethers } from "hardhat";

const getNetCfg = async function () {

    const networkConfig: any = await ethers.provider.send("net_config");
    if (!networkConfig || !networkConfig.L2MessageBusAddress) {
        throw new Error("Failed to retrieve L2MessageBusAddress from network config");
    }
    console.log(JSON.stringify(networkConfig, null, 2));
};

getNetCfg()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default getNetCfg;