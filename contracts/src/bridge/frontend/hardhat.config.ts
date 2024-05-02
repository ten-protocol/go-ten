/** @type import('hardhat/config').HardhatUserConfig */

require("dotenv").config();
require("@nomicfoundation/hardhat-toolbox");
require("@nomiclabs/hardhat-ethers");

const { PRIVATE_KEY, USER_KEY } = process.env;

task("deploy", "Deploys the ObscuroBridge contract").setAction(
  async (taskArgs, hre) => {
    const ObscuroBridge = await hre.ethers.getContractFactory("L1Bridge");
    const obscuroBridge = await ObscuroBridge.deploy();
    await ObscuroBridge.deployed();

    console.log("ObscuroBridge deployed to:", obscuroBridge.address);
  }
);

task("deploy", "Deploys the EthereumBridge contract").setAction(
  async (taskArgs, hre) => {
    const EthereumBridge = await hre.ethers.getContractFactory("L1Bridge");
    const ethereumBridge = await EthereumBridge.deploy();
    await EthereumBridge.deployed();

    console.log("EthereumBridge deployed to:", ethereumBridge.address);
  }
);

module.exports = {
  solidity: "0.8.19",
  paths: {
    sources: "./contracts",
    artifacts: "./src/assets/contract/artifacts",
  },
  networks: {
    ten: {
      chainId: 443,
      url: `https://testnet.ten.xyz/v1/${USER_KEY}`,
      gasPrice: 2000000000,
      accounts: [`0x${PRIVATE_KEY}`],
    },
  },
};
