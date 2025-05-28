// Requires: npm install axios
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import {MessageBus, TenBridge} from "../../../typechain-types";
import axios from 'axios';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();

    const rpcUrl = 'http://erpc.sepolia-testnet.obscu.ro:80';
    const rpcPayload = {
        jsonrpc: '2.0',
        method: 'ten_config',
        params: [],
        id: 1
    };
    const response = await axios.post(rpcUrl, rpcPayload, {
        headers: { 'Content-Type': 'application/json' }
    });

    // For 1.1 use this. Hardcoded to gnosis safe:
    // const messageBusAddress = response.data.result.L1MessageBus;
    // const messageBusContract = (await hre.ethers.getContractFactory('MessageBus')).attach(messageBusAddress) as MessageBus;
    // const tx = await messageBusContract.retrieveAllFunds("0xeA052c9635F1647A8a199c2315B9A66ce7d1e2a7");

    // After >1.2 has been released:
    const bridgeContractAddress = response.data.result.L1Bridge;
    const bridgeContract = (await hre.ethers.getContractFactory('TenBridge')).attach(bridgeContractAddress) as TenBridge;
    const tx = await bridgeContract.retrieveAllFunds();

    const receipt = await tx.wait();
    if (receipt && receipt.status === 1) {
        console.log("Successfully recovered funds from the bridge.");
    } else {
        console.log("Recovery transaction failed");
    }
};

export default func;
// No dependencies