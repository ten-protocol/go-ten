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

    // recover funds from the new bridge contract
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